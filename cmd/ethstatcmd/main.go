package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/shultziv/ethstat/internal/app"
	"github.com/shultziv/ethstat/internal/config"
)

func main() {
	etherScanConf := new(config.EtherScan)
	if err := env.Parse(etherScanConf); err != nil {
		fmt.Println(err)
		return
	}

	if err := app.EthStatCmdRun(etherScanConf); err != nil {
		fmt.Println(err)
		return
	}
}
