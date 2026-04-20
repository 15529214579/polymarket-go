package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDotEnv_MissingIsNotError(t *testing.T) {
	if err := LoadDotEnv(filepath.Join(t.TempDir(), "nope")); err != nil {
		t.Fatalf("missing file should be nil, got %v", err)
	}
}

func TestLoadDotEnv_ParsesAndRespectsExisting(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, ".env.local")
	content := `# comment
FOO=bar
QUOTED="hello world"
BLANK=
ALREADY=from-file
`
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ALREADY", "from-env") // should NOT be overwritten
	t.Setenv("FOO", "")             // cleared (empty string counts as set — keep it)
	os.Unsetenv("FOO")
	os.Unsetenv("QUOTED")
	os.Unsetenv("BLANK")
	t.Cleanup(func() {
		os.Unsetenv("FOO")
		os.Unsetenv("QUOTED")
		os.Unsetenv("BLANK")
	})
	if err := LoadDotEnv(p); err != nil {
		t.Fatalf("load: %v", err)
	}
	if got := os.Getenv("FOO"); got != "bar" {
		t.Errorf("FOO=%q", got)
	}
	if got := os.Getenv("QUOTED"); got != "hello world" {
		t.Errorf("QUOTED=%q", got)
	}
	if got := os.Getenv("ALREADY"); got != "from-env" {
		t.Errorf("ALREADY=%q (must not overwrite existing env)", got)
	}
}
