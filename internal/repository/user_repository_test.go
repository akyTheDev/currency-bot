package repository

import (
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/akyTheDev/currency-bot/internal/models"
)

func TestPosgresUserRepository_CreateUser(t *testing.T) {
	tests := []struct {
		name                string
		mockSetup           func(mock sqlmock.Sqlmock)
		chatID              int64
		expectedErrorString string
	}{
		{
			name:   "Success",
			chatID: 12345,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(
					`INSERT INTO users (chat_id) VALUES ($1)
	ON CONFLICT (chat_id) DO NOTHING`,
				)).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedErrorString: "",
		},
		{
			name:   "AlreadyExists",
			chatID: 12345,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(
					`INSERT INTO users (chat_id) VALUES ($1)
	ON CONFLICT (chat_id) DO NOTHING`,
				)).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedErrorString: "already exists",
		},
		{
			name:   "ExecError",
			chatID: 12345,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(
					`INSERT INTO users (chat_id) VALUES ($1)
	ON CONFLICT (chat_id) DO NOTHING`,
				)).WillReturnError(errors.New("ERROR"))
			},
			expectedErrorString: "CreateUser exec: ",
		},
		{
			name:   "RowsAffectedError",
			chatID: 12345,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(
					`INSERT INTO users (chat_id) VALUES ($1)
	ON CONFLICT (chat_id) DO NOTHING`,
				)).WillReturnResult(sqlmock.NewErrorResult(errors.New("rowsAffected failed")))
			},
			expectedErrorString: "CreateUser rows affected: ",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dbMock, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create sqlmock :%v", err)
			}
			defer dbMock.Close()

			tc.mockSetup(mock)

			repo := NewPostgresUserRepository(dbMock)
			err = repo.CreateUser(tc.chatID)

			if tc.expectedErrorString == "" {
				if err != nil {
					t.Errorf("Expected no error, got :%v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error string %s, got nil", tc.expectedErrorString)
				}
				if !strings.Contains(err.Error(), tc.expectedErrorString) {
					t.Errorf("error = %q; want it to contain %s", err.Error(), tc.expectedErrorString)
				}
			}
		})
	}
}

func TestPosgresUserRepository_DeleteUser(t *testing.T) {
	tests := []struct {
		name                string
		mockSetup           func(mock sqlmock.Sqlmock)
		chatID              int64
		expectedErrorString string
	}{
		{
			name:   "Success",
			chatID: 12345,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(
					`DELETE FROM users WHERE chat_id = $1`,
				)).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedErrorString: "",
		},
		{
			name:   "NoRowsError",
			chatID: 12345,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(
					`DELETE FROM users WHERE chat_id = $1`,
				)).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedErrorString: "not found",
		},
		{
			name:   "ExecError",
			chatID: 12345,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(
					`DELETE FROM users WHERE chat_id = $1`,
				)).WillReturnError(errors.New("ERROR"))
			},
			expectedErrorString: "DeleteUser exec: ",
		},
		{
			name:   "RowsAffectedError",
			chatID: 12345,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(
					`DELETE FROM users WHERE chat_id = $1`,
				)).WillReturnResult(sqlmock.NewErrorResult(errors.New("rowsAffected failed")))
			},
			expectedErrorString: "DeleteUser rows affected: ",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dbMock, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create sqlmock :%v", err)
			}
			defer dbMock.Close()

			tc.mockSetup(mock)

			repo := NewPostgresUserRepository(dbMock)
			err = repo.DeleteUser(tc.chatID)

			if tc.expectedErrorString == "" {
				if err != nil {
					t.Errorf("Expected no error, got :%v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error string %s, got nil", tc.expectedErrorString)
				}
				if !strings.Contains(err.Error(), tc.expectedErrorString) {
					t.Errorf("error = %q; want it to contain %s", err.Error(), tc.expectedErrorString)
				}
			}
		})
	}
}

func TestPosgresUserRepository_GetAllUsers(t *testing.T) {
	tests := []struct {
		name                string
		mockSetup           func(mock sqlmock.Sqlmock)
		expected            []models.User
		expectedErrorString string
	}{
		{
			name: "Success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "chat_id"}).AddRow(123, 1230).AddRow(124, 1240)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, chat_id FROM users")).WillReturnRows(rows)
			},
			expected: []models.User{
				{ID: 123, ChatID: 1230},
				{ID: 124, ChatID: 1240},
			},
			expectedErrorString: "",
		},
		{
			name: "QueryError",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, chat_id FROM users")).WillReturnError(errors.New("ERROR"))
			},
			expected:            nil,
			expectedErrorString: "GetAllUsers query:",
		},
		{
			name: "ScanError",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "chat_id"}).AddRow("not_integer", "chat_user_1")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, chat_id FROM users")).WillReturnRows(rows)
			},
			expected:            nil,
			expectedErrorString: "GetAllUsers scan:",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dbMock, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create sqlmock :%v", err)
			}
			defer dbMock.Close()

			tc.mockSetup(mock)

			repo := NewPostgresUserRepository(dbMock)
			users, err := repo.GetAllUsers()

			if tc.expectedErrorString == "" {
				if err != nil {
					t.Errorf("Expected no error, got :%v", err)
				}

				if len(users) != len(tc.expected) {
					t.Errorf("Expected %d users, got %d", len(tc.expected), len(users))
				}

				for i := range users {
					if users[i] != tc.expected[i] {
						t.Errorf("user[%d] = %+v; want %+v", i, users[i], tc.expected[i])
					}
				}

			} else {
				if err == nil {
					t.Errorf("Expected error string %s, got nil", tc.expectedErrorString)
				}
				if !strings.Contains(err.Error(), tc.expectedErrorString) {
					t.Errorf("error = %q; want it to contain %s", err.Error(), tc.expectedErrorString)
				}
			}

		})
	}
}
