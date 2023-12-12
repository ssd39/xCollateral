package aws_kms

import (
	"context"
	"crypto/ecdsa"
	"encoding/asn1"
	"encoding/hex"
	"math/big"
	config "peersyst/bridge-witness-go/configs"
	"peersyst/bridge-witness-go/internal/common/cache"
	"time"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	kmsTypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
)

type AwsKmsService struct {
	region string
	keyId  string
	client *kms.Client
}

func (service *AwsKmsService) GetPublicKey() string {
	compressedPubKey := service.GetCompressedPublicKey()
	if compressedPubKey == nil {
		return ""
	}
	return hex.EncodeToString(compressedPubKey)
}

func (service *AwsKmsService) describeKey() (*kms.DescribeKeyOutput, error) {
	input := kms.DescribeKeyInput{KeyId: &service.keyId}
	return service.client.DescribeKey(context.Background(), &input)
}

func (service *AwsKmsService) Sign(message []byte) (r, s, signature []byte) {
	signInput := kms.SignInput{
		KeyId:            &service.keyId,
		Message:          message,
		SigningAlgorithm: kmsTypes.SigningAlgorithmSpecEcdsaSha256,
		MessageType:      kmsTypes.MessageTypeDigest,
	}

	signOutput, err := service.client.Sign(context.Background(), &signInput)
	if err != nil {
		log.Error().Msgf("Error signing message: '%s'", err)
		return nil, nil, nil
	}

	var parsedSig struct{ R, S *big.Int }
	if _, err := asn1.Unmarshal(signOutput.Signature, &parsedSig); err != nil {
		log.Error().Msgf("asn1.Unmarshal: %s", err)
		return nil, nil, nil
	}

	// left pad R and S with zeroes
	rBytes := parsedSig.R.Bytes()
	sBytes := parsedSig.S.Bytes()

	// To avoid replay attack with inverse signature, only one half of the curve is allowed as a valid signature
	secp256k1N := crypto.S256().Params().N
	secp256k1halfN := new(big.Int).Div(secp256k1N, big.NewInt(2))
	sBI := new(big.Int).SetBytes(sBytes)

	// sBI > secp256k1halfN
	if sBI.Cmp(secp256k1halfN) == 1 {
		sBytes = sBI.Sub(secp256k1N, sBI).Bytes()
	}

	return rBytes, sBytes, signOutput.Signature
}

func (service *AwsKmsService) GetUncompressedPublicKey() *ecdsa.PublicKey {
	var ecdsaPubKey ecdsa.PublicKey

	err := cache.GetAndSet(func() any {
		pubKeyInput := kms.GetPublicKeyInput{KeyId: &service.keyId}
		pubKey, err := service.client.GetPublicKey(context.Background(), &pubKeyInput)
		if err != nil {
			log.Error().Msgf("Error getting public key from kms: '%+v'", err)
			return nil
		}

		var asn1pubk EcdsaPublicKey
		_, err = asn1.Unmarshal(pubKey.PublicKey, &asn1pubk)
		if err != nil {
			log.Error().Msgf("Error decoding kms public key in der: '%+v'", err)
			return nil
		}

		ecdsaUnmarshalledPubKey, err := crypto.UnmarshalPubkey(asn1pubk.PublicKey.Bytes)
		if err != nil {
			log.Error().Msgf("Error unmarshal ecdsa public key: '%+v'", err)
			return nil
		}

		return *ecdsaUnmarshalledPubKey
	}, &ecdsaPubKey, time.Now().Add(time.Hour*24))
	if err != nil {
		log.Error().Msgf("Error getting cache value for uncompressed public key: '%+v'", err)
	}
	return &ecdsaPubKey
}

func (service *AwsKmsService) GetCompressedPublicKey() []byte {
	ecdsaPubKey := service.GetUncompressedPublicKey()
	if ecdsaPubKey == nil {
		return nil
	}
	return crypto.CompressPubkey(ecdsaPubKey)
}

func NewAwsKmsService(awsConfigObj config.AwsSigner) *AwsKmsService {
	cfg, _ := awsConfig.LoadDefaultConfig(context.TODO(), awsConfig.WithRegion(awsConfigObj.Region))
	client := kms.NewFromConfig(cfg)
	kmsService := &AwsKmsService{
		region: awsConfigObj.Region,
		keyId:  awsConfigObj.KeyId,
		client: client,
	}

	// Check if keyId exists
	_, err := kmsService.describeKey()
	if err != nil {
		log.Fatal().Msgf("Error with keyId: '%s'", err)
	}

	return kmsService
}
