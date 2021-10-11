package sql

import (
	"testing"
)

func TestParseAllEnvs(t *testing.T) {
	testcases := []struct {
		in     string
		test   bool
		expect map[string]string
	}{{
		in: `
		/*
		[sqlmbegin]
		[script] 
			- description: "updates a,b and c in d" 
		[dev] 
			- out_table: "A.B.Table1"  
			- out_table_2: "A.B.Table2"  
		[prod] 
			- out_table: "A.B.Table3"  
			- out_table_2: "A.B.Table4"  
		[int] 
			- out_table: "A.B.Table5"  
			- out_table_2: "A.B.Table6"  
		[local] 
			- out_table: "A.B.Table7"  
			- out_table_2: "A.B.Table8"  
		[sqlmend]
		*/
		`,
		test:   true,
		expect: map[string]string{"description": "updates a,b and c in d"},
	}}

	for _, tcase := range testcases {
		obj, _ := Parse(tcase.in)
		if obj.Description != tcase.expect["description"] {
			t.Errorf(" error name %s %s", obj.Description, tcase.expect["description"])
		}
	}
}

func TestParseKeywords(t *testing.T) {
	testcases := []struct {
		in     string
		test   bool
		expect map[string]string
	}{{
		in: `
		/*
		[sqlmbegin]
		[script] 
			- description: "updates a,b and c in d" 
		[dev] 
			- out_table_1: "A.B.Table"  
			- out_table_2: sqlmref("test1.sql")
			- out_table_3: sqlmfile("test2.sql")
			- out_table_4: "a.c.d.table" 
		[sqlmend]
		*/
		`,
		test:   true,
		expect: map[string]string{"out_table_1": "A.B.Table", "out_table_2": "test1.sql", "out_table_3": "test2.sql", "out_table_4": "a.c.d.table"},
	}}

	for _, tcase := range testcases {
		obj, _ := Parse(tcase.in)
		if obj.Dev.Keywords["out_table_1"] != tcase.expect["out_table_1"] {
			t.Errorf(" error name %s %s", obj.Dev.Keywords["out_table_1"], tcase.expect["out_table_1"])
		}
		if obj.Dev.Keywords["out_table_4"] != tcase.expect["out_table_4"] {
			t.Errorf(" error name %s %s", obj.Dev.Keywords, tcase.expect)
		}
		if obj.Dev.NestedSQL[1].File != tcase.expect["out_table_3"] {
			t.Errorf(" error name %s %s", obj.Description, tcase.expect["description"])
		}
	}
}
