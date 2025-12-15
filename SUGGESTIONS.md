## README suggestions

- Add a short “Why dependency on godump?” note in the examples section to clarify it is test/example-only and not required for library consumers.
- Include a “Quickstart” snippet that shows loading env files (`LoadEnvFileIfExists`) plus a couple of typed getters, mirroring the new doc comments.
- Document the `Dump` helper as the supported way to view example outputs instead of calling `godump.Dump` directly.
- Surface testing guidance: `go test ./... -cover` already hits 100%—calling that out reassures users about stability.
- Mention example generation/watch flow (docs/readme and examplegen) so contributors know how README/examples stay in sync.
- Add a small “Container detection matrix” table summarizing IsDocker/IsDockerHost/IsDockerInDocker/IsContainer/IsKubernetes behavior for quick reference.
