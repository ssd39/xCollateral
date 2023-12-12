package factory

import (
	config "peersyst/bridge-witness-go/configs"
	"peersyst/bridge-witness-go/internal/signer"
	awsKms "peersyst/bridge-witness-go/internal/signer/aws_kms"
	evmAwsKms "peersyst/bridge-witness-go/internal/signer/aws_kms/evm"
	xrpAwsKms "peersyst/bridge-witness-go/internal/signer/aws_kms/xrp"
	evmLocal "peersyst/bridge-witness-go/internal/signer/local/evm"
	xrpLocal "peersyst/bridge-witness-go/internal/signer/local/xrp"

	"github.com/rs/zerolog/log"
)

func NewSignerProviderFromConfig(chainType config.ChainType, chainConfig config.ChainConfig) signer.SignerProvider {
	switch chainConfig.Signer.Type {
	case "aws":
		signerSpec, ok := chainConfig.Signer.Spec.(*config.AwsSigner)
		if !ok {
			log.Fatal().Msgf("Error instantiating aws signer spec for config %v", chainConfig)
		}
		kmsService := awsKms.NewAwsKmsService(*signerSpec)
		log.Info().Msgf("Kms service : '%+v'", kmsService)

		switch chainType {
		case config.Xrp:
			return xrpAwsKms.NewXrpAwsKmsSignerProvider(kmsService)
		case config.Evm:
			return evmAwsKms.NewEvmAwsKmsSignerProvider(kmsService)
		}

	case "local":
		signerSpec, ok := chainConfig.Signer.Spec.(*config.LocalSigner)
		if !ok {
			log.Fatal().Msgf("Error instantiating local signer spec for config %v", chainConfig.Signer.Spec)
		}
		switch chainType {
		case config.Evm:
			return evmLocal.NewEvmLocalSignerProvider(*signerSpec)
		case config.Xrp:
			return xrpLocal.NewXrpLocalSignerProvider(*signerSpec)
		}
	}
	log.Fatal().Msgf("Unknown signer for %v with config %v", chainType, chainConfig)
	return nil
}
