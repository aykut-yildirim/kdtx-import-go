package storage

import "os"

func ReadLocalFile(
	path string,
) ([]byte, error) {

	return os.ReadFile(path)
}