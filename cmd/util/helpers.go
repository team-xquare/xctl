package util

import (
	"encoding/json"
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
		CheckErr(UsageErrorf(c, "error: unknown command %q", strings.Join(args, " ")))
	}
}

func UsageErrorf(cmd *cobra.Command, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s \nSee '%s -h' for help and examples", msg, cmd.CommandPath())
}

func CheckErr(err error) {
	if err == nil {
		return
	}

	msg := err.Error()
	if !strings.HasPrefix(msg, "error: ") {
		msg = fmt.Sprintf("error: %s", msg)
	}

	fatal(msg, DefaultErrorExitCode)
}

func fatal(msg string, code int) {
	if len(msg) > 0 {
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		fmt.Fprint(os.Stderr, msg)
	}
	os.Exit(code)
}

func PrintObject(obj interface{}) (err error) {
	jsonObj, err := json.MarshalIndent(obj, "", " ")
	if err != nil {
		return
	}
	fmt.Fprintf(os.Stdout, "%s\n", string(jsonObj))
	return
}
