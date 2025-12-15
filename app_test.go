package env

import (
	"os"
	"testing"
)

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
