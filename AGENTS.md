# AGENTS.md

## Project

- Module: `github.com/Pshimaf-Git/iniparse` (Go 1.24.4)
- Package: `iniparse` (single flat package, no sub-packages)
- Purpose: Parse INI files into Go structs via `ini:"tag"` struct tags

## Build & Verify

```bash
go build ./...
go vet ./...
go test ./...
```

Or use Taskfile:

```bash
task build        # go build
task test         # go test
task test:race    # go test -race
task test:cover   # coverage report
task lint         # golangci-lint
task vet          # go vet
task fmt          # gofmt -s -w .
task check        # vet + lint + test
```

Config files:
- `Taskfile.yaml` — task runner commands
- `.golangci.yml` — linter config (errcheck, govet, staticcheck, gocritic, misspell, etc.)

## Architecture

| File | Role |
|------|------|
| `token.go` | `TokenKind` enum and `Token` struct |
| `lexer.go` | Context-aware tokenizer — tracks `expectValue` after `=` |
| `parser.go` | Parses tokens into `*IniFile` AST |
| `ast.go` | `IniFile` and `Section` types with `Get`/`Set`/`Sections`/`Keys` |
| `decode.go` | Reflection-based struct mapping via `ini` tags |
| `encode.go` | Struct-to-INI encoding |
| `serialize.go` | Serializes `IniFile` to bytes; quotes values containing `;`, `#`, `"`, `\n`, `\r` |
| `iniparse.go` | Public API: `Unmarshal`, `Marshal`, `FromReader`, `ToWriter` |

## Public API

```go
func Unmarshal(input []byte, v any) error
func Marshal(v any) ([]byte, error)
func FromReader(r io.Reader, v any, limit ...int64) error  // default 10MB limit
func ToWriter(w io.Writer, v any) error
```

`FromReader` accepts an optional limit parameter (bytes). Default: 10MB.

## Struct Tag Format

```go
type Config struct {
    Name    string   `ini:"name"`           // default section, key "name"
    Port    int      `ini:"server.port"`    // section "server", key "port"
    Tags    []string `ini:"tags"`           // comma-separated slice
    Debug   bool     `ini:"debug"`          // true/false/yes/no/on/off
    Ignored string   `ini:"-"`              // skipped
}
```

- Tags with `.` prefix: `section.key`
- Tags with `-`: field is skipped
- Slices: values split by comma
- Pointer-to-struct fields: auto-initialized, recurse into section
- Missing keys: field left at zero value (no error)

## INI Format Support

- Sections: `[name]`
- Key-value: `key = value` (whitespace around `=` optional)
- Comments: `;` or `#` at line start
- Quoted strings: `"hello world"`
- Types: booleans (`true`/`false`/`yes`/`no`/`on`/`off`), integers, floats, strings
- Default section: keys before any `[section]` header

## Gotchas

- Lexer is context-aware: after `=`, next token is a value (not a key)
- `isFloat` requires digits on both sides of the dot (e.g. `3.14`, not `.5` or `5.`)
- `isInteger` accepts optional leading `-` or `+`
- Unterminated quoted strings return an error
- Unterminated section headers return an error
- `FromReader` has a 10MB default size limit (configurable via optional parameter)
- `Marshal` quotes values containing `;`, `#`, `"`, `\n`, `\r` for roundtrip safety
