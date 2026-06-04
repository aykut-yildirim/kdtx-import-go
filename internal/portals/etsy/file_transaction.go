package amazon

import (
	"fmt"

	"myapp/internal/helpers"
	"myapp/internal/models"
)

type TransactionFileImporter struct{}

func (i TransactionFileImporter) LoginControl(
	ctx *models.Context,
) error {
	fmt.Println("--LoginControl")

	return nil
}

func (i TransactionFileImporter) FileLoad(
	ctx *models.Context,
) error {

	data, err := helper.GetFile(
		*ctx,
	)

	if err != nil {
		return err
	}

	ctx.RawData = data

	return nil
}

func (i TransactionFileImporter) Fetch(
	ctx *models.Context,
) error {

	return nil
}


func (i TransactionFileImporter) Map(
	ctx *models.Context,
)  (interface{}, error)  {

	rows := ctx.Parsed.([][]string)

	var transactions []map[string]interface{}

	for _, row := range rows {

		transactions = append(
			transactions,
			map[string]interface{}{
				"row": row,
			},
		)
	}

	ctx.Result = transactions

	return ctx.Result, nil
}
