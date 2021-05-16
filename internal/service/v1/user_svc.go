package v1

import (
	"bank-system-go/internal/config"
	"bank-system-go/internal/model"
	"bank-system-go/internal/repository"
	"bank-system-go/internal/service"
	"context"
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var _ service.UserService = (*UserService)(nil)

func NewUserService(config config.Config, userRepo repository.UserRepository) service.UserService {
	return &UserService{
		config:   config,
		userRepo: userRepo,
	}
}

type UserService struct {
	config   config.Config
	userRepo repository.UserRepository
}

func (svc *UserService) Register(ctx context.Context, user model.User) (model.User, error) {
	encodePassword, err := bcrypt.GenerateFromPassword(user.Password, bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}
	user.UUID = uuid.New().String()
	user.Password = encodePassword
	return user, svc.userRepo.CreateUser(ctx, &user)
}

func (svc *UserService) VerifyPassword(ctx context.Context, user model.User, password string) error {
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}

func (svc *UserService) Login(ctx context.Context, account, password, ip string) (string, error) {
	user, err := svc.userRepo.GetUser(ctx, model.User{
		Account: account,
	})
	if err != nil {
		return "", err
	}
	err = svc.VerifyPassword(ctx, user, password)
	if err != nil {
		return "", err
	}

	accessLog := &model.AccessLog{
		UserID: user.ID,
		IP:     ip,
	}
	err = svc.userRepo.CreateAccessLog(ctx, accessLog)
	if err != nil {
		return "", err
	}
	user.LatestAccessAt = &accessLog.CreatedAt
	err = svc.userRepo.UpdateUser(ctx, model.User{ID: user.ID}, &user)
	if err != nil {
		return "", err
	}

	now := time.Now()
	claim := jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
		ExpiresAt: now.Add(svc.config.JWT.Age).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claim)
	return token.SigningString()
}

func (svc *UserService) VerifyJWT(ctx context.Context, tokenStr string) error {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return svc.config.JWT.Secret, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return errors.New("that's not even a token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return errors.New("token is expired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return errors.New("token not active yet")
			} else {
				return errors.New("invalid token")
			}
		}
	}
	if !token.Valid {
		return errors.New("invalid token")
	}

	return nil
}
