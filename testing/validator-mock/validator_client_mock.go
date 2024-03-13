// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/prysmaticlabs/prysm/v5/validator/client/iface (interfaces: ValidatorClient)
//
// Generated by this command:
//
//	mockgen -package=validator_mock -destination=testing/validator-mock/validator_client_mock.go github.com/prysmaticlabs/prysm/v5/validator/client/iface ValidatorClient
//

// Package validator_mock is a generated GoMock package.
package validator_mock

import (
	context "context"
	reflect "reflect"

	"github.com/prysmaticlabs/prysm/v5/api/client/beacon"
	"github.com/prysmaticlabs/prysm/v5/api/client/event"
	eth "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
	iface "github.com/prysmaticlabs/prysm/v5/validator/client/iface"
	gomock "go.uber.org/mock/gomock"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// MockValidatorClient is a mock of ValidatorClient interface.
type MockValidatorClient struct {
	ctrl     *gomock.Controller
	recorder *MockValidatorClientMockRecorder
}

// MockValidatorClientMockRecorder is the mock recorder for MockValidatorClient.
type MockValidatorClientMockRecorder struct {
	mock *MockValidatorClient
}

// NewMockValidatorClient creates a new mock instance.
func NewMockValidatorClient(ctrl *gomock.Controller) *MockValidatorClient {
	mock := &MockValidatorClient{ctrl: ctrl}
	mock.recorder = &MockValidatorClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockValidatorClient) EXPECT() *MockValidatorClientMockRecorder {
	return m.recorder
}

// CheckDoppelGanger mocks base method.
func (m *MockValidatorClient) CheckDoppelGanger(arg0 context.Context, arg1 *eth.DoppelGangerRequest) (*eth.DoppelGangerResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckDoppelGanger", arg0, arg1)
	ret0, _ := ret[0].(*eth.DoppelGangerResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckDoppelGanger indicates an expected call of CheckDoppelGanger.
func (mr *MockValidatorClientMockRecorder) CheckDoppelGanger(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckDoppelGanger", reflect.TypeOf((*MockValidatorClient)(nil).CheckDoppelGanger), arg0, arg1)
}

// DomainData mocks base method.
func (m *MockValidatorClient) DomainData(arg0 context.Context, arg1 *eth.DomainRequest) (*eth.DomainResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DomainData", arg0, arg1)
	ret0, _ := ret[0].(*eth.DomainResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DomainData indicates an expected call of DomainData.
func (mr *MockValidatorClientMockRecorder) DomainData(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DomainData", reflect.TypeOf((*MockValidatorClient)(nil).DomainData), arg0, arg1)
}

// EventStreamIsRunning mocks base method.
func (m *MockValidatorClient) EventStreamIsRunning() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EventStreamIsRunning")
	ret0, _ := ret[0].(bool)
	return ret0
}

// EventStreamIsRunning indicates an expected call of EventStreamIsRunning.
func (mr *MockValidatorClientMockRecorder) EventStreamIsRunning() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EventStreamIsRunning", reflect.TypeOf((*MockValidatorClient)(nil).EventStreamIsRunning))
}

// GetAggregatedSelections mocks base method.
func (m *MockValidatorClient) GetAggregatedSelections(arg0 context.Context, arg1 []iface.BeaconCommitteeSelection) ([]iface.BeaconCommitteeSelection, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAggregatedSelections", arg0, arg1)
	ret0, _ := ret[0].([]iface.BeaconCommitteeSelection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAggregatedSelections indicates an expected call of GetAggregatedSelections.
func (mr *MockValidatorClientMockRecorder) GetAggregatedSelections(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAggregatedSelections", reflect.TypeOf((*MockValidatorClient)(nil).GetAggregatedSelections), arg0, arg1)
}

// GetAggregatedSyncSelections mocks base method.
func (m *MockValidatorClient) GetAggregatedSyncSelections(arg0 context.Context, arg1 []iface.SyncCommitteeSelection) ([]iface.SyncCommitteeSelection, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAggregatedSyncSelections", arg0, arg1)
	ret0, _ := ret[0].([]iface.SyncCommitteeSelection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAggregatedSyncSelections indicates an expected call of GetAggregatedSyncSelections.
func (mr *MockValidatorClientMockRecorder) GetAggregatedSyncSelections(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAggregatedSyncSelections", reflect.TypeOf((*MockValidatorClient)(nil).GetAggregatedSyncSelections), arg0, arg1)
}

// GetAttestationData mocks base method.
func (m *MockValidatorClient) GetAttestationData(arg0 context.Context, arg1 *eth.AttestationDataRequest) (*eth.AttestationData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAttestationData", arg0, arg1)
	ret0, _ := ret[0].(*eth.AttestationData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAttestationData indicates an expected call of GetAttestationData.
func (mr *MockValidatorClientMockRecorder) GetAttestationData(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAttestationData", reflect.TypeOf((*MockValidatorClient)(nil).GetAttestationData), arg0, arg1)
}

// GetBeaconBlock mocks base method.
func (m *MockValidatorClient) GetBeaconBlock(arg0 context.Context, arg1 *eth.BlockRequest) (*eth.GenericBeaconBlock, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBeaconBlock", arg0, arg1)
	ret0, _ := ret[0].(*eth.GenericBeaconBlock)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBeaconBlock indicates an expected call of GetBeaconBlock.
func (mr *MockValidatorClientMockRecorder) GetBeaconBlock(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBeaconBlock", reflect.TypeOf((*MockValidatorClient)(nil).GetBeaconBlock), arg0, arg1)
}

// GetDuties mocks base method.
func (m *MockValidatorClient) GetDuties(arg0 context.Context, arg1 *eth.DutiesRequest) (*eth.DutiesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDuties", arg0, arg1)
	ret0, _ := ret[0].(*eth.DutiesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDuties indicates an expected call of GetDuties.
func (mr *MockValidatorClientMockRecorder) GetDuties(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDuties", reflect.TypeOf((*MockValidatorClient)(nil).GetDuties), arg0, arg1)
}

// GetFeeRecipientByPubKey mocks base method.
func (m *MockValidatorClient) GetFeeRecipientByPubKey(arg0 context.Context, arg1 *eth.FeeRecipientByPubKeyRequest) (*eth.FeeRecipientByPubKeyResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFeeRecipientByPubKey", arg0, arg1)
	ret0, _ := ret[0].(*eth.FeeRecipientByPubKeyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFeeRecipientByPubKey indicates an expected call of GetFeeRecipientByPubKey.
func (mr *MockValidatorClientMockRecorder) GetFeeRecipientByPubKey(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeeRecipientByPubKey", reflect.TypeOf((*MockValidatorClient)(nil).GetFeeRecipientByPubKey), arg0, arg1)
}

// GetSyncCommitteeContribution mocks base method.
func (m *MockValidatorClient) GetSyncCommitteeContribution(arg0 context.Context, arg1 *eth.SyncCommitteeContributionRequest) (*eth.SyncCommitteeContribution, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSyncCommitteeContribution", arg0, arg1)
	ret0, _ := ret[0].(*eth.SyncCommitteeContribution)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSyncCommitteeContribution indicates an expected call of GetSyncCommitteeContribution.
func (mr *MockValidatorClientMockRecorder) GetSyncCommitteeContribution(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSyncCommitteeContribution", reflect.TypeOf((*MockValidatorClient)(nil).GetSyncCommitteeContribution), arg0, arg1)
}

// GetSyncMessageBlockRoot mocks base method.
func (m *MockValidatorClient) GetSyncMessageBlockRoot(arg0 context.Context, arg1 *emptypb.Empty) (*eth.SyncMessageBlockRootResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSyncMessageBlockRoot", arg0, arg1)
	ret0, _ := ret[0].(*eth.SyncMessageBlockRootResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSyncMessageBlockRoot indicates an expected call of GetSyncMessageBlockRoot.
func (mr *MockValidatorClientMockRecorder) GetSyncMessageBlockRoot(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSyncMessageBlockRoot", reflect.TypeOf((*MockValidatorClient)(nil).GetSyncMessageBlockRoot), arg0, arg1)
}

// GetSyncSubcommitteeIndex mocks base method.
func (m *MockValidatorClient) GetSyncSubcommitteeIndex(arg0 context.Context, arg1 *eth.SyncSubcommitteeIndexRequest) (*eth.SyncSubcommitteeIndexResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSyncSubcommitteeIndex", arg0, arg1)
	ret0, _ := ret[0].(*eth.SyncSubcommitteeIndexResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSyncSubcommitteeIndex indicates an expected call of GetSyncSubcommitteeIndex.
func (mr *MockValidatorClientMockRecorder) GetSyncSubcommitteeIndex(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSyncSubcommitteeIndex", reflect.TypeOf((*MockValidatorClient)(nil).GetSyncSubcommitteeIndex), arg0, arg1)
}

// MultipleValidatorStatus mocks base method.
func (m *MockValidatorClient) MultipleValidatorStatus(arg0 context.Context, arg1 *eth.MultipleValidatorStatusRequest) (*eth.MultipleValidatorStatusResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MultipleValidatorStatus", arg0, arg1)
	ret0, _ := ret[0].(*eth.MultipleValidatorStatusResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MultipleValidatorStatus indicates an expected call of MultipleValidatorStatus.
func (mr *MockValidatorClientMockRecorder) MultipleValidatorStatus(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MultipleValidatorStatus", reflect.TypeOf((*MockValidatorClient)(nil).MultipleValidatorStatus), arg0, arg1)
}

// PrepareBeaconProposer mocks base method.
func (m *MockValidatorClient) PrepareBeaconProposer(arg0 context.Context, arg1 *eth.PrepareBeaconProposerRequest) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrepareBeaconProposer", arg0, arg1)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PrepareBeaconProposer indicates an expected call of PrepareBeaconProposer.
func (mr *MockValidatorClientMockRecorder) PrepareBeaconProposer(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrepareBeaconProposer", reflect.TypeOf((*MockValidatorClient)(nil).PrepareBeaconProposer), arg0, arg1)
}

// ProposeAttestation mocks base method.
func (m *MockValidatorClient) ProposeAttestation(arg0 context.Context, arg1 *eth.Attestation) (*eth.AttestResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProposeAttestation", arg0, arg1)
	ret0, _ := ret[0].(*eth.AttestResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProposeAttestation indicates an expected call of ProposeAttestation.
func (mr *MockValidatorClientMockRecorder) ProposeAttestation(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProposeAttestation", reflect.TypeOf((*MockValidatorClient)(nil).ProposeAttestation), arg0, arg1)
}

// ProposeBeaconBlock mocks base method.
func (m *MockValidatorClient) ProposeBeaconBlock(arg0 context.Context, arg1 *eth.GenericSignedBeaconBlock) (*eth.ProposeResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProposeBeaconBlock", arg0, arg1)
	ret0, _ := ret[0].(*eth.ProposeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProposeBeaconBlock indicates an expected call of ProposeBeaconBlock.
func (mr *MockValidatorClientMockRecorder) ProposeBeaconBlock(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProposeBeaconBlock", reflect.TypeOf((*MockValidatorClient)(nil).ProposeBeaconBlock), arg0, arg1)
}

// ProposeExit mocks base method.
func (m *MockValidatorClient) ProposeExit(arg0 context.Context, arg1 *eth.SignedVoluntaryExit) (*eth.ProposeExitResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProposeExit", arg0, arg1)
	ret0, _ := ret[0].(*eth.ProposeExitResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProposeExit indicates an expected call of ProposeExit.
func (mr *MockValidatorClientMockRecorder) ProposeExit(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProposeExit", reflect.TypeOf((*MockValidatorClient)(nil).ProposeExit), arg0, arg1)
}

// StartEventStream mocks base method.
func (m *MockValidatorClient) StartEventStream(arg0 context.Context, arg1 []string, arg2 chan<- *event.Event){
	m.ctrl.T.Helper()
	_ = m.ctrl.Call(m, "StartEventStream", arg0,arg1,arg2)
}

// StartEventStream indicates an expected call of StartEventStream.
func (mr *MockValidatorClientMockRecorder) StartEventStream(arg0,arg1,arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartEventStream", reflect.TypeOf((*MockValidatorClient)(nil).StartEventStream), arg0, arg1, arg2)
}

// ProcessEvent mocks base method.
func (m *MockValidatorClient) ProcessEvent(arg0 *event.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessEvent", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProcessEvent indicates an expected call of ProcessEvent.
func (mr *MockValidatorClientMockRecorder) ProcessEvent(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessEvent", reflect.TypeOf((*MockValidatorClient)(nil).ProcessEvent), arg0)
}

// NodeIsHealthy mocks base method.
func (m *MockValidatorClient) NodeIsHealthy(arg0 context.Context) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NodeIsHealthy",arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// NodeIsHealthy indicates an expected call of NodeIsHealthy.
func (mr *MockValidatorClientMockRecorder) NodeIsHealthy(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NodeIsHealthy", reflect.TypeOf((*MockValidatorClient)(nil).NodeIsHealthy), arg0)
}

// NodeHealthTracker mocks base method.
func (m *MockValidatorClient) NodeHealthTracker() *beacon.NodeHealthTracker {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NodeHealthTracker")
	ret0, _ := ret[0].(*beacon.NodeHealthTracker)
	return ret0
}

// NodeHealthTracker indicates an expected call of NodeHealthTracker.
func (mr *MockValidatorClientMockRecorder) NodeHealthTracker() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NodeHealthTracker", reflect.TypeOf((*MockValidatorClient)(nil).NodeHealthTracker))
}

// SubmitAggregateSelectionProof mocks base method.
func (m *MockValidatorClient) SubmitAggregateSelectionProof(arg0 context.Context, arg1 *eth.AggregateSelectionRequest) (*eth.AggregateSelectionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubmitAggregateSelectionProof", arg0, arg1)
	ret0, _ := ret[0].(*eth.AggregateSelectionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubmitAggregateSelectionProof indicates an expected call of SubmitAggregateSelectionProof.
func (mr *MockValidatorClientMockRecorder) SubmitAggregateSelectionProof(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitAggregateSelectionProof", reflect.TypeOf((*MockValidatorClient)(nil).SubmitAggregateSelectionProof), arg0, arg1)
}

// SubmitSignedAggregateSelectionProof mocks base method.
func (m *MockValidatorClient) SubmitSignedAggregateSelectionProof(arg0 context.Context, arg1 *eth.SignedAggregateSubmitRequest) (*eth.SignedAggregateSubmitResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubmitSignedAggregateSelectionProof", arg0, arg1)
	ret0, _ := ret[0].(*eth.SignedAggregateSubmitResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubmitSignedAggregateSelectionProof indicates an expected call of SubmitSignedAggregateSelectionProof.
func (mr *MockValidatorClientMockRecorder) SubmitSignedAggregateSelectionProof(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitSignedAggregateSelectionProof", reflect.TypeOf((*MockValidatorClient)(nil).SubmitSignedAggregateSelectionProof), arg0, arg1)
}

// SubmitSignedContributionAndProof mocks base method.
func (m *MockValidatorClient) SubmitSignedContributionAndProof(arg0 context.Context, arg1 *eth.SignedContributionAndProof) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubmitSignedContributionAndProof", arg0, arg1)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubmitSignedContributionAndProof indicates an expected call of SubmitSignedContributionAndProof.
func (mr *MockValidatorClientMockRecorder) SubmitSignedContributionAndProof(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitSignedContributionAndProof", reflect.TypeOf((*MockValidatorClient)(nil).SubmitSignedContributionAndProof), arg0, arg1)
}

// SubmitSyncMessage mocks base method.
func (m *MockValidatorClient) SubmitSyncMessage(arg0 context.Context, arg1 *eth.SyncCommitteeMessage) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubmitSyncMessage", arg0, arg1)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubmitSyncMessage indicates an expected call of SubmitSyncMessage.
func (mr *MockValidatorClientMockRecorder) SubmitSyncMessage(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitSyncMessage", reflect.TypeOf((*MockValidatorClient)(nil).SubmitSyncMessage), arg0, arg1)
}

// SubmitValidatorRegistrations mocks base method.
func (m *MockValidatorClient) SubmitValidatorRegistrations(arg0 context.Context, arg1 *eth.SignedValidatorRegistrationsV1) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubmitValidatorRegistrations", arg0, arg1)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubmitValidatorRegistrations indicates an expected call of SubmitValidatorRegistrations.
func (mr *MockValidatorClientMockRecorder) SubmitValidatorRegistrations(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitValidatorRegistrations", reflect.TypeOf((*MockValidatorClient)(nil).SubmitValidatorRegistrations), arg0, arg1)
}

// SubscribeCommitteeSubnets mocks base method.
func (m *MockValidatorClient) SubscribeCommitteeSubnets(arg0 context.Context, arg1 *eth.CommitteeSubnetsSubscribeRequest, arg2 []*eth.DutiesResponse_Duty) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeCommitteeSubnets", arg0, arg1, arg2)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubscribeCommitteeSubnets indicates an expected call of SubscribeCommitteeSubnets.
func (mr *MockValidatorClientMockRecorder) SubscribeCommitteeSubnets(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeCommitteeSubnets", reflect.TypeOf((*MockValidatorClient)(nil).SubscribeCommitteeSubnets), arg0, arg1, arg2)
}

// ValidatorIndex mocks base method.
func (m *MockValidatorClient) ValidatorIndex(arg0 context.Context, arg1 *eth.ValidatorIndexRequest) (*eth.ValidatorIndexResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidatorIndex", arg0, arg1)
	ret0, _ := ret[0].(*eth.ValidatorIndexResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidatorIndex indicates an expected call of ValidatorIndex.
func (mr *MockValidatorClientMockRecorder) ValidatorIndex(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidatorIndex", reflect.TypeOf((*MockValidatorClient)(nil).ValidatorIndex), arg0, arg1)
}

// ValidatorStatus mocks base method.
func (m *MockValidatorClient) ValidatorStatus(arg0 context.Context, arg1 *eth.ValidatorStatusRequest) (*eth.ValidatorStatusResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidatorStatus", arg0, arg1)
	ret0, _ := ret[0].(*eth.ValidatorStatusResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidatorStatus indicates an expected call of ValidatorStatus.
func (mr *MockValidatorClientMockRecorder) ValidatorStatus(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidatorStatus", reflect.TypeOf((*MockValidatorClient)(nil).ValidatorStatus), arg0, arg1)
}

// WaitForActivation mocks base method.
func (m *MockValidatorClient) WaitForActivation(arg0 context.Context, arg1 *eth.ValidatorActivationRequest) (eth.BeaconNodeValidator_WaitForActivationClient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitForActivation", arg0, arg1)
	ret0, _ := ret[0].(eth.BeaconNodeValidator_WaitForActivationClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WaitForActivation indicates an expected call of WaitForActivation.
func (mr *MockValidatorClientMockRecorder) WaitForActivation(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitForActivation", reflect.TypeOf((*MockValidatorClient)(nil).WaitForActivation), arg0, arg1)
}

// WaitForChainStart mocks base method.
func (m *MockValidatorClient) WaitForChainStart(arg0 context.Context, arg1 *emptypb.Empty) (*eth.ChainStartResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitForChainStart", arg0, arg1)
	ret0, _ := ret[0].(*eth.ChainStartResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WaitForChainStart indicates an expected call of WaitForChainStart.
func (mr *MockValidatorClientMockRecorder) WaitForChainStart(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitForChainStart", reflect.TypeOf((*MockValidatorClient)(nil).WaitForChainStart), arg0, arg1)
}
