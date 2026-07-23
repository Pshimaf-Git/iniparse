// Package iniparse
package iniparse

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func decode(ini *iniFile, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("decode: expected non-nil pointer, got %T", v)
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("decode: expected struct pointer, got %s", rv.Kind())
	}
	return decodeStruct(ini, "", rv)
}

func decodeStruct(ini *iniFile, sec string, rv reflect.Value) error {
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

		// Parse tag: `ini:"key"` or `ini:"section.key"` or `ini:"section,optional"`
		tagSection, key, opts := parseIniTag(iniTag)
		_ = opts

		// If tag specifies a section, use it; otherwise use the passed section
		useSection := sec
		if tagSection != "" {
			useSection = tagSection
		}

		// If the field is a struct itself, recurse with the correct section
		if fieldVal.Kind() == reflect.Struct {
			// For tags like "database" (no dot), recurse into section "database"
			// For tags like "server.host" (dot), recurse into section "server"
			if tagSection != "" {
				if err := decodeStruct(ini, tagSection, fieldVal); err != nil {
					return err
				}
			} else {
				if err := decodeStruct(ini, key, fieldVal); err != nil {
					return err
				}
			}
			continue
		}

		// If the field is a pointer to a struct, recurse
		if fieldVal.Kind() == reflect.Ptr {
			if fieldVal.IsNil() {
				fieldVal.Set(reflect.New(fieldVal.Type().Elem()))
			}
			if fieldVal.Elem().Kind() == reflect.Struct {
				if tagSection != "" {
					if err := decodeStruct(ini, tagSection, fieldVal.Elem()); err != nil {
						return err
					}
				} else {
					if err := decodeStruct(ini, key, fieldVal.Elem()); err != nil {
						return err
					}
				}
				continue
			}
		}

		val, ok := ini.get(useSection, key)
		if !ok {
			continue
		}

		if err := setField(fieldVal, val); err != nil {
			return fmt.Errorf("decode: field %s: %w", field.Name, err)
		}
	}
	return nil
}

func parseIniTag(tag string) (tagSection, key, opts string) {
	// Format: "section.key" or "key" or "key,optional"
	parts := strings.SplitN(tag, ".", 2)
	if len(parts) == 2 {
		tagSection = parts[0]
		key = parts[1]
	} else {
		key = parts[0]
	}

	// Check for opts after comma
	if idx := strings.Index(key, ","); idx >= 0 {
		opts = key[idx+1:]
		key = key[:idx]
	}

	return tagSection, key, opts
}

func setField(rv reflect.Value, val string) error {
	switch rv.Kind() {
	case reflect.String:
		rv.SetString(val)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(val, 10, rv.Type().Bits())
		if err != nil {
			return err
		}
		rv.SetInt(n)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(val, 10, rv.Type().Bits())
		if err != nil {
			return err
		}
		rv.SetUint(n)

	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(val, rv.Type().Bits())
		if err != nil {
			return err
		}
		rv.SetFloat(n)

	case reflect.Bool:
		b, err := parseBool(val)
		if err != nil {
			return err
		}
		rv.SetBool(b)

	case reflect.Slice:
		return setSliceField(rv, val)

	default:
		return fmt.Errorf("unsupported type: %s", rv.Kind())
	}
	return nil
}

func setSliceField(rv reflect.Value, val string) error {
	// Split by comma for slice fields
	items := strings.Split(val, ",")
	slice := reflect.MakeSlice(rv.Type(), len(items), len(items))
	for i, item := range items {
		item = strings.TrimSpace(item)
		elem := slice.Index(i)
		if err := setField(elem, item); err != nil {
			return err
		}
	}
	rv.Set(slice)
	return nil
}

func parseBool(val string) (bool, error) {
	switch strings.ToLower(val) {
	case "true", "yes", "on", "1":
		return true, nil
	case "false", "no", "off", "0":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean: %q", val)
	}
}
