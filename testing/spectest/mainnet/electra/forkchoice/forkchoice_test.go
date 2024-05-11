package forkchoice

import (
	"testing"

	"github.com/prysmaticlabs/prysm/v5/runtime/version"
	"github.com/prysmaticlabs/prysm/v5/testing/spectest/shared/common/forkchoice"
)

func TestMainnet_Electra_Forkchoice(t *testing.T) {
	forkchoice.Run(t, "mainnet", version.Electra)
}
