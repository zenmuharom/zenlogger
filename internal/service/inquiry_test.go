package service_test

import (
	"strings"
	"testing"

	zenlogger "github.com/zenmuharom/zenlogger"
)

func TestCallerPath_FromNestedPackage(t *testing.T) {
	logger := zenlogger.NewZenlogger()
	logStr := logger.Info("nested package caller")

	if !strings.Contains(logStr, "\"caller\":\"service/inquiry_test.go:") {
		t.Fatalf("expected nested caller path, got: %s", logStr)
	}
}
