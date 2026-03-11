//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env/v2"
	"os"
)

func main() {
	// WithPrefix returns a scope rooted at prefix after minimal normalization.

	// Example: root scope access
	_ = os.Setenv("STORAGE_DRIVER", "local")
	_ = os.Setenv("STORAGE_ROOT", "storage/app/private")

	storage := env.WithPrefix(" STORAGE ")
	env.Dump(
		storage.Key("DRIVER"),
		storage.Get("DRIVER", "s3"),
		storage.Get("ROOT", "storage/app/private"),
	)
	// #string "STORAGE_DRIVER"
	// #string "local"
	// #string "storage/app/private"
}
