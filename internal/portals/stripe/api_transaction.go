package amazon

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"myapp/internal/helpers"
	"myapp/internal/models"
	"myapp/internal/services"
)

const stripeBaseURL = "https://api.stripe.com/v1"

type TransactionApiImporter struct{}

// ---------------------------------------------------------------------------
// Importer interface
// ---------------------------------------------------------------------------

func (i TransactionApiImporter) FileLoad(
	ctx *models.Context,
) error {
	return nil
}

func (i TransactionApiImporter) LoginControl(
	ctx *models.Context,
) error {
	token, err := getToken(ctx)
	if err != nil {
		return err
	}

	resp, err := stripeGet(token, stripeBaseURL+"/customers", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("login_control failed (%d): %s", resp.StatusCode, string(body))
	}

	services.Logger().STATUS("LoginControl OK")
	return nil
}

func (i TransactionApiImporter) Fetch(
	ctx *models.Context,
) error {
	return nil
}

func (i TransactionApiImporter) Map(
	ctx *models.Context,
) error {

	token, err := getToken(ctx)
	if err != nil {
		return  err
	}

	// -----------------------------------------------------------------------
	// Fetch all payouts
	// -----------------------------------------------------------------------
	services.Logger().SYSTEM("Downloading Payouts")

	allPayouts, err := fetchAllPayouts(token, ctx.Task.FetchStartDate, ctx.Task.FetchEndDate)
	if err != nil {
		return err
	}

	var settlementList []models.Settlement

	for _, payout := range allPayouts {
		payoutID, _ := payout["id"].(string)

		payoutTxList, err := fetchAllTransactionsForPayout(
			token,
			fmt.Sprintf("&payout=%s", payoutID),
			"",
		)
		if err != nil {
			return err
		}

		services.Logger().STATUS(fmt.Sprintf(
			"%s Related to payout -> %d transaction_no found",
			payoutID, len(payoutTxList),
		))

		var txIDs []string
		for _, tx := range payoutTxList {
			if id, ok := tx["id"].(string); ok {
				txIDs = append(txIDs, id)
			}
		}

		amount := helpers.ToFloat64(payout["amount"]) / 100
		currency, _ := payout["currency"].(string)
		created := helpers.FromUnixToTime(payout["created"])
		description := "Stripe Settlement "

		settlementList = append(settlementList, models.Settlement{
			SettlementNo:    payoutID,
			Amount:          amount,
			CurrencyCode:    currency,
			StartDate:       created,
			EndDate:         created,
			Description:     &description,
			TransactionList: txIDs,
			RawData:         payout,
		})
	}

	ctx.Data = append(ctx.Data, models.DataAll{
		SettlementList: settlementList,
	})

	// -----------------------------------------------------------------------
	// Fetch all transactions in date range
	// -----------------------------------------------------------------------
	startFilter := ""
	endFilter := ""
	if ctx.Task.FetchStartDate != nil {
		startFilter = fmt.Sprintf("&created[gt]=%d", helpers.ToUnixTimestamp(*ctx.Task.FetchStartDate))
	}
	if ctx.Task.FetchEndDate != nil {
		endFilter = fmt.Sprintf("&created[lt]=%d", helpers.ToUnixTimestamp(*ctx.Task.FetchEndDate))
	}

	allTransactions, err := fetchAllTransactionsForPayout(token, "", startFilter+endFilter)
	if err != nil {
		return err
	}

	// Duplicate check: collect all IDs first
	allIDs := make([]string, 0, len(allTransactions))
	for _, tx := range allTransactions {
		if id, ok := tx["id"].(string); ok {
			allIDs = append(allIDs, id)
		}
	}

	// Store IDs for duplicate filtering (caller / pipeline can use DataAll.AllIDs)
	ctx.Data = append(ctx.Data, models.DataAll{
		AllIDs: allIDs,
	})

	services.Logger().STATUS(fmt.Sprintf(
		"Transactions count (total): %d",
		len(allTransactions),
	))

	// -----------------------------------------------------------------------
	// Map transactions
	// -----------------------------------------------------------------------
	for _, tx := range allTransactions {
		txID, _ := tx["id"].(string)
		currency, _ := tx["currency"].(string)
		amount := helpers.ToFloat64(tx["amount"]) / 100
		net := helpers.ToFloat64(tx["net"]) / 100
		description, _ := tx["description"].(string)
		txType, _ := tx["type"].(string)
		source, _ := tx["source"].(string)
		created := helpers.FromUnixToTime(tx["created"])

		transaction := models.Transaction{
			Date:          created,
			TransactionNo: &txID,
			ReferenceNo:   &txID,
			CurrencyCode:  currency,
			Amount:        amount,
			OpenAmount:    net,
			ExternalReferenceData: map[string]interface{}{
				"paymentId": source,
			},
			RawData: tx,
		}

		// Main detail (net amount)
		details := []models.TransactionDetail{
			{
				Amount:          net,
				IsMatchRequired: true,
				Description:     &description,
				PatternTemplate: &txType,
			},
		}

		// Fee details
		if feeDetails, ok := tx["fee_details"].([]interface{}); ok {
			for _, fd := range feeDetails {
				feeMap, ok := fd.(map[string]interface{})
				if !ok {
					continue
				}
				feeAmount := helpers.ToFloat64(feeMap["amount"]) / 100
				feeDesc, _ := feeMap["description"].(string)
				feeType, _ := feeMap["type"].(string)

				details = append(details, models.TransactionDetail{
					Amount:          feeAmount,
					Description:     &feeDesc,
					IsMatchRequired: true,
					PatternTemplate: &feeType,
				})
			}
		}

		ctx.Data = append(ctx.Data, models.DataAll{
			Transaction:        &transaction,
			TransactionDetails: details,
		})
	}

	return nil
}

// ---------------------------------------------------------------------------
// HTTP helpers
// ---------------------------------------------------------------------------

func getToken(ctx *models.Context) (string, error) {
	services.Logger().STATUS(ctx.Task.Credentials["api_key"].(string))
	token, ok := ctx.Task.Credentials["api_key"].(string)
	if !ok || token == "" {
		return "", fmt.Errorf("stripe: missing credentials.token")
	}
	return token, nil
}

// stripeGet performs an authenticated GET request to the Stripe API.
func stripeGet(token, url string, params map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(token + ":"))
	req.Header.Set("Authorization", "Basic "+encoded)

	if len(params) > 0 {
		q := req.URL.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	return http.DefaultClient.Do(req)
}

// ---------------------------------------------------------------------------
// Pagination helpers
// ---------------------------------------------------------------------------

func fetchAllPayouts(token string, startDate, endDate *string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/payouts?limit=100", stripeBaseURL)
	if startDate != nil {
		url += fmt.Sprintf("&created[gt]=%d", helpers.ToUnixTimestamp(*startDate))
	}
	if endDate != nil {
		url += fmt.Sprintf("&created[lt]=%d", helpers.ToUnixTimestamp(*endDate))
	}

	var all []map[string]interface{}

	for {
		resp, err := stripeGet(token, url, nil)
		if err != nil {
			return nil, err
		}

		var data map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp.Body.Close()

		if items, ok := data["data"].([]interface{}); ok {
			for _, item := range items {
				if m, ok := item.(map[string]interface{}); ok {
					all = append(all, m)
				}
			}
		}

		if hasMore, _ := data["has_more"].(bool); !hasMore {
			break
		}

		// Pagination: starting_after = last item id
		if items, ok := data["data"].([]interface{}); ok && len(items) > 0 {
			last := items[len(items)-1].(map[string]interface{})
			lastID, _ := last["id"].(string)
			url = url + "&starting_after=" + lastID
		}
	}

	return all, nil
}

func fetchAllTransactionsForPayout(token, payoutFilter, dateFilter string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/balance_transactions?limit=100%s%s", stripeBaseURL, payoutFilter, dateFilter)

	var all []map[string]interface{}

	for {
		resp, err := stripeGet(token, url, nil)
		if err != nil {
			return nil, err
		}

		var data map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp.Body.Close()

		if items, ok := data["data"].([]interface{}); ok {
			for _, item := range items {
				if m, ok := item.(map[string]interface{}); ok {
					all = append(all, m)
				}
			}
		}

		if hasMore, _ := data["has_more"].(bool); !hasMore {
			break
		}

		if items, ok := data["data"].([]interface{}); ok && len(items) > 0 {
			last := items[len(items)-1].(map[string]interface{})
			lastID, _ := last["id"].(string)
			url = url + "&starting_after=" + lastID
		}
	}

	return all, nil
}
