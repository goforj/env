package env

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

// prepareLoaderTest isolates process-wide loader state and shims for one test.
func prepareLoaderTest(t *testing.T, keys ...string) {
	t.Helper()
	restoreEnvironment := snapshotEnv(append(keys, "APP_ENV"))
	originalDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}

	originalGetwd := envFileGetwd
	originalStat := envFileStat
	originalRead := envFileRead
	originalLookup := envLookup
	originalSet := envSet
	originalUnset := envUnset
	originalStatFile := statFile
	originalReadFile := readFile
	originalGetEnv := getEnv

	processEnvironmentLoader.mu.Lock()
	originalLoaded := processEnvironmentLoader.loaded
	originalValues := cloneLoadedEnvironmentValues(processEnvironmentLoader.values)
	originalBaseline := cloneEnvironmentSnapshots(processEnvironmentLoader.baseline)
	processEnvironmentLoader.loaded = false
	processEnvironmentLoader.values = make(map[string]loadedEnvironmentValue)
	processEnvironmentLoader.baseline = make(map[string]environmentSnapshot)
	processEnvironmentLoader.mu.Unlock()

	t.Cleanup(func() {
		envFileGetwd = originalGetwd
		envFileStat = originalStat
		envFileRead = originalRead
		envLookup = originalLookup
		envSet = originalSet
		envUnset = originalUnset
		statFile = originalStatFile
		readFile = originalReadFile
		getEnv = originalGetEnv
		_ = os.Chdir(originalDirectory)
		restoreEnvironment()

		processEnvironmentLoader.mu.Lock()
		processEnvironmentLoader.loaded = originalLoaded
		processEnvironmentLoader.values = originalValues
		processEnvironmentLoader.baseline = originalBaseline
		processEnvironmentLoader.mu.Unlock()
	})
}

// changeWorkingDirectory moves the process into directory until the test cleanup runs.
func changeWorkingDirectory(t *testing.T, directory string) {
	t.Helper()
	changeWorkingDirectoryWithin(t, directory, directory)
}

// changeWorkingDirectoryWithin prevents unrelated ancestor files from contaminating loader tests.
func changeWorkingDirectoryWithin(t *testing.T, directory, searchRoot string) {
	t.Helper()
	if err := os.Chdir(directory); err != nil {
		t.Fatalf("change working directory: %v", err)
	}
	root := filepath.Clean(searchRoot)
	stat := envFileStat
	envFileStat = func(path string) (os.FileInfo, error) {
		relative, err := filepath.Rel(root, filepath.Clean(path))
		if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
			return nil, os.ErrNotExist
		}
		return stat(path)
	}
}

// writeEnvFile creates a test env file with predictable permissions.
func writeEnvFile(t *testing.T, directory, name, contents string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(directory, name), []byte(contents), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

// TestLoadAppliesLayersInOrder ensures later environment-specific files override earlier base values.
func TestLoadAppliesLayersInOrder(t *testing.T) {
	prepareLoaderTest(t, "ENV_QPASS_SHARED", "ENV_QPASS_BASE", "ENV_QPASS_LAYER", "ENV_QPASS_TEST")
	t.Setenv("APP_ENV", Staging)
	t.Setenv("ENV_QPASS_SHARED", "ambient")
	directory := t.TempDir()
	writeEnvFile(t, directory, fileEnv, "ENV_QPASS_BASE=base\nENV_QPASS_SHARED=base\n")
	writeEnvFile(t, directory, envFileStaging, "ENV_QPASS_LAYER=staging\nENV_QPASS_SHARED=staging\n")
	writeEnvFile(t, directory, envFileTesting, "ENV_QPASS_TEST=test\nENV_QPASS_SHARED=testing\n")
	changeWorkingDirectory(t, directory)

	if err := Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}

	want := map[string]string{
		"ENV_QPASS_BASE":   "base",
		"ENV_QPASS_LAYER":  "staging",
		"ENV_QPASS_TEST":   "test",
		"ENV_QPASS_SHARED": "testing",
	}
	for key, expected := range want {
		if got := os.Getenv(key); got != expected {
			t.Fatalf("expected %s=%q, got %q", key, expected, got)
		}
	}
	if !IsEnvLoaded() {
		t.Fatal("expected successful Load to publish loaded state")
	}
}

// TestLoadBaseSelectsApplicationLayer ensures APP_ENV selects the matching application-specific file.
func TestLoadBaseSelectsApplicationLayer(t *testing.T) {
	prepareLoaderTest(t, "ENV_QPASS_LAYER")
	t.Setenv("APP_ENV", Local)
	directory := t.TempDir()
	writeEnvFile(t, directory, fileEnv, "APP_ENV=production\n")
	writeEnvFile(t, directory, envFileProd, "ENV_QPASS_LAYER=production\n")
	changeWorkingDirectory(t, directory)

	if err := Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := os.Getenv("ENV_QPASS_LAYER"); got != Production {
		t.Fatalf("expected base APP_ENV to select production, got %q", got)
	}
}

// TestLoadSearchesEachLayerIndependently ensures each filename uses its own nearest ancestor match.
func TestLoadSearchesEachLayerIndependently(t *testing.T) {
	prepareLoaderTest(t, "ENV_QPASS_BASE", "ENV_QPASS_LAYER")
	t.Setenv("APP_ENV", Local)
	parent := t.TempDir()
	child := filepath.Join(parent, "child")
	if err := os.Mkdir(child, 0o755); err != nil {
		t.Fatalf("make child: %v", err)
	}
	writeEnvFile(t, parent, fileEnv, "ENV_QPASS_BASE=parent\n")
	writeEnvFile(t, child, envFileLocal, "ENV_QPASS_LAYER=child\n")
	changeWorkingDirectoryWithin(t, child, parent)

	if err := Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if os.Getenv("ENV_QPASS_BASE") != "parent" || os.Getenv("ENV_QPASS_LAYER") != "child" {
		t.Fatalf("expected independent nearest-file lookup")
	}
}

// TestLoadAppliesHostLayer ensures host-specific values participate at their documented precedence.
func TestLoadAppliesHostLayer(t *testing.T) {
	prepareLoaderTest(t, "ENV_QPASS_HOST")
	t.Setenv("APP_ENV", Local)
	directory := t.TempDir()
	writeEnvFile(t, directory, fileEnvHost, "ENV_QPASS_HOST=host\n")
	changeWorkingDirectory(t, directory)
	statFile = func(string) (os.FileInfo, error) { return nil, os.ErrNotExist }
	readFile = func(string) ([]byte, error) { return []byte("0::/user.slice"), nil }

	if err := Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := os.Getenv("ENV_QPASS_HOST"); got != "host" {
		t.Fatalf("expected host layer, got %q", got)
	}
}

// TestLoadDefaultsAppEnvWithoutOwningIt ensures a synthesized development mode remains caller-owned state.
func TestLoadDefaultsAppEnvWithoutOwningIt(t *testing.T) {
	prepareLoaderTest(t, "ENV_QPASS_PRODUCTION")
	_ = os.Unsetenv("APP_ENV")
	directory := t.TempDir()
	writeEnvFile(t, directory, envFileProd, "ENV_QPASS_PRODUCTION=yes\n")
	changeWorkingDirectory(t, directory)

	if err := Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := os.Getenv("APP_ENV"); got != Local {
		t.Fatalf("expected default APP_ENV=%q, got %q", Local, got)
	}

	t.Setenv("APP_ENV", Production)
	if err := Reload(); err != nil {
		t.Fatalf("Reload: %v", err)
	}
	if got := os.Getenv("ENV_QPASS_PRODUCTION"); got != "yes" {
		t.Fatalf("expected caller APP_ENV to select production, got %q", got)
	}
}

// TestLoadIsIdempotent ensures repeated loads do not drift the effective environment.
func TestLoadIsIdempotent(t *testing.T) {
	prepareLoaderTest(t, "ENV_QPASS_IDEMPOTENT")
	t.Setenv("APP_ENV", Local)
	directory := t.TempDir()
	writeEnvFile(t, directory, fileEnv, "ENV_QPASS_IDEMPOTENT=file\n")
	changeWorkingDirectory(t, directory)

	if err := Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	t.Setenv("ENV_QPASS_IDEMPOTENT", "runtime")
	if err := Load(); err != nil {
		t.Fatalf("second Load: %v", err)
	}
	if got := os.Getenv("ENV_QPASS_IDEMPOTENT"); got != "runtime" {
		t.Fatalf("expected repeated Load to be a no-op, got %q", got)
	}
}

// TestReloadReappliesAndRestoresFileOwnedKeys ensures removed file values return to their pre-load baseline.
func TestReloadReappliesAndRestoresFileOwnedKeys(t *testing.T) {
	prepareLoaderTest(t, "ENV_QPASS_RELOAD", "ENV_QPASS_RESTORE", "ENV_QPASS_ABSENT", "ENV_QPASS_UNRELATED")
	t.Setenv("APP_ENV", Local)
	t.Setenv("ENV_QPASS_RESTORE", "ambient")
	_ = os.Unsetenv("ENV_QPASS_ABSENT")
	t.Setenv("ENV_QPASS_UNRELATED", "before")
	directory := t.TempDir()
	writeEnvFile(t, directory, fileEnv, "ENV_QPASS_RELOAD=first\nENV_QPASS_RESTORE=file\nENV_QPASS_ABSENT=file\n")
	changeWorkingDirectory(t, directory)

	if err := Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	t.Setenv("ENV_QPASS_RELOAD", "runtime")
	t.Setenv("ENV_QPASS_UNRELATED", "after")
	writeEnvFile(t, directory, fileEnv, "ENV_QPASS_RELOAD=second\n")
	if err := Reload(); err != nil {
		t.Fatalf("Reload: %v", err)
	}

	if got := os.Getenv("ENV_QPASS_RELOAD"); got != "second" {
		t.Fatalf("expected file to replace runtime edit, got %q", got)
	}
	if got := os.Getenv("ENV_QPASS_RESTORE"); got != "ambient" {
		t.Fatalf("expected original ambient restoration, got %q", got)
	}
	if _, present := os.LookupEnv("ENV_QPASS_ABSENT"); present {
		t.Fatal("expected originally absent key to become absent again")
	}
	if got := os.Getenv("ENV_QPASS_UNRELATED"); got != "after" {
		t.Fatalf("expected unrelated key to remain untouched, got %q", got)
	}
}

// TestReloadRestoresBaselineForKeysClaimedLater ensures newly managed keys retain their original ambient value.
func TestReloadRestoresBaselineForKeysClaimedLater(t *testing.T) {
	prepareLoaderTest(t, "ENV_QPASS_LATE_CLAIM")
	t.Setenv("APP_ENV", Local)
	_ = os.Unsetenv("ENV_QPASS_LATE_CLAIM")
	directory := t.TempDir()
	changeWorkingDirectory(t, directory)

	if err := Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	t.Setenv("ENV_QPASS_LATE_CLAIM", "runtime-after-load")
	writeEnvFile(t, directory, fileEnv, "ENV_QPASS_LATE_CLAIM=file\n")
	if err := Reload(); err != nil {
		t.Fatalf("Reload claiming key: %v", err)
	}
	writeEnvFile(t, directory, fileEnv, "")
	if err := Reload(); err != nil {
		t.Fatalf("Reload removing key: %v", err)
	}
	if _, present := os.LookupEnv("ENV_QPASS_LATE_CLAIM"); present {
		t.Fatal("expected a later-claimed key to restore the pre-first-load absent state")
	}
}

// TestReloadRefreshesFileOwnedAppEnv ensures a file-controlled APP_ENV can select a new layer transactionally.
func TestReloadRefreshesFileOwnedAppEnv(t *testing.T) {
	prepareLoaderTest(t, "ENV_QPASS_STAGE", "ENV_QPASS_PROD")
	t.Setenv("APP_ENV", Local)
	directory := t.TempDir()
	writeEnvFile(t, directory, fileEnv, "APP_ENV=staging\n")
	writeEnvFile(t, directory, envFileStaging, "ENV_QPASS_STAGE=yes\n")
	writeEnvFile(t, directory, envFileProd, "ENV_QPASS_PROD=yes\n")
	changeWorkingDirectory(t, directory)

	if err := Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	writeEnvFile(t, directory, fileEnv, "APP_ENV=production\n")
	if err := Reload(); err != nil {
		t.Fatalf("Reload: %v", err)
	}
	if got := os.Getenv("APP_ENV"); got != Production {
		t.Fatalf("expected refreshed APP_ENV, got %q", got)
	}
	if got := os.Getenv("ENV_QPASS_PROD"); got != "yes" {
		t.Fatalf("expected refreshed APP_ENV to select production, got %q", got)
	}
	if _, present := os.LookupEnv("ENV_QPASS_STAGE"); present {
		t.Fatal("expected staging-only key to be restored")
	}

	writeEnvFile(t, directory, fileEnv, "")
	if err := Reload(); err != nil {
		t.Fatalf("Reload removing file-owned APP_ENV: %v", err)
	}
	if got := os.Getenv("APP_ENV"); got != Local {
		t.Fatalf("expected removed file-owned APP_ENV to restore caller value, got %q", got)
	}
}

// TestLoadReturnsDiscoveryAndParseErrorsWithoutMutation ensures failed discovery or parsing cannot partially change process state.
func TestLoadReturnsDiscoveryAndParseErrorsWithoutMutation(t *testing.T) {
	t.Run("working directory error", func(t *testing.T) {
		prepareLoaderTest(t, "ENV_QPASS_UNCHANGED")
		t.Setenv("ENV_QPASS_UNCHANGED", "ambient")
		workingDirectoryErr := errors.New("injected working directory failure")
		envFileGetwd = func() (string, error) {
			return "", workingDirectoryErr
		}

		if err := Load(); !errors.Is(err, workingDirectoryErr) {
			t.Fatalf("expected working directory error, got %v", err)
		}
		if os.Getenv("ENV_QPASS_UNCHANGED") != "ambient" || IsEnvLoaded() {
			t.Fatal("expected failed working directory lookup to leave environment and state unchanged")
		}
	})

	t.Run("stat error", func(t *testing.T) {
		prepareLoaderTest(t, "ENV_QPASS_UNCHANGED")
		t.Setenv("APP_ENV", Local)
		t.Setenv("ENV_QPASS_UNCHANGED", "ambient")
		directory := t.TempDir()
		changeWorkingDirectory(t, directory)
		envFileStat = func(path string) (os.FileInfo, error) {
			if filepath.Base(path) == fileEnv {
				return nil, os.ErrPermission
			}
			return nil, os.ErrNotExist
		}

		if err := Load(); !errors.Is(err, os.ErrPermission) {
			t.Fatalf("expected permission error, got %v", err)
		}
		if os.Getenv("ENV_QPASS_UNCHANGED") != "ambient" || IsEnvLoaded() {
			t.Fatal("expected failed discovery to leave environment and state unchanged")
		}
	})

	t.Run("application layer stat error", func(t *testing.T) {
		prepareLoaderTest(t, "ENV_QPASS_UNCHANGED")
		t.Setenv("APP_ENV", Local)
		t.Setenv("ENV_QPASS_UNCHANGED", "ambient")
		directory := t.TempDir()
		changeWorkingDirectory(t, directory)
		applicationLayerErr := errors.New("injected application layer failure")
		envFileStat = func(path string) (os.FileInfo, error) {
			if filepath.Base(path) == envFileLocal {
				return nil, applicationLayerErr
			}
			return nil, os.ErrNotExist
		}

		if err := Load(); !errors.Is(err, applicationLayerErr) {
			t.Fatalf("expected application layer error, got %v", err)
		}
		if os.Getenv("ENV_QPASS_UNCHANGED") != "ambient" || IsEnvLoaded() {
			t.Fatal("expected failed application layer lookup to leave environment and state unchanged")
		}
	})

	t.Run("host layer stat error", func(t *testing.T) {
		prepareLoaderTest(t, "ENV_QPASS_UNCHANGED")
		t.Setenv("APP_ENV", Local)
		t.Setenv("ENV_QPASS_UNCHANGED", "ambient")
		directory := t.TempDir()
		changeWorkingDirectory(t, directory)
		hostLayerErr := errors.New("injected host layer failure")
		envFileStat = func(path string) (os.FileInfo, error) {
			if filepath.Base(path) == fileEnvHost {
				return nil, hostLayerErr
			}
			return nil, os.ErrNotExist
		}
		statFile = func(string) (os.FileInfo, error) {
			return nil, os.ErrNotExist
		}
		readFile = func(string) ([]byte, error) {
			return nil, os.ErrNotExist
		}

		if err := Load(); !errors.Is(err, hostLayerErr) {
			t.Fatalf("expected host layer error, got %v", err)
		}
		if os.Getenv("ENV_QPASS_UNCHANGED") != "ambient" || IsEnvLoaded() {
			t.Fatal("expected failed host layer lookup to leave environment and state unchanged")
		}
	})

	t.Run("testing layer stat error", func(t *testing.T) {
		prepareLoaderTest(t, "ENV_QPASS_UNCHANGED")
		t.Setenv("APP_ENV", Local)
		t.Setenv("ENV_QPASS_UNCHANGED", "ambient")
		directory := t.TempDir()
		changeWorkingDirectory(t, directory)
		testingLayerErr := errors.New("injected testing layer failure")
		envFileStat = func(path string) (os.FileInfo, error) {
			if filepath.Base(path) == envFileTesting {
				return nil, testingLayerErr
			}
			return nil, os.ErrNotExist
		}
		statFile = func(path string) (os.FileInfo, error) {
			if path == fileDockerEnv {
				return nil, nil
			}
			return nil, os.ErrNotExist
		}

		if err := Load(); !errors.Is(err, testingLayerErr) {
			t.Fatalf("expected testing layer error, got %v", err)
		}
		if os.Getenv("ENV_QPASS_UNCHANGED") != "ambient" || IsEnvLoaded() {
			t.Fatal("expected failed testing layer lookup to leave environment and state unchanged")
		}
	})

	t.Run("parse error", func(t *testing.T) {
		prepareLoaderTest(t, "ENV_QPASS_UNCHANGED")
		t.Setenv("APP_ENV", Local)
		t.Setenv("ENV_QPASS_UNCHANGED", "ambient")
		directory := t.TempDir()
		writeEnvFile(t, directory, fileEnv, "ENV_QPASS_UNCHANGED='unterminated\n")
		changeWorkingDirectory(t, directory)

		if err := Load(); err == nil {
			t.Fatal("expected malformed env file error")
		}
		if os.Getenv("ENV_QPASS_UNCHANGED") != "ambient" || IsEnvLoaded() {
			t.Fatal("expected failed parse to leave environment and state unchanged")
		}
	})

	t.Run("non-regular file", func(t *testing.T) {
		prepareLoaderTest(t)
		t.Setenv("APP_ENV", Local)
		directory := t.TempDir()
		if err := os.Mkdir(filepath.Join(directory, fileEnv), 0o755); err != nil {
			t.Fatalf("make env directory: %v", err)
		}
		changeWorkingDirectory(t, directory)

		if err := Load(); err == nil || !strings.Contains(err.Error(), "not a regular file") {
			t.Fatalf("expected regular-file error, got %v", err)
		}
	})
}

// TestLoadRollsBackApplicationFailure ensures an application-layer failure restores every earlier mutation.
func TestLoadRollsBackApplicationFailure(t *testing.T) {
	prepareLoaderTest(t, "ENV_QPASS_A", "ENV_QPASS_B")
	t.Setenv("APP_ENV", Local)
	t.Setenv("ENV_QPASS_A", "ambient-a")
	t.Setenv("ENV_QPASS_B", "ambient-b")
	directory := t.TempDir()
	writeEnvFile(t, directory, fileEnv, "ENV_QPASS_A=file-a\nENV_QPASS_B=file-b\n")
	changeWorkingDirectory(t, directory)

	failed := false
	envSet = func(key, value string) error {
		if key == "ENV_QPASS_B" && value == "file-b" && !failed {
			failed = true
			return errors.New("injected set failure")
		}
		return os.Setenv(key, value)
	}
	if err := Load(); err == nil {
		t.Fatal("expected application error")
	}
	if os.Getenv("ENV_QPASS_A") != "ambient-a" || os.Getenv("ENV_QPASS_B") != "ambient-b" {
		t.Fatal("expected failed Load to roll back all affected keys")
	}
	if IsEnvLoaded() {
		t.Fatal("expected failed Load to leave state unpublished")
	}
}

// TestLoadJoinsApplicationAndRollbackFailures ensures callers retain both causes when recovery also fails.
func TestLoadJoinsApplicationAndRollbackFailures(t *testing.T) {
	prepareLoaderTest(t, "ENV_QPASS_A", "ENV_QPASS_B")
	t.Setenv("APP_ENV", Local)
	t.Setenv("ENV_QPASS_A", "ambient-a")
	t.Setenv("ENV_QPASS_B", "ambient-b")
	directory := t.TempDir()
	writeEnvFile(t, directory, fileEnv, "ENV_QPASS_A=file-a\nENV_QPASS_B=file-b\n")
	changeWorkingDirectory(t, directory)

	applyErr := errors.New("injected apply failure")
	rollbackErr := errors.New("injected rollback failure")
	envSet = func(key, value string) error {
		if key == "ENV_QPASS_B" {
			switch value {
			case "file-b":
				return applyErr
			case "ambient-b":
				return rollbackErr
			}
		}
		return os.Setenv(key, value)
	}

	err := Load()
	if !errors.Is(err, applyErr) || !errors.Is(err, rollbackErr) {
		t.Fatalf("expected joined apply and rollback errors, got %v", err)
	}
	if IsEnvLoaded() {
		t.Fatal("expected failed Load to leave state unpublished")
	}
}

// TestReloadFailurePreservesPreviousConfiguration ensures a failed refresh leaves the last valid configuration active.
func TestReloadFailurePreservesPreviousConfiguration(t *testing.T) {
	prepareLoaderTest(t, "ENV_QPASS_A", "ENV_QPASS_B")
	t.Setenv("APP_ENV", Local)
	directory := t.TempDir()
	writeEnvFile(t, directory, fileEnv, "ENV_QPASS_A=old-a\nENV_QPASS_B=old-b\n")
	changeWorkingDirectory(t, directory)
	if err := Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}

	writeEnvFile(t, directory, fileEnv, "ENV_QPASS_A=new-a\nENV_QPASS_B=new-b\n")
	failed := false
	envSet = func(key, value string) error {
		if key == "ENV_QPASS_B" && value == "new-b" && !failed {
			failed = true
			return errors.New("injected reload failure")
		}
		return os.Setenv(key, value)
	}
	if err := Reload(); err == nil {
		t.Fatal("expected Reload application error")
	}
	if os.Getenv("ENV_QPASS_A") != "old-a" || os.Getenv("ENV_QPASS_B") != "old-b" {
		t.Fatal("expected failed Reload to preserve previous successful configuration")
	}
	if !IsEnvLoaded() {
		t.Fatal("expected previous successful loaded state to remain published")
	}

	envSet = os.Setenv
	if err := Reload(); err != nil {
		t.Fatalf("retry Reload: %v", err)
	}
	if os.Getenv("ENV_QPASS_A") != "new-a" || os.Getenv("ENV_QPASS_B") != "new-b" {
		t.Fatal("expected retry to apply new configuration")
	}
}

// TestFindEnvFileUsesBoundedNearestSearch ensures discovery prefers proximity and stops at the documented ancestor limit.
func TestFindEnvFileUsesBoundedNearestSearch(t *testing.T) {
	t.Run("finds ninth ancestor", func(t *testing.T) {
		root := t.TempDir()
		writeEnvFile(t, root, fileEnv, "A=1\n")
		start := root
		for index := 0; index < MaxDirectorySeekLevels-1; index++ {
			start = filepath.Join(start, "child")
			if err := os.Mkdir(start, 0o755); err != nil {
				t.Fatalf("make nested directory: %v", err)
			}
		}
		path, found, err := findEnvFile(start, fileEnv)
		if err != nil || !found || path != filepath.Join(root, fileEnv) {
			t.Fatalf("expected ninth-ancestor file, got path=%q found=%v err=%v", path, found, err)
		}
	})

	t.Run("does not inspect tenth ancestor", func(t *testing.T) {
		root := t.TempDir()
		writeEnvFile(t, root, fileEnv, "A=1\n")
		start := root
		for index := 0; index < MaxDirectorySeekLevels; index++ {
			start = filepath.Join(start, "child")
			if err := os.Mkdir(start, 0o755); err != nil {
				t.Fatalf("make nested directory: %v", err)
			}
		}
		if path, found, err := findEnvFile(start, fileEnv); err != nil || found || path != "" {
			t.Fatalf("expected bounded miss, got path=%q found=%v err=%v", path, found, err)
		}
	})
}

// TestLoadFollowsRegularFileSymlink ensures deployed symlinked dotenv files remain valid when their targets are regular files.
func TestLoadFollowsRegularFileSymlink(t *testing.T) {
	prepareLoaderTest(t, "ENV_QPASS_SYMLINK")
	t.Setenv("APP_ENV", Local)
	directory := t.TempDir()
	target := filepath.Join(directory, "values.env")
	if err := os.WriteFile(target, []byte("ENV_QPASS_SYMLINK=yes\n"), 0o644); err != nil {
		t.Fatalf("write symlink target: %v", err)
	}
	if err := os.Symlink(target, filepath.Join(directory, fileEnv)); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	changeWorkingDirectory(t, directory)

	if err := Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := os.Getenv("ENV_QPASS_SYMLINK"); got != "yes" {
		t.Fatalf("expected symlinked env value, got %q", got)
	}
}

// TestLoadReloadAndStateAreConcurrentSafe ensures process-wide environment ownership remains coherent under concurrent access.
func TestLoadReloadAndStateAreConcurrentSafe(t *testing.T) {
	prepareLoaderTest(t, "ENV_QPASS_CONCURRENT")
	t.Setenv("APP_ENV", Local)
	directory := t.TempDir()
	writeEnvFile(t, directory, fileEnv, "ENV_QPASS_CONCURRENT=value\n")
	changeWorkingDirectory(t, directory)

	const workers = 60
	errorsFound := make(chan error, workers)
	done := make(chan struct{})
	go func() {
		var wait sync.WaitGroup
		for index := 0; index < workers; index++ {
			wait.Add(1)
			go func(index int) {
				defer wait.Done()
				var err error
				if index%3 == 0 {
					err = Reload()
				} else {
					err = Load()
				}
				_ = IsEnvLoaded()
				if err != nil {
					errorsFound <- err
				}
			}(index)
		}
		wait.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("concurrent loader operations exceeded hard deadline")
	}
	close(errorsFound)
	for err := range errorsFound {
		t.Fatalf("concurrent loader operation: %v", err)
	}
	if got := os.Getenv("ENV_QPASS_CONCURRENT"); got != "value" {
		t.Fatalf("expected stable concurrent value, got %q", got)
	}
}

// TestLoadDebugOutputDoesNotExposeValues ensures diagnostics cannot disclose loaded secrets.
func TestLoadDebugOutputDoesNotExposeValues(t *testing.T) {
	prepareLoaderTest(t, "ENV_QPASS_SECRET")
	t.Setenv("APP_ENV", Local)
	directory := t.TempDir()
	writeEnvFile(t, directory, fileEnv, "ENV_DEBUG=3\nENV_QPASS_SECRET=do-not-print-me\n")
	changeWorkingDirectory(t, directory)

	output := captureStdout(t, func() {
		if err := Load(); err != nil {
			t.Fatalf("Load: %v", err)
		}
	})
	if strings.Contains(output, "do-not-print-me") || strings.Contains(output, "ENV_QPASS_SECRET") {
		t.Fatalf("debug output exposed env data: %q", output)
	}
	if !strings.Contains(output, "env [local]") || !strings.Contains(output, filepath.Join(directory, fileEnv)) {
		t.Fatalf("expected paths and APP_ENV in debug output, got %q", output)
	}
}

// TestLoadEnvFileIfExistsAliasesLoad ensures the compatibility entry point retains Load semantics.
func TestLoadEnvFileIfExistsAliasesLoad(t *testing.T) {
	prepareLoaderTest(t, "ENV_QPASS_ALIAS")
	t.Setenv("APP_ENV", Local)
	directory := t.TempDir()
	writeEnvFile(t, directory, fileEnv, "ENV_QPASS_ALIAS=yes\n")
	changeWorkingDirectory(t, directory)

	if err := LoadEnvFileIfExists(); err != nil {
		t.Fatalf("LoadEnvFileIfExists: %v", err)
	}
	if got := os.Getenv("ENV_QPASS_ALIAS"); got != "yes" {
		t.Fatalf("expected alias to load value, got %q", got)
	}
}

// TestEnvFileForAppEnv ensures application modes map to stable dotenv filenames.
func TestEnvFileForAppEnv(t *testing.T) {
	cases := []struct {
		appEnv string
		file   string
		found  bool
	}{
		{appEnv: Local, file: envFileLocal, found: true},
		{appEnv: Staging, file: envFileStaging, found: true},
		{appEnv: Production, file: envFileProd, found: true},
		{appEnv: Testing},
		{appEnv: "unknown"},
	}
	for _, test := range cases {
		got, found := envFileForAppEnv(test.appEnv)
		if got != test.file || found != test.found {
			t.Fatalf("envFileForAppEnv(%q) = %q, %v; want %q, %v", test.appEnv, got, found, test.file, test.found)
		}
	}
}

// TestEnvironmentPlanKeysAreDeterministic ensures transactional application order is reproducible.
func TestEnvironmentPlanKeysAreDeterministic(t *testing.T) {
	previous := map[string]loadedEnvironmentValue{"B": {}, "A": {}}
	plan := environmentLoadPlan{fileValues: map[string]string{"C": ""}, defaults: map[string]string{"D": ""}}
	if got, want := environmentPlanKeys(previous, plan), []string{"A", "B", "C", "D"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("expected keys %v, got %v", want, got)
	}
}

// TestEnvironmentPlanTargetReturnsAbsentForUnownedKey guards the helper's safe fallback for direct callers.
func TestEnvironmentPlanTargetReturnsAbsentForUnownedKey(t *testing.T) {
	plan := environmentLoadPlan{fileValues: map[string]string{}, defaults: map[string]string{}}
	if got := environmentPlanTarget("UNOWNED", nil, plan); got.present || got.value != "" {
		t.Fatalf("expected absent snapshot, got %+v", got)
	}
}

// captureStdout captures process output for loader diagnostics tests.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	original := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = writer
	t.Cleanup(func() { os.Stdout = original })

	done := make(chan string)
	go func() {
		var output strings.Builder
		_, _ = io.Copy(&output, reader)
		done <- output.String()
	}()

	fn()
	_ = writer.Close()
	output := <-done
	_ = reader.Close()
	os.Stdout = original
	return output
}
