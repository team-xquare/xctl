package create

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/xctl/pkg/api"
	cmdutil "github.com/xctl/pkg/cmd/util"
	gitops "github.com/xctl/pkg/gitops"
	"github.com/xctl/pkg/gitops/github"
)

var (
	cmdDescription = `
	Delete an application with the specified name.`
	cmdExample = `
	#Delete an application named testapp on staging environment
	xctl delete app app-test-staging
	
	#Delete an application named testapp on production environment
	xctl delete app app-test-production 
	`
)

type DeleteAppOptions struct {
	Name        string
	Command     []string
	Environment string
	environment string
	Type        string

	Gitops gitops.GitopsInterface
}

func NewDeleteAppOptions() *DeleteAppOptions {
	return &DeleteAppOptions{
		Environment: "staging",
		Type:        "backend",
	}
}

func NewCmdDeleteApp() *cobra.Command {
	o := NewDeleteAppOptions()
	cmd := &cobra.Command{
		Use:     "app NAME",
		Short:   cmdDescription,
		Long:    cmdDescription,
		Example: cmdExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.Run())
		},
	}

	cmd.Flags().StringVarP(&o.Environment, "environment", "e", o.Environment, "The environment to create an application. Default is staging")
	cmd.Flags().StringVarP(&o.Type, "type", "t", o.Type, "The type of service. default is backend")

	return cmd
}

func (o *DeleteAppOptions) Complete(cmd *cobra.Command, args []string) error {
	name, err := NameFromCommandArgs(cmd, args)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(name, "app-") {
		name = "app-" + name
	}

	o.Type, err = api.CheckApplicationType(o.Type)
	if !strings.HasSuffix(name, "-"+o.Type) {
		name = name + "-" + o.Type
	}

	o.Name = name
	if len(args) > 1 {
		o.Command = args[1:]
	}

	if err != nil {
		return err
	}

	o.Environment, err = api.CheckApplicationEnvironment(o.Environment)
	if err != nil {
		return err
	}

	client, err := github.NewGithubClient(o.Environment)
	if err != nil {
		return err
	}
	o.Gitops = gitops.NewGitops(client)

	return nil
}

func (o *DeleteAppOptions) Validate() error {
	return nil
}

func (o *DeleteAppOptions) Run() error {
	apps, err := o.Gitops.Application(nil).Get(context.Background())
	if err != nil {
		return err
	}

	for _, app := range apps {
		if fmt.Sprintf("app-%s-%s", app.Name, app.Type) == o.Name {
			err = o.Gitops.Application(app).Delete(context.Background())
			if err != nil {
				return err
			}
			cmdutil.PrintObject(app)
			return nil
		}
	}

	return fmt.Errorf("cannot find a specific application name: %s", o.Name)
}
