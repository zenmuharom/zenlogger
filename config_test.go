package zenlogger

import "testing"

func Test_config(t *testing.T) {
	logger := NewZenlogger()
	config := Config{
		Severity: Severity{
			Label:   "level",
			Info:    "This is for information",
			Warning: "Please Attention To this log",
			Error:   "oh no error",
		},
		BeautifyJson: true,
	}
	logger.SetConfig(config)
	logger.Info("This is info")
	logger.Debug("this is debug")
	logger.Warning("This is warning")
	logger.Error("this is error")

}
