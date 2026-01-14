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
	if IsAppEnvLocal() || IsAppEnvTestingOrLocal() {
		t.Fatalf("unexpected local/testing match")
	}

	_ = os.Setenv("APP_ENV", Production)
	if !IsAppEnvProduction() {
		t.Fatalf("IsAppEnvProduction should be true")
	}

	_ = os.Setenv("APP_ENV", Local)
	if !IsAppEnvLocalOrStaging() {
		t.Fatalf("IsAppEnvLocalOrStaging should be true for local")
	}

	_ = os.Setenv("APP_ENV", Staging)
	if !IsAppEnvStaging() {
		t.Fatalf("IsAppEnvStaging should be true")
	}
}

func TestSetAppEnv(t *testing.T) {
	t.Cleanup(func() { _ = os.Unsetenv("APP_ENV") })

	if err := SetAppEnv(Staging); err != nil {
		t.Fatalf("SetAppEnv: %v", err)
	}
	if got := os.Getenv("APP_ENV"); got != Staging {
		t.Fatalf("expected APP_ENV=%s, got %q", Staging, got)
	}

	_ = os.Setenv("APP_ENV", Local)
	if err := SetAppEnv("invalid"); err == nil {
		t.Fatalf("expected error for invalid APP_ENV")
	}
	if got := os.Getenv("APP_ENV"); got != Local {
		t.Fatalf("expected APP_ENV to remain %q, got %q", Local, got)
	}
}

func TestSetAppEnvHelpers(t *testing.T) {
	t.Cleanup(func() { _ = os.Unsetenv("APP_ENV") })

	cases := []struct {
		name   string
		setter func() error
		want   string
	}{
		{name: "local", setter: SetAppEnvLocal, want: Local},
		{name: "staging", setter: SetAppEnvStaging, want: Staging},
		{name: "production", setter: SetAppEnvProduction, want: Production},
		{name: "testing", setter: SetAppEnvTesting, want: Testing},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.setter(); err != nil {
				t.Fatalf("setter: %v", err)
			}
			if got := os.Getenv("APP_ENV"); got != tc.want {
				t.Fatalf("expected APP_ENV=%s, got %q", tc.want, got)
			}
		})
	}
}
