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
	envFileLocal   = ".env.local"
	envFileStaging = ".env.staging"
	envFileProd    = ".env.production"
)

// envLoaded is a flag to check if the environment file has been loaded
var envLoaded = false

// LoadEnvFileIfExists loads .env with optional layering for .env.local/.env.staging/.env.production,
// plus .env.testing/.env.host when present.
// @group Environment loading
// @behavior mutates-process-env
//
// Behavior:
//   - Sets APP_ENV=local when unset.
//   - Chooses .env.testing when APP_ENV indicates tests (or Go test flags are present).
//   - Loads .env first when present; .env.<app-env> overlays for local/staging/production.
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
	if os.Getenv("APP_ENV") == "" {
		_ = os.Setenv("APP_ENV", Local)
	}

	// avoid re-loading env files
	if envLoaded {
		return nil
	}

	// load base env first; layer testing/host overrides afterward
	var loadedFiles []string

	// load top-level .env
	if ok, path := loadEnvFile(fileEnv); ok {
		loadedFiles = append(loadedFiles, path)
	}

	if envFile, ok := envFileForAppEnv(os.Getenv("APP_ENV")); ok {
		if ok, path := loadEnvFile(envFile); ok {
			loadedFiles = append(loadedFiles, path)
		}
	}

	// search for global .env.host
	// we're likely talking from host -> container network
	// used from IDEs
	if IsHostEnvironment() || IsDockerInDocker() {
		if ok, path := loadEnvFile(fileEnvHost); ok {
			loadedFiles = append(loadedFiles, path)
		}
	}

	// use testing envs when the environment indicates tests
	if IsAppEnvTesting() {
		if ok, path := loadEnvFile(envFileTesting); ok {
			loadedFiles = append(loadedFiles, path)
		}
	}

	// display loaded env files
	if GetInt("APP_DEBUG", "0") >= 3 {
		printLoadedEnvFiles(loadedFiles)
	}

	// mark as loaded
	envLoaded = true

	return nil
}

// envFileForAppEnv returns the layered env filename for the given APP_ENV.
func envFileForAppEnv(appEnv string) (string, bool) {
	switch appEnv {
	case Local:
		return envFileLocal, true
	case Staging:
		return envFileStaging, true
	case Production:
		return envFileProd, true
	default:
		return "", false
	}
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
func loadEnvFile(envFile string) (bool, string) {
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

	return found, path
}

// ANSI color codes
const (
	colorGray  = "\033[90m"
	colorReset = "\033[0m"
)

// debugMark returns a gray dot symbol for debug output.
func debugMark() string {
	return colorMark(colorGray, "·")
}

// colorMark wraps a symbol in the provided ANSI color.
func colorMark(color, symbol string) string {
	return fmt.Sprintf("%s%s%s", color, symbol, colorReset)
}

// printLoadedEnvFiles outputs loaded env files to stdout
func printLoadedEnvFiles(paths []string) {
	if len(paths) == 0 {
		return
	}
	for _, path := range paths {
		fmt.Printf(" %s .env file loader · env [%v] file [%v]\n", debugMark(), os.Getenv("APP_ENV"), path)
	}
}
