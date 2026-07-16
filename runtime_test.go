package env

import "testing"

// resetRuntime restores runtime shims so tests remain isolated.
func resetRuntime() {
	goos = "linux"
	goarch = "amd64"
}

// TestOSAndArch ensures runtime identity helpers expose the active platform values.
func TestOSAndArch(t *testing.T) {
	defer resetRuntime()

	goos = "windows"
	goarch = "arm64"

	if OS() != "windows" {
		t.Fatalf("expected OS() == windows, got %s", OS())
	}

	if Arch() != "arm64" {
		t.Fatalf("expected Arch() == arm64, got %s", Arch())
	}
}

// TestIsLinux ensures Linux detection does not overlap unrelated operating systems.
func TestIsLinux(t *testing.T) {
	defer resetRuntime()

	goos = "linux"
	if !IsLinux() {
		t.Fatalf("expected IsLinux() == true")
	}
	if !IsUnix() {
		t.Fatalf("linux should be unix-like")
	}
	if IsMac() || IsWindows() {
		t.Fatalf("linux should not be mac/windows")
	}
}

// TestIsMac ensures Darwin detection does not overlap unrelated operating systems.
func TestIsMac(t *testing.T) {
	defer resetRuntime()

	goos = "darwin"
	if !IsMac() {
		t.Fatalf("expected IsMac() == true")
	}
	if !IsUnix() {
		t.Fatalf("darwin should be unix-like")
	}
	if IsLinux() || IsWindows() {
		t.Fatalf("mac should not be linux/windows")
	}
}

// TestIsWindows ensures Windows detection does not overlap unrelated operating systems.
func TestIsWindows(t *testing.T) {
	defer resetRuntime()

	goos = "windows"
	if !IsWindows() {
		t.Fatalf("expected IsWindows() == true")
	}
	if IsUnix() {
		t.Fatalf("windows should not be unix-like")
	}
	if IsLinux() || IsMac() {
		t.Fatalf("windows should not be linux/mac")
	}
}

// TestIsBSD ensures every supported BSD runtime is recognized.
func TestIsBSD(t *testing.T) {
	defer resetRuntime()

	for _, bsd := range []string{"freebsd", "openbsd", "netbsd", "dragonfly"} {
		goos = bsd

		if !IsBSD() {
			t.Fatalf("expected IsBSD() == true for %s", bsd)
		}
		if !IsUnix() {
			t.Fatalf("bsd should be unix-like: %s", bsd)
		}
		if IsLinux() || IsWindows() || IsMac() {
			t.Fatalf("bsd should not be linux/windows/mac: %s", bsd)
		}
	}
}

// TestIsBSD_False ensures non-BSD Unix systems are not misclassified.
func TestIsBSD_False(t *testing.T) {
	defer resetRuntime()

	goos = "linux"
	if IsBSD() {
		t.Fatalf("expected IsBSD() == false for linux")
	}
}

// TestIsUnix ensures the documented Unix family is recognized without including Windows.
func TestIsUnix(t *testing.T) {
	defer resetRuntime()

	unixSystems := []string{
		"linux", "darwin", "freebsd", "openbsd", "netbsd", "dragonfly",
		"solaris", "aix",
	}

	for _, sys := range unixSystems {
		goos = sys
		if !IsUnix() {
			t.Fatalf("expected IsUnix() == true for %s", sys)
		}
	}

	nonUnix := []string{"windows", "plan9"}
	for _, sys := range nonUnix {
		goos = sys
		if IsUnix() {
			t.Fatalf("expected IsUnix() == false for %s", sys)
		}
	}
}

// TestIsContainerOS ensures container-oriented runtime identifiers remain distinguishable from host kernels.
func TestIsContainerOS(t *testing.T) {
	defer resetRuntime()

	goos = "linux"
	if !IsContainerOS() {
		t.Fatalf("expected linux to be container OS")
	}

	goos = "darwin"
	if IsContainerOS() {
		t.Fatalf("expected macOS NOT to be container OS")
	}

	goos = "windows"
	if IsContainerOS() {
		t.Fatalf("expected windows NOT to be container OS")
	}
}
