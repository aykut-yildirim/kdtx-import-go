package importer

import "myapp/internal/models"

type Importer interface {
	Load(*Context) error
	Fetch(*Context) error
	Parse(*Context) error
	Transform(*Context) error
	Response(*Context) (interface{}, error)
}

type Pipeline interface {
	Run(task models.Task) (interface{}, error)
}