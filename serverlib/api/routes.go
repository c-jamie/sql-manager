package api

import (
	"io"
	"time"

	"github.com/c-jamie/sql-manager/serverlib/log"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Logrus(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now().UTC()
		path := c.Request.URL.Path
		c.Next()
		end := time.Now().UTC()
		latency := end.Sub(start)
		logger.WithFields(logrus.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       path,
			"ip":         c.ClientIP(),
			"duration":   latency,
			"user_agent": c.Request.UserAgent(),
		}).Info()
	}
}

// WriteFunc convert func to io.Writer.
type writeFunc func([]byte) (int, error)

func (fn writeFunc) Write(data []byte) (int, error) {
	return fn(data)
}

func newLogrusWrite() io.Writer {
	return writeFunc(func(data []byte) (int, error) {
		log.Debugf("%s", data)
		return 0, nil
	})
}

func (app *Application) Routes() *gin.Engine {

	gin.DefaultWriter = newLogrusWrite()
	router := gin.New()
	router.Use(Logrus(log.Log))

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found", "uri": c.Request.RequestURI})
	})

	public := router.Group("/" + app.Config.Version)

	public.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"version": app.Config.Version})
	})

	private := router.Group("/" + app.Config.Version)

	authenticate := func() gin.HandlerFunc { return app.Middleware.Authenticate }

	private.Use(authenticate())
	private.GET("/healthcheck", app.Middleware.Authorize("/users-read"), func(c *gin.Context) {
		c.JSON(200, gin.H{"version": app.Config.Version})
	})
	private.GET("/files/list", app.Middleware.Authorize("/users-read"), app.listFilesHandeler)
	private.GET("/files", app.Middleware.Authorize("/users-write"), app.getFilesHandeler)
	private.POST("/files", app.Middleware.Authorize("/users-write"), app.registerFilesHandeler)

	private.GET("/projects", app.Middleware.Authorize("/users-write"), app.getProjectsHandeler)

	private.POST("/migrations", app.Middleware.Authorize("/users-write"), app.addMigrationsHandeler)
	private.GET("/migrations", app.Middleware.Authorize("/users-write"), app.getMigrationsHandeler)
	private.DELETE("/migrations", app.Middleware.Authorize("/users-write"), app.deleteMigrationsHandeler)
	private.PATCH("/migrations", app.Middleware.Authorize("/users-write"), app.updateMigrationsHandeler)
	private.GET("/migrations/table", app.Middleware.Authorize("/users-write"), app.getMigrationTablesHandeler)
	private.POST("/migrations/latest", app.Middleware.Authorize("/users-write"), app.setLatestMigrationHandeler)

	return router
}
