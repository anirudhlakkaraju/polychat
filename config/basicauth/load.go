package basicauth

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/anirudhlakkaraju/polychat/config/props"
	"github.com/anirudhlakkaraju/polychat/log"
)

var (
	onceInit       = new(sync.Once)
	implementation = make(map[string]interface{})
	basicAuthKey   = "basicAuthKey"
	basicAuthRealm = "Restricted"
)

// Init initializes the basic auth configuration once,
// loading credentials from properties and storing them in memory.
func Init(ctx context.Context) error {
	var err error
	onceInit.Do(func() {
		logger := log.LoggerFromContext(ctx)
		p := props.GetProps()

		if implementation[basicAuthKey] == nil {
			basicAuth := p.MustGetString("basic.auth.credentials") // Expected format: "username:password"
			creds := strings.SplitN(basicAuth, ":", 2)
			if len(creds) != 2 {
				err = fmt.Errorf("invalid basic auth credentials format")
			}

			users := gin.Accounts{
				creds[0]: creds[1], // username -> password
			}

			implementation[basicAuthKey] = &BasicAuthImpl{
				Users: users,
				Realm: basicAuthRealm,
			}

			logger.Info("BasicAuth initialized successfully")
		}
	})
	return err
}

// GetBasicAuth returns the BasicAuth object
func GetBasicAuth() BasicAuth {
	v := implementation[basicAuthKey]
	return v.(BasicAuth)
}
