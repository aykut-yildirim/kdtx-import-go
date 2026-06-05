package helpers

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"myapp/internal/models"
	"myapp/internal/services"
	"os"

	"github.com/joho/godotenv"
	"strings"
)

func GetFile(ctx models.Context) ([]byte, error) {
	services.Logger().STATUS("--GetFile")
	if ctx.Task.FileMinioPath != nil {
		services.Logger().STATUS("--GetFile--FileMinioPath")
		godotenv.Load(".env")
		bucket := os.Getenv("MINIO_BUCKET")
		key := *ctx.Task.FileMinioPath
		minioSvc, err := services.NewMinioService()
		if err != nil {
			return []byte{}, fmt.Errorf("failed to create minio service: %w", err)
		}
		fileBytes, err := minioSvc.GetFileByPath(ctx, bucket, key)

		if err != nil {
			return nil, err
		}
		return fileBytes, nil
	}

	if ctx.Task.LocalPath != nil {
		services.Logger().STATUS("--GetFile--LocalPath")
		// services.Logger().STATUS(ctx.Task.LocalPath)
		services.Logger().STATUS(*ctx.Task.LocalPath)

		fileBytes, err := os.ReadFile(*ctx.Task.LocalPath)

		if err != nil {
			return nil, err
		}
		return fileBytes, nil
	}

	return nil, errors.New("missing file information")
}

func ParseDatetime(val string) (time.Time, error) {
	layouts := []string{
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05-0700",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}
	var lastErr error
	for _, l := range layouts {
		t, err := time.Parse(l, val)
		if err == nil {
			return t, nil
		}
		lastErr = err
	}
	return time.Time{}, lastErr
}

var moneyRegex = regexp.MustCompile(`[^\d,.\-]`)

func MoneyToFloat(value any) (float64, error) {
	if value == nil {
		return 0, nil
	}

	switch v := value.(type) {
	case int:
		return float64(v), nil
	case float64:
		return v, nil
	}

	str := fmt.Sprint(value)
	clean := moneyRegex.ReplaceAllString(str, "")
	clean = strings.ReplaceAll(clean, ",", ".")

	if !regexp.MustCompile(`\d`).MatchString(clean) {
		return 0, nil
	}

	return strconv.ParseFloat(clean, 64)
}

func ToStringCombine(sep string, values ...any) string {
	out := []string{}
	for _, v := range values {
		if v == nil {
			continue
		}
		out = append(out, fmt.Sprint(v))
	}
	return strings.Join(out, sep)
}

func ToFloatTotal(values ...string) float64 {
	var total float64
	for _, val := range values {
		if val != "" {
			f, err := MoneyToFloat(val)
			if err == nil {
				total += f
			}
		}
	}
	return total
}

func ReadTableFromBytes(content []byte, sep string) ([]map[string]string, error) {
	lines := strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return nil, errors.New("empty table content")
	}

	headers := strings.Split(lines[0], sep)
	for i, h := range headers {
		headers[i] = strings.TrimSpace(h)
	}

	var records []map[string]string
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		cols := strings.Split(lines[i], sep)
		row := make(map[string]string)
		for j, header := range headers {
			if j < len(cols) {
				row[header] = strings.TrimSpace(cols[j])
			} else {
				row[header] = ""
			}
		}
		records = append(records, row)
	}
	return records, nil
}
