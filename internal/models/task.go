package models

type Task struct {
	PortalKeyName string `json:"portal_key_name" validate:"required"`
	InputType     string `json:"input_type" validate:"required,oneof=api file"`
	PortalType    string `json:"portal_type" validate:"required,oneof=transaction invoice_out invoice_in"`

	AccountID int `json:"account_id"`
	ClientID  int `json:"client_id"`
	TaskID    int `json:"task_id"`

	FetchStartDate *string `json:"fetch_start_date,omitempty"`
	FetchEndDate   *string `json:"fetch_end_date,omitempty"`

	Credentials map[string]interface{} `json:"credentials,omitempty"`

	FileMinioPath *string `json:"file_minio_path,omitempty"`
	LocalPath     *string `json:"local_path,omitempty"`

	AccountingPatterns map[string]interface{} `json:"accounting_patterns,omitempty"`

	IsAPI      bool `json:"is_api"`
	IsResponse bool `json:"is_response"`

	Log *string `json:"log,omitempty"`
}
