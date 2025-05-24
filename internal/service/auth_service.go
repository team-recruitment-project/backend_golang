package service

import (
	config "backend_golang/configs"
	"backend_golang/ent"
	"backend_golang/internal/domain"
	"backend_golang/internal/repository"
	"backend_golang/internal/service/models"
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Login(c context.Context) models.LoginResponse
	GoogleCallback(c *gin.Context, code string) (string, error)
	Signup(c *gin.Context, userID string, signup models.SignupMember) (error, string)
}

type authService struct {
	authRepository repository.AuthRepository
}

func NewAuthService(authRepository repository.AuthRepository) AuthService {
	return &authService{
		authRepository: authRepository,
	}
}

func (a *authService) Login(c context.Context) models.LoginResponse {
	url := config.OAuthConfig.AuthCodeURL("state")
	return models.LoginResponse{
		URL: url,
	}
}

func (a *authService) GoogleCallback(c *gin.Context, code string) (string, error) {
	token, err := config.OAuthConfig.GetAccessToken(c, code)
	if err != nil {
		return "", err
	}

	member, err := config.OAuthConfig.GetMember(c, token)
	if err != nil {
		return "", err
	}

	if _, err := a.authRepository.GetTransientMemberByID(c, member.ID); err != nil {
		if ent.IsNotFound(err) {
			_, err := a.authRepository.CreateTransientMember(c, &domain.TransientMember{
				ID:       member.ID,
				Email:    member.Email,
				Picture:  member.Picture,
				Nickname: member.Nickname,
			})
			if err != nil {
				return "", err
			}
		}
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "team-recruitment",
		"sub": member.ID,
		"exp": time.Now().Add(time.Minute * 30).Unix(),
	}).SignedString(config.JWTConfig.GetSecretKey())
	if err != nil {
		log.Printf("error creating access token: %v", err)
		return "", err
	}

	return accessToken, nil
}

func (a *authService) Signup(c *gin.Context, userID string, signup models.SignupMember) (error, string) {
	member, err := a.authRepository.CreateMember(c, &domain.Member{
		ID:            userID,
		Bio:           signup.Bio,
		PreferredRole: string(signup.PreferredRole),
	})
	if err != nil {
		return err, ""
	}

	return nil, member.ID
}
