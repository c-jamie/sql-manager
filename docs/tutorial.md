# SQL Manager tutorial

This tutorial covers some of the featurs of the platform namely:

1. Grabbing SQL scripts from the server
2. Working with environments
3. Working with migrations

## Setup

clone `sql-manager` and `auth-manager`

run the servers

** Check you've got a valid git repo + token in the .env file **

```
~/auth-manager make ENV=dev run-server
```

```
~/sql-manager make ENV=dev run-server
```

create two databases, "dev" and "prod"

```
~/sql-manager make ENV=dev test-db-dev
```

```
~/sql-manager make ENV=dev test-db-prod
```

## Create an account and login

Follow the steps here

```
~/code/sql-manager$ make ENV=dev run-client args='register'
```

Now login

```
~/code/sql-manager$ make ENV=dev run-client args='login'
```

## Migrations

The server is now running and set up.

The dev.env file has a git repo which has been loaded into memory.

We can now register these files with the server.

Register 2 migrations with the dev env.

```
~/code/sql-manager$ make ENV=dev run-client args='migration add -e dev -t abc tutorial/migrations/dev/1_m.sql'
```

```
~/code/sql-manager$ make ENV=dev run-client args='migration add -e dev -t abc tutorial/migrations/dev/2_m.sql'
```

Now, tell the platform you want to set migration "2" as the active migration.

If you want to roll back migrations, run `migration set dev` again and you can set the active migration to migration 1.

Then the next time you run migrations, it will rollback to the first migration.

```
(base) cillairne@cillairne:~/code/sql-manager$ make ENV=dev run-client args='migration set dev'
```

Register the fist migration with the prod env.

```
~/code/sql-manager$ make ENV=dev run-client args='migration add -e prod -t abc tutorial/migrations/dev/1_m.sql'
```

## Scripts

Now we're going to register our SQL script with the platform.

```
~/code/sql-manager$ make ENV=dev run-client args='script register tutorial/example.sql'
```

NB: the name of the SQL file is a slugified version of the script location in git, in this instance `tutorial-example-sql`

## Running our Migrations

Run out dev enviroment migrations.

```
~/code/sql-manager$ make ENV=dev run-client args='-v 1 migration run -e dev -c postgres://test:welcome@0.0.0.0:5436/test?sslmode=disable -d postgres abc'
```

Run out prod enviornment migrations.

```
~/code/sql-manager$ make ENV=dev run-client args='-v 1 migration run -e prod -c postgres://test:welcome@0.0.0.0:5437/test?sslmode=disable -d postgres abc'
```

## Executing SQL

The templated SQL file looks like this.

The platform will parse the tags in the comment, and will provide them as variables to reference in SQL file.

Other cool tags like `- ref1: sqlmref('script-slug')` allow you to reference other scripts which have been registered with the system.

```
/*
  [sqlmbegin]
  [script]
    - description: "updates a, b and C in d"
  [dev]
    - table1: "project1"
  [prod]
    - table1: "project1"
  [sqlmend]
*/

{{if .abc_2}}
-- if the second migration is active, we have access to this column
select name, name_type from {{.table1}}
{{else}}
select name,  null as name_type from {{.table1}}
{{end}}

```

Now, lets run our SQL against the dev env.

```
~/code/sql-manager$ make ENV=dev run-client args='script gc -e prod tutorial-example-sql'
```

You'll notice that our templated SQL has access to all golang macros, and also the migration scripts which have been run against the database.

Because migration script two has been run, we have access to column name_type.

```
/*
  [sqlmbegin]
  [script]
    - description: "updates a, b and C in d"
  [dev]
    - table1: "project1"
  [prod]
    - table1: "project1"
  [sqlmend]
*/


-- if the second migration is active, we have access to this column
select name, name_type from project1

```

Now, lets look at the prod enviornment.

```
~/code/sql-manager$ make ENV=dev run-client args='script gc -e dev tutorial-example-sql'
```

The second migration hasn't been run in the prod enviornment, therefore the SQL script will use a null value instead for this column.

```
/*
  [sqlmbegin]
  [script]
    - description: "updates a, b and C in d"
  [dev]
    - table1: "project1"
  [prod]
    - table1: "project1"
  [sqlmend]
*/


select name,  null as name_type from project1

```