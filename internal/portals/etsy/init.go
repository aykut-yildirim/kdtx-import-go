package amazon

import "myapp/internal/importer"

func init() {

	importer.Register(
		"etsy:file:transaction",
		TransactionFileImporter{},
	)

}