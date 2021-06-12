package app

import (
	"context"
	"fmt"
	"github.com/shultziv/ethstat/internal/config"
	"github.com/shultziv/ethstat/internal/repo/etherscan"
	"github.com/shultziv/ethstat/internal/service/ethstat"
)

const lastCountBlocks = 100

func EthStatCmdRun(etherScanConfig *config.EtherScan) (err error) {
	proxyURLs, err := etherScanConfig.GetProxyURLs()
	if err != nil {
		return err
	}

	etherScan := etherscan.New(proxyURLs...)
	ethStat := ethstat.New(etherScan)

	ctx := context.Background()
	addrWithMaxChange, err := ethStat.GetAddrBiggestBalanceChange(ctx, lastCountBlocks)
	if err != nil {
		return err
	}

	fmt.Println(addrWithMaxChange)
	return nil
}
