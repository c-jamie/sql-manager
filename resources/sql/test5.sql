/*
  [sqlmbegin]
  [script]
    - description: "updates a, b and C in d"
  [dev]
    - table1: "A.C.Table1"
    - ref1: sqlmref("proj1-test3-sql")
  [sqlmend]
*/
select * from '{{.table1}}'
union all 
{{.ref1}}