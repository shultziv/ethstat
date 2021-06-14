package app

import (
	"github.com/shultziv/ethstat/internal/config"
	"github.com/shultziv/ethstat/internal/delivery/cmd"
	"github.com/shultziv/ethstat/internal/repo/etherscan"
	"github.com/shultziv/ethstat/internal/service/ethstat"
)

func EthStatCmdRun(etherScanConfig *config.EtherScan) (err error) {
	proxyURLs, err := etherScanConfig.GetProxyURLs()
	if err != nil {
		return err
	}

	etherScan := etherscan.New(proxyURLs...)
	ethStat := ethstat.New(etherScan)

	cmdHandler := cmd.New(ethStat)
	cmdHandler.Run()
	return nil
}
