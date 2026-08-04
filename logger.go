package zenlogger

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/heisenbergbrbat/uuid"
)

var replacer = strings.NewReplacer("\r", "", "\n", "")

type ZenLog []ZenField

type Zenlogger interface {
	SetConfig(config Config)
	GetConfig() Config
	WithPid(pid string) Zenlogger
	GetPid() string
	Access(message string, fields ...ZenField) string
	Info(message string, fields ...ZenField) string
	Notice(message string, fields ...ZenField) string
	Query(message string, fields ...ZenField) string
	Debug(message string, fields ...ZenField) string
	Warning(message string, fields ...ZenField) string
	Error(message string, fields ...ZenField) string
	Critical(message string, fields ...ZenField) string
	Alert(message string, fields ...ZenField) string
	Emergency(message string, fields ...ZenField) string
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
			Label:     "severity",
			Access:    DEFAULT_ACCESS,
			Info:      DEFAULT_INFO,
			Debug:     DEFAULT_DEBUG,
			Notice:    DEFAULT_NOTICE,
			Warning:   DFEAULT_WARNING,
			Error:     DEFAULT_ERROR,
			Query:     DEFAULT_QUERY,
			Critical:  DEFAULT_CRITICAL,
			Alert:     DEFAULT_ALERT,
			Emergency: DEFAULT_EMERGENCY,
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
		Sensitive: SensitiveFieldConfig{
			Enabled:         false,
			CaseInsensitive: true,
			Rules:           nil,
		},
		BeautifyJson: false,
		Level:        LevelDebug,
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
	pc, file, no, _ := runtime.Caller(zenlog.config.Caller.Level + 2)
	caller := fmt.Sprintf("%s:%d", formatCallerPackagePath(pc, file), no)

	config := zenlog.config

	// parse log structure
	newlog := ZenLog{
		{Key: config.Pid.Label, Value: zenlog.pid},
		{Key: config.Severity.Label, Value: Type},
		{Key: config.DateTime.Label, Value: time.Now().Format(zenlog.config.DateTime.Format)},
		{Key: config.Caller.Label, Value: caller},
	}

	// parse message structure
	if len(fields) > 0 {
		newlog = append(newlog, ZenField{Key: config.Message.Label, Value: ZenLog{
			{Key: config.Message.Title.Label, Value: msgStr},
			{Key: config.Message.Values.Label, Value: zenlog.parse(fields...)},
		}})
	} else {
		newlog = append(newlog, ZenField{Key: config.Message.Label, Value: msgStr})
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

func formatCallerPackagePath(pc uintptr, filePath string) string {
	baseFile := filepath.Base(filePath)
	if baseFile == "" {
		baseFile = filePath
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return formatCallerPath(filePath)
	}

	funcName := fn.Name()
	if funcName == "" {
		return formatCallerPath(filePath)
	}

	segment := funcName
	if idx := strings.LastIndex(segment, "/"); idx >= 0 {
		segment = segment[idx+1:]
	}

	pkgName := segment
	if idx := strings.Index(pkgName, "."); idx >= 0 {
		pkgName = pkgName[:idx]
	}

	pkgName = strings.TrimSuffix(pkgName, "_test")
	if pkgName == "" {
		return formatCallerPath(filePath)
	}

	return pkgName + "/" + baseFile
}

func formatCallerPath(absPath string) string {
	if absPath == "" {
		return ""
	}

	if !filepath.IsAbs(absPath) {
		cleanPath := filepath.ToSlash(filepath.Clean(absPath))
		if strings.Contains(cleanPath, "/") {
			return cleanPath
		}

		if resolvedPath, ok := resolveCallerPathFromFileName(cleanPath); ok {
			return resolvedPath
		}

		return cleanPath
	}

	absPath = normalizePath(absPath)

	wd, err := os.Getwd()
	if err == nil {
		wd = normalizePath(wd)

		if moduleRoot, found := findGoModuleRoot(wd); found {
			moduleRoot = normalizePath(moduleRoot)
			relPath, relErr := filepath.Rel(moduleRoot, absPath)
			if relErr == nil && relPath != "" && relPath != ".." && !strings.HasPrefix(relPath, ".."+string(os.PathSeparator)) {
				return filepath.ToSlash(relPath)
			}
		}

		relPath, relErr := filepath.Rel(wd, absPath)
		if relErr == nil && relPath != "" && relPath != ".." && !strings.HasPrefix(relPath, ".."+string(os.PathSeparator)) {
			return filepath.ToSlash(relPath)
		}
	}

	return filepath.ToSlash(filepath.Base(absPath))
}

func resolveCallerPathFromFileName(fileName string) (string, bool) {
	wd, err := os.Getwd()
	if err != nil {
		return "", false
	}

	moduleRoot, found := findGoModuleRoot(normalizePath(wd))
	if !found {
		return "", false
	}

	matches := make([]string, 0)
	_ = filepath.WalkDir(moduleRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}

		if d.IsDir() {
			switch d.Name() {
			case ".git", "vendor", "logs":
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Base(path) != fileName {
			return nil
		}

		relPath, relErr := filepath.Rel(moduleRoot, path)
		if relErr == nil && relPath != "" {
			matches = append(matches, filepath.ToSlash(relPath))
		}

		return nil
	})

	if len(matches) == 1 {
		return matches[0], true
	}

	return "", false
}

func normalizePath(path string) string {
	resolvedPath, err := filepath.EvalSymlinks(path)
	if err == nil {
		return resolvedPath
	}

	return path
}

func findGoModuleRoot(startDir string) (string, bool) {
	current := startDir
	for {
		if _, err := os.Stat(filepath.Join(current, "go.mod")); err == nil {
			return current, true
		}

		parent := filepath.Dir(current)
		if parent == current {
			return "", false
		}

		current = parent
	}
}
