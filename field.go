package zenlogger

type ZenField struct {
	Key   string
	Value interface{}
}

type ZenMessage struct {
	Title  string                 `json:"title,omitempty"`
	Values map[string]interface{} `json:"values,omitempty"`
}
