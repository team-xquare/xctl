package create

import (
	"github.com/spf13/cobra"

	cmdutil "github.com/xctl/cmd/util"
)

// createCmd represents the create command
func NewCmdCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a resource of xquare cluster",
		Long:  "Create a resource of xquare cluster.",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.RequireNoArguments(cmd, args)
			cmd.Help()
		},
	}

	cmd.AddCommand(NewCmdCreateApp())

	return cmd
}

func NameFromCommandArgs(cmd *cobra.Command, args []string) (string, error) {
	argsLen := cmd.ArgsLenAtDash()
	if argsLen == -1 {
		argsLen = len(args)
	}
	if argsLen != 1 {
		return "", cmdutil.UsageErrorf(cmd, "exactly one NAME is required, got %d", argsLen)
	}
	return args[0], nil
}
