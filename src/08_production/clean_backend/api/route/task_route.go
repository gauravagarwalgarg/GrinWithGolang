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

// NewTaskRouter wires the task routes (protected).
func NewTaskRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	tr := repository.NewTaskRepository(db, domain.CollectionTask)
	tc := &controller.TaskController{
		TaskUsecase: usecase.NewTaskUsecase(tr, timeout),
	}
	group.GET("/task", tc.Fetch)
	group.POST("/task", tc.Create)
}
