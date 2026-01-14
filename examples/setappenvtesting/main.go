//go:build ignore
// +build ignore

package main

import "github.com/goforj/env/v2"

func main() {
	// SetAppEnvTesting sets APP_ENV to "testing".

	_ = env.SetAppEnvTesting()
	env.Dump(env.GetAppEnv())

	// #string "testing"
}
