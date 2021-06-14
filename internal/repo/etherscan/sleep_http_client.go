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

	lastTimeReqCh chan time.Time // Буферезированный канала с емкостью 1, туда кладется время последнего запроса
	sleep         time.Duration  // Кол-во времени которое необходимо подождать псле предидущего запроса
}

func NewSleepHCWithProxy(proxy *url.URL, sleep time.Duration, timeout time.Duration) *SleepHttpClient {
	lastTimeReqCh := make(chan time.Time, 1)
	lastTimeReqCh <- time.Now().Add(-sleep)

	return &SleepHttpClient{
		lastTimeReqCh: lastTimeReqCh,
		sleep:         sleep,
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
func (srt *SleepHttpClient) Do(req *http.Request) (response *http.Response, err error) {
	lastTimeReq := <-srt.lastTimeReqCh
	sleep := srt.sleep - time.Now().Sub(lastTimeReq) // Кол-во времени прошедшего после последнего запроса
	time.Sleep(sleep)

	response, err = srt.client.Do(req)
	srt.lastTimeReqCh <- time.Now()
	return
}
