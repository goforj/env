package env

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestLoadEnvFileIfExists_testingEnv(t *testing.T) {
	tempDir := t.TempDir()
	dotEnvFile := tempDir + "/.env.testing"
	baseEnvFile := tempDir + "/.env"

	// Write mock .env.testing
	err := os.WriteFile(dotEnvFile, []byte("FAKE_ENV_TESTING=testing_value\nAPP_DEBUG=0"), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp .env.testing: %v", err)
	}
	if err := os.WriteFile(baseEnvFile, []byte("FAKE_ENV_BASE=base_value\nFAKE_ENV_TESTING=base_override\n"), 0644); err != nil {
		t.Fatalf("Failed to create temp .env: %v", err)
	}

	// Save original working dir to restore later
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalDir) // restore after test

	// Move to the temp directory where our mock .env.testing exists
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}

	// Set environment to "testing" to trigger .env.testing logic
	_ = os.Setenv("APP_ENV", "testing")

	// Reset internal state
	envLoaded = false

	err = LoadEnvFileIfExists()
	if err != nil {
		t.Fatalf("LoadEnvFileIfExists failed: %v", err)
	}

	val := os.Getenv("FAKE_ENV_TESTING")
	if val != "testing_value" {
		t.Errorf("Expected FAKE_ENV_TESTING to be 'testing_value', got %s", val)
	}
	baseVal := os.Getenv("FAKE_ENV_BASE")
	if baseVal != "base_value" {
		t.Errorf("Expected FAKE_ENV_BASE to be 'base_value', got %s", baseVal)
	}

	if !IsEnvLoaded() {
		t.Error("Expected IsEnvLoaded to return true")
	}
}

func TestLoadEnvFileIfExists_NoFile(t *testing.T) {
	wd, _ := os.Getwd()
	t.Cleanup(func() {
		envLoaded = false
		_ = os.Chdir(wd)
	})

	tmp := t.TempDir()
	_ = os.Chdir(tmp)

	if err := LoadEnvFileIfExists(); err != nil {
		t.Fatalf("expected no error when env file missing: %v", err)
	}
	if !IsEnvLoaded() {
		t.Fatalf("expected envLoaded flag set")
	}
}

func TestLoadEnvFileIfExists_WithDotEnvHostBranch(t *testing.T) {
	defer reset()

	wd, _ := os.Getwd()
	t.Cleanup(func() {
		envLoaded = false
		_ = os.Chdir(wd)
		_ = os.Unsetenv("HOST_BRANCH")
	})

	tmp := t.TempDir()
	if err := os.WriteFile(tmp+"/.env.testing", []byte("APP_DEBUG=3"), 0o644); err != nil {
		t.Fatalf("write .env.testing: %v", err)
	}
	if err := os.WriteFile(tmp+"/.env.host", []byte("HOST_BRANCH=hit\nAPP_DEBUG=3"), 0o644); err != nil {
		t.Fatalf("write .env.host: %v", err)
	}

	statFile = func(path string) (os.FileInfo, error) {
		switch path {
		case fileDockerEnv, fileDockerSock:
			return nil, nil
		default:
			return nil, os.ErrNotExist
		}
	}
	readFile = func(path string) ([]byte, error) { return []byte("0::/docker/xyz"), nil }
	mockEnv(nil)

	_ = os.Chdir(tmp)
	envLoaded = false
	if err := LoadEnvFileIfExists(); err != nil {
		t.Fatalf("LoadEnvFileIfExists: %v", err)
	}
	if os.Getenv("HOST_BRANCH") != "hit" {
		t.Fatalf("expected HOST_BRANCH to load from .env.host")
	}
}

func TestLoadEnvFile_NotFound(t *testing.T) {
	if ok, _ := loadEnvFile("does-not-exist"); ok {
		t.Fatalf("expected false when file missing")
	}
}

func TestLoadEnvFile_PanicsOnBadFile(t *testing.T) {
	tmp := t.TempDir()
	_ = os.Mkdir(tmp+"/.env.testing", 0o755)

	wd, _ := os.Getwd()
	_ = os.Chdir(tmp)
	t.Cleanup(func() { _ = os.Chdir(wd) })

	expectPanic(t, "loadEnvFile panic", func() {
		loadEnvFile(".env.testing")
	})
}

func TestPrintLoadedEnvFiles_NoPaths(t *testing.T) {
	output := captureStdout(t, func() {
		printLoadedEnvFiles(nil)
	})
	if output != "" {
		t.Fatalf("expected no output, got %q", output)
	}
}

func TestPrintLoadedEnvFiles_WithPaths(t *testing.T) {
	t.Setenv("APP_ENV", Testing)
	output := captureStdout(t, func() {
		printLoadedEnvFiles([]string{"./.env", "./.env.testing"})
	})
	if !strings.Contains(output, "env [testing]") {
		t.Fatalf("expected output to include APP_ENV, got %q", output)
	}
	if !strings.Contains(output, "file [./.env]") || !strings.Contains(output, "file [./.env.testing]") {
		t.Fatalf("expected output to include paths, got %q", output)
	}
}

func TestEnvFileForAppEnv(t *testing.T) {
	cases := map[string]string{
		Local:      ".env.local",
		Staging:    ".env.staging",
		Production: ".env.production",
	}

	for appEnv, expected := range cases {
		t.Run(appEnv, func(t *testing.T) {
			got, ok := envFileForAppEnv(appEnv)
			if !ok {
				t.Fatalf("expected env file for %s", appEnv)
			}
			if got != expected {
				t.Fatalf("expected %s, got %s", expected, got)
			}
		})
	}

	if _, ok := envFileForAppEnv(Testing); ok {
		t.Fatalf("expected no env file for %s", Testing)
	}
	if _, ok := envFileForAppEnv("unknown"); ok {
		t.Fatalf("expected no env file for unknown")
	}
}

func TestLoadEnvFileIfExists_LayeredLocal(t *testing.T) {
	wd, _ := os.Getwd()
	t.Cleanup(func() {
		envLoaded = false
		_ = os.Chdir(wd)
	})

	tmp := t.TempDir()
	if err := os.WriteFile(tmp+"/.env", []byte("BASE=base\nLAYER=base\n"), 0o644); err != nil {
		t.Fatalf("write .env: %v", err)
	}
	if err := os.WriteFile(tmp+"/.env.local", []byte("LAYER=local\nLOCAL_ONLY=1\n"), 0o644); err != nil {
		t.Fatalf("write .env.local: %v", err)
	}

	_ = os.Chdir(tmp)
	_ = os.Setenv("APP_ENV", Local)
	envLoaded = false

	if err := LoadEnvFileIfExists(); err != nil {
		t.Fatalf("LoadEnvFileIfExists: %v", err)
	}
	if got := os.Getenv("LAYER"); got != "local" {
		t.Fatalf("expected LAYER to be local, got %q", got)
	}
	if got := os.Getenv("LOCAL_ONLY"); got != "1" {
		t.Fatalf("expected LOCAL_ONLY to be set, got %q", got)
	}
	if got := os.Getenv("BASE"); got != "base" {
		t.Fatalf("expected BASE to be base, got %q", got)
	}
}

func TestLoadEnvFileIfExists_LayeredStaging(t *testing.T) {
	wd, _ := os.Getwd()
	t.Cleanup(func() {
		envLoaded = false
		_ = os.Chdir(wd)
	})

	tmp := t.TempDir()
	if err := os.WriteFile(tmp+"/.env", []byte("SHARED=base\n"), 0o644); err != nil {
		t.Fatalf("write .env: %v", err)
	}
	if err := os.WriteFile(tmp+"/.env.staging", []byte("SHARED=staging\nSTAGING_ONLY=1\n"), 0o644); err != nil {
		t.Fatalf("write .env.staging: %v", err)
	}

	_ = os.Chdir(tmp)
	_ = os.Setenv("APP_ENV", Staging)
	envLoaded = false

	if err := LoadEnvFileIfExists(); err != nil {
		t.Fatalf("LoadEnvFileIfExists: %v", err)
	}
	if got := os.Getenv("SHARED"); got != "staging" {
		t.Fatalf("expected SHARED to be staging, got %q", got)
	}
	if got := os.Getenv("STAGING_ONLY"); got != "1" {
		t.Fatalf("expected STAGING_ONLY to be set, got %q", got)
	}
}

func TestLoadEnvFileIfExists_LayeredProduction(t *testing.T) {
	wd, _ := os.Getwd()
	t.Cleanup(func() {
		envLoaded = false
		_ = os.Chdir(wd)
	})

	tmp := t.TempDir()
	if err := os.WriteFile(tmp+"/.env", []byte("SHARED=base\n"), 0o644); err != nil {
		t.Fatalf("write .env: %v", err)
	}
	if err := os.WriteFile(tmp+"/.env.production", []byte("SHARED=prod\nPROD_ONLY=1\n"), 0o644); err != nil {
		t.Fatalf("write .env.production: %v", err)
	}

	_ = os.Chdir(tmp)
	_ = os.Setenv("APP_ENV", Production)
	envLoaded = false

	if err := LoadEnvFileIfExists(); err != nil {
		t.Fatalf("LoadEnvFileIfExists: %v", err)
	}
	if got := os.Getenv("SHARED"); got != "prod" {
		t.Fatalf("expected SHARED to be prod, got %q", got)
	}
	if got := os.Getenv("PROD_ONLY"); got != "1" {
		t.Fatalf("expected PROD_ONLY to be set, got %q", got)
	}
}

func TestLoadEnvFileIfExists_NoReload(t *testing.T) {
	wd, _ := os.Getwd()
	t.Cleanup(func() {
		envLoaded = false
		_ = os.Chdir(wd)
	})

	tmp := t.TempDir()
	if err := os.WriteFile(tmp+"/.env", []byte("SHOULD_NOT_LOAD=1\n"), 0o644); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	_ = os.Chdir(tmp)
	envLoaded = true

	if err := LoadEnvFileIfExists(); err != nil {
		t.Fatalf("LoadEnvFileIfExists: %v", err)
	}
	if got := os.Getenv("SHOULD_NOT_LOAD"); got != "" {
		t.Fatalf("expected SHOULD_NOT_LOAD to remain unset, got %q", got)
	}
}

func TestLoadEnvFileIfExists_DefaultsAppEnv(t *testing.T) {
	wd, _ := os.Getwd()
	t.Cleanup(func() {
		envLoaded = false
		_ = os.Chdir(wd)
		_ = os.Unsetenv("APP_ENV")
	})

	tmp := t.TempDir()
	_ = os.Chdir(tmp)
	_ = os.Unsetenv("APP_ENV")
	envLoaded = false

	if err := LoadEnvFileIfExists(); err != nil {
		t.Fatalf("LoadEnvFileIfExists: %v", err)
	}
	if got := os.Getenv("APP_ENV"); got != Local {
		t.Fatalf("expected APP_ENV to default to %s, got %q", Local, got)
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	original := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w
	defer func() {
		os.Stdout = original
	}()

	done := make(chan string)
	go func() {
		var buf strings.Builder
		_, _ = io.Copy(&buf, r)
		done <- buf.String()
	}()

	fn()
	_ = w.Close()
	output := <-done
	return output
}
