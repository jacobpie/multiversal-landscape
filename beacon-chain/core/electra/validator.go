package electra

import (
	"context"

	"github.com/prysmaticlabs/prysm/v5/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/v5/config/params"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/primitives"
)

// QueueEntireBalanceAndResetValidator queues the entire balance and resets the validator. This is used in electra fork logic.
//
// Spec definition:
//
//	def queue_entire_balance_and_reset_validator(state: BeaconState, index: ValidatorIndex) -> None:
//	    balance = state.balances[index]
//	    state.balances[index] = 0
//	    validator = state.validators[index]
//	    validator.effective_balance = 0
//	    validator.activation_eligibility_epoch = FAR_FUTURE_EPOCH
//	    state.pending_balance_deposits.append(
//	        PendingBalanceDeposit(index=index, amount=balance)
//	    )
func QueueEntireBalanceAndResetValidator(ctx context.Context, s state.BeaconState, idx primitives.ValidatorIndex) error {
	bal, err := s.BalanceAtIndex(idx)
	if err != nil {
		return err
	}

	if err := s.UpdateBalancesAtIndex(idx, 0); err != nil {
		return err
	}

	v, err := s.ValidatorAtIndex(idx)
	if err != nil {
		return err
	}

	v.EffectiveBalance = 0
	v.ActivationEligibilityEpoch = params.BeaconConfig().FarFutureEpoch
	if err := s.UpdateValidatorAtIndex(idx, v); err != nil {
		return err
	}

	return s.AppendPendingBalanceDeposit(idx, bal)
}
