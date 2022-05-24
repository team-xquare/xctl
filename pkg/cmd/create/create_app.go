package create

import (
	"context"
	"fmt"
	"regexp"

	"github.com/spf13/cobra"

	cmdutil "github.com/xctl/cmd/util"
	"github.com/xctl/pkg/api"
	gitops "github.com/xctl/pkg/gitops"
	"github.com/xctl/pkg/gitops/github"
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
	AppTemplatePath    = "template/app-template"
	ValidHostnameRegex = "^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$"
)

type CreateAppOption struct {
	Name          string
	Command       []string
	Type          string
	Host          string
	ImageRegistry string
	ImageTag      string
	ContainerPort int32
	CommitMessage string
	Environment   string

	Gitops gitops.GitopsInterface
}

func NewCmdCreateApp() *cobra.Command {
	o := NewCreateAppOption()
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
	cmd.Flags().StringVar(&o.Host, "host", o.Host, "The host name of service. Default is api.xquare.app")
	cmd.Flags().StringVarP(&o.ImageRegistry, "registry", "r", o.ImageRegistry, "The container registry url. Default is registry.hub.docker.com")
	cmd.Flags().StringVarP(&o.Environment, "environment", "e", o.Environment, "The environment to create an application. Default is staging")
	cmd.Flags().StringVar(&o.ImageTag, "tag", o.ImageTag, "The tag name of a image at start. Default is latest")
	cmd.Flags().Int32Var(&o.ContainerPort, "port", o.ContainerPort, "The port number to run in Docker Container. Default is 8080")

	return cmd
}

func NewCreateAppOption() *CreateAppOption {
	return &CreateAppOption{
		Type:          "backend",
		Host:          "api.xquare.app",
		ImageRegistry: "registry.hub.docker.com",
		ImageTag:      "latest",
		ContainerPort: 8080,
		Environment:   "staging",
	}
}

func (o *CreateAppOption) Complete(cmd *cobra.Command, args []string) error {
	name, err := NameFromCommandArgs(cmd, args)
	if err != nil {
		return err
	}
	o.Name = name
	if len(args) > 1 {
		o.Command = args[1:]
	}

	var client github.GithubClient
	switch o.Environment {
	case "staging":
		client = github.NewGithubClient(github.StagingRepo)
	case "production":
		client = github.NewGithubClient(github.ProductionRepo)
	default:
		return fmt.Errorf("cannot find a specific environment name: %s", o.Environment)
	}
	o.Gitops = gitops.NewGitops(client)

	return nil
}

func (o *CreateAppOption) Validate() error {
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

func (o *CreateAppOption) Run() error {
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

func (o *CreateAppOption) createApplication() (*api.Application, error) {
	app := &api.Application{
		Name:          o.Name,
		Type:          o.Type,
		Host:          o.Host,
		ImageRegistry: o.ImageRegistry,
		ImageTag:      o.ImageTag,
		ContainerPort: o.ContainerPort,
		Environment:   o.Environment,
	}

	return app, nil
}
