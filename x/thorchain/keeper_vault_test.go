package thorchain

import (
	"github.com/cosmos/cosmos-sdk/types"
	. "gopkg.in/check.v1"

	"gitlab.com/thorchain/thornode/common"
)

type KeeperVaultSuite struct{}

var _ = Suite(&KeeperVaultSuite{})

func (s *KeeperVaultSuite) TestVault(c *C) {
	ctx, k := setupKeeperForTest(c)
	pubKey := GetRandomPubKey()
	yggdrasil := NewVault(ctx.BlockHeight(), ActiveVault, YggdrasilVault, pubKey)
	err := k.SetVault(ctx, yggdrasil)
	c.Assert(err, IsNil)
	c.Assert(k.VaultExists(ctx, pubKey), Equals, true)
	pubKey1 := GetRandomPubKey()
	yggdrasil1 := NewVault(ctx.BlockHeight(), ActiveVault, YggdrasilVault, pubKey1)
	yggdrasil1.Coins = common.Coins{
		common.NewCoin(common.BNBAsset, types.NewUint(100)),
	}
	c.Assert(k.SetVault(ctx, yggdrasil1), IsNil)
	ygg, err := k.GetVault(ctx, pubKey1)
	c.Assert(err, IsNil)
	c.Assert(ygg.IsEmpty(), Equals, false)
	hasYgg, err := k.HasValidVaultPools(ctx)
	c.Assert(err, IsNil)
	c.Assert(hasYgg, Equals, true)

	asgards, err := k.GetActiveAsgardVaults(ctx)
	c.Assert(err, IsNil)
	c.Assert(asgards, HasLen, 0)
	pubKey = GetRandomPubKey()
	asgard := NewVault(ctx.BlockHeight(), ActiveVault, AsgardVault, pubKey)
	c.Assert(k.SetVault(ctx, asgard), IsNil)
	asgard2 := NewVault(ctx.BlockHeight(), InactiveVault, AsgardVault, GetRandomPubKey())
	c.Assert(k.SetVault(ctx, asgard2), IsNil)
	asgards, err = k.GetActiveAsgardVaults(ctx)
	c.Assert(err, IsNil)
	c.Assert(asgards, HasLen, 1)
	c.Check(asgards[0].PubKey.Equals(pubKey), Equals, true)

	c.Assert(k.DeleteVault(ctx, pubKey), IsNil)
	c.Assert(k.DeleteVault(ctx, pubKey), IsNil) // second time should also not error
	asgards, err = k.GetActiveAsgardVaults(ctx)
	c.Assert(err, IsNil)
	c.Assert(asgards, HasLen, 0)
}
