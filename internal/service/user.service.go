package service

import (
	"errors"

	"github.com/hrnnsx/libra/internal/model"
	"github.com/hrnnsx/libra/internal/repository"
	"gorm.io/gorm"
)

type UserService interface {
	GetProfile(userID int64) (*model.User, error)
	UpdateProfile(userID int64, username, email *string) (*model.User, error)
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
	username, email *string,
) (*model.User, error) {

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if username != nil && *username != user.Username {
		existingUser, err := s.userRepo.FindByUsername(*username)

		if err == nil && existingUser.ID != userID {
			return nil, ErrUsernameExists
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		user.Username = *username
	}

	if email != nil && *email != user.Email {
		existingUser, err := s.userRepo.FindByEmail(*email)

		if err == nil && existingUser.ID != userID {
			return nil, ErrEmailExists
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		user.Email = *email
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
