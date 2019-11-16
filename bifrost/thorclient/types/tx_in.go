package types

import "gitlab.com/thorchain/bepswap/thornode/common"

type TxIn struct {
	BlockHeight string       `json:"blockHeight"`
	Count       string       `json:"count"`
	Chain       common.Chain `json:"chain"`
	TxArray     []TxInItem   `json:"txArray"`
}

type TxInItem struct {
	Tx                  string       `json:"tx"`
	Memo                string       `json:"MEMO"`
	Sender              string       `json:"sender"`
	To                  string       `json:"to"` // to adddress
	Coins               common.Coins `json:"coins"`
	ObservedPoolAddress string       `json:"observed_pool_address"`
}
type TxInStatus byte

const (
	Processing TxInStatus = iota
	Failed
)