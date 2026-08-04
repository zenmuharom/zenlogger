package zenlogger

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type logEnvelope struct {
	Message struct {
		Values map[string]interface{} `json:"values"`
	} `json:"message"`
}

func TestMasking_FullMasked(t *testing.T) {
	logger := newTestLogger(t)

	logStr := logger.Info("db connection", ZenField{
		Key:   "password",
		Value: "mysecret",
		Type:  FULL_MASKED,
	})

	var envelope logEnvelope
	err := json.Unmarshal([]byte(logStr), &envelope)
	require.NoError(t, err)
	require.Equal(t, "********", envelope.Message.Values["password"])
}

func TestMasking_FirstMasked(t *testing.T) {
	logger := newTestLogger(t)

	logStr := logger.Info("db connection", ZenField{
		Key:       "username",
		Value:     "finpayuser",
		Type:      FIRST_MASKED,
		MaskCount: 3,
	})

	var envelope logEnvelope
	err := json.Unmarshal([]byte(logStr), &envelope)
	require.NoError(t, err)
	require.Equal(t, "***payuser", envelope.Message.Values["username"])
}

func TestMasking_LastMasked(t *testing.T) {
	logger := newTestLogger(t)

	logStr := logger.Info("db connection", ZenField{
		Key:       "phone",
		Value:     "08123456789",
		Type:      LAST_MASKED,
		MaskCount: 4,
	})

	var envelope logEnvelope
	err := json.Unmarshal([]byte(logStr), &envelope)
	require.NoError(t, err)
	require.Equal(t, "0812345****", envelope.Message.Values["phone"])
}

func TestMasking_FirstLastMasked(t *testing.T) {
	logger := newTestLogger(t)

	logStr := logger.Info("db connection", ZenField{
		Key:       "token",
		Value:     "abcde12345xyz",
		Type:      FIRST_LAST_MASKED,
		MaskCount: 2,
	})

	var envelope logEnvelope
	err := json.Unmarshal([]byte(logStr), &envelope)
	require.NoError(t, err)
	require.Equal(t, "**cde12345x**", envelope.Message.Values["token"])
}

func TestMasking_DefaultCount(t *testing.T) {
	logger := newTestLogger(t)

	logStr := logger.Info("db connection", ZenField{
		Key:   "username",
		Value: "finpayuser",
		Type:  FIRST_MASKED,
	})

	var envelope logEnvelope
	err := json.Unmarshal([]byte(logStr), &envelope)
	require.NoError(t, err)
	require.Equal(t, "*inpayuser", envelope.Message.Values["username"])
}

func TestMasking_AutoByKey(t *testing.T) {
	logger := newTestLogger(t)
	setTestConfig(t, logger, Config{
		Sensitive: SensitiveFieldConfig{
			Enabled:         true,
			CaseInsensitive: true,
			Rules: map[string]SensitiveFieldRule{
				"password": {
					Type: FULL_MASKED,
				},
				"token": {
					Type:      FIRST_LAST_MASKED,
					MaskCount: 2,
				},
			},
		},
	})

	logStr := logger.Info("db connection",
		ZenField{Key: "username", Value: "myuser"},
		ZenField{Key: "password", Value: "mysecret"},
		ZenField{Key: "TOKEN", Value: "abcde12345xyz"},
	)

	var envelope logEnvelope
	err := json.Unmarshal([]byte(logStr), &envelope)
	require.NoError(t, err)
	require.Equal(t, "myuser", envelope.Message.Values["username"])
	require.Equal(t, "********", envelope.Message.Values["password"])
	require.Equal(t, "**cde12345x**", envelope.Message.Values["TOKEN"])
}

func TestMasking_Redacted(t *testing.T) {
	logger := newTestLogger(t)

	logStr := logger.Info("db connection", ZenField{
		Key:   "apiKey",
		Value: "abcd-efgh",
		Type:  REDACTED,
	})

	var envelope logEnvelope
	err := json.Unmarshal([]byte(logStr), &envelope)
	require.NoError(t, err)
	require.Equal(t, "[REDACTED]", envelope.Message.Values["apiKey"])
}

func TestMasking_HashSHA256(t *testing.T) {
	logger := newTestLogger(t)

	logStr := logger.Info("db connection", ZenField{
		Key:   "token",
		Value: "abc123",
		Type:  HASH_SHA256,
	})

	var envelope logEnvelope
	err := json.Unmarshal([]byte(logStr), &envelope)
	require.NoError(t, err)
	require.Equal(t, "6ca13d52ca70c883e0f0bb101e425a89e8624de51db2d2392593af6a84118090", envelope.Message.Values["token"])
}
