package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// getProjectRoot returns the project root directory.
func getProjectRoot(t *testing.T) string {
	// Get current file's directory
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not get caller information")
	}
	// Navigate up to project root: main_test.go -> pocket-vm-lab -> cmd -> project root
	return filepath.Dir(filepath.Dir(filepath.Dir(thisFile)))
}

func buildBinary(t *testing.T, projectRoot string) string {
	binaryName := "pocket-vm-lab-test"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	binaryPath := filepath.Join(t.TempDir(), binaryName)

	cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/pocket-vm-lab")
	cmd.Dir = projectRoot
	cmd.Env = append(os.Environ(), "PATH="+os.Getenv("PATH"))

	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, output)
	}

	return binaryPath
}

func TestDemoFlag(t *testing.T) {
	projectRoot := getProjectRoot(t)
	binaryPath := buildBinary(t, projectRoot)

	cmd := exec.Command(binaryPath, "-demo")
	cmd.Dir = projectRoot

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("demo failed: %v\nstderr: %s", err, stderr.String())
	}

	expected := "8\n"
	if stdout.String() != expected {
		t.Fatalf("expected %q, got %q", expected, stdout.String())
	}
}

func TestFileFlag(t *testing.T) {
	projectRoot := getProjectRoot(t)
	binaryPath := buildBinary(t, projectRoot)

	cmd := exec.Command(binaryPath, "-file", "testdata/sample.bin")
	cmd.Dir = projectRoot

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("file execution failed: %v\nstderr: %s", err, stderr.String())
	}

	expected := "8\n"
	if stdout.String() != expected {
		t.Fatalf("expected %q, got %q", expected, stdout.String())
	}
}

func TestNoFlags(t *testing.T) {
	projectRoot := getProjectRoot(t)
	binaryPath := buildBinary(t, projectRoot)

	cmd := exec.Command(binaryPath)
	cmd.Dir = projectRoot

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		t.Fatal("expected error when no flags provided")
	}

	// Should exit with non-zero status
	if cmd.ProcessState.Success() {
		t.Fatal("expected non-zero exit code")
	}
}
