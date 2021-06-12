package etherscan

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/shultziv/ethstat/internal/domain"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

var (
	sleep   = 10
	timeout = 10
	apiUrl  = "https://api.etherscan.io/api"
)

type EtherScan struct {
	hc httpClient
}

func New(proxyUrls ...*url.URL) *EtherScan {
	circleHC := NewCircleHC()
	for _, proxyUrl := range proxyUrls {
		circleHC.AddHC(NewSleepHCWithProxy(proxyUrl, time.Duration(sleep)*time.Second, time.Duration(timeout)*time.Second))
	}

	if len(proxyUrls) == 0 {
		circleHC.AddHC(NewSleepHC(time.Duration(sleep)*time.Second, time.Duration(timeout)*time.Second))
	}

	return &EtherScan{
		hc: circleHC,
	}
}

func (e *EtherScan) request(req *http.Request) (responseData []byte, err error) {
	response, err := e.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	errResponseData := new(errResponse)
	err = json.Unmarshal(body, errResponseData)
	if err != nil {
		// Ответ json формата, но это и не ошибка
		return body, nil
	}

	if errResponseData.Status == 0 && errResponseData.Message != "" {
		return nil, &errAccessToApi{
			Url:        req.URL.String(),
			ErrMessage: errResponseData.Message,
		}
	}

	return body, nil
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

	blockNumber, err = strconv.ParseUint(blockNumberResponseData.Number, 0, 64)
	if err != nil {
		return 0, &errInvalidResponseFromApi{
			Url:      requestUrl,
			Response: responseData,
		}
	}

	return
}

func (e *EtherScan) GetBlockInfo(ctx context.Context, blockNumber uint64) (blockInfo *domain.BlockInfo, err error) {
	requestData := url.Values{
		"module":  {"proxy"},
		"action":  {"eth_getBlockByNumber"},
		"tag":     {fmt.Sprintf("0x%x", blockNumber)},
		"boolean": {"true"},
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

	blockInfoResponseData := new(blockInfoResponse)
	if err = json.Unmarshal(responseData, blockInfoResponseData); err != nil {
		return nil, &errInvalidResponseFromApi{
			Url:      requestUrl,
			Response: responseData,
		}
	}

	if blockInfoResponseData.Result == nil {
		return nil, &errInvalidResponseFromApi{
			Url:      requestUrl,
			Response: responseData,
		}
	}

	return blockInfoResponseData.Result, nil
}

func (e *EtherScan) GetBlocksInfoInRange(ctx context.Context, startBlockNumber uint64, endBlockNumber uint64) (blocksInfo []*domain.BlockInfo, err error) {
	if startBlockNumber >= endBlockNumber {
		return nil, &ErrInvalidRangeBlockNumbers{
			StartBlockNumber: startBlockNumber,
			EndBlockNumber:   endBlockNumber,
		}
	}

	countBlock := int(endBlockNumber - startBlockNumber)
	blocksInfo = make([]*domain.BlockInfo, countBlock, countBlock)

	wg := sync.WaitGroup{}
	for i := 0; i < countBlock; i++ {
		wg.Add(1)
		go func(blockNumber uint64, blockIndex int) {
			blockInfo, err := e.GetBlockInfo(ctx, blockNumber)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
			blocksInfo[blockIndex] = blockInfo
			wg.Done()
		}(startBlockNumber+uint64(i), i)
	}

	wg.Wait()

	for blockIndex, blockInfo := range blocksInfo {
		if blockInfo == nil {
			return nil, &ErrReceiveBlockInfo{
				BlockNumber: startBlockNumber + uint64(blockIndex),
			}
		}
	}

	return blocksInfo, nil
}
