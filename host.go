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
//	godump.Dump(env.IsHostEnvironment())
//
//	// #bool true  (on bare-metal/VM hosts)
//	// #bool false (inside containers)
func IsHostEnvironment() bool {
	return !IsContainer() &&
		!IsDockerInDocker()
}
