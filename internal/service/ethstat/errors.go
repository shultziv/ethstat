package ethstat

import (
	"errors"
	"fmt"
)

var (
	ErrReceiveInfoFromRepo  = errors.New("Error receive info from repo")
	ErrNotFoundTransactions = errors.New("Not found transactions")
)

type ErrInvalidCountLastBlocks struct {
	CountLastBlocks uint64
}

func (e ErrInvalidCountLastBlocks) Error() string {
	return fmt.Sprintf("Ivalid value count last blocks (%d)", e.CountLastBlocks)
}
