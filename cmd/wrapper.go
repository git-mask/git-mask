package cmd

import (
	"os"
	"strings"
)

// Entrypoint for git-mask running in Wrapper Mode, i.e., git-mask is aliasing
// the real git.
func ExecuteWrapper() {
	// handle global options as in https://github.com/git/git/blob/master/git.c
	args := os.Args[1:]

	execPath, _ := os.Getwd()

	for len(args) > 0 {
		cmd := args[0]
		if !strings.HasPrefix(cmd, "-") {
			break
		}
		if after, found := strings.CutPrefix(cmd, "--exec-path"); found {
			if path, found := strings.CutPrefix(after, "="); found {
				execPath = path
			}
		}
	}
}
