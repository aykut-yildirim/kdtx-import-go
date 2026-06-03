package amazon

import (
	"bytes"
	"encoding/csv"
	"fmt"

	"myapp/internal/importer"
	"myapp/internal/storage"
)

type TransactionFileImporter struct{}

func (i TransactionFileImporter) Load(
	ctx *importer.Context,
) error {
	fmt.Println("--Load")

	data, err := storage.ReadLocalFile(
		*ctx.Task.LocalPath,
	)

	if err != nil {
		return err
	}

	ctx.RawData = data

	return nil
}

func (i TransactionFileImporter) Fetch(
	ctx *importer.Context,
) error {
	fmt.Println("--Fetch")

	return nil
}

func (i TransactionFileImporter) Parse(
	ctx *importer.Context,
) error {
	fmt.Println("--Parse")
	reader := csv.NewReader(
		bytes.NewReader(ctx.RawData),
	)

	rows, err := reader.ReadAll()

	if err != nil {
		return err
	}

	ctx.Parsed = rows

	return nil
}

func (i TransactionFileImporter) Transform(
	ctx *importer.Context,
) error {
	fmt.Println("--Transform")

	rows := ctx.Parsed.([][]string)

	var transactions []map[string]interface{}

	for index, row := range rows {

		fmt.Println(index, row)

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

func (i TransactionFileImporter) Response(
	ctx *importer.Context,
) (interface{}, error) {
	fmt.Println("--Response")

	return ctx.Result, nil
}