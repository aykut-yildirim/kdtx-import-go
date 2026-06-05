package amazon

import (
	"errors"
	"fmt"
	"math"

	"strings"
	"time"

	"myapp/internal/helpers"
	"myapp/internal/models"
	"myapp/internal/services"

	"github.com/go-playground/validator/v10"
)

type TransactionFileImporter struct{}

func (i TransactionFileImporter) FileLoad(
	ctx *models.Context,
) error {
	services.Logger().STATUS("- FileLoad -")
	data, err := helpers.GetFile(*ctx)
	if err != nil {
		return err
	}

	ctx.RawData = data
	return nil
}

func (i TransactionFileImporter) LoginControl(
	ctx *models.Context,
) error {
	return nil
}

func (i TransactionFileImporter) Fetch(
	ctx *models.Context,
) error {
	return nil
}

type SettlementTransactionV1 struct {
	TransactionType           string `validate:"required"`
	OrderID                   string
	TotalAmount               string
	SettlementStartDate       string
	SettlementEndDate         string
	Currency                  string
	DepositDate               string
	PostedDate                string
	ShipmentFeeAmount         string
	OtherFeeAmount            string
	PriceAmount               string
	ItemRelatedFeeAmount      string
	MiscFeeAmount             string
	PromotionAmount           string
	DirectPaymentAmount       string
	OtherAmount               string
	OrderItemCode             string
	PriceType                 string
	ItemRelatedFeeType        string
	OtherFeeReasonDescription string
	PromotionType             string
}

func mapToSettlementTransactionV1(row map[string]string) SettlementTransactionV1 {
	return SettlementTransactionV1{
		TransactionType:           row["transaction-type"],
		OrderID:                   row["order-id"],
		TotalAmount:               row["total-amount"],
		SettlementStartDate:       row["settlement-start-date"],
		SettlementEndDate:         row["settlement-end-date"],
		Currency:                  row["currency"],
		DepositDate:               row["deposit-date"],
		PostedDate:                row["posted-date"],
		ShipmentFeeAmount:         row["shipment-fee-amount"],
		OtherFeeAmount:            row["other-fee-amount"],
		PriceAmount:               row["price-amount"],
		ItemRelatedFeeAmount:      row["item-related-fee-amount"],
		MiscFeeAmount:             row["misc-fee-amount"],
		PromotionAmount:           row["promotion-amount"],
		DirectPaymentAmount:       row["direct-payment-amount"],
		OtherAmount:               row["other-amount"],
		OrderItemCode:             row["order-item-code"],
		PriceType:                 row["price-type"],
		ItemRelatedFeeType:        row["item-related-fee-type"],
		OtherFeeReasonDescription: row["other-fee-reason-description"],
		PromotionType:             row["promotion-type"],
	}
}

func toColumnName(fieldName string) string {
	switch fieldName {
	case "TransactionType":
		return "transaction-type"
	case "OrderID":
		return "order-id"
	case "TotalAmount":
		return "total-amount"
	case "SettlementStartDate":
		return "settlement-start-date"
	case "SettlementEndDate":
		return "settlement-end-date"
	case "Currency":
		return "currency"
	case "DepositDate":
		return "deposit-date"
	case "PostedDate":
		return "posted-date"
	case "ShipmentFeeAmount":
		return "shipment-fee-amount"
	case "OtherFeeAmount":
		return "other-fee-amount"
	case "PriceAmount":
		return "price-amount"
	case "ItemRelatedFeeAmount":
		return "item-related-fee-amount"
	case "MiscFeeAmount":
		return "misc-fee-amount"
	case "PromotionAmount":
		return "promotion-amount"
	case "DirectPaymentAmount":
		return "direct-payment-amount"
	case "OtherAmount":
		return "other-amount"
	case "OrderItemCode":
		return "order-item-code"
	case "PriceType":
		return "price-type"
	case "ItemRelatedFeeType":
		return "item-related-fee-type"
	case "OtherFeeReasonDescription":
		return "other-fee-reason-description"
	case "PromotionType":
		return "promotion-type"
	default:
		return strings.ToLower(fieldName)
	}
}

func (i TransactionFileImporter) Map(
	ctx *models.Context,
) (interface{}, error) {
	dfDict, err := helpers.ReadTableFromBytes(ctx.RawData, "\t")
	if err != nil {
		return nil, fmt.Errorf("failed to read table: %w", err)
	}
	_transactions_dict := make(map[string][]map[string]string)
	validate := validator.New()

	currency := "EUR"

	push := func(transaction models.Transaction, details []models.TransactionDetail) {
		ctx.Data = append(ctx.Data, models.DataAll{
			Transaction:        &transaction,
			TransactionDetails: details,
		})
	}

	for idx, row := range dfDict {
		if idx == 0 {
			row["transaction-type"] = "Amazon / Payout"
		}

		s := mapToSettlementTransactionV1(row)
		err = validate.Struct(s)
		if err != nil {
			var validationErrors validator.ValidationErrors
			if errors.As(err, &validationErrors) {
				for _, fieldErr := range validationErrors {
					fieldName := fieldErr.Field()
					columnName := toColumnName(fieldName)
					errorMsg := fmt.Sprintf("failed validation on tag '%s'", fieldErr.Tag())
					return nil, fmt.Errorf("Validation error at row %d (data row %d), column '%s': %s", idx+2, idx+1, columnName, errorMsg)
				}
			}
			return nil, fmt.Errorf("Validation error at row %d (data row %d): %w", idx+2, idx+1, err)
		}

		if (row["transaction-type"] == "Order" || row["transaction-type"] == "Refund") && row["order-id"] != "" {
			key := row["transaction-type"] + row["order-id"]
			_transactions_dict[key] = append(_transactions_dict[key], row)

		} else if row["total-amount"] != "" && row["settlement-start-date"] != "" && row["settlement-end-date"] != "" {
			currency = row["currency"]

			var tDate time.Time
			var dateErr error
			if row["deposit-date"] != "" {
				tDate, dateErr = helpers.ParseDatetime(row["deposit-date"])
			} else {
				tDate, dateErr = helpers.ParseDatetime(row["settlement-end-date"])
			}
			if dateErr != nil {
				tDate = time.Now()
			}

			startDate, _ := helpers.ParseDatetime(row["settlement-start-date"])
			endDate, _ := helpers.ParseDatetime(row["settlement-end-date"])
			desc := helpers.ToStringCombine(" / ",
				"Amazon Payout",
				startDate.Format("2006-01-02"),
				endDate.Format("2006-01-02"),
			)

			amount, _ := helpers.MoneyToFloat(row["total-amount"])
			extRef := map[string]interface{}{
				"settlement_no": row["transaction-type"],
			}

			_transaction := models.Transaction{
				Date:                  tDate,
				CurrencyCode:          currency,
				Description:           &desc,
				Amount:                amount,
				OpenAmount:            amount,
				ExternalReferenceData: extRef,
				RawData:               row,
			}

			pattern := "Amazon / Bank Amount"
			_details := []models.TransactionDetail{
				{
					Amount:          amount,
					Description:     &desc,
					IsMatchRequired: true,
					PatternTemplate: &pattern,
				},
			}
			push(_transaction, _details)

		} else {
			totalAmount := helpers.ToFloatTotal(
				row["total-amount"],
				row["shipment-fee-amount"],
				row["other-fee-amount"],
				row["price-amount"],
				row["item-related-fee-amount"],
				row["misc-fee-amount"],
				row["other-fee-amount"],
				row["promotion-amount"],
				row["direct-payment-amount"],
				row["other-amount"],
			)

			var tDate time.Time
			var dateErr error
			if row["posted-date"] != "" {
				tDate, dateErr = helpers.ParseDatetime(row["posted-date"])
			}
			if dateErr != nil || row["posted-date"] == "" {
				tDate = time.Now()
			}

			desc := fmt.Sprintf("Other / %s", row["transaction-type"])
			_transaction := models.Transaction{
				Date:         tDate,
				CurrencyCode: currency,
				Description:  &desc,
				Amount:       totalAmount,
				OpenAmount:   totalAmount,
				RawData:      row,
			}

			detailDesc := "Other"
			pattern := row["transaction-type"]
			_details := []models.TransactionDetail{
				{
					Amount:          totalAmount,
					Description:     &detailDesc,
					IsMatchRequired: false,
					PatternTemplate: &pattern,
				},
			}
			push(_transaction, _details)
		}
	}

	for _, details := range _transactions_dict {
		if len(details) == 0 {
			continue
		}
		isDebit := false
		refundStr := ""
		if details[0]["transaction-type"] == "Refund" {
			refundStr = " - Refund"
			isDebit = true
		}

		tDate, err := helpers.ParseDatetime(details[0]["posted-date"])
		if err != nil {
			tDate = time.Now()
		}

		desc := details[0]["transaction-type"]
		orderNo := details[0]["order-id"]
		transactionNo := details[0]["order-id"] + refundStr
		referenceNo := fmt.Sprintf("%s-%s", details[0]["order-id"], details[0]["transaction-type"])

		_transaction := models.Transaction{
			IsDebit:       isDebit,
			Date:          tDate,
			CurrencyCode:  currency,
			Amount:        0.0,
			OpenAmount:    0.0,
			Description:   &desc,
			RawData:       details,
			OrderNo:       &orderNo,
			TransactionNo: &transactionNo,
			ReferenceNo:   &referenceNo,
		}

		_details_dict := make(map[string][]map[string]string)
		for _, row := range details {
			itemCode := row["order-item-code"]
			_details_dict[itemCode] = append(_details_dict[itemCode], row)
		}

		var _details []models.TransactionDetail
		for _, itemData := range _details_dict {
			if len(itemData) == 0 {
				continue
			}

			dPricePrincipal := 0.0
			dPriceTax := 0.0

			dPriceDesc := helpers.ToStringCombine(" ", itemData[0]["order-id"], itemData[0]["order-item-code"])
			dPricePattern := helpers.ToStringCombine(" / ", itemData[0]["transaction-type"], "price-amount")

			_d_price := models.TransactionDetail{
				Amount:          0.0,
				Description:     &dPriceDesc,
				IsMatchRequired: true,
				PatternTemplate: &dPricePattern,
			}

			for _, row := range itemData {
				if row["price-amount"] != "" {
					val, _ := helpers.MoneyToFloat(row["price-amount"])
					if row["price-type"] == "Principal" {
						dPricePrincipal = val
					}
					if row["price-type"] == "Tax" {
						dPriceTax = val
					}

				} else if row["item-related-fee-amount"] != "" {
					val, _ := helpers.MoneyToFloat(row["item-related-fee-amount"])
					desc := "item-related-fee-amount"
					pattern := helpers.ToStringCombine(" / ",
						row["transaction-type"],
						"item-related-fee-amount",
						row["item-related-fee-type"],
					)

					_d := models.TransactionDetail{
						Amount:          val,
						Description:     &desc,
						IsMatchRequired: false,
						PatternTemplate: &pattern,
					}
					_details = append(_details, _d)
					_transaction.Amount += val

				} else if row["other-fee-amount"] != "" {
					val, _ := helpers.MoneyToFloat(row["other-fee-amount"])
					desc := "other-fee-amount"
					pattern := helpers.ToStringCombine(" / ",
						row["transaction-type"],
						"other-fee-amount",
						row["other-fee-reason-description"],
					)

					_d := models.TransactionDetail{
						Amount:          val,
						Description:     &desc,
						IsMatchRequired: false,
						PatternTemplate: &pattern,
					}
					_details = append(_details, _d)
					_transaction.Amount += val

				} else if row["promotion-amount"] != "" {
					val, _ := helpers.MoneyToFloat(row["promotion-amount"])
					desc := "promotion-amount"
					pattern := helpers.ToStringCombine(" / ",
						row["transaction-type"],
						"promotion-amount",
						row["promotion-type"],
					)

					_d := models.TransactionDetail{
						Amount:          val,
						Description:     &desc,
						IsMatchRequired: false,
						PatternTemplate: &pattern,
					}
					_details = append(_details, _d)
					_transaction.Amount += val

				} else if row["direct-payment-amount"] != "" {
					val, _ := helpers.MoneyToFloat(row["direct-payment-amount"])
					desc := "direct-payment-amount"
					pattern := helpers.ToStringCombine(" / ",
						row["transaction-type"],
						"direct-payment-amount",
						row["other-fee-reason-description"],
					)

					_d := models.TransactionDetail{
						Amount:          val,
						Description:     &desc,
						IsMatchRequired: false,
						PatternTemplate: &pattern,
					}
					_details = append(_details, _d)
					_transaction.Amount += val
				}
			}

			if dPriceTax != 0.0 || dPricePrincipal != 0.0 {
				_d_price.Amount = dPricePrincipal + dPriceTax
				if dPricePrincipal != 0.0 {
					_d_price.TaxRate = math.Round((dPriceTax / dPricePrincipal) * 100)
				}
				_transaction.Amount = _d_price.Amount
				_details = append(_details, _d_price)
			}
		}
		_transaction.OpenAmount = _transaction.Amount
		push(_transaction, _details)
	}
	ctx.Result = ctx.Data
	return ctx, nil
}
