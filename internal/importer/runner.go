package importer

import (
	"myapp/internal/models"
	"myapp/internal/services"
	"sync"
)

// type Pipeline interface {
// 	Run(task models.Task) (interface{}, error)
// }

var (
	mu       sync.RWMutex
	registry = map[string]models.Importer{}
)

func Register(
	key string,
	imp models.Importer,
) {
	services.Logger().STATUS("--Register" + key)

	mu.Lock()
	defer mu.Unlock()

	registry[key] = imp
}

func Get(
	key string,
) (models.Importer, bool) {
	services.Logger().STATUS("-- Get / " + key)
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
