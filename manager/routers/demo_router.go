package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Demo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": "demo",
	})
}
