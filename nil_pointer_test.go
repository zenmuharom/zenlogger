package zenlogger

import (
	"strings"
	"testing"
)

func Test_nil_pointer(t *testing.T) {
	logger := NewZenlogger()
	logger.SetConfig(Config{BeautifyJson: true})

	// Test with nil pointer
	var nilPointer *string = nil
	result := logger.Info("testing nil pointer", ZenField{Key: "string", Value: nilPointer})
	
	if !strings.Contains(result, "null") {
		t.Errorf("Expected null value for nil pointer, got: %s", result)
	}

	// Test with valid pointer
	validString := "test value"
	validPointer := &validString
	result2 := logger.Info("testing valid pointer", ZenField{Key: "string", Value: validPointer})
	
	if !strings.Contains(result2, "test value") {
		t.Errorf("Expected 'test value' for valid pointer, got: %s", result2)
	}

	// Test with nil pointer to struct
	var nilPerson *Person = nil
	result3 := logger.Info("testing nil struct pointer", ZenField{Key: "person", Value: nilPerson})
	
	if !strings.Contains(result3, "null") {
		t.Errorf("Expected null value for nil struct pointer, got: %s", result3)
	}

	t.Log("All nil pointer tests passed")
}
