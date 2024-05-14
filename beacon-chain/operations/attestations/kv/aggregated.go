package kv

import (
	"context"
	"runtime"
	"sync"

	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v5/crypto/hash"
	ethpb "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
	attaggregation "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1/attestation/aggregation/attestations"
	"github.com/prysmaticlabs/prysm/v5/runtime/version"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

// AggregateUnaggregatedAttestations aggregates the unaggregated attestations and saves the
// newly aggregated attestations in the pool.
// It tracks the unaggregated attestations that weren't able to aggregate to prevent
// the deletion of unaggregated attestations in the pool.
func (c *AttCaches) AggregateUnaggregatedAttestations(ctx context.Context) error {
	ctx, span := trace.StartSpan(ctx, "operations.attestations.kv.AggregateUnaggregatedAttestations")
	defer span.End()
	unaggregatedAtts, err := c.UnaggregatedAttestations()
	if err != nil {
		return err
	}
	return c.aggregateUnaggregatedAtts(ctx, unaggregatedAtts)
}

func (c *AttCaches) aggregateUnaggregatedAtts(ctx context.Context, unaggregatedAtts []ethpb.Att) error {
	_, span := trace.StartSpan(ctx, "operations.attestations.kv.aggregateUnaggregatedAtts")
	defer span.End()

	attsByVerAndDataRoot := make(map[AttestationId][]ethpb.Att, len(unaggregatedAtts))
	for _, att := range unaggregatedAtts {
		var attDataRoot [32]byte
		var err error
		if att.Version() == version.Phase0 {
			attDataRoot, err = att.GetData().HashTreeRoot()
			if err != nil {
				return err
			}
		} else {
			data := ethpb.CopyAttestationData(att.GetData())
			data.CommitteeIndex = primitives.CommitteeIndex(att.GetCommitteeBitsVal().BitIndices()[0])
			attDataRoot, err = data.HashTreeRoot()
			if err != nil {
				return err
			}
		}

		key := NewAttestationId(att, attDataRoot)
		attsByVerAndDataRoot[key] = append(attsByVerAndDataRoot[key], att)
	}

	// Aggregate unaggregated attestations from the pool and save them in the pool.
	// Track the unaggregated attestations that aren't able to aggregate.
	leftOverUnaggregatedAtt := make(map[AttestationId]bool)

	leftOverUnaggregatedAtt = c.aggregateParallel(attsByVerAndDataRoot, leftOverUnaggregatedAtt)

	// Remove the unaggregated attestations from the pool that were successfully aggregated.
	for _, att := range unaggregatedAtts {
		h, err := hashFn(att)
		if err != nil {
			return err
		}

		if leftOverUnaggregatedAtt[NewAttestationId(att, h)] {
			continue
		}
		if err := c.DeleteUnaggregatedAttestation(att); err != nil {
			return err
		}
	}
	return nil
}

// aggregateParallel aggregates attestations in parallel for `atts` and saves them in the pool,
// returns the unaggregated attestations that weren't able to aggregate.
// Given `n` CPU cores, it creates a channel of size `n` and spawns `n` goroutines to aggregate attestations
func (c *AttCaches) aggregateParallel(atts map[AttestationId][]ethpb.Att, leftOver map[AttestationId]bool) map[AttestationId]bool {
	var leftoverLock sync.Mutex
	wg := sync.WaitGroup{}

	n := runtime.GOMAXPROCS(0) // defaults to the value of runtime.NumCPU
	ch := make(chan []ethpb.Att, n)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for as := range ch {
				aggregated, err := attaggregation.AggregateDisjointOneBitAtts(as)
				if err != nil {
					log.WithError(err).Error("could not aggregate unaggregated attestations")
					continue
				}
				if aggregated == nil {
					log.Error("nil aggregated attestation")
					continue
				}
				if helpers.IsAggregated(aggregated) {
					if err := c.SaveAggregatedAttestations([]ethpb.Att{aggregated}); err != nil {
						log.WithError(err).Error("could not save aggregated attestation")
						continue
					}
				} else {
					h, err := hashFn(aggregated)
					if err != nil {
						log.WithError(err).Error("could not hash attestation")
						continue
					}
					leftoverLock.Lock()
					leftOver[NewAttestationId(aggregated, h)] = true
					leftoverLock.Unlock()
				}
			}
		}()
	}

	for _, as := range atts {
		ch <- as
	}

	close(ch)
	wg.Wait()

	return leftOver
}

// SaveAggregatedAttestation saves an aggregated attestation in cache.
func (c *AttCaches) SaveAggregatedAttestation(att ethpb.Att) error {
	if err := helpers.ValidateNilAttestation(att); err != nil {
		return err
	}
	if !helpers.IsAggregated(att) {
		return errors.New("attestation is not aggregated")
	}
	has, err := c.HasAggregatedAttestation(att)
	if err != nil {
		return err
	}
	if has {
		return nil
	}

	seen, err := c.hasSeenBit(att)
	if err != nil {
		return err
	}
	if seen {
		return nil
	}

	var r [32]byte
	if att.Version() == version.Phase0 {
		r, err = hash.Proto(att.GetData())
		if err != nil {
			return err
		}
	} else {
		data := ethpb.CopyAttestationData(att.GetData())
		data.CommitteeIndex = primitives.CommitteeIndex(att.GetCommitteeBitsVal().BitIndices()[0])
		r, err = hash.Proto(data)
		if err != nil {
			return err
		}
	}
	key := NewAttestationId(att, r)

	copiedAtt := att.Copy()
	c.aggregatedAttLock.Lock()
	defer c.aggregatedAttLock.Unlock()
	atts, ok := c.aggregatedAtt[key]
	if !ok {
		atts := []ethpb.Att{copiedAtt}
		c.aggregatedAtt[key] = atts
		return nil
	}

	atts, err = attaggregation.Aggregate(append(atts, copiedAtt))
	if err != nil {
		return err
	}
	c.aggregatedAtt[key] = atts

	return nil
}

// SaveAggregatedAttestations saves a list of aggregated attestations in cache.
func (c *AttCaches) SaveAggregatedAttestations(atts []ethpb.Att) error {
	for _, att := range atts {
		if err := c.SaveAggregatedAttestation(att); err != nil {
			log.WithError(err).Debug("Could not save aggregated attestation")
			if err := c.DeleteAggregatedAttestation(att); err != nil {
				log.WithError(err).Debug("Could not delete aggregated attestation")
			}
		}
	}
	return nil
}

// AggregatedAttestations returns the aggregated attestations in cache.
func (c *AttCaches) AggregatedAttestations() []ethpb.Att {
	c.aggregatedAttLock.RLock()
	defer c.aggregatedAttLock.RUnlock()

	atts := make([]ethpb.Att, 0)

	for _, a := range c.aggregatedAtt {
		atts = append(atts, a...)
	}

	return atts
}

// AggregatedAttestationsBySlotIndex returns the aggregated attestations in cache,
// filtered by committee index and slot.
func (c *AttCaches) AggregatedAttestationsBySlotIndex(
	ctx context.Context,
	slot primitives.Slot,
	committeeIndex primitives.CommitteeIndex,
) []*ethpb.Attestation {
	_, span := trace.StartSpan(ctx, "operations.attestations.kv.AggregatedAttestationsBySlotIndex")
	defer span.End()

	atts := make([]*ethpb.Attestation, 0)

	c.aggregatedAttLock.RLock()
	defer c.aggregatedAttLock.RUnlock()
	for _, as := range c.aggregatedAtt {
		if as[0].Version() == version.Phase0 && slot == as[0].GetData().Slot && committeeIndex == as[0].GetData().CommitteeIndex {
			for _, a := range as {
				att, ok := a.(*ethpb.Attestation)
				// This will never fail in practice because we asserted the version
				if ok {
					atts = append(atts, att)
				}
			}
		}
	}

	return atts
}

// AggregatedAttestationsBySlotIndexElectra returns the aggregated attestations in cache,
// filtered by committee index and slot.
func (c *AttCaches) AggregatedAttestationsBySlotIndexElectra(
	ctx context.Context,
	slot primitives.Slot,
	committeeIndex primitives.CommitteeIndex,
) []*ethpb.AttestationElectra {
	_, span := trace.StartSpan(ctx, "operations.attestations.kv.AggregatedAttestationsBySlotIndexElectra")
	defer span.End()

	atts := make([]*ethpb.AttestationElectra, 0)

	c.aggregatedAttLock.RLock()
	defer c.aggregatedAttLock.RUnlock()
	for _, as := range c.aggregatedAtt {
		if as[0].Version() == version.Electra && slot == as[0].GetData().Slot && as[0].GetCommitteeBitsVal().BitAt(uint64(committeeIndex)) {
			for _, a := range as {
				att, ok := a.(*ethpb.AttestationElectra)
				// This will never fail in practice because we asserted the version
				if ok {
					atts = append(atts, att)
				}
			}
		}
	}

	return atts
}

// DeleteAggregatedAttestation deletes the aggregated attestations in cache.
func (c *AttCaches) DeleteAggregatedAttestation(att ethpb.Att) error {
	if err := helpers.ValidateNilAttestation(att); err != nil {
		return err
	}
	if !helpers.IsAggregated(att) {
		return errors.New("attestation is not aggregated")
	}
	var r [32]byte
	var err error
	if att.Version() == version.Phase0 {
		r, err = hash.Proto(att.GetData())
		if err != nil {
			return err
		}
	} else {
		data := ethpb.CopyAttestationData(att.GetData())
		data.CommitteeIndex = primitives.CommitteeIndex(att.GetCommitteeBitsVal().BitIndices()[0])
		r, err = hash.Proto(data)
		if err != nil {
			return err
		}
	}
	key := NewAttestationId(att, r)

	if err := c.insertSeenBit(att); err != nil {
		return err
	}

	c.aggregatedAttLock.Lock()
	defer c.aggregatedAttLock.Unlock()
	attList, ok := c.aggregatedAtt[key]
	if !ok {
		return nil
	}

	filtered := make([]ethpb.Att, 0)
	for _, a := range attList {
		if a.Version() != att.Version() {
			continue
		}
		if c, err := att.GetAggregationBits().Contains(a.GetAggregationBits()); err != nil {
			return err
		} else if !c {
			filtered = append(filtered, a)
		}
	}
	if len(filtered) == 0 {
		delete(c.aggregatedAtt, key)
	} else {
		c.aggregatedAtt[key] = filtered
	}

	return nil
}

// HasAggregatedAttestation checks if the input attestations has already existed in cache.
func (c *AttCaches) HasAggregatedAttestation(att ethpb.Att) (bool, error) {
	if err := helpers.ValidateNilAttestation(att); err != nil {
		return false, err
	}
	var r [32]byte
	var err error
	if att.Version() == version.Phase0 {
		r, err = hash.Proto(att.GetData())
		if err != nil {
			return false, err
		}
	} else {
		data := ethpb.CopyAttestationData(att.GetData())
		ci := primitives.CommitteeIndex(att.GetCommitteeBitsVal().BitIndices()[0])
		data.CommitteeIndex = ci
		r, err = hash.Proto(data)
		if err != nil {
			return false, err
		}
	}
	key := NewAttestationId(att, r)

	c.aggregatedAttLock.RLock()
	defer c.aggregatedAttLock.RUnlock()
	if atts, ok := c.aggregatedAtt[key]; ok {
		for _, a := range atts {
			if a.Version() != att.Version() {
				continue
			}
			if c, err := a.GetAggregationBits().Contains(att.GetAggregationBits()); err != nil {
				return false, err
			} else if c {
				return true, nil
			}
		}
	}

	c.blockAttLock.RLock()
	defer c.blockAttLock.RUnlock()
	if atts, ok := c.blockAtt[key]; ok {
		for _, a := range atts {
			if a.Version() != att.Version() {
				continue
			}
			if c, err := a.GetAggregationBits().Contains(att.GetAggregationBits()); err != nil {
				return false, err
			} else if c {
				return true, nil
			}
		}
	}

	return false, nil
}

// AggregatedAttestationCount returns the number of aggregated attestations key in the pool.
func (c *AttCaches) AggregatedAttestationCount() int {
	c.aggregatedAttLock.RLock()
	defer c.aggregatedAttLock.RUnlock()
	return len(c.aggregatedAtt)
}
