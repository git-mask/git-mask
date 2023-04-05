package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var (
	GIT_DIR string
)

// Entrypoint for git-mask running in Wrapper Mode, i.e., git-mask is aliasing
// the real git.
func ExecuteWrapper() {
	// we want to find out the top level command to execute. we need to strip
	// global git options as in https://github.com/git/git/blob/master/git.c
	args := os.Args[1:]
	if len(args) == 0 {
		runGit(nil, nil)
	}
	var flags []string
	var cmd *string
	for len(args) > 0 {
		if !strings.HasPrefix(args[0], "-") {
			cmd = &args[0]
			args = args[1:]
			break
		}
		var flag string
		flag, args = args[0], args[1:]
		flags = append(flags, flag)
		if path, found := strings.CutPrefix(flag, "--git-dir="); found {
			setGitDir(path)
		} else if flag == "--git-dir" {
			if len(args) > 0 {
				setGitDir(args[0])
				flags = append(flags, args[0])
				args = args[1:]
			}
		} else if flag == "--bare" {
			pwd, _ := os.Getwd()
			setGitDir(pwd)
		} else if flag == "-C" {
			if len(args) > 0 {
				err := os.Chdir(args[0])
				cobra.CheckErr(err)
				flags = append(flags, args[0])
				args = args[1:]
			}
		} else if flag == "--exec-path" ||
			flag == "--namespace" ||
			flag == "--work-tree" ||
			flag == "--bare" ||
			flag == "-c" ||
			flag == "--config-env" ||
			flag == "--shallow-file" {
			if len(args) > 0 {
				flags = append(flags, args[0])
				args = args[1:]
			}
		}
	}
	runGit(flags, cmd, args...)
}

func setGitDir(val string) {
	GIT_DIR = val
	os.Setenv("GIT_DIR", val)
}

func runGit(flags []string, cmd *string, args ...string) {
	if cmd != nil {
		args = append(append(flags, *cmd), args...)
	} else {
		args = append(flags, args...)
	}
	fmt.Printf("> %s %s\n", REAL_GIT, strings.Join(args, " "))
	c := exec.Command(REAL_GIT, args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	var exitError *exec.ExitError
	if err == nil {
		os.Exit(0)
	} else if errors.As(err, &exitError) {
		os.Exit(exitError.ExitCode())
	}
}
