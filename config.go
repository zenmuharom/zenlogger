package zenlogger

type ZenConf struct {
	Label string
}

type Severity struct {
	Label  string
	Access string
	Info   string
	Debug  string
	Query  string
	Error  string
}

type Message struct {
	Label  string
	Title  ZenConf
	Values ZenConf
}

type DateTime struct {
	Label  string
	Format string
}

type Caller struct {
	Label string
	Level int
}

type Output struct {
	Path       string
	Format     string
	BufferSize int64
}

type Config struct {
	Pid          ZenConf
	Severity     Severity
	DateTime     DateTime
	Caller       Caller
	Message      Message
	BeautifyJson bool
	Production   bool
	Output       Output
}

func (zenlog *DefaultZenlogger) SetConfig(newConfig Config) {

	// config PID
	zenlog.ConfPid(newConfig.Pid)

	// config Severity
	zenlog.ConfSeverity(newConfig.Severity)

	// config DateTime
	zenlog.ConfDateTime(newConfig.DateTime)

	// config Caller
	zenlog.ConfCaller(newConfig.Caller)

	// config Message
	zenlog.ConfMessage(newConfig.Message)

	// config BeautifyJson
	zenlog.config.BeautifyJson = newConfig.BeautifyJson

	// config production
	zenlog.config.Production = newConfig.Production

	// config output
	zenlog.config.Output.Path = newConfig.Output.Path
	zenlog.config.Output.Format = newConfig.Output.Format
	zenlog.config.Output.BufferSize = newConfig.Output.BufferSize
}

func (zenlog *DefaultZenlogger) ConfPid(pidConf ZenConf) {
	if pidConf.Label != "" {
		zenlog.config.Pid = pidConf
	}
}

func (zenlog *DefaultZenlogger) ConfSeverity(ServerityConf Severity) {
	if ServerityConf.Label != "" {
		zenlog.config.Severity.Label = ServerityConf.Label
	}

	if ServerityConf.Access != "" {
		zenlog.config.Severity.Access = ServerityConf.Access
	}

	if ServerityConf.Info != "" {
		zenlog.config.Severity.Info = ServerityConf.Info
	}

	if ServerityConf.Debug != "" {
		zenlog.config.Severity.Debug = ServerityConf.Debug
	}

	if ServerityConf.Error != "" {
		zenlog.config.Severity.Error = ServerityConf.Error
	}

	if ServerityConf.Query != "" {
		zenlog.config.Severity.Query = ServerityConf.Query
	}
}

func (zenlog *DefaultZenlogger) ConfDateTime(DateTimeConf DateTime) {
	if DateTimeConf.Label != "" {
		zenlog.config.DateTime.Label = DateTimeConf.Label
	}

	if DateTimeConf.Format != "" {
		zenlog.config.DateTime.Format = DateTimeConf.Format
	}

}

func (zenlog *DefaultZenlogger) ConfCaller(CallerConf Caller) {
	if CallerConf.Label != "" {
		zenlog.config.Caller.Label = CallerConf.Label
	}

	zenlog.config.Caller.Level = CallerConf.Level

}

func (zenlog *DefaultZenlogger) ConfMessage(MessageConf Message) {
	if MessageConf.Label != "" {
		zenlog.config.Message.Label = MessageConf.Label
	}

	if MessageConf.Title.Label != "" {
		zenlog.config.Message.Title = MessageConf.Title
	}

	if MessageConf.Values.Label != "" {
		zenlog.config.Message.Values = MessageConf.Values
	}

}

func (zenlog *DefaultZenlogger) GetConfig() Config {
	return zenlog.config
}
