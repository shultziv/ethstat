package ethstat

import (
	"context"
	"github.com/shultziv/ethstat/internal/domain"
)

type ethRepoMock struct {
	LastBlockNumber   uint64
	BlockNumberToInfo map[uint64]*domain.BlockInfo
	BlocksInfoInRange []*domain.BlockInfo
	Err               error
}

func (e *ethRepoMock) GetLastBlockNumber(ctx context.Context) (blockNumber uint64, err error) {
	return e.LastBlockNumber, e.Err
}

func (e *ethRepoMock) GetBlockInfo(ctx context.Context, blockNumber uint64) (blockInfo *domain.BlockInfo, err error) {
	return e.BlockNumberToInfo[blockNumber], err
}

func (e *ethRepoMock) GetBlocksInfoInRange(ctx context.Context, startBlockNumber uint64, endBlockNumber uint64) (blocksInfo []*domain.BlockInfo, err error) {
	return e.BlocksInfoInRange, e.Err
}
