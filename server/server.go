package server

import (
	"context"

	"github.com/anirudhlakkaraju/polychat/config/basicauth"
	"github.com/anirudhlakkaraju/polychat/handlers/base"
	"github.com/gin-gonic/gin"
)

// ServerConfig contains the config, gin handler funcs and port
type ServerConfig struct {
	Auth      basicauth.BasicAuth
	Metrics   gin.HandlerFunc
	ReqLogger gin.HandlerFunc

	Port string
}

// CreateServer returns a gin server that listens on the desired endpoints
func (s *ServerConfig) CreateServer(ctx context.Context) *gin.Engine {

	router := gin.New()
	s.registerBaseHandlers(router)
	s.exposeEndpoints(ctx, router)

	return router
}

func (s *ServerConfig) exposeEndpoints(ctx context.Context, rtr *gin.Engine) {
}

func (s *ServerConfig) registerBaseHandlers(router *gin.Engine) {
	_ = router.SetTrustedProxies(nil)
	router.GET("/app-info", base.GetAppInfo)
	router.GET("/health", base.HealthCheck)
}
