package zenlogger

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/heisenbergbrbat/uuid"
)

var re = regexp.MustCompile(`^(.*/)?(?:$|(.+?)(?:(\.[^.]*$)|$))`)
var replacer = strings.NewReplacer("\r", "", "\n", "")

type ZenLog []ZenField

type Zenlogger interface {
	SetConfig(config Config)
	GetConfig() Config
	WithPid(pid string) Zenlogger
	GetPid() string
	Access(message string, fields ...ZenField) string
	Info(message string, fields ...ZenField) string
	Query(message string, fields ...ZenField) string
	Debug(message string, fields ...ZenField) string
	Warning(message string, fields ...ZenField) string
	Error(message string, fields ...ZenField) string
}

type DefaultZenlogger struct {
	pid    string
	config Config
}

func NewZenlogger(pid ...string) Zenlogger {
	pid0 := ""
	if len(pid) == 0 {
		pid0 = strings.Replace(uuid.New().String(), "-", "", -1)
	} else {
		pid0 = pid[0]
	}

	config := Config{
		Pid: ZenConf{
			Label: "pid",
		},
		Severity: Severity{
			Label:   "severity",
			Access:  DEFAULT_ACCESS,
			Info:    DEFAULT_INFO,
			Debug:   DEFAULT_DEBUG,
			Warning: DFEAULT_WARNING,
			Error:   DEFAULT_ERROR,
			Query:   DEFAULT_QUERY,
		},
		DateTime: DateTime{
			Label:  "timestamp",
			Format: "2006-01-02T15:04:05-0700",
		},
		Caller: Caller{
			Label: "caller",
		},
		Message: Message{
			Label: "message",
			Title: ZenConf{
				Label: "title",
			},
			Values: ZenConf{
				Label: "values",
			},
		},
		BeautifyJson: false,
	}

	return &DefaultZenlogger{
		pid:    pid0,
		config: config,
	}
}

func (zenlog *DefaultZenlogger) WithPid(pid string) Zenlogger {
	zenlog.pid = pid
	return zenlog
}

func (zenlog *DefaultZenlogger) GetPid() string {
	return zenlog.pid
}

func (zenlog *DefaultZenlogger) write(Type string, msgStr string, fields ...ZenField) (log string) {
	_, file, no, _ := runtime.Caller(zenlog.config.Caller.Level + 2)
	fileNameOnly := re.FindStringSubmatch(file)[2]
	caller := fmt.Sprintf("%v.go:%v", fileNameOnly, no)

	config := zenlog.config

	// parse log structure
	newlog := ZenLog{
		{config.Pid.Label, zenlog.pid},
		{config.Severity.Label, Type},
		{config.DateTime.Label, time.Now().Format(zenlog.config.DateTime.Format)},
		{config.Caller.Label, caller},
	}

	// parse message structure
	if len(fields) > 0 {
		newlog = append(newlog, ZenField{config.Message.Label, ZenLog{
			{config.Message.Title.Label, msgStr},
			{config.Message.Values.Label, zenlog.parse(fields...)},
		}})
	} else {
		newlog = append(newlog, ZenField{config.Message.Label, msgStr})
	}

	var logStr []byte

	logStr, _ = jsonMarshal(newlog, zenlog.config.BeautifyJson)

	if zenlog.config.Output.Path == "" {
		fmt.Print(string(logStr))
	} else {
		fileName := fmt.Sprintf("%s.log", time.Now().Format(zenlog.config.Output.Format))
		filePath := filepath.Join(zenlog.config.Output.Path, fileName)
		err := os.MkdirAll(zenlog.config.Output.Path, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
		}

		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer file.Close()
		if err != nil {
			fmt.Println(err.Error())
		}

		datawriter := bufio.NewWriter(file)
		defer datawriter.Flush()
		_, err = datawriter.Write(logStr)
		if err != nil {
			fmt.Println(err.Error())
		}
		_, err = datawriter.WriteString("\n")
		if err != nil {
			fmt.Println(err.Error())
		}

	}

	log = string(logStr)

	return
}
