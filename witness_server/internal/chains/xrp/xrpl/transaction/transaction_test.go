package transaction

import (
	"github.com/stretchr/testify/require"
	"testing"
)

var addClaimAttestation = `{"Account":"rQGde68DA5avZXsUbjwjzyxQSJpGb4yREw","Amount":"10000000","AttestationRewardAccount":"rQGde68DA5avZXsUbjwjzyxQSJpGb4yREw","AttestationSignerAccount":"rQGde68DA5avZXsUbjwjzyxQSJpGb4yREw","Destination":"rs7Cxh3HytVTdrfkQRvw4vJ7eHbWrdgMoy","Fee":"10","Flags":0,"LastLedgerSequence":1207969,"OtherChainSource":"rhvkhs5wqJm1fQRMYn2s9nFcJSeYsiKmRN","PublicKey":"03EF6C5C1CFF617E08107FC723D376480A9C7CCDCCFAEEDE7424E015B3A640DC88","Sequence":1038537,"Signature":"3044022016670B2CE7C05137FBD8BF443E9400F0917B1C4A9A95A827B68607D3201C30F202205B672001D28EA1ACF525A84D9C1C110E6AD1037C9BF386C9E6FF2429C9B4BD90","SigningPubKey":"03EF6C5C1CFF617E08107FC723D376480A9C7CCDCCFAEEDE7424E015B3A640DC88","TransactionType":"XChainAddClaimAttestation","TxnSignature":"3045022100E8A965C77E9ABDD62A66D9AD837396BDE34243A89210789E7813BFA6CCFF36A502201A6E92AFF873F88D567ACB79E375423BBF17333947BC80CFE6F912A73E7BA228","WasLockingChainSend":0,"XChainBridge":{"IssuingChainDoor":"rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh","IssuingChainIssue":{"currency":"XRP"},"LockingChainDoor":"rayv9pKSvSuWaEU5gJiQRqsLXP5XBV1n5Y","LockingChainIssue":{"currency":"XRP"}},"XChainClaimID":"6"}`
var addAccountCreateAttestation = `{"Account":"rM8co4v5iExhFECoCD8KqPqPRYNwgXpT61","Amount":"9000000000","Destination":"rhqAvtMjKmKHyFDSu6e39fBDWQ2gmgVpQh","Fee":"12","Flags":0,"LastLedgerSequence":1208140,"Sequence":1208099,"SignatureReward":"1000000","SigningPubKey":"ED33D21D567934E2D3721B118CF39BF65C7CFC08A444E1EAFA3BC84BD72282E071","TransactionType":"XChainAccountCreateCommit","TxnSignature":"5356DB7711463BB13BE58E79788AFF215F8BA76FFFEE2CFD57D38F9BA10E70B24F6CB0C8209097AF55F97E998D0CBB0C7473AF9F9F06035A073E8A008C8EB00C","XChainBridge":{"IssuingChainDoor":"rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh","IssuingChainIssue":{"currency":"XRP"},"LockingChainDoor":"rayv9pKSvSuWaEU5gJiQRqsLXP5XBV1n5Y","LockingChainIssue":{"currency":"XRP"}}}`
var payment = `{"Account":"rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh","Amount":"1000000000000","Destination":"rpSspP5yYyomcSrgsohyKMCnu5oJsTMkYP","Fee":"12","Flags":2147483648,"LastLedgerSequence":424508,"Memos":[],"Sequence":300,"SigningPubKey":"0330E7FC9D56BB25D6893BA3F317AE5BCF33B3291BD63DB32654A313222F7FD020","TransactionType":"Payment","TxnSignature":"3045022100BDD838D8DA64E8AD433F1F8193060F945F63522710E12F2A7315877BDF86E4D502202618C47F38C87E1845E5649F18B6F990BE8E2316E1057620F9C4F0AD8F8CBDF9","date":725037291,"hash":"8C40D760AA032FD938AFF5BF70803A665626ABBD55CECF98D74FBE7470B93470","inLedger":424505,"ledger_index":424505}`
var createClaim = `{"Account":"rPEhYYszi4sdqUJVhQUUGuxh7dVyE19V9B","Fee":"12","Flags":0,"LastLedgerSequence":401357,"OtherChainSource":"rKf5XDdWGaSS5k66JvAkqrZHikBANqvTkN","Sequence":401334,"SignatureReward":"1000","SigningPubKey":"ED29265BB86D900CDED019DC926C2C3E03742A5CE5066447F991A5B57105EB85A4","TransactionType":"XChainCreateClaimID","TxnSignature":"C5C0CF5F10C22C9671F36F41A647930DC430EF847165F326C3836D550CEECDD7F8C6A8F702A3163EB67CA9E63187D3D423EF6D9D301DA23C6DDFDF37D2896A00","XChainBridge":{"IssuingChainDoor":"rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh","IssuingChainIssue":{"currency":"XRP"},"LockingChainDoor":"rhaY1Jxh8wiezQrRNrDdnfpVMKZJZd4ipt","LockingChainIssue":{"currency":"XRP"}},"date":724967072,"hash":"4DF18B10A8BA62ACAD1D4E8F54EBFA21C118A38972799C6E419D1D8AEEF8BF9A","inLedger":401339,"ledger_index":401339}`
var createBridge = `{"Account":"rhaY1Jxh8wiezQrRNrDdnfpVMKZJZd4ipt","Fee":"12","Flags":0,"LastLedgerSequence":401091,"MinAccountCreateAmount":"1","Sequence":401072,"SignatureReward":"1000","SigningPubKey":"EDB1BB1A1C3EF9244C4426DB726CC2EC03AC94531EC6CA05197E6F1A2EC82E9B4D","TransactionType":"XChainCreateBridge","TxnSignature":"9B91D336800CD48C3D1E31438D08AD0D61A74377C7CC6BB8B501006AD8AD3651940356A45AEFF48830F4B670DD89E2699DDBB76C2A3B01DED400AB0F0586C001","XChainBridge":{"IssuingChainDoor":"rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh","IssuingChainIssue":{"currency":"XRP"},"LockingChainDoor":"rhaY1Jxh8wiezQrRNrDdnfpVMKZJZd4ipt","LockingChainIssue":{"currency":"XRP"}},"date":724966273,"hash":"B927193BD05FBAD100523A60D0E82A218E44FF37E55FDF3CA02B95ADCD9C1E61","inLedger":401073,"ledger_index":401073}`
var signerListSet = `{"Account":"rhaY1Jxh8wiezQrRNrDdnfpVMKZJZd4ipt","Fee":"12","Flags":0,"LastLedgerSequence":401091,"Sequence":401071,"SignerEntries":[{"SignerEntry":{"Account":"rpSspP5yYyomcSrgsohyKMCnu5oJsTMkYP","SignerWeight":1}},{"SignerEntry":{"Account":"rUFDiADdSDbgbvYubaDeYnoqLyVjfrnjLB","SignerWeight":1}}],"SignerQuorum":1,"SigningPubKey":"EDB1BB1A1C3EF9244C4426DB726CC2EC03AC94531EC6CA05197E6F1A2EC82E9B4D","TransactionType":"SignerListSet","TxnSignature":"0AF8AC3AA1A34C08D322B4830519A2DC553E9012292E96650E8E7AFBA3CDCAA88977506E9F1B90B9FD0AE537059B19BE4DB3B17093736F57EBC3CBA8EF93CE0E","date":724966273,"hash":"FB8342AF66A936A40D7AB39EA1B85807A8232915C28757489A4D868CCEA08C78","inLedger":401073,"ledger_index":401073}`
var accountCreateCommit = `{"Account":"rKzspyP7z9qEuak2YVnNn7TCinnWnpxFma","Amount":"500000000","Destination":"rwrRS1UYjVi3pNLg7QqtuzMoaTtueJ2Gim","Fee":"12","Flags":0,"LastLedgerSequence":426397,"Sequence":426377,"SignatureReward":"1000000","SigningPubKey":"ED4E5282AAEA179F2261E73AC4D220365F00933AE95C89FB4D14B224E6CE17D544","TransactionType":"XChainAccountCreateCommit","TxnSignature":"A17B173FD8785EF9395FE166B2EC15734FF768B0E481D8E2554FC0D3E29312050CBE2C5E80A532CBDFD6D5BF2F47E6814B2ABEB4C69B76D2FB7B5EEE84EC3006","XChainBridge":{"IssuingChainDoor":"rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh","IssuingChainIssue":{"currency":"XRP"},"LockingChainDoor":"rapLiFbSsEhWszvgFViv9aB4LXzGaHqFd8","LockingChainIssue":{"currency":"XRP"}},"date":725042971,"hash":"7239148696836E8B7024D847D54604862FCE5814601AE97A8506A4CECCE89C59","inLedger":426379,"ledger_index":426379}`
var commit = `{"Account":"rKpteb8hRJtFWWxgZozoyrFxTM36W8uiWy","Amount":"50000000","Fee":"12","Flags":0,"LastLedgerSequence":424413,"OtherChainDestination":"rh2yN6Epe5nDt1BA62whQr355rBu9ZGnRb","Sequence":424390,"SigningPubKey":"EDDEECD806375495D33EF2FFE1B3494DE3E2D74AFDF3B2BE8BADD5F1A40C60F4BF","TransactionType":"XChainCommit","TxnSignature":"6C408C09D2F75684AC46C02F00E80F486D3D894310FBE417947402D081437B786E71F9279933588929C7DC762BC1C2EFEB16EF63DF45432CA70BA5870D769607","XChainBridge":{"IssuingChainDoor":"rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh","IssuingChainIssue":{"currency":"XRP"},"LockingChainDoor":"rapLiFbSsEhWszvgFViv9aB4LXzGaHqFd8","LockingChainIssue":{"currency":"XRP"}},"XChainClaimID":"1","date":725036961,"hash":"4DFAC10F20F19B85AA78777BA26F22F4A72BEB31025BFE0522C24BE9C311CD9C","inLedger":424395,"ledger_index":424395}`

func Test_UnmarshalClaimAttestationTransaction(t *testing.T) {
	tx, err := UnmarshalTransaction(addClaimAttestation)
	if err != nil {
		t.Errorf("Error unmarshaling %v", err)
	}

	marshaled, err := MarshalTransaction(tx)
	if err != nil {
		t.Errorf("Error marshaling %v", err)
	}

	require.JSONEq(t, marshaled, addClaimAttestation)
}

func Test_UnmarshalAccountCreateAttestationTransaction(t *testing.T) {
	tx, err := UnmarshalTransaction(addAccountCreateAttestation)
	if err != nil {
		t.Errorf("Error unmarshaling %v", err)
	}

	marshaled, err := MarshalTransaction(tx)
	if err != nil {
		t.Errorf("Error marshaling %v", err)
	}

	require.JSONEq(t, marshaled, addAccountCreateAttestation)
}

func Test_UnmarshalPaymentTransaction(t *testing.T) {
	tx, err := UnmarshalTransaction(payment)
	if err != nil {
		t.Errorf("Error unmarshaling %v", err)
	}

	marshaled, err := MarshalTransaction(tx)
	if err != nil {
		t.Errorf("Error marshaling %v", err)
	}
	require.JSONEq(t, marshaled, payment)
}

func Test_UnmarshalCreateClaimTransaction(t *testing.T) {
	tx, err := UnmarshalTransaction(createClaim)
	if err != nil {
		t.Errorf("Error unmarshaling %v", err)
	}

	marshaled, err := MarshalTransaction(tx)
	if err != nil {
		t.Errorf("Error marshaling %v", err)
	}

	require.JSONEq(t, marshaled, createClaim)
}

func Test_UnmarshalCreateBridgeTransaction(t *testing.T) {
	tx, err := UnmarshalTransaction(createBridge)
	if err != nil {
		t.Errorf("Error unmarshaling %v", err)
	}

	marshaled, err := MarshalTransaction(tx)
	if err != nil {
		t.Errorf("Error marshaling %v", err)
	}
	require.JSONEq(t, marshaled, createBridge)
}

func Test_UnmarshalSignerListSetTransaction(t *testing.T) {
	tx, err := UnmarshalTransaction(signerListSet)
	if err != nil {
		t.Errorf("Error unmarshaling %v", err)
	}

	marshaled, err := MarshalTransaction(tx)
	if err != nil {
		t.Errorf("Error marshaling %v", err)
	}
	require.JSONEq(t, marshaled, signerListSet)
}

func Test_UnmarshalAccountCreateCommitTransaction(t *testing.T) {
	tx, err := UnmarshalTransaction(accountCreateCommit)
	if err != nil {
		t.Errorf("Error unmarshaling %v", err)
	}

	marshaled, err := MarshalTransaction(tx)
	if err != nil {
		t.Errorf("Error marshaling %v", err)
	}
	require.JSONEq(t, marshaled, accountCreateCommit)
}

func Test_UnmarshalCommitTransaction(t *testing.T) {
	tx, err := UnmarshalTransaction(commit)
	if err != nil {
		t.Errorf("Error unmarshaling %v", err)
	}

	marshaled, err := MarshalTransaction(tx)
	if err != nil {
		t.Errorf("Error marshaling %v", err)
	}
	require.JSONEq(t, marshaled, commit)
}
