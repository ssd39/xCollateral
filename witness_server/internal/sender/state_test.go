package sender

import (
	"testing"
)

func TestState_SetAttested(t *testing.T) {
	AppAttestationState = AttestationState{
		LastAttestedBlocks: make(LastAttestedBlocksState),
		BlockAttestations:  make(BlockAttestationsState),
	}

	chainId := uint64(1)
	block := uint64(1000)
	id := uint64(0)
	AppAttestationState.BlockAttestations[chainId] = make(BlockAttestationState)
	AppAttestationState.BlockAttestations[chainId][block] = make(IndividualAttestationState)
	AppAttestationState.BlockAttestations[chainId][block][id] = false
	AppAttestationState.SetAttested(chainId, block, id)
	expected := true
	got := AppAttestationState.BlockAttestations[chainId][block][id]
	if got != expected {
		t.Errorf("expected %+v got %+v", expected, got)
	}
}

func TestState_AddAttestation(t *testing.T) {
	AppAttestationState = AttestationState{
		LastAttestedBlocks: make(LastAttestedBlocksState),
		BlockAttestations:  make(BlockAttestationsState),
	}

	chainId := uint64(2)
	block := uint64(2000)
	id := uint64(1)
	AppAttestationState.AddAttestation(chainId, block, id)
	got, found := AppAttestationState.BlockAttestations[chainId][block][id]
	expected := false
	if !found {
		t.Errorf("expected %+v got %+v", true, found)
	}
	if got != expected {
		t.Errorf("expected %+v got %+v", expected, got)
	}
}

func TestState_UpdateLastAttestedBlocks(t *testing.T) {
	AppAttestationState = AttestationState{
		LastAttestedBlocks: make(LastAttestedBlocksState),
		BlockAttestations:  make(BlockAttestationsState),
	}

	chainId := uint64(2)
	block := uint64(2000)
	id := uint64(1)
	got := AppAttestationState.UpdateLastAttestedBlocks(chainId)
	expected := uint64(0)
	if got != expected {
		t.Errorf("expected %+v got %+v", expected, got)
	}

	AppAttestationState.AddAttestation(chainId, block, id)
	got = AppAttestationState.UpdateLastAttestedBlocks(chainId)
	if got != expected {
		t.Errorf("expected %+v got %+v", expected, got)
	}

	AppAttestationState.SetAttested(chainId, block, id)
	got = AppAttestationState.UpdateLastAttestedBlocks(chainId)
	expectedBlock := uint64(2000)
	if got != expectedBlock {
		t.Errorf("expected %+v got %+v", expectedBlock, got)
	}

	block = uint64(1980)
	AppAttestationState.AddAttestation(chainId, block, id)
	AppAttestationState.SetAttested(chainId, block, id)
	got = AppAttestationState.UpdateLastAttestedBlocks(chainId)
	if got != expectedBlock {
		t.Errorf("expected %+v got %+v", expectedBlock, got)
	}

	block = uint64(2050)
	expectedBlock = uint64(2050)
	AppAttestationState.AddAttestation(chainId, block, id)
	AppAttestationState.SetAttested(chainId, block, id)
	got = AppAttestationState.UpdateLastAttestedBlocks(chainId)
	if got != expectedBlock {
		t.Errorf("expected %+v got %+v", expectedBlock, got)
	}
}
