package attestate

import (
	"math/big"
	"peersyst/bridge-witness-go/internal/chains"
	"peersyst/bridge-witness-go/internal/sender"
	"testing"
	"time"
)

func TestAttestate_resendToQueue(t *testing.T) {
	AttestateInSideChainQueue = make(chan *interface{}, 1000)
	AttestateInMainChainQueue = make(chan *interface{}, 1000)

	resendToQueue(AttestateInSideChainQueue, nil)
	if len(AttestateInSideChainQueue) != 1 {
		t.Errorf("error: resendToQUeue should add elem to queue expected %+v got %+v", 1, len(AttestateInSideChainQueue))
	}

	resendToQueue(AttestateInMainChainQueue, nil)
	if len(AttestateInMainChainQueue) != 1 {
		t.Errorf("error: resendToQUeue should add elem to queue expected %+v got %+v", 1, len(AttestateInMainChainQueue))
	}
}

func createClaim(block, claimId uint64, sender, amount, destination string) interface{} {
	return &struct {
		Block       uint64
		ClaimId     uint64
		Sender      string
		Amount      string
		Destination string
		Nonce       int
		Fee         int
		BridgeId    string
	}{Block: block, ClaimId: claimId, Sender: sender, Amount: amount, Destination: destination, Nonce: 1, Fee: 10}
}

func TestAttestate_AttestateInEvm_Claim(t *testing.T) {
	AttestateInSideChainQueue = make(chan *interface{}, 1000)
	go AttestateInSideChain(AttestateInSideChainQueue)

	chains.StartEvmTestProvider(0, 0, false, big.NewInt(117), nil)
	chains.StartXrpTestProvider(0, 0, true, big.NewInt(144), nil)

	claim := createClaim(100, 1, "", "", "")
	AttestateInSideChainQueue <- &claim
	time.Sleep(time.Millisecond)
	if len(AttestateInSideChainQueue) != 0 {
		t.Errorf("error: queue should be empty as not in signer list expected %+v got %+v", 0, len(AttestateInSideChainQueue))
	}
	if chains.EvmTestProvider.GetUnattestedClaimCalledTimes != 0 {
		t.Errorf("error: should not call checkSideChainClaim as not in signer list expected %+v got %+v", 0, chains.EvmTestProvider.GetUnattestedClaimCalledTimes)
	}

	chains.StartEvmTestProvider(0, 0, true, big.NewInt(117), nil)
	AttestateInSideChainQueue <- &claim
	time.Sleep(time.Millisecond)
	if chains.EvmTestProvider.GetUnattestedClaimCalledTimes != 1 {
		t.Errorf("error: should call checkSideChainClaim and not continue expected %+v got %+v", 1, chains.EvmTestProvider.GetUnattestedClaimCalledTimes)
	}
	if chains.EvmTestProvider.GetAttestClaimTxCalledTimes != 0 {
		t.Errorf("error: should not call GetAttestClaimTransaction expected %+v got %+v", 0, chains.EvmTestProvider.GetAttestClaimTxCalledTimes)
	}

	claim = createClaim(100, 0, "", "", "")
	AttestateInSideChainQueue <- &claim
	time.Sleep(time.Millisecond)
	if chains.EvmTestProvider.GetUnattestedClaimCalledTimes != 2 {
		t.Errorf("error: should call checkSideChainClaim expected %+v got %+v", 1, chains.EvmTestProvider.GetUnattestedClaimCalledTimes)
	}
	if chains.EvmTestProvider.GetAttestClaimTxCalledTimes != 0 {
		t.Errorf("error: should not call GetAttestClaimTransaction expected %+v got %+v", 0, chains.EvmTestProvider.GetAttestClaimTxCalledTimes)
	}
	close(AttestateInSideChainQueue)

	AttestateInSideChainQueue = make(chan *interface{}, 1000)
	go AttestateInSideChain(AttestateInSideChainQueue)
	claim = createClaim(100, 2, "r3an6Cz2MgHQT9q3Kj3QzwQv9ARkgkkxqo", "100", "rs99jCuSAjrXzdebKm1AgpErz9M2FwHQCE")
	AttestateInSideChainQueue <- &claim
	time.Sleep(time.Millisecond)
	if chains.EvmTestProvider.GetAttestClaimTxCalledTimes != 1 {
		t.Errorf("error: should call GetAttestClaimTransaction expected %+v got %+v", 1, chains.EvmTestProvider.GetAttestClaimTxCalledTimes)
	}
	close(AttestateInSideChainQueue)

	sender.AppAttestationState = sender.AttestationState{
		LastAttestedBlocks: make(sender.LastAttestedBlocksState),
		BlockAttestations:  make(sender.BlockAttestationsState),
	}
	AttestateInSideChainQueue = make(chan *interface{}, 1000)
	go AttestateInSideChain(AttestateInSideChainQueue)
	claim = createClaim(100, 2, "r3an6Cz2MgHQT9q3Kj3QzwQv9ARkgkkxqo", "100", "rDTZ46LPHmKSpEAEUEbFuFjy6D4C3ud8GC")
	AttestateInSideChainQueue <- &claim
	sender.BroadcastTransactionInQueue = make(chan *sender.BroadcastTransactionQueueItem, 3000)
	time.Sleep(time.Millisecond * 2)
	if chains.EvmTestProvider.GetAttestClaimTxCalledTimes != 2 {
		t.Errorf("error: should call GetAttestClaimTransaction expected %+v got %+v", 2, chains.EvmTestProvider.GetAttestClaimTxCalledTimes)
	}
	if len(sender.BroadcastTransactionInQueue) != 1 {
		t.Errorf("error: tx should be in broadcast queue expected %+v got %+v", 1, len(sender.BroadcastTransactionInQueue))
	}
	close(sender.BroadcastTransactionInQueue)
}

func TestAttestate_AttestateInXrp_Claim(t *testing.T) {
	AttestateInMainChainQueue = make(chan *interface{}, 1000)
	go AttestateInMainChain(AttestateInMainChainQueue)

	chains.StartXrpTestProvider(0, 0, false, big.NewInt(117), nil)
	chains.StartEvmTestProvider(0, 0, true, big.NewInt(144), nil)

	claim := createClaim(100, 1, "0xc2cD370bAdC28A01682394E8072824c1D7300D96", "100", "0x177adf17f5ac5df0178a24ba5b805a88a7a4be2a")
	AttestateInMainChainQueue <- &claim
	time.Sleep(time.Millisecond)
	if len(AttestateInMainChainQueue) != 0 {
		t.Errorf("error: queue should be empty as not in signer list expected %+v got %+v", 0, len(AttestateInMainChainQueue))
	}
	if chains.XrpTestProvider.GetUnattestedClaimCalledTimes != 0 {
		t.Errorf("error: should not call checkMainChainClaim as not in signer list expected %+v got %+v", 0, chains.XrpTestProvider.GetUnattestedClaimCalledTimes)
	}

	chains.StartXrpTestProvider(0, 0, true, big.NewInt(117), nil)
	AttestateInMainChainQueue <- &claim
	time.Sleep(time.Millisecond)
	if chains.XrpTestProvider.GetUnattestedClaimCalledTimes != 1 {
		t.Errorf("error: should call checkMainChainClaim and not continue expected %+v got %+v", 1, chains.XrpTestProvider.GetUnattestedClaimCalledTimes)
	}
	if chains.XrpTestProvider.GetAttestClaimTxCalledTimes != 0 {
		t.Errorf("error: should not call GetAttestClaimTransaction expected %+v got %+v", 0, chains.XrpTestProvider.GetAttestClaimTxCalledTimes)
	}

	claim = createClaim(100, 0, "", "", "")
	AttestateInMainChainQueue <- &claim
	time.Sleep(time.Millisecond)
	if chains.XrpTestProvider.GetUnattestedClaimCalledTimes != 2 {
		t.Errorf("error: should call checkMainChainClaim expected %+v got %+v", 1, chains.XrpTestProvider.GetUnattestedClaimCalledTimes)
	}
	if chains.XrpTestProvider.GetAttestClaimTxCalledTimes != 0 {
		t.Errorf("error: should not call GetAttestClaimTransaction expected %+v got %+v", 0, chains.XrpTestProvider.GetAttestClaimTxCalledTimes)
	}
	close(AttestateInMainChainQueue)

	AttestateInMainChainQueue = make(chan *interface{}, 1000)
	go AttestateInMainChain(AttestateInMainChainQueue)
	claim = createClaim(100, 1, "0x4DBeE27B94c970B6A7916628236ad6D9369a4518", "100", "0x177adf17f5ac5df0178a24ba5b805a88a7a4be2a")
	AttestateInMainChainQueue <- &claim
	time.Sleep(time.Millisecond)
	if chains.XrpTestProvider.GetAttestClaimTxCalledTimes != 1 {
		t.Errorf("error: should call GetAttestClaimTransaction expected %+v got %+v", 1, chains.XrpTestProvider.GetAttestClaimTxCalledTimes)
	}
	close(AttestateInMainChainQueue)

	sender.AppAttestationState = sender.AttestationState{
		LastAttestedBlocks: make(sender.LastAttestedBlocksState),
		BlockAttestations:  make(sender.BlockAttestationsState),
	}
	AttestateInMainChainQueue = make(chan *interface{}, 1000)
	go AttestateInMainChain(AttestateInMainChainQueue)
	claim = createClaim(100, 1, "0x4DBeE27B94c970B6A7916628236ad6D9369a4518", "100", "0xc2cD370bAdC28A01682394E8072824c1D7300D96")
	AttestateInMainChainQueue <- &claim
	sender.BroadcastTransactionInQueue = make(chan *sender.BroadcastTransactionQueueItem, 3000)
	time.Sleep(time.Millisecond * 2)
	if chains.XrpTestProvider.GetAttestClaimTxCalledTimes != 2 {
		t.Errorf("error: should call GetAttestClaimTransaction expected %+v got %+v", 2, chains.XrpTestProvider.GetAttestClaimTxCalledTimes)
	}
	if len(sender.BroadcastTransactionInQueue) != 1 {
		t.Errorf("error: tx should be in broadcast queue expected %+v got %+v", 1, len(sender.BroadcastTransactionInQueue))
	}
	close(sender.BroadcastTransactionInQueue)
}

func createCreateAccount(block uint64, sender, amount, destination, signatureReward string) interface{} {
	return &struct {
		Block           uint64
		Sender          string
		Amount          string
		Destination     string
		SignatureReward string
		Nonce           int
		Fee             int
		BridgeId        string
	}{Block: block, Sender: sender, Amount: amount, Destination: destination, SignatureReward: signatureReward, Nonce: 1, Fee: 10}
}

func TestAttestate_AttestateInEvm_CreateAccount(t *testing.T) {
	AttestateInSideChainQueue = make(chan *interface{}, 1000)
	go AttestateInSideChain(AttestateInSideChainQueue)

	chains.StartEvmTestProvider(0, 0, false, big.NewInt(117), nil)
	chains.StartXrpTestProvider(0, 0, true, big.NewInt(144), nil)

	claim := createCreateAccount(1000, "", "100", "", "1")
	AttestateInSideChainQueue <- &claim
	time.Sleep(time.Millisecond)
	if len(AttestateInSideChainQueue) != 0 {
		t.Errorf("error: queue should be empty as not in signer list expected %+v got %+v", 0, len(AttestateInSideChainQueue))
	}
	if chains.EvmTestProvider.CheckAccountCreatedCalledTimes != 0 {
		t.Errorf("error: should not call checkSideChainClaim as not in signer list expected %+v got %+v", 0, chains.EvmTestProvider.CheckAccountCreatedCalledTimes)
	}

	chains.StartEvmTestProvider(0, 0, true, big.NewInt(117), nil)
	AttestateInSideChainQueue <- &claim
	time.Sleep(time.Millisecond)
	if chains.EvmTestProvider.CheckAccountCreatedCalledTimes != 1 {
		t.Errorf("error: should call checkSideChainClaim and not continue expected %+v got %+v", 1, chains.EvmTestProvider.CheckAccountCreatedCalledTimes)
	}
	if chains.EvmTestProvider.GetAttestCreateAccountTxCalledTimes != 0 {
		t.Errorf("error: should not call GetAttestCreateAccountTransaction expected %+v got %+v", 0, chains.EvmTestProvider.GetAttestCreateAccountTxCalledTimes)
	}

	claim = createCreateAccount(1000, "", "100", "error-account", "1")
	AttestateInSideChainQueue <- &claim
	time.Sleep(time.Millisecond)
	if chains.EvmTestProvider.CheckAccountCreatedCalledTimes != 2 {
		t.Errorf("error: should call checkSideChainClaim expected %+v got %+v", 1, chains.EvmTestProvider.CheckAccountCreatedCalledTimes)
	}
	if chains.EvmTestProvider.GetAttestCreateAccountTxCalledTimes != 0 {
		t.Errorf("error: should not call GetAttestCreateAccountTransaction expected %+v got %+v", 0, chains.EvmTestProvider.GetAttestCreateAccountTxCalledTimes)
	}
	close(AttestateInSideChainQueue)

	AttestateInSideChainQueue = make(chan *interface{}, 1000)
	go AttestateInSideChain(AttestateInSideChainQueue)
	claim = createCreateAccount(1000, "rDTZ46LPHmKSpEAEUEbFuFjy6D4C3ud8GC", "100", "rs99jCuSAjrXzdebKm1AgpErz9M2FwHQCE", "1")
	AttestateInSideChainQueue <- &claim
	time.Sleep(time.Millisecond)
	if chains.EvmTestProvider.GetAttestCreateAccountTxCalledTimes != 1 {
		t.Errorf("error: should call GetAttestCreateAccountTransaction expected %+v got %+v", 1, chains.EvmTestProvider.GetAttestCreateAccountTxCalledTimes)
	}
	close(AttestateInSideChainQueue)

	sender.AppAttestationState = sender.AttestationState{
		LastAttestedBlocks: make(sender.LastAttestedBlocksState),
		BlockAttestations:  make(sender.BlockAttestationsState),
	}
	AttestateInSideChainQueue = make(chan *interface{}, 1000)
	go AttestateInSideChain(AttestateInSideChainQueue)
	claim = createCreateAccount(1000, "rDTZ46LPHmKSpEAEUEbFuFjy6D4C3ud8GC", "100", "rDTZ46LPHmKSpEAEUEbFuFjy6D4C3ud8GC", "1")
	AttestateInSideChainQueue <- &claim
	sender.BroadcastTransactionInQueue = make(chan *sender.BroadcastTransactionQueueItem, 3000)
	time.Sleep(time.Millisecond * 2)
	if chains.EvmTestProvider.GetAttestCreateAccountTxCalledTimes != 2 {
		t.Errorf("error: should call GetAttestCreateAccountTransaction expected %+v got %+v", 2, chains.EvmTestProvider.GetAttestCreateAccountTxCalledTimes)
	}
	if len(sender.BroadcastTransactionInQueue) != 1 {
		t.Errorf("error: tx should be in broadcast queue expected %+v got %+v", 1, len(sender.BroadcastTransactionInQueue))
	}
	close(sender.BroadcastTransactionInQueue)
}

func TestAttestate_AttestateInXrp_CreateAccount(t *testing.T) {
	AttestateInMainChainQueue = make(chan *interface{}, 1000)
	go AttestateInMainChain(AttestateInMainChainQueue)

	chains.StartXrpTestProvider(0, 0, false, big.NewInt(117), nil)
	chains.StartEvmTestProvider(0, 0, true, big.NewInt(144), nil)

	claim := createCreateAccount(1000, "0x4DBeE27B94c970B6A7916628236ad6D9369a4518", "100", "found-attested", "1")
	AttestateInMainChainQueue <- &claim
	time.Sleep(time.Millisecond)
	if len(AttestateInMainChainQueue) != 0 {
		t.Errorf("error: queue should be empty as not in signer list expected %+v got %+v", 0, len(AttestateInMainChainQueue))
	}
	if chains.XrpTestProvider.CheckAccountCreatedCalledTimes != 0 {
		t.Errorf("error: should not call checkMainChainClaim as not in signer list expected %+v got %+v", 0, chains.XrpTestProvider.CheckAccountCreatedCalledTimes)
	}

	chains.StartXrpTestProvider(0, 0, true, big.NewInt(117), nil)
	AttestateInMainChainQueue <- &claim
	time.Sleep(time.Millisecond)
	if chains.XrpTestProvider.CheckAccountCreatedCalledTimes != 1 {
		t.Errorf("error: should call checkMainChainClaim and not continue expected %+v got %+v", 1, chains.XrpTestProvider.CheckAccountCreatedCalledTimes)
	}
	if chains.XrpTestProvider.GetAttestCreateAccountTxCalledTimes != 0 {
		t.Errorf("error: should not call GetAttestCreateAccountTransaction expected %+v got %+v", 0, chains.XrpTestProvider.GetAttestCreateAccountTxCalledTimes)
	}

	claim = createCreateAccount(1000, "0x4DBeE27B94c970B6A7916628236ad6D9369a4518", "100", "error-account", "1")
	AttestateInMainChainQueue <- &claim
	time.Sleep(time.Millisecond)
	if chains.XrpTestProvider.CheckAccountCreatedCalledTimes != 2 {
		t.Errorf("error: should call checkMainChainClaim expected %+v got %+v", 1, chains.XrpTestProvider.CheckAccountCreatedCalledTimes)
	}
	if chains.XrpTestProvider.GetAttestCreateAccountTxCalledTimes != 0 {
		t.Errorf("error: should not call GetAttestCreateAccountTransaction expected %+v got %+v", 0, chains.XrpTestProvider.GetAttestCreateAccountTxCalledTimes)
	}
	close(AttestateInMainChainQueue)

	AttestateInMainChainQueue = make(chan *interface{}, 1000)
	go AttestateInMainChain(AttestateInMainChainQueue)
	claim = createCreateAccount(1000, "0x4DBeE27B94c970B6A7916628236ad6D9369a4518", "100", "0x177adf17f5ac5df0178a24ba5b805a88a7a4be2a", "1")
	AttestateInMainChainQueue <- &claim
	time.Sleep(time.Millisecond)
	if chains.XrpTestProvider.GetAttestCreateAccountTxCalledTimes != 1 {
		t.Errorf("error: should call GetAttestCreateAccountTransaction expected %+v got %+v", 1, chains.XrpTestProvider.GetAttestCreateAccountTxCalledTimes)
	}
	close(AttestateInMainChainQueue)

	sender.AppAttestationState = sender.AttestationState{
		LastAttestedBlocks: make(sender.LastAttestedBlocksState),
		BlockAttestations:  make(sender.BlockAttestationsState),
	}
	AttestateInMainChainQueue = make(chan *interface{}, 1000)
	go AttestateInMainChain(AttestateInMainChainQueue)
	claim = createCreateAccount(1000, "0x4DBeE27B94c970B6A7916628236ad6D9369a4518", "100", "0x4DBeE27B94c970B6A7916628236ad6D9369a4518", "1")
	AttestateInMainChainQueue <- &claim
	sender.CreateAccountQueue = make(chan *sender.CreateAccountQueueItem, 3000)
	time.Sleep(time.Millisecond * 2)
	if chains.XrpTestProvider.GetAttestCreateAccountTxCalledTimes != 2 {
		t.Errorf("error: should call GetAttestCreateAccountTransaction expected %+v got %+v", 2, chains.XrpTestProvider.GetAttestCreateAccountTxCalledTimes)
	}
	if len(sender.CreateAccountQueue) != 1 {
		t.Errorf("error: tx should be in broadcast queue expected %+v got %+v", 1, len(sender.CreateAccountQueue))
	}
	close(sender.CreateAccountQueue)
}

func TestAttestate_checkMainChainClaim(t *testing.T) {
	chains.StartEvmTestProvider(0, 0, true, nil, nil)
	chains.StartXrpTestProvider(0, 2, true, nil, nil)
	commit := &struct {
		Block       uint64
		ClaimId     uint64
		Sender      string
		Amount      string
		Destination string
		Nonce       int
		Fee         int
		BridgeId    string
	}{ClaimId: 0}

	_, err := checkMainChainClaim(commit)
	if err == nil {
		t.Errorf("error: should return error invalid claimId got %+v", err)
	}

	commit.ClaimId = 3
	claimExists, err := checkMainChainClaim(commit)
	if err != nil {
		t.Errorf("error: should not return error expected %+v got %+v", nil, err)
	}
	if claimExists != false {
		t.Errorf("error: expected %+v got %+v", false, claimExists)
	}

	commit.ClaimId = 1
	commit.Sender = "0xc2cD370bAdC28A01682394E8072824c1D7300D96"
	claimExists, err = checkMainChainClaim(commit)
	if err != nil {
		t.Errorf("error: should not return error expected %+v got %+v", nil, err)
	}
	if claimExists != false {
		t.Errorf("error: expected %+v got %+v", false, claimExists)
	}

	commit.Sender = "0x4DBeE27B94c970B6A7916628236ad6D9369a4518"
	claimExists, err = checkMainChainClaim(commit)
	if err != nil {
		t.Errorf("error: should not return error expected %+v got %+v", nil, err)
	}
	if claimExists != true {
		t.Errorf("error: expected %+v got %+v", true, claimExists)
	}
}

func TestAttestate_checkSideChainClaim(t *testing.T) {
	chains.StartEvmTestProvider(0, 0, true, nil, nil)
	chains.StartXrpTestProvider(0, 2, true, nil, nil)
	commit := &struct {
		Block       uint64
		ClaimId     uint64
		Sender      string
		Amount      string
		Destination string
		Nonce       int
		Fee         int
		BridgeId    string
	}{ClaimId: 0}

	_, err := checkSideChainClaim(commit)
	if err == nil {
		t.Errorf("error: should return error invalid claimId got %+v", err)
	}

	commit.ClaimId = 3
	claimExists, err := checkSideChainClaim(commit)
	if err != nil {
		t.Errorf("error: should not return error expected %+v got %+v", nil, err)
	}
	if claimExists != false {
		t.Errorf("error: expected %+v got %+v", false, claimExists)
	}

	commit.ClaimId = 2
	commit.Sender = "r4Sja8D6WM8hM53XksctnaYto6vhmh8dB"
	claimExists, err = checkSideChainClaim(commit)
	if err != nil {
		t.Errorf("error: should not return error expected %+v got %+v", nil, err)
	}
	if claimExists != false {
		t.Errorf("error: expected %+v got %+v", false, claimExists)
	}

	commit.Sender = "r3an6Cz2MgHQT9q3Kj3QzwQv9ARkgkkxqo"
	claimExists, err = checkSideChainClaim(commit)
	if err != nil {
		t.Errorf("error: should not return error expected %+v got %+v", nil, err)
	}
	if claimExists != true {
		t.Errorf("error: expected %+v got %+v", true, claimExists)
	}
}

func TestAttestate_checkMainChainCreateAccount(t *testing.T) {
	chains.StartXrpTestProvider(0, 2, true, nil, nil)
	account := "error-account"
	bridgeId := "XRP:XRP"

	_, err := checkMainChainCreateAccount(account, bridgeId)
	if err == nil {
		t.Errorf("error: should return error got %+v", err)
	}

	account = "found-account"
	canCreateAccount, err := checkMainChainCreateAccount(account, bridgeId)
	if err != nil {
		t.Errorf("error: should not return error expected %+v got %+v", nil, err)
	}
	if canCreateAccount != false {
		t.Errorf("error: expected %+v got %+v", false, canCreateAccount)
	}
	if chains.XrpTestProvider.CheckWitnessAttestedCalledTimes != 0 {
		t.Errorf("error: expected %+v got %+v", 0, chains.XrpTestProvider.CheckWitnessAttestedCalledTimes)
	}

	account = "error-attested"
	canCreateAccount, err = checkMainChainCreateAccount(account, bridgeId)
	if err == nil {
		t.Errorf("error: should return error got %+v", err)
	}
	if canCreateAccount != false {
		t.Errorf("error: expected %+v got %+v", false, canCreateAccount)
	}
	if chains.XrpTestProvider.CheckWitnessAttestedCalledTimes != 1 {
		t.Errorf("error: expected %+v got %+v", 1, chains.XrpTestProvider.CheckWitnessAttestedCalledTimes)
	}

	account = "found-attested"
	canCreateAccount, err = checkMainChainCreateAccount(account, bridgeId)
	if err != nil {
		t.Errorf("error: should not return error got %+v", err)
	}
	if canCreateAccount != false {
		t.Errorf("error: expected %+v got %+v", false, canCreateAccount)
	}
	if chains.XrpTestProvider.CheckWitnessAttestedCalledTimes != 2 {
		t.Errorf("error: expected %+v got %+v", 2, chains.XrpTestProvider.CheckWitnessAttestedCalledTimes)
	}

	account = "r4Sja8D6WM8hM53XksctnaYto6vhmh8dB"
	canCreateAccount, err = checkMainChainCreateAccount(account, bridgeId)
	if err != nil {
		t.Errorf("error: should not return error got %+v", err)
	}
	if canCreateAccount != true {
		t.Errorf("error: expected %+v got %+v", true, canCreateAccount)
	}
	if chains.XrpTestProvider.CheckWitnessAttestedCalledTimes != 3 {
		t.Errorf("error: expected %+v got %+v", 3, chains.XrpTestProvider.CheckWitnessAttestedCalledTimes)
	}
}

func TestAttestate_checkSideChainCreateAccount(t *testing.T) {
	chains.StartEvmTestProvider(0, 0, true, nil, nil)
	account := "error-account"
	bridgeId := "XRP:XRP"

	_, err := checkSideChainCreateAccount(account, bridgeId)
	if err == nil {
		t.Errorf("error: should return error got %+v", err)
	}

	account = "found-account"
	canCreateAccount, err := checkSideChainCreateAccount(account, bridgeId)
	if err != nil {
		t.Errorf("error: should not return error expected %+v got %+v", nil, err)
	}
	if canCreateAccount != false {
		t.Errorf("error: expected %+v got %+v", false, canCreateAccount)
	}
	if chains.EvmTestProvider.CheckWitnessAttestedCalledTimes != 0 {
		t.Errorf("error: expected %+v got %+v", 0, chains.EvmTestProvider.CheckWitnessAttestedCalledTimes)
	}

	account = "error-attested"
	canCreateAccount, err = checkSideChainCreateAccount(account, bridgeId)
	if err == nil {
		t.Errorf("error: should return error got %+v", err)
	}
	if canCreateAccount != false {
		t.Errorf("error: expected %+v got %+v", false, canCreateAccount)
	}
	if chains.EvmTestProvider.CheckWitnessAttestedCalledTimes != 1 {
		t.Errorf("error: expected %+v got %+v", 1, chains.EvmTestProvider.CheckWitnessAttestedCalledTimes)
	}

	account = "found-attested"
	canCreateAccount, err = checkSideChainCreateAccount(account, bridgeId)
	if err != nil {
		t.Errorf("error: should not return error got %+v", err)
	}
	if canCreateAccount != false {
		t.Errorf("error: expected %+v got %+v", false, canCreateAccount)
	}
	if chains.EvmTestProvider.CheckWitnessAttestedCalledTimes != 2 {
		t.Errorf("error: expected %+v got %+v", 2, chains.EvmTestProvider.CheckWitnessAttestedCalledTimes)
	}

	account = "0x4DBeE27B94c970B6A7916628236ad6D9369a4518"
	canCreateAccount, err = checkSideChainCreateAccount(account, bridgeId)
	if err != nil {
		t.Errorf("error: should not return error got %+v", err)
	}
	if canCreateAccount != true {
		t.Errorf("error: expected %+v got %+v", true, canCreateAccount)
	}
	if chains.EvmTestProvider.CheckWitnessAttestedCalledTimes != 3 {
		t.Errorf("error: expected %+v got %+v", 3, chains.EvmTestProvider.CheckWitnessAttestedCalledTimes)
	}
}
