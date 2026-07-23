package iniparse

import (
	"fmt"
	"io"
)

const defaultReaderLimit int64 = 10 * 1024 * 1024 // 10MB

func FromReader(r io.Reader, v any, _limit ...int64) error {
	limit := limit(_limit...)

	var (
		data []byte
		err  error
	)

	if limit <= 0 {
		data, err = io.ReadAll(r)
	} else {
		data, err = io.ReadAll(io.LimitReader(r, limit))
	}

	if err != nil {
		return fmt.Errorf("ini: read error: %w", err)
	}

	return Unmarshal(data, v)
}

func ToWriter(w io.Writer, v any) error {
	data, err := Marshal(v)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("ini: write error: %w", err)
	}

	return nil
}

// Unmarshal parses INI data and maps it to the provided struct.
func Unmarshal(input []byte, v any) error {
	p := newParser()
	return p.unmarshal(string(input), v)
}

func Marshal(v any) ([]byte, error) {
	ini, err := encode(v)
	if err != nil {
		return nil, err
	}
	return serialize(ini), nil
}

func limit(optionalLimit ...int64) int64 {
	if len(optionalLimit) == 0 {
		return defaultReaderLimit
	}
	return optionalLimit[0]
}
