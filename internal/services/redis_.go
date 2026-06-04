package services

import (
	"context"
	"encoding/json"
	"os"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	client *redis.Client

	stream string
	group  string
}

// constructor
func NewRedisService() (*RedisService, error) {

	url := os.Getenv("REDIS_URL")
	if url == "" {
		return nil, ErrMissingRedisURL
	}

	stream := os.Getenv("REDIS_STREAM")
	if stream == "" {
		stream = "import_stream"
	}

	group := os.Getenv("REDIS_GROUP")
	if group == "" {
		group = "import_group"
	}

	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	return &RedisService{
		client: client,
		stream: stream,
		group:  group,
	}, nil
}

// --------------------
// ERR
// --------------------

var ErrMissingRedisURL = redis.Nil

// --------------------
// HASH OPERATIONS
// --------------------

func (r *RedisService) GetTask(
	ctx context.Context,
	task map[string]string,
) (string, error) {

	return r.client.HGet(ctx, task["portal_key_name"], task["task_id"]).Result()
}

func (r *RedisService) SetDescription(
	ctx context.Context,
	taskID string,
	message string,
) error {

	return r.client.HSet(ctx, taskID, map[string]interface{}{
		"description": message,
	}).Err()
}

func (r *RedisService) DeleteTask(
	ctx context.Context,
	task map[string]string,
) error {

	return r.client.HDel(ctx, task["portal_key_name"], task["task_id"]).Err()
}

// --------------------
// STREAM OPERATIONS
// --------------------

func (r *RedisService) StreamAdd(
	ctx context.Context,
	data interface{},
) (string, error) {

	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: r.stream,
		ID:     "*",
		Values: map[string]interface{}{
			"data": string(b),
		},
	}).Result()
}

func (r *RedisService) StreamCreateGroup(
	ctx context.Context,
) error {

	err := r.client.XGroupCreateMkStream(
		ctx,
		r.stream,
		r.group,
		"0-0",
	).Err()

	// Python suppress error behavior equivalent
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return err
	}

	return nil
}

func (r *RedisService) StreamGetGroup(
	ctx context.Context,
	consumer string,
	count int64,
	block int,
) ([]redis.XStream, error) {

	return r.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    r.group,
		Consumer: consumer,
		Streams:  []string{r.stream, ">"},
		Count:    count,
		Block:    0, // 0 = block indefinitely (adjust if needed)
	}).Result()
}

func (r *RedisService) StreamAck(
	ctx context.Context,
	messageID string,
) (int64, error) {

	return r.client.XAck(ctx, r.stream, r.group, messageID).Result()
}

func (r *RedisService) StreamDelete(
	ctx context.Context,
	messageID string,
) (int64, error) {

	return r.client.XDel(ctx, r.stream, messageID).Result()
}

func (r *RedisService) StreamLen(
	ctx context.Context,
) (int64, error) {

	return r.client.XLen(ctx, r.stream).Result()
}

func (r *RedisService) StreamRead(
	ctx context.Context,
	start string,
	end string,
	count int64,
) ([]redis.XMessage, error) {

	return r.client.XRangeN(ctx, r.stream, start, end, count).Result()
}
