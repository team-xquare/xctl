/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"github.com/xctl/cmd"
	"github.com/xctl/config"
)

func init() {
	config.MustLoadConfig()
}

func main() {
	cmd.Execute()
}
