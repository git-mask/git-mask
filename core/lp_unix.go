//go:build unix

package core

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// ErrNotFound is the error resulting if a path search failed to find an executable file.
var ErrNotFound = errors.New("executable file not found in $PATH")

func findExecutable(file string) error {
	d, err := os.Stat(file)
	if err != nil {
		return err
	}
	m := d.Mode()
	if m.IsDir() {
		return syscall.EISDIR
	}
	if m&0111 != 0 {
		return nil
	}
	return fs.ErrPermission
}

// LookPath searches for an executable named file in the
// directories named by the PATH environment variable.
// If file contains a slash, it is tried directly and the PATH is not consulted.
// Otherwise, on success, the result is an absolute path.
//
// In older versions of Go, LookPath could return a path relative to the current directory.
// As of Go 1.19, LookPath will instead return that path along with an error satisfying
// errors.Is(err, ErrDot). See the package documentation for more details.
func LookPath(file string) ([]string, error) {
	// NOTE(rsc): I wish we could use the Plan 9 behavior here
	// (only bypass the path if file begins with / or ./ or ../)
	// but that would not match all the Unix shells.

	if strings.Contains(file, "/") {
		err := findExecutable(file)
		if err == nil {
			return []string{file}, nil
		}
		return nil, &LookPathError{file, err}
	}
	var paths []string
	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		if dir == "" {
			// Unix shell semantics: path element "" means "."
			dir = "."
		}
		path := filepath.Join(dir, file)
		if err := findExecutable(path); err == nil {
			paths = append(paths, path)
		}
	}
	if len(paths) > 0 {
		return paths, nil
	}
	return nil, &LookPathError{file, ErrNotFound}
}

func IsExecutable(path string) bool {
	return findExecutable(path) == nil
}
