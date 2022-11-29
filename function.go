package zenlogger

const (
	DEFAULT_ACCESS = "Access"
	DEFAULT_INFO   = "Info"
	DEFAULT_ERROR  = "Error"
	DEFAULT_DEBUG  = "Debug"
	DEFAULT_QUERY  = "Query"
)

func (zenlog *DefaultZenlogger) Access(message string, fields ...ZenField) {
	zenlog.Write(zenlog.config.Severity.Access, message, fields...)
}

func (zenlog *DefaultZenlogger) Info(message string, fields ...ZenField) {
	zenlog.Write(zenlog.config.Severity.Info, message, fields...)
}

func (zenlog *DefaultZenlogger) Query(message string, fields ...ZenField) {
	zenlog.Write(zenlog.config.Severity.Query, message, fields...)
}

func (zenlog *DefaultZenlogger) Debug(message string, fields ...ZenField) {
	if !zenlog.config.Production {
		zenlog.Write(zenlog.config.Severity.Debug, message, fields...)
	}
}

func (zenlog *DefaultZenlogger) Error(message string, fields ...ZenField) {
	zenlog.Write(zenlog.config.Severity.Error, message, fields...)
}
