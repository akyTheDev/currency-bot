package fetcher

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestFetchRate(t *testing.T) {
	tests := []struct {
		name            string
		xmlPayload      string
		statusCode      int
		timeoutSec      int
		expectedSelling float64
		expectedBuying  float64
		expectErrSub    string
		delayResponse   time.Duration
	}{
		{
			name: "ValidXML",
			xmlPayload: `
<Tarih_Date>
  <Currency>
	<CurrencyName>EURO</CurrencyName>
    <ForexSelling>22.2222</ForexSelling>
    <ForexBuying>21.2222</ForexBuying>
  </Currency>
  <Currency>
	<CurrencyName>USD</CurrencyName>
    <ForexSelling>20.1000</ForexSelling>
    <ForexBuying>19.1000</ForexBuying>
  </Currency>
</Tarih_Date>`,
			statusCode:      http.StatusOK,
			timeoutSec:      2,
			expectedSelling: 22.2222,
			expectedBuying:  21.2222,
		},
		{
			name:         "Non200Status",
			xmlPayload:   "",
			statusCode:   http.StatusServiceUnavailable,
			timeoutSec:   2,
			expectErrSub: "unexpected status",
		},
		{
			name: "MalformedXML",
			xmlPayload: `
<Tarih_Date>
  <Currency Kod="EUR"><ForexSelling>21.3333</ForexSelling>
  <!-- missing closing tags -->
</Tarih_Date>`,
			statusCode:   http.StatusOK,
			timeoutSec:   2,
			expectErrSub: "parse XML",
		},
		{
			name: "MissingEUR",
			xmlPayload: `
<Tarih_Date>
  <Currency Kod="USD">
    <ForexSelling>20.0000</ForexSelling>
  </Currency>
</Tarih_Date>`,
			statusCode:   http.StatusOK,
			timeoutSec:   2,
			expectErrSub: "EUR not found",
		},
		{
			name: "ClientTimeout",
			xmlPayload: `
<Tarih_Date>
  <Currency Kod="EUR">
    <ForexSelling>30.0000</ForexSelling>
  </Currency>
</Tarih_Date>`,
			statusCode:    http.StatusOK,
			timeoutSec:    1,
			expectErrSub:  "http GET",
			delayResponse: 2 * time.Second,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(tc.delayResponse)
				w.WriteHeader(tc.statusCode)
				if tc.statusCode == http.StatusOK {
					w.Header().Set("Content-Type", "application/xml")
					fmt.Fprint(w, tc.xmlPayload)
				}
			}))
			defer ts.Close()

			client := NewTCMBClient(
				ts.URL,
				tc.timeoutSec,
			)

			rate, err := client.FetchRate()

			if tc.expectErrSub == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if rate.Selling != tc.expectedSelling || rate.Buying != tc.expectedBuying {
					t.Errorf("rate = %v; want %v", rate, Rate{Buying: tc.expectedBuying, Selling: tc.expectedSelling})
				}
			} else {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tc.expectErrSub) {
					t.Errorf("error %q does not contain %q", err.Error(), tc.expectErrSub)
				}
			}
		})
	}
}
