package sync

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	pubsubpb "github.com/libp2p/go-libp2p-pubsub/pb"
	gcache "github.com/patrickmn/go-cache"
	types "github.com/prysmaticlabs/eth2-types"
	"github.com/prysmaticlabs/prysm/async/abool"
	mock "github.com/prysmaticlabs/prysm/beacon-chain/blockchain/testing"
	"github.com/prysmaticlabs/prysm/beacon-chain/core"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/transition"
	dbtest "github.com/prysmaticlabs/prysm/beacon-chain/db/testing"
	"github.com/prysmaticlabs/prysm/beacon-chain/operations/attestations"
	"github.com/prysmaticlabs/prysm/beacon-chain/p2p"
	p2ptest "github.com/prysmaticlabs/prysm/beacon-chain/p2p/testing"
	"github.com/prysmaticlabs/prysm/beacon-chain/state/stategen"
	mockSync "github.com/prysmaticlabs/prysm/beacon-chain/sync/initial-sync/testing"
	lruwrpr "github.com/prysmaticlabs/prysm/cache/lru"
	"github.com/prysmaticlabs/prysm/config/params"
	"github.com/prysmaticlabs/prysm/crypto/bls"
	"github.com/prysmaticlabs/prysm/encoding/bytesutil"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1/wrapper"
	"github.com/prysmaticlabs/prysm/testing/assert"
	"github.com/prysmaticlabs/prysm/testing/require"
	"github.com/prysmaticlabs/prysm/testing/util"
	logTest "github.com/sirupsen/logrus/hooks/test"
)

// General note for writing validation tests: Use a random value for any field
// on the beacon block to avoid hitting shared global cache conditions across
// tests in this package.

func TestValidateBeaconBlockPubSub_InvalidSignature(t *testing.T) {
	db := dbtest.SetupDB(t)
	p := p2ptest.NewTestP2P(t)
	ctx := context.Background()
	beaconState, _ := util.DeterministicGenesisState(t, 100)
	parentBlock := util.NewBeaconBlock()
	require.NoError(t, db.SaveBlock(ctx, wrapper.WrappedPhase0SignedBeaconBlock(parentBlock)))
	bRoot, err := parentBlock.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, db.SaveState(ctx, beaconState, bRoot))
	require.NoError(t, db.SaveStateSummary(ctx, &ethpb.StateSummary{Root: bRoot[:]}))
	copied := beaconState.Copy()
	require.NoError(t, copied.SetSlot(1))
	proposerIdx, err := helpers.BeaconProposerIndex(copied)
	require.NoError(t, err)
	msg := util.NewBeaconBlock()
	msg.Block.ParentRoot = bRoot[:]
	msg.Block.Slot = 1
	msg.Block.ProposerIndex = proposerIdx
	msg.Signature = bytesutil.PadTo([]byte("fake"), 96)

	stateGen := stategen.New(db)
	chainService := &mock.ChainService{Genesis: time.Unix(time.Now().Unix()-int64(params.BeaconConfig().SecondsPerSlot), 0),
		FinalizedCheckPoint: &ethpb.Checkpoint{
			Epoch: 0,
			Root:  make([]byte, 32),
		}}
	r := &Service{
		cfg: &Config{
			DB:            db,
			P2P:           p,
			InitialSync:   &mockSync.Sync{IsSyncing: false},
			Chain:         chainService,
			BlockNotifier: chainService.BlockNotifier(),
			StateGen:      stateGen,
		},
		seenBlockCache: lruwrpr.New(10),
		badBlockCache:  lruwrpr.New(10),
	}

	buf := new(bytes.Buffer)
	_, err = p.Encoding().EncodeGossip(buf, msg)
	require.NoError(t, err)
	topic := p2p.GossipTypeMapping[reflect.TypeOf(msg)]
	m := &pubsub.Message{
		Message: &pubsubpb.Message{
			Data:  buf.Bytes(),
			Topic: &topic,
		},
	}
	res, err := r.validateBeaconBlockPubSub(ctx, "", m)
	assert.NotNil(t, err)
	result := res == pubsub.ValidationReject
	assert.Equal(t, true, result)
}

func TestValidateBeaconBlockPubSub_BlockAlreadyPresentInDB(t *testing.T) {
	db := dbtest.SetupDB(t)
	ctx := context.Background()

	p := p2ptest.NewTestP2P(t)
	msg := util.NewBeaconBlock()
	msg.Block.Slot = 100
	msg.Block.ParentRoot = util.Random32Bytes(t)
	require.NoError(t, db.SaveBlock(context.Background(), wrapper.WrappedPhase0SignedBeaconBlock(msg)))

	chainService := &mock.ChainService{Genesis: time.Now()}
	r := &Service{
		cfg: &Config{
			DB:            db,
			P2P:           p,
			InitialSync:   &mockSync.Sync{IsSyncing: false},
			Chain:         chainService,
			BlockNotifier: chainService.BlockNotifier(),
		},
		seenBlockCache: lruwrpr.New(10),
		badBlockCache:  lruwrpr.New(10),
	}

	buf := new(bytes.Buffer)
	_, err := p.Encoding().EncodeGossip(buf, msg)
	require.NoError(t, err)

	topic := p2p.GossipTypeMapping[reflect.TypeOf(msg)]
	m := &pubsub.Message{
		Message: &pubsubpb.Message{
			Data:  buf.Bytes(),
			Topic: &topic,
		},
	}
	res, err := r.validateBeaconBlockPubSub(ctx, "", m)
	_ = err
	result := res == pubsub.ValidationAccept
	assert.Equal(t, false, result)
}

func TestValidateBeaconBlockPubSub_CanRecoverStateSummary(t *testing.T) {
	db := dbtest.SetupDB(t)
	p := p2ptest.NewTestP2P(t)
	ctx := context.Background()
	beaconState, privKeys := util.DeterministicGenesisState(t, 100)
	parentBlock := util.NewBeaconBlock()
	require.NoError(t, db.SaveBlock(ctx, wrapper.WrappedPhase0SignedBeaconBlock(parentBlock)))
	bRoot, err := parentBlock.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, db.SaveState(ctx, beaconState, bRoot))
	copied := beaconState.Copy()
	require.NoError(t, copied.SetSlot(1))
	proposerIdx, err := helpers.BeaconProposerIndex(copied)
	require.NoError(t, err)
	msg := util.NewBeaconBlock()
	msg.Block.ParentRoot = bRoot[:]
	msg.Block.Slot = 1
	msg.Block.ProposerIndex = proposerIdx
	msg.Signature, err = helpers.ComputeDomainAndSign(beaconState, 0, msg.Block, params.BeaconConfig().DomainBeaconProposer, privKeys[proposerIdx])
	require.NoError(t, err)

	stateGen := stategen.New(db)
	chainService := &mock.ChainService{Genesis: time.Unix(time.Now().Unix()-int64(params.BeaconConfig().SecondsPerSlot), 0),
		State: beaconState,
		FinalizedCheckPoint: &ethpb.Checkpoint{
			Epoch: 0,
			Root:  make([]byte, 32),
		},
	}
	r := &Service{
		cfg: &Config{
			DB:            db,
			P2P:           p,
			InitialSync:   &mockSync.Sync{IsSyncing: false},
			Chain:         chainService,
			BlockNotifier: chainService.BlockNotifier(),
			StateGen:      stateGen,
		},
		seenBlockCache:      lruwrpr.New(10),
		badBlockCache:       lruwrpr.New(10),
		slotToPendingBlocks: gcache.New(time.Second, 2*time.Second),
		seenPendingBlocks:   make(map[[32]byte]bool),
	}
	buf := new(bytes.Buffer)
	_, err = p.Encoding().EncodeGossip(buf, msg)
	require.NoError(t, err)
	topic := p2p.GossipTypeMapping[reflect.TypeOf(msg)]
	digest, err := r.currentForkDigest()
	assert.NoError(t, err)
	topic = r.addDigestToTopic(topic, digest)
	m := &pubsub.Message{
		Message: &pubsubpb.Message{
			Data:  buf.Bytes(),
			Topic: &topic,
		},
	}
	res, err := r.validateBeaconBlockPubSub(ctx, "", m)
	assert.NoError(t, err)
	result := res == pubsub.ValidationAccept
	assert.Equal(t, true, result)
	assert.NotNil(t, m.ValidatorData, "Decoded message was not set on the message validator data")
}

func TestValidateBeaconBlockPubSub_IsInCache(t *testing.T) {
	db := dbtest.SetupDB(t)
	p := p2ptest.NewTestP2P(t)
	ctx := context.Background()
	beaconState, privKeys := util.DeterministicGenesisState(t, 100)
	parentBlock := util.NewBeaconBlock()
	bRoot, err := parentBlock.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, db.SaveState(ctx, beaconState, bRoot))
	require.NoError(t, db.SaveStateSummary(ctx, &ethpb.StateSummary{Root: bRoot[:]}))
	copied := beaconState.Copy()
	require.NoError(t, copied.SetSlot(1))
	proposerIdx, err := helpers.BeaconProposerIndex(copied)
	require.NoError(t, err)
	msg := util.NewBeaconBlock()
	msg.Block.ParentRoot = bRoot[:]
	msg.Block.Slot = 1
	msg.Block.ProposerIndex = proposerIdx
	msg.Signature, err = helpers.ComputeDomainAndSign(beaconState, 0, msg.Block, params.BeaconConfig().DomainBeaconProposer, privKeys[proposerIdx])
	require.NoError(t, err)

	stateGen := stategen.New(db)
	chainService := &mock.ChainService{Genesis: time.Unix(time.Now().Unix()-int64(params.BeaconConfig().SecondsPerSlot), 0),
		State: beaconState,
		FinalizedCheckPoint: &ethpb.Checkpoint{
			Epoch: 0,
			Root:  make([]byte, 32),
		},
		InitSyncBlockRoots: map[[32]byte]bool{bRoot: true},
	}
	r := &Service{
		cfg: &Config{
			DB:            db,
			P2P:           p,
			InitialSync:   &mockSync.Sync{IsSyncing: false},
			Chain:         chainService,
			BlockNotifier: chainService.BlockNotifier(),
			StateGen:      stateGen,
		},
		seenBlockCache:      lruwrpr.New(10),
		badBlockCache:       lruwrpr.New(10),
		slotToPendingBlocks: gcache.New(time.Second, 2*time.Second),
		seenPendingBlocks:   make(map[[32]byte]bool),
	}
	buf := new(bytes.Buffer)
	_, err = p.Encoding().EncodeGossip(buf, msg)
	require.NoError(t, err)
	topic := p2p.GossipTypeMapping[reflect.TypeOf(msg)]
	digest, err := r.currentForkDigest()
	assert.NoError(t, err)
	topic = r.addDigestToTopic(topic, digest)
	m := &pubsub.Message{
		Message: &pubsubpb.Message{
			Data:  buf.Bytes(),
			Topic: &topic,
		},
	}
	res, err := r.validateBeaconBlockPubSub(ctx, "", m)
	assert.NoError(t, err)
	result := res == pubsub.ValidationAccept
	assert.Equal(t, true, result)
	assert.NotNil(t, m.ValidatorData, "Decoded message was not set on the message validator data")
}

func TestValidateBeaconBlockPubSub_ValidProposerSignature(t *testing.T) {
	db := dbtest.SetupDB(t)
	p := p2ptest.NewTestP2P(t)
	ctx := context.Background()
	beaconState, privKeys := util.DeterministicGenesisState(t, 100)
	parentBlock := util.NewBeaconBlock()
	require.NoError(t, db.SaveBlock(ctx, wrapper.WrappedPhase0SignedBeaconBlock(parentBlock)))
	bRoot, err := parentBlock.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, db.SaveState(ctx, beaconState, bRoot))
	require.NoError(t, db.SaveStateSummary(ctx, &ethpb.StateSummary{Root: bRoot[:]}))
	copied := beaconState.Copy()
	require.NoError(t, copied.SetSlot(1))
	proposerIdx, err := helpers.BeaconProposerIndex(copied)
	require.NoError(t, err)
	msg := util.NewBeaconBlock()
	msg.Block.ParentRoot = bRoot[:]
	msg.Block.Slot = 1
	msg.Block.ProposerIndex = proposerIdx
	msg.Signature, err = helpers.ComputeDomainAndSign(beaconState, 0, msg.Block, params.BeaconConfig().DomainBeaconProposer, privKeys[proposerIdx])
	require.NoError(t, err)

	stateGen := stategen.New(db)
	chainService := &mock.ChainService{Genesis: time.Unix(time.Now().Unix()-int64(params.BeaconConfig().SecondsPerSlot), 0),
		State: beaconState,
		FinalizedCheckPoint: &ethpb.Checkpoint{
			Epoch: 0,
			Root:  make([]byte, 32),
		},
	}
	r := &Service{
		cfg: &Config{
			DB:            db,
			P2P:           p,
			InitialSync:   &mockSync.Sync{IsSyncing: false},
			Chain:         chainService,
			BlockNotifier: chainService.BlockNotifier(),
			StateGen:      stateGen,
		},
		seenBlockCache:      lruwrpr.New(10),
		badBlockCache:       lruwrpr.New(10),
		slotToPendingBlocks: gcache.New(time.Second, 2*time.Second),
		seenPendingBlocks:   make(map[[32]byte]bool),
	}
	buf := new(bytes.Buffer)
	_, err = p.Encoding().EncodeGossip(buf, msg)
	require.NoError(t, err)
	topic := p2p.GossipTypeMapping[reflect.TypeOf(msg)]
	digest, err := r.currentForkDigest()
	assert.NoError(t, err)
	topic = r.addDigestToTopic(topic, digest)
	m := &pubsub.Message{
		Message: &pubsubpb.Message{
			Data:  buf.Bytes(),
			Topic: &topic,
		},
	}
	res, err := r.validateBeaconBlockPubSub(ctx, "", m)
	assert.NoError(t, err)
	result := res == pubsub.ValidationAccept
	assert.Equal(t, true, result)
	assert.NotNil(t, m.ValidatorData, "Decoded message was not set on the message validator data")
}

func TestValidateBeaconBlockPubSub_WithLookahead(t *testing.T) {
	db := dbtest.SetupDB(t)
	p := p2ptest.NewTestP2P(t)
	ctx := context.Background()
	beaconState, privKeys := util.DeterministicGenesisState(t, 100)
	parentBlock := util.NewBeaconBlock()
	require.NoError(t, db.SaveBlock(ctx, wrapper.WrappedPhase0SignedBeaconBlock(parentBlock)))
	bRoot, err := parentBlock.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, db.SaveState(ctx, beaconState, bRoot))
	require.NoError(t, db.SaveStateSummary(ctx, &ethpb.StateSummary{Root: bRoot[:]}))
	copied := beaconState.Copy()
	// The next block is only 1 epoch ahead so as to not induce a new seed.
	blkSlot := params.BeaconConfig().SlotsPerEpoch.Mul(uint64(core.NextEpoch(copied)))
	copied, err = transition.ProcessSlots(context.Background(), copied, blkSlot)
	require.NoError(t, err)
	proposerIdx, err := helpers.BeaconProposerIndex(copied)
	require.NoError(t, err)
	msg := util.NewBeaconBlock()
	msg.Block.ProposerIndex = proposerIdx
	msg.Block.Slot = blkSlot
	msg.Block.ParentRoot = bRoot[:]
	msg.Signature, err = helpers.ComputeDomainAndSign(beaconState, 0, msg.Block, params.BeaconConfig().DomainBeaconProposer, privKeys[proposerIdx])
	require.NoError(t, err)

	stateGen := stategen.New(db)
	offset := int64(blkSlot.Mul(params.BeaconConfig().SecondsPerSlot))
	chainService := &mock.ChainService{Genesis: time.Unix(time.Now().Unix()-offset, 0),
		State: beaconState,
		FinalizedCheckPoint: &ethpb.Checkpoint{
			Epoch: 0,
		}}
	r := &Service{
		cfg: &Config{
			DB:            db,
			P2P:           p,
			InitialSync:   &mockSync.Sync{IsSyncing: false},
			Chain:         chainService,
			BlockNotifier: chainService.BlockNotifier(),
			StateGen:      stateGen,
		},
		seenBlockCache:      lruwrpr.New(10),
		badBlockCache:       lruwrpr.New(10),
		slotToPendingBlocks: gcache.New(time.Second, 2*time.Second),
		seenPendingBlocks:   make(map[[32]byte]bool),
		subHandler:          newSubTopicHandler(),
	}
	buf := new(bytes.Buffer)
	_, err = p.Encoding().EncodeGossip(buf, msg)
	require.NoError(t, err)
	topic := p2p.GossipTypeMapping[reflect.TypeOf(msg)]
	digest, err := r.currentForkDigest()
	assert.NoError(t, err)
	topic = r.addDigestToTopic(topic, digest)
	m := &pubsub.Message{
		Message: &pubsubpb.Message{
			Data:  buf.Bytes(),
			Topic: &topic,
		},
	}
	res, err := r.validateBeaconBlockPubSub(ctx, "", m)
	assert.NoError(t, err)
	result := res == pubsub.ValidationAccept
	assert.Equal(t, true, result)
	assert.NotNil(t, m.ValidatorData, "Decoded message was not set on the message validator data")
}

func TestValidateBeaconBlockPubSub_AdvanceEpochsForState(t *testing.T) {
	db := dbtest.SetupDB(t)
	p := p2ptest.NewTestP2P(t)
	ctx := context.Background()
	beaconState, privKeys := util.DeterministicGenesisState(t, 100)
	parentBlock := util.NewBeaconBlock()
	require.NoError(t, db.SaveBlock(ctx, wrapper.WrappedPhase0SignedBeaconBlock(parentBlock)))
	bRoot, err := parentBlock.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, db.SaveState(ctx, beaconState, bRoot))
	require.NoError(t, db.SaveStateSummary(ctx, &ethpb.StateSummary{Root: bRoot[:]}))
	copied := beaconState.Copy()
	// The next block is at least 2 epochs ahead to induce shuffling and a new seed.
	blkSlot := params.BeaconConfig().SlotsPerEpoch * 2
	copied, err = transition.ProcessSlots(context.Background(), copied, blkSlot)
	require.NoError(t, err)
	proposerIdx, err := helpers.BeaconProposerIndex(copied)
	require.NoError(t, err)
	msg := util.NewBeaconBlock()
	msg.Block.ProposerIndex = proposerIdx
	msg.Block.Slot = blkSlot
	msg.Block.ParentRoot = bRoot[:]
	msg.Signature, err = helpers.ComputeDomainAndSign(beaconState, 0, msg.Block, params.BeaconConfig().DomainBeaconProposer, privKeys[proposerIdx])
	require.NoError(t, err)

	stateGen := stategen.New(db)
	offset := int64(blkSlot.Mul(params.BeaconConfig().SecondsPerSlot))
	chainService := &mock.ChainService{Genesis: time.Unix(time.Now().Unix()-offset, 0),
		State: beaconState,
		FinalizedCheckPoint: &ethpb.Checkpoint{
			Epoch: 0,
		}}
	r := &Service{
		cfg: &Config{
			DB:            db,
			P2P:           p,
			InitialSync:   &mockSync.Sync{IsSyncing: false},
			Chain:         chainService,
			BlockNotifier: chainService.BlockNotifier(),
			StateGen:      stateGen,
		},
		seenBlockCache:      lruwrpr.New(10),
		badBlockCache:       lruwrpr.New(10),
		slotToPendingBlocks: gcache.New(time.Second, 2*time.Second),
		seenPendingBlocks:   make(map[[32]byte]bool),
	}
	buf := new(bytes.Buffer)
	_, err = p.Encoding().EncodeGossip(buf, msg)
	require.NoError(t, err)
	topic := p2p.GossipTypeMapping[reflect.TypeOf(msg)]
	digest, err := r.currentForkDigest()
	assert.NoError(t, err)
	topic = r.addDigestToTopic(topic, digest)
	m := &pubsub.Message{
		Message: &pubsubpb.Message{
			Data:  buf.Bytes(),
			Topic: &topic,
		},
	}
	res, err := r.validateBeaconBlockPubSub(ctx, "", m)
	assert.NoError(t, err)
	result := res == pubsub.ValidationAccept
	assert.Equal(t, true, result)
	assert.NotNil(t, m.ValidatorData, "Decoded message was not set on the message validator data")
}

func TestValidateBeaconBlockPubSub_Syncing(t *testing.T) {
	db := dbtest.SetupDB(t)
	p := p2ptest.NewTestP2P(t)
	ctx := context.Background()
	b := []byte("sk")
	b32 := bytesutil.ToBytes32(b)
	sk, err := bls.SecretKeyFromBytes(b32[:])
	require.NoError(t, err)
	msg := util.NewBeaconBlock()
	msg.Block.ParentRoot = util.Random32Bytes(t)
	msg.Signature = sk.Sign([]byte("data")).Marshal()
	chainService := &mock.ChainService{
		Genesis: time.Now(),
		FinalizedCheckPoint: &ethpb.Checkpoint{
			Epoch: 0,
		}}
	r := &Service{
		cfg: &Config{
			DB:            db,
			P2P:           p,
			InitialSync:   &mockSync.Sync{IsSyncing: true},
			Chain:         chainService,
			BlockNotifier: chainService.BlockNotifier(),
		},
	}

	buf := new(bytes.Buffer)
	_, err = p.Encoding().EncodeGossip(buf, msg)
	require.NoError(t, err)
	topic := p2p.GossipTypeMapping[reflect.TypeOf(msg)]
	m := &pubsub.Message{
		Message: &pubsubpb.Message{
			Data:  buf.Bytes(),
			Topic: &topic,
		},
	}
	res, err := r.validateBeaconBlockPubSub(ctx, "", m)
	_ = err
	result := res == pubsub.ValidationAccept
	assert.Equal(t, false, result)
}

func TestValidateBeaconBlockPubSub_AcceptBlocksFromNearFuture(t *testing.T) {
	db := dbtest.SetupDB(t)
	p := p2ptest.NewTestP2P(t)
	ctx := context.Background()

	beaconState, privKeys := util.DeterministicGenesisState(t, 100)
	parentBlock := util.NewBeaconBlock()
	require.NoError(t, db.SaveBlock(ctx, wrapper.WrappedPhase0SignedBeaconBlock(parentBlock)))
	bRoot, err := parentBlock.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, db.SaveState(ctx, beaconState, bRoot))
	require.NoError(t, db.SaveStateSummary(ctx, &ethpb.StateSummary{Root: bRoot[:]}))
	copied := beaconState.Copy()
	require.NoError(t, copied.SetSlot(1))
	proposerIdx, err := helpers.BeaconProposerIndex(copied)
	require.NoError(t, err)

	msg := util.NewBeaconBlock()
	msg.Block.Slot = 2 // two slots in future
	msg.Block.ParentRoot = bRoot[:]
	msg.Block.ProposerIndex = proposerIdx
	msg.Signature, err = helpers.ComputeDomainAndSign(beaconState, 0, msg.Block, params.BeaconConfig().DomainBeaconProposer, privKeys[proposerIdx])
	require.NoError(t, err)

	stateGen := stategen.New(db)
	chainService := &mock.ChainService{Genesis: time.Now(),
		FinalizedCheckPoint: &ethpb.Checkpoint{
			Epoch: 0,
			Root:  make([]byte, 32),
		}}
	r := &Service{
		cfg: &Config{
			P2P:           p,
			DB:            db,
			InitialSync:   &mockSync.Sync{IsSyncing: false},
			Chain:         chainService,
			BlockNotifier: chainService.BlockNotifier(),
			StateGen:      stateGen,
		},
		chainStarted:        abool.New(),
		seenBlockCache:      lruwrpr.New(10),
		badBlockCache:       lruwrpr.New(10),
		slotToPendingBlocks: gcache.New(time.Second, 2*time.Second),
		seenPendingBlocks:   make(map[[32]byte]bool),
	}

	buf := new(bytes.Buffer)
	_, err = p.Encoding().EncodeGossip(buf, msg)
	require.NoError(t, err)
	topic := p2p.GossipTypeMapping[reflect.TypeOf(msg)]
	digest, err := r.currentForkDigest()
	assert.NoError(t, err)
	topic = r.addDigestToTopic(topic, digest)
	m := &pubsub.Message{
		Message: &pubsubpb.Message{
			Data:  buf.Bytes(),
			Topic: &topic,
		},
	}
	res, err := r.validateBeaconBlockPubSub(ctx, "", m)
	_ = err
	result := res == pubsub.ValidationIgnore
	assert.Equal(t, true, result)

	// check if the block is inserted in the Queue
	assert.Equal(t, true, len(r.pendingBlocksInCache(msg.Block.Slot)) == 1)
}

func TestValidateBeaconBlockPubSub_RejectBlocksFromFuture(t *testing.T) {
	db := dbtest.SetupDB(t)
	p := p2ptest.NewTestP2P(t)
	ctx := context.Background()
	b := []byte("sk")
	b32 := bytesutil.ToBytes32(b)
	sk, err := bls.SecretKeyFromBytes(b32[:])
	require.NoError(t, err)
	msg := util.NewBeaconBlock()
	msg.Block.Slot = 3
	msg.Block.ParentRoot = util.Random32Bytes(t)
	msg.Signature = sk.Sign([]byte("data")).Marshal()

	chainService := &mock.ChainService{Genesis: time.Now()}
	r := &Service{
		cfg: &Config{
			P2P:           p,
			DB:            db,
			InitialSync:   &mockSync.Sync{IsSyncing: false},
			Chain:         chainService,
			BlockNotifier: chainService.BlockNotifier(),
		},
		chainStarted:        abool.New(),
		seenBlockCache:      lruwrpr.New(10),
		badBlockCache:       lruwrpr.New(10),
		slotToPendingBlocks: gcache.New(time.Second, 2*time.Second),
		seenPendingBlocks:   make(map[[32]byte]bool),
	}

	buf := new(bytes.Buffer)
	_, err = p.Encoding().EncodeGossip(buf, msg)
	require.NoError(t, err)
	topic := p2p.GossipTypeMapping[reflect.TypeOf(msg)]
	m := &pubsub.Message{
		Message: &pubsubpb.Message{
			Data:  buf.Bytes(),
			Topic: &topic,
		},
	}
	res, err := r.validateBeaconBlockPubSub(ctx, "", m)
	_ = err
	result := res == pubsub.ValidationAccept
	assert.Equal(t, false, result)
}

func TestValidateBeaconBlockPubSub_RejectBlocksFromThePast(t *testing.T) {
	db := dbtest.SetupDB(t)
	b := []byte("sk")
	b32 := bytesutil.ToBytes32(b)
	p := p2ptest.NewTestP2P(t)
	ctx := context.Background()
	sk, err := bls.SecretKeyFromBytes(b32[:])
	require.NoError(t, err)
	msg := util.NewBeaconBlock()
	msg.Block.ParentRoot = util.Random32Bytes(t)
	msg.Block.Slot = 10
	msg.Signature = sk.Sign([]byte("data")).Marshal()

	genesisTime := time.Now()
	chainService := &mock.ChainService{
		Genesis: time.Unix(genesisTime.Unix()-1000, 0),
		FinalizedCheckPoint: &ethpb.Checkpoint{
			Epoch: 1,
		},
	}
	r := &Service{
		cfg: &Config{
			DB:            db,
			P2P:           p,
			InitialSync:   &mockSync.Sync{IsSyncing: false},
			Chain:         chainService,
			BlockNotifier: chainService.BlockNotifier(),
		},
		seenBlockCache: lruwrpr.New(10),
		badBlockCache:  lruwrpr.New(10),
	}

	buf := new(bytes.Buffer)
	_, err = p.Encoding().EncodeGossip(buf, msg)
	require.NoError(t, err)
	topic := p2p.GossipTypeMapping[reflect.TypeOf(msg)]
	m := &pubsub.Message{
		Message: &pubsubpb.Message{
			Data:  buf.Bytes(),
			Topic: &topic,
		},
	}
	res, err := r.validateBeaconBlockPubSub(ctx, "", m)
	_ = err
	result := res == pubsub.ValidationAccept
	assert.Equal(t, false, result)
}

func TestValidateBeaconBlockPubSub_SeenProposerSlot(t *testing.T) {
	db := dbtest.SetupDB(t)
	p := p2ptest.NewTestP2P(t)
	ctx := context.Background()
	beaconState, privKeys := util.DeterministicGenesisState(t, 100)
	parentBlock := util.NewBeaconBlock()
	require.NoError(t, db.SaveBlock(ctx, wrapper.WrappedPhase0SignedBeaconBlock(parentBlock)))
	bRoot, err := parentBlock.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, db.SaveState(ctx, beaconState, bRoot))
	proposerIdx, err := helpers.BeaconProposerIndex(beaconState)
	require.NoError(t, err)

	msg := util.NewBeaconBlock()
	msg.Block.Slot = 1
	msg.Block.ProposerIndex = proposerIdx
	msg.Block.ParentRoot = bRoot[:]
	msg.Signature, err = helpers.ComputeDomainAndSign(beaconState, 0, msg.Block, params.BeaconConfig().DomainBeaconProposer, privKeys[proposerIdx])
	require.NoError(t, err)

	chainService := &mock.ChainService{Genesis: time.Unix(time.Now().Unix()-int64(params.BeaconConfig().SecondsPerSlot), 0),
		State: beaconState,
		FinalizedCheckPoint: &ethpb.Checkpoint{
			Epoch: 0,
			Root:  make([]byte, 32),
		},
	}
	r := &Service{
		cfg: &Config{
			DB:            db,
			P2P:           p,
			InitialSync:   &mockSync.Sync{IsSyncing: false},
			Chain:         chainService,
			BlockNotifier: chainService.BlockNotifier(),
		},
		seenBlockCache:      lruwrpr.New(10),
		badBlockCache:       lruwrpr.New(10),
		slotToPendingBlocks: gcache.New(time.Second, 2*time.Second),
		seenPendingBlocks:   make(map[[32]byte]bool),
	}

	buf := new(bytes.Buffer)
	_, err = p.Encoding().EncodeGossip(buf, msg)
	require.NoError(t, err)
	topic := p2p.GossipTypeMapping[reflect.TypeOf(msg)]
	m := &pubsub.Message{
		Message: &pubsubpb.Message{
			Data:  buf.Bytes(),
			Topic: &topic,
		},
	}
	r.setSeenBlockIndexSlot(msg.Block.Slot, msg.Block.ProposerIndex)
	time.Sleep(10 * time.Millisecond) // Wait for cached value to pass through buffers.
	res, err := r.validateBeaconBlockPubSub(ctx, "", m)
	_ = err
	result := res == pubsub.ValidationAccept
	assert.Equal(t, false, result)
}

func TestValidateBeaconBlockPubSub_FilterByFinalizedEpoch(t *testing.T) {
	hook := logTest.NewGlobal()
	db := dbtest.SetupDB(t)
	p := p2ptest.NewTestP2P(t)

	parent := util.NewBeaconBlock()
	require.NoError(t, db.SaveBlock(context.Background(), wrapper.WrappedPhase0SignedBeaconBlock(parent)))
	parentRoot, err := parent.Block.HashTreeRoot()
	require.NoError(t, err)
	chain := &mock.ChainService{Genesis: time.Unix(time.Now().Unix()-int64(params.BeaconConfig().SecondsPerSlot), 0),
		FinalizedCheckPoint: &ethpb.Checkpoint{
			Epoch: 1,
		},
		ValidatorsRoot: [32]byte{},
	}

	r := &Service{
		cfg: &Config{
			DB:            db,
			P2P:           p,
			Chain:         chain,
			BlockNotifier: chain.BlockNotifier(),
			AttPool:       attestations.NewPool(),
			InitialSync:   &mockSync.Sync{IsSyncing: false},
		},
		seenBlockCache: lruwrpr.New(10),
		badBlockCache:  lruwrpr.New(10),
	}

	b := util.NewBeaconBlock()
	b.Block.Slot = 1
	b.Block.ParentRoot = parentRoot[:]
	buf := new(bytes.Buffer)
	_, err = p.Encoding().EncodeGossip(buf, b)
	require.NoError(t, err)
	digest, err := helpers.ComputeForkDigest(params.BeaconConfig().GenesisForkVersion, make([]byte, 32))
	assert.NoError(t, err)
	topic := fmt.Sprintf(p2p.GossipTypeMapping[reflect.TypeOf(b)], digest)
	m := &pubsub.Message{
		Message: &pubsubpb.Message{
			Data:  buf.Bytes(),
			Topic: &topic,
		},
	}

	res, err := r.validateBeaconBlockPubSub(context.Background(), "", m)
	_ = err
	assert.Equal(t, pubsub.ValidationIgnore, res)

	hook.Reset()
	b.Block.Slot = params.BeaconConfig().SlotsPerEpoch
	buf = new(bytes.Buffer)
	_, err = p.Encoding().EncodeGossip(buf, b)
	require.NoError(t, err)
	m = &pubsub.Message{
		Message: &pubsubpb.Message{
			Data:  buf.Bytes(),
			Topic: &topic,
		},
	}

	res, err = r.validateBeaconBlockPubSub(context.Background(), "", m)
	assert.NoError(t, err)
	assert.Equal(t, pubsub.ValidationIgnore, res)
}

func TestValidateBeaconBlockPubSub_ParentNotFinalizedDescendant(t *testing.T) {
	db := dbtest.SetupDB(t)
	p := p2ptest.NewTestP2P(t)
	ctx := context.Background()
	beaconState, privKeys := util.DeterministicGenesisState(t, 100)
	parentBlock := util.NewBeaconBlock()
	require.NoError(t, db.SaveBlock(ctx, wrapper.WrappedPhase0SignedBeaconBlock(parentBlock)))
	bRoot, err := parentBlock.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, db.SaveState(ctx, beaconState, bRoot))
	require.NoError(t, db.SaveStateSummary(ctx, &ethpb.StateSummary{Root: bRoot[:]}))
	copied := beaconState.Copy()
	require.NoError(t, copied.SetSlot(1))
	proposerIdx, err := helpers.BeaconProposerIndex(copied)
	require.NoError(t, err)
	msg := util.NewBeaconBlock()
	msg.Block.Slot = 1
	msg.Block.ProposerIndex = proposerIdx
	msg.Block.ParentRoot = bRoot[:]
	msg.Signature, err = helpers.ComputeDomainAndSign(beaconState, 0, msg.Block, params.BeaconConfig().DomainBeaconProposer, privKeys[proposerIdx])
	require.NoError(t, err)

	stateGen := stategen.New(db)
	chainService := &mock.ChainService{Genesis: time.Unix(time.Now().Unix()-int64(params.BeaconConfig().SecondsPerSlot), 0),
		State: beaconState,
		FinalizedCheckPoint: &ethpb.Checkpoint{
			Epoch: 0,
			Root:  make([]byte, 32),
		},
		VerifyBlkDescendantErr: errors.New("not part of finalized chain"),
	}
	r := &Service{
		cfg: &Config{
			DB:            db,
			P2P:           p,
			InitialSync:   &mockSync.Sync{IsSyncing: false},
			Chain:         chainService,
			BlockNotifier: chainService.BlockNotifier(),
			StateGen:      stateGen,
		},
		seenBlockCache:      lruwrpr.New(10),
		badBlockCache:       lruwrpr.New(10),
		slotToPendingBlocks: gcache.New(time.Second, 2*time.Second),
		seenPendingBlocks:   make(map[[32]byte]bool),
	}
	buf := new(bytes.Buffer)
	_, err = p.Encoding().EncodeGossip(buf, msg)
	require.NoError(t, err)
	topic := p2p.GossipTypeMapping[reflect.TypeOf(msg)]
	digest, err := r.currentForkDigest()
	assert.NoError(t, err)
	topic = r.addDigestToTopic(topic, digest)
	m := &pubsub.Message{
		Message: &pubsubpb.Message{
			Data:  buf.Bytes(),
			Topic: &topic,
		},
	}
	res, err := r.validateBeaconBlockPubSub(ctx, "", m)
	assert.Equal(t, pubsub.ValidationReject, res, "Wrong validation result returned")
	require.ErrorContains(t, "not part of finalized chain", err)
}

func TestValidateBeaconBlockPubSub_InvalidParentBlock(t *testing.T) {
	db := dbtest.SetupDB(t)
	p := p2ptest.NewTestP2P(t)
	ctx := context.Background()
	beaconState, privKeys := util.DeterministicGenesisState(t, 100)
	parentBlock := util.NewBeaconBlock()
	require.NoError(t, db.SaveBlock(ctx, wrapper.WrappedPhase0SignedBeaconBlock(parentBlock)))
	bRoot, err := parentBlock.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, db.SaveState(ctx, beaconState, bRoot))
	require.NoError(t, db.SaveStateSummary(ctx, &ethpb.StateSummary{Root: bRoot[:]}))
	copied := beaconState.Copy()
	require.NoError(t, copied.SetSlot(1))
	proposerIdx, err := helpers.BeaconProposerIndex(copied)
	require.NoError(t, err)
	msg := util.NewBeaconBlock()
	msg.Block.ProposerIndex = proposerIdx
	msg.Block.Slot = 1
	msg.Block.ParentRoot = bRoot[:]
	msg.Signature, err = helpers.ComputeDomainAndSign(beaconState, 0, msg.Block, params.BeaconConfig().DomainBeaconProposer, privKeys[proposerIdx])
	require.NoError(t, err)

	// Mutate Signature
	copy(msg.Signature[:4], []byte{1, 2, 3, 4})
	currBlockRoot, err := msg.Block.HashTreeRoot()
	require.NoError(t, err)

	stateGen := stategen.New(db)
	chainService := &mock.ChainService{Genesis: time.Unix(time.Now().Unix()-int64(params.BeaconConfig().SecondsPerSlot), 0),
		State: beaconState,
		FinalizedCheckPoint: &ethpb.Checkpoint{
			Epoch: 0,
		}}
	r := &Service{
		cfg: &Config{
			DB:            db,
			P2P:           p,
			InitialSync:   &mockSync.Sync{IsSyncing: false},
			Chain:         chainService,
			BlockNotifier: chainService.BlockNotifier(),
			StateGen:      stateGen,
		},
		seenBlockCache:      lruwrpr.New(10),
		badBlockCache:       lruwrpr.New(10),
		slotToPendingBlocks: gcache.New(time.Second, 2*time.Second),
		seenPendingBlocks:   make(map[[32]byte]bool),
	}
	buf := new(bytes.Buffer)
	_, err = p.Encoding().EncodeGossip(buf, msg)
	require.NoError(t, err)
	topic := p2p.GossipTypeMapping[reflect.TypeOf(msg)]
	m := &pubsub.Message{
		Message: &pubsubpb.Message{
			Data:  buf.Bytes(),
			Topic: &topic,
		},
	}
	res, err := r.validateBeaconBlockPubSub(ctx, "", m)
	_ = err
	result := res == pubsub.ValidationAccept
	assert.Equal(t, false, result)

	require.NoError(t, copied.SetSlot(2))
	proposerIdx, err = helpers.BeaconProposerIndex(copied)
	require.NoError(t, err)

	msg = util.NewBeaconBlock()
	msg.Block.Slot = 2
	msg.Block.ProposerIndex = proposerIdx
	msg.Block.ParentRoot = currBlockRoot[:]

	buf = new(bytes.Buffer)
	_, err = p.Encoding().EncodeGossip(buf, msg)
	require.NoError(t, err)
	m = &pubsub.Message{
		Message: &pubsubpb.Message{
			Data:  buf.Bytes(),
			Topic: &topic,
		},
	}

	res, err = r.validateBeaconBlockPubSub(ctx, "", m)
	_ = err
	result = res == pubsub.ValidationAccept
	// Expect block with bad parent to fail too
	assert.Equal(t, false, result)
}

func TestValidateBeaconBlockPubSub_RejectEvilBlocksFromFuture(t *testing.T) {
	db := dbtest.SetupDB(t)
	p := p2ptest.NewTestP2P(t)
	ctx := context.Background()

	beaconState, privKeys := util.DeterministicGenesisState(t, 100)
	parentBlock := util.NewBeaconBlock()
	require.NoError(t, db.SaveBlock(ctx, wrapper.WrappedPhase0SignedBeaconBlock(parentBlock)))
	bRoot, err := parentBlock.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, db.SaveState(ctx, beaconState, bRoot))
	require.NoError(t, db.SaveStateSummary(ctx, &ethpb.StateSummary{Root: bRoot[:]}))

	copied := beaconState.Copy()
	// The next block is at least 2 epochs ahead to induce shuffling and a new seed.
	blkSlot := params.BeaconConfig().SlotsPerEpoch * 2
	copied, err = transition.ProcessSlots(context.Background(), copied, blkSlot)
	require.NoError(t, err)
	proposerIdx, err := helpers.BeaconProposerIndex(copied)
	require.NoError(t, err)

	msg := util.NewBeaconBlock()
	msg.Block.ProposerIndex = proposerIdx
	msg.Block.Slot = blkSlot

	perSlot := params.BeaconConfig().SecondsPerSlot
	// current slot time
	slotsSinceGenesis := types.Slot(1000)
	// max uint, divided by slot time. But avoid losing precision too much.
	overflowBase := (1 << 63) / (perSlot >> 1)
	msg.Block.Slot = slotsSinceGenesis.Add(overflowBase)

	// valid block
	msg.Block.ParentRoot = bRoot[:]
	msg.Signature, err = helpers.ComputeDomainAndSign(beaconState, 0, msg.Block, params.BeaconConfig().DomainBeaconProposer, privKeys[proposerIdx])
	require.NoError(t, err)

	genesisTime := time.Now()

	stateGen := stategen.New(db)
	chainService := &mock.ChainService{
		Genesis: time.Unix(genesisTime.Unix()-int64(slotsSinceGenesis.Mul(perSlot)), 0),
		FinalizedCheckPoint: &ethpb.Checkpoint{
			Epoch: 0,
		},
	}
	r := &Service{
		cfg: &Config{
			DB:            db,
			P2P:           p,
			InitialSync:   &mockSync.Sync{IsSyncing: false},
			Chain:         chainService,
			BlockNotifier: chainService.BlockNotifier(),
			StateGen:      stateGen,
		},
		seenBlockCache:      lruwrpr.New(10),
		badBlockCache:       lruwrpr.New(10),
		slotToPendingBlocks: gcache.New(time.Second, 2*time.Second),
		seenPendingBlocks:   make(map[[32]byte]bool),
	}

	buf := new(bytes.Buffer)
	_, err = p.Encoding().EncodeGossip(buf, msg)
	require.NoError(t, err)
	topic := p2p.GossipTypeMapping[reflect.TypeOf(msg)]
	m := &pubsub.Message{
		Message: &pubsubpb.Message{
			Data:  buf.Bytes(),
			Topic: &topic,
		},
	}
	res, err := r.validateBeaconBlockPubSub(ctx, "", m)
	_ = err
	result := res == pubsub.ValidationAccept
	assert.Equal(t, false, result)
}

func TestService_setBadBlock_DoesntSetWithContextErr(t *testing.T) {
	s := Service{}
	s.initCaches()

	root := [32]byte{'b', 'a', 'd'}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	s.setBadBlock(ctx, root)
	if s.hasBadBlock(root) {
		t.Error("Set bad root with cancelled context")
	}
}

func TestService_isBlockQueueable(t *testing.T) {
	currentTime := time.Now().Round(time.Second)
	genesisTime := uint64(currentTime.Unix() - int64(params.BeaconConfig().SecondsPerSlot))
	blockSlot := types.Slot(1)

	// slot time within MAXIMUM_GOSSIP_CLOCK_DISPARITY, so dont queue the block.
	receivedTime := currentTime.Add(-400 * time.Millisecond)
	result := isBlockQueueable(genesisTime, blockSlot, receivedTime)
	assert.Equal(t, false, result)

	// slot time just above MAXIMUM_GOSSIP_CLOCK_DISPARITY, so queue the block.
	receivedTime = currentTime.Add(-600 * time.Millisecond)
	result = isBlockQueueable(genesisTime, blockSlot, receivedTime)
	assert.Equal(t, true, result)
}
