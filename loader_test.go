package env

import (
	"os"
	"testing"
)

func TestLoadEnvFileIfExists_testingEnv(t *testing.T) {
	tempDir := t.TempDir()
	dotEnvFile := tempDir + "/.env.testing"

	// Write mock .env.testing
	err := os.WriteFile(dotEnvFile, []byte("FAKE_ENV_TESTING=testing_value\nAPP_DEBUG=0"), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp .env.testing: %v", err)
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
	if loadEnvFile("does-not-exist") {
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
