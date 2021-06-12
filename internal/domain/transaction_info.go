package domain

type TransactionInfo struct {
	FromAddress string `json:"from"`
	ToAddress   string `json:"to"`
	Value       string `json:"value"`
}
