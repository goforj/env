//go:build ignore
// +build ignore

package main

import "github.com/goforj/env"

func main() {
	// SetAppEnvDev sets APP_ENV to "dev".

	_ = env.SetAppEnvDev()
	env.Dump(env.GetAppEnv())

	// #string "dev"
}
