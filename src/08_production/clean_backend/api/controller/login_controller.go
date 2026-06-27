package controller

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/GauravAgarwalGarg/grin-with-golang/src/08_production/clean_backend/bootstrap"
	"github.com/GauravAgarwalGarg/grin-with-golang/src/08_production/clean_backend/domain"
	"github.com/gin-gonic/gin"
)

// LoginController handles user authentication.
type LoginController struct {
	LoginUsecase domain.LoginUsecase
	Env          *bootstrap.Env
}

// Login handles POST /login.
// Flow: validate → find user → compare password → issue tokens
func (lc *LoginController) Login(c *gin.Context) {
	var request domain.LoginRequest

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	user, err := lc.LoginUsecase.GetUserByEmail(c, request.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "User not found with the given email"})
		return
	}

	// bcrypt.CompareHashAndPassword is timing-safe (prevents timing attacks)
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)) != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Invalid credentials"})
		return
	}

	accessToken, err := lc.LoginUsecase.CreateAccessToken(&user, lc.Env.AccessTokenSecret, lc.Env.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	refreshToken, err := lc.LoginUsecase.CreateRefreshToken(&user, lc.Env.RefreshTokenSecret, lc.Env.RefreshTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
