package amazon

import "myapp/internal/importer"

func init() {

	importer.Register(
		"stripe:api:transaction",
		TransactionApiImporter{},
	)
	
	// importer.Register(
	// 	"stripe:file:transaction",
	// 	TransactionFileImporter{},
	// )
	
}