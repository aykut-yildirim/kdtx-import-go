package importer

import (
	"fmt"
	"myapp/internal/models"
	"myapp/internal/services"
)

type DefaultPipeline struct {
	Importer models.Importer
}

func (p DefaultPipeline) Run(
	task models.Task,
) (interface{}, error) {

	ctx := &models.Context{
		Task: task,
	}

	services.Logger().STATUS(task.PortalKeyName + " - FileLoad -")
	if err := p.Importer.FileLoad(ctx); err != nil {

		return nil, err

	}
	services.Logger().STATUS(task.PortalKeyName + " - LoginControl -")
	if err := p.Importer.LoginControl(ctx); err != nil {

		return nil, err
	}

	services.Logger().STATUS(task.PortalKeyName + " - Fetch -")
	if err := p.Importer.Fetch(ctx); err != nil {

		return nil, err
	}

	services.Logger().STATUS(task.PortalKeyName + " - Map -")
		if err := p.Importer.Map(ctx); err != nil {

		return nil, err
	}
	if ctx.Task.IsMinioSaved {

	}

	if ctx.Task.IsRawDataAdded {

	}
	if ctx.Task.IsResponse {
		return ctx.Data, nil

	}

	// return ctx.Data, nil
	return map[string]interface{}{
		"IsMinioSaved":           ctx.Task.IsMinioSaved,
		"IsRawDataAdded":         ctx.Task.IsRawDataAdded,
		"IsResponse":             ctx.Task.IsResponse,
		"success":                true,
}, nil
}

type Service struct{}

func (s Service) Execute(
	task models.Task,
) (interface{}, error) {

	services.Logger().STATUS("-- Service Execute")
	// services.Logger().STATUS(fmt.Sprintf("%+v", task))
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
