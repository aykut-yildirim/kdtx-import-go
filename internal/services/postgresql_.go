package services

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// --------------------
// SERVICE
// --------------------

type PostgresService struct {
	pool      *pgxpool.Pool
	schema    string
	tableName string
}

// constructor
func NewPostgresService(ctx context.Context, schema string, tableName string) (*PostgresService, error) {

	if schema == "" {
		schema = "public"
	}

	if tableName == "" {
		tableName = "logs"
	}

	connStr := os.Getenv("POSTGRES_URI")

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}

	s := &PostgresService{
		pool:      pool,
		schema:    schema,
		tableName: tableName,
	}

	if err := s.createTableIfNotExists(ctx); err != nil {
		return nil, err
	}

	return s, nil
}

// --------------------
// CREATE TABLE
// --------------------

func (p *PostgresService) createTableIfNotExists(ctx context.Context) error {

	query := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s.%s (
id SERIAL PRIMARY KEY,
level VARCHAR(20) NOT NULL,
message TEXT NOT NULL,
client_message TEXT,
task JSONB,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
`, p.schema, p.tableName)

	_, err := p.pool.Exec(ctx, query)
	return err
}

// --------------------
// INSERT ONE (dynamic)
// --------------------

func (p *PostgresService) InsertOne(
	ctx context.Context,
	tableName string,
	data map[string]interface{},
) error {

	cols := []string{}
	placeholders := []string{}
	values := []interface{}{}

	i := 1
	for k, v := range data {
		cols = append(cols, k)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		values = append(values, v)
		i++
	}

	query := fmt.Sprintf(
		"INSERT INTO %s.%s (%s) VALUES (%s)",
		p.schema,
		tableName,
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "),
	)

	_, err := p.pool.Exec(ctx, query, values...)
	return err
}

// --------------------
// INSERT MANY
// --------------------

func (p *PostgresService) InsertMany(
	ctx context.Context,
	tableName string,
	data []map[string]interface{},
) error {

	if len(data) == 0 {
		return nil
	}

	for _, row := range data {
		if err := p.InsertOne(ctx, tableName, row); err != nil {
			return err
		}
	}

	return nil
}

// --------------------
// INSERT LOG
// --------------------

func (p *PostgresService) InsertLog(
	ctx context.Context,
	level string,
	message string,
	clientMessage *string,
	task interface{},
) error {

	query := fmt.Sprintf(`
INSERT INTO %s.%s
(level, message, client_message, task)
VALUES ($1, $2, $3, $4)
`, p.schema, p.tableName)

	_, err := p.pool.Exec(
		ctx,
		query,
		level,
		message,
		clientMessage,
		task,
	)

	return err
}

// --------------------
// GET LOGS (WHERE)
// --------------------

func (p *PostgresService) GetLogs(
	ctx context.Context,
	where map[string]interface{},
) ([]map[string]interface{}, error) {

	base := fmt.Sprintf("SELECT * FROM %s.%s", p.schema, p.tableName)

	values := []interface{}{}
	conditions := []string{}

	i := 1

	for k, v := range where {
		conditions = append(conditions, fmt.Sprintf("%s = $%d", k, i))
		values = append(values, v)
		i++
	}

	query := base

	if len(where) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := p.pool.Query(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []map[string]interface{}{}

	for rows.Next() {

		row, err := rows.Values()
		if err != nil {
			return nil, err
		}

		results = append(results, map[string]interface{}{
			"row": row,
		})
	}

	return results, nil
}

// --------------------
// CLOSE
// --------------------

func (p *PostgresService) Close() {
	p.pool.Close()
}
