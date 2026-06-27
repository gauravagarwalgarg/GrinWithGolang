package controller

import (
	"net/http"

	"github.com/GauravAgarwalGarg/grin-with-golang/src/08_production/clean_backend/domain"
	"github.com/gin-gonic/gin"
)

// ProfileController handles user profile retrieval.
type ProfileController struct {
	ProfileUsecase domain.ProfileUsecase
}

// Fetch handles GET /profile.
// The user ID is extracted from the JWT by the auth middleware.
func (pc *ProfileController) Fetch(c *gin.Context) {
	userID := c.GetString("x-user-id")

	profile, err := pc.ProfileUsecase.GetProfileByID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}
