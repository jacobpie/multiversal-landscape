package sync

import (
	"context"
	"errors"
	"fmt"

	"github.com/prysmaticlabs/prysm/v5/beacon-chain/core/helpers"
	ethpb "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
	"google.golang.org/protobuf/proto"
)

// beaconAggregateProofSubscriber forwards the incoming validated aggregated attestation and proof to the
// attestation pool for processing.
func (s *Service) beaconAggregateProofSubscriber(_ context.Context, msg proto.Message) error {
	a, ok := msg.(ethpb.SignedAggregateAttAndProof)
	if !ok {
		return fmt.Errorf("message was not type ethpb.SignedAggregateAttAndProof, type=%T", msg)
	}

	if a.GetAggregateAttestationAndProof().GetAggregateVal() == nil || a.GetAggregateAttestationAndProof().GetAggregateVal().GetData() == nil {
		return errors.New("nil aggregate")
	}

	// An unaggregated attestation can make it here. It’s valid, the aggregator it just itself, although it means poor performance for the subnet.
	if !helpers.IsAggregated(a.GetAggregateAttestationAndProof().GetAggregateVal()) {
		return s.cfg.attPool.SaveUnaggregatedAttestation(a.GetAggregateAttestationAndProof().GetAggregateVal())
	}

	return s.cfg.attPool.SaveAggregatedAttestation(a.GetAggregateAttestationAndProof().GetAggregateVal())
}
