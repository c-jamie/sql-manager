package mocks

import "github.com/gin-gonic/gin"

type MockMiddleware struct {
}

func (mi *MockMiddleware) Authenticate(c *gin.Context) {
	c.Next()
}

func (mi *MockMiddleware) Authorize(code string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func (mi *MockMiddleware) HasPermission(code string, permissions []string) bool {
	return true
}