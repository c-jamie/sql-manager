/*
  [sqlmbegin]
  [script]
    - description: "updates a, b and C in d"
  [dev]
    - table1: "A.C.Table1"
    - table2: "A.C.Table2"
    - ref1: sqlmfile("../../resources/sql/test2.sql")
  [sqlmend]
*/
{{if .a_b_c_1}}
-- uncommented
{{end}}
select * from '{{.table1}}' union all '{{.table2}}' 
union all 
{{.ref1}}