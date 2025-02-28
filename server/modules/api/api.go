package api

import (
	"goGinRest/modules/check"

	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	o := router.Group("/o")
	r := router.Group("/r")
	check.Init(r, o)
}
