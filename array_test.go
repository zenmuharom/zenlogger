package zenlogger

import (
	"testing"
)

func Test_array(t *testing.T) {
	logger := NewZenlogger()
	logger.SetConfig(Config{BeautifyJson: true})

	family := make([]Person, 0)
	family = append(family, Person{
		Name:         "Gin - Gin",
		Gender:       "male",
		Age:          23,
		Relationship: "Brother",
		IsMarried:    false,
	})
	family = append(family, Person{
		Name:         "Dad",
		Gender:       "male",
		Age:          60,
		Relationship: "Father",
		IsMarried:    true,
	})
	family = append(family, Person{
		Name:         "Mom",
		Gender:       "female",
		Age:          48,
		Relationship: "Mother",
		IsMarried:    true,
	})

	logger.Warning("family member", ZenField{Key: "my family", Value: family})
}
