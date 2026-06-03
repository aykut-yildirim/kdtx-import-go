package importer

import (
	"fmt"
	"myapp/internal/models"
)

type Service struct{}

func (s Service) Execute(
	task models.Task,
) (interface{}, error) {

	fmt.Println("-- Service Execute")
	fmt.Println(task)
	key := fmt.Sprintf(
		"%s:%s:%s",
		task.PortalKeyName,
		task.InputType,
		task.PortalType,
	)

	imp, ok := Get(key)

	if !ok {

		return nil,
			fmt.Errorf(
				"portal not found: %s",
				key,
			)
	}

	pipeline := DefaultPipeline{
		Importer: imp,
	}

	return pipeline.Run(task)
}