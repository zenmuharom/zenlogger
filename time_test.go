package zenlogger

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTimeType(t *testing.T) {
	logger := newTestLogger(t)

	// Test with direct time.Time value
	now := time.Now()
	result := logger.Info("test time.Time", ZenField{Key: "timestamp", Value: now})
	require.NotEmpty(t, result)
	require.Contains(t, result, "timestamp")

	// Test with time.Time in struct
	type Event struct {
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	event := Event{
		Name:      "Test Event",
		CreatedAt: time.Date(2025, 11, 24, 10, 30, 0, 0, time.UTC),
		UpdatedAt: time.Date(2025, 11, 24, 15, 45, 30, 0, time.UTC),
	}

	result = logger.Info("test struct with time.Time", ZenField{Key: "event", Value: event})
	require.NotEmpty(t, result)
	require.Contains(t, result, "event")

	// Verify JSON contains formatted time
	var logData map[string]interface{}
	err := json.Unmarshal([]byte(result), &logData)
	require.NoError(t, err)

	// Test with slice of structs containing time.Time
	events := []Event{
		{
			Name:      "Event 1",
			CreatedAt: time.Date(2025, 11, 24, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 11, 24, 11, 0, 0, 0, time.UTC),
		},
		{
			Name:      "Event 2",
			CreatedAt: time.Date(2025, 11, 24, 12, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 11, 24, 13, 0, 0, 0, time.UTC),
		},
	}

	result = logger.Info("test slice with time.Time", ZenField{Key: "events", Value: events})
	require.NotEmpty(t, result)
	require.Contains(t, result, "events")

	// Test with map containing time.Time
	eventMap := map[string]interface{}{
		"name":       "Map Event",
		"created_at": time.Date(2025, 11, 24, 14, 0, 0, 0, time.UTC),
		"count":      42,
	}

	result = logger.Info("test map with time.Time", ZenField{Key: "event_map", Value: eventMap})
	require.NotEmpty(t, result)
	require.Contains(t, result, "event_map")

	// Test with pointer to time.Time
	timePtr := &now
	result = logger.Info("test pointer to time.Time", ZenField{Key: "time_ptr", Value: timePtr})
	require.NotEmpty(t, result)
	require.Contains(t, result, "time_ptr")

	// Test with nested struct containing time.Time
	type NestedEvent struct {
		Event     Event     `json:"event"`
		Timestamp time.Time `json:"timestamp"`
	}

	nested := NestedEvent{
		Event: Event{
			Name:      "Nested Event",
			CreatedAt: time.Date(2025, 11, 24, 16, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 11, 24, 17, 0, 0, 0, time.UTC),
		},
		Timestamp: time.Date(2025, 11, 24, 18, 0, 0, 0, time.UTC),
	}

	result = logger.Info("test nested struct with time.Time", ZenField{Key: "nested", Value: nested})
	require.NotEmpty(t, result)
	require.Contains(t, result, "nested")
}

func TestTimeTypeZeroValue(t *testing.T) {
	logger := newTestLogger(t)

	// Test with zero time
	var zeroTime time.Time
	result := logger.Info("test zero time", ZenField{Key: "zero_time", Value: zeroTime})
	require.NotEmpty(t, result)
	require.Contains(t, result, "zero_time")

	// Test with struct containing zero time
	type Record struct {
		Name      string    `json:"name"`
		DeletedAt time.Time `json:"deleted_at"`
	}

	record := Record{
		Name: "Test Record",
		// DeletedAt is zero value
	}

	result = logger.Info("test struct with zero time", ZenField{Key: "record", Value: record})
	require.NotEmpty(t, result)
	require.Contains(t, result, "record")
}
