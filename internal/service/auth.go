package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	vk_test_task "vk-test-task"
	"vk-test-task/internal/storage"
)

type Auth interface {
	Authentication(user vk_test_task.User) (vk_test_task.Tokens, error)
	CreateUser(newUser vk_test_task.CreateUser) error
	ParseToken(accessToken string) (*ParseDataUser, error)
}

type AuthService struct {
	repo *storage.StorageServ
}

func NewAuthService(repo *storage.StorageServ) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (a *AuthService) Authentication(user vk_test_task.User) (vk_test_task.Tokens, error) {
	const op = "service.Authentication"

	user.Password = HashPassword(user.Password)

	userData, err := a.repo.Auth.GetUserInfo(user)
	if err != nil {
		return vk_test_task.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	accessToken, err := CreateAccessToken(userData)
	if err != nil {
		return vk_test_task.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	refreshToken, err := CreateRefreshToken()
	if err != nil {
		return vk_test_task.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	tokens := vk_test_task.Tokens{
		Access_token:  accessToken,
		Refresh_token: refreshToken,
	}

	a.repo.MemoryStorage.SaveTokens(userData.ID, tokens)

	return tokens, nil
}

func (a *AuthService) CreateUser(newUser vk_test_task.CreateUser) error {
	const op = "service.CreateUser"

	newUser.Password = HashPassword(newUser.Password)
	err := a.repo.Auth.CreateUser(newUser)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func HashPassword(password string) string {
	data := []byte(password + salt)
	hashData := sha256.Sum256(data)
	hashString := hex.EncodeToString(hashData[:])

	return hashString
}
