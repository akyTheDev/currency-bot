package fetcher

type Rate struct {
	Selling float64
	Buying  float64
}

type RateFetcher interface {
	FetchRate() (*Rate, error)
}
