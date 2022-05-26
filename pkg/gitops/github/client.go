package github

import (
	"context"
	"time"

	"github.com/google/go-github/github"
	"github.com/xctl/pkg/gitops/config"
	"golang.org/x/oauth2"
)

var (
	DefaultRepoOwner      = "team-xquare"
	DefaultCommitterName  = "xquare-admin"
	DefaultCommitterEmail = "teamxquare@gmail.com"

	StagingRepo    = "xquare-gitops-repo-staging"
	ProductionRepo = "xquare-gitops-repo-production"

	FileType = github.String("blob")
	ReadMode = github.String("100644")
)

type GithubClient struct {
	Client         *github.Client
	Repo           string
	RepoOwner      string
	CommitterName  string
	CommitterEmail string
}

func NewGithubClient(repo string) (*GithubClient, error) {
	credential, err := config.GetCredential()
	if err != nil {
		return nil, err
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: credential.GithubToken,
	})
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	return &GithubClient{
		Client:         client,
		Repo:           repo,
		RepoOwner:      DefaultRepoOwner,
		CommitterName:  DefaultCommitterName,
		CommitterEmail: DefaultCommitterEmail,
	}, nil
}

func (g *GithubClient) GetRef(ctx context.Context) (ref *github.Reference, err error) {
	ref, _, err = g.Client.Git.GetRef(ctx, g.RepoOwner, g.Repo, "refs/heads/master")
	if err != nil {
		return ref, err
	}

	return ref, nil
}

func (g *GithubClient) GetHeadFromRef(ctx context.Context, ref *github.Reference) (head *github.RepositoryCommit, err error) {
	head, _, err = g.Client.Repositories.GetCommit(ctx, g.RepoOwner, g.Repo, *ref.Object.SHA)
	if err != nil {
		return nil, err
	}
	head.Commit.SHA = head.SHA

	return head, err
}

func (g *GithubClient) CreateCommit(ctx context.Context, message string, tree *github.Tree, head *github.Commit) (newCommit *github.Commit, err error) {
	date := time.Now()
	author := &github.CommitAuthor{Date: &date, Name: &g.CommitterName, Email: &g.CommitterEmail}
	commit := &github.Commit{Author: author, Message: &message, Tree: tree, Parents: []github.Commit{*head}}
	newCommit, _, err = g.Client.Git.CreateCommit(ctx, g.RepoOwner, g.Repo, commit)
	if err != nil {
		return nil, err
	}
	return
}

func (g *GithubClient) CreateTreeFromEntries(ctx context.Context, ref *github.Reference, entries []github.TreeEntry) (tree *github.Tree, err error) {
	tree, _, err = g.Client.Git.CreateTree(ctx, g.RepoOwner, g.Repo, *ref.Object.SHA, entries)
	return tree, err
}

func (g *GithubClient) UpdateRef(ctx context.Context, ref *github.Reference, newCommit *github.Commit) (err error) {
	ref.Object.SHA = newCommit.SHA
	_, _, err = g.Client.Git.UpdateRef(ctx, g.RepoOwner, g.Repo, ref, false)
	return
}

func (g *GithubClient) GetContents(ctx context.Context, path string) (*github.RepositoryContent, []*github.RepositoryContent, error) {
	fileContent, directoryContent, _, err := g.Client.Repositories.GetContents(ctx, g.RepoOwner, g.Repo, path, nil)
	return fileContent, directoryContent, err
}
