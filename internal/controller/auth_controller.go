package controller

import (
	"backend_golang/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Login(c *gin.Context)
	GoogleCallback(c *gin.Context)
}

type authController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) AuthController {
	return &authController{
		authService: authService,
	}
}

func (a *authController) Login(c *gin.Context) {
	response := a.authService.Login(c)
	c.JSON(http.StatusOK, response)
}

func (a *authController) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	accessToken, err := a.authService.GoogleCallback(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		// Secure:   false, // TODO : 本番環境では true にする
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   int((time.Minute * 30).Seconds()),
	})

	c.Redirect(http.StatusTemporaryRedirect, "/")
}
