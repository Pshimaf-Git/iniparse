package iniparse

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func encode(v any) (*iniFile, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return nil, fmt.Errorf("encode: expected non-nil pointer, got %T", v)
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("encode: expected struct pointer, got %s", rv.Kind())
	}

	ini := newIniFile()
	if err := encodeStruct(ini, "", rv); err != nil {
		return nil, err
	}
	return ini, nil
}

func encodeStruct(ini *iniFile, sec string, rv reflect.Value) error {
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fieldVal := rv.Field(i)

		if !field.IsExported() {
			continue
		}

		iniTag := field.Tag.Get("ini")
		if iniTag == "" || iniTag == "-" {
			continue
		}

		tagSection, key, _ := parseIniTag(iniTag)

		useSection := sec
		if tagSection != "" {
			useSection = tagSection
		}

		// struct field — recurse, section = tag key
		if fieldVal.Kind() == reflect.Struct {
			if err := encodeStruct(ini, key, fieldVal); err != nil {
				return err
			}
			continue
		}

		// pointer to struct — recurse if non-nil
		if fieldVal.Kind() == reflect.Ptr {
			if fieldVal.IsNil() {
				continue
			}
			if fieldVal.Elem().Kind() == reflect.Struct {
				if err := encodeStruct(ini, key, fieldVal.Elem()); err != nil {
					return err
				}
				continue
			}
		}

		val, err := formatValue(fieldVal)
		if err != nil {
			return fmt.Errorf("encode: field %s: %w", field.Name, err)
		}
		ini.set(useSection, key, val)
	}
	return nil
}

func formatValue(rv reflect.Value) (string, error) {
	switch rv.Kind() {
	case reflect.String:
		return rv.String(), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10), nil

	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'f', -1, rv.Type().Bits()), nil

	case reflect.Bool:
		return strconv.FormatBool(rv.Bool()), nil

	case reflect.Slice:
		if rv.IsNil() {
			return "", nil
		}
		parts := make([]string, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			s, err := formatValue(rv.Index(i))
			if err != nil {
				return "", err
			}
			parts[i] = s
		}
		return strings.Join(parts, ","), nil

	default:
		return "", fmt.Errorf("unsupported type: %s", rv.Kind())
	}
}
