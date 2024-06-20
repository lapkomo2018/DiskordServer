package service

import (
	"errors"
	"github.com/lapkomo2018/DiskordServer/internal/core"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
	"time"
)

type UserStorage interface {
	Create(user *core.User) error
	Exists(email string) error
	First(user *core.User, cond ...interface{}) error
	LoadFiles(user *core.User) error
}

type TokenManager interface {
	NewJWT(userId string, ttl time.Duration) (string, error)
	Parse(accessToken string) (string, error)
}

type UserService struct {
	storage      UserStorage
	tokenManager TokenManager

	accessTokenTTL time.Duration
}

func NewUserService(storage UserStorage, tokenManager TokenManager, accessTokenTTL time.Duration) *UserService {
	log.Printf("Created user service")
	return &UserService{
		storage:        storage,
		tokenManager:   tokenManager,
		accessTokenTTL: accessTokenTTL,
	}
}

func (us *UserService) Create(email, password string) (*core.User, error) {
	// look up for user
	if err := us.storage.Exists(email); err == nil {
		return nil, errors.New("user already exists")
	}

	//hash pass
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &core.User{
		Email:    email,
		Password: string(passHash),
	}

	if err := us.storage.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) Login(email, password string) (string, error) {
	user := &core.User{
		Email: email,
	}
	if err := us.storage.First(user); err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}

	token, err := us.tokenManager.NewJWT(strconv.Itoa(int(user.ID)), us.accessTokenTTL)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (us *UserService) LoadFiles(user *core.User) error {
	return us.storage.LoadFiles(user)
}

func (us *UserService) GetUserFromToken(token string) (*core.User, error) {
	userId, err := us.tokenManager.Parse(token)
	if err != nil {
		return nil, err
	}

	user := &core.User{}
	if err := us.storage.First(user, userId); err != nil {
		return nil, err
	}

	return user, nil
}
