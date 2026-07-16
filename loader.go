package env

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

// MaxDirectorySeekLevels bounds env-file discovery to the working directory and nine ancestors.
const MaxDirectorySeekLevels int = 10

const (
	fileEnv        = ".env"
	fileEnvHost    = ".env.host"
	envFileTesting = ".env.testing"
	envFileLocal   = ".env.local"
	envFileStaging = ".env.staging"
	envFileProd    = ".env.production"
)

var (
	envFileGetwd = os.Getwd
	envFileStat  = os.Stat
	envFileRead  = godotenv.Read
	envLookup    = os.LookupEnv
	envSet       = os.Setenv
	envUnset     = os.Unsetenv
)

// environmentSnapshot retains both presence and value because an empty variable differs from an unset one.
type environmentSnapshot struct {
	value   string
	present bool
}

// loadedEnvironmentValue records the key's ambient value before the first successful load.
type loadedEnvironmentValue struct {
	original environmentSnapshot
}

// environmentLoaderState serializes loading and protects ownership metadata and the ambient baseline.
type environmentLoaderState struct {
	mu       sync.Mutex
	loaded   bool
	values   map[string]loadedEnvironmentValue
	baseline map[string]environmentSnapshot
}

var processEnvironmentLoader = environmentLoaderState{
	values:   make(map[string]loadedEnvironmentValue),
	baseline: make(map[string]environmentSnapshot),
}

// environmentFile contains one parsed file before any process environment mutation occurs.
type environmentFile struct {
	path   string
	values map[string]string
}

// environmentLoadPlan is the complete, deterministic result of discovery and layering.
type environmentLoadPlan struct {
	fileValues map[string]string
	defaults   map[string]string
	files      []string
	appEnv     string
}

// Load loads the nearest env files with deterministic layering.
//
// Load applies once per process. Files override ambient values, and later files override earlier
// files. Discovery and parsing complete before the process environment changes; errors leave both
// the environment and loader state unchanged. Use Reload to re-read files.
//
// @group Environment loading
// @behavior mutates-process-env
//
// Layer order:
//   - .env
//   - .env.local, .env.staging, or .env.production selected after parsing .env
//   - .env.host on hosts and Docker-in-Docker
//   - .env.testing when APP_ENV or the process identifies a test
//
// Each filename is searched independently from the working directory through at most nine
// ancestors. APP_ENV defaults to local when neither the ambient environment nor a file sets it.
//
// Example: test-specific env file
//
//	tmp, _ := os.MkdirTemp("", "envdoc")
//	defer os.RemoveAll(tmp)
//	originalDirectory, _ := os.Getwd()
//	defer os.Chdir(originalDirectory)
//	_ = os.WriteFile(filepath.Join(tmp, ".env.testing"), []byte("PORT=9090\nENV_DEBUG=0"), 0o644)
//	_ = os.Chdir(tmp)
//	_ = os.Setenv("APP_ENV", env.Testing)
//
//	_ = env.Load()
//	env.Dump(os.Getenv("PORT"))
//	// #string "9090"
func Load() error {
	return load(false)
}

// Reload re-discovers and transactionally reapplies env files even after Load has run.
//
// Keys loaded from files remain file-owned: Reload replaces runtime edits to those keys. When a
// key disappears from all files, Reload restores the ambient value (including unset versus empty)
// that existed before the first successful Load. Unrelated process variables are never changed.
//
// @group Environment loading
// @behavior mutates-process-env
//
// Example: refresh changed env files
//
//	tmp, _ := os.MkdirTemp("", "envdoc")
//	defer os.RemoveAll(tmp)
//	originalDirectory, _ := os.Getwd()
//	defer os.Chdir(originalDirectory)
//	_ = os.Chdir(tmp)
//	_ = os.WriteFile(filepath.Join(tmp, ".env"), []byte("SERVICE=api"), 0o644)
//	_ = env.Load()
//	_ = os.WriteFile(filepath.Join(tmp, ".env"), []byte("SERVICE=worker"), 0o644)
//	_ = env.Reload()
//	env.Dump(os.Getenv("SERVICE"))
//	// #string "worker"
func Reload() error {
	return load(true)
}

// load serializes discovery, application, and state publication as one loader operation.
func load(force bool) error {
	processEnvironmentLoader.mu.Lock()
	defer processEnvironmentLoader.mu.Unlock()

	if processEnvironmentLoader.loaded && !force {
		return nil
	}

	workingDirectory, err := envFileGetwd()
	if err != nil {
		return fmt.Errorf("get working directory for env loading: %w", err)
	}

	previous := cloneLoadedEnvironmentValues(processEnvironmentLoader.values)
	baseline := cloneEnvironmentSnapshots(processEnvironmentLoader.baseline)
	if !processEnvironmentLoader.loaded {
		baseline = snapshotProcessEnvironment()
	}
	plan, err := buildEnvironmentLoadPlan(filepath.Clean(workingDirectory), previous)
	if err != nil {
		return err
	}

	next, err := applyEnvironmentLoadPlan(previous, baseline, plan)
	if err != nil {
		return err
	}

	processEnvironmentLoader.values = next
	processEnvironmentLoader.baseline = baseline
	processEnvironmentLoader.loaded = true

	if environmentPlanInt(plan, previous, "ENV_DEBUG") >= 3 {
		printLoadedEnvFiles(plan.files, plan.appEnv)
	}
	return nil
}

// buildEnvironmentLoadPlan parses every selected layer before process-wide mutation begins.
func buildEnvironmentLoadPlan(startDirectory string, previous map[string]loadedEnvironmentValue) (environmentLoadPlan, error) {
	plan := environmentLoadPlan{
		fileValues: make(map[string]string),
		defaults:   make(map[string]string),
	}

	appEnv := effectiveEnvironmentValue("APP_ENV", plan.fileValues, previous)
	if appEnv == "" {
		appEnv = Local
	}

	base, found, err := loadEnvFile(startDirectory, fileEnv)
	if err != nil {
		return environmentLoadPlan{}, err
	}
	if found {
		mergeEnvironmentFile(&plan, base)
		if value, ok := plan.fileValues["APP_ENV"]; ok {
			appEnv = value
		}
	}

	if appEnvFile, ok := envFileForAppEnv(appEnv); ok {
		layer, found, err := loadEnvFile(startDirectory, appEnvFile)
		if err != nil {
			return environmentLoadPlan{}, err
		}
		if found {
			mergeEnvironmentFile(&plan, layer)
		}
	}

	lookup := func(key string) string {
		return effectiveEnvironmentValue(key, plan.fileValues, previous)
	}
	if isHostEnvironmentWithEnv(lookup) || IsDockerInDocker() {
		host, found, err := loadEnvFile(startDirectory, fileEnvHost)
		if err != nil {
			return environmentLoadPlan{}, err
		}
		if found {
			mergeEnvironmentFile(&plan, host)
		}
	}

	appEnv = effectiveEnvironmentValue("APP_ENV", plan.fileValues, previous)
	if appEnv == "" {
		appEnv = Local
	}
	if isAppEnvTestingValue(appEnv) {
		testing, found, err := loadEnvFile(startDirectory, envFileTesting)
		if err != nil {
			return environmentLoadPlan{}, err
		}
		if found {
			mergeEnvironmentFile(&plan, testing)
		}
	}

	if _, fileOwnsAppEnv := plan.fileValues["APP_ENV"]; !fileOwnsAppEnv {
		ambient := originalEnvironmentSnapshot("APP_ENV", previous)
		if !ambient.present || ambient.value == "" {
			plan.defaults["APP_ENV"] = Local
		}
	}
	plan.appEnv = environmentPlanValue(plan, previous, "APP_ENV")
	return plan, nil
}

// mergeEnvironmentFile applies one already parsed file to the in-memory layer map.
func mergeEnvironmentFile(plan *environmentLoadPlan, file environmentFile) {
	plan.files = append(plan.files, file.path)
	for key, value := range file.values {
		plan.fileValues[key] = value
	}
}

// effectiveEnvironmentValue uses stable originals for owned keys and live ambient values otherwise.
func effectiveEnvironmentValue(key string, fileValues map[string]string, previous map[string]loadedEnvironmentValue) string {
	if value, ok := fileValues[key]; ok {
		return value
	}
	snapshot := originalEnvironmentSnapshot(key, previous)
	if !snapshot.present {
		return ""
	}
	return snapshot.value
}

// originalEnvironmentSnapshot returns the stable ambient value for file-owned keys and the live value otherwise.
func originalEnvironmentSnapshot(key string, previous map[string]loadedEnvironmentValue) environmentSnapshot {
	if loaded, ok := previous[key]; ok {
		return loaded.original
	}
	value, present := envLookup(key)
	return environmentSnapshot{value: value, present: present}
}

// environmentPlanValue returns the value that will be visible after a successful plan application.
func environmentPlanValue(plan environmentLoadPlan, previous map[string]loadedEnvironmentValue, key string) string {
	if value, ok := plan.fileValues[key]; ok {
		return value
	}
	if value, ok := plan.defaults[key]; ok {
		return value
	}
	snapshot := originalEnvironmentSnapshot(key, previous)
	if !snapshot.present {
		return ""
	}
	return snapshot.value
}

// environmentPlanInt parses a planned integer without consulting partially applied process state.
func environmentPlanInt(plan environmentLoadPlan, previous map[string]loadedEnvironmentValue, key string) int {
	value := environmentPlanValue(plan, previous, key)
	parsed, _ := strconv.Atoi(value)
	return parsed
}

// applyEnvironmentLoadPlan rolls back every affected variable when any application step fails.
func applyEnvironmentLoadPlan(previous map[string]loadedEnvironmentValue, baseline map[string]environmentSnapshot, plan environmentLoadPlan) (map[string]loadedEnvironmentValue, error) {
	keys := environmentPlanKeys(previous, plan)
	before := make(map[string]environmentSnapshot, len(keys))
	for _, key := range keys {
		value, present := envLookup(key)
		before[key] = environmentSnapshot{value: value, present: present}
	}

	for _, key := range keys {
		target := environmentPlanTarget(key, previous, plan)
		if snapshotsEqual(before[key], target) {
			continue
		}
		if err := writeEnvironmentSnapshot(key, target); err != nil {
			rollbackErr := restoreEnvironmentSnapshots(keys, before)
			applyErr := fmt.Errorf("apply env value %s: %w", key, err)
			if rollbackErr != nil {
				return nil, errors.Join(applyErr, rollbackErr)
			}
			return nil, applyErr
		}
	}

	next := make(map[string]loadedEnvironmentValue, len(plan.fileValues))
	for key := range plan.fileValues {
		if loaded, ok := previous[key]; ok {
			next[key] = loaded
			continue
		}
		next[key] = loadedEnvironmentValue{original: baseline[key]}
	}
	return next, nil
}

// environmentPlanKeys returns affected keys in stable order for deterministic application and rollback.
func environmentPlanKeys(previous map[string]loadedEnvironmentValue, plan environmentLoadPlan) []string {
	set := make(map[string]struct{}, len(previous)+len(plan.fileValues)+len(plan.defaults))
	for key := range previous {
		set[key] = struct{}{}
	}
	for key := range plan.fileValues {
		set[key] = struct{}{}
	}
	for key := range plan.defaults {
		set[key] = struct{}{}
	}
	keys := make([]string, 0, len(set))
	for key := range set {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// environmentPlanTarget resolves a key to a file value, default, or restored ambient snapshot.
func environmentPlanTarget(key string, previous map[string]loadedEnvironmentValue, plan environmentLoadPlan) environmentSnapshot {
	if value, ok := plan.fileValues[key]; ok {
		return environmentSnapshot{value: value, present: true}
	}
	if value, ok := plan.defaults[key]; ok {
		return environmentSnapshot{value: value, present: true}
	}
	if loaded, ok := previous[key]; ok {
		return loaded.original
	}
	return environmentSnapshot{}
}

// restoreEnvironmentSnapshots restores the pre-application process state after a failed plan.
func restoreEnvironmentSnapshots(keys []string, snapshots map[string]environmentSnapshot) error {
	var rollbackErrors []error
	for index := len(keys) - 1; index >= 0; index-- {
		key := keys[index]
		if err := writeEnvironmentSnapshot(key, snapshots[key]); err != nil {
			rollbackErrors = append(rollbackErrors, fmt.Errorf("restore env value %s: %w", key, err))
		}
	}
	if len(rollbackErrors) == 0 {
		return nil
	}
	return fmt.Errorf("roll back process environment: %w", errors.Join(rollbackErrors...))
}

// writeEnvironmentSnapshot preserves unset and explicitly empty values as distinct states.
func writeEnvironmentSnapshot(key string, snapshot environmentSnapshot) error {
	if snapshot.present {
		return envSet(key, snapshot.value)
	}
	return envUnset(key)
}

// snapshotsEqual avoids unnecessary process-wide environment writes.
func snapshotsEqual(left, right environmentSnapshot) bool {
	return left.present == right.present && (!left.present || left.value == right.value)
}

// cloneLoadedEnvironmentValues prevents failed operations from mutating published loader ownership state.
func cloneLoadedEnvironmentValues(values map[string]loadedEnvironmentValue) map[string]loadedEnvironmentValue {
	clone := make(map[string]loadedEnvironmentValue, len(values))
	for key, value := range values {
		clone[key] = value
	}
	return clone
}

// snapshotProcessEnvironment records the exact ambient baseline before the first successful load.
func snapshotProcessEnvironment() map[string]environmentSnapshot {
	snapshots := make(map[string]environmentSnapshot)
	for _, entry := range os.Environ() {
		key, value, _ := strings.Cut(entry, "=")
		snapshots[key] = environmentSnapshot{value: value, present: true}
	}
	return snapshots
}

// cloneEnvironmentSnapshots prevents a failed operation from changing the published baseline.
func cloneEnvironmentSnapshots(snapshots map[string]environmentSnapshot) map[string]environmentSnapshot {
	clone := make(map[string]environmentSnapshot, len(snapshots))
	for key, snapshot := range snapshots {
		clone[key] = snapshot
	}
	return clone
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

// IsEnvLoaded reports whether a Load or Reload completed successfully in this process.
// @group Environment loading
// @behavior readonly
//
// Example:
//
//	env.Dump(env.IsEnvLoaded())
//	// #bool true  (after Load)
//	// #bool false (otherwise)
func IsEnvLoaded() bool {
	processEnvironmentLoader.mu.Lock()
	defer processEnvironmentLoader.mu.Unlock()
	return processEnvironmentLoader.loaded
}

// loadEnvFile returns the nearest parsed regular file without changing the process environment.
func loadEnvFile(startDirectory, name string) (environmentFile, bool, error) {
	path, found, err := findEnvFile(startDirectory, name)
	if err != nil || !found {
		return environmentFile{}, found, err
	}
	values, err := envFileRead(path)
	if err != nil {
		return environmentFile{}, false, fmt.Errorf("read env file %s: %w", path, err)
	}
	return environmentFile{path: path, values: values}, true, nil
}

// findEnvFile performs exactly the documented bounded nearest-ancestor search and follows regular-file symlinks.
func findEnvFile(startDirectory, name string) (string, bool, error) {
	directory := filepath.Clean(startDirectory)
	for level := 0; level < MaxDirectorySeekLevels; level++ {
		candidate := filepath.Join(directory, name)
		info, err := envFileStat(candidate)
		switch {
		case err == nil:
			if !info.Mode().IsRegular() {
				return "", false, fmt.Errorf("env file %s is not a regular file", candidate)
			}
			return candidate, true, nil
		case errors.Is(err, os.ErrNotExist):
			// Missing candidates are the only errors that permit ancestor fallback.
		default:
			return "", false, fmt.Errorf("stat env file %s: %w", candidate, err)
		}

		parent := filepath.Dir(directory)
		if parent == directory {
			break
		}
		directory = parent
	}
	return "", false, nil
}

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

// printLoadedEnvFiles reports filenames and APP_ENV only; file values are intentionally never logged.
func printLoadedEnvFiles(paths []string, appEnv string) {
	for _, path := range paths {
		fmt.Fprintf(os.Stdout, " %s .env file loader · env [%v] file [%v]\n", debugMark(), appEnv, path)
	}
}
