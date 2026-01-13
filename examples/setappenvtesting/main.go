//go:build ignore
// +build ignore

package main

import "github.com/goforj/env"

func main() {
	// SetAppEnvTesting sets APP_ENV to "testing".

	_ = env.SetAppEnvTesting()
	env.Dump(env.GetAppEnv())

	// #string "testing"
}
