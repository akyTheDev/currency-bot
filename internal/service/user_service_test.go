package service

import (
	"errors"
	"log"
	"os"
	"testing"

	"github.com/akyTheDev/currency-bot/internal/domain"
	"github.com/akyTheDev/currency-bot/internal/models"
)

type fakeUserRepo struct {
	createErr  error
	deleteErr  error
	lastChatId int64
}

func (f *fakeUserRepo) CreateUser(chatID int64) error {
	f.lastChatId = chatID
	return f.createErr
}

func (f *fakeUserRepo) DeleteUser(chatID int64) error {
	f.lastChatId = chatID
	return f.deleteErr
}

func (f *fakeUserRepo) GetAllUsers() ([]models.User, error) {
	return nil, nil
}

func TestUserServiceRegister(t *testing.T) {
	tests := []struct {
		name        string
		repoErr     error
		expectedErr error
	}{
		{
			name:        "Success",
			repoErr:     nil,
			expectedErr: nil,
		},
		{
			name:        "AlreadyExists",
			repoErr:     domain.ErrUserAlreadyExists,
			expectedErr: domain.ErrUserAlreadyExists,
		},
		{
			name:        "Other Error",
			repoErr:     errors.New("other error"),
			expectedErr: domain.ErrGeneric,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := &fakeUserRepo{
				createErr: tc.repoErr,
				deleteErr: nil,
			}
			u := NewUserService(f, log.New(os.Stdout, "", 0))

			err := u.Register(12345)

			if tc.expectedErr == nil {
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}
			} else {
				if err != tc.expectedErr {
					t.Fatalf("Expected error: %v, got %v", tc.expectedErr, err)
				}
			}

			if f.lastChatId != 12345 {
				t.Errorf("Expected last chat id: %d, got %d", 12345, f.lastChatId)
			}
		})
	}
}

func TestUserServiceDelete(t *testing.T) {
	tests := []struct {
		name        string
		repoErr     error
		expectedErr error
	}{
		{
			name:        "Success",
			repoErr:     nil,
			expectedErr: nil,
		},
		{
			name:        "NotFound",
			repoErr:     domain.ErrUserNotFound,
			expectedErr: domain.ErrUserNotFound,
		},
		{
			name:        "Other Error",
			repoErr:     errors.New("other error"),
			expectedErr: domain.ErrGeneric,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := &fakeUserRepo{
				deleteErr: tc.repoErr,
				createErr: nil,
			}
			u := NewUserService(f, log.New(os.Stdout, "", 0))

			err := u.Delete(12345)

			if tc.expectedErr == nil {
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}
			} else {
				if err != tc.expectedErr {
					t.Fatalf("Expected error: %v, got %v", tc.expectedErr, err)
				}
			}

			if f.lastChatId != 12345 {
				t.Errorf("Expected last chat id: %d, got %d", 12345, f.lastChatId)
			}
		})
	}
}
