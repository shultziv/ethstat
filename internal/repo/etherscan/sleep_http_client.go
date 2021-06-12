package etherscan

import (
	"net/http"
	"net/url"
	"time"
)

// Http Client с паузой перед запросом, с поддержкой прокси
// Потокобезопасный
type SleepHttpClient struct {
	client *http.Client
	ticker *time.Ticker
}

func NewSleepHCWithProxy(proxy *url.URL, sleep time.Duration, timeout time.Duration) *SleepHttpClient {
	return &SleepHttpClient{
		ticker: time.NewTicker(sleep),
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
			Timeout: timeout,
		},
	}
}

func NewSleepHC(sleep time.Duration, timeout time.Duration) *SleepHttpClient {
	return NewSleepHCWithProxy(nil, sleep, timeout)
}

// Перед вызовом запросом выдерживает паузу
func (srt *SleepHttpClient) Do(req *http.Request) (*http.Response, error) {
	<-srt.ticker.C
	return srt.client.Do(req)
}
