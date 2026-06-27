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

// NewSignupRouter wires the signup route: repo → usecase → controller → route.
func NewSignupRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db, domain.CollectionUser)
	sc := controller.SignupController{
		SignupUsecase: usecase.NewSignupUsecase(ur, timeout),
		Env:          env,
	}
	group.POST("/signup", sc.Signup)
}
