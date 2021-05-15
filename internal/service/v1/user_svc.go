package v1

import (
	"bank-system-go/internal/model"
	"bank-system-go/internal/repository"
	"bank-system-go/internal/service"
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var _ service.UserService = (*UserService)(nil)

func NewUserService(userRepo repository.UserRepository) service.UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

type UserService struct {
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
