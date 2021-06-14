package ethstat

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/shultziv/ethstat/internal/domain"
	"github.com/shultziv/ethstat/internal/service/ethstat/mocks"
	"testing"
)

func TestEthStat_GetAddrBiggestBalanceChange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEthRepo := mocks.NewMockEthRepo(ctrl)

	var lastBlockNumber uint64 = 255
	var countLastBlocks uint64 = 1
	blocksInfoInRange := []*domain.BlockInfo{
		{
			Number: "0x1",
			Transactions: []*domain.TransactionInfo{
				{
					FromAddress: "0x1",
					ToAddress:   "0x2",
					Value:       "0x1",
				},
				{
					FromAddress: "0x1",
					ToAddress:   "0x3",
					Value:       "0x2",
				},
				{
					FromAddress: "0x3",
					ToAddress:   "0x2",
					Value:       "0x1",
				},
			},
		},
	}
	expectedAddr := "0x1"

	mockEthRepo.
		EXPECT().
		GetLastBlockNumber(gomock.Any()).
		Return(lastBlockNumber, nil)

	mockEthRepo.
		EXPECT().
		GetBlocksInfoInRange(gomock.Any(), lastBlockNumber-countLastBlocks, lastBlockNumber).
		Return(blocksInfoInRange, nil)

	eStat := New(mockEthRepo)

	ctx := context.Background()
	addr, err := eStat.GetAddrBiggestBalanceChange(ctx, countLastBlocks)
	if err != nil {
		t.Error(err)
		return
	}

	if expectedAddr != addr {
		t.Errorf("Expected addr: %s, but result: %s", expectedAddr, addr)
	}
}
