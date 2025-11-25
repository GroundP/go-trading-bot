package api

import (
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	v1Group := router.Group("/api/v1")
	{
		v1Group.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "OK"})
		})
	}
	return router
}
