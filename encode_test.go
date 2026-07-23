package iniparse

import (
	"strings"
	"testing"
)

func TestMarshalSimpleStruct(t *testing.T) {
	type User struct {
		Name string `ini:"name"`
		Age  int    `ini:"age"`
	}

	u := User{Name: "Alice", Age: 30}
	out, err := Marshal(&u)
	if err != nil {
		t.Fatal(err)
	}

	got := string(out)
	if !strings.Contains(got, "name = Alice") {
		t.Errorf("expected name = Alice, got:\n%s", got)
	}
	if !strings.Contains(got, "age = 30") {
		t.Errorf("expected age = 30, got:\n%s", got)
	}
}

func TestMarshalSection(t *testing.T) {
	type Database struct {
		Host string `ini:"host"`
		Port int    `ini:"port"`
	}

	type Config struct {
		DB Database `ini:"database"`
	}

	cfg := Config{DB: Database{Host: "localhost", Port: 5432}}
	out, err := Marshal(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	got := string(out)
	if !strings.Contains(got, "[database]") {
		t.Errorf("expected [database] section, got:\n%s", got)
	}
	if !strings.Contains(got, "host = localhost") {
		t.Errorf("expected host = localhost, got:\n%s", got)
	}
	if !strings.Contains(got, "port = 5432") {
		t.Errorf("expected port = 5432, got:\n%s", got)
	}
}

func TestMarshalMultipleSections(t *testing.T) {
	type Database struct {
		Host string `ini:"host"`
	}

	type Server struct {
		Host string `ini:"host"`
	}

	type Config struct {
		DB     Database `ini:"database"`
		Server Server   `ini:"server"`
	}

	cfg := Config{
		DB:     Database{Host: "db.example.com"},
		Server: Server{Host: "api.example.com"},
	}
	out, err := Marshal(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	got := string(out)
	if !strings.Contains(got, "[database]") {
		t.Errorf("expected [database] section, got:\n%s", got)
	}
	if !strings.Contains(got, "[server]") {
		t.Errorf("expected [server] section, got:\n%s", got)
	}
}

func TestMarshalBoolean(t *testing.T) {
	type Settings struct {
		Debug   bool `ini:"debug"`
		Verbose bool `ini:"verbose"`
	}

	s := Settings{Debug: true, Verbose: false}
	out, err := Marshal(&s)
	if err != nil {
		t.Fatal(err)
	}

	got := string(out)
	if !strings.Contains(got, "debug = true") {
		t.Errorf("expected debug = true, got:\n%s", got)
	}
	if !strings.Contains(got, "verbose = false") {
		t.Errorf("expected verbose = false, got:\n%s", got)
	}
}

func TestMarshalFloat(t *testing.T) {
	type Metrics struct {
		Rate float64 `ini:"rate"`
	}

	m := Metrics{Rate: 3.14}
	out, err := Marshal(&m)
	if err != nil {
		t.Fatal(err)
	}

	got := string(out)
	if !strings.Contains(got, "rate = 3.14") {
		t.Errorf("expected rate = 3.14, got:\n%s", got)
	}
}

func TestMarshalNegativeInt(t *testing.T) {
	type Offset struct {
		Value int `ini:"value"`
	}

	o := Offset{Value: -42}
	out, err := Marshal(&o)
	if err != nil {
		t.Fatal(err)
	}

	got := string(out)
	if !strings.Contains(got, "value = -42") {
		t.Errorf("expected value = -42, got:\n%s", got)
	}
}

func TestMarshalUint(t *testing.T) {
	type Config struct {
		Port uint16 `ini:"port"`
	}

	cfg := Config{Port: 8080}
	out, err := Marshal(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	got := string(out)
	if !strings.Contains(got, "port = 8080") {
		t.Errorf("expected port = 8080, got:\n%s", got)
	}
}

func TestMarshalSlice(t *testing.T) {
	type Config struct {
		Tags []string `ini:"tags"`
	}

	cfg := Config{Tags: []string{"go", "ini", "parser"}}
	out, err := Marshal(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	got := string(out)
	if !strings.Contains(got, "tags = go,ini,parser") {
		t.Errorf("expected tags = go,ini,parser, got:\n%s", got)
	}
}

func TestMarshalSkipsIgnoredFields(t *testing.T) {
	type Config struct {
		Name    string `ini:"name"`
		Ignored string `ini:"-"`
	}

	cfg := Config{Name: "Alice", Ignored: "secret"}
	out, err := Marshal(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	got := string(out)
	if strings.Contains(got, "secret") {
		t.Errorf("expected ignored field to be skipped, got:\n%s", got)
	}
	if !strings.Contains(got, "name = Alice") {
		t.Errorf("expected name = Alice, got:\n%s", got)
	}
}

func TestMarshalTagWithSection(t *testing.T) {
	type Config struct {
		Host string `ini:"server.host"`
		Port int    `ini:"server.port"`
	}

	cfg := Config{Host: "0.0.0.0", Port: 8080}
	out, err := Marshal(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	got := string(out)
	if !strings.Contains(got, "[server]") {
		t.Errorf("expected [server] section, got:\n%s", got)
	}
	if !strings.Contains(got, "host = 0.0.0.0") {
		t.Errorf("expected host = 0.0.0.0, got:\n%s", got)
	}
	if !strings.Contains(got, "port = 8080") {
		t.Errorf("expected port = 8080, got:\n%s", got)
	}
}

func TestMarshalPointerToStruct(t *testing.T) {
	type Inner struct {
		Value string `ini:"value"`
	}

	type Outer struct {
		Inner *Inner `ini:"section"`
	}

	outer := Outer{Inner: &Inner{Value: "hello"}}
	out, err := Marshal(&outer)
	if err != nil {
		t.Fatal(err)
	}

	got := string(out)
	if !strings.Contains(got, "[section]") {
		t.Errorf("expected [section], got:\n%s", got)
	}
	if !strings.Contains(got, "value = hello") {
		t.Errorf("expected value = hello, got:\n%s", got)
	}
}

func TestMarshalNilPointerSkipped(t *testing.T) {
	type Inner struct {
		Value string `ini:"value"`
	}

	type Outer struct {
		Inner *Inner `ini:"section"`
	}

	outer := Outer{Inner: nil}
	out, err := Marshal(&outer)
	if err != nil {
		t.Fatal(err)
	}

	got := string(out)
	if strings.Contains(got, "[section]") {
		t.Errorf("expected section to be skipped for nil pointer, got:\n%s", got)
	}
}

func TestMarshalEmptyStruct(t *testing.T) {
	type Empty struct{}
	e := Empty{}
	out, err := Marshal(&e)
	if err != nil {
		t.Fatal(err)
	}

	if len(out) != 0 {
		t.Errorf("expected empty output, got:\n%s", string(out))
	}
}

func TestMarshalNonStructPointer(t *testing.T) {
	s := "hello"
	_, err := Marshal(&s)
	if err == nil {
		t.Fatal("expected error for non-struct pointer")
	}
}

func TestMarshalNilPointerArg(t *testing.T) {
	_, err := Marshal(nil)
	if err == nil {
		t.Fatal("expected error for nil pointer")
	}
}

func TestMarshalUint64(t *testing.T) {
	type Config struct {
		Big uint64 `ini:"big"`
	}

	cfg := Config{Big: 18446744073709551615}
	out, err := Marshal(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	got := string(out)
	if !strings.Contains(got, "big = 18446744073709551615") {
		t.Errorf("expected big = 18446744073709551615, got:\n%s", got)
	}
}

func TestMarshalFloat32(t *testing.T) {
	type Config struct {
		Rate float32 `ini:"rate"`
	}

	cfg := Config{Rate: 1.5}
	out, err := Marshal(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	got := string(out)
	if !strings.Contains(got, "rate = 1.5") {
		t.Errorf("expected rate = 1.5, got:\n%s", got)
	}
}

func TestMarshalNilSlice(t *testing.T) {
	type Config struct {
		Tags []string `ini:"tags"`
	}

	cfg := Config{Tags: nil}
	out, err := Marshal(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	got := string(out)
	if !strings.Contains(got, "tags = ") {
		t.Errorf("expected tags = (empty), got:\n%s", got)
	}
}

func TestMarshalInt8(t *testing.T) {
	type Config struct {
		Val int8 `ini:"val"`
	}

	cfg := Config{Val: 127}
	out, err := Marshal(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	got := string(out)
	if !strings.Contains(got, "val = 127") {
		t.Errorf("expected val = 127, got:\n%s", got)
	}
}

func TestMarshalDefaultSectionBeforeNamed(t *testing.T) {
	type Config struct {
		Name string `ini:"name"`
		DB   struct {
			Host string `ini:"host"`
		} `ini:"database"`
	}

	cfg := Config{Name: "app", DB: struct {
		Host string `ini:"host"`
	}{Host: "localhost"}}
	out, err := Marshal(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	got := string(out)
	idxName := strings.Index(got, "name = app")
	idxSection := strings.Index(got, "[database]")
	if idxName < 0 || idxSection < 0 {
		t.Fatalf("expected name and [database] in output:\n%s", got)
	}
	if idxName > idxSection {
		t.Errorf("default section keys should appear before named sections:\n%s", got)
	}
}
