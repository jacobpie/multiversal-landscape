package blocks

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/core/electra"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/core/signing"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/v5/config/params"
	"github.com/prysmaticlabs/prysm/v5/container/trie"
	"github.com/prysmaticlabs/prysm/v5/contracts/deposit"
	"github.com/prysmaticlabs/prysm/v5/crypto/bls"
	"github.com/prysmaticlabs/prysm/v5/encoding/bytesutil"
	"github.com/prysmaticlabs/prysm/v5/math"
	ethpb "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v5/runtime/version"
)

// ProcessPreGenesisDeposits processes a deposit for the beacon state before chainstart.
func ProcessPreGenesisDeposits(
	ctx context.Context,
	beaconState state.BeaconState,
	deposits []*ethpb.Deposit,
) (state.BeaconState, error) {
	var err error
	beaconState, err = ProcessDeposits(ctx, beaconState, deposits)
	if err != nil {
		return nil, errors.Wrap(err, "could not process deposit")
	}
	beaconState, err = ActivateValidatorWithEffectiveBalance(beaconState, deposits)
	if err != nil {
		return nil, err
	}
	return beaconState, nil
}

// ActivateValidatorWithEffectiveBalance updates validator's effective balance, and if it's above MaxEffectiveBalance, validator becomes active in genesis.
func ActivateValidatorWithEffectiveBalance(beaconState state.BeaconState, deposits []*ethpb.Deposit) (state.BeaconState, error) {
	for _, d := range deposits {
		pubkey := d.Data.PublicKey
		index, ok := beaconState.ValidatorIndexByPubkey(bytesutil.ToBytes48(pubkey))
		// In the event of the pubkey not existing, we continue processing the other
		// deposits.
		if !ok {
			continue
		}
		balance, err := beaconState.BalanceAtIndex(index)
		if err != nil {
			return nil, err
		}
		validator, err := beaconState.ValidatorAtIndex(index)
		if err != nil {
			return nil, err
		}
		validator.EffectiveBalance = math.Min(balance-balance%params.BeaconConfig().EffectiveBalanceIncrement, params.BeaconConfig().MaxEffectiveBalance)
		if validator.EffectiveBalance ==
			params.BeaconConfig().MaxEffectiveBalance {
			validator.ActivationEligibilityEpoch = 0
			validator.ActivationEpoch = 0
		}
		if err := beaconState.UpdateValidatorAtIndex(index, validator); err != nil {
			return nil, err
		}
	}
	return beaconState, nil
}

// ProcessDeposits is one of the operations performed on each processed
// beacon block to verify queued validators from the Ethereum 1.0 Deposit Contract
// into the beacon chain.
//
// Spec pseudocode definition:
//
//	For each deposit in block.body.deposits:
//	  process_deposit(state, deposit)
func ProcessDeposits(
	ctx context.Context,
	beaconState state.BeaconState,
	deposits []*ethpb.Deposit,
) (state.BeaconState, error) {
	// Attempt to verify all deposit signatures at once, if this fails then fall back to processing
	// individual deposits with signature verification enabled.
	batchVerified, err := BatchVerifyDepositsSignatures(ctx, deposits)
	if err != nil {
		return nil, err
	}

	for _, d := range deposits {
		if d == nil || d.Data == nil {
			return nil, errors.New("got a nil deposit in block")
		}
		beaconState, err = ProcessDeposit(beaconState, d, batchVerified)
		if err != nil {
			return nil, errors.Wrapf(err, "could not process deposit from %#x", bytesutil.Trunc(d.Data.PublicKey))
		}
	}
	return beaconState, nil
}

// BatchVerifyDepositsSignatures batch verifies deposit signatures.
func BatchVerifyDepositsSignatures(ctx context.Context, deposits []*ethpb.Deposit) (bool, error) {
	var err error
	domain, err := signing.ComputeDomain(params.BeaconConfig().DomainDeposit, nil, nil)
	if err != nil {
		return false, err
	}

	verified := false
	if err := verifyDepositDataWithDomain(ctx, deposits, domain); err != nil {
		log.WithError(err).Debug("Failed to batch verify deposits signatures, will try individual verify")
		verified = true
	}
	return verified, nil
}

// ProcessDeposit takes in a deposit object and inserts it
// into the registry as a new validator or balance change.
// Returns the resulting state, a boolean to indicate whether or not the deposit
// resulted in a new validator entry into the beacon state, and any error.
//
// Spec pseudocode definition:
// def process_deposit(state: BeaconState, deposit: Deposit) -> None:
//
//		# Verify the Merkle branch
//		assert is_valid_merkle_branch(
//		    leaf=hash_tree_root(deposit.data),
//		    branch=deposit.proof,
//		    depth=DEPOSIT_CONTRACT_TREE_DEPTH + 1,  # Add 1 for the List length mix-in
//		    index=state.eth1_deposit_index,
//		    root=state.eth1_data.deposit_root,
//		)
//
//	 # Deposits must be processed in order
//	 state.eth1_deposit_index += 1
//
//	 apply_deposit(
//	  state=state,
//	  pubkey=deposit.data.pubkey,
//	  withdrawal_credentials=deposit.data.withdrawal_credentials,
//	  amount=deposit.data.amount,
//	  signature=deposit.data.signature,
//	 )
func ProcessDeposit(beaconState state.BeaconState, deposit *ethpb.Deposit, verifySignature bool) (state.BeaconState, error) {
	if err := verifyDeposit(beaconState, deposit); err != nil {
		if deposit == nil || deposit.Data == nil {
			return nil, err
		}
		return nil, errors.Wrapf(err, "could not verify deposit from %#x", bytesutil.Trunc(deposit.Data.PublicKey))
	}
	if err := beaconState.SetEth1DepositIndex(beaconState.Eth1DepositIndex() + 1); err != nil {
		return nil, err
	}
	return ApplyDeposit(beaconState, deposit.Data, verifySignature)
}

// ApplyDeposit
// def apply_deposit(state: BeaconState, pubkey: BLSPubkey, withdrawal_credentials: Bytes32, amount: uint64, signature: BLSSignature) -> None:
// validator_pubkeys = [v.pubkey for v in state.validators]
// if pubkey not in validator_pubkeys:
//
//	# Verify the deposit signature (proof of possession) which is not checked by the deposit contract
//	if is_valid_deposit_signature(pubkey, withdrawal_credentials, amount, signature):
//	  add_validator_to_registry(state, pubkey, withdrawal_credentials, amount)
//
// else:
//
//	# Increase balance by deposit amount
//	index = ValidatorIndex(validator_pubkeys.index(pubkey))
//	state.pending_balance_deposits.append(PendingBalanceDeposit(index=index, amount=amount))  # [Modified in Electra:EIP-7251]
//	# Check if valid deposit switch to compounding credentials
//
// if ( is_compounding_withdrawal_credential(withdrawal_credentials) and has_eth1_withdrawal_credential(state.validators[index])
//
//	 and is_valid_deposit_signature(pubkey, withdrawal_credentials, amount, signature)
//	):
//	 switch_to_compounding_validator(state, index)
func ApplyDeposit(beaconState state.BeaconState, data *ethpb.Deposit_Data, verifySignature bool) (state.BeaconState, error) {
	pubKey := data.PublicKey
	amount := data.Amount
	withdrawalCredentials := data.WithdrawalCredentials
	index, ok := beaconState.ValidatorIndexByPubkey(bytesutil.ToBytes48(pubKey))
	if !ok {
		if verifySignature {
			valid, err := IsValidDepositSignature(data)
			if err != nil {
				return nil, err
			}
			if !valid {
				return beaconState, nil
			}
		}
		if err := AddValidatorToRegistry(beaconState, pubKey, withdrawalCredentials, amount); err != nil {
			return nil, err
		}
	} else {
		if beaconState.Version() >= version.Electra {
			//if err := beaconState.AppendPendingBalanceDeposit(index, amount); err != nil {
			//		return nil, newValidator, err
			//	}
			val, err := beaconState.ValidatorAtIndex(index)
			if err != nil {
				return nil, err
			}

			hasValidSignature, err := IsValidDepositSignature(data)
			if err != nil {
				return nil, err
			}
			if helpers.IsCompoundingWithdrawalCredential(withdrawalCredentials) && helpers.HasETH1WithdrawalCredential(val) && hasValidSignature {
				if err := electra.SwitchToCompoundingValidator(beaconState, index); err != nil {
					return nil, err
				}
			}
		} else {
			if err := helpers.IncreaseBalance(beaconState, index, amount); err != nil {
				return nil, err
			}
		}
	}

	return beaconState, nil
}

// AddValidatorToRegistry
// def add_validator_to_registry(state: BeaconState,
//
//	                          pubkey: BLSPubkey,
//	                          withdrawal_credentials: Bytes32,
//	                          amount: uint64) -> None:
//	index = get_index_for_new_validator(state)
//	validator = get_validator_from_deposit(pubkey, withdrawal_credentials)
//	set_or_append_list(state.validators, index, validator)
//	set_or_append_list(state.balances, index, 0)  # [Modified in Electra:EIP7251]
//	set_or_append_list(state.previous_epoch_participation, index, ParticipationFlags(0b0000_0000))
//	set_or_append_list(state.current_epoch_participation, index, ParticipationFlags(0b0000_0000))
//	set_or_append_list(state.inactivity_scores, index, uint64(0))
//	state.pending_balance_deposits.append(PendingBalanceDeposit(index=index, amount=amount))  # [New in Electra:EIP7251]
func AddValidatorToRegistry(beaconState state.BeaconState, pubKey []byte, withdrawalCredentials []byte, amount uint64) error {
	val := GetValidatorFromDeposit(beaconState.Version(), pubKey, withdrawalCredentials, amount)
	if err := beaconState.AppendValidator(val); err != nil {
		return err
	}
	if beaconState.Version() >= version.Electra {
		if err := beaconState.AppendBalance(0); err != nil {
			return err
		}
		// In specs this function is at the end function
		//	if err := beaconState.AppendPendingBalanceDeposit(index, amount); err != nil {
		//		return nil, newValidator, err
		//	}
	} else {
		if err := beaconState.AppendBalance(amount); err != nil {
			return err
		}
	}

	if err := beaconState.AppendPreviousParticipationBits(0); err != nil {
		return err
	}
	if err := beaconState.AppendCurrentParticipationBits(0); err != nil {
		return err
	}
	return beaconState.AppendInactivityScore(0)
}

// GetValidatorFromDeposit gets a new validator object with provided paramaters
//
// def get_validator_from_deposit(pubkey: BLSPubkey, withdrawal_credentials: Bytes32) -> Validator:
//
//	return Validator(
//	pubkey=pubkey,
//	withdrawal_credentials=withdrawal_credentials,
//	activation_eligibility_epoch=FAR_FUTURE_EPOCH,
//	activation_epoch=FAR_FUTURE_EPOCH,
//	exit_epoch=FAR_FUTURE_EPOCH,
//	withdrawable_epoch=FAR_FUTURE_EPOCH,
//	effective_balance=0,  # [Modified in Electra:EIP7251]
//
// )
func GetValidatorFromDeposit(stateVersion int, pubKey []byte, withdrawalCredentials []byte, amount uint64) *ethpb.Validator {
	var effectiveBalance uint64
	if stateVersion >= version.Electra {
		effectiveBalance = 0 // [Modified in Electra:EIP7251]
	} else {
		effectiveBalance = amount - (amount % params.BeaconConfig().EffectiveBalanceIncrement)
		if params.BeaconConfig().MaxEffectiveBalance < effectiveBalance {
			effectiveBalance = params.BeaconConfig().MaxEffectiveBalance
		}
	}

	return &ethpb.Validator{
		PublicKey:                  pubKey,
		WithdrawalCredentials:      withdrawalCredentials,
		ActivationEligibilityEpoch: params.BeaconConfig().FarFutureEpoch,
		ActivationEpoch:            params.BeaconConfig().FarFutureEpoch,
		ExitEpoch:                  params.BeaconConfig().FarFutureEpoch,
		WithdrawableEpoch:          params.BeaconConfig().FarFutureEpoch,
		EffectiveBalance:           effectiveBalance,
	}
}

// IsValidDepositSignature returns whether deposit_data is valid
// def is_valid_deposit_signature(pubkey: BLSPubkey, withdrawal_credentials: Bytes32, amount: uint64, signature: BLSSignature) -> bool:
//
//	deposit_message = DepositMessage( pubkey=pubkey, withdrawal_credentials=withdrawal_credentials, amount=amount, )
//	domain = compute_domain(DOMAIN_DEPOSIT)  # Fork-agnostic domain since deposits are valid across forks
//	signing_root = compute_signing_root(deposit_message, domain)
//	return bls.Verify(pubkey, signing_root, signature)
func IsValidDepositSignature(data *ethpb.Deposit_Data) (bool, error) {
	domain, err := signing.ComputeDomain(params.BeaconConfig().DomainDeposit, nil, nil)
	if err != nil {
		return false, err
	}
	if err := verifyDepositDataSigningRoot(data, domain); err != nil {
		// Ignore this error as in the spec pseudo code.
		log.WithError(err).Debug("Skipping deposit: could not verify deposit data signature")
		return false, nil
	}
	return true, nil
}

func verifyDeposit(beaconState state.ReadOnlyBeaconState, deposit *ethpb.Deposit) error {
	// Verify Merkle proof of deposit and deposit trie root.
	if deposit == nil || deposit.Data == nil {
		return errors.New("received nil deposit or nil deposit data")
	}
	eth1Data := beaconState.Eth1Data()
	if eth1Data == nil {
		return errors.New("received nil eth1data in the beacon state")
	}

	receiptRoot := eth1Data.DepositRoot
	leaf, err := deposit.Data.HashTreeRoot()
	if err != nil {
		return errors.Wrap(err, "could not tree hash deposit data")
	}
	if ok := trie.VerifyMerkleProofWithDepth(
		receiptRoot,
		leaf[:],
		beaconState.Eth1DepositIndex(),
		deposit.Proof,
		params.BeaconConfig().DepositContractTreeDepth,
	); !ok {
		return fmt.Errorf(
			"deposit merkle branch of deposit root did not verify for root: %#x",
			receiptRoot,
		)
	}

	return nil
}

func verifyDepositDataSigningRoot(obj *ethpb.Deposit_Data, domain []byte) error {
	return deposit.VerifyDepositSignature(obj, domain)
}

func verifyDepositDataWithDomain(ctx context.Context, deps []*ethpb.Deposit, domain []byte) error {
	if len(deps) == 0 {
		return nil
	}
	pks := make([]bls.PublicKey, len(deps))
	sigs := make([][]byte, len(deps))
	msgs := make([][32]byte, len(deps))
	for i, dep := range deps {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if dep == nil || dep.Data == nil {
			return errors.New("nil deposit")
		}
		dpk, err := bls.PublicKeyFromBytes(dep.Data.PublicKey)
		if err != nil {
			return err
		}
		pks[i] = dpk
		sigs[i] = dep.Data.Signature
		depositMessage := &ethpb.DepositMessage{
			PublicKey:             dep.Data.PublicKey,
			WithdrawalCredentials: dep.Data.WithdrawalCredentials,
			Amount:                dep.Data.Amount,
		}
		sr, err := signing.ComputeSigningRoot(depositMessage, domain)
		if err != nil {
			return err
		}
		msgs[i] = sr
	}
	verify, err := bls.VerifyMultipleSignatures(sigs, msgs, pks)
	if err != nil {
		return errors.Errorf("could not verify multiple signatures: %v", err)
	}
	if !verify {
		return errors.New("one or more deposit signatures did not verify")
	}
	return nil
}
