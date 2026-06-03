package importer

import (
	"fmt"
	"sync"
)

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
	fmt.Println("Get", key)
	mu.RLock()
	defer mu.RUnlock()

	imp, ok := registry[key]

	return imp, ok
}