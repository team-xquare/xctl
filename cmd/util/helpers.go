package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const (
	DefaultErrorExitCode = 1
)

func RequireNoArguments(c *cobra.Command, args []string) {
	if len(args) > 0 {
		fatal(UsageErrorf(c, "error: unknown command %q", strings.Join(args, " ")))
	}
}

func UsageErrorf(cmd *cobra.Command, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s \nSee '%s -h' for help and examples", msg, cmd.CommandPath())
}

func fatal(err error) {
	fmt.Fprint(os.Stderr, err.Error())
	os.Exit(DefaultErrorExitCode)
}
