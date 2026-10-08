package service

import (
	"github.com/hrnnsx/libra/internal/model"
	"github.com/hrnnsx/libra/internal/repository"
)

type UserService interface {
	GetProfile(userID int64) (*model.User, error)
	UpdateProfile(userID int64, username, email string) (*model.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetProfile(userID int64) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *userService) UpdateProfile(
	userID int64,
	username string,
	email string,
) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	user.Username = username
	user.Email = email

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
