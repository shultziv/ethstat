package domain

type BlockInfo struct {
	Number       string             `json:"number"`
	Transactions []*TransactionInfo `json:"transactions"`
}
