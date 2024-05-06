//go:build !minimal

package state_native

import (
	"encoding/json"
	"sync"

	"github.com/prysmaticlabs/go-bitfield"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/state/fieldtrie"
	customtypes "github.com/prysmaticlabs/prysm/v5/beacon-chain/state/state-native/custom-types"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/state/state-native/types"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/state/stateutil"
	"github.com/prysmaticlabs/prysm/v5/config/features"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v5/math"
	enginev1 "github.com/prysmaticlabs/prysm/v5/proto/engine/v1"
	ethpb "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
)

// BeaconState defines a struct containing utilities for the Ethereum Beacon Chain state, defining
// getters and setters for its respective values and helpful functions such as HashTreeRoot().
type BeaconState struct {
	version                             int
	genesisTime                         uint64
	genesisValidatorsRoot               [32]byte
	slot                                primitives.Slot
	fork                                *ethpb.Fork
	latestBlockHeader                   *ethpb.BeaconBlockHeader
	blockRoots                          customtypes.BlockRoots
	blockRootsMultiValue                *MultiValueBlockRoots
	stateRoots                          customtypes.StateRoots
	stateRootsMultiValue                *MultiValueStateRoots
	historicalRoots                     customtypes.HistoricalRoots
	historicalSummaries                 []*ethpb.HistoricalSummary
	eth1Data                            *ethpb.Eth1Data
	eth1DataVotes                       []*ethpb.Eth1Data
	eth1DepositIndex                    uint64
	validators                          []*ethpb.Validator
	validatorsMultiValue                *MultiValueValidators
	balances                            []uint64
	balancesMultiValue                  *MultiValueBalances
	randaoMixes                         customtypes.RandaoMixes
	randaoMixesMultiValue               *MultiValueRandaoMixes
	slashings                           []uint64
	previousEpochAttestations           []*ethpb.PendingAttestation
	currentEpochAttestations            []*ethpb.PendingAttestation
	previousEpochParticipation          []byte
	currentEpochParticipation           []byte
	justificationBits                   bitfield.Bitvector4
	previousJustifiedCheckpoint         *ethpb.Checkpoint
	currentJustifiedCheckpoint          *ethpb.Checkpoint
	finalizedCheckpoint                 *ethpb.Checkpoint
	inactivityScores                    []uint64
	inactivityScoresMultiValue          *MultiValueInactivityScores
	currentSyncCommittee                *ethpb.SyncCommittee
	nextSyncCommittee                   *ethpb.SyncCommittee
	latestExecutionPayloadHeader        *enginev1.ExecutionPayloadHeader
	latestExecutionPayloadHeaderCapella *enginev1.ExecutionPayloadHeaderCapella
	latestExecutionPayloadHeaderDeneb   *enginev1.ExecutionPayloadHeaderDeneb
	latestExecutionPayloadHeaderElectra *enginev1.ExecutionPayloadHeaderElectra
	nextWithdrawalIndex                 uint64
	nextWithdrawalValidatorIndex        primitives.ValidatorIndex

	// Electra fields
	depositReceiptsStartIndex     uint64
	depositBalanceToConsume       math.Gwei
	exitBalanceToConsume          math.Gwei
	earliestExitEpoch             primitives.Epoch
	consolidationBalanceToConsume math.Gwei
	earliestConsolidationEpoch    primitives.Epoch
	pendingBalanceDeposits        []*ethpb.PendingBalanceDeposit    // pending_balance_deposits: List[PendingBalanceDeposit, PENDING_BALANCE_DEPOSITS_LIMIT]
	pendingPartialWithdrawals     []*ethpb.PendingPartialWithdrawal // pending_partial_withdrawals: List[PartialWithdrawal, PENDING_PARTIAL_WITHDRAWALS_LIMIT]
	pendingConsolidations         []*ethpb.PendingConsolidation     // pending_consolidations: List[PendingConsolidation, PENDING_CONSOLIDATIONS_LIMIT]

	id                    uint64
	lock                  sync.RWMutex
	dirtyFields           map[types.FieldIndex]bool
	dirtyIndices          map[types.FieldIndex][]uint64
	stateFieldLeaves      map[types.FieldIndex]*fieldtrie.FieldTrie
	rebuildTrie           map[types.FieldIndex]bool
	valMapHandler         *stateutil.ValidatorMapHandler
	validatorIndexCache   *finalizedValidatorIndexCache
	merkleLayers          [][][]byte
	sharedFieldReferences map[types.FieldIndex]*stateutil.Reference
}

type beaconStateMarshalable struct {
	Version                             int                                     `json:"version" yaml:"version"`
	GenesisTime                         uint64                                  `json:"genesis_time" yaml:"genesis_time"`
	GenesisValidatorsRoot               [32]byte                                `json:"genesis_validators_root" yaml:"genesis_validators_root"`
	Slot                                primitives.Slot                         `json:"slot" yaml:"slot"`
	Fork                                *ethpb.Fork                             `json:"fork" yaml:"fork"`
	LatestBlockHeader                   *ethpb.BeaconBlockHeader                `json:"latest_block_header" yaml:"latest_block_header"`
	BlockRoots                          customtypes.BlockRoots                  `json:"block_roots" yaml:"block_roots"`
	StateRoots                          customtypes.StateRoots                  `json:"state_roots" yaml:"state_roots"`
	HistoricalRoots                     customtypes.HistoricalRoots             `json:"historical_roots" yaml:"historical_roots"`
	HistoricalSummaries                 []*ethpb.HistoricalSummary              `json:"historical_summaries" yaml:"historical_summaries"`
	Eth1Data                            *ethpb.Eth1Data                         `json:"eth_1_data" yaml:"eth_1_data"`
	Eth1DataVotes                       []*ethpb.Eth1Data                       `json:"eth_1_data_votes" yaml:"eth_1_data_votes"`
	Eth1DepositIndex                    uint64                                  `json:"eth_1_deposit_index" yaml:"eth_1_deposit_index"`
	Validators                          []*ethpb.Validator                      `json:"validators" yaml:"validators"`
	Balances                            []uint64                                `json:"balances" yaml:"balances"`
	RandaoMixes                         customtypes.RandaoMixes                 `json:"randao_mixes" yaml:"randao_mixes"`
	Slashings                           []uint64                                `json:"slashings" yaml:"slashings"`
	PreviousEpochAttestations           []*ethpb.PendingAttestation             `json:"previous_epoch_attestations" yaml:"previous_epoch_attestations"`
	CurrentEpochAttestations            []*ethpb.PendingAttestation             `json:"current_epoch_attestations" yaml:"current_epoch_attestations"`
	PreviousEpochParticipation          []byte                                  `json:"previous_epoch_participation" yaml:"previous_epoch_participation"`
	CurrentEpochParticipation           []byte                                  `json:"current_epoch_participation" yaml:"current_epoch_participation"`
	JustificationBits                   bitfield.Bitvector4                     `json:"justification_bits" yaml:"justification_bits"`
	PreviousJustifiedCheckpoint         *ethpb.Checkpoint                       `json:"previous_justified_checkpoint" yaml:"previous_justified_checkpoint"`
	CurrentJustifiedCheckpoint          *ethpb.Checkpoint                       `json:"current_justified_checkpoint" yaml:"current_justified_checkpoint"`
	FinalizedCheckpoint                 *ethpb.Checkpoint                       `json:"finalized_checkpoint" yaml:"finalized_checkpoint"`
	InactivityScores                    []uint64                                `json:"inactivity_scores" yaml:"inactivity_scores"`
	CurrentSyncCommittee                *ethpb.SyncCommittee                    `json:"current_sync_committee" yaml:"current_sync_committee"`
	NextSyncCommittee                   *ethpb.SyncCommittee                    `json:"next_sync_committee" yaml:"next_sync_committee"`
	LatestExecutionPayloadHeader        *enginev1.ExecutionPayloadHeader        `json:"latest_execution_payload_header" yaml:"latest_execution_payload_header"`
	LatestExecutionPayloadHeaderCapella *enginev1.ExecutionPayloadHeaderCapella `json:"latest_execution_payload_header_capella" yaml:"latest_execution_payload_header_capella"`
	LatestExecutionPayloadHeaderDeneb   *enginev1.ExecutionPayloadHeaderDeneb   `json:"latest_execution_payload_header_deneb" yaml:"latest_execution_payload_header_deneb"`
	LatestExecutionPayloadHeaderElectra *enginev1.ExecutionPayloadHeaderElectra `json:"latest_execution_payload_header_electra" yaml:"latest_execution_payload_header_electra"`
	NextWithdrawalIndex                 uint64                                  `json:"next_withdrawal_index" yaml:"next_withdrawal_index"`
	NextWithdrawalValidatorIndex        primitives.ValidatorIndex               `json:"next_withdrawal_validator_index" yaml:"next_withdrawal_validator_index"`
	DepositReceiptsStartIndex           uint64                                  `json:"deposit_receipts_start_index" yaml:"deposit_receipts_start_index"`
	DepositBalanceToConsume             math.Gwei                               `json:"deposit_balance_to_consume" yaml:"deposit_balance_to_consume"`
	ExitBalanceToConsume                math.Gwei                               `json:"exit_balance_to_consume" yaml:"exit_balance_to_consume"`
	EarliestExitEpoch                   primitives.Epoch                        `json:"earliest_exit_epoch" yaml:"earliest_exit_epoch"`
	ConsolidationBalanceToConsume       math.Gwei                               `json:"consolidation_balance_to_consume" yaml:"consolidation_balance_to_consume"`
	EarliestConsolidationEpoch          primitives.Epoch                        `json:"earliest_consolidation_epoch" yaml:"earliest_consolidation_epoch"`
	PendingBalanceDeposits              []*ethpb.PendingBalanceDeposit          `json:"pending_balance_deposits" yaml:"pending_balance_deposits"`
	PendingPartialWithdrawals           []*ethpb.PendingPartialWithdrawal       `json:"pending_partial_withdrawals" yaml:"pending_partial_withdrawals"`
	PendingConsolidations               []*ethpb.PendingConsolidation           `json:"pending_consolidations" yaml:"pending_consolidations"`
}

func (b *BeaconState) MarshalJSON() ([]byte, error) {
	var bRoots customtypes.BlockRoots
	var sRoots customtypes.StateRoots
	var mixes customtypes.RandaoMixes
	var balances []uint64
	var inactivityScores []uint64
	var vals []*ethpb.Validator

	if features.Get().EnableExperimentalState {
		bRoots = b.blockRootsMultiValue.Value(b)
		sRoots = b.stateRootsMultiValue.Value(b)
		mixes = b.randaoMixesMultiValue.Value(b)
		balances = b.balancesMultiValue.Value(b)
		inactivityScores = b.inactivityScoresMultiValue.Value(b)
		vals = b.validatorsMultiValue.Value(b)
	} else {
		bRoots = b.blockRoots
		sRoots = b.stateRoots
		mixes = b.randaoMixes
		balances = b.balances
		inactivityScores = b.inactivityScores
		vals = b.validators
	}

	marshalable := &beaconStateMarshalable{
		Version:                             b.version,
		GenesisTime:                         b.genesisTime,
		GenesisValidatorsRoot:               b.genesisValidatorsRoot,
		Slot:                                b.slot,
		Fork:                                b.fork,
		LatestBlockHeader:                   b.latestBlockHeader,
		BlockRoots:                          bRoots,
		StateRoots:                          sRoots,
		HistoricalRoots:                     b.historicalRoots,
		HistoricalSummaries:                 b.historicalSummaries,
		Eth1Data:                            b.eth1Data,
		Eth1DataVotes:                       b.eth1DataVotes,
		Eth1DepositIndex:                    b.eth1DepositIndex,
		Validators:                          vals,
		Balances:                            balances,
		RandaoMixes:                         mixes,
		Slashings:                           b.slashings,
		PreviousEpochAttestations:           b.previousEpochAttestations,
		CurrentEpochAttestations:            b.currentEpochAttestations,
		PreviousEpochParticipation:          b.previousEpochParticipation,
		CurrentEpochParticipation:           b.currentEpochParticipation,
		JustificationBits:                   b.justificationBits,
		PreviousJustifiedCheckpoint:         b.previousJustifiedCheckpoint,
		CurrentJustifiedCheckpoint:          b.currentJustifiedCheckpoint,
		FinalizedCheckpoint:                 b.finalizedCheckpoint,
		InactivityScores:                    inactivityScores,
		CurrentSyncCommittee:                b.currentSyncCommittee,
		NextSyncCommittee:                   b.nextSyncCommittee,
		LatestExecutionPayloadHeader:        b.latestExecutionPayloadHeader,
		LatestExecutionPayloadHeaderCapella: b.latestExecutionPayloadHeaderCapella,
		LatestExecutionPayloadHeaderDeneb:   b.latestExecutionPayloadHeaderDeneb,
		LatestExecutionPayloadHeaderElectra: b.latestExecutionPayloadHeaderElectra,
		NextWithdrawalIndex:                 b.nextWithdrawalIndex,
		NextWithdrawalValidatorIndex:        b.nextWithdrawalValidatorIndex,
		DepositReceiptsStartIndex:           b.depositReceiptsStartIndex,
		DepositBalanceToConsume:             b.depositBalanceToConsume,
		ExitBalanceToConsume:                b.exitBalanceToConsume,
		EarliestExitEpoch:                   b.earliestExitEpoch,
		ConsolidationBalanceToConsume:       b.consolidationBalanceToConsume,
		EarliestConsolidationEpoch:          b.earliestConsolidationEpoch,
		PendingBalanceDeposits:              b.pendingBalanceDeposits,
		PendingPartialWithdrawals:           b.pendingPartialWithdrawals,
		PendingConsolidations:               b.pendingConsolidations,
	}
	return json.Marshal(marshalable)
}
