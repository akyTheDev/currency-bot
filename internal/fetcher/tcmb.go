package fetcher

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type TCMBClient struct {
	client *http.Client
	url    string
}

type tcmbDate struct {
	XMLName    xml.Name   `xml:"Tarih_Date"`
	Currencies []currency `xml:"Currency"`
}

type currency struct {
	Unit         int     `xml:"Unit"`
	CurrencyName string  `xml:"CurrencyName"`
	ForexBuying  float64 `xml:"ForexBuying"`
	ForexSelling float64 `xml:"ForexSelling"`
}

const TcmbUrl = "https://www.tcmb.gov.tr/kurlar/today.xml"

func NewTCMBClient(url string, timeoutSeconds int) *TCMBClient {
	return &TCMBClient{
		url: url,
		client: &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
	}
}

func (c *TCMBClient) FetchRate() (*Rate, error) {
	var result Rate

	resp, err := c.client.Get(c.url)
	if err != nil {
		return &result, fmt.Errorf("http GET: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &result, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &result, fmt.Errorf("read body: %w", err)
	}

	var parsedBody tcmbDate
	if err := xml.Unmarshal(body, &parsedBody); err != nil {
		return &result, fmt.Errorf("parse XML: %w", err)
	}

	for _, cur := range parsedBody.Currencies {
		if cur.CurrencyName == "EURO" {
			result.Buying = cur.ForexBuying
			result.Selling = cur.ForexSelling
			return &result, nil
		}
	}

	return &result, errors.New("EUR not found")
}
