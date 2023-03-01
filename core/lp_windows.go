package core

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ErrNotFound is the error resulting if a path search failed to find an executable file.
var ErrNotFound = errors.New("executable file not found in %PATH%")

func chkStat(file string) error {
	d, err := os.Stat(file)
	if err != nil {
		return err
	}
	if d.IsDir() {
		return fs.ErrPermission
	}
	return nil
}

func hasExt(file string) bool {
	i := strings.LastIndex(file, ".")
	if i < 0 {
		return false
	}
	return strings.LastIndexAny(file, `:\/`) < i
}

func findExecutable(file string, exts []string) (string, error) {
	if len(exts) == 0 {
		return file, chkStat(file)
	}
	if hasExt(file) {
		if chkStat(file) == nil {
			return file, nil
		}
	}
	for _, e := range exts {
		if f := file + e; chkStat(f) == nil {
			return f, nil
		}
	}
	return "", fs.ErrNotExist
}

// LookPath searches for an executable named file in the
// directories named by the PATH environment variable.
// LookPath also uses PATHEXT environment variable to match
// a suitable candidate.
// If file contains a slash, it is tried directly and the PATH is not consulted.
// Otherwise, on success, the result is an absolute path.
func LookPath(file string) ([]string, error) {
	var exts []string
	var paths []string
	x := os.Getenv(`PATHEXT`)
	if x != "" {
		for _, e := range strings.Split(strings.ToLower(x), `;`) {
			if e == "" {
				continue
			}
			if e[0] != '.' {
				e = "." + e
			}
			exts = append(exts, e)
		}
	} else {
		exts = []string{".com", ".exe", ".bat", ".cmd"}
	}

	if strings.ContainsAny(file, `:\/`) {
		f, err := findExecutable(file, exts)
		if err == nil {
			return []string{f}, nil
		}
		return nil, &LookPathError{file, err}
	}

	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		if f, err := findExecutable(filepath.Join(dir, file), exts); err == nil {
			paths = append(paths, f)
		} else if !errors.Is(err, fs.ErrNotExist) {
			return nil, &LookPathError{file, err}
		}
	}

	if len(paths) > 0 {
		return paths, nil
	}

	return nil, &LookPathError{file, ErrNotFound}
}

func IsExecutable(path string) bool {
	return chkStat(path) == nil
}
