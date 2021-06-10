package etherscan

import (
	"net/http"
	"net/url"
	"time"
)

type SleepRoundTripper struct {
	ticker    *time.Ticker
	transport *http.Transport
}

func NewSleepRTWithProxy(proxy *url.URL, sleep time.Duration) *SleepRoundTripper {
	return &SleepRoundTripper{
		ticker: time.NewTicker(sleep),
		transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}
}

func (srt *SleepRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	<-srt.ticker.C
	return srt.transport.RoundTrip(req)
}
