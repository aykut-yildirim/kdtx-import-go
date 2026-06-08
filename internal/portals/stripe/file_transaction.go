package amazon

import (
	"bytes"
	"encoding/csv"

	"myapp/internal/helpers"
	"myapp/internal/models"
)

type TransactionFileImporter struct{}

func (i TransactionFileImporter) FileLoad(
	ctx *models.Context,
) error {
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

func (i TransactionFileImporter) Map(
	ctx *models.Context,
)  error {

	var transactions []map[string]interface{}

	reader := csv.NewReader(
		bytes.NewReader(ctx.RawData),
	)

	rows, err := reader.ReadAll()

	if err != nil {
		return  err
	}

	for index, row := range rows {
		transactions = append(
			transactions,
			map[string]interface{}{
				"row_index": index,
				"row_value": row,
			},
		)
	}

	ctx.Result = transactions

	return nil
}
