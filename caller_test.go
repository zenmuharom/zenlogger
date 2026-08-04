package zenlogger

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_formatCallerPath_RelativePath(t *testing.T) {
	wd, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("failed to resolve cwd: %v", err)
	}

	input := filepath.Join(wd, "internal", "service", "inquiry.go")
	got := formatCallerPath(input)
	want := "internal/service/inquiry.go"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func Test_formatCallerPath_RelativeInput(t *testing.T) {
	input := filepath.Join("internal", "service", "inquiry.go")
	got := formatCallerPath(input)
	want := "internal/service/inquiry.go"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func Test_formatCallerPath_FromNestedCwdUsesModuleRoot(t *testing.T) {
	tmpDir := t.TempDir()
	moduleRoot := filepath.Join(tmpDir, "example")
	nestedDir := filepath.Join(moduleRoot, "internal", "web")
	serviceDir := filepath.Join(moduleRoot, "internal", "service")
	callerFile := filepath.Join(moduleRoot, "internal", "service", "inquiry.go")

	if err := os.MkdirAll(nestedDir, 0o755); err != nil {
		t.Fatalf("failed to create nested dir: %v", err)
	}
	if err := os.MkdirAll(serviceDir, 0o755); err != nil {
		t.Fatalf("failed to create service dir: %v", err)
	}

	goModPath := filepath.Join(moduleRoot, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module example\n\ngo 1.19\n"), 0o644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}
	if err := os.WriteFile(callerFile, []byte("package service\n"), 0o644); err != nil {
		t.Fatalf("failed to write caller file: %v", err)
	}

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer func() {
		_ = os.Chdir(oldWd)
	}()

	if err := os.Chdir(nestedDir); err != nil {
		t.Fatalf("failed to change cwd: %v", err)
	}

	got := formatCallerPath(callerFile)
	want := "internal/service/inquiry.go"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func Test_formatCallerPath_BasenameInputResolvesToModulePath(t *testing.T) {
	tmpDir := t.TempDir()
	moduleRoot := filepath.Join(tmpDir, "example")
	nestedDir := filepath.Join(moduleRoot, "internal", "web")
	serviceDir := filepath.Join(moduleRoot, "internal", "service")
	callerFile := filepath.Join(moduleRoot, "internal", "service", "inquiry.go")

	if err := os.MkdirAll(nestedDir, 0o755); err != nil {
		t.Fatalf("failed to create nested dir: %v", err)
	}
	if err := os.MkdirAll(serviceDir, 0o755); err != nil {
		t.Fatalf("failed to create service dir: %v", err)
	}

	goModPath := filepath.Join(moduleRoot, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module example\n\ngo 1.19\n"), 0o644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}
	if err := os.WriteFile(callerFile, []byte("package service\n"), 0o644); err != nil {
		t.Fatalf("failed to write caller file: %v", err)
	}

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer func() {
		_ = os.Chdir(oldWd)
	}()

	if err := os.Chdir(nestedDir); err != nil {
		t.Fatalf("failed to change cwd: %v", err)
	}

	got := formatCallerPath("inquiry.go")
	want := "internal/service/inquiry.go"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
