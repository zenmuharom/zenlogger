package zenlogger

import (
	"testing"
)

func Test_pointer(t *testing.T) {
	logger := newTestLogger(t)
	setTestConfig(t, logger, Config{BeautifyJson: true})

	patrick := &Person{
		Name:         "Patrick",
		Gender:       "male",
		Age:          28,
		IsMarried:    false,
		Relationship: "single",
		Supporter:    []Person{},
	}

	logger.Info("this is american", ZenField{Key: "people", Value: patrick})

}
