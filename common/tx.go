package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

type (
	TxID  string
	TxIDs []TxID
)

var BlankTxID = TxID("0000000000000000000000000000000000000000000000000000000000000000")

func NewTxID(hash string) (TxID, error) {
	switch len(hash) {
	case 64:
		// do nothing
	case 66: // ETH check
		if !strings.HasPrefix(hash, "0x") {
			err := fmt.Errorf("TxID Error: Must be 66 characters (got %d)", len(hash))
			return TxID(""), err
		}
	default:
		err := fmt.Errorf("TxID Error: Must be 64 characters (got %d)", len(hash))
		return TxID(""), err
	}

	return TxID(strings.ToUpper(hash)), nil
}

func (tx TxID) Equals(tx2 TxID) bool {
	return strings.EqualFold(tx.String(), tx2.String())
}

func (tx TxID) IsEmpty() bool {
	return strings.TrimSpace(tx.String()) == ""
}

func (tx TxID) String() string {
	return string(tx)
}

type Tx struct {
	ID          TxID    `json:"id"`
	Chain       Chain   `json:"chain"`
	FromAddress Address `json:"from_address"`
	ToAddress   Address `json:"to_address"`
	Coins       Coins   `json:"coins"`
	Gas         Gas     `json:"gas"`
	Memo        string  `json:"memo"`
}

type Txs []Tx

func GetRagnarokTx(chain Chain) Tx {
	return Tx{
		Chain:       chain,
		ID:          BlankTxID,
		FromAddress: RagnarokAddr,
		ToAddress:   RagnarokAddr,
		Coins: Coins{
			// used for ragnarok, so doesn't really matter
			NewCoin(BNBAsset, sdk.OneUint()),
		},
		Gas: Gas{
			// used for ragnarok, so doesn't really matter
			NewCoin(BNBAsset, sdk.OneUint()),
		},
		Memo: "Ragnarok",
	}
}

func NewTx(txID TxID, from Address, to Address, coins Coins, gas Gas, memo string) Tx {
	var chain Chain
	for _, coin := range coins {
		chain = coin.Asset.Chain
		break
	}
	return Tx{
		ID:          txID,
		Chain:       chain,
		FromAddress: from,
		ToAddress:   to,
		Coins:       coins,
		Gas:         gas,
		Memo:        memo,
	}
}

func (tx Tx) String() string {
	return fmt.Sprintf("%s: %s ==> %s (Memo: %s) %s", tx.ID, tx.FromAddress, tx.ToAddress, tx.Memo, tx.Coins)
}

func (tx Tx) Bytes() ([]byte, error) {
	buf, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}
	return sdk.SortJSON(buf)
}

func (tx Tx) Sign(priv crypto.PrivKey) ([]byte, error) {
	buf, err := tx.Bytes()
	if err != nil {
		return buf, err
	}
	return priv.Sign(buf)
}

func (tx Tx) Verify(pk crypto.PubKey, sig []byte) (bool, error) {
	buf, err := tx.Bytes()
	if err != nil {
		return false, err
	}
	return pk.VerifyBytes(buf, sig), nil
}

func (tx Tx) IsEmpty() bool {
	return tx.ID.IsEmpty()
}

func (tx1 Tx) Equals(tx2 Tx) bool {
	if !tx1.ID.Equals(tx2.ID) {
		return false
	}
	if !tx1.Chain.Equals(tx2.Chain) {
		return false
	}
	if !tx1.FromAddress.Equals(tx2.FromAddress) {
		return false
	}
	if !tx1.ToAddress.Equals(tx2.ToAddress) {
		return false
	}
	if !tx1.Coins.Equals(tx2.Coins) {
		return false
	}
	if !tx1.Gas.Equals(tx2.Gas) {
		return false
	}
	if !strings.EqualFold(tx1.Memo, tx2.Memo) {
		return false
	}
	return true
}

func (tx Tx) IsValid() error {
	if tx.ID.IsEmpty() {
		return errors.New("Tx ID cannot be empty")
	}
	if tx.FromAddress.IsEmpty() {
		return errors.New("From address cannot be empty")
	}
	if tx.ToAddress.IsEmpty() {
		return errors.New("To address cannot be empty")
	}
	if tx.Chain.IsEmpty() {
		return errors.New("Chain cannot be empty")
	}
	if len(tx.Coins) == 0 {
		return errors.New("Must have at least 1 coin")
	}
	if err := tx.Coins.IsValid(); err != nil {
		return err
	}
	if len(tx.Gas) == 0 {
		return errors.New("Must have at least 1 gas coin")
	}
	if err := tx.Gas.IsValid(); err != nil {
		return err
	}
	return nil
}
