package epoch_processing

import (
	"testing"

	"github.com/prysmaticlabs/prysm/v5/testing/spectest/shared/deneb/epoch_processing"
)

func TestMainnet_Deneb_EpochProcessing_ParticipationFlag(t *testing.T) {
	epoch_processing.RunParticipationFlagUpdatesTests(t, "mainnet")
}
