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

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Login(c context.Context) models.LoginResponse
	GoogleCallback(c context.Context, code string) (string, error)
	Signup(c context.Context, userID string, signup models.SignupMember) (string, error)
	GetMember(c context.Context, userID string) (*models.UserResponse, error)
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

func (a *authService) GoogleCallback(c context.Context, code string) (string, error) {
	token, err := config.OAuthConfig.GetAccessToken(c, code)
	if err != nil {
		return "", err
	}

	member, err := config.OAuthConfig.GetMember(c, token)
	if err != nil {
		return "", err
	}

	if foundMember, err := a.authRepository.GetMemberByID(c, member.ID); foundMember != nil && err == nil {
		return a.generateAccessToken(foundMember.ID)
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

	return a.generateAccessToken(member.ID)
}

func (*authService) generateAccessToken(id string) (string, error) {
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "team-recruitment",
		"sub": id,
		"exp": time.Now().Add(time.Minute * 30).Unix(),
	}).SignedString(config.JWTConfig.GetSecretKey())
	if err != nil {
		log.Printf("error creating access token: %v", err)
		return "", err
	}

	return accessToken, nil
}

func (a *authService) Signup(c context.Context, userID string, signup models.SignupMember) (string, error) {
	foundMember, err := a.authRepository.GetTransientMemberByID(c, userID)
	if err != nil {
		return "", err
	}

	err = a.authRepository.DeleteTransientMemberByID(c, foundMember.ID)
	if err != nil {
		return "", err
	}

	member, err := a.authRepository.CreateMember(c, &domain.Member{
		ID:            userID,
		Bio:           signup.Bio,
		PreferredRole: string(signup.PreferredRole),
	})
	if err != nil {
		return "", err
	}

	return member.ID, nil
}

func (a *authService) GetMember(c context.Context, userID string) (*models.UserResponse, error) {
	member, err := a.authRepository.GetMemberByID(c, userID)

	if err != nil {
		if ent.IsNotFound(err) {
			transientMember, err := a.authRepository.GetTransientMemberByID(c, userID)
			if err != nil {
				return nil, err
			}
			return &models.UserResponse{
				ID:        transientMember.ID,
				Transient: true,
			}, nil
		}
		return nil, err
	}

	return &models.UserResponse{
		ID:            member.ID,
		Email:         member.Email,
		Nickname:      member.Nickname,
		Picture:       member.Picture,
		Bio:           member.Bio,
		PreferredRole: member.PreferredRole,
		Transient:     false,
	}, nil
}
