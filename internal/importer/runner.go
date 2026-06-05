package importer

import (
	"myapp/internal/models"
	"myapp/internal/services"
	"sync"
)

type Importer interface {
	FileLoad(*models.Context) error
	LoginControl(*models.Context) error
	Fetch(*models.Context) error
	Map(*models.Context) (interface{}, error)
}

// type Pipeline interface {
// 	Run(task models.Task) (interface{}, error)
// }

var (
	mu       sync.RWMutex
	registry = map[string]Importer{}
)

func Register(
	key string,
	imp Importer,
) {
	mu.Lock()
	defer mu.Unlock()

	registry[key] = imp
}

func Get(
	key string,
) (Importer, bool) {
	services.Logger().STATUS("Get" + key)
	mu.RLock()
	defer mu.RUnlock()

	imp, ok := registry[key]

	return imp, ok
}

// func BuildKey(task models.Task) string {

// 	return fmt.Sprintf(
// 		"%s:%s:%s",
// 		task.PortalKeyName,
// 		task.InputType,
// 		task.PortalType,
// 	)
// }
