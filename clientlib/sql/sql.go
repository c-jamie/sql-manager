package sql

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"text/template"

	"github.com/c-jamie/sql-manager/clientlib/log"
	sqlMig "github.com/c-jamie/sql-manager/clientlib/migration"
	"github.com/c-jamie/sql-manager/clientlib/utils"
)

// SQLStack represents the domain of a parsed SQL script, including any nested SQL referenced
type SQLStack struct {
	SQL        string
	Directives *SQLDirectives
	Nested     bool
	Level      int
	SQuery     *NestedSQLQuery
}

type getter func(table string) (string, error)


// SQLMngr represents the domain of a SQL scripts combined with any migrations for that table
type SQLMngr struct {
	Raw                   string
	Parsed                string
	Directives            *SQLDirectives
	DirectiveKeyOverrides map[string]string
	Err                   error
	Env                   string
	Migrations            map[string][]*sqlMig.SQLMigrationStrategy
	Stack                 []*SQLStack
	StackDepth            int
	Root                  map[string]string
	Getter                getter
}

// New Returns a new SQLMngr
func New(script string, env string, migrations map[string][]*sqlMig.SQLMigrationStrategy, get getter) SQLMngr {
	var sql SQLMngr
	sql.Raw = script
	sql.Env = env
	sql.StackDepth = -1
	sql.Migrations = migrations
	sql.Getter = get
	return sql
}

// stackPush adds any nested SQL to the stack
func (sql *SQLMngr) stackPush(dir *SQLDirectives, sqlScript string, level int, query *NestedSQLQuery) {
	sql.Stack = append(sql.Stack, &SQLStack{SQL: sqlScript, Directives: dir, SQuery: query, Level: level})
	sql.StackDepth += 1
}


// stackPop returns the last nested SQL in the stack
func (sql *SQLMngr) stackPop() *SQLStack {
	sql.StackDepth -= 1
	return sql.Stack[sql.StackDepth+1]
}


// parseRawDirectives grabs the SQLDirectives for a given SQL
func (sql *SQLMngr) parseRawDirectives() {
	sql.Directives, sql.Err = Parse(sql.Raw)
}


// getScript returns the referenced SQL Script
func (sql *SQLMngr) getScript(file string, key string, name string) string {
	if file != "" {
		return sql.fileLoader(file)
	} else {
		return sql.serverLoader(name)
	}
}


// fileLoader loads a file from disk
func (sql *SQLMngr) fileLoader(location string) string {
	sqlScript, err := utils.ReadFile(location)
	if err != nil {
		return ""
	} else {
		return string(sqlScript)
	}
}

// serverLoader loads a file from the SQL manager server
func (sql *SQLMngr) serverLoader(name string) string {
	sqlScript, err := sql.Getter(name)
	if err != nil {
		return ""
	} else {
		return string(sqlScript)
	}
}


// parseDirectives gets the nested SQL and any keywords from the SQLDirectives based on the Env
func (sql *SQLMngr) parseDirectives(dir *SQLDirectives) ([]*NestedSQLQuery, map[string]string) {
	if sql.Env == "int" {
		return dir.Int.NestedSQL, dir.Int.Keywords
	} else if sql.Env == "prod" {
		return dir.Prod.NestedSQL, dir.Prod.Keywords
	} else if sql.Env == "dev" {
		return dir.Dev.NestedSQL, dir.Dev.Keywords
	} else if sql.Env == "local" {
		return dir.Local.NestedSQL, dir.Local.Keywords
	} else {
		panic("bad env declared")
	}
}


// compileFragments parses all nested SQL in the stack
func (sql *SQLMngr) compileFragments(sqlScript string, nestLevel int, dirUp bool, query *NestedSQLQuery) {
	if sql.StackDepth > 100 {
		panic("recursion error compiling SQL, depth exceeded")
	}
	var nSQL []*NestedSQLQuery
	var sqlIn string
	var err error
	var dir *SQLDirectives
	dir, err = Parse(sqlScript)
	if err != nil {
		panic(err)
	}
	sql.stackPush(dir, sqlScript, nestLevel, query)
	nSQL, _ = sql.parseDirectives(dir)
	for _, j := range nSQL {
		log.Debug("nest level ", nestLevel)
		sqlIn = sql.getScript(j.File, j.Key, j.Name)
		sql.compileFragments(sqlIn, nestLevel+1, false, j)
	}
}


// mergeFragments combines the SQL stack together
func (sql *SQLMngr) mergeFragments() {
	context := make([]map[string]string, len(sql.Stack))
	for i := 0; i < len(sql.Stack); i++ {
		log.Debug("merge ", i)
		stk := sql.stackPop()
		nestedSQL, keywords := sql.parseDirectives(stk.Directives)
		t := template.Must(template.New("sql").Parse(stk.SQL))
		var out bytes.Buffer
		insert := make(map[string]string)
		if stk.Level != 0 {
			if nestedSQL == nil {
				t.Execute(&out, keywords)
				insert[stk.SQuery.Key] = out.String()
				context[stk.Level] = insert
			} else {
				t := template.Must(template.New("sql").Parse(stk.SQL))
				log.Debug("stk.level", stk.Level)
				for k, v := range context[stk.Level+1] {
					keywords[k] = v
				}
				t.Execute(&out, keywords)
				insert[stk.SQuery.Key] = out.String()
				context[stk.Level] = insert
			}
		}
	}
	if len(sql.Stack) > 1 {
		sql.Root = context[1]
	}
}

// finalise generates the final SQL script
func (sql *SQLMngr) finalise() {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		fmt.Println("error", err)
	}
	var keywords map[string]string
	if len(sql.Stack) > 1 {
		keywords = sql.Root
	} else {
		keywords = make(map[string]string)
	}
	log.Debug(sql.Root)
	t := template.Must(template.New("sql").Parse(sql.Raw))
	var out bytes.Buffer

	_, topKey := sql.parseDirectives(sql.Directives)

	for k, v := range topKey {
		keywords[k] = v
	}
	for k, v := range sql.DirectiveKeyOverrides {
		keywords[k] = v
	}
	if sql.Migrations[sql.Env] != nil {
		for _, v := range sql.Migrations[sql.Env] {
			for _, j := range v.MigrationsUp {
				table := reg.ReplaceAllString(j.SourceTable, "_") + "_" + strconv.Itoa(j.FileOrder)
				keywords[table] = "true"
			}
		}
	}
	t.Execute(&out, keywords)
	sql.Parsed = out.String()
}


// Compile returns the parsed SQL
func (sql *SQLMngr) Compile() string {
	sql.parseRawDirectives()
	sql.compileFragments(sql.Raw, 0, false, nil)
	sql.mergeFragments()
	sql.finalise()
	return sql.Parsed
}
