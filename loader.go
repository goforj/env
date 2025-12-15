package env

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

// MaxDirectorySeekLevels is the number of directory
// levels a .env file needs to be searched in
const MaxDirectorySeekLevels int = 10

const (
	runtimeDarwin  = "darwin"
	runtimeWindows = "windows"

	// file
	fileEnv        = ".env"
	fileEnvHost    = ".env.host"
	envFileTesting = ".env.testing"
)

// envLoaded is a flag to check if the environment file has been loaded
var envLoaded = false

// LoadEnvFileIfExists loads .env/.env.testing/.env.host when present.
// @group Environment loading
// @behavior mutates-process-env
//
// Behavior:
//   - Sets APP_ENV=local on macOS/Windows when unset.
//   - Chooses .env.testing when APP_ENV indicates tests (or Go test flags are present).
//   - Loads .env.host for host-to-container networking when running on the host or DinD.
//   - Idempotent: subsequent calls no-op after the first load.
//
// Example: test-specific env file
//
//	tmp, _ := os.MkdirTemp("", "envdoc")
//	_ = os.WriteFile(filepath.Join(tmp, ".env.testing"), []byte("PORT=9090\nAPP_DEBUG=0"), 0o644)
//	_ = os.Chdir(tmp)
//	_ = os.Setenv("APP_ENV", env.Testing)
//
//	_ = env.LoadEnvFileIfExists()
//	env.Dump(os.Getenv("PORT"))
//
//	// #string "9090"
//
// Example: default .env on a host
//
//	_ = os.WriteFile(".env", []byte("SERVICE=api\nAPP_DEBUG=3"), 0o644)
//	_ = env.LoadEnvFileIfExists()
//	env.Dump(os.Getenv("SERVICE"))
//
//	// #string "api"
func LoadEnvFileIfExists() error {
	_ = os.Setenv("APP_ENV", Local)

	if !envLoaded {
		// use dev or testing envs depending on the environment
		envLoadMsg := ""
		envFile := fileEnv
		if IsAppEnvTesting() {
			envFile = envFileTesting
		}

		// load top-level .env
		if loadEnvFile(envFile) {
			envLoadMsg = fmt.Sprintf("[LoadEnv] APP_ENV [%v] ENV_FILE [%v]", os.Getenv("APP_ENV"), envFile)
		}

		// display env: [LoadEnv] APP_ENV [local] ENV_FILE [.env]
		if GetInt("APP_DEBUG", "0") >= 3 {
			fmt.Println(envLoadMsg)
		}

		envLoaded = true

		// search for global .env.host
		// we're likely talking from host -> container network
		// used from IDEs
		if IsHostEnvironment() || IsDockerInDocker() {
			env := fileEnvHost
			if loadEnvFile(env) {
				if GetInt("APP_DEBUG", "0") > 0 {
					fmt.Println(fmt.Sprintf("Loaded environment [env] APP_ENV [%v] ENV_FILE [%v]", os.Getenv("APP_ENV"), env))
				}
			}
		}
	}

	return nil
}

// IsEnvLoaded reports whether LoadEnvFileIfExists was executed in this process.
// @group Environment loading
// @behavior readonly
//
// Example:
//
//	env.Dump(env.IsEnvLoaded())
//
//	// #bool true  (after LoadEnvFileIfExists)
//	// #bool false (otherwise)
func IsEnvLoaded() bool {
	return envLoaded
}

// searches for .env file through directory traversal
// loads .env file if found
func loadEnvFile(envFile string) bool {
	var path string
	found := false
	for i := 0; i < MaxDirectorySeekLevels; i++ {
		if _, err := os.Stat(path + envFile); err == nil {
			path += envFile
			found = true
			break
		}
		path += "../"
	}

	if found {
		if err := godotenv.Overload(path); err != nil {
			panic(err)
		}
	}

	return found
}
