package amazon

import (
	"bytes"
	"encoding/csv"

	"myapp/internal/helpers"
	"myapp/internal/models"
	"myapp/internal/services"
)

type TransactionFileImporter struct{}

func (i TransactionFileImporter) FileLoad(
	ctx *models.Context,
) error {
	services.Logger().STATUS("- FileLoad -")
	data, err := helpers.GetFile(*ctx)

	services.Logger().STATUS("test")

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
) (interface{}, error) {

	var transactions []map[string]interface{}

	reader := csv.NewReader(
		bytes.NewReader(ctx.RawData),
	)

	rows, err := reader.ReadAll()

	if err != nil {
		return transactions, err
	}

	for index, row := range rows {

		// services.Logger().STATUS(index, row)

		transactions = append(
			transactions,
			map[string]interface{}{
				"row_index": index,
				"row_value": row,
			},
		)
	}

	ctx.Result = transactions

	return ctx.Result, nil
}
