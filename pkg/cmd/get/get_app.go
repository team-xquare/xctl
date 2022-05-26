package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/xctl/pkg/api"
	cmdutil "github.com/xctl/pkg/cmd/util"
	"github.com/xctl/pkg/gitops"
	"github.com/xctl/pkg/gitops/github"
)

var (
	appDescription = `
	Get application list from the specified environment`
	appExample = `
	#Get application list from staging environment
	xctl get app -e stag or xctl get app -e staging
	`
)

type GetAppOptions struct {
	Command     []string
	Environment string

	Gitops gitops.GitopsInterface
}

func NewGetAppOptions() *GetAppOptions {
	return &GetAppOptions{
		Environment: "staging",
	}
}

func NewCmdGetApp() *cobra.Command {
	o := NewGetAppOptions()
	cmd := &cobra.Command{
		Use:     "app",
		Short:   appDescription,
		Long:    appDescription,
		Example: appExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.Run())
		},
	}
	cmd.Flags().StringVarP(&o.Environment, "environment", "e", o.Environment, "The environment to get the resouces from. Default is staging")

	return cmd
}

func (o *GetAppOptions) Complete(cmd *cobra.Command, args []string) error {
	var repo string
	if o.Environment == "stag" {
		o.Environment = "staging"
	}
	if o.Environment == "prod" {
		o.Environment = "production"
	}

	switch o.Environment {
	case "staging":
		repo = github.StagingRepo
	case "production":
		repo = github.StagingRepo
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

func (o *GetAppOptions) Validate() error {
	return nil
}

func (o *GetAppOptions) Run() error {
	frontend, backend, err := o.Gitops.Application().Get(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("Environemnt: %s\n", o.Environment)

	fmt.Println("frontend applications")
	o.printApplications(frontend)

	fmt.Println("backend applications")
	o.printApplications(backend)
	return nil
}

func (o *GetAppOptions) printApplications(apps []*api.Application) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "\t Name\t Base Url\t Image Version\t")
	for _, app := range apps {
		fmt.Fprintf(w, "\t app-%s-%s\t %s\t %s\t\n",
			app.Name,
			app.Type,
			app.Host+app.Prefix,
			app.ImageTag)
	}
	w.Flush()
}
