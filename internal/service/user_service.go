package service

import (
	"log"

	"github.com/akyTheDev/currency-bot/internal/domain"
	"github.com/akyTheDev/currency-bot/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepository
	logger   *log.Logger
}

func NewUserService(userRepo repository.UserRepository, logger *log.Logger) *UserService {
	return &UserService{userRepo: userRepo, logger: logger}
}

func (s *UserService) Register(chatID int64) error {
	err := s.userRepo.CreateUser(chatID)
	if err != nil {
		s.logger.Printf("ERROR: UserService:Register: %v\n", err)
		if err == domain.ErrUserAlreadyExists {
			return err
		}
		return domain.ErrGeneric
	}
	return nil
}

func (s *UserService) Delete(chatID int64) error {
	err := s.userRepo.DeleteUser(chatID)
	if err != nil {
		s.logger.Printf("ERROR: UserService:Delete: %v\n", err)
		if err == domain.ErrUserNotFound {
			return err
		}
		return domain.ErrGeneric
	}
	return nil
}
