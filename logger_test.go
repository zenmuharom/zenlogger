package zenlogger

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
)

func TestLogger(t *testing.T) {
	logger := NewZenlogger()
	logger.Info("test info 0")
	logger.Error("test error 0")
	logger.WithPid("00as0df0").Info("test info 1")
	logger.Error("test error 1")
	// logger.SetFormatDateTime("2006-01-02T15:04:05")
	logger.Info("test info 2")
	logger.Error("test error 2")
	config0 := Config{
		Caller: Caller{
			Label: "trace",
			Level: 0,
		},
		Message: Message{
			Label: "pesan",
			Title: ZenConf{
				Label: "judul",
			},
			Values: ZenConf{
				Label: "isian",
			},
		},
	}
	logger.SetConfig(config0)
	logger.Info("test with json 2", ZenField{
		Key: "trx",
		Value: `{
			"userId": 1,
			"id": 1,
			"title": "delectus aut autem",
			"completed": false
		}`,
	})
	logger.Info("changing pid label to insertId....")
	config := Config{
		Pid: ZenConf{
			Label: "insertId",
		},
	}
	logger.SetConfig(config)
	logger.Info("test info 3")

	logger.Info("changing caller label & level....")
	config = Config{
		Caller: Caller{
			Label: "trace",
			Level: 0,
		},
	}
	logger.SetConfig(config)
	logger.Info("test info 4")

	logger.Info("changing datetime....")
	config = Config{
		DateTime: DateTime{
			Label:  "tanggal",
			Format: "2 Jan 2006 15:04:05",
		},
	}
	logger.SetConfig(config)
	logger.Info("test info 4")

	logger.Info("changing severity label....")
	config = Config{
		Severity: Severity{
			Label:  "seperity",
			Access: "akses",
			Info:   "inpo",
			Debug:  "gedebug",
			Error:  "rusak",
			Query:  "kueri",
		},
		BeautifyJson: true,
	}
	logger.SetConfig(config)
	logger.Info("test info 5")
	logger.Debug("test debug 5")
	logger.Error("test error 5")
	logger.Query("test kueri 5")
	sqlStmt := `SELECT a.NO_VA, a.PASSWD, b.NO_HP, a.ADD_INFO FROM TUSER a, ACCINFO b WHERE a.USERNAME = '` + "FINC0203" + `' AND b.NO_HP = '` + "+62" + `' AND a.NO_VA = b.NO_VA AND a.KODE = '` + "finks020" + `'`
	logger.Query("Test query", ZenField{Value: sqlStmt})

	xmlLog := `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ins="https://www.finnet.com/schemas/mCORE/shared/InSchema.xsd"> <soapenv:Header></soapenv:Header> <soapenv:Body> <ins:mTransaction> <ins:MTI>0200</ins:MTI> <ins:bit1></ins:bit1> <ins:bit2>20221115145359</ins:bit2> <ins:bit3>55</ins:bit3> <ins:bit4>0</ins:bit4> <ins:bit7>20221115145359</ins:bit7> <ins:bit11>798080</ins:bit11> <ins:bit12>145359</ins:bit12> <ins:bit13>20221115</ins:bit13> <ins:bit14>0000</ins:bit14> <ins:bit15>20221116</ins:bit15> <ins:bit18>6012</ins:bit18> <ins:bit22></ins:bit22> <ins:bit25></ins:bit25> <ins:bit26></ins:bit26> <ins:bit27>0</ins:bit27> <ins:bit28></ins:bit28> <ins:bit29></ins:bit29> <ins:bit32>770</ins:bit32> <ins:bit33>770924</ins:bit33> <ins:bit35>1668498839764</ins:bit35> <ins:bit37>145359798080</ins:bit37> <ins:bit38></ins:bit38> <ins:bit41>FINNET</ins:bit41> <ins:bit42>Loket Finnet</ins:bit42> <ins:bit43>Bidakara Pancoran</ins:bit43> <ins:bit44></ins:bit44> <ins:bit45>0</ins:bit45> <ins:bit46>17-11-2011 15:50:01:425931</ins:bit46> <ins:bit47></ins:bit47> <ins:bit48></ins:bit48> <ins:bit49>360</ins:bit49> <ins:bit52></ins:bit52> <ins:bit54>0</ins:bit54> <ins:bit55></ins:bit55> <ins:bit56></ins:bit56> <ins:bit57></ins:bit57> <ins:bit58></ins:bit58> <ins:bit59></ins:bit59> <ins:bit60></ins:bit60> <ins:bit61>{"outlet_desc":"outlet desc from api","outlet_id":"001/002/FNC089","request_type":"add_outlet"}</ins:bit61> <ins:bit62>006MW001</ins:bit62> <ins:bit63></ins:bit63> <ins:bit74></ins:bit74> <ins:bit90></ins:bit90> <ins:bit102></ins:bit102> <ins:bit103>030026</ins:bit103> <ins:bit104>030026</ins:bit104> </ins:mTransaction> </soapenv:Body> </soapenv:Envelope>`
	logger.Query("Test xml", ZenField{Value: xmlLog})
	logger.Info("test int", ZenField{Key: "integer", Value: 69})
	type Anjing struct {
		Name   string
		Gender string
		Age    int
	}

	type AnjingDB struct {
		Name   sql.NullString
		Female sql.NullBool
		Age    sql.NullInt64
	}

	anjing := Anjing{
		Name:   "R",
		Gender: "P",
		Age:    28,
	}
	logger.Info("test struct", ZenField{Key: "si anjing", Value: anjing})

	anjing2 := AnjingDB{
		Name:   sql.NullString{String: "R"},
		Female: sql.NullBool{Bool: true},
		Age:    sql.NullInt64{Int64: 28},
	}
	logger.Info("test struct DB", ZenField{Key: "si anjing DB", Value: anjing2})

	anjings := make([]AnjingDB, 0)
	anjing3 := AnjingDB{
		Name:   sql.NullString{String: "R"},
		Female: sql.NullBool{Bool: true},
		Age:    sql.NullInt64{Int64: 28},
	}
	anjings = append(anjings, anjing2)
	anjings = append(anjings, anjing3)
	logger.Info("test struct DB", ZenField{Key: "si anjings DB list", Value: anjings})

	config.Output.Path = "logs"
	config.Output.Format = "2006"
	config.BeautifyJson = true
	logger.SetConfig(config)
	fmt.Println(config.Output.Path)
	logger.Info("test info 6")
	logger.Debug("test debug 6")
	logger.Error("test error 6")
	logger.Query("test kueri 6")

	var args interface{}
	argsData := `{"trim": [{"substring": ["$middle_response_id", 13, 2]}]}`

	err := json.Unmarshal([]byte(argsData), &args)
	if err != nil {
		logger.Error(err.Error())
	} else {
		logger.Info("args", ZenField{Key: "data", Value: args})
	}
}
