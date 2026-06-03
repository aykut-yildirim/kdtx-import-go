package amazon

import (
	"encoding/csv"
	"bytes"

	"myapp/internal/importer"
	"myapp/internal/storage"
)

type TransactionFileImporter struct{}

func (i TransactionFileImporter) Load(
	ctx *importer.Context,
) error {

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

	return nil
}

func (i TransactionFileImporter) Parse(
	ctx *importer.Context,
) error {

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

	return nil
}

func (i TransactionFileImporter) Response(
	ctx *importer.Context,
) (interface{}, error) {

	return ctx.Result, nil
}