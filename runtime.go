package env

import "runtime"

// goos and goarch are internal shims that allow tests to override the
// detected operating system and architecture.
//
// In production builds, they default to runtime.GOOS and runtime.GOARCH.
// In tests, they can be temporarily replaced to simulate other platforms.
var (
	goos   = runtime.GOOS
	goarch = runtime.GOARCH
)

// OS returns the current operating system identifier.
// @group Runtime
// @behavior readonly
//
// Mirrors runtime.GOOS; tests may override via the internal shim.
//
// Example: inspect GOOS
//
//	env.Dump(env.OS())
//
//	// #string "linux"   (on Linux)
//	// #string "darwin"  (on macOS)
//	// #string "windows" (on Windows)
func OS() string {
	return goos
}

// Arch returns the CPU architecture the binary is running on.
// @group Runtime
// @behavior readonly
//
// Mirrors runtime.GOARCH; tests may override via the internal shim.
//
// Example: print GOARCH
//
//	env.Dump(env.Arch())
//
//	// #string "amd64"
//	// #string "arm64"
func Arch() string {
	return goarch
}

// IsLinux reports whether the runtime OS is Linux.
// @group Runtime
// @behavior readonly
//
// Example:
//
//	env.Dump(env.IsLinux())
//
//	// #bool true  (on Linux)
//	// #bool false (on other OSes)
func IsLinux() bool {
	return goos == "linux"
}

// IsMac reports whether the runtime OS is macOS (Darwin).
// @group Runtime
// @behavior readonly
//
// Example:
//
//	env.Dump(env.IsMac())
//
//	// #bool true  (on macOS)
//	// #bool false (elsewhere)
func IsMac() bool {
	return goos == "darwin"
}

// IsWindows reports whether the runtime OS is Windows.
// @group Runtime
// @behavior readonly
//
// Example:
//
//	env.Dump(env.IsWindows())
//
//	// #bool true  (on Windows)
//	// #bool false (elsewhere)
func IsWindows() bool {
	return goos == "windows"
}

// IsBSD reports whether the runtime OS is any BSD variant.
// @group Runtime
// @behavior readonly
//
// BSD identifiers include: freebsd, openbsd, netbsd, dragonfly.
//
// Example:
//
//	env.Dump(env.IsBSD())
//
//	// #bool true  (on BSD variants)
//	// #bool false (elsewhere)
func IsBSD() bool {
	switch goos {
	case "freebsd", "openbsd", "netbsd", "dragonfly":
		return true
	}
	return false
}

// IsUnix reports whether the OS is Unix-like.
// @group Runtime
// @behavior readonly
//
// Returns true for Linux, macOS, BSD, Solaris, and AIX identifiers.
//
// Example:
//
//	env.Dump(env.IsUnix())
//
//	// #bool true  (on Unix-like OSes)
//	// #bool false (e.g., on Windows or Plan 9)
func IsUnix() bool {
	switch goos {
	case "linux", "darwin", "freebsd", "openbsd", "netbsd", "dragonfly", "solaris", "aix":
		return true
	}
	return false
}

// IsContainerOS reports whether this OS is *typically* used as a container base.
// @group Runtime
// @behavior readonly
//
// This does NOT indicate you are inside a container — only that the OS is usually
// the base for container images (currently Linux).
//
// Example:
//
//	env.Dump(env.IsContainerOS())
//
//	// #bool true  (on Linux)
//	// #bool false (on macOS/Windows)
func IsContainerOS() bool {
	return goos == "linux"
}
