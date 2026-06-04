package importer

import (
	"fmt"
	"myapp/internal/models"
)
type DefaultPipeline struct {
	Importer Importer
}

func (p DefaultPipeline) Run(
	task models.Task,
) (interface{}, error) {

	ctx := &models.Context{
		Task: task,
	}

	fmt.Println(task.PortalKeyName, "- FileLoad -")
	if err := p.Importer.FileLoad(ctx); err != nil {

		return nil, err

	}

	fmt.Println(task.PortalKeyName, "- LoginControl -")
	if err := p.Importer.LoginControl(ctx); err != nil {

		return nil, err
	}
	
	fmt.Println(task.PortalKeyName, "- Fetch -")
	if err := p.Importer.Fetch(ctx); err != nil {
		
		return nil, err
	}
	
	fmt.Println(task.PortalKeyName, "- Map -")
	return p.Importer.Map(ctx)
}

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