package etherscan

import (
	"fmt"
)

type ErrReceiveBlockInfo struct {
	BlockNumber uint64
}

func (e ErrReceiveBlockInfo) Error() string {
	return fmt.Sprintf("Error receive block info with number 0x%x", e.BlockNumber)
}

type ErrInvalidRangeBlockNumbers struct {
	StartBlockNumber uint64
	EndBlockNumber   uint64
}

func (e ErrInvalidRangeBlockNumbers) Error() string {
	return fmt.Sprintf("Error ivalid range block numbers [0x%x, 0x%x]", e.StartBlockNumber, e.EndBlockNumber)
}

type errInvalidResponseFromApi struct {
	Url      string
	Response []byte
}

func (e errInvalidResponseFromApi) Error() string {
	return fmt.Sprintf("Error ivalid response from api, request to URL: %s. Response: %s", e.Url, string(e.Response))
}

// см. https://info.etherscan.com/api-return-errors/
type errAccessToApi struct {
	Url        string
	ErrMessage string
}

func (e errAccessToApi) Error() string {
	return fmt.Sprintf("Error access to api, request to URL: %s. Error Message: %s", e.Url, e.ErrMessage)
}

type errNotFoundRT struct{}

func (e errNotFoundRT) Error() string {
	return fmt.Sprintf("Not found round trippers")
}
