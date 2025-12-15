package env

import (
	"bytes"
	"os"
	"testing"
)

func expectPanic(t *testing.T, name string, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic in %s", name)
		}
	}()
	fn()
}

func TestIsTestSuffixFromArguments(t *testing.T) {
	original := os.Args
	defer func() { os.Args = original }()

	os.Args = []string{"cmd", "example.test"}
	if !isTestSuffixFromArguments() {
		t.Fatalf("expected true when args include .test")
	}

	os.Args = []string{"cmd", "run"}
	if isTestSuffixFromArguments() {
		t.Fatalf("expected false when no test args present")
	}
}

func TestAppEnvHelpers(t *testing.T) {
	t.Cleanup(func() { _ = os.Unsetenv("APP_ENV") })

	_ = os.Setenv("APP_ENV", Staging)
	if got := GetAppEnv(); got != Staging {
		t.Fatalf("GetAppEnv mismatch: %s", got)
	}
	if !IsAppEnv(Staging, Production) {
		t.Fatalf("IsAppEnv should match staging")
	}
	if IsAppEnvDev() || IsAppEnvLocal() || IsAppEnvTestingOrLocal() {
		t.Fatalf("unexpected dev/local/testing match")
	}

	_ = os.Setenv("APP_ENV", Production)
	if !IsAppEnvProduction() {
		t.Fatalf("IsAppEnvProduction should be true")
	}

	_ = os.Setenv("APP_ENV", Local)
	if !IsAppEnvLocalOrStaging() {
		t.Fatalf("IsAppEnvLocalOrStaging should be true for local")
	}
	if IsAppEnvDev() {
		t.Fatalf("IsAppEnvDev should be false for local")
	}

	_ = os.Setenv("APP_ENV", Staging)
	if !IsAppEnvStaging() {
		t.Fatalf("IsAppEnvStaging should be true")
	}
}

func TestTypedGetters_InvalidPanics(t *testing.T) {
	cases := []struct {
		key string
		val string
		fn  func()
	}{
		{"BAD_INT", "nope", func() { GetInt("BAD_INT", "") }},
		{"BAD_INT64", "xx", func() { GetInt64("BAD_INT64", "") }},
		{"BAD_UINT", "-1", func() { GetUint("BAD_UINT", "") }},
		{"BAD_UINT64", "-1", func() { GetUint64("BAD_UINT64", "") }},
		{"BAD_FLOAT", "not-float", func() { GetFloat("BAD_FLOAT", "") }},
		{"BAD_BOOL", "maybe", func() { GetBool("BAD_BOOL", "") }},
		{"BAD_DURATION", "sometimes", func() { GetDuration("BAD_DURATION", "") }},
	}

	for _, tt := range cases {
		t.Run(tt.key, func(t *testing.T) {
			_ = os.Setenv(tt.key, tt.val)
			expectPanic(t, tt.key, tt.fn)
		})
	}
}

func TestSliceAndMapEmptyFallbacks(t *testing.T) {
	withEnv("EMPTY_SLICE", "", func() {
		got := GetSlice("EMPTY_SLICE", "")
		if len(got) != 0 {
			t.Fatalf("expected empty slice, got %v", got)
		}
	})

	withEnv("EMPTY_MAP", "   ", func() {
		got := GetMap("EMPTY_MAP", "")
		if len(got) != 0 {
			t.Fatalf("expected empty map, got %v", got)
		}
	})
}

func TestDumpUsesWriter(t *testing.T) {
	original := dumpWriter
	defer func() { dumpWriter = original }()

	var buf bytes.Buffer
	setDumpWriter(&buf)
	Dump("status", 42)

	out := buf.String()
	if out == "" || !bytes.Contains([]byte(out), []byte("status")) || !bytes.Contains([]byte(out), []byte("42")) {
		t.Fatalf("expected output to include values, got: %q", out)
	}
}

func TestIsDockerInDocker_NoDockerenv(t *testing.T) {
	defer reset()

	statFile = func(path string) (os.FileInfo, error) { return nil, os.ErrNotExist }
	if IsDockerInDocker() {
		t.Fatalf("expected false when /.dockerenv missing")
	}
}

func TestIsDockerInDocker_DockerenvNoSocket(t *testing.T) {
	defer reset()

	statFile = func(path string) (os.FileInfo, error) {
		if path == fileDockerEnv {
			return nil, nil
		}
		return nil, os.ErrNotExist
	}
	if IsDockerInDocker() {
		t.Fatalf("expected false when dockerenv exists but docker.sock missing")
	}
}

func TestIsDockerHost_NoSocket(t *testing.T) {
	defer reset()

	statFile = func(path string) (os.FileInfo, error) { return nil, os.ErrNotExist }
	if IsDockerHost() {
		t.Fatalf("expected false when docker.sock missing")
	}
}

func TestIsDockerHost_ReadError(t *testing.T) {
	defer reset()

	statFile = func(path string) (os.FileInfo, error) {
		if path == fileDockerSock {
			return nil, nil
		}
		return nil, os.ErrNotExist
	}
	readFile = func(path string) ([]byte, error) { return nil, os.ErrPermission }

	if IsDockerHost() {
		t.Fatalf("expected false when cgroup read fails")
	}
}

func TestIsContainer_ReadErrorButEnvPresent(t *testing.T) {
	defer reset()

	statFile = func(path string) (os.FileInfo, error) { return nil, os.ErrNotExist }
	readFile = func(path string) ([]byte, error) { return nil, os.ErrPermission }
	mockEnv(map[string]string{"KUBERNETES_SERVICE_HOST": "10.0.0.1"})

	if !IsContainer() {
		t.Fatalf("expected true when kube env present even if cgroup read fails")
	}
}

func TestIsContainer_ReadErrorNoEnv(t *testing.T) {
	defer reset()

	statFile = func(path string) (os.FileInfo, error) { return nil, os.ErrNotExist }
	readFile = func(path string) ([]byte, error) { return nil, os.ErrPermission }
	mockEnv(nil)

	if IsContainer() {
		t.Fatalf("expected false when cgroup read fails and no env hint")
	}
}

func TestIsBSD_FalseBranch(t *testing.T) {
	defer func() {
		goos = "linux"
	}()

	goos = "linux"
	if IsBSD() {
		t.Fatalf("expected false for linux")
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
	if IsEnvLoaded() == false {
		t.Fatalf("expected envLoaded flag set")
	}
}

func TestLoadEnvFileIfExists_WithDotEnv(t *testing.T) {
	wd, _ := os.Getwd()
	t.Cleanup(func() {
		envLoaded = false
		_ = os.Chdir(wd)
		_ = os.Unsetenv("IN_DOTENV")
		_ = os.Unsetenv("APP_ENV")
	})

	tmp := t.TempDir()
	if err := os.WriteFile(tmp+"/.env", []byte("IN_DOTENV=yes\nAPP_DEBUG=0"), 0o644); err != nil {
		t.Fatalf("write .env: %v", err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	envLoaded = false
	_ = os.Setenv("APP_ENV", Dev)
	if err := LoadEnvFileIfExists(); err != nil {
		t.Fatalf("load .env: %v", err)
	}
	if got := os.Getenv("IN_DOTENV"); got != "yes" {
		t.Fatalf("expected IN_DOTENV=yes, got %q", got)
	}
}

func TestLoadEnvFile_NotFound(t *testing.T) {
	if loadEnvFile("does-not-exist") {
		t.Fatalf("expected false when file missing")
	}
}
