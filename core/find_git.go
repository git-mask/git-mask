package core

import (
	"os"
	"strings"
)

func FindGit() (string, error) {
	// check $GIT_MASK_REAL_GIT
	if git, ok := os.LookupEnv("GIT_MASK_REAL_GIT"); ok && IsExecutable(git) {
		return git, nil
	}
	self, err := os.Executable()
	if err != nil {
		return "", err
	}
	paths, err := LookPath("git")
	if err != nil {
		return "", err
	}
	for _, path := range paths {
		if strings.EqualFold(path, self) {
			continue
		}
		if IsExecutable(path) {
			return path, nil
		}
	}
	return "", ErrNotFound
}
