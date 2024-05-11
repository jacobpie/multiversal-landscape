package blocks

import (
	"bytes"
	"errors"
	"math/big"

	fastssz "github.com/prysmaticlabs/fastssz"
	fieldparams "github.com/prysmaticlabs/prysm/v5/config/fieldparams"
	consensus_types "github.com/prysmaticlabs/prysm/v5/consensus-types"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/interfaces"
	"github.com/prysmaticlabs/prysm/v5/encoding/bytesutil"
	"github.com/prysmaticlabs/prysm/v5/encoding/ssz"
	"github.com/prysmaticlabs/prysm/v5/math"
	enginev1 "github.com/prysmaticlabs/prysm/v5/proto/engine/v1"
	"google.golang.org/protobuf/proto"
)

// executionPayload is a convenience wrapper around a beacon block body's execution payload data structure
// This wrapper allows us to conform to a common interface so that beacon
// blocks for future forks can also be applied across Prysm without issues.
type executionPayload struct {
	p *enginev1.ExecutionPayload
}

// NewWrappedExecutionData creates an appropriate execution payload wrapper based on the incoming type.
func NewWrappedExecutionData(v proto.Message, weiValue math.Wei) (interfaces.ExecutionData, error) {
	if weiValue == nil {
		weiValue = new(big.Int).SetInt64(0)
	}
	switch pbStruct := v.(type) {
	case *enginev1.ExecutionPayload:
		return WrappedExecutionPayload(pbStruct)
	case *enginev1.ExecutionPayloadCapella:
		return WrappedExecutionPayloadCapella(pbStruct, weiValue)
	case *enginev1.ExecutionPayloadDeneb:
		return WrappedExecutionPayloadDeneb(pbStruct, weiValue)
	case *enginev1.ExecutionPayloadElectra:
		return WrappedExecutionPayloadElectra(pbStruct, weiValue)
	default:
		return nil, ErrUnsupportedVersion
	}
}

var _ interfaces.ExecutionData = &executionPayload{}

// WrappedExecutionPayload is a constructor which wraps a protobuf execution payload into an interface.
func WrappedExecutionPayload(p *enginev1.ExecutionPayload) (interfaces.ExecutionData, error) {
	w := executionPayload{p: p}
	if w.IsNil() {
		return nil, consensus_types.ErrNilObjectWrapped
	}
	return w, nil
}

// IsNil checks if the underlying data is nil.
func (e executionPayload) IsNil() bool {
	return e.p == nil
}

// IsBlinded returns true if the underlying data is blinded.
func (executionPayload) IsBlinded() bool {
	return false
}

// MarshalSSZ --
func (e executionPayload) MarshalSSZ() ([]byte, error) {
	return e.p.MarshalSSZ()
}

// MarshalSSZTo --
func (e executionPayload) MarshalSSZTo(dst []byte) ([]byte, error) {
	return e.p.MarshalSSZTo(dst)
}

// SizeSSZ --
func (e executionPayload) SizeSSZ() int {
	return e.p.SizeSSZ()
}

// UnmarshalSSZ --
func (e executionPayload) UnmarshalSSZ(buf []byte) error {
	return e.p.UnmarshalSSZ(buf)
}

// HashTreeRoot --
func (e executionPayload) HashTreeRoot() ([32]byte, error) {
	return e.p.HashTreeRoot()
}

// HashTreeRootWith --
func (e executionPayload) HashTreeRootWith(hh *fastssz.Hasher) error {
	return e.p.HashTreeRootWith(hh)
}

// Proto --
func (e executionPayload) Proto() proto.Message {
	return e.p
}

// ParentHash --
func (e executionPayload) ParentHash() []byte {
	return e.p.ParentHash
}

// FeeRecipient --
func (e executionPayload) FeeRecipient() []byte {
	return e.p.FeeRecipient
}

// StateRoot --
func (e executionPayload) StateRoot() []byte {
	return e.p.StateRoot
}

// ReceiptsRoot --
func (e executionPayload) ReceiptsRoot() []byte {
	return e.p.ReceiptsRoot
}

// LogsBloom --
func (e executionPayload) LogsBloom() []byte {
	return e.p.LogsBloom
}

// PrevRandao --
func (e executionPayload) PrevRandao() []byte {
	return e.p.PrevRandao
}

// BlockNumber --
func (e executionPayload) BlockNumber() uint64 {
	return e.p.BlockNumber
}

// GasLimit --
func (e executionPayload) GasLimit() uint64 {
	return e.p.GasLimit
}

// GasUsed --
func (e executionPayload) GasUsed() uint64 {
	return e.p.GasUsed
}

// Timestamp --
func (e executionPayload) Timestamp() uint64 {
	return e.p.Timestamp
}

// ExtraData --
func (e executionPayload) ExtraData() []byte {
	return e.p.ExtraData
}

// BaseFeePerGas --
func (e executionPayload) BaseFeePerGas() []byte {
	return e.p.BaseFeePerGas
}

// BlockHash --
func (e executionPayload) BlockHash() []byte {
	return e.p.BlockHash
}

// Transactions --
func (e executionPayload) Transactions() ([][]byte, error) {
	return e.p.Transactions, nil
}

// TransactionsRoot --
func (executionPayload) TransactionsRoot() ([]byte, error) {
	return nil, consensus_types.ErrUnsupportedField
}

// Withdrawals --
func (executionPayload) Withdrawals() ([]*enginev1.Withdrawal, error) {
	return nil, consensus_types.ErrUnsupportedField
}

// WithdrawalsRoot --
func (executionPayload) WithdrawalsRoot() ([]byte, error) {
	return nil, consensus_types.ErrUnsupportedField
}

// BlobGasUsed --
func (e executionPayload) BlobGasUsed() (uint64, error) {
	return 0, consensus_types.ErrUnsupportedField
}

// ExcessBlobGas --
func (e executionPayload) ExcessBlobGas() (uint64, error) {
	return 0, consensus_types.ErrUnsupportedField
}

// ValueInWei --
func (executionPayload) ValueInWei() (math.Wei, error) {
	return nil, consensus_types.ErrUnsupportedField
}

// ValueInGwei --
func (executionPayload) ValueInGwei() (uint64, error) {
	return 0, consensus_types.ErrUnsupportedField
}

// executionPayloadHeader is a convenience wrapper around a blinded beacon block body's execution header data structure
// This wrapper allows us to conform to a common interface so that beacon
// blocks for future forks can also be applied across Prysm without issues.
type executionPayloadHeader struct {
	p *enginev1.ExecutionPayloadHeader
}

var _ interfaces.ExecutionData = &executionPayloadHeader{}

// WrappedExecutionPayloadHeader is a constructor which wraps a protobuf execution header into an interface.
func WrappedExecutionPayloadHeader(p *enginev1.ExecutionPayloadHeader) (interfaces.ExecutionData, error) {
	w := executionPayloadHeader{p: p}
	if w.IsNil() {
		return nil, consensus_types.ErrNilObjectWrapped
	}
	return w, nil
}

// IsNil checks if the underlying data is nil.
func (e executionPayloadHeader) IsNil() bool {
	return e.p == nil
}

// IsBlinded returns true if the underlying data is a header.
func (executionPayloadHeader) IsBlinded() bool {
	return true
}

// MarshalSSZ --
func (e executionPayloadHeader) MarshalSSZ() ([]byte, error) {
	return e.p.MarshalSSZ()
}

// MarshalSSZTo --
func (e executionPayloadHeader) MarshalSSZTo(dst []byte) ([]byte, error) {
	return e.p.MarshalSSZTo(dst)
}

// SizeSSZ --
func (e executionPayloadHeader) SizeSSZ() int {
	return e.p.SizeSSZ()
}

// UnmarshalSSZ --
func (e executionPayloadHeader) UnmarshalSSZ(buf []byte) error {
	return e.p.UnmarshalSSZ(buf)
}

// HashTreeRoot --
func (e executionPayloadHeader) HashTreeRoot() ([32]byte, error) {
	return e.p.HashTreeRoot()
}

// HashTreeRootWith --
func (e executionPayloadHeader) HashTreeRootWith(hh *fastssz.Hasher) error {
	return e.p.HashTreeRootWith(hh)
}

// Proto --
func (e executionPayloadHeader) Proto() proto.Message {
	return e.p
}

// ParentHash --
func (e executionPayloadHeader) ParentHash() []byte {
	return e.p.ParentHash
}

// FeeRecipient --
func (e executionPayloadHeader) FeeRecipient() []byte {
	return e.p.FeeRecipient
}

// StateRoot --
func (e executionPayloadHeader) StateRoot() []byte {
	return e.p.StateRoot
}

// ReceiptsRoot --
func (e executionPayloadHeader) ReceiptsRoot() []byte {
	return e.p.ReceiptsRoot
}

// LogsBloom --
func (e executionPayloadHeader) LogsBloom() []byte {
	return e.p.LogsBloom
}

// PrevRandao --
func (e executionPayloadHeader) PrevRandao() []byte {
	return e.p.PrevRandao
}

// BlockNumber --
func (e executionPayloadHeader) BlockNumber() uint64 {
	return e.p.BlockNumber
}

// GasLimit --
func (e executionPayloadHeader) GasLimit() uint64 {
	return e.p.GasLimit
}

// GasUsed --
func (e executionPayloadHeader) GasUsed() uint64 {
	return e.p.GasUsed
}

// Timestamp --
func (e executionPayloadHeader) Timestamp() uint64 {
	return e.p.Timestamp
}

// ExtraData --
func (e executionPayloadHeader) ExtraData() []byte {
	return e.p.ExtraData
}

// BaseFeePerGas --
func (e executionPayloadHeader) BaseFeePerGas() []byte {
	return e.p.BaseFeePerGas
}

// BlockHash --
func (e executionPayloadHeader) BlockHash() []byte {
	return e.p.BlockHash
}

// Transactions --
func (executionPayloadHeader) Transactions() ([][]byte, error) {
	return nil, consensus_types.ErrUnsupportedField
}

// TransactionsRoot --
func (e executionPayloadHeader) TransactionsRoot() ([]byte, error) {
	return e.p.TransactionsRoot, nil
}

// Withdrawals --
func (executionPayloadHeader) Withdrawals() ([]*enginev1.Withdrawal, error) {
	return nil, consensus_types.ErrUnsupportedField
}

// WithdrawalsRoot --
func (executionPayloadHeader) WithdrawalsRoot() ([]byte, error) {
	return nil, consensus_types.ErrUnsupportedField
}

// BlobGasUsed --
func (e executionPayloadHeader) BlobGasUsed() (uint64, error) {
	return 0, consensus_types.ErrUnsupportedField
}

// ExcessBlobGas --
func (e executionPayloadHeader) ExcessBlobGas() (uint64, error) {
	return 0, consensus_types.ErrUnsupportedField
}

// ValueInWei --
func (executionPayloadHeader) ValueInWei() (math.Wei, error) {
	return nil, consensus_types.ErrUnsupportedField
}

// ValueInGwei --
func (executionPayloadHeader) ValueInGwei() (uint64, error) {
	return 0, consensus_types.ErrUnsupportedField
}

// PayloadToHeader converts `payload` into execution payload header format.
func PayloadToHeader(payload interfaces.ExecutionData) (*enginev1.ExecutionPayloadHeader, error) {
	txs, err := payload.Transactions()
	if err != nil {
		return nil, err
	}
	txRoot, err := ssz.TransactionsRoot(txs)
	if err != nil {
		return nil, err
	}
	return &enginev1.ExecutionPayloadHeader{
		ParentHash:       bytesutil.SafeCopyBytes(payload.ParentHash()),
		FeeRecipient:     bytesutil.SafeCopyBytes(payload.FeeRecipient()),
		StateRoot:        bytesutil.SafeCopyBytes(payload.StateRoot()),
		ReceiptsRoot:     bytesutil.SafeCopyBytes(payload.ReceiptsRoot()),
		LogsBloom:        bytesutil.SafeCopyBytes(payload.LogsBloom()),
		PrevRandao:       bytesutil.SafeCopyBytes(payload.PrevRandao()),
		BlockNumber:      payload.BlockNumber(),
		GasLimit:         payload.GasLimit(),
		GasUsed:          payload.GasUsed(),
		Timestamp:        payload.Timestamp(),
		ExtraData:        bytesutil.SafeCopyBytes(payload.ExtraData()),
		BaseFeePerGas:    bytesutil.SafeCopyBytes(payload.BaseFeePerGas()),
		BlockHash:        bytesutil.SafeCopyBytes(payload.BlockHash()),
		TransactionsRoot: txRoot[:],
	}, nil
}

// executionPayloadCapella is a convenience wrapper around a beacon block body's execution payload data structure
// This wrapper allows us to conform to a common interface so that beacon
// blocks for future forks can also be applied across Prysm without issues.
type executionPayloadCapella struct {
	p         *enginev1.ExecutionPayloadCapella
	weiValue  math.Wei
	gweiValue uint64
}

var _ interfaces.ExecutionData = &executionPayloadCapella{}

// WrappedExecutionPayloadCapella is a constructor which wraps a protobuf execution payload into an interface.
func WrappedExecutionPayloadCapella(p *enginev1.ExecutionPayloadCapella, value math.Wei) (interfaces.ExecutionData, error) {
	w := executionPayloadCapella{p: p, weiValue: value, gweiValue: uint64(math.WeiToGwei(value))}
	if w.IsNil() {
		return nil, consensus_types.ErrNilObjectWrapped
	}
	return w, nil
}

// IsNil checks if the underlying data is nil.
func (e executionPayloadCapella) IsNil() bool {
	return e.p == nil
}

// IsBlinded returns true if the underlying data is blinded.
func (executionPayloadCapella) IsBlinded() bool {
	return false
}

// MarshalSSZ --
func (e executionPayloadCapella) MarshalSSZ() ([]byte, error) {
	return e.p.MarshalSSZ()
}

// MarshalSSZTo --
func (e executionPayloadCapella) MarshalSSZTo(dst []byte) ([]byte, error) {
	return e.p.MarshalSSZTo(dst)
}

// SizeSSZ --
func (e executionPayloadCapella) SizeSSZ() int {
	return e.p.SizeSSZ()
}

// UnmarshalSSZ --
func (e executionPayloadCapella) UnmarshalSSZ(buf []byte) error {
	return e.p.UnmarshalSSZ(buf)
}

// HashTreeRoot --
func (e executionPayloadCapella) HashTreeRoot() ([32]byte, error) {
	return e.p.HashTreeRoot()
}

// HashTreeRootWith --
func (e executionPayloadCapella) HashTreeRootWith(hh *fastssz.Hasher) error {
	return e.p.HashTreeRootWith(hh)
}

// Proto --
func (e executionPayloadCapella) Proto() proto.Message {
	return e.p
}

// ParentHash --
func (e executionPayloadCapella) ParentHash() []byte {
	return e.p.ParentHash
}

// FeeRecipient --
func (e executionPayloadCapella) FeeRecipient() []byte {
	return e.p.FeeRecipient
}

// StateRoot --
func (e executionPayloadCapella) StateRoot() []byte {
	return e.p.StateRoot
}

// ReceiptsRoot --
func (e executionPayloadCapella) ReceiptsRoot() []byte {
	return e.p.ReceiptsRoot
}

// LogsBloom --
func (e executionPayloadCapella) LogsBloom() []byte {
	return e.p.LogsBloom
}

// PrevRandao --
func (e executionPayloadCapella) PrevRandao() []byte {
	return e.p.PrevRandao
}

// BlockNumber --
func (e executionPayloadCapella) BlockNumber() uint64 {
	return e.p.BlockNumber
}

// GasLimit --
func (e executionPayloadCapella) GasLimit() uint64 {
	return e.p.GasLimit
}

// GasUsed --
func (e executionPayloadCapella) GasUsed() uint64 {
	return e.p.GasUsed
}

// Timestamp --
func (e executionPayloadCapella) Timestamp() uint64 {
	return e.p.Timestamp
}

// ExtraData --
func (e executionPayloadCapella) ExtraData() []byte {
	return e.p.ExtraData
}

// BaseFeePerGas --
func (e executionPayloadCapella) BaseFeePerGas() []byte {
	return e.p.BaseFeePerGas
}

// BlockHash --
func (e executionPayloadCapella) BlockHash() []byte {
	return e.p.BlockHash
}

// Transactions --
func (e executionPayloadCapella) Transactions() ([][]byte, error) {
	return e.p.Transactions, nil
}

// TransactionsRoot --
func (executionPayloadCapella) TransactionsRoot() ([]byte, error) {
	return nil, consensus_types.ErrUnsupportedField
}

// Withdrawals --
func (e executionPayloadCapella) Withdrawals() ([]*enginev1.Withdrawal, error) {
	return e.p.Withdrawals, nil
}

// WithdrawalsRoot --
func (executionPayloadCapella) WithdrawalsRoot() ([]byte, error) {
	return nil, consensus_types.ErrUnsupportedField
}

// BlobGasUsed --
func (e executionPayloadCapella) BlobGasUsed() (uint64, error) {
	return 0, consensus_types.ErrUnsupportedField
}

// ExcessBlobGas --
func (e executionPayloadCapella) ExcessBlobGas() (uint64, error) {
	return 0, consensus_types.ErrUnsupportedField
}

// ValueInWei --
func (e executionPayloadCapella) ValueInWei() (math.Wei, error) {
	return e.weiValue, nil
}

// ValueInGwei --
func (e executionPayloadCapella) ValueInGwei() (uint64, error) {
	return e.gweiValue, nil
}

// executionPayloadHeaderCapella is a convenience wrapper around a blinded beacon block body's execution header data structure
// This wrapper allows us to conform to a common interface so that beacon
// blocks for future forks can also be applied across Prysm without issues.
type executionPayloadHeaderCapella struct {
	p         *enginev1.ExecutionPayloadHeaderCapella
	weiValue  math.Wei
	gweiValue uint64
}

var _ interfaces.ExecutionData = &executionPayloadHeaderCapella{}

// WrappedExecutionPayloadHeaderCapella is a constructor which wraps a protobuf execution header into an interface.
func WrappedExecutionPayloadHeaderCapella(p *enginev1.ExecutionPayloadHeaderCapella, value math.Wei) (interfaces.ExecutionData, error) {
	w := executionPayloadHeaderCapella{p: p, weiValue: value, gweiValue: uint64(math.WeiToGwei(value))}
	if w.IsNil() {
		return nil, consensus_types.ErrNilObjectWrapped
	}
	return w, nil
}

// IsNil checks if the underlying data is nil.
func (e executionPayloadHeaderCapella) IsNil() bool {
	return e.p == nil
}

// IsBlinded returns true if the underlying data is blinded.
func (executionPayloadHeaderCapella) IsBlinded() bool {
	return true
}

// MarshalSSZ --
func (e executionPayloadHeaderCapella) MarshalSSZ() ([]byte, error) {
	return e.p.MarshalSSZ()
}

// MarshalSSZTo --
func (e executionPayloadHeaderCapella) MarshalSSZTo(dst []byte) ([]byte, error) {
	return e.p.MarshalSSZTo(dst)
}

// SizeSSZ --
func (e executionPayloadHeaderCapella) SizeSSZ() int {
	return e.p.SizeSSZ()
}

// UnmarshalSSZ --
func (e executionPayloadHeaderCapella) UnmarshalSSZ(buf []byte) error {
	return e.p.UnmarshalSSZ(buf)
}

// HashTreeRoot --
func (e executionPayloadHeaderCapella) HashTreeRoot() ([32]byte, error) {
	return e.p.HashTreeRoot()
}

// HashTreeRootWith --
func (e executionPayloadHeaderCapella) HashTreeRootWith(hh *fastssz.Hasher) error {
	return e.p.HashTreeRootWith(hh)
}

// Proto --
func (e executionPayloadHeaderCapella) Proto() proto.Message {
	return e.p
}

// ParentHash --
func (e executionPayloadHeaderCapella) ParentHash() []byte {
	return e.p.ParentHash
}

// FeeRecipient --
func (e executionPayloadHeaderCapella) FeeRecipient() []byte {
	return e.p.FeeRecipient
}

// StateRoot --
func (e executionPayloadHeaderCapella) StateRoot() []byte {
	return e.p.StateRoot
}

// ReceiptsRoot --
func (e executionPayloadHeaderCapella) ReceiptsRoot() []byte {
	return e.p.ReceiptsRoot
}

// LogsBloom --
func (e executionPayloadHeaderCapella) LogsBloom() []byte {
	return e.p.LogsBloom
}

// PrevRandao --
func (e executionPayloadHeaderCapella) PrevRandao() []byte {
	return e.p.PrevRandao
}

// BlockNumber --
func (e executionPayloadHeaderCapella) BlockNumber() uint64 {
	return e.p.BlockNumber
}

// GasLimit --
func (e executionPayloadHeaderCapella) GasLimit() uint64 {
	return e.p.GasLimit
}

// GasUsed --
func (e executionPayloadHeaderCapella) GasUsed() uint64 {
	return e.p.GasUsed
}

// Timestamp --
func (e executionPayloadHeaderCapella) Timestamp() uint64 {
	return e.p.Timestamp
}

// ExtraData --
func (e executionPayloadHeaderCapella) ExtraData() []byte {
	return e.p.ExtraData
}

// BaseFeePerGas --
func (e executionPayloadHeaderCapella) BaseFeePerGas() []byte {
	return e.p.BaseFeePerGas
}

// BlockHash --
func (e executionPayloadHeaderCapella) BlockHash() []byte {
	return e.p.BlockHash
}

// Transactions --
func (executionPayloadHeaderCapella) Transactions() ([][]byte, error) {
	return nil, consensus_types.ErrUnsupportedField
}

// TransactionsRoot --
func (e executionPayloadHeaderCapella) TransactionsRoot() ([]byte, error) {
	return e.p.TransactionsRoot, nil
}

// Withdrawals --
func (executionPayloadHeaderCapella) Withdrawals() ([]*enginev1.Withdrawal, error) {
	return nil, consensus_types.ErrUnsupportedField
}

// WithdrawalsRoot --
func (e executionPayloadHeaderCapella) WithdrawalsRoot() ([]byte, error) {
	return e.p.WithdrawalsRoot, nil
}

// BlobGasUsed --
func (e executionPayloadHeaderCapella) BlobGasUsed() (uint64, error) {
	return 0, consensus_types.ErrUnsupportedField
}

// ExcessBlobGas --
func (e executionPayloadHeaderCapella) ExcessBlobGas() (uint64, error) {
	return 0, consensus_types.ErrUnsupportedField
}

// ValueInWei --
func (e executionPayloadHeaderCapella) ValueInWei() (math.Wei, error) {
	return e.weiValue, nil
}

// ValueInGwei --
func (e executionPayloadHeaderCapella) ValueInGwei() (uint64, error) {
	return e.gweiValue, nil
}

// PayloadToHeaderCapella converts `payload` into execution payload header format.
func PayloadToHeaderCapella(payload interfaces.ExecutionData) (*enginev1.ExecutionPayloadHeaderCapella, error) {
	txs, err := payload.Transactions()
	if err != nil {
		return nil, err
	}
	txRoot, err := ssz.TransactionsRoot(txs)
	if err != nil {
		return nil, err
	}
	withdrawals, err := payload.Withdrawals()
	if err != nil {
		return nil, err
	}
	withdrawalsRoot, err := ssz.WithdrawalSliceRoot(withdrawals, fieldparams.MaxWithdrawalsPerPayload)
	if err != nil {
		return nil, err
	}

	return &enginev1.ExecutionPayloadHeaderCapella{
		ParentHash:       bytesutil.SafeCopyBytes(payload.ParentHash()),
		FeeRecipient:     bytesutil.SafeCopyBytes(payload.FeeRecipient()),
		StateRoot:        bytesutil.SafeCopyBytes(payload.StateRoot()),
		ReceiptsRoot:     bytesutil.SafeCopyBytes(payload.ReceiptsRoot()),
		LogsBloom:        bytesutil.SafeCopyBytes(payload.LogsBloom()),
		PrevRandao:       bytesutil.SafeCopyBytes(payload.PrevRandao()),
		BlockNumber:      payload.BlockNumber(),
		GasLimit:         payload.GasLimit(),
		GasUsed:          payload.GasUsed(),
		Timestamp:        payload.Timestamp(),
		ExtraData:        bytesutil.SafeCopyBytes(payload.ExtraData()),
		BaseFeePerGas:    bytesutil.SafeCopyBytes(payload.BaseFeePerGas()),
		BlockHash:        bytesutil.SafeCopyBytes(payload.BlockHash()),
		TransactionsRoot: txRoot[:],
		WithdrawalsRoot:  withdrawalsRoot[:],
	}, nil
}

// PayloadToHeaderDeneb converts `payload` into execution payload header format.
func PayloadToHeaderDeneb(payload interfaces.ExecutionData) (*enginev1.ExecutionPayloadHeaderDeneb, error) {
	txs, err := payload.Transactions()
	if err != nil {
		return nil, err
	}
	txRoot, err := ssz.TransactionsRoot(txs)
	if err != nil {
		return nil, err
	}
	withdrawals, err := payload.Withdrawals()
	if err != nil {
		return nil, err
	}
	withdrawalsRoot, err := ssz.WithdrawalSliceRoot(withdrawals, fieldparams.MaxWithdrawalsPerPayload)
	if err != nil {
		return nil, err
	}
	blobGasUsed, err := payload.BlobGasUsed()
	if err != nil {
		return nil, err
	}
	excessBlobGas, err := payload.ExcessBlobGas()
	if err != nil {
		return nil, err
	}

	return &enginev1.ExecutionPayloadHeaderDeneb{
		ParentHash:       bytesutil.SafeCopyBytes(payload.ParentHash()),
		FeeRecipient:     bytesutil.SafeCopyBytes(payload.FeeRecipient()),
		StateRoot:        bytesutil.SafeCopyBytes(payload.StateRoot()),
		ReceiptsRoot:     bytesutil.SafeCopyBytes(payload.ReceiptsRoot()),
		LogsBloom:        bytesutil.SafeCopyBytes(payload.LogsBloom()),
		PrevRandao:       bytesutil.SafeCopyBytes(payload.PrevRandao()),
		BlockNumber:      payload.BlockNumber(),
		GasLimit:         payload.GasLimit(),
		GasUsed:          payload.GasUsed(),
		Timestamp:        payload.Timestamp(),
		ExtraData:        bytesutil.SafeCopyBytes(payload.ExtraData()),
		BaseFeePerGas:    bytesutil.SafeCopyBytes(payload.BaseFeePerGas()),
		BlockHash:        bytesutil.SafeCopyBytes(payload.BlockHash()),
		TransactionsRoot: txRoot[:],
		WithdrawalsRoot:  withdrawalsRoot[:],
		BlobGasUsed:      blobGasUsed,
		ExcessBlobGas:    excessBlobGas,
	}, nil
}

// PayloadToHeaderElectra converts `payload` into execution payload header format.
func PayloadToHeaderElectra(payload interfaces.ExecutionDataElectra) (*enginev1.ExecutionPayloadHeaderElectra, error) {
	txs, err := payload.Transactions()
	if err != nil {
		return nil, err
	}
	txRoot, err := ssz.TransactionsRoot(txs)
	if err != nil {
		return nil, err
	}
	withdrawals, err := payload.Withdrawals()
	if err != nil {
		return nil, err
	}
	withdrawalsRoot, err := ssz.WithdrawalSliceRoot(withdrawals, fieldparams.MaxWithdrawalsPerPayload)
	if err != nil {
		return nil, err
	}
	blobGasUsed, err := payload.BlobGasUsed()
	if err != nil {
		return nil, err
	}
	excessBlobGas, err := payload.ExcessBlobGas()
	if err != nil {
		return nil, err
	}

	depositReceipts := payload.DepositReceipts()
	depositReceiptsRoot, err := ssz.DepositReceiptSliceRoot(depositReceipts, fieldparams.MaxDepositReceiptsPerPayload)
	if err != nil {
		return nil, err
	}

	withdrawalRequests := payload.WithdrawalRequests()
	withdrawalRequestsRoot, err := ssz.WithdrawalRequestSliceRoot(withdrawalRequests, fieldparams.MaxWithdrawalRequestsPerPayload)
	if err != nil {
		return nil, err
	}

	return &enginev1.ExecutionPayloadHeaderElectra{
		ParentHash:             bytesutil.SafeCopyBytes(payload.ParentHash()),
		FeeRecipient:           bytesutil.SafeCopyBytes(payload.FeeRecipient()),
		StateRoot:              bytesutil.SafeCopyBytes(payload.StateRoot()),
		ReceiptsRoot:           bytesutil.SafeCopyBytes(payload.ReceiptsRoot()),
		LogsBloom:              bytesutil.SafeCopyBytes(payload.LogsBloom()),
		PrevRandao:             bytesutil.SafeCopyBytes(payload.PrevRandao()),
		BlockNumber:            payload.BlockNumber(),
		GasLimit:               payload.GasLimit(),
		GasUsed:                payload.GasUsed(),
		Timestamp:              payload.Timestamp(),
		ExtraData:              bytesutil.SafeCopyBytes(payload.ExtraData()),
		BaseFeePerGas:          bytesutil.SafeCopyBytes(payload.BaseFeePerGas()),
		BlockHash:              bytesutil.SafeCopyBytes(payload.BlockHash()),
		TransactionsRoot:       txRoot[:],
		WithdrawalsRoot:        withdrawalsRoot[:],
		BlobGasUsed:            blobGasUsed,
		ExcessBlobGas:          excessBlobGas,
		DepositReceiptsRoot:    depositReceiptsRoot[:],
		WithdrawalRequestsRoot: withdrawalRequestsRoot[:],
	}, nil
}

// IsEmptyExecutionData checks if an execution data is empty underneath. If a single field has
// a non-zero value, this function will return false.
func IsEmptyExecutionData(data interfaces.ExecutionData) (bool, error) {
	if data == nil {
		return true, nil
	}
	if !bytes.Equal(data.ParentHash(), make([]byte, fieldparams.RootLength)) {
		return false, nil
	}
	if !bytes.Equal(data.FeeRecipient(), make([]byte, fieldparams.FeeRecipientLength)) {
		return false, nil
	}
	if !bytes.Equal(data.StateRoot(), make([]byte, fieldparams.RootLength)) {
		return false, nil
	}
	if !bytes.Equal(data.ReceiptsRoot(), make([]byte, fieldparams.RootLength)) {
		return false, nil
	}
	if !bytes.Equal(data.LogsBloom(), make([]byte, fieldparams.LogsBloomLength)) {
		return false, nil
	}
	if !bytes.Equal(data.PrevRandao(), make([]byte, fieldparams.RootLength)) {
		return false, nil
	}
	if !bytes.Equal(data.BaseFeePerGas(), make([]byte, fieldparams.RootLength)) {
		return false, nil
	}
	if !bytes.Equal(data.BlockHash(), make([]byte, fieldparams.RootLength)) {
		return false, nil
	}

	txs, err := data.Transactions()
	switch {
	case errors.Is(err, consensus_types.ErrUnsupportedField):
	case err != nil:
		return false, err
	default:
		if len(txs) != 0 {
			return false, nil
		}
	}

	if len(data.ExtraData()) != 0 {
		return false, nil
	}
	if data.BlockNumber() != 0 {
		return false, nil
	}
	if data.GasLimit() != 0 {
		return false, nil
	}
	if data.GasUsed() != 0 {
		return false, nil
	}
	if data.Timestamp() != 0 {
		return false, nil
	}

	epe, postElectra := data.(interfaces.ExecutionDataElectra)
	if postElectra {
		drs := epe.DepositReceipts()
		if len(drs) != 0 {
			return false, nil
		}

		wrs := epe.WithdrawalRequests()
		if len(wrs) != 0 {
			return false, nil
		}
	}

	return true, nil
}

// executionPayloadHeaderDeneb is a convenience wrapper around a blinded beacon block body's execution header data structure.
// This wrapper allows us to conform to a common interface so that beacon
// blocks for future forks can also be applied across Prysm without issues.
type executionPayloadHeaderDeneb struct {
	p         *enginev1.ExecutionPayloadHeaderDeneb
	weiValue  math.Wei
	gweiValue uint64
}

var _ interfaces.ExecutionData = &executionPayloadHeaderDeneb{}

// WrappedExecutionPayloadHeaderDeneb is a constructor which wraps a protobuf execution header into an interface.
func WrappedExecutionPayloadHeaderDeneb(p *enginev1.ExecutionPayloadHeaderDeneb, value math.Wei) (interfaces.ExecutionData, error) {
	w := executionPayloadHeaderDeneb{p: p, weiValue: value, gweiValue: uint64(math.WeiToGwei(value))}
	if w.IsNil() {
		return nil, consensus_types.ErrNilObjectWrapped
	}
	return w, nil
}

// IsNil checks if the underlying data is nil.
func (e executionPayloadHeaderDeneb) IsNil() bool {
	return e.p == nil
}

// MarshalSSZ --
func (e executionPayloadHeaderDeneb) MarshalSSZ() ([]byte, error) {
	return e.p.MarshalSSZ()
}

// MarshalSSZTo --
func (e executionPayloadHeaderDeneb) MarshalSSZTo(dst []byte) ([]byte, error) {
	return e.p.MarshalSSZTo(dst)
}

// SizeSSZ --
func (e executionPayloadHeaderDeneb) SizeSSZ() int {
	return e.p.SizeSSZ()
}

// UnmarshalSSZ --
func (e executionPayloadHeaderDeneb) UnmarshalSSZ(buf []byte) error {
	return e.p.UnmarshalSSZ(buf)
}

// HashTreeRoot --
func (e executionPayloadHeaderDeneb) HashTreeRoot() ([32]byte, error) {
	return e.p.HashTreeRoot()
}

// HashTreeRootWith --
func (e executionPayloadHeaderDeneb) HashTreeRootWith(hh *fastssz.Hasher) error {
	return e.p.HashTreeRootWith(hh)
}

// Proto --
func (e executionPayloadHeaderDeneb) Proto() proto.Message {
	return e.p
}

// ParentHash --
func (e executionPayloadHeaderDeneb) ParentHash() []byte {
	return e.p.ParentHash
}

// FeeRecipient --
func (e executionPayloadHeaderDeneb) FeeRecipient() []byte {
	return e.p.FeeRecipient
}

// StateRoot --
func (e executionPayloadHeaderDeneb) StateRoot() []byte {
	return e.p.StateRoot
}

// ReceiptsRoot --
func (e executionPayloadHeaderDeneb) ReceiptsRoot() []byte {
	return e.p.ReceiptsRoot
}

// LogsBloom --
func (e executionPayloadHeaderDeneb) LogsBloom() []byte {
	return e.p.LogsBloom
}

// PrevRandao --
func (e executionPayloadHeaderDeneb) PrevRandao() []byte {
	return e.p.PrevRandao
}

// BlockNumber --
func (e executionPayloadHeaderDeneb) BlockNumber() uint64 {
	return e.p.BlockNumber
}

// GasLimit --
func (e executionPayloadHeaderDeneb) GasLimit() uint64 {
	return e.p.GasLimit
}

// GasUsed --
func (e executionPayloadHeaderDeneb) GasUsed() uint64 {
	return e.p.GasUsed
}

// Timestamp --
func (e executionPayloadHeaderDeneb) Timestamp() uint64 {
	return e.p.Timestamp
}

// ExtraData --
func (e executionPayloadHeaderDeneb) ExtraData() []byte {
	return e.p.ExtraData
}

// BaseFeePerGas --
func (e executionPayloadHeaderDeneb) BaseFeePerGas() []byte {
	return e.p.BaseFeePerGas
}

// BlockHash --
func (e executionPayloadHeaderDeneb) BlockHash() []byte {
	return e.p.BlockHash
}

// Transactions --
func (executionPayloadHeaderDeneb) Transactions() ([][]byte, error) {
	return nil, consensus_types.ErrUnsupportedField
}

// TransactionsRoot --
func (e executionPayloadHeaderDeneb) TransactionsRoot() ([]byte, error) {
	return e.p.TransactionsRoot, nil
}

// Withdrawals --
func (e executionPayloadHeaderDeneb) Withdrawals() ([]*enginev1.Withdrawal, error) {
	return nil, consensus_types.ErrUnsupportedField
}

// WithdrawalsRoot --
func (e executionPayloadHeaderDeneb) WithdrawalsRoot() ([]byte, error) {
	return e.p.WithdrawalsRoot, nil
}

// BlobGasUsed --
func (e executionPayloadHeaderDeneb) BlobGasUsed() (uint64, error) {
	return e.p.BlobGasUsed, nil
}

// ExcessBlobGas --
func (e executionPayloadHeaderDeneb) ExcessBlobGas() (uint64, error) {
	return e.p.ExcessBlobGas, nil
}

// ValueInWei --
func (e executionPayloadHeaderDeneb) ValueInWei() (math.Wei, error) {
	return e.weiValue, nil
}

// ValueInGwei --
func (e executionPayloadHeaderDeneb) ValueInGwei() (uint64, error) {
	return e.gweiValue, nil
}

// IsBlinded returns true if the underlying data is blinded.
func (e executionPayloadHeaderDeneb) IsBlinded() bool {
	return true
}

// executionPayloadDeneb is a convenience wrapper around a beacon block body's execution payload data structure
// This wrapper allows us to conform to a common interface so that beacon
// blocks for future forks can also be applied across Prysm without issues.
type executionPayloadDeneb struct {
	p         *enginev1.ExecutionPayloadDeneb
	weiValue  math.Wei
	gweiValue uint64
}

var _ interfaces.ExecutionData = &executionPayloadDeneb{}

// WrappedExecutionPayloadDeneb is a constructor which wraps a protobuf execution payload into an interface.
func WrappedExecutionPayloadDeneb(p *enginev1.ExecutionPayloadDeneb, value math.Wei) (interfaces.ExecutionData, error) {
	w := executionPayloadDeneb{p: p, weiValue: value, gweiValue: uint64(math.WeiToGwei(value))}
	if w.IsNil() {
		return nil, consensus_types.ErrNilObjectWrapped
	}
	return w, nil
}

// IsNil checks if the underlying data is nil.
func (e executionPayloadDeneb) IsNil() bool {
	return e.p == nil
}

// MarshalSSZ --
func (e executionPayloadDeneb) MarshalSSZ() ([]byte, error) {
	return e.p.MarshalSSZ()
}

// MarshalSSZTo --
func (e executionPayloadDeneb) MarshalSSZTo(dst []byte) ([]byte, error) {
	return e.p.MarshalSSZTo(dst)
}

// SizeSSZ --
func (e executionPayloadDeneb) SizeSSZ() int {
	return e.p.SizeSSZ()
}

// UnmarshalSSZ --
func (e executionPayloadDeneb) UnmarshalSSZ(buf []byte) error {
	return e.p.UnmarshalSSZ(buf)
}

// HashTreeRoot --
func (e executionPayloadDeneb) HashTreeRoot() ([32]byte, error) {
	return e.p.HashTreeRoot()
}

// HashTreeRootWith --
func (e executionPayloadDeneb) HashTreeRootWith(hh *fastssz.Hasher) error {
	return e.p.HashTreeRootWith(hh)
}

// Proto --
func (e executionPayloadDeneb) Proto() proto.Message {
	return e.p
}

// ParentHash --
func (e executionPayloadDeneb) ParentHash() []byte {
	return e.p.ParentHash
}

// FeeRecipient --
func (e executionPayloadDeneb) FeeRecipient() []byte {
	return e.p.FeeRecipient
}

// StateRoot --
func (e executionPayloadDeneb) StateRoot() []byte {
	return e.p.StateRoot
}

// ReceiptsRoot --
func (e executionPayloadDeneb) ReceiptsRoot() []byte {
	return e.p.ReceiptsRoot
}

// LogsBloom --
func (e executionPayloadDeneb) LogsBloom() []byte {
	return e.p.LogsBloom
}

// PrevRandao --
func (e executionPayloadDeneb) PrevRandao() []byte {
	return e.p.PrevRandao
}

// BlockNumber --
func (e executionPayloadDeneb) BlockNumber() uint64 {
	return e.p.BlockNumber
}

// GasLimit --
func (e executionPayloadDeneb) GasLimit() uint64 {
	return e.p.GasLimit
}

// GasUsed --
func (e executionPayloadDeneb) GasUsed() uint64 {
	return e.p.GasUsed
}

// Timestamp --
func (e executionPayloadDeneb) Timestamp() uint64 {
	return e.p.Timestamp
}

// ExtraData --
func (e executionPayloadDeneb) ExtraData() []byte {
	return e.p.ExtraData
}

// BaseFeePerGas --
func (e executionPayloadDeneb) BaseFeePerGas() []byte {
	return e.p.BaseFeePerGas
}

// BlockHash --
func (e executionPayloadDeneb) BlockHash() []byte {
	return e.p.BlockHash
}

// Transactions --
func (e executionPayloadDeneb) Transactions() ([][]byte, error) {
	return e.p.Transactions, nil
}

// TransactionsRoot --
func (e executionPayloadDeneb) TransactionsRoot() ([]byte, error) {
	return nil, consensus_types.ErrUnsupportedField
}

// Withdrawals --
func (e executionPayloadDeneb) Withdrawals() ([]*enginev1.Withdrawal, error) {
	return e.p.Withdrawals, nil
}

// WithdrawalsRoot --
func (e executionPayloadDeneb) WithdrawalsRoot() ([]byte, error) {
	return nil, consensus_types.ErrUnsupportedField
}

func (e executionPayloadDeneb) BlobGasUsed() (uint64, error) {
	return e.p.BlobGasUsed, nil
}

func (e executionPayloadDeneb) ExcessBlobGas() (uint64, error) {
	return e.p.ExcessBlobGas, nil
}

// ValueInWei --
func (e executionPayloadDeneb) ValueInWei() (math.Wei, error) {
	return e.weiValue, nil
}

// ValueInGwei --
func (e executionPayloadDeneb) ValueInGwei() (uint64, error) {
	return e.gweiValue, nil
}

// IsBlinded returns true if the underlying data is blinded.
func (e executionPayloadDeneb) IsBlinded() bool {
	return false
}

// executionPayloadHeaderElectra is a convenience wrapper around a blinded beacon block body's execution header data structure.
// This wrapper allows us to conform to a common interface so that beacon
// blocks for future forks can also be applied across Prysm without issues.
type executionPayloadHeaderElectra struct {
	p         *enginev1.ExecutionPayloadHeaderElectra
	weiValue  math.Wei
	gweiValue uint64
}

var _ interfaces.ExecutionData = &executionPayloadElectra{}
var _ interfaces.ExecutionDataElectra = &executionPayloadElectra{}

// WrappedExecutionPayloadHeaderElectra is a constructor which wraps a protobuf execution header into an interface.
func WrappedExecutionPayloadHeaderElectra(p *enginev1.ExecutionPayloadHeaderElectra, value math.Wei) (interfaces.ExecutionData, error) {
	w := executionPayloadHeaderElectra{p: p, weiValue: value, gweiValue: uint64(math.WeiToGwei(value))}
	if w.IsNil() {
		return nil, consensus_types.ErrNilObjectWrapped
	}
	return w, nil
}

// IsNil checks if the underlying data is nil.
func (e executionPayloadHeaderElectra) IsNil() bool {
	return e.p == nil
}

// MarshalSSZ --
func (e executionPayloadHeaderElectra) MarshalSSZ() ([]byte, error) {
	return e.p.MarshalSSZ()
}

// MarshalSSZTo --
func (e executionPayloadHeaderElectra) MarshalSSZTo(dst []byte) ([]byte, error) {
	return e.p.MarshalSSZTo(dst)
}

// SizeSSZ --
func (e executionPayloadHeaderElectra) SizeSSZ() int {
	return e.p.SizeSSZ()
}

// UnmarshalSSZ --
func (e executionPayloadHeaderElectra) UnmarshalSSZ(buf []byte) error {
	return e.p.UnmarshalSSZ(buf)
}

// HashTreeRoot --
func (e executionPayloadHeaderElectra) HashTreeRoot() ([32]byte, error) {
	return e.p.HashTreeRoot()
}

// HashTreeRootWith --
func (e executionPayloadHeaderElectra) HashTreeRootWith(hh *fastssz.Hasher) error {
	return e.p.HashTreeRootWith(hh)
}

// Proto --
func (e executionPayloadHeaderElectra) Proto() proto.Message {
	return e.p
}

// ParentHash --
func (e executionPayloadHeaderElectra) ParentHash() []byte {
	return e.p.ParentHash
}

// FeeRecipient --
func (e executionPayloadHeaderElectra) FeeRecipient() []byte {
	return e.p.FeeRecipient
}

// StateRoot --
func (e executionPayloadHeaderElectra) StateRoot() []byte {
	return e.p.StateRoot
}

// ReceiptsRoot --
func (e executionPayloadHeaderElectra) ReceiptsRoot() []byte {
	return e.p.ReceiptsRoot
}

// LogsBloom --
func (e executionPayloadHeaderElectra) LogsBloom() []byte {
	return e.p.LogsBloom
}

// PrevRandao --
func (e executionPayloadHeaderElectra) PrevRandao() []byte {
	return e.p.PrevRandao
}

// BlockNumber --
func (e executionPayloadHeaderElectra) BlockNumber() uint64 {
	return e.p.BlockNumber
}

// GasLimit --
func (e executionPayloadHeaderElectra) GasLimit() uint64 {
	return e.p.GasLimit
}

// GasUsed --
func (e executionPayloadHeaderElectra) GasUsed() uint64 {
	return e.p.GasUsed
}

// Timestamp --
func (e executionPayloadHeaderElectra) Timestamp() uint64 {
	return e.p.Timestamp
}

// ExtraData --
func (e executionPayloadHeaderElectra) ExtraData() []byte {
	return e.p.ExtraData
}

// BaseFeePerGas --
func (e executionPayloadHeaderElectra) BaseFeePerGas() []byte {
	return e.p.BaseFeePerGas
}

// BlockHash --
func (e executionPayloadHeaderElectra) BlockHash() []byte {
	return e.p.BlockHash
}

// Transactions --
func (executionPayloadHeaderElectra) Transactions() ([][]byte, error) {
	return nil, consensus_types.ErrUnsupportedField
}

// TransactionsRoot --
func (e executionPayloadHeaderElectra) TransactionsRoot() ([]byte, error) {
	return e.p.TransactionsRoot, nil
}

// Withdrawals --
func (e executionPayloadHeaderElectra) Withdrawals() ([]*enginev1.Withdrawal, error) {
	return nil, consensus_types.ErrUnsupportedField
}

// WithdrawalsRoot --
func (e executionPayloadHeaderElectra) WithdrawalsRoot() ([]byte, error) {
	return e.p.WithdrawalsRoot, nil
}

// BlobGasUsed --
func (e executionPayloadHeaderElectra) BlobGasUsed() (uint64, error) {
	return e.p.BlobGasUsed, nil
}

// ExcessBlobGas --
func (e executionPayloadHeaderElectra) ExcessBlobGas() (uint64, error) {
	return e.p.ExcessBlobGas, nil
}

// ValueInWei --
func (e executionPayloadHeaderElectra) ValueInWei() (math.Wei, error) {
	return e.weiValue, nil
}

// ValueInGwei --
func (e executionPayloadHeaderElectra) ValueInGwei() (uint64, error) {
	return e.gweiValue, nil
}

// DepositReceipts --
func (e executionPayloadHeaderElectra) DepositReceipts() ([]*enginev1.DepositReceipt, error) {
	return nil, consensus_types.ErrUnsupportedField
}

// WithdrawalRequests --
func (e executionPayloadHeaderElectra) WithdrawalRequests() ([]*enginev1.ExecutionLayerWithdrawalRequest, error) {
	return nil, consensus_types.ErrUnsupportedField
}

// IsBlinded returns true if the underlying data is blinded.
func (e executionPayloadHeaderElectra) IsBlinded() bool {
	return true
}

// executionPayloadElectra is a convenience wrapper around a beacon block body's execution payload data structure
// This wrapper allows us to conform to a common interface so that beacon
// blocks for future forks can also be applied across Prysm without issues.
type executionPayloadElectra struct {
	p         *enginev1.ExecutionPayloadElectra
	weiValue  math.Wei
	gweiValue uint64
}

// WrappedExecutionPayloadElectra is a constructor which wraps a protobuf execution payload into an interface.
func WrappedExecutionPayloadElectra(p *enginev1.ExecutionPayloadElectra, value math.Wei) (interfaces.ExecutionDataElectra, error) {
	w := executionPayloadElectra{p: p, weiValue: value, gweiValue: uint64(math.WeiToGwei(value))}
	if w.IsNil() {
		return nil, consensus_types.ErrNilObjectWrapped
	}
	return w, nil
}

var _ interfaces.ExecutionData = &executionPayloadElectra{}
var _ interfaces.ExecutionDataElectra = &executionPayloadElectra{}

// IsNil checks if the underlying data is nil.
func (e executionPayloadElectra) IsNil() bool {
	return e.p == nil
}

// MarshalSSZ --
func (e executionPayloadElectra) MarshalSSZ() ([]byte, error) {
	return e.p.MarshalSSZ()
}

// MarshalSSZTo --
func (e executionPayloadElectra) MarshalSSZTo(dst []byte) ([]byte, error) {
	return e.p.MarshalSSZTo(dst)
}

// SizeSSZ --
func (e executionPayloadElectra) SizeSSZ() int {
	return e.p.SizeSSZ()
}

// UnmarshalSSZ --
func (e executionPayloadElectra) UnmarshalSSZ(buf []byte) error {
	return e.p.UnmarshalSSZ(buf)
}

// HashTreeRoot --
func (e executionPayloadElectra) HashTreeRoot() ([32]byte, error) {
	return e.p.HashTreeRoot()
}

// HashTreeRootWith --
func (e executionPayloadElectra) HashTreeRootWith(hh *fastssz.Hasher) error {
	return e.p.HashTreeRootWith(hh)
}

// Proto --
func (e executionPayloadElectra) Proto() proto.Message {
	return e.p
}

// ParentHash --
func (e executionPayloadElectra) ParentHash() []byte {
	return e.p.ParentHash
}

// FeeRecipient --
func (e executionPayloadElectra) FeeRecipient() []byte {
	return e.p.FeeRecipient
}

// StateRoot --
func (e executionPayloadElectra) StateRoot() []byte {
	return e.p.StateRoot
}

// ReceiptsRoot --
func (e executionPayloadElectra) ReceiptsRoot() []byte {
	return e.p.ReceiptsRoot
}

// LogsBloom --
func (e executionPayloadElectra) LogsBloom() []byte {
	return e.p.LogsBloom
}

// PrevRandao --
func (e executionPayloadElectra) PrevRandao() []byte {
	return e.p.PrevRandao
}

// BlockNumber --
func (e executionPayloadElectra) BlockNumber() uint64 {
	return e.p.BlockNumber
}

// GasLimit --
func (e executionPayloadElectra) GasLimit() uint64 {
	return e.p.GasLimit
}

// GasUsed --
func (e executionPayloadElectra) GasUsed() uint64 {
	return e.p.GasUsed
}

// Timestamp --
func (e executionPayloadElectra) Timestamp() uint64 {
	return e.p.Timestamp
}

// ExtraData --
func (e executionPayloadElectra) ExtraData() []byte {
	return e.p.ExtraData
}

// BaseFeePerGas --
func (e executionPayloadElectra) BaseFeePerGas() []byte {
	return e.p.BaseFeePerGas
}

// BlockHash --
func (e executionPayloadElectra) BlockHash() []byte {
	return e.p.BlockHash
}

// Transactions --
func (e executionPayloadElectra) Transactions() ([][]byte, error) {
	return e.p.Transactions, nil
}

// TransactionsRoot --
func (e executionPayloadElectra) TransactionsRoot() ([]byte, error) {
	return nil, consensus_types.ErrUnsupportedField
}

// Withdrawals --
func (e executionPayloadElectra) Withdrawals() ([]*enginev1.Withdrawal, error) {
	return e.p.Withdrawals, nil
}

// WithdrawalsRoot --
func (e executionPayloadElectra) WithdrawalsRoot() ([]byte, error) {
	return nil, consensus_types.ErrUnsupportedField
}

func (e executionPayloadElectra) BlobGasUsed() (uint64, error) {
	return e.p.BlobGasUsed, nil
}

func (e executionPayloadElectra) ExcessBlobGas() (uint64, error) {
	return e.p.ExcessBlobGas, nil
}

// ValueInWei --
func (e executionPayloadElectra) ValueInWei() (math.Wei, error) {
	return e.weiValue, nil
}

// ValueInGwei --
func (e executionPayloadElectra) ValueInGwei() (uint64, error) {
	return e.gweiValue, nil
}

// DepositReceipts --
func (e executionPayloadElectra) DepositReceipts() []*enginev1.DepositReceipt {
	return e.p.DepositReceipts
}

// WithdrawalRequests --
func (e executionPayloadElectra) WithdrawalRequests() []*enginev1.ExecutionLayerWithdrawalRequest {
	return e.p.WithdrawalRequests
}

// IsBlinded returns true if the underlying data is blinded.
func (e executionPayloadElectra) IsBlinded() bool {
	return false
}

// PayloadValueToWei returns a Wei value given the payload's value
func PayloadValueToWei(value []byte) math.Wei {
	// We have to convert big endian to little endian because the value is coming from the execution layer.
	return big.NewInt(0).SetBytes(bytesutil.ReverseByteOrder(value))
}

// PayloadValueToGwei returns a Gwei value given the payload's value
func PayloadValueToGwei(value []byte) math.Gwei {
	// We have to convert big endian to little endian because the value is coming from the execution layer.
	v := big.NewInt(0).SetBytes(bytesutil.ReverseByteOrder(value))
	return math.WeiToGwei(v)
}
