package service

import (
	"bank-system-go/internal/model"
	"context"

	"github.com/dgrijalva/jwt-go"
)

type UserService interface {
	Register(ctx context.Context, user model.User) (model.User, error)
	VerifyPassword(ctx context.Context, user model.User, password string) error
	Login(ctx context.Context, account, password, ip string) (string, error)
	VerifyJWT(ctx context.Context, tokenStr string) (jwt.StandardClaims, error)
	GetUser(ctx context.Context, user model.User) (model.User, error)
}
