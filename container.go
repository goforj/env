package env

import (
	"os"
)

var (
	// These are shims that tests override.
	statFile = os.Stat
	readFile = os.ReadFile
	getEnv   = os.Getenv
)

const (
	// files
	fileDockerSock = "/var/run/docker.sock"
	fileDockerEnv  = "/.dockerenv"
	fileCgroup     = "/proc/1/cgroup"

	// cgroup names
	cgroupContainer      = "container"
	cgroupNameDocker     = "docker"
	cgroupNameKube       = "kubepods"
	cgroupNameContainerd = "containerd"
	cgroupNamePodman     = "podman"
	cgroupNameLibpod     = "libpod"
)

// IsDocker reports whether the current process is running in a Docker container.
// @group Container detection
// @behavior readonly
//
// Heuristics: presence of /.dockerenv or Docker-related cgroup markers.
//
// Example: typical host
//
//	env.Dump(env.IsDocker())
//
//	// #bool false (unless inside Docker)
func IsDocker() bool {
	// Check /.dockerenv
	if _, err := statFile(fileDockerEnv); err == nil {
		return true
	}

	// Check cgroup
	cgroup, err := readFile(fileCgroup)
	if err == nil && containsAny(cgroup, cgroupNameDocker, cgroupNameContainerd, cgroupNamePodman) {
		return true
	}

	return false
}

// IsDockerInDocker reports whether we are inside a Docker-in-Docker environment.
// @group Container detection
// @behavior readonly
//
// Requires /.dockerenv to be present and a docker.sock exposed to the container.
//
// Example:
//
//	env.Dump(env.IsDockerInDocker())
//
//	// #bool true  (inside DinD containers)
//	// #bool false (on hosts or non-DinD containers)
func IsDockerInDocker() bool {
	// If /.dockerenv does not exist → not a Docker *container* at all.
	if _, err := statFile(fileDockerEnv); err != nil {
		return false
	}

	// If docker.sock exists → this IS an inner DinD container.
	if _, err := statFile(fileDockerSock); err == nil {
		return true
	}

	return false
}

// IsDockerHost reports whether this container behaves like a Docker host.
// @group Container detection
// @behavior readonly
//
// True when docker.sock is available but container-level cgroups are absent.
//
// Example:
//
//	env.Dump(env.IsDockerHost())
//
//	// #bool true  (when acting as Docker host)
//	// #bool false (for normal containers/hosts)
func IsDockerHost() bool {
	if _, err := statFile(fileDockerSock); err != nil {
		return false
	}

	cgroup, err := readFile(fileCgroup)
	if err != nil {
		return false
	}

	// Docker host should *not* have container-scoped cgroups
	if !containsAny(cgroup, cgroupNameDocker, cgroupNameKube, cgroupNameContainerd) {
		return true
	}

	return false
}

// IsContainer detects common container runtimes (Docker, containerd, Kubernetes, Podman).
// @group Container detection
// @behavior readonly
//
// Example: host vs container
//
//	env.Dump(env.IsContainer())
//
//	// #bool true  (inside most containers)
//	// #bool false (on bare-metal/VM hosts)
func IsContainer() bool {
	if IsDocker() {
		return true
	}

	cgroup, err := readFile(fileCgroup)
	if err == nil && containsAny(cgroup,
		cgroupContainer,
		cgroupNameKube,
		cgroupNameLibpod,
		cgroupNameContainerd,
	) {
		return true
	}

	if getEnv("KUBERNETES_SERVICE_HOST") != "" {
		return true
	}

	return false
}

// IsKubernetes reports whether the process is running inside Kubernetes.
// @group Container detection
// @behavior readonly
//
// Checks the KUBERNETES_SERVICE_HOST env var and kubepods cgroup markers.
//
// Example:
//
//	env.Dump(env.IsKubernetes())
//
//	// #bool true  (inside Kubernetes pods)
//	// #bool false (elsewhere)
func IsKubernetes() bool {
	if getEnv("KUBERNETES_SERVICE_HOST") != "" {
		return true
	}

	cgroup, err := readFile(fileCgroup)
	if err == nil && containsAny(cgroup, cgroupNameKube) {
		return true
	}

	return false
}
