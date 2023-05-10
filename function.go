package zenlogger

const (
	DEFAULT_ACCESS  = "Access"
	DEFAULT_INFO    = "Info"
	DFEAULT_WARNING = "Warning"
	DEFAULT_ERROR   = "Error"
	DEFAULT_DEBUG   = "Debug"
	DEFAULT_QUERY   = "Query"
)

func (zenlog *DefaultZenlogger) Access(message string, fields ...ZenField) string {
	return zenlog.write(zenlog.config.Severity.Access, message, fields...)
}

func (zenlog *DefaultZenlogger) Info(message string, fields ...ZenField) string {
	return zenlog.write(zenlog.config.Severity.Info, message, fields...)
}

func (zenlog *DefaultZenlogger) Query(message string, fields ...ZenField) string {
	return zenlog.write(zenlog.config.Severity.Query, message, fields...)
}

func (zenlog *DefaultZenlogger) Debug(message string, fields ...ZenField) string {
	str := ""
	if !zenlog.config.Production {
		str = zenlog.write(zenlog.config.Severity.Debug, message, fields...)
	}

	return str
}

func (zenlog *DefaultZenlogger) Warning(message string, fields ...ZenField) string {
	return zenlog.write(zenlog.config.Severity.Warning, message, fields...)
}

func (zenlog *DefaultZenlogger) Error(message string, fields ...ZenField) string {
	return zenlog.write(zenlog.config.Severity.Error, message, fields...)
}
