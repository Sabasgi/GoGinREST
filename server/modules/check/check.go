package check

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Init(r, o *gin.RouterGroup) {
	o.GET("/server/check", CheckStatus)
}
func CheckStatus(g *gin.Context) {
	log.Println("Server check", time.Now())
	g.JSON(http.StatusOK, gin.H{
		"ServerTime": time.Now(),
		"STATUS":     "Running",
	})
}
