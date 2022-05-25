/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"

	cmdutil "github.com/xctl/pkg/cmd/util"
)

// createCmd represents the create command
func NewCmdSet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set a credential information of xctl",
		Long:  "Set a credential information of xctl.",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.RequireNoArguments(cmd, args)
			cmd.Help()
		},
	}

	cmd.AddCommand(NewCmdSetCredential())

	return cmd
}
