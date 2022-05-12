/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
func NewCmdCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create NAME",
		Short: "Create and deploy app for XQUARE k8s manifest",
		Long:  "Create and deploy app for XQUARE k8s manifest.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("create called")
		},
	}

	cmd.AddCommand(NewCmdCreateApp())

	return cmd
}
