package zenlogger

import "testing"

func Test_config(t *testing.T) {
	logger := newTestLogger(t)
	config := Config{
		Severity: Severity{
			Label:   "level",
			Info:    "This is for information",
			Warning: "Please Attention To this log",
			Error:   "oh no error",
		},
		Caller: Caller{
			Level: 0,
		},
		BeautifyJson: true,
	}
	setTestConfig(t, logger, config)
	logger.Info("This is info")
	logger.Debug("this is debug")
	logger.Warning("This is warning")
	logger.Error("this is error")

}
