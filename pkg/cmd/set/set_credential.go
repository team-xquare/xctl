package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	cmdutil "github.com/xctl/pkg/cmd/util"
	"github.com/xctl/pkg/gitops/config"
)

var (
	cmdDescription = `
	Set a credential of xctl`
	cmdExample = `
	#Set a credential
	xctl set credential example-token-awefivdfv1fd
	`
)

type SetCredentialOptions struct {
	GithubToken string
}

func NewGetCredentialOptions() *SetCredentialOptions {
	return &SetCredentialOptions{}
}

func NewCmdSetCredential() *cobra.Command {
	o := NewGetCredentialOptions()
	cmd := &cobra.Command{
		Use:     "credential TOKEN",
		Short:   cmdDescription,
		Long:    cmdDescription,
		Example: cmdExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.Run())
		},
	}

	cmd.Flags().StringVarP(&o.GithubToken, "token", "t", o.GithubToken, "The api token of github account")

	return cmd
}

func (o *SetCredentialOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

func (o *SetCredentialOptions) Validate() error {
	if len(o.GithubToken) == 0 {
		return fmt.Errorf("you missed required token value, please specify -t or --token option")
	}
	return nil
}

func (o *SetCredentialOptions) Run() error {
	credential := o.createCredential()
	err := config.SetCredential(credential)
	if err != nil {
		return err
	}

	cmdutil.PrintObject(credential)

	return nil
}

func (o *SetCredentialOptions) createCredential() *config.Credential {
	credential := &config.Credential{GithubToken: o.GithubToken}
	return credential
}
