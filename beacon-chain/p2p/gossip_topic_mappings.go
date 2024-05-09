package p2p

import (
	"reflect"

	"github.com/prysmaticlabs/prysm/v5/config/params"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/primitives"
	ethpb "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
	"google.golang.org/protobuf/proto"
)

// gossipTopicMappings represent the protocol ID to protobuf message type map for easy
// lookup.
var gossipTopicMappings = map[string]proto.Message{
	BlockSubnetTopicFormat:                    &ethpb.SignedBeaconBlock{},
	AttestationSubnetTopicFormat:              &ethpb.Attestation{},
	ExitSubnetTopicFormat:                     &ethpb.SignedVoluntaryExit{},
	ProposerSlashingSubnetTopicFormat:         &ethpb.ProposerSlashing{},
	AttesterSlashingSubnetTopicFormat:         &ethpb.AttesterSlashing{},
	AggregateAndProofSubnetTopicFormat:        &ethpb.SignedAggregateAttestationAndProof{},
	SyncContributionAndProofSubnetTopicFormat: &ethpb.SignedContributionAndProof{},
	SyncCommitteeSubnetTopicFormat:            &ethpb.SyncCommitteeMessage{},
	BlsToExecutionChangeSubnetTopicFormat:     &ethpb.SignedBLSToExecutionChange{},
	BlobSubnetTopicFormat:                     &ethpb.BlobSidecar{},
}

// GossipTopicMappings is a function to return the assigned data type
// versioned by epoch.
func GossipTopicMappings(topic string, epoch primitives.Epoch) proto.Message {
	// TODO: not pretty, would be nicer if the map itself supported multiple types
	switch topic {
	case BlockSubnetTopicFormat:
		if epoch >= params.BeaconConfig().ElectraForkEpoch {
			return &ethpb.SignedBeaconBlockElectra{}
		}
		if epoch >= params.BeaconConfig().DenebForkEpoch {
			return &ethpb.SignedBeaconBlockDeneb{}
		}
		if epoch >= params.BeaconConfig().CapellaForkEpoch {
			return &ethpb.SignedBeaconBlockCapella{}
		}
		if epoch >= params.BeaconConfig().BellatrixForkEpoch {
			return &ethpb.SignedBeaconBlockBellatrix{}
		}
		if epoch >= params.BeaconConfig().AltairForkEpoch {
			return &ethpb.SignedBeaconBlockAltair{}
		}
		return gossipTopicMappings[topic]
	case AttestationSubnetTopicFormat:
		if epoch >= params.BeaconConfig().ElectraForkEpoch {
			return &ethpb.AttestationElectra{}
		}
		return gossipTopicMappings[topic]
	case AttesterSlashingSubnetTopicFormat:
		if epoch >= params.BeaconConfig().ElectraForkEpoch {
			return &ethpb.AttesterSlashingElectra{}
		}
		return gossipTopicMappings[topic]
	case AggregateAndProofSubnetTopicFormat:
		if epoch >= params.BeaconConfig().ElectraForkEpoch {
			return &ethpb.SignedAggregateAttestationAndProofElectra{}
		}
		return gossipTopicMappings[topic]
	default:
		return gossipTopicMappings[topic]
	}
}

// AllTopics returns all topics stored in our
// gossip mapping.
func AllTopics() []string {
	var topics []string
	for k := range gossipTopicMappings {
		topics = append(topics, k)
	}
	return topics
}

// GossipTypeMapping is the inverse of GossipTopicMappings so that an arbitrary protobuf message
// can be mapped to a protocol ID string.
var GossipTypeMapping = make(map[reflect.Type]string, len(gossipTopicMappings))

func init() {
	for k, v := range gossipTopicMappings {
		GossipTypeMapping[reflect.TypeOf(v)] = k
	}
	// Specially handle Altair objects.
	GossipTypeMapping[reflect.TypeOf(&ethpb.SignedBeaconBlockAltair{})] = BlockSubnetTopicFormat
	// Specially handle Bellatrix objects.
	GossipTypeMapping[reflect.TypeOf(&ethpb.SignedBeaconBlockBellatrix{})] = BlockSubnetTopicFormat
	// Specially handle Capella objects.
	GossipTypeMapping[reflect.TypeOf(&ethpb.SignedBeaconBlockCapella{})] = BlockSubnetTopicFormat
	// Specially handle Deneb objects.
	GossipTypeMapping[reflect.TypeOf(&ethpb.SignedBeaconBlockDeneb{})] = BlockSubnetTopicFormat
	// Specifically handle Electra objects.
	GossipTypeMapping[reflect.TypeOf(&ethpb.SignedBeaconBlockElectra{})] = BlockSubnetTopicFormat
	GossipTypeMapping[reflect.TypeOf(&ethpb.AttestationElectra{})] = AttestationSubnetTopicFormat
	GossipTypeMapping[reflect.TypeOf(&ethpb.SignedAggregateAttestationAndProofElectra{})] = AggregateAndProofSubnetTopicFormat
	GossipTypeMapping[reflect.TypeOf(&ethpb.AttesterSlashingElectra{})] = AttesterSlashingSubnetTopicFormat
}
