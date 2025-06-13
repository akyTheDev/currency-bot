package service

import (
	"log"

	"github.com/akyTheDev/currency-bot/internal/domain"
	"github.com/akyTheDev/currency-bot/internal/fetcher"
	"github.com/akyTheDev/currency-bot/internal/repository"
)

type NotifyService struct {
	logger         *log.Logger
	userRepository repository.UserRepository
	rateFetch      fetcher.RateFetcher
}

func NewNotifyService(logger *log.Logger, userRepository repository.UserRepository, rateFetch fetcher.RateFetcher) *NotifyService {
	return &NotifyService{
		logger:         logger,
		userRepository: userRepository,
		rateFetch:      rateFetch,
	}
}

func (ns *NotifyService) GetUsersAndCurrencyRate() ([]int64, *fetcher.Rate, error) {
	var rate *fetcher.Rate
	rate, err := ns.rateFetch.FetchRate()
	if err != nil {
		ns.logger.Printf("NotifyService: GetUsersAndCurrencyRate: FetchRate %v\n", err)
		return nil, rate, domain.ErrGeneric
	}

	users, err := ns.userRepository.GetAllUsers()
	if err != nil {
		ns.logger.Printf("NotifyService: GetUsersAndCurrencyRate: GetAllUsers %v\n", err)
		return nil, rate, domain.ErrGeneric
	}

	var ids []int64
	for _, user := range users {
		ids = append(ids, user.ChatID)
	}

	if len(ids) == 0 {
		ns.logger.Println("NotifyService: GetUsersAndCurrencyRate: No subscribers found")
		return nil, rate, nil
	}

	return ids, rate, nil
}
