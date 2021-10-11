package api

import (
	"net/http"

	"github.com/c-jamie/sql-manager/serverlib/internal/data"
	"github.com/c-jamie/sql-manager/serverlib/internal/validator"
	"github.com/gin-gonic/gin"
)

func (app *Application) listFilesHandeler(c *gin.Context) {
	qs := c.Request.URL.Query()
	name := app.readString(qs, "project", "")

	v := validator.New()
	v.Check(name != "", "project", "must not be empty")

	if !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}
	sqlScript, err := app.Models.SQLScript.GetAll(name)

	if err != nil {
		app.badRequest(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": sqlScript})
}

func (app *Application) getFilesHandeler(c *gin.Context) {
	qs := c.Request.URL.Query()
	name := app.readString(qs, "name", "")

	v := validator.New()
	v.Check(name != "", "name", "must not be empty")

	if !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}
	sqlScript, err := app.Models.SQLScript.Get(name)

	if err != nil {
		app.badRequest(c, err)
		return
	}

	sql, err := app.GIT.GetFile(sqlScript.FileLocation)
	if err != nil {
		app.badRequest(c, err)
		return
	}

	sqlScript.Script = sql

	c.JSON(http.StatusOK, gin.H{"files": sqlScript})
}

func (app *Application) registerFilesHandeler(c *gin.Context) {
	var input struct {
		File    string `json:"file"`
		Project string `json:"project"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		app.badRequest(c, err)
	}

	v := validator.New()
	v.Check(input.File != "", "file", "must not be empty")
	v.Check(input.Project != "", "project", "must not be empty")

	if !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}

	sqlScript := data.SQLScript{Project: input.Project, FileLocation: input.File}
	err := app.Models.SQLScript.Register(&sqlScript)

	if err != nil {
		app.badRequest(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"files": sqlScript})

}

func (app *Application) getProjectsHandeler(c *gin.Context) {
	qs := c.Request.URL.Query()
	name := app.readString(qs, "name", "")

	v := validator.New()
	v.Check(name != "", "name", "must not be empty")

	if !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}
	project, err := app.Models.Project.Get(name)
	if err != nil {
		app.badRequest(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"projects": project})
}

func (app *Application) addMigrationsHandeler(c *gin.Context) {
	var input struct {
		File  string `json:"file"`
		Env   string `json:"env"`
		Table string `json:"table"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		app.badRequest(c, err)
	}

	v := validator.New()
	v.Check(input.File != "", "file", "must not be empty")
	v.Check(input.Env != "", "env", "must not be empty")
	v.Check(input.Table != "", "table", "must not be empty")

	if !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}

	mig := data.SQLMigration{File: input.File, Env: input.Env, SourceTable: input.Table}
	err := app.Models.SQLMigration.Add(&mig)

	if err != nil {
		app.badRequest(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"migrations": mig})
}

func (app *Application) getMigrationTablesHandeler(c *gin.Context) {
	qs := c.Request.URL.Query()
	env := app.readString(qs, "env", "")

	v := validator.New()
	v.Check(env != "", "env", "must not be empty")

	if !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}

	migTables, err := app.Models.SQLMigrationTables.Get(env)

	if err != nil {
		app.badRequest(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"migrations": migTables})
}

func (app *Application) getMigrationsHandeler(c *gin.Context) {
	qs := c.Request.URL.Query()
	env := app.readString(qs, "env", "")
	table := app.readString(qs, "table", "")

	v := validator.New()
	v.Check(env != "", "env", "must not be empty")
	v.Check(table != "", "table", "must not be empty")

	if !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}

	migs, err := app.Models.SQLMigrationGroup.Get(env, table)
	if err != nil {
		app.badRequest(c, err)
		return
	}

	for _, m := range migs.MigrationsUp {
		sql, err := app.GIT.GetFile(m.File)

		if err != nil {
			m.Script = ""
		} else {
			m.Script = sql
		}
	}

	for _, m := range migs.MigrationsDown {
		sql, err := app.GIT.GetFile(m.File)

		if err != nil {
			m.Script = ""
		} else {
			m.Script = sql
		}

		m.Script = sql

	}
	if err != nil {
		app.badRequest(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"migrations": migs})
}

func (app *Application) setLatestMigrationHandeler(c *gin.Context) {
	var input struct {
		SQLMigrationID int    `json:"sql_migration_id"`
		Env            string `json:"env"`
		Table          string `json:"table"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		app.badRequest(c, err)
	}

	v := validator.New()
	v.Check(input.SQLMigrationID != 0, "sql_migration_id", "must not be empty")
	v.Check(input.Env != "", "env", "must not be empty")
	v.Check(input.Table != "", "table", "must not be empty")

	if !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}

	mig := data.SQLMigrationsLatest{SQLMigrationsID: input.SQLMigrationID, Env: input.Env, Table: input.Table}
	err := app.Models.SQLMigrationsLatest.Update(&mig)

	if err != nil {
		app.badRequest(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"migrations_latest": mig})
}

func (app *Application) deleteMigrationsHandeler(c *gin.Context) {
	var input struct {
		SQLMigrationID int    `json:"sql_migration_id"`
		Env            string `json:"env"`
		Table          string `json:"table"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		app.badRequest(c, err)
	}

	v := validator.New()
	v.Check(input.SQLMigrationID != 0, "sql_migration_id", "must not be empty")
	v.Check(input.Env != "", "env", "must not be empty")
	v.Check(input.Table != "", "table", "must not be empty")

	if !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}

	mig := data.SQLMigration{ID: input.SQLMigrationID, Env: input.Env, SourceTable: input.Table}
	err := app.Models.SQLMigration.Remove(&mig)

	if err != nil {
		app.badRequest(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"migrations": mig})
}

func (app *Application) updateMigrationsHandeler(c *gin.Context) {
	var input struct {
		SQLMigrationID int    `json:"sql_migration_id"`
		Env            string `json:"env"`
		MigratedAtNull bool   `json:"migrated_at_null"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		app.badRequest(c, err)
	}

	v := validator.New()
	v.Check(input.SQLMigrationID != 0, "sql_migration_id", "must not be empty")
	v.Check(input.Env != "", "env", "must not be empty")

	if !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}

	mig := data.SQLMigration{ID: input.SQLMigrationID, Env: input.Env, MigratedAtNull: input.MigratedAtNull}
	err := app.Models.SQLMigration.Update(&mig)

	if err != nil {
		app.badRequest(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"migrations": mig})
}
