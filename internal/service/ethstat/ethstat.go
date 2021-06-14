package ethstat

import (
	"context"
	"github.com/shultziv/ethstat/internal/domain"
	"math"
	"strconv"
)

//go:generate mockgen -destination=mocks/mock_eth_repo.go -package=mocks . EthRepo
type EthRepo interface {
	GetLastBlockNumber(ctx context.Context) (blockNumber uint64, err error)
	GetBlockInfo(ctx context.Context, blockNumber uint64) (blockInfo *domain.BlockInfo, err error)
	GetBlocksInfoInRange(ctx context.Context, startBlockNumber uint64, endBlockNumber uint64) (blocksInfo []*domain.BlockInfo, err error)
}

type EthStat struct {
	eRepo EthRepo
}

func New(eRepo EthRepo) *EthStat {
	return &EthStat{
		eRepo: eRepo,
	}
}

func (es *EthStat) GetBalanceChanges(ctx context.Context, countLastBlocks uint64) (addrToChange map[string]int, err error) {
	lastBlockNumber, err := es.eRepo.GetLastBlockNumber(ctx)
	if err != nil {
		return nil, ErrReceiveInfoFromRepo
	}

	var startBlockNumber uint64
	if lastBlockNumber-countLastBlocks < 0 {
		return nil, ErrInvalidCountLastBlocks{
			CountLastBlocks: countLastBlocks,
		}
	}

	startBlockNumber = lastBlockNumber - countLastBlocks

	blocksInfo, err := es.eRepo.GetBlocksInfoInRange(ctx, startBlockNumber, lastBlockNumber)
	if err != nil {
		return nil, ErrReceiveInfoFromRepo
	}

	addrToChange = make(map[string]int)

	for _, blockInfo := range blocksInfo {

		for _, transaction := range blockInfo.Transactions {
			value, err := strconv.ParseUint(transaction.Value, 0, 64)
			if err != nil {
				// log
			}

			addrToChange[transaction.FromAddress] -= int(value)
			addrToChange[transaction.ToAddress] += int(value)
		}
	}

	return addrToChange, nil
}

func (es *EthStat) GetAddrBiggestBalanceChange(ctx context.Context, countLastBlocks uint64) (addrWithMaxChange string, err error) {
	addrToChange, err := es.GetBalanceChanges(ctx, countLastBlocks)
	if err != nil {
		return "", err
	}

	if len(addrToChange) == 0 {
		return "", ErrNotFoundTransactions
	}

	for addr, change := range addrToChange {
		if math.Abs(float64(change)) > math.Abs(float64(addrToChange[addrWithMaxChange])) {
			addrWithMaxChange = addr
		}
	}

	return addrWithMaxChange, nil
}
