package etherscan

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/shultziv/ethstat/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEtherScan_GetLastBlockNumber(t *testing.T) {
	var expectedBlockNumber uint64 = 25
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		blockNumberData := blockNumberResponse{
			Id:     1,
			Number: fmt.Sprintf("0x%x", expectedBlockNumber),
		}

		responseData, err := json.Marshal(blockNumberData)
		if err != nil {
			t.Error(err)
			return
		}

		_, err = w.Write(responseData)
		if err != nil {
			t.Error(err)
			return
		}
	}))
	defer ts.Close()

	apiUrl = ts.URL
	sleep = 1

	etherScan := New()
	ctx := context.Background()
	blockNumber, err := etherScan.GetLastBlockNumber(ctx)
	if err != nil {
		t.Error(err)
		return
	}

	if expectedBlockNumber != blockNumber {
		t.Errorf("Expected %d, but result %d", expectedBlockNumber, blockNumber)
	}
}

func TestEtherScan_GetBlockInfo(t *testing.T) {
	expectedBlockInfo := &domain.BlockInfo{
		Number:       "0x1",
		Transactions: nil,
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		blockInfoResponse := &blockInfoResponse{
			Id:     0,
			Result: expectedBlockInfo,
		}

		responseData, err := json.Marshal(blockInfoResponse)
		if err != nil {
			t.Error(err)
			return
		}

		_, err = w.Write(responseData)
		if err != nil {
			t.Error(err)
			return
		}
	}))
	defer ts.Close()

	apiUrl = ts.URL
	sleep = 1

	etherScan := New()
	ctx := context.Background()
	blockInfo, err := etherScan.GetBlockInfo(ctx, 1)
	if err != nil {
		t.Error(err)
		return
	}

	if blockInfo.Number != expectedBlockInfo.Number {
		t.Errorf("Expected %v, but result %v", expectedBlockInfo, blockInfo)
	}
}

func TestEtherScan_GetBlocksInfoInRange(t *testing.T) {
	startBlockNumber := 1
	endBlockNumber := 14

	countBlocks := endBlockNumber - startBlockNumber
	expectedBlocksInfo := make([]*domain.BlockInfo, countBlocks, countBlocks)
	for i := 0; i < countBlocks; i++ {
		expectedBlocksInfo[i] = &domain.BlockInfo{
			Number:       fmt.Sprintf("0x%x", startBlockNumber+i),
			Transactions: nil,
		}
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		blocksNumber, ok := r.URL.Query()["tag"]
		if !ok || len(blocksNumber) == 0 {
			t.Errorf("Invalid request to API: %s", r.URL.String())
			return
		}

		blockNumber := blocksNumber[0]

		var response *blockInfoResponse
		for index, expectedBlockInfo := range expectedBlocksInfo {
			if expectedBlockInfo.Number == blockNumber {
				response = &blockInfoResponse{
					Id:     uint8(index),
					Result: expectedBlockInfo,
				}
			}
		}

		if response == nil {
			t.Errorf("Block with number %s not found.", blockNumber)
			return
		}

		responseData, err := json.Marshal(response)
		if err != nil {
			t.Error(err)
			return
		}

		_, err = w.Write(responseData)
		if err != nil {
			t.Error(err)
			return
		}
	}))
	defer ts.Close()

	apiUrl = ts.URL
	sleep = 1

	etherScan := New()
	ctx := context.Background()
	blocksInfo, err := etherScan.GetBlocksInfoInRange(ctx, uint64(startBlockNumber), uint64(endBlockNumber))
	if err != nil {
		t.Error(err)
		return
	}

	if len(blocksInfo) != len(expectedBlocksInfo) {
		t.Errorf("Expected blocks info: %v, but result: %v", expectedBlocksInfo, blocksInfo)
		return
	}

	for i := range blocksInfo {
		if blocksInfo[i].Number != expectedBlocksInfo[i].Number {
			t.Errorf("Expected block info: %v, but result: %v", expectedBlocksInfo[i], blocksInfo[i])
			return
		}
	}
}
