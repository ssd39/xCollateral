package sender

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"os"
	"sort"
	"time"
)

type AttestationState struct {
	LastAttestedBlocks LastAttestedBlocksState
	BlockAttestations  BlockAttestationsState
}

type BlockAttestationsState map[uint64]BlockAttestationState

type BlockAttestationState map[uint64]IndividualAttestationState

type IndividualAttestationState map[uint64]bool

type LastAttestedBlocksState map[uint64]uint64

var AppAttestationState AttestationState

func (state *AttestationState) SetAttested(chainId uint64, block uint64, id uint64) {
	_, found := (*state).BlockAttestations[chainId]
	if !found {
		log.Warn().Msgf("Chain Id not found when trying to set it as attested for chainId %v - block %v - id %v", chainId, block, id)
		return
	}
	_, found = (*state).BlockAttestations[chainId][block]
	if !found {
		log.Warn().Msgf("BlockAttestations not found when trying to set it as attested for chainId %v - block %v - id %v", chainId, block, id)
		return
	}
	_, found = (*state).BlockAttestations[chainId][block][id]
	if !found {
		log.Warn().Msgf("Attestation not found when trying to set it as attested for chain %v and id %v", chainId, id)
		return
	}
	(*state).BlockAttestations[chainId][block][id] = true
	state.UpdateLastAttestedBlocks(chainId)
}

func (state *AttestationState) AddAttestation(chainId uint64, block uint64, id uint64) {
	_, found := (*state).BlockAttestations[chainId]
	if !found {
		(*state).BlockAttestations[chainId] = make(BlockAttestationState)
	}
	_, found = (*state).BlockAttestations[chainId][block]
	if !found {
		(*state).BlockAttestations[chainId][block] = make(IndividualAttestationState)
	}
	_, found = (*state).BlockAttestations[chainId][block][id]
	if !found {
		(*state).BlockAttestations[chainId][block][id] = false
	}
}

func (state *AttestationState) UpdateLastAttestedBlocks(chainId uint64) uint64 {
	_, found := (*state).BlockAttestations[chainId]
	if !found {
		(*state).BlockAttestations = map[uint64]BlockAttestationState{chainId: {}}
	}
	blocks := (*state).BlockAttestations[chainId]
	keys := make([]int, 0, len(blocks))
	for key := range blocks {
		keys = append(keys, int(key))
	}
	sort.Ints(keys)

	var lastAttestedBlock uint64 = 0
	for _, key := range keys {
		fullyAttested := true
		for _, attested := range blocks[uint64(key)] {
			if !attested {
				fullyAttested = false
				break
			}
		}
		if fullyAttested {
			lastAttestedBlock = uint64(key)
		} else {
			break
		}
	}
	state.LastAttestedBlocks[chainId] = lastAttestedBlock
	return lastAttestedBlock
}

func (state *AttestationState) SaveAttestationState() {
	res, err := json.Marshal((*state).LastAttestedBlocks)
	if err != nil {
		log.Error().Msgf("Error when marshaling attestation state %v", err)
	}
	err = os.WriteFile("state.lock", res, 0644)
	if err != nil {
		log.Error().Msgf("Error saving attestation state %v", err)
	}
}

func LoadAttestationState() *AttestationState {
	b, err := os.ReadFile("state.lock")
	if err != nil {
		log.Info().Msgf("Attestation state not found: %v", err)
		AppAttestationState = AttestationState{
			LastAttestedBlocks: make(LastAttestedBlocksState),
			BlockAttestations:  make(BlockAttestationsState),
		}
		return &AppAttestationState
	}
	savedState := LastAttestedBlocksState{}
	err = json.Unmarshal(b, &savedState)
	if err != nil {
		log.Error().Msgf("Error when unmarshaling attestation state %v", err)
	}
	if savedState == nil {
		savedState = make(LastAttestedBlocksState)
	}
	AppAttestationState = AttestationState{
		LastAttestedBlocks: savedState,
		BlockAttestations:  make(BlockAttestationsState),
	}
	return &AppAttestationState
}

func StartSavingState() {
	ticker := time.NewTicker(time.Second * 15)
	for range ticker.C {
		AppAttestationState.SaveAttestationState()
	}
}
