package main

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestStoreAddCreatesDirAndNotebook(t *testing.T) {
	s := Store{Dir: filepath.Join(t.TempDir(), "notes")}

	if _, err := s.Add("work", "buy milk"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Add("work", "call bob"); err != nil {
		t.Fatal(err)
	}

	got, err := s.Read("work")
	if err != nil {
		t.Fatal(err)
	}
	want := "# work\n\n- buy milk\n- call bob\n"
	if string(got) != want {
		t.Errorf("notebook contents = %q, want %q", got, want)
	}
}

func TestStoreEnsureKeepsExisting(t *testing.T) {
	s := Store{Dir: t.TempDir()}
	if _, err := s.Add("a", "first"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Ensure("a"); err != nil {
		t.Fatal(err)
	}
	got, _ := s.Read("a")
	if !strings.Contains(string(got), "- first") {
		t.Errorf("Ensure overwrote existing notebook: %q", got)
	}
}

func TestStoreRejectsBadNames(t *testing.T) {
	s := Store{Dir: t.TempDir()}
	for _, name := range []string{"", ".", "..", "../escape", "a/b", `a\b`} {
		if _, err := s.Path(name); err == nil {
			t.Errorf("Path(%q) succeeded, want error", name)
		}
	}
}

func TestStoreReadMissing(t *testing.T) {
	s := Store{Dir: t.TempDir()}
	if _, err := s.Read("nope"); err == nil || !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("Read of missing notebook: err = %v", err)
	}
}

func TestStoreList(t *testing.T) {
	dir := t.TempDir()
	s := Store{Dir: dir}
	for _, n := range []string{"zeta", "alpha"} {
		if _, err := s.Ensure(n); err != nil {
			t.Fatal(err)
		}
	}
	os.WriteFile(filepath.Join(dir, "ignore.txt"), nil, 0o644)
	os.Mkdir(filepath.Join(dir, "subdir.md"), 0o755)

	got, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"alpha", "zeta"}; !reflect.DeepEqual(got, want) {
		t.Errorf("List() = %v, want %v", got, want)
	}

	empty, err := Store{Dir: filepath.Join(dir, "missing")}.List()
	if err != nil || len(empty) != 0 {
		t.Errorf("List() of missing dir = %v, %v", empty, err)
	}
}

// setupEnv points HOME and the config file at a temp dir and returns the
// notes directory noteit will use.
func setupEnv(t *testing.T, config string) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	cfgPath := filepath.Join(home, "config.json")
	t.Setenv("NOTEIT_CONFIG", cfgPath)
	if config != "" {
		if err := os.WriteFile(cfgPath, []byte(config), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return filepath.Join(home, "noteit")
}

func runCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	var out bytes.Buffer
	err := run(args, &out)
	return out.String(), err
}

func TestRunAddAndView(t *testing.T) {
	dir := setupEnv(t, "")

	if _, err := runCmd(t, "add", "-n", "work", "buy", "milk"); err != nil {
		t.Fatal(err)
	}
	if _, err := runCmd(t, "add", "default", "note"); err != nil {
		t.Fatal(err)
	}

	out, err := runCmd(t, "view", "work")
	if err != nil {
		t.Fatal(err)
	}
	if out != "# work\n\n- buy milk\n" {
		t.Errorf("view work = %q", out)
	}

	out, err = runCmd(t, "view")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "- default note") {
		t.Errorf("view default = %q", out)
	}

	out, err = runCmd(t, "list")
	if err != nil {
		t.Fatal(err)
	}
	if out != "default\nwork\n" {
		t.Errorf("list = %q", out)
	}

	if _, err := os.Stat(filepath.Join(dir, "work.md")); err != nil {
		t.Errorf("expected notebook in default dir: %v", err)
	}
}

func TestRunErrors(t *testing.T) {
	setupEnv(t, "")
	for _, args := range [][]string{
		{"add"},
		{"add", "-n", "../x", "note"},
		{"view", "missing"},
		{"view", "a", "b"},
		{"bogus"},
	} {
		if _, err := runCmd(t, args...); err == nil {
			t.Errorf("run(%q) succeeded, want error", args)
		}
	}
}

func TestLoadConfig(t *testing.T) {
	setupEnv(t, `{"path": "~/notes", "editor": "code --wait", "defaultNotebook": "inbox"}`)
	home := os.Getenv("HOME")

	cfg, err := loadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(home, "notes"); cfg.Path != want {
		t.Errorf("Path = %q, want %q", cfg.Path, want)
	}
	if cfg.DefaultNotebook != "inbox" {
		t.Errorf("DefaultNotebook = %q", cfg.DefaultNotebook)
	}
	if got := cfg.editorCommand(); !reflect.DeepEqual(got, []string{"code", "--wait"}) {
		t.Errorf("editorCommand() = %v", got)
	}
}

func TestLoadConfigInvalid(t *testing.T) {
	setupEnv(t, `{not json`)
	if _, err := loadConfig(); err == nil {
		t.Error("loadConfig succeeded on invalid JSON")
	}
}

func TestEditorCommandFallback(t *testing.T) {
	t.Setenv("VISUAL", "")
	t.Setenv("EDITOR", "nano")
	if got := (Config{}).editorCommand(); !reflect.DeepEqual(got, []string{"nano"}) {
		t.Errorf("editorCommand() = %v, want [nano]", got)
	}
	t.Setenv("EDITOR", "")
	if got := (Config{}).editorCommand(); !reflect.DeepEqual(got, []string{"vi"}) {
		t.Errorf("editorCommand() = %v, want [vi]", got)
	}
}

func TestRunEdit(t *testing.T) {
	dir := setupEnv(t, `{"editor": "true"}`)
	if _, err := runCmd(t, "edit", "journal"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "journal.md")); err != nil {
		t.Errorf("edit did not create notebook: %v", err)
	}
}
