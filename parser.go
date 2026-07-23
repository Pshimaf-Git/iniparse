package iniparse

import "fmt"

type parser struct{}

func newParser() *parser {
	return &parser{}
}

// parse parses an INI string into an iniFile.
func (p *parser) parse(input string) (*iniFile, error) {
	lex := newLexer(input)
	tokens, err := lex.tokenize()
	if err != nil {
		return nil, err
	}

	ini := newIniFile()
	currentSection := ""

	i := 0
LOOP:
	for i < len(tokens) {
		tok := tokens[i]

		switch tok.kind {
		case tokenSection:
			currentSection = tok.literal
			// ensure section exists even if empty
			ini.getOrCreateSection(currentSection)

		case tokenKey:
			key := tok.literal
			// expect assign next
			i++
			if i >= len(tokens) || tokens[i].kind != tokenAssign {
				return nil, fmt.Errorf("%d:%d: expected '=' after key %q", tok.line, tok.col, key)
			}
			// expect value after assign
			i++
			if i >= len(tokens) {
				return nil, fmt.Errorf("%d:%d: expected value after '=' for key %q", tok.line, tok.col, key)
			}
			valTok := tokens[i]
			var value string
			switch valTok.kind {
			case tokenString, tokenValue, tokenInteger, tokenFloat, tokenBoolean:
				value = valTok.literal
			default:
				return nil, fmt.Errorf("%d:%d: unexpected token after '=': %s", valTok.line, valTok.col, valTok.kind)
			}
			ini.set(currentSection, key, value)

		case tokenComment, tokenAssign:
			// skip comments and stray assigns

		case tokenEOF:
			break LOOP

		default:
			return nil, fmt.Errorf("%d:%d: unexpected token: %s", tok.line, tok.col, tok.kind)
		}

		i++
	}

	return ini, nil
}

// unmarshal parses INI data and maps it to the provided struct.
func (p *parser) unmarshal(input string, v any) error {
	ini, err := p.parse(input)
	if err != nil {
		return err
	}
	return decode(ini, v)
}
