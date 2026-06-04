package services

import (
	"fmt"
	"os"
	"sync"
	"time"
	"google.golang.org/grpc"
)

var (
	grpcEnabled = os.Getenv("GRPC_ENABLED") == "true"
	grpcAddr    = os.Getenv("GRPC_IMPORT_WORKER_ADDR")
)

// ---------------- COLORS ----------------

var colors = map[string]string{
	"-INFO-":   "\033[92m",
	"SUCCESS":  "\033[92m",
	"ERROR":    "\033[91m",
	"STARTED":  "\033[94m",
	"STATUS":   "\033[94m",
	"FINISHED": "\033[94m",
	"CRITICAL": "\033[1;91m",
	"SYSTEM":   "\033[93m",
}

const reset = "\033[0m"

type LoggerService struct{}

func (l *LoggerService) write(level, message string) {
	color := colors[level]
	now := time.Now().UTC()

	fmt.Printf("%s %s[%s]%s %s\n",
		now.Format(time.RFC3339),
		color,
		level,
		reset,
		message,
	)
}

type Log struct {
	logger *LoggerService
	
	mu   sync.Mutex
	conn *grpc.ClientConn
}

func NewLog() *Log {
	return &Log{
		logger: &LoggerService{},
	}
}

func (l *Log) INFO(message string) {
	l.logger.write("-INFO-", message)
}

func (l *Log) SUCCESS(message string) {
	l.logger.write("SUCCESS", message)
}

func (l *Log) SYSTEM(message string) {
	l.logger.write("SYSTEM", message)
}

func (l *Log) STATUS(message string) {
	l.logger.write("STATUS", message)
}


func (l *Log) STARTED(taskID string) {
	l.logger.write("STARTED", "Portal Started" + taskID )

}

func (l *Log) FINISHED(taskID string) {

	l.logger.write("FINISHED", "Portal Finished" + taskID)
}

func (l *Log) ERROR(err error) {
	l.logger.write("ERROR", err.Error())

}

var LoggerServiceInstance = NewLog()

func LoggerSystem(message string) {
	LoggerServiceInstance.SYSTEM(message)
}

func LoggerStatus(message string, taskID int) {
	LoggerServiceInstance.STATUS(fmt.Sprintf("%s (TaskID: %d)", message, taskID))
}
