package server

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

type entryResource struct{}

func (es entryResource) mountRoutes(h huma.API) {
	huma.Get(h, "/entries", es.listEntries)
}

type listEntriesOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

func (es entryResource) listEntries(ctx context.Context, _ *struct{}) (*listEntriesOutput, error) {
	logger := getLogger(ctx)
	logger.Info("eyyy here")
	resp := &listEntriesOutput{}
	resp.Body.Message = "hello world"
	return resp, nil
}
