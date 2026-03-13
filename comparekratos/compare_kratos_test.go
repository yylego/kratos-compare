package comparekratos

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestComparePath_SameFile(t *testing.T) {
	tmpRoot := t.TempDir()

	fileA := filepath.Join(tmpRoot, "a.go")
	fileB := filepath.Join(tmpRoot, "b.go")
	require.NoError(t, os.WriteFile(fileA, []byte("package main\n"), 0644))
	require.NoError(t, os.WriteFile(fileB, []byte("package main\n"), 0644))

	output := ComparePath(fileA, fileB)
	require.Empty(t, output)
}

func TestComparePath_DiffFile(t *testing.T) {
	tmpRoot := t.TempDir()

	fileA := filepath.Join(tmpRoot, "a.go")
	fileB := filepath.Join(tmpRoot, "b.go")
	require.NoError(t, os.WriteFile(fileA, []byte("package main\n"), 0644))
	require.NoError(t, os.WriteFile(fileB, []byte("package main\n\nfunc init() {}\n"), 0644))

	output := ComparePath(fileA, fileB)
	require.NotEmpty(t, output)
	require.Contains(t, string(output), "func init")
}

func TestComparePath_Same(t *testing.T) {
	tmpRoot := t.TempDir()

	pathA := filepath.Join(tmpRoot, "a")
	pathB := filepath.Join(tmpRoot, "b")
	require.NoError(t, os.MkdirAll(pathA, 0755))
	require.NoError(t, os.MkdirAll(pathB, 0755))

	require.NoError(t, os.WriteFile(filepath.Join(pathA, "main.go"), []byte("package main\n"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(pathB, "main.go"), []byte("package main\n"), 0644))

	output := ComparePath(pathA, pathB)
	require.Empty(t, output)
}

func TestComparePath_Diff(t *testing.T) {
	tmpRoot := t.TempDir()

	pathA := filepath.Join(tmpRoot, "a")
	pathB := filepath.Join(tmpRoot, "b")
	require.NoError(t, os.MkdirAll(pathA, 0755))
	require.NoError(t, os.MkdirAll(pathB, 0755))

	require.NoError(t, os.WriteFile(filepath.Join(pathA, "main.go"), []byte("package main\n"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(pathB, "main.go"), []byte("package main\n\nfunc init() {}\n"), 0644))

	output := ComparePath(pathA, pathB)
	require.NotEmpty(t, output)
	require.Contains(t, string(output), "func init")
}

func TestShowReadableChanges_Same(t *testing.T) {
	tmpRoot := t.TempDir()

	pathA := filepath.Join(tmpRoot, "a")
	pathB := filepath.Join(tmpRoot, "b")
	require.NoError(t, os.MkdirAll(pathA, 0755))
	require.NoError(t, os.MkdirAll(pathB, 0755))

	require.NoError(t, os.WriteFile(filepath.Join(pathA, "main.go"), []byte("package main\n"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(pathB, "main.go"), []byte("package main\n"), 0644))

	output := ShowReadableChanges(pathA, pathB)
	require.Empty(t, output)
}

func TestShowReadableChanges_Diff(t *testing.T) {
	tmpRoot := t.TempDir()

	pathA := filepath.Join(tmpRoot, "a")
	pathB := filepath.Join(tmpRoot, "b")
	require.NoError(t, os.MkdirAll(pathA, 0755))
	require.NoError(t, os.MkdirAll(pathB, 0755))

	require.NoError(t, os.WriteFile(filepath.Join(pathA, "main.go"), []byte("package main\n"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(pathB, "main.go"), []byte("package main\n\nfunc init() {}\n"), 0644))

	output := ShowReadableChanges(pathA, pathB)
	require.NotEmpty(t, output)
	require.Contains(t, string(output), "func init")
}

func TestGenerateChangesFile_Same(t *testing.T) {
	tmpRoot := t.TempDir()

	pathA := filepath.Join(tmpRoot, "a")
	pathB := filepath.Join(tmpRoot, "b")
	require.NoError(t, os.MkdirAll(pathA, 0755))
	require.NoError(t, os.MkdirAll(pathB, 0755))

	require.NoError(t, os.WriteFile(filepath.Join(pathA, "main.go"), []byte("package main\n"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(pathB, "main.go"), []byte("package main\n"), 0644))

	outputPath := filepath.Join(tmpRoot, "changes.md")
	GenerateChangesFile(pathA, pathB, outputPath)

	content, e := os.ReadFile(outputPath)
	require.NoError(t, e)
	require.Contains(t, string(content), "NO CHANGES")
}

func TestGenerateChangesFile_Diff(t *testing.T) {
	tmpRoot := t.TempDir()

	pathA := filepath.Join(tmpRoot, "a")
	pathB := filepath.Join(tmpRoot, "b")
	require.NoError(t, os.MkdirAll(pathA, 0755))
	require.NoError(t, os.MkdirAll(pathB, 0755))

	require.NoError(t, os.WriteFile(filepath.Join(pathA, "main.go"), []byte("package main\n"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(pathB, "main.go"), []byte("package main\n\nfunc init() {}\n"), 0644))

	outputPath := filepath.Join(tmpRoot, "changes.md")
	GenerateChangesFile(pathA, pathB, outputPath)

	content, e := os.ReadFile(outputPath)
	require.NoError(t, e)
	require.Contains(t, string(content), "Changes")
	require.Contains(t, string(content), "```diff")
}
