package config

import (
	"backend_golang/internal/domain"
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var OAuthConfig *OAuth
var JWTConfig *JWT

type OAuth struct {
	config oauth2.Config
}

type JWT struct {
	secret string
}

func NewOAuth() *OAuth {
	scopes := strings.Split(os.Getenv("OAUTH_SCOPES"), ",")
	return &OAuth{
		config: oauth2.Config{
			ClientID:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),
			RedirectURL:  os.Getenv("OAUTH_REDIRECT_URL"),
			Scopes:       scopes,
			Endpoint:     google.Endpoint,
		},
	}
}

func NewJWT() *JWT {
	return &JWT{
		secret: os.Getenv("JWT_SIGN_KEY"),
	}
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	OAuthConfig = NewOAuth()
	JWTConfig = NewJWT()

	log.Println("OAuthConfig", OAuthConfig)
	log.Println("JWTConfig", JWTConfig)
}

func (o *OAuth) GetAccessToken(c context.Context, code string) (*oauth2.Token, error) {
	return o.config.Exchange(c, code)
}

func (o *OAuth) AuthCodeURL(state string) string {
	return o.config.AuthCodeURL(state)
}

func (o *OAuth) GetMember(c context.Context, token *oauth2.Token) (*domain.Member, error) {
	client := o.config.Client(c, token)
	// TODO : 環境変数に定義する
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		log.Println("Failed to get user info from google", err)
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]any
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		log.Println("Failed to decode user info from google", err)
		return nil, err
	}

	log.Println("userInfo", userInfo)
	return &domain.Member{
		ID:       userInfo["sub"].(string),
		Email:    userInfo["email"].(string),
		Picture:  userInfo["picture"].(string),
		Nickname: userInfo["name"].(string),
	}, nil
}

func (j *JWT) GetSecretKey() []byte {
	return []byte(j.secret)
}
