package apimiddleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/api/gateway/apimiddleware"
	"github.com/prysmaticlabs/prysm/encoding/bytesutil"
	ethpbv2 "github.com/prysmaticlabs/prysm/proto/eth/v2"
)

// https://ethereum.github.io/beacon-apis/#/Beacon/submitPoolAttestations expects posting a top-level array.
// We make it more proto-friendly by wrapping it in a struct with a 'data' field.
func wrapAttestationsArray(endpoint apimiddleware.Endpoint, _ http.ResponseWriter, req *http.Request) apimiddleware.ErrorJson {
	if _, ok := endpoint.PostRequest.(*submitAttestationRequestJson); ok {
		atts := make([]*attestationJson, 0)
		if err := json.NewDecoder(req.Body).Decode(&atts); err != nil {
			return apimiddleware.InternalServerErrorWithMessage(err, "could not decode body")
		}
		j := &submitAttestationRequestJson{Data: atts}
		b, err := json.Marshal(j)
		if err != nil {
			return apimiddleware.InternalServerErrorWithMessage(err, "could not marshal wrapped body")
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(b))
	}
	return nil
}

// Some endpoints e.g. https://ethereum.github.io/beacon-apis/#/Validator/getAttesterDuties expect posting a top-level array.
// We make it more proto-friendly by wrapping it in a struct with an 'Index' field.
func wrapValidatorIndicesArray(endpoint apimiddleware.Endpoint, _ http.ResponseWriter, req *http.Request) apimiddleware.ErrorJson {
	if _, ok := endpoint.PostRequest.(*dutiesRequestJson); ok {
		indices := make([]string, 0)
		if err := json.NewDecoder(req.Body).Decode(&indices); err != nil {
			return apimiddleware.InternalServerErrorWithMessage(err, "could not decode body")
		}
		j := &dutiesRequestJson{Index: indices}
		b, err := json.Marshal(j)
		if err != nil {
			return apimiddleware.InternalServerErrorWithMessage(err, "could not marshal wrapped body")
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(b))
	}
	return nil
}

// https://ethereum.github.io/beacon-apis/#/Validator/publishAggregateAndProofs expects posting a top-level array.
// We make it more proto-friendly by wrapping it in a struct with a 'data' field.
func wrapSignedAggregateAndProofArray(endpoint apimiddleware.Endpoint, _ http.ResponseWriter, req *http.Request) apimiddleware.ErrorJson {
	if _, ok := endpoint.PostRequest.(*submitAggregateAndProofsRequestJson); ok {
		data := make([]*signedAggregateAttestationAndProofJson, 0)
		if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
			return apimiddleware.InternalServerErrorWithMessage(err, "could not decode body")
		}
		j := &submitAggregateAndProofsRequestJson{Data: data}
		b, err := json.Marshal(j)
		if err != nil {
			return apimiddleware.InternalServerErrorWithMessage(err, "could not marshal wrapped body")
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(b))
	}
	return nil
}

// https://ethereum.github.io/beacon-apis/#/Validator/prepareBeaconCommitteeSubnet expects posting a top-level array.
// We make it more proto-friendly by wrapping it in a struct with a 'data' field.
func wrapBeaconCommitteeSubscriptionsArray(endpoint apimiddleware.Endpoint, _ http.ResponseWriter, req *http.Request) apimiddleware.ErrorJson {
	if _, ok := endpoint.PostRequest.(*submitBeaconCommitteeSubscriptionsRequestJson); ok {
		data := make([]*beaconCommitteeSubscribeJson, 0)
		if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
			return apimiddleware.InternalServerErrorWithMessage(err, "could not decode body")
		}
		j := &submitBeaconCommitteeSubscriptionsRequestJson{Data: data}
		b, err := json.Marshal(j)
		if err != nil {
			return apimiddleware.InternalServerErrorWithMessage(err, "could not marshal wrapped body")
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(b))
	}
	return nil
}

// https://ethereum.github.io/beacon-APIs/#/Validator/prepareSyncCommitteeSubnets expects posting a top-level array.
// We make it more proto-friendly by wrapping it in a struct with a 'data' field.
func wrapSyncCommitteeSubscriptionsArray(endpoint apimiddleware.Endpoint, _ http.ResponseWriter, req *http.Request) apimiddleware.ErrorJson {
	if _, ok := endpoint.PostRequest.(*submitSyncCommitteeSubscriptionRequestJson); ok {
		data := make([]*syncCommitteeSubscriptionJson, 0)
		if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
			return apimiddleware.InternalServerErrorWithMessage(err, "could not decode body")
		}
		j := &submitSyncCommitteeSubscriptionRequestJson{Data: data}
		b, err := json.Marshal(j)
		if err != nil {
			return apimiddleware.InternalServerErrorWithMessage(err, "could not marshal wrapped body")
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(b))
	}
	return nil
}

// https://ethereum.github.io/beacon-APIs/#/Beacon/submitPoolSyncCommitteeSignatures expects posting a top-level array.
// We make it more proto-friendly by wrapping it in a struct with a 'data' field.
func wrapSyncCommitteeSignaturesArray(endpoint apimiddleware.Endpoint, _ http.ResponseWriter, req *http.Request) apimiddleware.ErrorJson {
	if _, ok := endpoint.PostRequest.(*submitSyncCommitteeSignaturesRequestJson); ok {
		data := make([]*syncCommitteeMessageJson, 0)
		if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
			return apimiddleware.InternalServerErrorWithMessage(err, "could not decode body")
		}
		j := &submitSyncCommitteeSignaturesRequestJson{Data: data}
		b, err := json.Marshal(j)
		if err != nil {
			return apimiddleware.InternalServerErrorWithMessage(err, "could not marshal wrapped body")
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(b))
	}
	return nil
}

// https://ethereum.github.io/beacon-APIs/#/Validator/publishContributionAndProofs expects posting a top-level array.
// We make it more proto-friendly by wrapping it in a struct with a 'data' field.
func wrapSignedContributionAndProofsArray(endpoint apimiddleware.Endpoint, _ http.ResponseWriter, req *http.Request) apimiddleware.ErrorJson {
	if _, ok := endpoint.PostRequest.(*submitContributionAndProofsRequestJson); ok {
		data := make([]*signedContributionAndProofJson, 0)
		if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
			return apimiddleware.InternalServerErrorWithMessage(err, "could not decode body")
		}
		j := &submitContributionAndProofsRequestJson{Data: data}
		b, err := json.Marshal(j)
		if err != nil {
			return apimiddleware.InternalServerErrorWithMessage(err, "could not marshal wrapped body")
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(b))
	}
	return nil
}

// Posted graffiti needs to have length of 32 bytes, but client is allowed to send data of any length.
func prepareGraffiti(endpoint apimiddleware.Endpoint, _ http.ResponseWriter, _ *http.Request) apimiddleware.ErrorJson {
	if block, ok := endpoint.PostRequest.(*signedBeaconBlockContainerJson); ok {
		b := bytesutil.ToBytes32([]byte(block.Message.Body.Graffiti))
		block.Message.Body.Graffiti = hexutil.Encode(b[:])
	}
	return nil
}

type tempSyncCommitteesResponseJson struct {
	Data *tempSyncCommitteeValidatorsJson `json:"data"`
}

type tempSyncCommitteeValidatorsJson struct {
	Validators          []string                              `json:"validators"`
	ValidatorAggregates []*tempSyncSubcommitteeValidatorsJson `json:"validator_aggregates"`
}

type tempSyncSubcommitteeValidatorsJson struct {
	Validators []string `json:"validators"`
}

// https://ethereum.github.io/beacon-APIs/?urls.primaryName=v2.0.0#/Beacon/getEpochSyncCommittees returns validator_aggregates as a nested array.
// grpc-gateway returns a struct with nested fields which we have to transform into a plain 2D array.
func prepareValidatorAggregates(body []byte, responseContainer interface{}) (bool, apimiddleware.ErrorJson) {
	tempContainer := &tempSyncCommitteesResponseJson{}
	if err := json.Unmarshal(body, tempContainer); err != nil {
		return false, apimiddleware.InternalServerErrorWithMessage(err, "could not unmarshal response into temp container")
	}
	container, ok := responseContainer.(*syncCommitteesResponseJson)
	if !ok {
		return false, apimiddleware.InternalServerError(errors.New("container is not of the correct type"))
	}

	container.Data = &syncCommitteeValidatorsJson{}
	container.Data.Validators = tempContainer.Data.Validators
	container.Data.ValidatorAggregates = make([][]string, len(tempContainer.Data.ValidatorAggregates))
	for i, srcValAgg := range tempContainer.Data.ValidatorAggregates {
		dstValAgg := make([]string, len(srcValAgg.Validators))
		copy(dstValAgg, tempContainer.Data.ValidatorAggregates[i].Validators)
		container.Data.ValidatorAggregates[i] = dstValAgg
	}

	return true, nil
}

type phase0BlockResponseJson struct {
	Version string                          `json:"version"`
	Data    *signedBeaconBlockContainerJson `json:"data"`
}

type altairBlockResponseJson struct {
	Version string                                `json:"version"`
	Data    *signedBeaconBlockAltairContainerJson `json:"data"`
}

func serializeV2Block(response interface{}) (bool, []byte, apimiddleware.ErrorJson) {
	respContainer, ok := response.(*blockV2ResponseJson)
	if !ok {
		return false, nil, apimiddleware.InternalServerError(errors.New("container is not of the correct type"))
	}

	var actualRespContainer interface{}
	if strings.EqualFold(respContainer.Version, strings.ToLower(ethpbv2.Version_PHASE0.String())) {
		actualRespContainer = &phase0BlockResponseJson{
			Version: respContainer.Version,
			Data: &signedBeaconBlockContainerJson{
				Message:   respContainer.Data.Phase0Block,
				Signature: respContainer.Data.Signature,
			},
		}
	} else {
		actualRespContainer = &altairBlockResponseJson{
			Version: respContainer.Version,
			Data: &signedBeaconBlockAltairContainerJson{
				Message:   respContainer.Data.AltairBlock,
				Signature: respContainer.Data.Signature,
			},
		}
	}

	j, err := json.Marshal(actualRespContainer)
	if err != nil {
		return false, nil, apimiddleware.InternalServerErrorWithMessage(err, "could not marshal response")
	}
	return true, j, nil
}

type phase0StateResponseJson struct {
	Version string           `json:"version"`
	Data    *beaconStateJson `json:"data"`
}

type altairStateResponseJson struct {
	Version string             `json:"version"`
	Data    *beaconStateV2Json `json:"data"`
}

func serializeV2State(response interface{}) (bool, []byte, apimiddleware.ErrorJson) {
	respContainer, ok := response.(*beaconStateV2ResponseJson)
	if !ok {
		return false, nil, apimiddleware.InternalServerError(errors.New("container is not of the correct type"))
	}

	var actualRespContainer interface{}
	if strings.EqualFold(respContainer.Version, strings.ToLower(ethpbv2.Version_PHASE0.String())) {
		actualRespContainer = &phase0StateResponseJson{
			Version: respContainer.Version,
			Data:    respContainer.Data.Phase0State,
		}
	} else {
		actualRespContainer = &altairStateResponseJson{
			Version: respContainer.Version,
			Data:    respContainer.Data.AltairState,
		}
	}

	j, err := json.Marshal(actualRespContainer)
	if err != nil {
		return false, nil, apimiddleware.InternalServerErrorWithMessage(err, "could not marshal response")
	}
	return true, j, nil
}

type phase0ProduceBlockResponseJson struct {
	Version string           `json:"version"`
	Data    *beaconBlockJson `json:"data"`
}

type altairProduceBlockResponseJson struct {
	Version string                 `json:"version"`
	Data    *beaconBlockAltairJson `json:"data"`
}

func serializeProducedV2Block(response interface{}) (bool, []byte, apimiddleware.ErrorJson) {
	respContainer, ok := response.(*produceBlockResponseV2Json)
	if !ok {
		return false, nil, apimiddleware.InternalServerError(errors.New("container is not of the correct type"))
	}

	var actualRespContainer interface{}
	if strings.EqualFold(respContainer.Version, strings.ToLower(ethpbv2.Version_PHASE0.String())) {
		actualRespContainer = &phase0ProduceBlockResponseJson{
			Version: respContainer.Version,
			Data:    respContainer.Data.Phase0Block,
		}
	} else {
		actualRespContainer = &altairProduceBlockResponseJson{
			Version: respContainer.Version,
			Data:    respContainer.Data.AltairBlock,
		}
	}

	j, err := json.Marshal(actualRespContainer)
	if err != nil {
		return false, nil, apimiddleware.InternalServerErrorWithMessage(err, "could not marshal response")
	}
	return true, j, nil
}
