/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	appLong    = "Create a XQUARE application"
	appExample = ""
)

type CreateAppOptions struct {
	Name          string
	ServiceType   ServiceType
	Environment   Envrionment
	ImageRegistry string
	Host          string
	ContainerPort int32
}

func NewCreateAppOptions() *CreateAppOptions {
	return &CreateAppOptions{
		ServiceType:   BACKEND,
		Environment:   STAGING,
		Host:          "default.xquare.app",
		ImageRegistry: "registry.hub.docker.com",
		ContainerPort: 8080,
	}
}

func (o *CreateAppOptions) Validate() error {
	if len(o.Name) < 1 {
		return fmt.Errorf("Error on validating")
	}
	return nil
}

func (o *CreateAppOptions) Run() error {
	app := 	
}


type ServiceType int32

const (
	BACKEND ServiceType = iota + 1
	FRONTEND
)

type Envrionment int32

const (
	PRODUCTION Envrionment = iota + 1
	STAGING
)

// createAppCmd represents the createApp command
func NewCmdCreateApp() *cobra.Command {
	o := NewCreateAppOptions()
	cmd := &cobra.Command{
		Use:     "app NAME",
		Short:   "Create and deploy app for XQUARE k8s manifest",
		Long:    "Create and deploy app for XQUARE k8s manifest.",
		Example: appExample,
		Run: func(cmd *cobra.Command, args []string) {
			o.Run()
		},
	}

	cmd.Flags().Int32Var(&o.ContainerPort, "containerPort", o.ContainerPort, "port number to run in Docker Container")

	return cmd
}
