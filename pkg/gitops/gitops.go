package gitops

import (
	"github.com/xctl/pkg/api"
	"github.com/xctl/pkg/gitops/application"
	"github.com/xctl/pkg/gitops/github"
)

type GitopsInterface interface {
	Application(app *api.Application) application.ApplicationInterface
}

type Gitops struct {
	client *github.GithubClient
}

func NewGitops(client *github.GithubClient) *Gitops {
	return &Gitops{
		client: client,
	}
}

func (c *Gitops) Application(app *api.Application) application.ApplicationInterface {
	return application.NewApplications(c.client, app)
}
