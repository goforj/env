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
var loadedEnvKeys = map[string]struct{}{}

// Load loads .env with optional layering for .env.local/.env.staging/.env.production,
// plus .env.testing/.env.host when present. It only applies once per process;
// subsequent calls return without reloading because the result is cached. Use
// Reload to re-read env files after the first load.
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
//	_ = os.WriteFile(filepath.Join(tmp, ".env.testing"), []byte("PORT=9090\nENV_DEBUG=0"), 0o644)
//	_ = os.Chdir(tmp)
//	_ = os.Setenv("APP_ENV", env.Testing)
//
//	_ = env.Load()
//	env.Dump(os.Getenv("PORT"))
//	// #string "9090"
//
// Example: default .env on a host
//
//	_ = os.WriteFile(".env", []byte("SERVICE=api\nENV_DEBUG=3"), 0o644)
//	_ = env.Load()
//	env.Dump(os.Getenv("SERVICE"))
//	// #string "api"
func Load() error {
	return load(false)
}

// Reload re-applies the same layered env loading as Load, even if Load already
// ran earlier in the same process.
// @group Environment loading
// @behavior mutates-process-env
//
// Behavior:
//   - Sets APP_ENV=local when unset.
//   - Re-runs the same .env/.env.<app-env>/.env.host/.env.testing layering.
//   - Uses overload semantics, so reloaded values replace previously loaded ones.
//
// Example: refresh changed env files
//
//	_ = os.WriteFile(".env", []byte("SERVICE=api"), 0o644)
//	_ = env.Load()
//	_ = os.WriteFile(".env", []byte("SERVICE=worker"), 0o644)
//	_ = env.Reload()
//	env.Dump(os.Getenv("SERVICE"))
//	// #string "worker"
func Reload() error {
	return load(true)
}

func load(force bool) error {
	if force {
		clearLoadedEnvKeys()
	}

	if os.Getenv("APP_ENV") == "" {
		_ = os.Setenv("APP_ENV", Local)
	}

	// avoid re-loading env files
	if envLoaded && !force {
		return nil
	}

	// load base env first; layer testing/host overrides afterward
	var loadedFiles []string

	// load top-level .env
	if ok, path, keys := loadEnvFile(fileEnv); ok {
		loadedFiles = append(loadedFiles, path)
		recordLoadedEnvKeys(keys)
	}

	if envFile, ok := envFileForAppEnv(os.Getenv("APP_ENV")); ok {
		if ok, path, keys := loadEnvFile(envFile); ok {
			loadedFiles = append(loadedFiles, path)
			recordLoadedEnvKeys(keys)
		}
	}

	// search for global .env.host
	// we're likely talking from host -> container network
	// used from IDEs
	if IsHostEnvironment() || IsDockerInDocker() {
		if ok, path, keys := loadEnvFile(fileEnvHost); ok {
			loadedFiles = append(loadedFiles, path)
			recordLoadedEnvKeys(keys)
		}
	}

	// use testing envs when the environment indicates tests
	if IsAppEnvTesting() {
		if ok, path, keys := loadEnvFile(envFileTesting); ok {
			loadedFiles = append(loadedFiles, path)
			recordLoadedEnvKeys(keys)
		}
	}

	// display loaded env files
	if GetInt("ENV_DEBUG", "0") >= 3 {
		printLoadedEnvFiles(loadedFiles)
	}

	// mark as loaded
	envLoaded = true

	return nil
}

// LoadEnvFileIfExists is a compatibility alias for Load.
// @group Environment loading
//
// Example:
//
//	_ = env.LoadEnvFileIfExists()
func LoadEnvFileIfExists() error {
	return Load()
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

// IsEnvLoaded reports whether Load or LoadEnvFileIfExists was executed in this process.
// @group Environment loading
// @behavior readonly
//
// Example:
//
//	env.Dump(env.IsEnvLoaded())
//	// #bool true  (after Load)
//	// #bool false (otherwise)
func IsEnvLoaded() bool {
	return envLoaded
}

// searches for .env file through directory traversal
// loads .env file if found
func loadEnvFile(envFile string) (bool, string, []string) {
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
		values, err := godotenv.Read(path)
		if err != nil {
			panic(err)
		}
		if err := godotenv.Overload(path); err != nil {
			panic(err)
		}
		return true, path, mapKeys(values)
	}

	return found, path, nil
}

func clearLoadedEnvKeys() {
	for key := range loadedEnvKeys {
		_ = os.Unsetenv(key)
	}
	loadedEnvKeys = map[string]struct{}{}
}

func recordLoadedEnvKeys(keys []string) {
	for _, key := range keys {
		loadedEnvKeys[key] = struct{}{}
	}
}

func mapKeys(values map[string]string) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	return keys
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
