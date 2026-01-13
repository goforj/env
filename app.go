package env

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// environment helpers
const (
	Testing    = "testing"
	Local      = "local"
	Dev        = "dev"
	Staging    = "staging"
	Production = "production"
)

// IsAppEnvTesting reports whether APP_ENV is "testing" or the process looks like `go test`.
// @group Application environment
// @behavior readonly
//
// Checks APP_ENV, the -test.v flag, and arguments ending with ".test" or "-test.run".
//
// Example: APP_ENV explicitly testing
//
//	_ = os.Setenv("APP_ENV", env.Testing)
//	env.Dump(env.IsAppEnvTesting())
//
//	// #bool true
//
// Example: no test markers
//
//	_ = os.Unsetenv("APP_ENV")
//	env.Dump(env.IsAppEnvTesting())
//
//	// #bool false (outside of test binaries)
func IsAppEnvTesting() bool {
	return os.Getenv("APP_ENV") == Testing ||
		flag.Lookup("test.v") != nil ||
		isTestSuffixFromArguments()
}

// isTestSuffixFromArguments checks if the test suffix is present in the command line arguments
func isTestSuffixFromArguments() bool {
	anyArgumentContainsTestSuffix := false

	for _, arg := range os.Args {
		if strings.HasSuffix(arg, ".test") || strings.HasSuffix(arg, "-test.run") {
			anyArgumentContainsTestSuffix = true
		}
	}

	return anyArgumentContainsTestSuffix
}

// GetAppEnv returns the current APP_ENV (empty string if unset).
// @group Application environment
// @behavior readonly
//
// Example: simple retrieval
//
//	_ = os.Setenv("APP_ENV", "staging")
//	env.Dump(env.GetAppEnv())
//
//	// #string "staging"
func GetAppEnv() string {
	return os.Getenv("APP_ENV")
}

// IsAppEnv checks if APP_ENV matches any of the provided environments.
// @group Application environment
// @behavior readonly
//
// Example: match any allowed environment
//
//	_ = os.Setenv("APP_ENV", "staging")
//	env.Dump(env.IsAppEnv(env.Production, env.Staging))
//
//	// #bool true
//
// Example: unmatched environment
//
//	_ = os.Setenv("APP_ENV", "local")
//	env.Dump(env.IsAppEnv(env.Production, env.Staging))
//
//	// #bool false
func IsAppEnv(envs ...string) bool {
	current := os.Getenv("APP_ENV")
	for _, env := range envs {
		if current == env {
			return true
		}
	}
	return false
}

// IsAppEnvProduction checks if APP_ENV is "production".
// @group Application environment
// @behavior readonly
//
// Example:
//
//	_ = os.Setenv("APP_ENV", env.Production)
//	env.Dump(env.IsAppEnvProduction())
//
//	// #bool true
func IsAppEnvProduction() bool {
	return IsAppEnv(Production)
}

// IsAppEnvStaging checks if APP_ENV is "staging".
// @group Application environment
// @behavior readonly
//
// Example:
//
//	_ = os.Setenv("APP_ENV", env.Staging)
//	env.Dump(env.IsAppEnvStaging())
//
//	// #bool true
func IsAppEnvStaging() bool {
	return IsAppEnv(Staging)
}

// IsAppEnvLocalOrStaging checks if APP_ENV is either "local" or "staging".
// @group Application environment
// @behavior readonly
//
// Example:
//
//	_ = os.Setenv("APP_ENV", env.Local)
//	env.Dump(env.IsAppEnvLocalOrStaging())
//
//	// #bool true
func IsAppEnvLocalOrStaging() bool {
	return IsAppEnv(Local, Staging)
}

// IsAppEnvLocal checks if APP_ENV is "local".
// @group Application environment
// @behavior readonly
//
// Example:
//
//	_ = os.Setenv("APP_ENV", env.Local)
//	env.Dump(env.IsAppEnvLocal())
//
//	// #bool true
func IsAppEnvLocal() bool {
	return IsAppEnv(Local)
}

// IsAppEnvDev checks if APP_ENV is "dev".
// @group Application environment
// @behavior readonly
//
// Example:
//
//	_ = os.Setenv("APP_ENV", env.Dev)
//	env.Dump(env.IsAppEnvDev())
//
//	// #bool true
func IsAppEnvDev() bool {
	return IsAppEnv(Dev)
}

// IsAppEnvTestingOrLocal checks if APP_ENV is "testing" or "local".
// @group Application environment
// @behavior readonly
//
// Example:
//
//	_ = os.Setenv("APP_ENV", env.Testing)
//	env.Dump(env.IsAppEnvTestingOrLocal())
//
//	// #bool true
func IsAppEnvTestingOrLocal() bool {
	return IsAppEnv(Testing, Local)
}

// SetAppEnv sets APP_ENV to a supported value.
// @group Application environment
// @behavior mutates-process-env
//
// Returns an error when the environment is unsupported.
//
// Example: set a supported environment
//
//	_ = env.SetAppEnv(env.Staging)
//	env.Dump(env.GetAppEnv())
//
//	// #string "staging"
func SetAppEnv(appEnv string) error {
	if !isValidAppEnv(appEnv) {
		return fmt.Errorf("invalid APP_ENV: %s", appEnv)
	}
	return os.Setenv("APP_ENV", appEnv)
}

// SetAppEnvLocal sets APP_ENV to "local".
// @group Application environment
// @behavior mutates-process-env
//
// Example:
//
//	_ = env.SetAppEnvLocal()
//	env.Dump(env.GetAppEnv())
//
//	// #string "local"
func SetAppEnvLocal() error {
	return SetAppEnv(Local)
}

// SetAppEnvDev sets APP_ENV to "dev".
// @group Application environment
// @behavior mutates-process-env
//
// Example:
//
//	_ = env.SetAppEnvDev()
//	env.Dump(env.GetAppEnv())
//
//	// #string "dev"
func SetAppEnvDev() error {
	return SetAppEnv(Dev)
}

// SetAppEnvStaging sets APP_ENV to "staging".
// @group Application environment
// @behavior mutates-process-env
//
// Example:
//
//	_ = env.SetAppEnvStaging()
//	env.Dump(env.GetAppEnv())
//
//	// #string "staging"
func SetAppEnvStaging() error {
	return SetAppEnv(Staging)
}

// SetAppEnvProduction sets APP_ENV to "production".
// @group Application environment
// @behavior mutates-process-env
//
// Example:
//
//	_ = env.SetAppEnvProduction()
//	env.Dump(env.GetAppEnv())
//
//	// #string "production"
func SetAppEnvProduction() error {
	return SetAppEnv(Production)
}

// SetAppEnvTesting sets APP_ENV to "testing".
// @group Application environment
// @behavior mutates-process-env
//
// Example:
//
//	_ = env.SetAppEnvTesting()
//	env.Dump(env.GetAppEnv())
//
//	// #string "testing"
func SetAppEnvTesting() error {
	return SetAppEnv(Testing)
}

func isValidAppEnv(appEnv string) bool {
	switch appEnv {
	case Testing, Local, Dev, Staging, Production:
		return true
	default:
		return false
	}
}
