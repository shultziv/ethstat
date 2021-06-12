package ethstat

import (
	"context"
	"github.com/shultziv/ethstat/internal/domain"
	"testing"
)

func TestEthStat_GetAddrBiggestBalanceChange(t *testing.T) {
	expectedAddr := "0x1"
	eRepo := &ethRepoMock{
		LastBlockNumber: 255,
		BlocksInfoInRange: []*domain.BlockInfo{
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
		},
		Err: nil,
	}

	eStat := New(eRepo)

	ctx := context.Background()
	addr, err := eStat.GetAddrBiggestBalanceChange(ctx, 2)
	if err != nil {
		t.Error(err)
		return
	}

	if expectedAddr != addr {
		t.Errorf("Expected addr: %s, but result: %s", expectedAddr, addr)
	}
}
