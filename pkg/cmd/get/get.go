/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"

	cmdutil "github.com/xctl/pkg/cmd/util"
)

// createCmd represents the create command
func NewCmdGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a resource of xquare cluster",
		Long:  "Get a resource of xquare cluster.",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.RequireNoArguments(cmd, args)
			cmd.Help()
		},
	}

	cmd.AddCommand(NewCmdGetCredential())
	cmd.AddCommand(NewCmdGetApp())

	return cmd
}
