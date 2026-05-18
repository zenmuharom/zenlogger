package zenlogger

type LogLevel int

const (
	LevelDebug     LogLevel = 100
	LevelInfo      LogLevel = 200
	LevelAccess    LogLevel = 200
	LevelQuery     LogLevel = 200
	LevelNotice    LogLevel = 300
	LevelWarning   LogLevel = 400
	LevelError     LogLevel = 500
	LevelCritical  LogLevel = 600
	LevelAlert     LogLevel = 700
	LevelEmergency LogLevel = 800
)

const (
	DEFAULT_ACCESS    = "Access"
	DEFAULT_INFO      = "Info"
	DEFAULT_NOTICE    = "Notice"
	DFEAULT_WARNING   = "Warning"
	DEFAULT_ERROR     = "Error"
	DEFAULT_DEBUG     = "Debug"
	DEFAULT_QUERY     = "Query"
	DEFAULT_CRITICAL  = "Critical"
	DEFAULT_ALERT     = "Alert"
	DEFAULT_EMERGENCY = "Emergency"
)

func (zenlog *DefaultZenlogger) Access(message string, fields ...ZenField) string {
	if LevelAccess < zenlog.config.Level {
		return ""
	}
	return zenlog.write(zenlog.config.Severity.Access, message, fields...)
}

func (zenlog *DefaultZenlogger) Info(message string, fields ...ZenField) string {
	if LevelInfo < zenlog.config.Level {
		return ""
	}
	return zenlog.write(zenlog.config.Severity.Info, message, fields...)
}

func (zenlog *DefaultZenlogger) Query(message string, fields ...ZenField) string {
	if LevelQuery < zenlog.config.Level {
		return ""
	}
	return zenlog.write(zenlog.config.Severity.Query, message, fields...)
}

func (zenlog *DefaultZenlogger) Debug(message string, fields ...ZenField) string {
	if LevelDebug < zenlog.config.Level {
		return ""
	}
	if zenlog.config.Production {
		return ""
	}
	return zenlog.write(zenlog.config.Severity.Debug, message, fields...)
}

func (zenlog *DefaultZenlogger) Warning(message string, fields ...ZenField) string {
	if LevelWarning < zenlog.config.Level {
		return ""
	}
	return zenlog.write(zenlog.config.Severity.Warning, message, fields...)
}

func (zenlog *DefaultZenlogger) Error(message string, fields ...ZenField) string {
	if LevelError < zenlog.config.Level {
		return ""
	}
	return zenlog.write(zenlog.config.Severity.Error, message, fields...)
}

func (zenlog *DefaultZenlogger) Notice(message string, fields ...ZenField) string {
	if LevelNotice < zenlog.config.Level {
		return ""
	}
	return zenlog.write(zenlog.config.Severity.Notice, message, fields...)
}

func (zenlog *DefaultZenlogger) Critical(message string, fields ...ZenField) string {
	if LevelCritical < zenlog.config.Level {
		return ""
	}
	return zenlog.write(zenlog.config.Severity.Critical, message, fields...)
}

func (zenlog *DefaultZenlogger) Alert(message string, fields ...ZenField) string {
	if LevelAlert < zenlog.config.Level {
		return ""
	}
	return zenlog.write(zenlog.config.Severity.Alert, message, fields...)
}

func (zenlog *DefaultZenlogger) Emergency(message string, fields ...ZenField) string {
	if LevelEmergency < zenlog.config.Level {
		return ""
	}
	return zenlog.write(zenlog.config.Severity.Emergency, message, fields...)
}
