/*
  [sqlmbegin]
  [script]
    - description: "updates a, b and C in d"
  [dev]
    - table1: "A.C.Table8"
    - ref1: sqlmfile("../../resources/sql/test3.sql")
  [sqlmend]
*/
select top 10 * from {{.table1}} intersection select 1 from {{.ref1}}
