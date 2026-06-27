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

// NewProfileRouter wires the profile route (protected).
func NewProfileRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db, domain.CollectionUser)
	pc := &controller.ProfileController{
		ProfileUsecase: usecase.NewProfileUsecase(ur, timeout),
	}
	group.GET("/profile", pc.Fetch)
}
