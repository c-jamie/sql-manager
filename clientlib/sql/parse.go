package sql

import (
	"errors"
	"fmt"
)

// SQLQuery represents the domain for a SQL query when referenced by another query
type NestedSQLQuery struct {
	Key  string
	File string
	Name string
}

// SQLDirectives represents the domain for a fully parsed SQL script
type SQLDirectives struct {
	Name        string
	Description string
	Prod        struct {
		Keywords  map[string]string
		NestedSQL []*NestedSQLQuery
	}
	Dev struct {
		Keywords  map[string]string
		NestedSQL []*NestedSQLQuery
	}
	Local struct {
		Keywords  map[string]string
		NestedSQL []*NestedSQLQuery
	}
	Int struct {
		Keywords  map[string]string
		NestedSQL []*NestedSQLQuery
	}
}


// Lexer represents the domain for our lexer
type Lexer struct {
	tkn       *Tokenizer
	cur_id    int
	cur_byte  []byte
	next_id   int
	next_byte []byte
}

// NewLexer returns a new lexer
func NewLexer(tkn *Tokenizer) *Lexer {
	cur_id, cur_byte := tkn.Scan()
	next_id, next_byte := tkn.Scan()
	return &Lexer{tkn: tkn, cur_id: cur_id, cur_byte: cur_byte, next_id: next_id, next_byte: next_byte}
}


// Next returns the proceeding ID and byte in the stream
func (lex *Lexer) Next() (int, []byte) {
	next_id, next_byte := lex.tkn.Scan()
	lex.cur_byte = lex.next_byte
	lex.cur_id = lex.next_id
	lex.next_id = next_id
	lex.next_byte = next_byte
	return lex.cur_id, lex.cur_byte
}

// Peek allows you to see the next ID and byte wihout moving forward in the stream
func (lex *Lexer) Peek() (int, []byte) {
	return lex.next_id, lex.next_byte
}


// Parse returns SQLDirectives from a SQL script
func Parse(sql string) (*SQLDirectives, error) {
	var sqlDir SQLDirectives
	sqlDir.Dev.Keywords = make(map[string]string)
	tokenizer := NewStringTokenizer(sql)
	lexer := NewLexer(tokenizer)
	id := lexer.cur_id

	if id != int('[') {
		return nil, errors.New("invalid sql script no [")
	}

	id, _ = lexer.Next()
	if id != PD_BEGIN {
		return nil, errors.New("invalid sql script, no PD_BEGIN")
	}

	id, _ = lexer.Next()
	if id != int(']') {
		return nil, errors.New("invalid sql script  no ]")
	}

	if !sqlDir.parseScript(lexer) {
		return nil, errors.New("unable to parse sql script")
	}
	return &sqlDir, nil
}


// parseScript returns the SQLDirectives from a script
func (pds *SQLDirectives) parseScript(lex *Lexer) bool {
	id, _ := lex.Next()
	if id != 91 {
		return false
	}

	id, _ = lex.Next()
	if id != SCRIPT {
		return false
	}

	id, _ = lex.Next()
	if id != 93 {
		return false
	}

	if !pds.parseDescription(lex) {
		fmt.Println("parse description failed")
	}

	for lex.cur_id != PD_END {
		id, _ := lex.Next()
		if id != 91 {
			return false
		}
		if !pds.parseEnv(lex, DEV, "dev") {
			if !pds.parseEnv(lex, PROD, "prod") {
				if !pds.parseEnv(lex, INT, "int") {
					if !pds.parseEnv(lex, LOCAL, "local") {
						id, _ := lex.Next()
						if id == PD_END {
							return true
						}
						continue
					}
				}
			}
		}
	}

	return true
}


// parseEnv retruns the key values from an Env
func (pds *SQLDirectives) parseEnv(lex *Lexer, envID int, env string) bool {
	id, _ := lex.Peek()
	if id != envID {
		return false
	}
	id, _ = lex.Next()
	id, _ = lex.Next()
	if id != 93 {
		return false
	}

	for {
		id, _ = lex.Peek()
		if id != 45 {
			break
		} else {
			pds.parseKeyValues(lex, env)
			continue
		}
	}

	return true
}


// parseKeyValues fills out parsed key values from a SQL Script
func (pds *SQLDirectives) parseKeyValues(lex *Lexer, env string) bool {
	key := ""
	id, buf := lex.Next()
	if id != 45 {
		return false
	}
	id, buf = lex.Next()
	if id != ID {
		return false
	} else {
		key = string(buf)
	}
	id, buf = lex.Next()
	if id != 58 {
		return false
	}
	id, buf = lex.Peek()
	if id == PD_REF {
		id, buf = lex.Next()
		id, buf = lex.Next()

		if id != 40 {
			return false
		}

		id, buf = lex.Next()
		if id != STRING {
			return false
		} else {
			if env == "dev" {
				pds.Dev.NestedSQL = append(pds.Dev.NestedSQL, &NestedSQLQuery{key, "", string(buf)})
			} else if env == "prod" {
				pds.Prod.NestedSQL = append(pds.Prod.NestedSQL, &NestedSQLQuery{key, "", string(buf)})
			} else if env == "local" {
				pds.Local.NestedSQL = append(pds.Local.NestedSQL, &NestedSQLQuery{key, "", string(buf)})
			} else if env == "int" {
				pds.Int.NestedSQL = append(pds.Int.NestedSQL, &NestedSQLQuery{key, "", string(buf)})
			}
		}

		id, buf = lex.Next()
		return id == 41
	} else if id == PD_FILE {
		id, buf = lex.Next()
		id, buf = lex.Next()

		if id != 40 {
			return false
		}

		id, buf = lex.Next()
		if id != STRING {
			return false
		} else {
			if env == "dev" {
				pds.Dev.NestedSQL = append(pds.Dev.NestedSQL, &NestedSQLQuery{key, string(buf), ""})
			} else if env == "prod" {
				pds.Prod.NestedSQL = append(pds.Prod.NestedSQL, &NestedSQLQuery{key, string(buf), ""})
			} else if env == "int" {
				pds.Int.NestedSQL = append(pds.Int.NestedSQL, &NestedSQLQuery{key, string(buf), ""})
			} else if env == "local" {
				pds.Local.NestedSQL = append(pds.Local.NestedSQL, &NestedSQLQuery{key, string(buf), ""})
			}
		}

		id, buf = lex.Next()
		if id != 41 {
			return false
		}
		return true
	}

	id, _ = lex.Peek()
	if id != STRING {
		return false
	} else {
		id, buf = lex.Next()
		pds.Dev.Keywords[key] = string(buf)
	}
	return true
}


// parseName fills out the Name from a SQL script
func (pds *SQLDirectives) parseName(lex *Lexer) bool {
	id, _ := lex.Next()
	if id != 45 {
		return false
	}
	id, _ = lex.Next()
	if id != NAME {
		return false
	}
	id, _ = lex.Next()
	if id != 58 {
		return false
	} else {
		id, buf := lex.Next()
		out := buf
		for {
			id, buf = lex.Peek()
			if id != STRING {
				break
			} else {
				id, buf = lex.Next()
				out = append(out, " "...)
				out = append(out, buf...)
				continue
			}
		}
		pds.Name = string(out)
		return true
	}
}


// parseDescription fills out the description from a SQL script
func (pds *SQLDirectives) parseDescription(lex *Lexer) bool {
	id, _ := lex.Next()
	if id != 45 {
		return false
	}
	id, _ = lex.Next()
	if id != DESCRIPTION {
		return false
	}
	id, _ = lex.Next()
	if id != 58 {
		return false
	} else {
		id, buf := lex.Next()
		out := buf
		for {
			id, buf = lex.Peek()
			if id != STRING {
				break
			} else {
				id, buf = lex.Next()
				out = append(out, " "...)
				out = append(out, buf...)
				continue
			}
		}
		pds.Description = string(out)
		return true
	}
}
