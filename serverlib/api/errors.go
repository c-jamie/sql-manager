package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (app *Application) failedValidationResponse(c *gin.Context, errors map[string]string) {
	c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errors, "message": "invalid json provided"})
}

func (app *Application) badRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
}