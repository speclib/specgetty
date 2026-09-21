package scanner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// --- 1.2 the shorthands ---

func TestExpandPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SOMEWHERE", "/srv/exports")

	cases := []struct{ in, want string }{
		{"~", home},
		{"~/Downloads", filepath.Join(home, "Downloads")},
		{"$SOMEWHERE", "/srv/exports"},
		{"$SOMEWHERE/zips", "/srv/exports/zips"},
		{"  ~/Downloads  ", filepath.Join(home, "Downloads")},
		{"/absolute/path", "/absolute/path"},
		{"", ""},
	}
	for _, tc := range cases {
		if got := ExpandPath(tc.in); got != tc.want {
			t.Errorf("ExpandPath(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestExpandPathLeavesATildeInside(t *testing.T) {
	// Only a leading tilde is a home directory. One in the middle is a
	// character in a path.
	if got := ExpandPath("/srv/a~b"); got != "/srv/a~b" {
		t.Errorf("got %q", got)
	}
}

// --- 1.1 the configured default ---

func TestDefaultExportDir(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	t.Run("configured", func(t *testing.T) {
		cfg := &Config{ExportDir: "~/Downloads"}
		if got := DefaultExportDir(cfg); got != filepath.Join(home, "Downloads") {
			t.Errorf("got %q", got)
		}
	})
	t.Run("not configured", func(t *testing.T) {
		if got := DefaultExportDir(&Config{}); got != home {
			t.Errorf("got %q, want the home directory", got)
		}
	})
	t.Run("configured blank", func(t *testing.T) {
		if got := DefaultExportDir(&Config{ExportDir: "   "}); got != home {
			t.Errorf("got %q, want the home directory", got)
		}
	})
	t.Run("no config at all", func(t *testing.T) {
		if got := DefaultExportDir(nil); got != home {
			t.Errorf("got %q, want the home directory", got)
		}
	})
}

// --- 3.1 to 3.3 refusing what cannot be used ---

func TestCheckExportDir(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "a-file")
	if err := os.WriteFile(file, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}

	if p := CheckExportDir(dir); p != nil {
		t.Errorf("a real directory is usable: %v", p)
	}

	missing := filepath.Join(dir, "not-here")
	p := CheckExportDir(missing)
	if p == nil {
		t.Fatal("a directory that does not exist must be refused")
	}
	if !strings.Contains(p.Detail, missing) {
		t.Errorf("the message must name it: %q", p.Detail)
	}
	if _, err := os.Stat(missing); !os.IsNotExist(err) {
		t.Error("nothing may be created: a typo is not an intention")
	}

	if p := CheckExportDir(file); p == nil || !strings.Contains(p.Detail, "not a directory") {
		t.Errorf("a path that is not a directory must be refused, got %v", p)
	}
	if p := CheckExportDir("   "); p == nil {
		t.Error("an empty directory must be refused")
	}
	var nilProblem *ExportDirProblem
	if nilProblem.Error() != "" {
		t.Error("a nil problem says nothing")
	}
}

// --- 2.3 completion ---

func TestCompleteDir(t *testing.T) {
	base := t.TempDir()
	for _, n := range []string{"docs", "downloads", "drafts", "other"} {
		if err := os.Mkdir(filepath.Join(base, n), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(base, "dofile"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}

	got := CompleteDir(filepath.Join(base, "do"))
	want := map[string]bool{
		filepath.Join(base, "docs"):      true,
		filepath.Join(base, "downloads"): true,
	}
	if len(got) != 2 {
		t.Fatalf("got %v, want the two directories starting with do", got)
	}
	for _, g := range got {
		if !want[g] {
			t.Errorf("unexpected suggestion %q", g)
		}
	}

	if all := CompleteDir(base + "/"); len(all) != 4 {
		t.Errorf("got %v, want every directory under it", all)
	}
	if none := CompleteDir(filepath.Join(base, "zzz")); len(none) != 0 {
		t.Errorf("got %v, want nothing", none)
	}
	if none := CompleteDir("/no/such/place/x"); len(none) != 0 {
		t.Errorf("got %v, want nothing for a parent that is not there", none)
	}
}

// --- 1.3 the retired key ---

func TestEditCommandIsRetired(t *testing.T) {
	if RetiredConfigKeys["edit_command"] == "" {
		t.Error("edit_command must carry a reason to report")
	}

	path := filepath.Join(t.TempDir(), "config.yml")
	body := "scandirs:\n  include: [x]\nedit_command: code %WORKING_DIRECTORY\nchange_mode: active\n"
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	got := RetiredKeysIn(path)
	if len(got) != 2 || got[0] != "change_mode" || got[1] != "edit_command" {
		t.Errorf("got %v, want both in a stable order", got)
	}
}

func TestExportDirIsParsedFromTheConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yml")
	if err := os.WriteFile(path, []byte("export_dir: /srv/exports\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := ParseConfigFile(path, "")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ExportDir != "/srv/exports" {
		t.Errorf("got %q", cfg.ExportDir)
	}
}

func TestExportDirProblemError(t *testing.T) {
	p := &ExportDirProblem{Detail: "/nowhere does not exist"}
	if p.Error() != "/nowhere does not exist" {
		t.Errorf("got %q", p.Error())
	}
}

func TestDefaultExportDirWithNoHome(t *testing.T) {
	// A machine with no resolvable home is the one case where there is nothing
	// better to fall back to than the working directory.
	t.Setenv("HOME", "")
	if got := DefaultExportDir(&Config{}); got != "." && got == "" {
		t.Errorf("got %q, want a usable directory", got)
	}
}
