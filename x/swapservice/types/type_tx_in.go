package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	common "gitlab.com/thorchain/bepswap/common"
)

type status string
type TxInIndex []common.TxID

const (
	Incomplete status = "incomplete"
	Done       status = "done"
	Reverted   status = "reverted"
)

// Meant to track if we have processed a specific binance tx
type TxIn struct {
	Status  status            `json:"status"`
	Done    common.TxID       `json:"txhash"` // completed binance chain tx hash. This is a slice to track if we've "double spent" an input
	Memo    string            `json:"memo"`   // memo
	Coins   common.Coins      `json:"coins"`  // coins sent in tx
	Sender  common.BnbAddress `json:"sender"`
	Signers []sdk.AccAddress  `json:"signers"` // trust accounts saw this tx
}

func NewTxIn(coins common.Coins, memo string, sender common.BnbAddress) TxIn {
	return TxIn{
		Coins:  coins,
		Memo:   memo,
		Sender: sender,
		Status: Incomplete,
	}
}

func (tx TxIn) Valid() error {
	if tx.Sender.IsEmpty() {
		return fmt.Errorf("Sender cannot be empty")
	}
	if len(tx.Coins) == 0 {
		return fmt.Errorf("Coins cannot be empty")
	}
	if len(tx.Memo) == 0 {
		return fmt.Errorf("Memo cannot be empty")
	}

	return nil
}

func (tx TxIn) Empty() bool {
	return tx.Sender.IsEmpty()
}

func (tx1 TxIn) Equals(tx2 TxIn) bool {
	if tx1.Memo != tx2.Memo {
		return false
	}
	if !tx1.Sender.Equals(tx2.Sender) {
		return false
	}
	if len(tx1.Coins) != len(tx2.Coins) {
		return false
	}
	for i, _ := range tx1.Coins {
		if !tx1.Coins[i].Denom.Equals(tx2.Coins[i].Denom) {
			return false
		}
		if !tx1.Coins[i].Amount.Equals(tx2.Coins[i].Amount) {
			return false
		}
	}

	return true
}

func (tx TxIn) String() string {
	return tx.Done.String()
}

func (tx *TxIn) Sign(signer sdk.AccAddress) {
	for _, sign := range tx.Signers {
		if sign.Equals(signer) {
			return // do nothing
		}
	}
	tx.Signers = append(tx.Signers, signer)
}

func (tx *TxIn) SetDone(hash common.TxID) {
	tx.Status = Done
	tx.Done = hash
}

func (tx *TxIn) SetReverted(hash common.TxID) {
	tx.Status = Reverted
	tx.Done = hash
}

type TxInVoter struct {
	TxID        common.TxID `json:"tx_id"`
	Txs         []TxIn      `json:"txs"`
	IsProcessed bool        `json:"is_processed"`
}

func NewTxInVoter(txID common.TxID, txs []TxIn) TxInVoter {
	return TxInVoter{
		TxID: txID,
		Txs:  txs,
	}
}

func (tx TxInVoter) Key() common.TxID {
	return tx.TxID
}

func (tx TxInVoter) String() string {
	return tx.TxID.String()
}

func (tx *TxInVoter) SetDone(hash common.TxID) {
	for i, _ := range tx.Txs {
		tx.Txs[i].SetDone(hash)
	}
}

func (tx *TxInVoter) Add(txIn TxIn, signer sdk.AccAddress) {
	// check if this signer has already signed, no take backs allowed
	for _, transaction := range tx.Txs {
		for _, siggy := range transaction.Signers {
			if siggy.Equals(signer) {
				return
			}
		}
	}

	for i, _ := range tx.Txs {
		if tx.Txs[i].Equals(txIn) {
			tx.Txs[i].Sign(signer)
			return
		}
	}

	txIn.Sign(signer)
	tx.Txs = append(tx.Txs, txIn)
}

func (tx *TxInVoter) Adds(txs []TxIn, signer sdk.AccAddress) {
	for _, txIn := range txs {
		tx.Add(txIn, signer)
	}
}

func (tx TxInVoter) HasConensus(totalTrusted int) bool {
	for _, txIn := range tx.Txs {
		if HasMajority(len(txIn.Signers), totalTrusted) {
			return true
		}
	}

	return false
}

func (tx TxInVoter) GetTx(totalTrusted int) TxIn {
	for _, txIn := range tx.Txs {
		if HasMajority(len(txIn.Signers), totalTrusted) {
			return txIn
		}
	}

	return TxIn{}
}
