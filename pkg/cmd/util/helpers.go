package util

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
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

func MarshalObject(obj interface{}) string {
	jsonObj, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return ""
	}
	return string(jsonObj)
}

func PrintObject(obj interface{}) {
	jsonObj := MarshalObject(obj)
	fmt.Fprintf(os.Stdout, "%s\n", jsonObj)
}

func GetHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		home = strings.Replace(home, "\\", "/", -1)
		return home
	} else if runtime.GOOS == "linux" {
		home := os.Getenv("XDG_CONFIG_HOME")
		if home != "" {
			return home
		}
	}
	return os.Getenv("HOME")
}
