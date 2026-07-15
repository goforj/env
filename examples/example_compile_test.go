package examples

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestExamplesBuild ensures every generated standalone program remains valid outside the workspace.
func TestExamplesBuild(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("cannot read examples directory: %v", err)
	}

	buildSlots := make(chan struct{}, 4)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			buildSlots <- struct{}{}
			defer func() { <-buildSlots }()

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if err := buildExampleWithoutTags(ctx, name); err != nil {
				t.Fatalf("example %q failed to build:\n%s", name, err)
			}
		})
	}
}

// buildExampleWithoutTags overlays generated source so ignored examples are still compiled in CI.
func buildExampleWithoutTags(ctx context.Context, exampleName string) error {
	orig := filepath.Join(exampleName, "main.go")
	originalPath, err := filepath.Abs(orig)
	if err != nil {
		return fmt.Errorf("resolve example path: %w", err)
	}

	src, err := os.ReadFile(originalPath)
	if err != nil {
		return fmt.Errorf("read main.go: %w", err)
	}

	clean := stripBuildTags(src)

	tmpDir, err := os.MkdirTemp("", "example-overlay-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	tmpFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(tmpFile, clean, 0o644); err != nil {
		return err
	}

	temporaryPath, err := filepath.Abs(tmpFile)
	if err != nil {
		return fmt.Errorf("resolve overlay path: %w", err)
	}
	overlay := map[string]any{
		"Replace": map[string]string{
			originalPath: temporaryPath,
		},
	}

	overlayJSON, err := json.Marshal(overlay)
	if err != nil {
		return err
	}

	overlayPath := filepath.Join(tmpDir, "overlay.json")
	if err := os.WriteFile(overlayPath, overlayJSON, 0o644); err != nil {
		return err
	}

	cmd := exec.CommandContext(
		ctx,
		"go", "build",
		"-overlay", overlayPath,
		"-o", os.DevNull,
		"./"+exampleName,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("go build exceeded deadline: %w", ctx.Err())
		}
		return fmt.Errorf("go build: %w\n%s", err, stderr.String())
	}

	return nil
}

// stripBuildTags removes only the leading constraints that intentionally hide generated programs.
func stripBuildTags(src []byte) []byte {
	lines := strings.Split(string(src), "\n")

	i := 0
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])

		if strings.HasPrefix(line, "//go:build") ||
			strings.HasPrefix(line, "// +build") ||
			line == "" {
			i++
			continue
		}

		break
	}

	return []byte(strings.Join(lines[i:], "\n"))
}
