package swapservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

func validateUnstake(ctx sdk.Context, keeper poolStorage, msg MsgSetUnStake) error {
	if msg.PublicAddress.Empty() {
		return errors.New("empty public address")
	}
	if msg.WithdrawBasisPoints.Empty() {
		return errors.New("empty withdraw basis points")
	}
	if msg.RequestTxHash.Empty() {
		return errors.New("request tx hash is empty")
	}
	if msg.Ticker.Empty() {
		return errors.New("empty ticker")
	}
	withdrawBasisPoints := msg.WithdrawBasisPoints.Float64()
	if withdrawBasisPoints <= 0 || withdrawBasisPoints > MaxWithdrawBasisPoints {
		return errors.Errorf("withdraw basis points %s is invalid", msg.WithdrawBasisPoints)
	}
	if !keeper.PoolExist(ctx, msg.Ticker) {
		// pool doesn't exist
		return errors.Errorf("pool-%s doesn't exist", msg.Ticker)
	}
	return nil
}

// unstake withdraw all the asset
func unstake(ctx sdk.Context, keeper poolStorage, msg MsgSetUnStake) (Amount, Amount, error) {
	if err := validateUnstake(ctx, keeper, msg); nil != err {
		return "0", "0", err
	}
	withdrawPercentage := msg.WithdrawBasisPoints.Float64() / 100 // convert from basis point to percentage
	// here fBalance should be valid , because we did the validation above
	pool := keeper.GetPool(ctx, msg.Ticker)
	poolStaker, err := keeper.GetPoolStaker(ctx, msg.Ticker)
	if nil != err {
		return "0", "0", errors.Wrap(err, "can't find pool staker")

	}
	stakerPool, err := keeper.GetStakerPool(ctx, msg.PublicAddress)
	if nil != err {
		return "0", "0", errors.Wrap(err, "can't find staker pool")
	}

	poolUnits := pool.PoolUnits.Float64()
	poolRune := pool.BalanceRune.Float64()
	poolToken := pool.BalanceToken.Float64()
	stakerUnit := poolStaker.GetStakerUnit(msg.PublicAddress)
	fStakerUnit := stakerUnit.Units.Float64()
	if !stakerUnit.Units.GreaterThen(0) {
		return "0", "0", errors.New("nothing to withdraw")
	}

	ctx.Logger().Info("pool before unstake", "pool unit", poolUnits, "balance RUNE", poolRune, "balance token", poolToken)
	ctx.Logger().Info("staker before withdraw", "staker unit", fStakerUnit)
	withdrawRune, withDrawToken, unitAfter, err := calculateUnstake(poolUnits, poolRune, poolToken, fStakerUnit, withdrawPercentage)
	if err != nil {
		return "0", "0", err
	}
	ctx.Logger().Info("client withdraw", "RUNE", withdrawRune, "token", withDrawToken, "units left", unitAfter)
	// update pool
	pool.PoolUnits = NewAmountFromFloat(poolUnits - fStakerUnit + unitAfter)
	pool.BalanceRune = NewAmountFromFloat(poolRune - withdrawRune)
	pool.BalanceToken = NewAmountFromFloat(poolToken - withDrawToken)
	ctx.Logger().Info("pool after unstake", "pool unit", pool.PoolUnits, "balance RUNE", pool.BalanceRune, "balance token", pool.BalanceToken)
	// update pool staker
	poolStaker.TotalUnits = pool.PoolUnits
	if unitAfter == 0 {
		// just remove it
		poolStaker.RemoveStakerUnit(msg.PublicAddress)
	} else {
		stakerUnit.Units = NewAmountFromFloat(unitAfter)
		poolStaker.UpsertStakerUnit(stakerUnit)
	}
	if unitAfter <= 0 {
		stakerPool.RemoveStakerPoolItem(msg.Ticker)
	} else {
		spi := stakerPool.GetStakerPoolItem(msg.Ticker)
		spi.Units = NewAmountFromFloat(unitAfter)
		stakerPool.UpsertStakerPoolItem(spi)
	}
	// update staker pool
	keeper.SetPool(ctx, pool)
	keeper.SetPoolStaker(ctx, msg.Ticker, poolStaker)
	keeper.SetStakerPool(ctx, msg.PublicAddress, stakerPool)
	return NewAmountFromFloat(withdrawRune), NewAmountFromFloat(withDrawToken), nil
}

func calculateUnstake(poolUnit, poolRune, poolToken, stakerUnit, percentage float64) (float64, float64, float64, error) {
	if poolUnit <= 0 {
		return 0, 0, 0, errors.New("poolUnits can't be zero or negative")
	}
	if poolRune <= 0 {
		return 0, 0, 0, errors.New("pool rune balance can't be zero or negative")
	}
	if poolToken <= 0 {
		return 0, 0, 0, errors.New("pool token balance can't be zero or negative")
	}
	if stakerUnit < 0 {
		return 0, 0, 0, errors.New("staker unit can't be negative")
	}
	if percentage < 0 || percentage > 100 {
		return 0, 0, 0, errors.Errorf("percentage %f is not valid", percentage)
	}
	stakerOwnership := stakerUnit / poolUnit
	withdrawRune := stakerOwnership * percentage / 100 * poolRune
	withdrawToken := stakerOwnership * percentage / 100 * poolToken
	unitAfter := stakerUnit * (100 - percentage) / 100
	return withdrawRune, withdrawToken, unitAfter, nil
}
