# Zenlogger :robot:

[![GoDoc][doc-img]][doc]

This library is built based on internal & personal needs.
So, this lib can be not suit for your need.

**This library support such as**

- Logging to file (including customize naming file)
- Structured log to satisfy Gcloud logs structure
- Customize structure & keys
- Full GCP severity levels (Debug, Info, Notice, Warning, Error, Critical, Alert, Emergency)
- Log level filtering — execute only logs at or above configured level
- Caller metadata includes relative file path and line number (example: `internal/service/inquiry.go:32`)

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

### Installation :rocket:

### Short Usage

```
package main

import "github.com/zenmuharom/zenlogger"

func main() {
	logger := zenlogger.NewZenlogger()
	logger.SetConfig(zenlogger.Config{
		Output: zenlogger.Output{Path: "logs", Format: "20060102"},
		Sensitive: zenlogger.SensitiveFieldConfig{
			Enabled:         true,
			CaseInsensitive: true,
			Rules: map[string]zenlogger.SensitiveFieldRule{
				"password": {Type: zenlogger.FULL_MASKED},
				"token":    {Type: zenlogger.FIRST_LAST_MASKED, MaskCount: 2},
			},
		},
	})

	logger.Info("db connection",
		zenlogger.Field("username", "finpayuser"),
		zenlogger.Field("password", "mysecret"),
		zenlogger.Field("token", "abcde12345xyz"),
	)
}
```

Caller output format is `package/file.go:line` (example: `service/inquiry.go:39`).

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

<br/>

### Config

This is the sample of config you can customize

```
	logger := zenlogger.NewZenlogger()
	config := zenlogger.Config{
		Pid: zenlogger.ZenConf{
			Label: "insertId",
		},
		Severity: zenlogger.Severity{
			Label:     "Level",
			Access:    "API",
			Info:      "THIS IS INFO",
			Debug:     "DEBUG",
			Notice:    "NOTICE",
			Warning:   "Please Attention To This",
			Error:     "Fault",
			Query:     "DB",
			Critical:  "CRITICAL",
			Alert:     "ALERT",
			Emergency: "EMERGENCY",
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

### Set log to file

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

### Config Property :robot:

| Property     | sub             | description                                                                                                          |
| ------------ | --------------- | -------------------------------------------------------------------------------------------------------------------- |
| Pid          |                 | process id is label of the key that differentiate with another process in your program                               |
| Severity     | Label           | The label of severity key                                                                                            |
| Severity     | Access          | The value of severity when you call zenlogger.Access()                                                               |
| Severity     | Info            | The value of severity when you call zenlogger.Info()                                                                 |
| Severity     | Debug           | The value of severity when you call zenlogger.Debug()                                                                |
| Severity     | Warning         | The value of severity when you call zenlogger.Warning()                                                              |
| Severity     | Error           | The value of severity when you call zenlogger.Error()                                                                |
| Severity     | Notice          | The value of severity when you call zenlogger.Notice()                                                               |
| Severity     | Query           | The value of severity when you call zenlogger.Query()                                                                |
| Severity     | Critical        | The value of severity when you call zenlogger.Critical()                                                             |
| Severity     | Alert           | The value of severity when you call zenlogger.Alert()                                                                |
| Severity     | Emergency       | The value of severity when you call zenlogger.Emergency()                                                            |
| DateTime     | Label           | The label of datetime key when log written                                                                           |
| DateTime     | Format          | The output datetime format, using standard GO datetime format                                                        |
| Caller       | Label           | The label of caller key                                                                                              |
| Caller       | Level           | The level of caller key. the default is 0                                                                            |
| Message      | Label           | The label of message key                                                                                             |
| Message      | Title.Label     | The label of title key                                                                                               |
| Message      | Values.Label    | The label of values key                                                                                              |
| BeautifyJson |                 | The beautify config, true to make it beautify formatted, otherwise, set it to false                                  |
| Level        |                 | Minimum log level to execute. Logs below this level are skipped entirely. Default is `LevelDebug` (all logs execute) |
| Output       | Path            | The output path of log files                                                                                         |
| Output       | Format          | The output file name using datetime GO format                                                                        |
| Sensitive    | Enabled         | Enable automatic sensitive field protection by key                                                                   |
| Sensitive    | CaseInsensitive | Match sensitive keys case-insensitively                                                                              |
| Sensitive    | Rules           | Key-to-rule map for automatic protection (`Type`, optional `MaskCount`)                                              |

### Log Level Filtering

Control which logs execute by setting a minimum level. Logs with a level **below** the configured value are skipped entirely with zero cost (no JSON marshaling, no reflection).

```
const (
    LevelDebug     LogLevel = 100  // Access & Query share this tier
    LevelInfo       LogLevel = 200
    LevelNotice     LogLevel = 300
    LevelWarning    LogLevel = 400
    LevelError      LogLevel = 500
    LevelCritical   LogLevel = 600
    LevelAlert      LogLevel = 700
    LevelEmergency  LogLevel = 800
)
```

Example — only execute Warning and above:

```
logger := zenlogger.NewZenlogger()
logger.SetConfig(zenlogger.Config{
    Level: zenlogger.LevelWarning, // skips Debug, Info, Notice, Access, Query
})

logger.Debug("skipped")    // no-op
logger.Info("skipped")     // no-op
logger.Warning("executed") // written
logger.Critical("executed") // written
```

Default is `LevelDebug` — all logs execute unless you configure otherwise.
<br><br><br>

### Sensitive Field Masking

You can mask sensitive values directly per field.

```
logger.Info("db connection",
		zenlogger.ZenField{Key: "username", Value: username, Type: zenlogger.FIRST_MASKED, MaskCount: 2},
		zenlogger.ZenField{Key: "password", Value: password, Type: zenlogger.FULL_MASKED},
)
```

Or use helper constructors for cleaner call sites:

```
logger.Info("db connection",
	zenlogger.Field("username", username),
	zenlogger.MaskedField("password", password, zenlogger.FULL_MASKED, 0),
)
```

Or configure automatic protection by key so caller does not need to set `Type` on every field:

```
logger := zenlogger.NewZenlogger()
logger.SetConfig(zenlogger.Config{
	Sensitive: zenlogger.SensitiveFieldConfig{
		Enabled:         true,
		CaseInsensitive: true,
		Rules: map[string]zenlogger.SensitiveFieldRule{
			"password": {Type: zenlogger.FULL_MASKED},
			"token":    {Type: zenlogger.FIRST_LAST_MASKED, MaskCount: 2},
			"api_key":  {Type: zenlogger.REDACTED},
		},
	},
})

logger.Info("db connection",
	zenlogger.Field("username", username),
	zenlogger.Field("password", password),
	zenlogger.Field("token", token),
	zenlogger.Field("api_key", apiKey),
)
```

Available protection types:

- `FULL_MASKED`: mask the entire value
- `FIRST_MASKED`: mask first `MaskCount` characters
- `LAST_MASKED`: mask last `MaskCount` characters
- `FIRST_LAST_MASKED`: mask first and last `MaskCount` characters
- `REDACTED`: replace value with `[REDACTED]`
- `HASH_SHA256`: replace value with SHA-256 hash (hex encoded)

Notes:

- `MaskCount` defaults to `1` when not set or `<= 0`
- If `MaskCount` is longer than the value length, masking is capped to value length
- Protection is applied to the final string representation of `Value`

Example output:

```
{
	"message": {
		"title": "db connection",
		"values": {
			"username": "**npayuser",
			"password": "********"
		}
	}
}
```

You can set or try it in the link below here:
:point_right: [GO Play](https://goplay.tools/snippet/i9cDLZ8yVHf "GO Play"). :point_left:
<br><br>

This library is released under: [MIT License](https://github.com/zenmuharom/zenlogger/blob/master/LICENSE.txt "MIT License").

[doc-img]: https://pkg.go.dev/badge/github.com/zenmuharom/zenlogger
[doc]: https://pkg.go.dev/github.com/zenmuharom/zenlogger
