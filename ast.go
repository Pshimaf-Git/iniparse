package iniparse

// IniFile represents a parsed INI file.
type iniFile struct {
	sections map[string]*section
	order    []string
	default_ *section
}

// Section represents an INI section.
type section struct {
	name   string
	values map[string]string
	order  []string
}

// newIniFile creates a new empty IniFile.
func newIniFile() *iniFile {
	return &iniFile{
		sections: make(map[string]*section),
		default_: newSection(""),
	}
}

func newSection(name string) *section {
	return &section{
		name:   name,
		values: make(map[string]string),
	}
}

// get returns the value for the given key in the specified section.
// If section is empty, the default (unnamed) section is used.
// Returns ("", false) if not found.
func (f *iniFile) get(sec, key string) (string, bool) {
	s := f.getSection(sec)
	if s == nil {
		return "", false
	}
	val, ok := s.values[key]
	return val, ok
}

// set sets a value in the specified section.
// If section is empty, the default section is used.
func (f *iniFile) set(sec, key, value string) {
	s := f.getOrCreateSection(sec)
	if _, exists := s.values[key]; !exists {
		s.order = append(s.order, key)
	}
	s.values[key] = value
}

// sectionNames returns all section names in order.
func (f *iniFile) sectionNames() []string {
	return append([]string{}, f.order...)
}

// keys returns all keys in a section in order.
func (f *iniFile) keys(sec string) []string {
	s := f.getSection(sec)
	if s == nil {
		return nil
	}
	return append([]string{}, s.order...)
}

// hasSection returns true if the section exists.
func (f *iniFile) hasSection(name string) bool {
	if name == "" {
		return true
	}
	_, ok := f.sections[name]
	return ok
}

func (f *iniFile) getSection(name string) *section {
	if name == "" {
		return f.default_
	}
	return f.sections[name]
}

func (f *iniFile) getOrCreateSection(name string) *section {
	if name == "" {
		return f.default_
	}
	if s, ok := f.sections[name]; ok {
		return s
	}
	s := newSection(name)
	f.sections[name] = s
	f.order = append(f.order, name)
	return s
}
