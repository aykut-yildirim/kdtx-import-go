package amazon

import "myapp/internal/importer"

func init() {

	importer.Register(
		"amazon:file:transaction",
		TransactionFileImporter{},
	)

	// importer.Register(
	// 	"amazon:api:transaction",
	// 	TransactionAPIImporter{},
	// )

	// importer.Register(
	// 	"amazon:file:invoice_out",
	// 	InvoiceOutFileImporter{},
	// )
}