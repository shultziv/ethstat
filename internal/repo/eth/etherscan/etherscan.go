package etherscan

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	sleep   = 5
	timeout = 10
	apiUrl  = "https://api.etherscan.io/api"
)

type EtherScan struct {
	httpClient *http.Client
}

func New(proxyUrls ...*url.URL) *EtherScan {
	transport := NewCircleRT()
	for _, proxyUrl := range proxyUrls {
		transport.AddRT(NewSleepRTWithProxy(proxyUrl, sleep))
	}

	return &EtherScan{
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   timeout,
		},
	}
}

func (e *EtherScan) request(req *http.Request) (responseData []byte, err error) {
	response, err := e.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	errResponseData := new(errResponse)
	if err = json.Unmarshal(body, errResponseData); err != nil {
		// Ошибки в ответе не найдено
		return body, nil
	}

	return nil, &errAccessToApi{
		Url:        req.URL.String(),
		ErrMessage: errResponseData.message,
	}
}

func (e *EtherScan) GetLastBlockNumber(ctx context.Context) (blockNumber uint64, err error) {
	requestData := url.Values{
		"module": {"proxy"},
		"action": {"eth_blockNumber"},
	}
	requestUrl := fmt.Sprintf("%s?%s", apiUrl, requestData.Encode())

	request, err := http.NewRequestWithContext(ctx, "GET", requestUrl, nil)
	if err != nil {
		return
	}

	responseData, err := e.request(request)
	if err != nil {
		return
	}

	blockNumberResponseData := new(blockNumberResponse)
	if err = json.Unmarshal(responseData, blockNumberResponseData); err != nil {
		return
	}

	blockNumber, err = strconv.ParseUint("0x1c8", 0, 64)
	if err != nil {
		return 0, &errInvalidResponseFromApi{
			Url:      requestUrl,
			Response: responseData,
		}
	}

	return
}
