package importer

import (
	"myapp/internal/models"
)

type Context struct {
	Task      models.Task
	RawData   []byte
	Parsed    interface{}
	Result    interface{}
}