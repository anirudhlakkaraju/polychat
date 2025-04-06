package server

import (
	"context"
	"sync"

	"github.com/anirudhlakkaraju/polychat/config/basicauth"
	"github.com/anirudhlakkaraju/polychat/config/metrics"
	"github.com/anirudhlakkaraju/polychat/config/props"
	"github.com/anirudhlakkaraju/polychat/log"
)

var (
	onceInit        = new(sync.Once)
	implementations = make(map[string]interface{})
	serverKey       = "serverKey"
)

// Init initializes the server configuration once,
// loading properties, basic auth, metrics, logger, etc.
func Init(ctx context.Context) error {
	var err error
	onceInit.Do(func() {
		logger := log.GetCustomLogger()
		p := props.GetProps()

		if implementations[serverKey] == nil {
			port := p.MustGetString("app.port")

			conf := &ServerConfig{
				Auth:      basicauth.GetBasicAuth(),
				Metrics:   metrics.GetMetricsTracker().Log(),
				ReqLogger: log.AttachLoggerToContext(logger),
				Port:      port,
			}

			implementations[serverKey] = conf
			logger.Info("Server intialized successfully")
		}
	})
	return err
}

// GetServerConfig returns the Server Configuration
func GetServerConfig() *ServerConfig {
	v := implementations[serverKey]
	return v.(*ServerConfig)
}
