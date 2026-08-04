package zenlogger

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func (zenlog *DefaultZenlogger) unmarshalMap(v reflect.Value) map[string]interface{} {
	newMap := make(map[string]interface{}, 0)
	iter := v.MapRange()
	for iter.Next() {
		iV := iter.Value()
		switch iV.Kind() {
		case reflect.String:
			newMap[iter.Key().String()] = iter.Value().String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			newMap[iter.Key().String()] = iter.Value().Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			newMap[iter.Key().String()] = iter.Value().Uint()
		case reflect.Float32, reflect.Float64:
			newMap[iter.Key().String()] = iter.Value().Float()
		case reflect.Bool:
			newMap[iter.Key().String()] = iter.Value().Bool()
		case reflect.Struct:
			newMap[iter.Key().String()] = zenlog.unmarshalStruct(iter.Value().Interface())
		case reflect.Slice:
			newMap[iter.Key().String()] = zenlog.unmarshalSliceAndArray(iter.Value())
		case reflect.Map:
			newMap[iter.Key().String()] = zenlog.unmarshalMap(iter.Value())
		case reflect.Interface:
			newMap[iter.Key().String()] = zenlog.unmarshalInterface(iter.Value().Interface())
		default:
			// Fallback for any unhandled types
			newMap[iter.Key().String()] = fmt.Sprintf("%v", iter.Value().Interface())
		}
	}

	return newMap
}

func (zenlog *DefaultZenlogger) unmarshalSliceAndArray(vRef reflect.Value) []interface{} {
	newMap := make([]interface{}, 0)

	for i := 0; i < vRef.Len(); i++ {
		iV := reflect.ValueOf(vRef.Index(i))

		var realVal interface{}
		switch iV.Kind() {
		case reflect.String:
			if json.Valid([]byte(vRef.Index(i).String())) {
				json.Unmarshal([]byte(vRef.Index(i).String()), &realVal)
			} else if isValidXML(vRef.Index(i).String()) {
				realVal = fmt.Sprintf("%s", vRef.Index(i))
			} else {
				realVal = replacer.Replace(fmt.Sprintf("%v", vRef.Index(i).String()))
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			realVal = vRef.Index(i).Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			realVal = vRef.Index(i).Uint()
		case reflect.Float32, reflect.Float64:
			realVal = vRef.Index(i).Float()
		case reflect.Bool:
			realVal = vRef.Index(i).Bool()
		case reflect.Interface:
			realVal = vRef.Index(i).Interface()
		case reflect.Array:
			realVal = zenlog.unmarshalSliceAndArray(vRef.Index(i))
		case reflect.Struct:
			realVal = zenlog.unmarshalStruct(vRef.Index(i).Interface())
		case reflect.Slice:
			realVal = zenlog.unmarshalSliceAndArray(vRef.Index(i))
		case reflect.Map:
			realVal = zenlog.unmarshalMap(vRef.Index(i))
		default:
			realVal = nil
		}
		newMap = append(newMap, realVal)
	}

	return newMap
}

func (zenlog *DefaultZenlogger) unmarshalInterface(value interface{}) (realVal interface{}) {
	if value == nil {
		return nil
	}

	// Handle json.Number type specifically
	if jsonNum, ok := value.(json.Number); ok {
		return jsonNum.String()
	}

	// Handle time.Time type specifically
	if t, ok := value.(time.Time); ok {
		return t.Format(time.RFC3339)
	}

	v := reflect.ValueOf(value)
	if !v.IsValid() {
		return nil
	}

	switch v.Kind() {
	case reflect.String:
		if json.Valid([]byte(value.(string))) {
			json.Unmarshal([]byte(value.(string)), &value)
		} else if isValidXML(value.(string)) {
			value = fmt.Sprintf("%s", value)
		} else {
			value = replacer.Replace(fmt.Sprintf("%v", v.String()))
		}
		realVal = value
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		realVal = v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		realVal = v.Uint()
	case reflect.Float32, reflect.Float64:
		realVal = v.Float()
	case reflect.Bool:
		realVal = v.Bool()
	case reflect.Struct:
		realVal = zenlog.unmarshalStruct(value)
	case reflect.Slice, reflect.Array:
		realVal = zenlog.unmarshalSliceAndArray(v)
	case reflect.Map:
		realVal = zenlog.unmarshalMap(v)
	default:
		// Fallback for any unhandled types - convert to string representation
		realVal = fmt.Sprintf("%v", value)
	}
	return
}

func (zenlog *DefaultZenlogger) unmarshalStruct(structToParse interface{}) interface{} {

	v := reflect.ValueOf(structToParse)
	fieldValues := reflect.ValueOf(structToParse)

	var parsedVal interface{}

	switch v.Kind() {
	case reflect.String:
		parsedVal = fieldValues.String()
		return parsedVal
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsedVal = fieldValues.Int()
		return parsedVal
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parsedVal = fieldValues.Uint()
		return parsedVal
	case reflect.Float32, reflect.Float64:
		parsedVal = fieldValues.Float()
		return parsedVal
	case reflect.Bool:
		parsedVal = fieldValues.Bool()
		return parsedVal
	case reflect.Slice:
		parsedVal = zenlog.unmarshalSliceAndArray(fieldValues)
		return parsedVal
	case reflect.Map:
		parsedVal = zenlog.unmarshalMap(fieldValues)
		return parsedVal
	case reflect.Interface:
		parsedVal = zenlog.unmarshalInterface(fieldValues.Interface())
		return parsedVal
	default:
		fields := fieldValues.Type()
		parsedStruct := make(map[string]interface{})
		for i := 0; i < fieldValues.NumField(); i++ {
			// Skip unexported fields
			if !fieldValues.Field(i).CanInterface() {
				continue
			}

			tag := fields.Field(i).Tag.Get("json")
			if tag == "" {
				tag = fields.Field(i).Tag.Get("db")
			}
			if tag == "" {
				tag = fields.Field(i).Tag.Get("zenlogger")
			}
			if tag == "" {
				tag = fields.Field(i).Name
			}

			refValue := fieldValues.Field(i).Interface()

			switch f := refValue.(type) {
			case sql.NullBool:
				parsedStruct[tag] = f.Bool
			case sql.NullByte:
				parsedStruct[tag] = f.Byte
			case sql.NullInt16:
				parsedStruct[tag] = f.Int16
			case sql.NullInt32:
				parsedStruct[tag] = f.Int32
			case sql.NullInt64:
				parsedStruct[tag] = f.Int64
			case sql.NullString:
				parsedStruct[tag] = f.String
			case sql.NullFloat64:
				parsedStruct[tag] = f.Float64
			case sql.NullTime:
				parsedStruct[tag] = f.Time
			case time.Time:
				parsedStruct[tag] = f.Format(time.RFC3339)
			case json.Number:
				// Handle json.Number type
				parsedStruct[tag] = f.String()
			default:
				parsedStruct[tag] = zenlog.unmarshalInterface(f)
			}
		}

		return parsedStruct
	}
}

func isValidXML(s string) bool {
	return xml.Unmarshal([]byte(s), new(interface{})) == nil
}

func (zenlog *DefaultZenlogger) parse(fields ...ZenField) map[string]interface{} {

	parsed := make(map[string]interface{})
	for i, field := range fields {

		var value any

		// check if key is null then add index as key
		if field.Key == "" {
			field.Key = strconv.Itoa(i)
		}

		// Use reflection to check if the value in the interface is a pointer
		refValue := reflect.ValueOf(field.Value)
		if !refValue.IsValid() {
			value = nil
		} else if refValue.Kind() == reflect.Ptr {
			// Check if the pointer is nil
			if refValue.IsNil() {
				value = nil
			} else {
				// It's a pointer and not nil, so we can access its value
				value = refValue.Elem().Interface()
			}
		} else {
			value = field.Value
		}

		parsedValue := zenlog.unmarshalInterface(value)
		parsed[field.Key] = zenlog.applyFieldMask(parsedValue, field)
	}
	return parsed
}

func (zenlog *DefaultZenlogger) applyFieldMask(value interface{}, field ZenField) interface{} {
	fieldType := field.Type
	maskCount := field.MaskCount

	if fieldType == "" {
		rule, ok := zenlog.resolveSensitiveRule(field.Key)
		if ok {
			fieldType = rule.Type
			if maskCount <= 0 {
				maskCount = rule.MaskCount
			}
		}
	}

	if fieldType == "" {
		return value
	}

	strValue := fmt.Sprintf("%v", value)
	if strValue == "<nil>" {
		return value
	}

	return protectString(strValue, fieldType, maskCount)
}

func (zenlog *DefaultZenlogger) resolveSensitiveRule(key string) (SensitiveFieldRule, bool) {
	sensitiveConf := zenlog.config.Sensitive
	if !sensitiveConf.Enabled || key == "" || sensitiveConf.Rules == nil {
		return SensitiveFieldRule{}, false
	}

	if rule, ok := sensitiveConf.Rules[key]; ok {
		return rule, true
	}

	if !sensitiveConf.CaseInsensitive {
		return SensitiveFieldRule{}, false
	}

	if rule, ok := sensitiveConf.Rules[strings.ToLower(key)]; ok {
		return rule, true
	}

	for ruleKey, rule := range sensitiveConf.Rules {
		if strings.EqualFold(ruleKey, key) {
			return rule, true
		}
	}

	return SensitiveFieldRule{}, false
}

func protectString(value string, maskType MaskType, maskCount int) string {
	if value == "" {
		return value
	}

	runes := []rune(value)
	length := len(runes)

	switch maskType {
	case FULL_MASKED:
		return strings.Repeat("*", length)
	case FIRST_MASKED:
		count := normalizedMaskCount(maskCount, length)
		return strings.Repeat("*", count) + string(runes[count:])
	case LAST_MASKED:
		count := normalizedMaskCount(maskCount, length)
		return string(runes[:length-count]) + strings.Repeat("*", count)
	case FIRST_LAST_MASKED:
		count := normalizedMaskCount(maskCount, length)
		if count*2 >= length {
			return strings.Repeat("*", length)
		}
		return strings.Repeat("*", count) + string(runes[count:length-count]) + strings.Repeat("*", count)
	case REDACTED:
		return "[REDACTED]"
	case HASH_SHA256:
		hash := sha256.Sum256([]byte(value))
		return hex.EncodeToString(hash[:])
	default:
		return value
	}
}

func normalizedMaskCount(maskCount int, max int) int {
	if maskCount <= 0 {
		return 1
	}
	if maskCount > max {
		return max
	}
	return maskCount
}

func jsonMarshal(t interface{}, indentation bool) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	if indentation {
		encoder.SetIndent(" ", "\t")
	}
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

func (zmap ZenLog) MarshalJSON() ([]byte, error) {

	var buf bytes.Buffer

	buf.WriteString("{")
	for i, kv := range zmap {
		if i != 0 {
			buf.WriteString(",")
		}
		// marshal key
		key, err := json.Marshal(kv.Key)
		if err != nil {
			return nil, err
		}
		buf.Write(key)
		buf.WriteString(":")
		// marshal value
		val, err := jsonMarshal(kv.Value, false)
		if err != nil {
			return nil, err
		}
		buf.Write(val)
	}

	buf.WriteString("}")
	return buf.Bytes(), nil
}
