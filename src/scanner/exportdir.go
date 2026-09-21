package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

// ExpandPath resolves the shorthands people write a path with.
//
// `~` and environment variables are both how a person writes a directory they
// mean, and a configured value and a typed one deserve the same treatment. The
// scan directories already expand environment variables; this adds the tilde,
// for the one field a person types interactively.
func ExpandPath(p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return ""
	}
	p = os.ExpandEnv(p)
	if p == "~" || strings.HasPrefix(p, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			p = filepath.Join(home, strings.TrimPrefix(strings.TrimPrefix(p, "~"), "/"))
		}
	}
	return filepath.Clean(p)
}

// DefaultExportDir is where an export goes when nothing is configured: the home
// directory, which is where every export went before it could be chosen.
func DefaultExportDir(cfg *Config) string {
	if cfg != nil && strings.TrimSpace(cfg.ExportDir) != "" {
		return ExpandPath(cfg.ExportDir)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return home
}

// ExportDirProblem is a destination that cannot be used.
type ExportDirProblem struct {
	Detail string
}

func (p *ExportDirProblem) Error() string {
	if p == nil {
		return ""
	}
	return p.Detail
}

// CheckExportDir reports why a directory cannot be exported into, or nil.
//
// Nothing is created: a directory that does not exist is a typo far more often
// than it is an intention, and creating it silently turns one into the other.
func CheckExportDir(dir string) *ExportDirProblem {
	if strings.TrimSpace(dir) == "" {
		return &ExportDirProblem{Detail: "no directory given"}
	}
	info, err := os.Stat(dir)
	if err != nil {
		return &ExportDirProblem{Detail: dir + " does not exist"}
	}
	if !info.IsDir() {
		return &ExportDirProblem{Detail: dir + " is not a directory"}
	}
	return nil
}

// CompleteDir offers the directories that extend a partly typed path.
//
// The text input accepts suggestions and binds `tab` to them, so completing a
// path is a matter of handing it what is there rather than building anything.
func CompleteDir(typed string) []string {
	expanded := ExpandPath(typed)
	dir, prefix := expanded, ""
	if !strings.HasSuffix(typed, "/") {
		dir, prefix = filepath.Split(expanded)
	}
	entries, err := os.ReadDir(filepath.Clean(dir))
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		if !e.IsDir() || !strings.HasPrefix(e.Name(), prefix) {
			continue
		}
		out = append(out, filepath.Join(filepath.Clean(dir), e.Name()))
	}
	return out
}
