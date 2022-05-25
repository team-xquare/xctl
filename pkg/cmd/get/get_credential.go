package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cmdutil "github.com/xctl/pkg/cmd/util"
	"github.com/xctl/pkg/gitops/config"
)

var (
	cmdDescription = `
	Get a credential of xctl`
	cmdExample = `
	#Get a credential
	xctl get credential
	`
)

type GetCredentialOptions struct {
	Command []string
}

func NewGetCredentialOptions() *GetCredentialOptions {
	return &GetCredentialOptions{}
}

func NewCmdGetCredential() *cobra.Command {
	o := NewGetCredentialOptions()
	cmd := &cobra.Command{
		Use:     "credential",
		Short:   cmdDescription,
		Long:    cmdDescription,
		Example: cmdExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.Run())
		},
	}

	return cmd
}

func (o *GetCredentialOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

func (o *GetCredentialOptions) Validate() error {
	if _, err := os.Stat(config.CredentialPath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("cannot find a credential, please try this command:\n xctl set credential <githubToken>")
	}
	return nil
}

func (o *GetCredentialOptions) Run() error {
	credential, err := config.GetCredential()
	if err != nil {
		return err
	}

	cmdutil.PrintObject(credential)
	return nil
}
