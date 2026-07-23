package iniparse

import (
	"bytes"
	"strings"
)

func quoteValue(val string) string {
	if strings.ContainsAny(val, ";#\"\n\r") {
		val = strings.ReplaceAll(val, `\`, `\\`)
		val = strings.ReplaceAll(val, `"`, `\"`)
		val = strings.ReplaceAll(val, "\n", `\n`)
		val = strings.ReplaceAll(val, "\r", `\r`)
		return `"` + val + `"`
	}
	return val
}

func serialize(ini *iniFile) []byte {
	var buf bytes.Buffer

	// default section keys (no header)
	for _, key := range ini.default_.order {
		val := ini.default_.values[key]
		buf.WriteString(key)
		buf.WriteString(" = ")
		buf.WriteString(quoteValue(val))
		buf.WriteByte('\n')
	}

	// named sections in order
	for _, secName := range ini.order {
		s := ini.sections[secName]
		buf.WriteByte('[')
		buf.WriteString(secName)
		buf.WriteString("]\n")
		for _, key := range s.order {
			val := s.values[key]
			buf.WriteString(key)
			buf.WriteString(" = ")
			buf.WriteString(quoteValue(val))
			buf.WriteByte('\n')
		}
	}

	return buf.Bytes()
}
