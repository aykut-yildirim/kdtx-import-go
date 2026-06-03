package importer

import "myapp/internal/models"

type DefaultPipeline struct {
	Importer Importer
}

func (p DefaultPipeline) Run(
	task models.Task,
) (interface{}, error) {

	ctx := &Context{
		Task: task,
	}

	if err := p.Importer.Load(ctx); err != nil {
		return nil, err
	}

	if err := p.Importer.Fetch(ctx); err != nil {
		return nil, err
	}

	if err := p.Importer.Parse(ctx); err != nil {
		return nil, err
	}

	if err := p.Importer.Transform(ctx); err != nil {
		return nil, err
	}

	return p.Importer.Response(ctx)
}