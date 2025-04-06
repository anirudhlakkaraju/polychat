package base

import (
	"net/http"
	"os"
	"time"

	"github.com/anirudhlakkaraju/polychat/config/props"
	"github.com/gin-gonic/gin"
)

func GetAppInfo(c *gin.Context) {
	p := props.GetProps()
	host, _ := os.Hostname()

	c.JSON(http.StatusOK, gin.H{
		"name":    p.MustGetString("app.name"),
		"version": p.MustGetString("app.version"),
		"env":     os.Getenv("configEnvironment"),
		"host":    host,
		"upSince": time.Now().Format(time.RFC1123),
	})
}
