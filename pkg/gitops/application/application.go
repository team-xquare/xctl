package application

import (
	"context"

	"github.com/xctl/pkg/api"
	"github.com/xctl/pkg/gitops/github"
)

var appTemplatePath = "template/app-template"

type ApplicationInterface interface {
	Create(ctx context.Context) error
	Get(ctx context.Context) ([]*api.Application, error)
	GetFrontend(ctx context.Context) ([]*api.Application, error)
	GetBackend(ctx context.Context) ([]*api.Application, error)
	Delete(ctx context.Context) error
}

type applications struct {
	application *api.Application
	client      *github.GithubClient
}

func NewApplications(client *github.GithubClient, app *api.Application) *applications {
	return &applications{
		client:      client,
		application: app,
	}
}
