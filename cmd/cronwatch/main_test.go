package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestMain_ReportFlag verifies the binary exits cleanly when --report is passed
// with a valid config. Requires the binary to be built first via go test -run.
func TestMain_ReportFlag(t *testing.T) {
	if os.Getenv("CRONWATCH_INTEGRATION") == "" {
		t.Skip("set CRONWATCH_INTEGRATION=1 to run binary integration tests")
	}

	bin := buildBinary(t)
	cfgPath := writeMinimalConfig(t)

	cmd := exec.Command(bin, "--config", cfgPath, "--report")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected clean exit, got error: %v\noutput: %s", err, out)
	}
}

func TestMain_MissingConfig(t *testing.T) {
	if os.Getenv("CRONWATCH_INTEGRATION") == "" {
		t.Skip("set CRONWATCH_INTEGRATION=1 to run binary integration tests")
	}

	bin := buildBinary(t)
	cmd := exec.Command(bin, "--config", "/nonexistent/cronwatch.yaml", "--report")
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit for missing config, got nil")
	}
}

// TestMain_InvalidConfig verifies the binary exits with an error when the
// config file contains invalid YAML.
func TestMain_InvalidConfig(t *testing.T) {
	if os.Getenv("CRONWATCH_INTEGRATION") == "" {
		t.Skip("set CRONWATCH_INTEGRATION=1 to run binary integration tests")
	}

	bin := buildBinary(t)
	cfgPath := writeInvalidConfig(t)

	cmd := exec.Command(bin, "--config", cfgPath, "--report")
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit for invalid config YAML, got nil")
	}
}

func buildBinary(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	bin := filepath.Join(dir, "cronwatch")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return bin
}

func writeMinimalConfig(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "cronwatch.yaml")
	content := `store_path: ` + filepath.Join(dir, "store.json") + `
interval: 60
webhook: ""
jobs:
  - name: test-job
    cron: "* * * * *"
    grace_period: 5
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

func writeInvalidConfig(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "cronwatch.yaml")
	content := `this: is: not: valid: yaml: [
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write invalid config: %v", err)
	}
	return path
}
