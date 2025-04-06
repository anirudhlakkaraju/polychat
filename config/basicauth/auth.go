package basicauth

import (
	"context"

	"github.com/gin-gonic/gin"
)

type BasicAuth interface {
	Apply(ctx context.Context, grp *gin.RouterGroup) *gin.RouterGroup
}

// BasicAuth holds basic auth credentials
type BasicAuthImpl struct {
	Users gin.Accounts
	Realm string
}

// Apply applies basic auth to the given router group
func (b *BasicAuthImpl) Apply(ctx context.Context, grp *gin.RouterGroup) *gin.RouterGroup {
	return grp.Group("/", gin.BasicAuthForRealm(b.Users, b.Realm))
}
