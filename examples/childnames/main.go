//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env/v2"
	"os"
)

func main() {
	// ChildNames discovers named child scopes under the current prefix.

	// Example: discover child names
	_ = os.Setenv("STORAGE_DRIVER", "local")
	_ = os.Setenv("STORAGE_ROOT", "storage/app/private")
	_ = os.Setenv("STORAGE_PUBLIC_ROOT", "storage/app/public")
	_ = os.Setenv("STORAGE_AVATARS_BUCKET", "my-bucket")
	_ = os.Setenv("STORAGE_AVATARS_REGION", "us-east-1")

	names := env.WithPrefix("STORAGE").ChildNames([]string{
		"DRIVER",
		"ROOT",
		"BUCKET",
		"REGION",
	})
	env.Dump(names)
	// #[]string [
	//  0 => "AVATARS" #string
	//  1 => "PUBLIC" #string
	// ]
}
