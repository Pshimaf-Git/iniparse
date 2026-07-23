# iniparse

Go library for parsing INI files into Go structs and back, using `ini` struct tags.

## Install

```bash
go get github.com/Pshimaf-Git/iniparse
```

## Usage

### Unmarshal (INI to struct)

```go
type Config struct {
    Name    string `ini:"name"`
    Port    int    `ini:"server.port"`
    Debug   bool   `ini:"debug"`
}

var cfg Config
err := iniparse.Unmarshal([]byte(`
name = myapp
debug = true

[server]
port = 8080
`), &cfg)
// cfg.Name == "myapp", cfg.Port == 8080, cfg.Debug == true
```

### Marshal (struct to INI)

```go
cfg := Config{Name: "myapp", Port: 8080, Debug: true}
data, err := iniparse.Marshal(&cfg)
// data == "name = myapp\ndebug = true\n\n[server]\nport = 8080\n"
```

### FromReader / ToWriter

```go
// Read from file with default 10MB limit
err := iniparse.FromReader(file, &cfg)

// Read with custom limit (5MB)
err := iniparse.FromReader(file, &cfg, 5*1024*1024)

// Read without limit
err := iniparse.FromReader(file, &cfg, -1)

// Write to file
err := iniparse.ToWriter(file, &cfg)
```

## API

```go
func Unmarshal(input []byte, v any) error
func Marshal(v any) ([]byte, error)
func FromReader(r io.Reader, v any, limit ...int64) error
func ToWriter(w io.Writer, v any) error
```

## Struct Tags

| Tag | Description |
|-----|-------------|
| `ini:"key"` | Map to `key` in the default section |
| `ini:"section.key"` | Map to `key` in `[section]` |
| `ini:"-"` | Skip this field |
| `[]string` field | Comma-separated values |

Pointer-to-struct fields are auto-initialized and recurse into their section.

## Supported INI Features

- Sections: `[name]`
- Key-value: `key = value`
- Comments: `;` or `#` at line start
- Quoted strings: `"hello world"`
- Types: booleans (`true`/`false`/`yes`/`no`/`on`/`off`), integers, floats, strings
- Default section: keys before any `[section]` header

## Development

```bash
go build ./...
go test ./...
go vet ./...
```

Or with [Task](https://taskfile.dev/):

```bash
task build
task test
task lint
task check
```

## License

MIT (see [LICENSE](LICENSE))
