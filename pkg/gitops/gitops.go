package gitops

import "github.com/xctl/pkg/gitops/github"

type GitopsInterface interface {
	Application() ApplicationInterface
}

type Gitops struct {
	client *github.GithubClient
}

func NewGitops(client *github.GithubClient) *Gitops {
	return &Gitops{
		client: client,
	}
}

func (c *Gitops) Application() ApplicationInterface {
	return newApplication(c)
}
