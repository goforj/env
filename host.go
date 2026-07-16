package env

// IsHostEnvironment reports whether the process is running *outside* any
// container or orchestrated runtime.
// @group Container detection
// @behavior readonly
//
// Being a Docker host does NOT count as being in a container.
//
// Example:
//
//	env.Dump(env.IsHostEnvironment())
//	// #bool true  (on bare-metal/VM hosts)
//	// #bool false (inside containers)
func IsHostEnvironment() bool {
	return isHostEnvironmentWithEnv(getEnv)
}

// isHostEnvironmentWithEnv evaluates host detection against a caller-supplied environment view.
func isHostEnvironmentWithEnv(getenv func(string) string) bool {
	return !isContainerWithEnv(getenv) &&
		!IsDockerInDocker()
}
