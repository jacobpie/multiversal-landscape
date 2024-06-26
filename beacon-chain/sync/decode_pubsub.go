package sync

import (
	"reflect"
	"strings"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/pkg/errors"
	ssz "github.com/prysmaticlabs/fastssz"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/blockchain"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/core/signing"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/p2p"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/p2p/types"
	"github.com/prysmaticlabs/prysm/v5/config/params"
	"github.com/prysmaticlabs/prysm/v5/encoding/bytesutil"
	ethpb "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
	"google.golang.org/protobuf/proto"
)

var errNilPubsubMessage = errors.New("nil pubsub message")
var errInvalidTopic = errors.New("invalid topic format")

func (s *Service) decodePubsubMessage(msg *pubsub.Message) (ssz.Unmarshaler, error) {
	if msg == nil || msg.Topic == nil || *msg.Topic == "" {
		return nil, errNilPubsubMessage
	}
	topic := *msg.Topic
	fDigest, err := p2p.ExtractGossipDigest(topic)
	if err != nil {
		return nil, errors.Wrapf(err, "extraction failed for topic: %s", topic)
	}
	topic = strings.TrimSuffix(topic, s.cfg.p2p.Encoding().ProtocolSuffix())
	topic, err = s.replaceForkDigest(topic)
	if err != nil {
		return nil, err
	}
	// Specially handle subnet messages.
	switch {
	case strings.Contains(topic, p2p.GossipAttestationMessage):
		topic = p2p.GossipTypeMapping[reflect.TypeOf(&ethpb.Attestation{})]
		// Given that both sync message related subnets have the same message name, we have to
		// differentiate them below.
	case strings.Contains(topic, p2p.GossipSyncCommitteeMessage) && !strings.Contains(topic, p2p.SyncContributionAndProofSubnetTopicFormat):
		topic = p2p.GossipTypeMapping[reflect.TypeOf(&ethpb.SyncCommitteeMessage{})]
	case strings.Contains(topic, p2p.GossipBlobSidecarMessage):
		topic = p2p.GossipTypeMapping[reflect.TypeOf(&ethpb.BlobSidecar{})]
	}

	base := p2p.GossipTopicMappings(topic, 0)
	if base == nil {
		return nil, p2p.ErrMessageNotMapped
	}
	m, ok := proto.Clone(base).(ssz.Unmarshaler)
	if !ok {
		return nil, errors.Errorf("message of %T does not support marshaller interface", base)
	}

	// Handle different message types across forks.
	if topic == p2p.BlockSubnetTopicFormat {
		m, err = extractBlockDataType(fDigest[:], s.cfg.clock)
		if err != nil {
			return nil, err
		}
	}
	if topic == p2p.AttestationSubnetTopicFormat {
		m, err = extractAttestationType(fDigest[:], s.cfg.clock)
		if err != nil {
			return nil, err
		}
	}
	if topic == p2p.AggregateAndProofSubnetTopicFormat {
		m, err = extractAggregateAndProofType(fDigest[:], s.cfg.clock)
		if err != nil {
			return nil, err
		}
	}

	if err := s.cfg.p2p.Encoding().DecodeGossip(msg.Data, m); err != nil {
		return nil, err
	}
	return m, nil
}

func extractAttestationType(digest []byte, tor blockchain.TemporalOracle) (ethpb.Att, error) {
	if len(digest) == 0 {
		aFunc, ok := types.AttestationMap[bytesutil.ToBytes4(params.BeaconConfig().GenesisForkVersion)]
		if !ok {
			return nil, errors.New("no attestation type exists for the genesis fork version")
		}
		return aFunc()
	}
	if len(digest) != forkDigestLength {
		return nil, errors.Errorf("invalid digest returned, wanted a length of %d but received %d", forkDigestLength, len(digest))
	}
	vRoot := tor.GenesisValidatorsRoot()
	for k, aFunc := range types.AttestationMap {
		rDigest, err := signing.ComputeForkDigest(k[:], vRoot[:])
		if err != nil {
			return nil, err
		}
		if rDigest == bytesutil.ToBytes4(digest) {
			return aFunc()
		}
	}
	return nil, errors.Wrapf(
		ErrNoValidDigest,
		"could not extract attestation data type, saw digest=%#x, genesis=%v, vr=%#x",
		digest,
		tor.GenesisTime(),
		tor.GenesisValidatorsRoot(),
	)
}

func extractAggregateAndProofType(digest []byte, tor blockchain.TemporalOracle) (ethpb.SignedAggregateAttAndProof, error) {
	if len(digest) == 0 {
		aFunc, ok := types.AggregateAttestationMap[bytesutil.ToBytes4(params.BeaconConfig().GenesisForkVersion)]
		if !ok {
			return nil, errors.New("no aggregate attestation type exists for the genesis fork version")
		}
		return aFunc()
	}
	if len(digest) != forkDigestLength {
		return nil, errors.Errorf("invalid digest returned, wanted a length of %d but received %d", forkDigestLength, len(digest))
	}
	vRoot := tor.GenesisValidatorsRoot()
	for k, aFunc := range types.AggregateAttestationMap {
		rDigest, err := signing.ComputeForkDigest(k[:], vRoot[:])
		if err != nil {
			return nil, err
		}
		if rDigest == bytesutil.ToBytes4(digest) {
			return aFunc()
		}
	}
	return nil, errors.Wrapf(
		ErrNoValidDigest,
		"could not extract aggregate attestation data type, saw digest=%#x, genesis=%v, vr=%#x",
		digest,
		tor.GenesisTime(),
		tor.GenesisValidatorsRoot(),
	)
}

// Replaces our fork digest with the formatter.
func (_ *Service) replaceForkDigest(topic string) (string, error) {
	subStrings := strings.Split(topic, "/")
	if len(subStrings) != 4 {
		return "", errInvalidTopic
	}
	subStrings[2] = "%x"
	return strings.Join(subStrings, "/"), nil
}
