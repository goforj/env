package env

import (
	"os"
	"testing"
)

// Save real implementations so we can restore them.
var (
	realStatFile = statFile
	realReadFile = readFile
	realGetEnv   = getEnv
)

// reset restores process-detection shims so tests remain isolated.
func reset() {
	statFile = realStatFile
	readFile = realReadFile
	getEnv = realGetEnv
}

// mockEnv isolates environment lookup from the host running the tests.
func mockEnv(vars map[string]string) {
	getEnv = func(key string) string {
		return vars[key]
	}
}

// TestIsDocker_DockerenvExists ensures Docker's root marker is sufficient evidence of containment.
func TestIsDocker_DockerenvExists(t *testing.T) {
	defer reset()

	statFile = func(path string) (os.FileInfo, error) {
		if path == "/.dockerenv" {
			return nil, nil // simulate exists
		}
		return nil, os.ErrNotExist
	}

	readFile = func(path string) ([]byte, error) {
		return nil, os.ErrNotExist
	}

	if !IsDocker() {
		t.Fatal("expected IsDocker() == true when /.dockerenv exists")
	}
}

// TestIsDocker_CgroupDetectsContainer ensures Docker cgroup membership works when the marker file is absent.
func TestIsDocker_CgroupDetectsContainer(t *testing.T) {
	defer reset()

	statFile = func(path string) (os.FileInfo, error) {
		return nil, os.ErrNotExist // no /.dockerenv
	}

	readFile = func(path string) ([]byte, error) {
		return []byte("0::/docker/12345"), nil
	}

	if !IsDocker() {
		t.Fatal("expected IsDocker() == true from docker cgroup")
	}
}

// TestIsDocker_False ensures clean host evidence is not misclassified as Docker.
func TestIsDocker_False(t *testing.T) {
	defer reset()

	statFile = func(path string) (os.FileInfo, error) {
		return nil, os.ErrNotExist
	}

	readFile = func(path string) ([]byte, error) {
		return []byte("0::/user.slice"), nil
	}

	if IsDocker() {
		t.Fatal("expected IsDocker() == false when no signals")
	}
}

// TestIsDind_True ensures a mounted Docker socket inside Docker identifies nested daemon access.
func TestIsDind_True(t *testing.T) {
	defer reset()

	statFile = func(path string) (os.FileInfo, error) {
		switch path {
		case fileDockerEnv:
			return nil, nil // must have /.dockerenv
		case fileDockerSock:
			return nil, nil // must have docker.sock
		}
		return nil, os.ErrNotExist
	}

	// DinD inner often has host-like cgroups
	readFile = func(path string) ([]byte, error) {
		return []byte("0::/user.slice"), nil
	}

	if !IsDockerInDocker() {
		t.Fatal("expected IsDockerInDocker() == true for DinD inner")
	}
}

// TestIsDind_FalseWhenNormalContainer ensures ordinary containers without a daemon socket are not marked Docker-in-Docker.
func TestIsDind_FalseWhenNormalContainer(t *testing.T) {
	defer reset()

	statFile = func(path string) (os.FileInfo, error) {
		if path == "/var/run/docker.sock" {
			return nil, nil
		}
		return nil, os.ErrNotExist
	}

	readFile = func(path string) ([]byte, error) {
		return []byte("0::/docker/12345"), nil
	}

	if IsDockerInDocker() {
		t.Fatal("expected IsDockerInDocker() == false in real container cgroup")
	}
}

// TestIsDockerInDocker_NoDockerenv ensures a host socket alone does not imply nested Docker.
func TestIsDockerInDocker_NoDockerenv(t *testing.T) {
	defer reset()

	statFile = func(path string) (os.FileInfo, error) { return nil, os.ErrNotExist }

	if IsDockerInDocker() {
		t.Fatal("expected false when /.dockerenv missing")
	}
}

// TestIsDockerInDocker_DockerenvNoSocket ensures a container marker alone does not imply nested daemon access.
func TestIsDockerInDocker_DockerenvNoSocket(t *testing.T) {
	defer reset()

	statFile = func(path string) (os.FileInfo, error) {
		if path == fileDockerEnv {
			return nil, nil
		}
		return nil, os.ErrNotExist
	}

	if IsDockerInDocker() {
		t.Fatal("expected false when /.dockerenv exists but docker.sock missing")
	}
}

// TestIsDockerHost_True ensures a reachable Docker socket with host cgroups identifies a Docker host.
func TestIsDockerHost_True(t *testing.T) {
	defer reset()

	statFile = func(path string) (os.FileInfo, error) {
		if path == "/var/run/docker.sock" {
			return nil, nil
		}
		return nil, os.ErrNotExist
	}

	readFile = func(path string) ([]byte, error) {
		return []byte("0::/user.slice"), nil
	}

	if !IsDockerHost() {
		t.Fatal("expected IsDockerHost() == true")
	}
}

// TestIsDockerHost_FalseWhenNamespaced ensures container cgroups take precedence over a mounted socket.
func TestIsDockerHost_FalseWhenNamespaced(t *testing.T) {
	defer reset()

	statFile = func(path string) (os.FileInfo, error) {
		if path == "/var/run/docker.sock" {
			return nil, nil
		}
		return nil, os.ErrNotExist
	}

	readFile = func(path string) ([]byte, error) {
		return []byte("0::/docker/xyz"), nil
	}

	if IsDockerHost() {
		t.Fatal("expected IsDockerHost() == false in container namespace")
	}
}

// TestIsDockerHostFalseForNonDockerContainerMarkers ensures other container runtimes are not mistaken for Docker hosts.
func TestIsDockerHostFalseForNonDockerContainerMarkers(t *testing.T) {
	markers := []string{"container", "kubepods", "containerd", "podman", "libpod"}
	for _, marker := range markers {
		t.Run(marker, func(t *testing.T) {
			defer reset()
			statFile = func(path string) (os.FileInfo, error) {
				if path == fileDockerSock {
					return nil, nil
				}
				return nil, os.ErrNotExist
			}
			readFile = func(string) ([]byte, error) {
				return []byte("0::/" + marker + "/workload"), nil
			}

			if IsDockerHost() {
				t.Fatalf("expected %s cgroup marker to reject Docker-host detection", marker)
			}
		})
	}
}

// TestIsDockerHost_NoSocket ensures host classification requires Docker daemon evidence.
func TestIsDockerHost_NoSocket(t *testing.T) {
	defer reset()

	statFile = func(path string) (os.FileInfo, error) { return nil, os.ErrNotExist }
	if IsDockerHost() {
		t.Fatal("expected false when docker.sock missing")
	}
}

// TestIsDockerHost_ReadError ensures unreadable cgroup state fails closed.
func TestIsDockerHost_ReadError(t *testing.T) {
	defer reset()

	statFile = func(path string) (os.FileInfo, error) {
		if path == fileDockerSock {
			return nil, nil
		}
		return nil, os.ErrNotExist
	}
	readFile = func(path string) ([]byte, error) { return nil, os.ErrPermission }

	if IsDockerHost() {
		t.Fatal("expected false when cgroup read fails")
	}
}

// TestIsKubernetes_EnvVar ensures the service environment marker identifies Kubernetes pods.
func TestIsKubernetes_EnvVar(t *testing.T) {
	defer reset()

	mockEnv(map[string]string{
		"KUBERNETES_SERVICE_HOST": "10.0.0.1",
	})

	if !IsKubernetes() {
		t.Fatal("expected IsKubernetes() == true via env")
	}
}

// TestIsKubernetes_Cgroup ensures pod cgroup membership identifies Kubernetes without environment injection.
func TestIsKubernetes_Cgroup(t *testing.T) {
	defer reset()

	mockEnv(nil)

	readFile = func(path string) ([]byte, error) {
		return []byte("0::/kubepods/besteffort"), nil
	}

	if !IsKubernetes() {
		t.Fatal("expected IsKubernetes() == true via cgroup")
	}
}

// TestIsKubernetes_False ensures clean host evidence is not misclassified as Kubernetes.
func TestIsKubernetes_False(t *testing.T) {
	defer reset()

	mockEnv(nil)

	readFile = func(path string) ([]byte, error) {
		return []byte("0::/user.slice"), nil
	}

	if IsKubernetes() {
		t.Fatal("expected false when no kube signal")
	}
}

// TestIsContainer_Docker ensures Docker evidence contributes to the aggregate container result.
func TestIsContainer_Docker(t *testing.T) {
	defer reset()

	statFile = func(path string) (os.FileInfo, error) {
		if path == "/.dockerenv" {
			return nil, nil
		}
		return nil, os.ErrNotExist
	}

	if !IsContainer() {
		t.Fatal("expected true because IsDocker() is true")
	}
}

// TestIsContainerPodmanMarker ensures Podman's standard marker participates in aggregate detection.
func TestIsContainerPodmanMarker(t *testing.T) {
	defer reset()

	statFile = func(path string) (os.FileInfo, error) {
		if path == fileContainerEnv {
			return nil, nil
		}
		return nil, os.ErrNotExist
	}
	readFile = func(string) ([]byte, error) { return nil, os.ErrNotExist }
	mockEnv(nil)

	if !IsContainer() {
		t.Fatal("expected /run/.containerenv to identify a Podman container")
	}
}

// TestIsContainer_ContainerCgroup ensures generic container cgroups are recognized without vendor markers.
func TestIsContainer_ContainerCgroup(t *testing.T) {
	defer reset()

	statFile = func(path string) (os.FileInfo, error) { return nil, os.ErrNotExist }

	readFile = func(path string) ([]byte, error) {
		return []byte("0::/containerd/xyz"), nil
	}

	if !IsContainer() {
		t.Fatal("expected true from containerd cgroup")
	}
}

// TestIsContainer_Kubernetes ensures Kubernetes evidence contributes to the aggregate container result.
func TestIsContainer_Kubernetes(t *testing.T) {
	defer reset()

	mockEnv(map[string]string{
		"KUBERNETES_SERVICE_HOST": "10.0.0.1",
	})

	if !IsContainer() {
		t.Fatal("expected IsContainer() == true in kube")
	}
}

// TestIsContainer_False ensures an ordinary host remains outside the aggregate container classification.
func TestIsContainer_False(t *testing.T) {
	defer reset()

	statFile = func(path string) (os.FileInfo, error) {
		return nil, os.ErrNotExist
	}

	readFile = func(path string) ([]byte, error) {
		return []byte("0::/user.slice"), nil
	}

	mockEnv(nil)

	if IsContainer() {
		t.Fatal("expected false when no container signals")
	}
}

// TestIsContainer_ReadErrorButEnvPresent ensures strong environment evidence survives unavailable cgroup data.
func TestIsContainer_ReadErrorButEnvPresent(t *testing.T) {
	defer reset()

	statFile = func(path string) (os.FileInfo, error) { return nil, os.ErrNotExist }
	readFile = func(path string) ([]byte, error) { return nil, os.ErrPermission }
	mockEnv(map[string]string{"KUBERNETES_SERVICE_HOST": "10.0.0.1"})

	if !IsContainer() {
		t.Fatal("expected true when kube env present even if cgroup read fails")
	}
}

// TestIsContainer_ReadErrorNoEnv ensures missing evidence and unreadable cgroups fail closed.
func TestIsContainer_ReadErrorNoEnv(t *testing.T) {
	defer reset()

	statFile = func(path string) (os.FileInfo, error) { return nil, os.ErrNotExist }
	readFile = func(path string) ([]byte, error) { return nil, os.ErrPermission }
	mockEnv(nil)

	if IsContainer() {
		t.Fatal("expected false when cgroup read fails and no env hint")
	}
}

// TestIsContainer_EnvWinsWhenCgroupClean ensures explicit container metadata outranks neutral cgroup text.
func TestIsContainer_EnvWinsWhenCgroupClean(t *testing.T) {
	defer reset()

	statFile = func(path string) (os.FileInfo, error) { return nil, os.ErrNotExist }
	readFile = func(path string) ([]byte, error) { return []byte("0::/user.slice"), nil }
	mockEnv(map[string]string{"KUBERNETES_SERVICE_HOST": "10.0.0.2"})

	if !IsContainer() {
		t.Fatal("expected true when kube env set even with clean cgroup")
	}
}

// TestIsContainer_GenericContainerCgroup ensures runtime-neutral cgroup markers remain supported.
func TestIsContainer_GenericContainerCgroup(t *testing.T) {
	defer reset()

	statFile = func(path string) (os.FileInfo, error) { return nil, os.ErrNotExist }
	readFile = func(path string) ([]byte, error) { return []byte("0::/container/abc"), nil }
	mockEnv(nil)

	if !IsContainer() {
		t.Fatal("expected true for generic container cgroup marker")
	}
}
