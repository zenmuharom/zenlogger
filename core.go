package zenlogger

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
	"strconv"
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
		case reflect.Struct:
			realVal = zenlog.unmarshalStruct(vRef.Index(i).Interface())
		case reflect.Slice, reflect.Array:
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
	v := reflect.ValueOf(value)
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
		realVal = nil
	}
	return
}

func (zenlog *DefaultZenlogger) unmarshalStruct(structToParse interface{}) map[string]interface{} {
	parsedStruct := make(map[string]interface{})
	fieldValues := reflect.ValueOf(structToParse)
	fields := fieldValues.Type()
	for i := 0; i < fieldValues.NumField(); i++ {
		tag := fields.Field(i).Tag.Get("json")
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
		default:
			parsedStruct[tag] = zenlog.unmarshalInterface(f)
		}
	}

	return parsedStruct
}

func isValidXML(s string) bool {
	return xml.Unmarshal([]byte(s), new(interface{})) == nil
}

func (zenlog *DefaultZenlogger) Parse(fields ...ZenField) map[string]interface{} {

	parsed := make(map[string]interface{})
	for i, field := range fields {

		// check if key is null then add index as key
		if field.Key == "" {
			field.Key = strconv.Itoa(i)
		}

		parsed[field.Key] = zenlog.unmarshalInterface(field.Value)
	}
	return parsed
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
