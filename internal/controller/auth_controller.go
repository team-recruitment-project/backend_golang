package controller

import (
	"backend_golang/internal/controller/request"
	"backend_golang/internal/models"
	"backend_golang/internal/service"
	smodels "backend_golang/internal/service/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthController interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	GoogleCallback(c *gin.Context)
	Signup(c *gin.Context)
	GetMember(c *gin.Context)
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

func (a *authController) Logout(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   -1,
	})
	c.Status(http.StatusOK)
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

	c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/")
}

func (a *authController) Signup(c *gin.Context) {
	userID := c.Value("userID").(string)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	req := &request.SignUpRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		validationErrors := make([]models.ValidationError, 0)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, models.NewValidationError(err))
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": validationErrors,
		})
		return
	}

	signup := smodels.SignupMember{
		Bio:           req.Bio,
		PreferredRole: models.Role(req.PreferredRole),
	}
	memberID, err := a.authService.Signup(c, userID, signup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"memberID": memberID})
}

func (a *authController) GetMember(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	member, err := a.authService.GetMember(c, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, member)
}
