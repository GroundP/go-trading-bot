package api

import (
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	v1Group := router.Group("/api/v1")
	{
		v1Group.GET("/positions", .)
	}
	return router
}
