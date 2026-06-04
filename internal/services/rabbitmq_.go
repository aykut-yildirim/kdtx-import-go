package services

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	rabbitURL   = os.Getenv("RABBITMQ_URL")
	rabbitQueue = os.Getenv("RABBITMQ_QUEUE")
)

// ---------------- SERVICE ----------------

type RabbitmqService struct {
	url   string
	queue string

	mu   sync.Mutex
	conn *amqp.Connection
}

// ---------------- CONSTRUCTOR ----------------

func NewRabbitmqService(url, queue string) (*RabbitmqService, error) {

	if url == "" {
		url = rabbitURL
	}
	if queue == "" {
		queue = rabbitQueue
	}

	if url == "" {
		return nil, errors.New("RABBITMQ_URL is not defined")
	}
	if queue == "" {
		return nil, errors.New("RABBITMQ_QUEUE is not defined")
	}

	return &RabbitmqService{
		url:   url,
		queue: queue,
	}, nil
}

// ---------------- CONNECTION ----------------

func (r *RabbitmqService) getConn() (*amqp.Connection, error) {

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.conn != nil {
		return r.conn, nil
	}

	conn, err := amqp.DialConfig(r.url, amqp.Config{
		Heartbeat: 60,
	})

	if err != nil {
		return nil, err
	}

	r.conn = conn
	return conn, nil
}

// ---------------- PUBLISH ----------------

func (r *RabbitmqService) Publish(data interface{}, queueOverride string) error {

	queue := queueOverride
	if queue == "" {
		queue = r.queue
	}

	conn, err := r.getConn()
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		queue,
		true, // durable
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ch.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
}

// ---------------- CLOSE ----------------

func (r *RabbitmqService) Close() error {
	if r.conn == nil {
		return nil
	}
	return r.conn.Close()
}

// ---------------- GLOBAL INSTANCE ----------------

var RabbitMQServiceInstance *RabbitmqService

func InitRabbitMQ() error {
	svc, err := NewRabbitmqService("", "")
	if err != nil {
		return err
	}

	RabbitMQServiceInstance = svc
	return nil
}

func RabbitPublish(data interface{}, queueOverride string) error {
	if RabbitMQServiceInstance == nil {
		if err := InitRabbitMQ(); err != nil {
			return err
		}
	}
	return RabbitMQServiceInstance.Publish(data, queueOverride)
}

// services.RabbitMQServiceInstance.Publish(
// 	map[string]interface{}{
// 		"task_id": "123",
// 		"status":  "created",
// 	},
// 	"",
// )
