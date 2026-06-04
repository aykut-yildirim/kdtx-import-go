package models

import (
	"time"
)

type Invoice struct {
	IsIncoming            bool                   `json:"is_incoming"`
	CurrencyCode          string                 `json:"currency_code" validate:"len=3"`
	ExchangeRate          float64                `json:"exchange_rate"`
	Amount                float64                `json:"amount"`
	NetAmount             float64                `json:"net_amount"`
	TaxAmount             float64                `json:"tax_amount"`
	OpenAmount            float64                `json:"open_amount"`
	Date                  time.Time              `json:"date"`
	DueDate               *time.Time             `json:"due_date,omitempty"`
	DeliveryDate          *time.Time             `json:"delivery_date,omitempty"`
	InvoiceNo             *string                `json:"invoice_no,omitempty"`
	SettlementNo          *string                `json:"settlement_no,omitempty"`
	TransactionNo         *string                `json:"transaction_no,omitempty"`
	ReferenceNo           *string                `json:"reference_no,omitempty"`
	OrderNo               *string                `json:"order_no,omitempty"`
	Description           *string                `json:"description,omitempty"`
	PaymentData           map[string]interface{} `json:"payment_data,omitempty"`
	ExternalReferenceData map[string]interface{} `json:"external_reference_data,omitempty"`
	PartyAccountNo        *int                   `json:"party_account_no,omitempty"`
	PartyName             *string                `json:"party_name,omitempty"`
	PartyVATID            *string                `json:"party_vat_id,omitempty"`
	PartyTaxID            *string                `json:"party_tax_id,omitempty"`
	PartyEmail            *string                `json:"party_email,omitempty"`
	PartyPhone            *string                `json:"party_phone,omitempty"`
	PartyAddress          *string                `json:"party_address,omitempty"`
	PartyCity             *string                `json:"party_city,omitempty"`
	PartyZipCode          *string                `json:"party_zip_code,omitempty"`
	PartyCountry          *string                `json:"party_country,omitempty"`
	PartyIBAN             *string                `json:"party_iban,omitempty"`
	PartyBIC              *string                `json:"party_bic,omitempty"`
	DepotVATID            *string                `json:"depot_vat_id,omitempty"`
	DocumentPaths         []string               `json:"document_paths"`
	RawData               interface{}            `json:"raw_data,omitempty"`
}

type InvoiceDetail struct {
	Amount          float64 `json:"amount"`
	NetAmount       float64 `json:"net_amount"`
	TaxAmount       float64 `json:"tax_amount"`
	TaxRate         float64 `json:"tax_rate"`
	Quantity        float64 `json:"quantity"`
	Description     *string `json:"description,omitempty"`
	IsService       *bool   `json:"is_service,omitempty"`
	LedgerAccountNo *string `json:"ledger_account_no,omitempty"`
}

type ExternalReferenceData struct {
	ExternalID       *string `json:"ExternalId,omitempty"`
	ExternalSourceID *string `json:"ExternalSourceId,omitempty"`
	CustomIdentifier *string `json:"CustomIdentifier,omitempty"`
	InvoiceNo        *string `json:"InvoiceNo,omitempty"`
	OrderNo          *string `json:"OrderNo,omitempty"`
	PaymentNo        *string `json:"PaymentNo,omitempty"`
	PaymentID        *string `json:"PaymentId,omitempty"`
	BookingText      *string `json:"BookingText,omitempty"`
	PartnerName      *string `json:"PartnerName,omitempty"`
}

type Transaction struct {
	Amount                float64                `json:"amount"`
	ExchangeRate          float64                `json:"exchange_rate"`
	OpenAmount            float64                `json:"open_amount"`
	CurrencyCode          string                 `json:"currency_code"`
	IsDebit               bool                   `json:"is_debit"`
	Date                  time.Time              `json:"date"`
	Description           *string                `json:"description,omitempty"`
	DocumentPaths         []string               `json:"document_paths"`
	ExternalReferenceData map[string]interface{} `json:"external_reference_data,omitempty"`
	TransactionNo         *string                `json:"transaction_no,omitempty"`
	OrderNo               *string                `json:"order_no,omitempty"`
	ReferenceNo           *string                `json:"reference_no,omitempty"`
	SettlementNo          *string                `json:"settlement_no,omitempty"`
	RawData               interface{}            `json:"raw_data,omitempty"`
}

type Settlement struct {
	SettlementNo    string      `json:"settlement_no"`
	Amount          float64     `json:"amount"`
	CurrencyCode    string      `json:"currency_code"`
	ExchangeRate    float64     `json:"exchange_rate"`
	StartDate       time.Time   `json:"start_date"`
	EndDate         time.Time   `json:"end_date"`
	Description     *string     `json:"description,omitempty"`
	Count           int         `json:"count"`
	TransactionList []string    `json:"transaction_list"`
	RawData         interface{} `json:"raw_data,omitempty"`
}

type TransactionDetail struct {
	Amount            float64 `json:"amount"`
	TaxRate           float64 `json:"tax_rate"`
	Description       *string `json:"description,omitempty"`
	IsMatchRequired   bool    `json:"is_match_required"`
	AccountingCodeIn  *int    `json:"accounting_code_in,omitempty"`
	AccountingCodeOut *int    `json:"accounting_code_out,omitempty"`
	PatternTemplate   *string `json:"pattern_template,omitempty"`
}
