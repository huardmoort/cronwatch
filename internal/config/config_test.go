package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/cronwatch/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "cronwatch-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_Valid(t *testing.T) {
	yaml := `
log_level: debug
state_dir: /tmp/cw
alerts:
  email: ops@example.com
jobs:
  - name: backup
    schedule: "0 2 * * *"
    timeout: 30m
    command: /usr/local/bin/backup.sh
`
	path := writeTemp(t, yaml)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("log_level: got %q, want %q", cfg.LogLevel, "debug")
	}
	if len(cfg.Jobs) != 1 {
		t.Fatalf("jobs count: got %d, want 1", len(cfg.Jobs))
	}
	if cfg.Jobs[0].Timeout != 30*time.Minute {
		t.Errorf("timeout: got %v, want 30m", cfg.Jobs[0].Timeout)
	}
}

func TestLoad_Defaults(t *testing.T) {
	yaml := `
jobs:
  - name: ping
    schedule: "* * * * *"
`
	path := writeTemp(t, yaml)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("default log_level: got %q", cfg.LogLevel)
	}
	if cfg.StateDir != "/var/lib/cronwatch" {
		t.Errorf("default state_dir: got %q", cfg.StateDir)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_NoJobs(t *testing.T) {
	path := writeTemp(t, "log_level: info\n")
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for empty jobs")
	}
}

func TestLoad_DuplicateJobName(t *testing.T) {
	yaml := `
jobs:
  - name: dup
    schedule: "* * * * *"
  - name: dup
    schedule: "0 * * * *"
`
	path := writeTemp(t, yaml)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for duplicate job name")
	}
}
