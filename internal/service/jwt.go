package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"math/rand"
	"time"
	"vk-test-task/internal/storage"
)

const (
	refresh_tokenTtl = time.Hour * 24 * 30
	access_tokenTtl  = time.Minute * 15
	signInKey        = "@(#tf53$*#$(RHfverib}#Rfrte)"
	salt             = "lsd2#tfv%2"
)

type JWT interface {
	CreateAccessToken(id int64) (string, error)
	CreateRefreshToken() (string, error)
	ParseToken(accessToken string) (*ParseDataUser, error)
}

type TokenClaims struct {
	jwt.StandardClaims
	Id    int64  `json:"id"`
	Role  string `json:"role"`
	Email string `json:"email"`
}

type ParseDataUser struct {
	ID    int64
	Email string
	Role  string
}

func CreateAccessToken(userData storage.UserData) (string, error) {
	const op = "service.CreateAccessToken"

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(access_tokenTtl).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Id:    userData.ID,
		Role:  userData.Role,
		Email: userData.Email,
	})

	signedAccessToken, err := accessToken.SignedString([]byte(signInKey))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return signedAccessToken, nil
}

func CreateRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	_, err := r.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *AuthService) ParseToken(accessToken string) (*ParseDataUser, error) {
	const op = "service.ParseToken"

	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s: %w", op, errors.New("token verification error"))
		}
		return []byte(signInKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, errors.New("failed to parse token claims")
	}

	if !ok {
		return nil, err
	}

	res := ParseDataUser{
		ID:    claims.Id,
		Email: claims.Email,
		Role:  claims.Role,
	}

	fmt.Println(res)
	return &res, nil
}
