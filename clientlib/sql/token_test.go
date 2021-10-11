package sql

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTokenLiteralID(t *testing.T) {
	testcases := []struct {
		in  string
		id  int
		out string
	}{{
		in:  "`aa`",
		id:  ID,
		out: "aa",
	}, {
		in:  "```a```",
		id:  ID,
		out: "`a`",
	}, {
		in:  "`a``b`",
		id:  ID,
		out: "a`b",
	}, {
		in:  "`a``b`c",
		id:  ID,
		out: "a`b",
	}, {
		in:  "`a``b",
		id:  LEX_ERROR,
		out: "a`b",
	}, {
		in:  "`a``b``",
		id:  LEX_ERROR,
		out: "a`b`",
	}, {
		in:  "``",
		id:  LEX_ERROR,
		out: "",
	}}

	for _, tcase := range testcases {
		tkn := NewStringTokenizer(tcase.in)
		id, out := tkn.Scan()
		if tcase.id != id || string(out) != tcase.out {
			t.Errorf("Scan(%s): %d, %s, want %d, %s", tcase.in, id, out, tcase.id, tcase.out)
		}
	}
}

func tokenName(id int) string {
	if id == STRING {
		return "STRING"
	} else if id == LEX_ERROR {
		return "LEX_ERROR"
	}
	return fmt.Sprintf("%d", id)
}

func TestTokenString(t *testing.T) {
	testcases := []struct {
		in   string
		id   int
		want string
	}{{
		in:   "''",
		id:   STRING,
		want: "",
	}, {
		in:   "''''",
		id:   STRING,
		want: "'",
	}, {
		in:   "'hello'",
		id:   STRING,
		want: "hello",
	}, {
		in:   "'\\n'",
		id:   STRING,
		want: "\n",
	}, {
		in:   "'\\nhello\\n'",
		id:   STRING,
		want: "\nhello\n",
	}, {
		in:   "'a''b'",
		id:   STRING,
		want: "a'b",
	}, {
		in:   "'a\\'b'",
		id:   STRING,
		want: "a'b",
	}, {
		in:   "'\\'",
		id:   LEX_ERROR,
		want: "'",
	}, {
		in:   "'",
		id:   LEX_ERROR,
		want: "",
	}, {
		in:   "'hello\\'",
		id:   LEX_ERROR,
		want: "hello'",
	}, {
		in:   "'hello",
		id:   LEX_ERROR,
		want: "hello",
	}, {
		in:   "'hello\\",
		id:   LEX_ERROR,
		want: "hello",
	}}

	for _, tcase := range testcases {
		id, got := NewStringTokenizer(tcase.in).Scan()
		if tcase.id != id || string(got) != tcase.want {
			t.Errorf("Scan(%q) = (%s, %q), want (%s, %q)", tcase.in, tokenName(id), got, tokenName(tcase.id), tcase.want)
		}
	}
}

func TestTokenSQLManager(t *testing.T) {
	testcases := []struct {
		in   string
		test bool
		out  []int
	}{{
		in: `/*
			[sqlmbegin]
			[script] 
				- description: "updates Δ a,b and c in d" 
				- version: "1.1"
			[dev] 
				- out_table: "A.B.Table2"  
				- out_script: sqlmref("test2.sql")
			[sqlmend]
			*/`,
		test: false,
	}, {
		in: `/*
			[sqlmbegin]
			[script] 
				- description: "updates Δ a,b and c in d" 
			[sqlmend]
		*/`,
		test: true,
		out:  []int{int('['), PD_BEGIN, int(']'), int('['), SCRIPT, int(']'), int('-'), DESCRIPTION, int(':'), 7002, int('['), PD_END},
	}, {
		in: `/*
			[sqlmbegin]
			[dev] 
				- out_script: sqlmref("test2.sql")
			[prod] 
				- out_script1: "aa.t" 
			[sqlmend]
			*/`,
		test: true,
		out:  []int{int('['), PD_BEGIN, int(']'), int('['), DEV, int(']'), int('-'), ID, int(':'), PD_REF, int('('), 7002, int(')'), int('['), PROD, int(']'), int('-'), ID, int(':'), 7002, int('['), PD_END},
	}}

	for _, tcase := range testcases {
		tokenizer := NewStringTokenizer(tcase.in)
		var ids []int
		for {
			id, _ := tokenizer.Scan()
			ids = append(ids, id)
			if id == PD_END {
				break
			}
		}
		if tcase.test == true {
			tst := reflect.DeepEqual(ids, tcase.out)
			if tst == false {
				t.Error("fail: SQL     ", tcase.in)
				t.Error("fail: expected", tcase.out)
				t.Error("fail: got     ", ids)
			}
		}
	}
}
