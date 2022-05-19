/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/google/go-github/v44/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	cmdutil "github.com/xctl/cmd/util"
	"github.com/xctl/config"
)

var (
	appDescription = `
	Create and deploy an application with the specified name.`
	appExample = `
	#Create and deploy an application named testapp
	xctl create app testapp
	
	#Create and deploy an application named testapp whose container exposed port 3030
	xctl create app testapp --containerPort 3030

	#Create and deploy an application named testapp to production environment
	xctl create app testapp --production
	`
	sourceOwner          = "team-xquare"
	defaultSourceRepo    = stagingSourceRepo
	stagingSourceRepo    = "xquare-gitops-repo-staging"
	productionSourceRepo = "xquare-gitops-repo-production"
	authorName           = "xquare-admin"
	authorEmail          = "teamxquare@gmail.com"
	charTemplatePath     = "template/chart-template/"
	appTemplatePath      = "template/app-template/"
	validHostnameRegex   = "^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$"
)

var (
	fileType = github.String("blob")
	readMode = github.String("100644")
)

type CreateAppOptions struct {
	Name          string
	Command       []string
	Type          string
	Host          string
	ImageRegistry string
	ImageTag      string
	ContainerPort int32
	Staging       bool
	Production    bool

	client      *github.Client
	ref         *github.Reference
	tree        *github.Tree
	head        *github.RepositoryCommit
	sourceRepos []string
}

func NewCmdCreateApp() *cobra.Command {
	o := NewCreateAppOptions()
	cmd := &cobra.Command{
		Use:     "app NAME -- [COMMAND] [args...]",
		Short:   appDescription,
		Long:    appDescription,
		Example: appExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.Run())
		},
	}

	cmd.Flags().StringVarP(&o.Type, "type", "t", o.Type, "The type of service. default is backend")
	cmd.Flags().BoolVarP(&o.Staging, "staging", "s", o.Staging, "If set this flag, application will deploy to staging environment.")
	cmd.Flags().BoolVarP(&o.Production, "production", "p", o.Production, "If set this flag, application will deploy to production environment.")
	cmd.Flags().StringVar(&o.Host, "host", o.Host, "The host name of service. Default is api.xquare.app")
	cmd.Flags().StringVarP(&o.ImageRegistry, "registry", "r", o.ImageRegistry, "The container registry url. Default is registry.hub.docker.com")
	cmd.Flags().StringVar(&o.ImageTag, "tag", o.ImageTag, "The tag name of a image at start. Default is latest")
	cmd.Flags().Int32Var(&o.ContainerPort, "port", o.ContainerPort, "The port number to run in Docker Container. Default is 8080")

	return cmd
}

func NewCreateAppOptions() *CreateAppOptions {
	return &CreateAppOptions{
		Type:          "backend",
		Host:          "api.xquare.app",
		ImageRegistry: "registry.hub.docker.com",
		ImageTag:      "latest",
		ContainerPort: 8080,
	}
}

func (o *CreateAppOptions) Complete(cmd *cobra.Command, args []string) error {
	name, err := NameFromCommandArgs(cmd, args)
	if err != nil {
		return err
	}
	o.Name = name
	if len(args) > 1 {
		o.Command = args[1:]
	}

	if o.Staging {
		o.sourceRepos = append(o.sourceRepos, stagingSourceRepo)
	} else if o.Production {
		o.sourceRepos = append(o.sourceRepos, productionSourceRepo)
	} else {
		o.sourceRepos = append(o.sourceRepos, defaultSourceRepo)
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: config.Config.GithubToken})
	tc := oauth2.NewClient(context.Background(), ts)
	o.client = github.NewClient(tc)

	return nil
}

func (o *CreateAppOptions) Validate() error {
	if o.ContainerPort < 0 || o.ContainerPort > 65535 {
		return fmt.Errorf("port number cannot out of range 0 to 65535")
	}

	matched, err := regexp.MatchString(validHostnameRegex, o.Host)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("require valid host name (%s)", o.Host)
	}

	return nil
}

func (o *CreateAppOptions) Run() error {
	for _, sourceRepo := range o.sourceRepos {
		newCommit, err := o.createAppToRepo(sourceRepo)
		if err != nil {
			return err
		}

		err = cmdutil.PrintObject(newCommit)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *CreateAppOptions) createAppToRepo(sourceRepo string) (newCommit *github.Commit, err error) {
	ctx := context.Background()

	o.ref, err = o.getRef(ctx, sourceRepo)
	if err != nil {
		return nil, err
	}
	o.tree, err = o.createTree(ctx, sourceRepo)
	if err != nil {
		return nil, err
	}
	o.head, err = o.getHeadFromRef(ctx, sourceRepo)
	if err != nil {
		return nil, err
	}
	newCommit, err = o.createCommit(ctx, sourceRepo)
	if err != nil {
		return nil, err
	}
	if err = o.pushCommit(ctx, newCommit, sourceRepo); err != nil {
		return nil, err
	}

	return newCommit, nil
}

func (o *CreateAppOptions) getRef(ctx context.Context, sourceRepo string) (ref *github.Reference, err error) {
	ref, _, err = o.client.Git.GetRef(ctx, sourceOwner, sourceRepo, "refs/heads/master")
	if err != nil {
		return ref, err
	}

	return ref, nil
}

func (o *CreateAppOptions) createTree(ctx context.Context, sourceRepo string) (tree *github.Tree, err error) {
	entries := []*github.TreeEntry{}

	var buf bytes.Buffer
	err = filepath.Walk(charTemplatePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			var paths []string
			if runtime.GOOS == "windows" {
				path = strings.Replace(path, "\\", "/", -1)
			}
			paths = strings.Split(path, "/")
			fileName := strings.Join(paths[2:], "/")

			tmpl, err := template.ParseGlob(fmt.Sprintf("%s/%s", charTemplatePath, fileName))
			if err != nil {
				return err
			}
			err = tmpl.Execute(&buf, &o)
			if err != nil {
				return err
			}

			content := github.String(buf.String())
			buf.Reset()
			entries = append(entries, &github.TreeEntry{
				Path:    github.String(fmt.Sprintf("charts/%s/%s/%s", o.Type, o.Name, fileName)),
				Type:    fileType,
				Content: content,
				Mode:    readMode,
			})
			return nil
		})
	if err != nil {
		return nil, err
	}

	tmpl, err := template.ParseFiles(appTemplatePath + "app-template.yaml")
	if err != nil {
		return nil, err
	}
	err = tmpl.Execute(&buf, &o)
	if err != nil {
		return nil, err
	}
	content := github.String(buf.String())
	entries = append(entries, &github.TreeEntry{
		Path:    github.String(fmt.Sprintf("apps/templates/%s/%s", o.Type, fmt.Sprintf("app-%s-%s.yaml", o.Name, o.Type))),
		Type:    fileType,
		Content: content,
		Mode:    readMode,
	})

	tree, _, err = o.client.Git.CreateTree(ctx, sourceOwner, sourceRepo, *o.ref.Object.SHA, entries)
	return tree, err
}

func (o *CreateAppOptions) getHeadFromRef(ctx context.Context, sourceRepo string) (head *github.RepositoryCommit, err error) {
	head, _, err = o.client.Repositories.GetCommit(ctx, sourceOwner, sourceRepo, *o.ref.Object.SHA, nil)
	if err != nil {
		return nil, err
	}
	head.Commit.SHA = head.SHA

	return head, err
}

func (o *CreateAppOptions) createCommit(ctx context.Context, sourceRepo string) (newCommit *github.Commit, err error) {
	date := time.Now()
	commitMessage := fmt.Sprintf("⚡️ :: %s 서비스 추가", o.Name+"-"+o.Type)
	author := &github.CommitAuthor{Date: &date, Name: &authorName, Email: &authorEmail}
	commit := &github.Commit{Author: author, Message: &commitMessage, Tree: o.tree, Parents: []*github.Commit{o.head.Commit}}
	newCommit, _, err = o.client.Git.CreateCommit(ctx, sourceOwner, sourceRepo, commit)
	if err != nil {
		return nil, err
	}
	return
}

func (o *CreateAppOptions) pushCommit(ctx context.Context, newCommit *github.Commit, sourceRepo string) (err error) {
	o.ref.Object.SHA = newCommit.SHA
	_, _, err = o.client.Git.UpdateRef(ctx, sourceOwner, sourceRepo, o.ref, false)
	return
}
