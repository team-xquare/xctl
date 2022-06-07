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

func (o *GetAppOptions) Complete(cmd *cobra.Command, args []string) (err error) {
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

func (o *GetAppOptions) Validate() error {
	return nil
}

func (o *GetAppOptions) Run() error {
	fmt.Printf("Environemnt: %s\n", o.Environment)
	frontend, err := o.Gitops.Application(nil).GetFrontend(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("frontend applications")
	o.printApplications(frontend)

	backend, err := o.Gitops.Application(nil).GetBackend(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("backend applications")
	o.printApplications(backend)
	return nil
}

func (o *GetAppOptions) printApplications(apps []*api.Application) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "\t Name\t Base Url\t Image Version\t")
	for _, app := range apps {
		var host string
		if app.Type == api.Backend {
			host = "api.xquare.app"
		} else if app.Type == api.Frontend {
			host = "webview.xquare.app"
		}
		if app.Environment != api.Production {
			host = app.Environment + "-" + host
		}
		fmt.Fprintf(w, "\t app-%s-%s\t %s\t %s\t\n",
			app.Name,
			app.Type,
			host+app.Prefix,
			app.ImageTag)
	}
	w.Flush()
}
