# Zenlogger
[![GoDoc][doc-img]][doc]

This library is built based on internal & personal needs.
So, this lib can be not suit for your need.



**This library support such as**
- Logging to file (including customize naming file)
- Structured log to satisfy Gcloud logs structure 
- Customize structure & keys

# Sample Output

###Sample structured output log

This is the sample of output output when beautifyJson set to `true`.
```
{
  "pid": "97732c84970d49a883193352f7da24f3",
  "severity": "Info",
  "timestamp": "2022-11-29T20:01:48+0700",
  "caller": "map_test.go:47",
  "message": {
    "title": "test map",
    "values": {
      "halo": {
        "2": "two",
        "developer": {
          "Age": 27,
          "Gender": "male",
          "IsMarried": false,
          "Name": "Zeni",
          "Relationship": "",
          "Supporter": [
            {
              "Age": 22,
              "Gender": "bunga",
              "IsMarried": true,
              "Name": "Mawar",
              "Relationship": "",
              "Supporter": []
            },
            {
              "Age": 20,
              "Gender": "warna",
              "IsMarried": false,
              "Name": "Hitam",
              "Relationship": "",
              "Supporter": []
            }
          ]
        },
        "four": 4,
        "one": 1,
        "three": "tilu"
      }
    }
  }
}
```

<br />

# Installation
**How to use**:
This is the sample code of zenlogger usage
```
package main

import (
	"github.com/zenmuharom/zenlogger"
)

type Person struct {
	Name         string
	Gender       string
	Age          int
	IsMarried    bool
	Relationship string
	Supporter    []Person
}

func Test_map(t *testing.T) {

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

```

**Config**
This is the sample of config you can customize
```
	config := zenlogger.Config{
		Pid: zenlogger.ZenConf{
			Label: "insertId",
		},
		Severity: zenlogger.Severity{
			Label:  "Level",
			Access: "API",
			Info:   "THIS IS INFO",
			Debug:  "DEBUG",
			Warning: "Please Attention To This"
			Error:  "OH MY GOD THIS IS ERROR",
			Query:  "QUERY",
		},
		Caller: zenlogger.Caller{
			Label: "trace_file",
			Level: 0,
		},
		Message: zenlogger.Message{
			Label: "pesan",
			Title: zenlogger.ZenConf{
				Label: "judul",
			},
			Values: zenlogger.ZenConf{
				Label: "isi",
			},
		},
		BeautifyJson: true,	
	}
	logger.SetConfig(config)
```

You can set the log to make zenlogger write into file by adding:
```
	config := zenlogger.Config{}
	config.Output.Path = "logs"
	config.Output.Format = "20060102"
	config.BeautifyJson = true
	logger.SetConfig(config)
```
Zenlogger will automatically make directory logs (if not exists), and write into file with golang timed format.
<br><br><br>



You can set or try it in the link below here: 
[GO Playground](https://go.dev/play/p/1wvZJefYor1 "GO Playground").
<br><br>

This library is released under: [MIT License](https://github.com/zenmuharom/zenlogger/blob/master/LICENSE.txt "MIT License").

[doc-img]: https://pkg.go.dev/badge/github.com/zenmuharom/zenlogger
[doc]: https://pkg.go.dev/github.com/zenmuharom/zenlogger