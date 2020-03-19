package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"gitlab.com/thorchain/thornode/common"
	. "gopkg.in/check.v1"
)

type EventSuite struct{}

var _ = Suite(&EventSuite{})

func (s EventSuite) TestSwapEvent(c *C) {
	evt := NewEventSwap(
		common.BNBAsset,
		sdk.NewUint(5),
		sdk.NewUint(5),
		sdk.NewUint(5),
	)
	c.Check(evt.Type(), Equals, "swap")
}

func (s EventSuite) TestStakeEvent(c *C) {
	evt := NewEventStake(
		common.BNBAsset,
		sdk.NewUint(5),
	)
	c.Check(evt.Type(), Equals, "stake")
}

func (s EventSuite) TestUnstakeEvent(c *C) {
	evt := NewEventUnstake(
		common.BNBAsset,
		sdk.NewUint(6),
		5000,
		sdk.NewDec(0),
	)
	c.Check(evt.Type(), Equals, "unstake")
}

func (s EventSuite) TestPool(c *C) {
	evt := NewEventPool(common.BNBAsset, Enabled)
	c.Check(evt.Type(), Equals, "pool")
	c.Check(evt.Pool.String(), Equals, common.BNBAsset.String())
	c.Check(evt.Status.String(), Equals, Enabled.String())
}

func (s EventSuite) TestReward(c *C) {
	evt := NewEventRewards(sdk.NewUint(300), []PoolAmt{
		{common.BNBAsset, 30},
		{common.BTCAsset, 40},
	})
	c.Check(evt.Type(), Equals, "rewards")
	c.Check(evt.BondReward.String(), Equals, "300")
	c.Assert(evt.PoolRewards, HasLen, 2)
	c.Check(evt.PoolRewards[0].Asset.Equals(common.BNBAsset), Equals, true)
	c.Check(evt.PoolRewards[0].Amount, Equals, int64(30))
	c.Check(evt.PoolRewards[1].Asset.Equals(common.BTCAsset), Equals, true)
	c.Check(evt.PoolRewards[1].Amount, Equals, int64(40))
}

func (s EventSuite) TestAdminConfig(c *C) {
	evt := NewEventAdminConfig("foo", "bar")
	c.Check(evt.Type(), Equals, "admin_config")
	c.Check(evt.Key, Equals, "foo")
	c.Check(evt.Value, Equals, "bar")
}

func (s EventSuite) TestEvent(c *C) {
	txID, err := common.NewTxID("A1C7D97D5DB51FFDBC3FE29FFF6ADAA2DAF112D2CEAADA0902822333A59BD218")
	c.Assert(err, IsNil)
	swap := NewEventSwap(
		common.BNBAsset,
		sdk.NewUint(5),
		sdk.NewUint(5),
		sdk.NewUint(5),
	)

	swapBytes, _ := json.Marshal(swap)
	evt := NewEvent(swap.Type(),
		12,
		common.NewTx(
			txID,
			GetRandomBNBAddress(),
			GetRandomBNBAddress(),
			common.Coins{
				common.NewCoin(common.BNBAsset, sdk.NewUint(320000000)),
				common.NewCoin(common.RuneAsset(), sdk.NewUint(420000000)),
			},
			common.BNBGasFeeSingleton,
			"SWAP:BNB.BNB",
		),
		swapBytes,
		Success,
	)

	c.Check(evt.Empty(), Equals, false)

	txID, err = common.NewTxID("B1C7D97D5DB51FFDBC3FE29FFF6ADAA2DAF112D2CEAADA0902822333A59BD218")
	c.Assert(err, IsNil)
	stake := NewEventStake(
		common.BNBAsset,
		sdk.NewUint(5),
	)
	stakeBytes, _ := json.Marshal(stake)
	evt2 := NewEvent(stake.Type(),
		12,
		common.NewTx(
			txID,
			GetRandomBNBAddress(),
			GetRandomBNBAddress(),
			common.Coins{
				common.NewCoin(common.BNBAsset, sdk.NewUint(320000000)),
				common.NewCoin(common.RuneAsset(), sdk.NewUint(420000000)),
			},
			common.BNBGasFeeSingleton,
			"SWAP:BNB.BNB",
		),
		stakeBytes,
		Success,
	)

	events := Events{evt, evt2}
	found, events := events.PopByInHash(txID)
	c.Assert(found, HasLen, 1)
	c.Check(found[0].Empty(), Equals, false)
	c.Check(found[0].Type, Equals, evt2.Type)
	c.Assert(events, HasLen, 1)
	c.Check(events[0].Type, Equals, evt.Type)

	c.Check(Event{}.Empty(), Equals, true)
}

func (s EventSuite) TestSlash(c *C) {
	evt := NewEventSlash(common.BNBAsset, []PoolAmt{
		{common.BNBAsset, -20},
		{common.RuneAsset(), 30},
	})
	c.Check(evt.Type(), Equals, "slash")
	c.Check(evt.Pool, Equals, common.BNBAsset)
	c.Assert(evt.SlashAmount, HasLen, 2)
	c.Check(evt.SlashAmount[0].Asset, Equals, common.BNBAsset)
	c.Check(evt.SlashAmount[0].Amount, Equals, int64(-20))
	c.Check(evt.SlashAmount[1].Asset, Equals, common.RuneAsset())
	c.Check(evt.SlashAmount[1].Amount, Equals, int64(30))
}
