package etherscan

type blockNumberResponse struct {
	id     uint8  `json:"id"`
	number string `json:"result"`
}

// Описание здесь: https://info.etherscan.com/api-return-errors/
type errResponse struct {
	status  uint8  `json:"status"`
	message string `json:"message"`
}
