package zenlogger

import (
	"os"
	"path/filepath"
	"testing"
)

func newTestLogger(t *testing.T) Zenlogger {
	t.Helper()

	logger := NewZenlogger()
	setTestConfig(t, logger, Config{})
	return logger
}

func setTestConfig(t *testing.T, logger Zenlogger, cfg Config) {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to resolve cwd: %v", err)
	}

	if cfg.Output.Path == "" {
		cfg.Output.Path = filepath.Join(wd, "logs")
	}
	if cfg.Output.Format == "" {
		cfg.Output.Format = "20060102"
	}

	logger.SetConfig(cfg)
}
