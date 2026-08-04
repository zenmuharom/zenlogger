package zenlogger

type MaskType string

const (
	FULL_MASKED       MaskType = "FULL_MASKED"
	FIRST_MASKED      MaskType = "FIRST_MASKED"
	LAST_MASKED       MaskType = "LAST_MASKED"
	FIRST_LAST_MASKED MaskType = "FIRST_LAST_MASKED"
	REDACTED          MaskType = "REDACTED"
	HASH_SHA256       MaskType = "HASH_SHA256"
)

const (
	MaskFull       MaskType = FULL_MASKED
	MaskFirst      MaskType = FIRST_MASKED
	MaskLast       MaskType = LAST_MASKED
	MaskFirstLast  MaskType = FIRST_LAST_MASKED
	MaskRedacted   MaskType = REDACTED
	MaskHashSHA256 MaskType = HASH_SHA256
)

type ZenField struct {
	Key       string
	Value     interface{}
	Type      MaskType
	MaskCount int
}

func Field(key string, value interface{}) ZenField {
	return ZenField{Key: key, Value: value}
}

func MaskedField(key string, value interface{}, maskType MaskType, maskCount int) ZenField {
	return ZenField{Key: key, Value: value, Type: maskType, MaskCount: maskCount}
}

type ZenMessage struct {
	Title  string                 `json:"title,omitempty"`
	Values map[string]interface{} `json:"values,omitempty"`
}
