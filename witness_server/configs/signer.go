package config

import (
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Signer struct {
	Type string      `yaml:"type"`
	Spec interface{} `yaml:"-"`
}

type AwsSigner struct {
	Region    string `yaml:"region"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	KeyId     string `yaml:"key_id"`
}

type LocalSigner struct {
	PrivateKey string `yaml:"private_key"`
}

func (s *Signer) UnmarshalYAML(node *yaml.Node) error {
	type S Signer
	type T struct {
		*S   `yaml:",inline"`
		Spec yaml.Node `yaml:"spec"`
	}

	obj := &T{S: (*S)(s)}
	if err := node.Decode(obj); err != nil {
		return err
	}

	switch s.Type {
	case "aws":
		s.Spec = new(AwsSigner)
	case "local":
		s.Spec = new(LocalSigner)
	default:
		log.Fatal().Msgf("Unknown signer type found in config %v", s.Type)
	}

	return obj.Spec.Decode(s.Spec)
}

func readSignerEnv(cfg *Config) {
	mainchainSignerType := ""
	if cfg.MainChain.Signer != nil {
		mainchainSignerType = cfg.MainChain.Signer.Type
	}
	envMainchainType := os.Getenv("MAINCHAIN_SIGNER_TYPE")
	if len(envMainchainType) != 0 {
		mainchainSignerType = envMainchainType
	}
	if mainchainSignerType != "" {
		switch mainchainSignerType {
		case "aws":
			cfg.MainChain.Signer.Spec = &AwsSigner{
				Region:    os.Getenv("MAINCHAIN_SIGNER_AWS_REGION"),
				AccessKey: os.Getenv("AWS_ACCESS_KEY"),
				SecretKey: os.Getenv("AWS_SECRET_KEY"),
				KeyId:     os.Getenv("MAINCHAIN_SIGNER_KMS_KEY_ID"),
			}
		case "local":
			if len(os.Getenv("MAINCHAIN_SIGNER_PRIVATE_KEY")) == 64 {
				cfg.MainChain.Signer.Spec = &LocalSigner{
					PrivateKey: os.Getenv("MAINCHAIN_SIGNER_PRIVATE_KEY"),
				}
			}
		default:
			log.Fatal().Msgf("Unknown MAINCHAIN signer type found in environment %v", mainchainSignerType)
		}
		cfg.MainChain.Signer.Type = mainchainSignerType
	}

	sidechainSignerType := os.Getenv("SIDECHAIN_SIGNER_TYPE")
	if sidechainSignerType != "" {
		switch sidechainSignerType {
		case "aws":
			cfg.SideChain.Signer.Spec = &AwsSigner{
				Region:    os.Getenv("SIDECHAIN_SIGNER_AWS_REGION"),
				AccessKey: os.Getenv("AWS_ACCESS_KEY"),
				SecretKey: os.Getenv("AWS_SECRET_KEY"),
				KeyId:     os.Getenv("SIDECHAIN_SIGNER_KMS_KEY_ID"),
			}
		case "local":
			if len(os.Getenv("SIDECHAIN_SIGNER_PRIVATE_KEY")) == 64 {
				cfg.SideChain.Signer.Spec = &LocalSigner{
					PrivateKey: os.Getenv("SIDECHAIN_SIGNER_PRIVATE_KEY"),
				}
			}
		default:
			log.Fatal().Msgf("Unknown SIDECHAIN signer type found in environment %v", sidechainSignerType)
		}
		cfg.SideChain.Signer.Type = sidechainSignerType
	}
}
