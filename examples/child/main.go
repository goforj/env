//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env/v2"
	"os"
)

func main() {
	// Child returns a new scope rooted at the current prefix plus name.

	// Example: named child scope
	_ = os.Setenv("STORAGE_PUBLIC_ROOT", "storage/app/public")

	public := env.WithPrefix("STORAGE").Child("PUBLIC")
	env.Dump(
		public.Key("ROOT"),
		public.Get("ROOT", "storage/app/public"),
	)
	// #string "STORAGE_PUBLIC_ROOT"
	// #string "storage/app/public"
}
