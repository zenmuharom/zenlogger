# Zenlogger
[![GoDoc][doc-img]][doc]

This library is built based on internal & personal needs.
So, this lib can be not suit for your need.



**This library support such as**
- Logging to file (including customize naming file)
- Structured log to satisfy Gcloud logs structure 
- Customize structure & keys

# Getting Started

###Sample structured output log

This is the sample of output output when beautifyJson set to `true`.
```
{
  "pid": "42c85c07958f4d6aaa7ea43f8aeeab1d",
  "severity": "INFO",
  "timestamp": "2022-11-18T14:32:45+0700",
  "caller": "middleware.go:45",
  "message": {
    "title": "Incoming request",
    "values": {
      "Payload": {
        "billpayment": {
          "inputBillPayment": {
            "amount": 420500,
            "billNumber": "B0D0230221031",
            "bit61": {
              "merchant_id": "DELIMA",
              "trans_number": "BDO223200163976",
              "outlet_id": "001/002/MNC096",
              "amount": 419000,
              "payment_date": "2022-11-14 17:46:53",
              "payment_code": "001",
              "ref_no": "22"
            },
            "feeAmount": 1500,
            "merchantCode": "MER778",
            "merchantNumber": "+6281000212009",
            "productCode": "030026",
            "terminal": "MER778",
            "timeStamp": "18-11-2022 14:32:46:000000",
            "transactionType": "50",
            "traxId": "202211181432401196",
            "userName": "dev"
          }
        }
      }
    }
  }
}
```

<br/>
    
**How to use**:
This is the sample code of zenlogger usage
```
package main

import (
	"github.com/zenmuharom/zenlogger"
)

type Planet struct {
	Name              string `json:"name"`
	Volume            string `json:"volume"`
	OrbitPeriodInDays int    `json:"orbit_period_in_days"`
}

func main() {
	logger := zenlogger.NewZenlogger()
	logger.SetConfig(zenlogger.Config{BeautifyJson: true})
	logger.Info("Hello World!")
	logger.Error("Hello Error!")
	planet1 := Planet{
		Name:              "Mars",
		Volume:            "1,6318×1011 km³",
		OrbitPeriodInDays: 6867,
	}
	logger.Info("Found 1 planet", zenfield.ZenField{Key: "Earth", Value: planet1})
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
			Label:  "Severity",
			Access: "API",
			Info:   "INFO",
			Debug:  "DEBUG",
			Error:  "ERROR",
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

<br/>
You can set or try it in the link below here: 
https://go.dev/play/p/1wvZJefYor1

This library is released under:
[MIT License](LICENSE.txt).

[doc-img]: https://pkg.go.dev/badge/zenmuharom/zenlogger
[doc]: https://pkg.go.dev/github.com/zenmuharom/zenlogger