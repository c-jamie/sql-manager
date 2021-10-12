# SQL Manager

## This project is an experiment in adding an API of sorts over your SQL.

* A CLI to retrieve your SQL from the server by table / env
* Group you SQL by project 

SQL Manager also has some helper features for generating SQL

* Leverage Golang Templates in SQL
* Reference other SQL scripts from a given script

SQL Manager also has a migration manager, but provides some extra features I've not seen, namely: 

* Access metadata on migrations which have been run on specific tables in you SQL scripts
* Full control of what migrations should be run from the CLI - you can add and rollback migrations as needed

These two things allow a data analyst or scientist to almost think of migrations as feature flags which can be switched on and off depending on a given environment.

## What is this comparable to?

dbt - use dbt in prod. It's much better than this project.

## Running the server

Step 1 - clone `auth-manager` 

This is used by the server for permissioning

Step 2 - run the `auth-manager` server 
```
cd ~/auth-manager
make ENV=dev run-server
```

Step 3 - run the server
```
cd ~/sql-manager
make ENV=dev run-server
```


## Tutorial

See `docs/`