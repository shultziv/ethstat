package etherscan

import "github.com/shultziv/ethstat/internal/domain"

type blockInfoResponse struct {
	Id     uint8             `json:"id"`
	Result *domain.BlockInfo `json:"result"`
}

type blockNumberResponse struct {
	Id     uint8  `json:"id"`
	Number string `json:"result"`
}

// Описание здесь: https://info.etherscan.com/api-return-errors/
type errResponse struct {
	Status  uint8  `json:"status"`
	Message string `json:"message"`
}
