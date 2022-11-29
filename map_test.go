package zenlogger

import "testing"

func Test_map(t *testing.T) {

	type Person struct {
		Name      string
		Gender    string
		Age       int
		IsMarried bool
		Supporter []Person
	}

	logger := NewZenlogger()
	logger.SetConfig(Config{BeautifyJson: true})

	person := Person{
		Name:      "Zeni",
		Gender:    "male",
		Age:       27,
		IsMarried: false,
		Supporter: []Person{
			{
				Name:      "Mawar",
				Age:       22,
				Gender:    "bunga",
				IsMarried: true,
			},
			{
				Name:      "Hitam",
				Age:       20,
				Gender:    "warna",
				IsMarried: false,
			},
		},
	}

	testMap := map[string]interface{}{
		"one":       "1",
		"2":         "two",
		"three":     "tilu",
		"four":      4,
		"developer": person,
	}
	logger.Info("test map", ZenField{Key: "halo", Value: testMap})
}
