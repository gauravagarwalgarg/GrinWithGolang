package route

import (
	"time"

	"github.com/GauravAgarwalGarg/grin-with-golang/src/08_production/clean_backend/api/controller"
	"github.com/GauravAgarwalGarg/grin-with-golang/src/08_production/clean_backend/bootstrap"
	"github.com/GauravAgarwalGarg/grin-with-golang/src/08_production/clean_backend/domain"
	"github.com/GauravAgarwalGarg/grin-with-golang/src/08_production/clean_backend/mongo"
	"github.com/GauravAgarwalGarg/grin-with-golang/src/08_production/clean_backend/repository"
	"github.com/GauravAgarwalGarg/grin-with-golang/src/08_production/clean_backend/usecase"
	"github.com/gin-gonic/gin"
)

// NewRefreshTokenRouter wires the token refresh route.
func NewRefreshTokenRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db, domain.CollectionUser)
	rtc := &controller.RefreshTokenController{
		RefreshTokenUsecase: usecase.NewRefreshTokenUsecase(ur, timeout),
		Env:                 env,
	}
	group.POST("/refresh", rtc.RefreshToken)
}
