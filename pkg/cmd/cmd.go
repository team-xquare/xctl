/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	create "github.com/xctl/pkg/cmd/create"
	get "github.com/xctl/pkg/cmd/get"
	set "github.com/xctl/pkg/cmd/set"
	cmdutil "github.com/xctl/pkg/cmd/util"
)

// rootCmd represents the base command when called without any subcommands
func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "xctl",
		Short: "새로운 앱을 쉽게 생성하고 추가해주는 CLI",
		Long: ` xctl는 gitops repository에 새로운 프로젝트를 추가하거나
생성된 프로젝트를 관리할 수 있는 커맨드 라인 애플리케이션이다`,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.RequireNoArguments(cmd, args)
			cmd.Help()
		},
	}

	cmd.AddCommand(create.NewCmdCreate())
	cmd.AddCommand(get.NewCmdGet())
	cmd.AddCommand(set.NewCmdSet())

	cmd.CompletionOptions.DisableDefaultCmd = true

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cmd := NewCmdRoot()

	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
