package models

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-playground/validator/v10"
)

type Task struct {
	PortalKeyName		string					`json:"portal_key_name" validate:"required"`
	InputType			string					`json:"input_type" validate:"required,oneof=api file"`
	PortalType			string					`json:"portal_type" validate:"required,oneof=transaction invoice_out invoice_in"`
	AccountID			int						`json:"account_id"`
	ClientID			int						`json:"client_id"`
	TaskID				int						`json:"task_id"`
	FetchStartDate		*string					`json:"fetch_start_date,omitempty"`
	FetchEndDate		*string					`json:"fetch_end_date,omitempty"`
	FileMinioPath		*string					`json:"file_minio_path,omitempty"`
	LocalPath			*string					`json:"local_path,omitempty"`
	Credentials			map[string]interface{}	`json:"credentials,omitempty"`
	AccountingPatterns	map[string]interface{}	`json:"accounting_patterns,omitempty"`
	IsResponse			bool					`json:"is_response"`
	IsMinioSaved		bool					`json:"Is_minio_saved"`
	IsRawDataAdded		bool					`json:"is_raw_data_added"`
	Log					*string					`json:"log,omitempty"`
}

type Context struct {
	Task		Task			
	RawData		[]byte			
	Data		[]DataAll 		`json:"data"`
	Parsed		interface{}		
	Result		interface{}		
}

type Importer interface {
	FileLoad(*Context) error
	LoginControl(*Context) error
	Fetch(*Context) error
	Map(*Context) error
}

type DataAll struct {
	Task               *Task               `json:"task,omitempty"`
	Invoice            *Invoice            `json:"invoice,omitempty"`
	InvoiceDetails     []InvoiceDetail     `json:"invoice_details,omitempty"`
	Transaction        *Transaction        `json:"transaction,omitempty"`
	TransactionDetails []TransactionDetail `json:"transaction_details,omitempty"`
	SettlementList     []Settlement        `json:"settlement_list,omitempty"`
	AllIDs             []string            `json:"all_ids,omitempty"`
	FilteredIDs        []string            `json:"filtered_ids,omitempty"`
}

type DataAllList struct {
	Task      Task      `json:"task"`
	Data      []DataAll `json:"data"`
	PDFFolder *string   `json:"pdf_folder,omitempty"`
}

var validate = validator.New()

func (t *Task) Validate() error {

	err := validate.Struct(t)

	if err != nil {
		return err
	}

	if t.InputType == "api" {

		if t.FetchStartDate == nil || t.FetchEndDate == nil {
			return errors.New(
				"fetch_start_date and fetch_end_date are required",
			)
		}

		startDate, err := time.Parse(
			"2006-01-02",
			*t.FetchStartDate,
		)

		if err != nil {
			return errors.New(
				"invalid fetch_start_date format",
			)
		}

		endDate, err := time.Parse(
			"2006-01-02",
			*t.FetchEndDate,
		)

		if err != nil {
			return errors.New(
				"invalid fetch_end_date format",
			)
		}

		if startDate.After(endDate) {
			return errors.New(
				"fetch_start_date cannot be greater than fetch_end_date",
			)
		}

		if t.Credentials == nil {
			return errors.New(
				"credentials are required for api input_type",
			)
		}
	}

	baseDir, _ := os.Getwd()

	modulePath := filepath.Join(
		baseDir,
		"worker",
		"portals",
		t.PortalKeyName,
		fmt.Sprintf(
			"%s_%s.py",
			t.InputType,
			t.PortalType,
		),
	)

	if _, err := os.Stat(modulePath); os.IsNotExist(err) {
		return fmt.Errorf(
			"service %s_%s for portal %s is not implemented",
			t.InputType,
			t.PortalType,
			t.PortalKeyName,
		)
	}

	return nil
}

