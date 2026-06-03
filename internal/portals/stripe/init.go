package amazon

import "myapp/internal/importer"

func init() {

	importer.Register(
		"stripe:file:transaction",
		TransactionFileImporter{},
	)
	
}