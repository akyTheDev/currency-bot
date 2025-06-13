package service

import (
	"errors"
	"log"
	"os"
	"testing"

	"github.com/akyTheDev/currency-bot/internal/domain"
	"github.com/akyTheDev/currency-bot/internal/fetcher"
	"github.com/akyTheDev/currency-bot/internal/models"
)

type fakeUserRepoNotifyService struct {
	users []models.User
	err   error
}

func (f *fakeUserRepoNotifyService) CreateUser(chatID int64) error       { return nil }
func (f *fakeUserRepoNotifyService) DeleteUser(chatID int64) error       { return nil }
func (f *fakeUserRepoNotifyService) GetAllUsers() ([]models.User, error) { return f.users, f.err }

type fakeRateFetcher struct {
	rate *fetcher.Rate
	err  error
}

func (f *fakeRateFetcher) FetchRate() (*fetcher.Rate, error) {
	return f.rate, f.err
}

var logger = log.New(os.Stdout, "", 0)

func TestGetUsersAndCurrencyRate(t *testing.T) {

	tests := []struct {
		name        string
		fetcherRate *fetcher.Rate
		fetcherErr  error
		repoUsers   []models.User
		repoErr     error
		wantIDs     []int64
		wantRate    *fetcher.Rate
		wantErr     error
	}{
		{
			name:        "FetchRateError",
			fetcherRate: nil,
			fetcherErr:  errors.New("fetch failed"),
			repoUsers:   nil,
			repoErr:     nil,
			wantIDs:     nil,
			wantRate:    nil,
			wantErr:     domain.ErrGeneric,
		},
		{
			name:        "GetAllUsersError",
			fetcherRate: &fetcher.Rate{Selling: 25.5, Buying: 24.5},
			fetcherErr:  nil,
			repoUsers:   nil,
			repoErr:     errors.New("db failed"),
			wantIDs:     nil,
			wantRate:    &fetcher.Rate{Selling: 25.5, Buying: 24.5},
			wantErr:     domain.ErrGeneric,
		},
		{
			name:        "NoSubscribers",
			fetcherRate: &fetcher.Rate{Selling: 25.5, Buying: 24.5},
			fetcherErr:  nil,
			repoUsers:   []models.User{},
			repoErr:     nil,
			wantIDs:     nil,
			wantRate:    &fetcher.Rate{Selling: 25.5, Buying: 24.5},
			wantErr:     nil,
		},
		{
			name:        "SomeSubscribers",
			fetcherRate: &fetcher.Rate{Selling: 26.5, Buying: 24.5},
			fetcherErr:  nil,
			repoUsers: []models.User{
				{ChatID: 101},
				{ChatID: 202},
				{ChatID: 303},
			},
			repoErr:  nil,
			wantIDs:  []int64{101, 202, 303},
			wantRate: &fetcher.Rate{Selling: 26.5, Buying: 24.5},
			wantErr:  nil,
		},
		// },
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ff := &fakeRateFetcher{
				rate: tc.fetcherRate,
				err:  tc.fetcherErr,
			}
			fr := &fakeUserRepoNotifyService{
				users: tc.repoUsers,
				err:   tc.repoErr,
			}

			ns := NewNotifyService(logger, fr, ff)

			ids, rate, err := ns.GetUsersAndCurrencyRate()

			if tc.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error: %v, got nil", tc.wantErr)
				}
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected error: %v, got %v", tc.wantErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}

			if len(ids) != len(tc.wantIDs) {
				t.Fatalf("Expected ids length: %d, got %d", len(tc.wantIDs), len(ids))
			}
			for i := range ids {
				if ids[i] != tc.wantIDs[i] {
					t.Errorf("ids[%d] = %d, want %d", i, ids[i], tc.wantIDs[i])
				}
			}

			if tc.wantRate == nil {
				if rate != nil {
					t.Errorf("expected nil rate, got %v", rate)
				}
			} else {
				if rate == nil {
					t.Errorf("expected rate %v, got nil", tc.wantRate)
				} else if *rate != *tc.wantRate {
					t.Errorf("expected rate %v, got %v", tc.wantRate, rate)
				}
			}
		})
	}
}
