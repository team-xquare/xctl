package create

import (
	"context"
	"fmt"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/xctl/pkg/api"
	cmdutil "github.com/xctl/pkg/cmd/util"
	gitops "github.com/xctl/pkg/gitops"
	"github.com/xctl/pkg/gitops/github"
)

var (
	cmdDescription = `
	Create and deploy an application with the specified name.`
	cmdExample = `
	#Create and deploy an application named testapp
	xctl create app testapp
	
	#Create and deploy an application named testapp whose container exposed port 3030
	xctl create app testapp --containerPort 3030

	#Create and deploy an application named testapp to production environment
	xctl create app testapp --production
	`
	AppTemplatePath    = "template/app-template"
	ValidHostnameRegex = "^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$"
)

type CreateAppOptions struct {
	Name          string
	Command       []string
	Type          string
	Host          string
	ImageRegistry string
	ImageTag      string
	ContainerPort int32
	CommitMessage string
	Environment   string
	Prefix        string

	Gitops gitops.GitopsInterface
}

func NewCreateAppOption() *CreateAppOptions {
	return &CreateAppOptions{
		Type:          "backend",
		Host:          "api.xquare.app",
		ImageRegistry: "registry.hub.docker.com",
		ImageTag:      "latest",
		/* 도커 파일의 외부 포트와 연결되는 포트.
		예를 들어 리액트 애플리케이션을 배포하는 경우, 리액트의 기본 포트인 3000을 지정해주면 된다.*/
		ContainerPort: 8080,
		Environment:   "staging",
		Prefix:        "/",
	}
}

func NewCmdCreateApp() *cobra.Command {
	o := NewCreateAppOption()
	cmd := &cobra.Command{
		Use:     "app NAME -- [COMMAND] [args...]",
		Short:   cmdDescription,
		Long:    cmdDescription,
		Example: cmdExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.Run())
		},
	}

	cmd.Flags().StringVarP(&o.Type, "type", "t", o.Type, "The type of service. default is backend")
	cmd.Flags().StringVar(&o.Host, "host", o.Host, "The host name of service. Default is api.xquare.app")
	cmd.Flags().StringVarP(&o.ImageRegistry, "registry", "r", o.ImageRegistry, "The container registry url. Default is registry.hub.docker.com")
	cmd.Flags().StringVarP(&o.Environment, "environment", "e", o.Environment, "The environment to create an application. Default is staging")
	cmd.Flags().StringVar(&o.ImageTag, "tag", o.ImageTag, "The tag name of a image at start. Default is latest")
	cmd.Flags().StringVarP(&o.Prefix, "prefix", "p", o.Prefix, "The prefix for service routing. Default is /")
	cmd.Flags().Int32Var(&o.ContainerPort, "port", o.ContainerPort, "The port number to run in Docker Container. Default is 8080")

	return cmd
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

	var repo string
	switch o.Environment {
	case "staging":
		repo = github.StagingRepo
	case "production":
		repo = github.ProductionRepo
	default:
		return fmt.Errorf("cannot find a specific environment name: %s", o.Environment)
	}
	client, err := github.NewGithubClient(repo)
	if err != nil {
		return err
	}
	o.Gitops = gitops.NewGitops(client)

	return nil
}

func (o *CreateAppOptions) Validate() error {
	if o.ContainerPort < 0 || o.ContainerPort > 65535 {
		return fmt.Errorf("port number cannot out of range 0 to 65535")
	}

	matched, err := regexp.MatchString(ValidHostnameRegex, o.Host)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("require valid host name (%s)", o.Host)
	}

	return nil
}

func (o *CreateAppOptions) Run() error {
	app, err := o.createApplication()
	if err != nil {
		return err
	}
	err = o.Gitops.Application().Create(context.Background(), app)
	if err != nil {
		return err
	}

	cmdutil.PrintObject(app)
	return nil
}

func (o *CreateAppOptions) createApplication() (*api.Application, error) {
	app := &api.Application{
		Name:          o.Name,
		Type:          o.Type,
		Host:          o.Host,
		ImageRegistry: o.ImageRegistry,
		ImageTag:      o.ImageTag,
		ContainerPort: o.ContainerPort,
		Environment:   o.Environment,
		Prefix:        o.Prefix,
	}

	return app, nil
}
