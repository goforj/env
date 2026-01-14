//go:build ignore
// +build ignore

package main

import "github.com/goforj/env/v2"

func main() {
	// SetAppEnvLocal sets APP_ENV to "local".

	_ = env.SetAppEnvLocal()
	env.Dump(env.GetAppEnv())

	// #string "local"
}
