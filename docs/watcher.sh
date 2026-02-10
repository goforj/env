#!/usr/bin/env bash

if ! command -v wgo >/dev/null 2>&1; then
  go install github.com/bokwoon95/wgo@latest
fi

echo "Watching for .go file changes to regenerate documentation..."

wgo -verbose -file=.go -xdir examples \
  go run ./docs/examplegen/main.go :: \
  go run ./docs/readme/main.go
