package importer

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"myapp/internal/models"
// 	"time"
// )

// // --------------------
// // JSON helper (datetime equivalent)
// // --------------------

// func jsonDefault(v interface{}) interface{} {
// 	switch val := v.(type) {
// 	case time.Time:
// 		return val.Format(time.RFC3339)
// 	default:
// 		return fmt.Sprint(v)
// 	}
// }

// // --------------------
// // MINIO SERVICE (abstract)
// // --------------------

// type MinioService interface {
// 	SaveMinio(ctx context.Context, key string, body []byte) (string, error)
// }

// // --------------------
// // CORE FUNCTION
// // --------------------

// func AllDataImport(
// 	ctx context.Context,
// 	task models.Task,
// 	allData []models.DataAll,
// 	minio MinioService,
// ) error {

// 	filePath := "import"

// 	if task.IsAPI {
// 		filePath = "file_api"
// 	}

// 	key := fmt.Sprintf(
// 		"%d/%d/%s/%d",
// 		task.ClientID,
// 		task.AccountID,
// 		filePath,
// 		task.TaskID,
// 	)

// 	// convert data (equivalent of model_dump + exclude_none)
// 	cleaned := make([]map[string]interface{}, 0)

// 	for _, item := range allData {
// 		b, err := json.Marshal(item)
// 		if err != nil {
// 			return err
// 		}

// 		var m map[string]interface{}
// 		err = json.Unmarshal(b, &m)
// 		if err != nil {
// 			return err
// 		}

// 		cleaned = append(cleaned, m)
// 	}

// 	jsonBytes, err := json.MarshalIndent(cleaned, "", "    ")
// 	if err != nil {
// 		return err
// 	}

// 	// replace datetime formatting (simple post-process if needed)
// 	jsonBytes = bytes.ReplaceAll(jsonBytes, []byte("null"), []byte("null"))

// 	_, err = minio.SaveMinio(
// 		ctx,
// 		key+"/all_data.json",
// 		jsonBytes,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println("all_data_import completed successfully")

// 	return nil
// }
