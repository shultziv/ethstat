package cmd

import (
	"context"
	"flag"
	"fmt"
	"github.com/shultziv/ethstat/internal/service/ethstat"
)

type Handler struct {
	ethStat *ethstat.EthStat
}

func New(ethStat *ethstat.EthStat) *Handler {
	return &Handler{
		ethStat: ethStat,
	}
}

func (h *Handler) Run() {
	var lastCountBlocks int
	flag.IntVar(&lastCountBlocks, "c", 100, "Count last blocks")

	flag.Parse()

	ctx := context.Background()
	addrWithMaxChange, err := h.ethStat.GetAddrBiggestBalanceChange(ctx, uint64(lastCountBlocks))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(addrWithMaxChange)
}
