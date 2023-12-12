package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type ChainType string

const (
	Xrp ChainType = "xrp"
	Evm ChainType = "evm"
)

type Server struct {
	QueuePeriod               int    `yaml:"queue_period"`
	BridgeListenerQueuePeriod int    `yaml:"bridge_listener_queue_period"`
	BridgeCreationQueuePeriod int    `yaml:"bridge_creation_queue_period"`
	LoggingLevel              string `yaml:"logging_level"`
	LogFilePath               string `yaml:"log_file_path"`
	LogFormat                 string `yaml:"log_format"`
	DynamicBridgeCreation     bool   `yaml:"dynamic_bridge_creation"`
	SequencerUrl              string `yaml:"sequencer_url"`
	MinBridgeSignatureReward  uint64 `yaml:"min_bridge_signature_reward"`
	MaxBridgeSignatureReward  uint64 `yaml:"max_bridge_signature_reward"`
	MaxCreateBridgeIterations uint64 `yaml:"max_create_bridge_iterations"`
}

type ChainConfig struct {
	Type              ChainType `yaml:"type"`
	Node              string    `yaml:"node"`
	DoorAddress       string    `yaml:"door_address"`
	StartingBlock     uint64    `yaml:"starting_block"`
	SignerListSeconds int64     `yaml:"signer_list_seconds"`
	MaxGasFactor      int64     `yaml:"max_gas_factor"`
	Signer            *Signer   `yaml:"signer"`
}

type Config struct {
	Server    `yaml:"server"`
	MainChain ChainConfig `yaml:"mainchain"`
	SideChain ChainConfig `yaml:"sidechain"`
}

func LoadConfig(filePath string) Config {
	var cfg Config
	readFile(&cfg, filePath)
	loadDotEnvVars()
	readEnv(&cfg)
	return cfg
}

func readFile(cfg *Config, filePath string) {
	if filePath == "" {
		filePath = "./configs/config.yml"
	}
	f, err := os.Open(filePath)
	if err != nil {
		log.Warn().Msg(err.Error())
		return
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func loadDotEnvVars() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Info().Msg(err.Error())
	}
}

func readEnv(cfg *Config) {
	logFormat := os.Getenv("SERVER_LOG_FORMAT")
	if logFormat != "" {
		cfg.Server.LogFormat = logFormat
	}

	serverQueuePeriod := os.Getenv("SERVER_QUEUE_PERIOD")
	if serverQueuePeriod != "" {
		period, err := strconv.Atoi(serverQueuePeriod)
		if err != nil {
			cfg.Server.QueuePeriod = period
		}
	}

	serverBridgeListenerQueuePeriod := os.Getenv("SERVER_BRIDGE_LISTENER_QUEUE_PERIOD")
	if serverBridgeListenerQueuePeriod != "" {
		period, err := strconv.Atoi(serverBridgeListenerQueuePeriod)
		if err != nil {
			cfg.Server.BridgeListenerQueuePeriod = period
		}
	}

	serverBridgeCreationQueuePeriod := os.Getenv("SERVER_BRIDGE_QUEUE_PERIOD")
	if serverBridgeCreationQueuePeriod != "" {
		period, err := strconv.Atoi(serverBridgeCreationQueuePeriod)
		if err != nil {
			cfg.Server.BridgeCreationQueuePeriod = period
		}
	}

	serverLoggingLevel := os.Getenv("SERVER_LOGGING_LEVEL")
	if serverLoggingLevel != "" {
		cfg.Server.LoggingLevel = serverLoggingLevel
	}

	serverLogFilePath := os.Getenv("SERVER_LOG_FILE_PATH")
	if serverLogFilePath != "" {
		cfg.Server.LogFilePath = serverLogFilePath
	}

	serverDynamicBridgeCreation := os.Getenv("SERVER_DYNAMIC_BRIDGE_CREATION")
	if serverDynamicBridgeCreation != "" {
		cfg.Server.DynamicBridgeCreation = !(serverDynamicBridgeCreation == "false")
	}

	serverSequencerUrl := os.Getenv("SERVER_SEQUENCER_URL")
	if serverSequencerUrl != "" {
		cfg.Server.SequencerUrl = serverSequencerUrl
	}

	serverMinBridgeSignatureReward := os.Getenv("SERVER_MIN_BRIDGE_SIGNATURE_REWARD")
	if serverMinBridgeSignatureReward != "" {
		minReward, err := strconv.Atoi(serverMinBridgeSignatureReward)
		if err != nil {
			cfg.Server.MinBridgeSignatureReward = uint64(minReward)
		}
	}

	serverMaxBridgeSignatureReward := os.Getenv("SERVER_MAX_BRIDGE_SIGNATURE_REWARD")
	if serverMaxBridgeSignatureReward != "" {
		maxReward, err := strconv.Atoi(serverMaxBridgeSignatureReward)
		if err != nil {
			cfg.Server.MaxBridgeSignatureReward = uint64(maxReward)
		}
	}

	serverMaxCreateBridgeIterations := os.Getenv("SERVER_MAX_CREATE_BRIDGE_ITERATIONS")
	if serverMaxCreateBridgeIterations != "" {
		maxIters, err := strconv.Atoi(serverMaxCreateBridgeIterations)
		if err != nil {
			cfg.Server.MaxCreateBridgeIterations = uint64(maxIters)
		}
	}

	mainchainType := os.Getenv("MAINCHAIN_TYPE")
	if mainchainType != "" {
		if mainchainType == "xrp" {
			cfg.MainChain.Type = Xrp
		} else if mainchainType == "evm" {
			cfg.MainChain.Type = Evm
		}
	}

	mainchainNode := os.Getenv("MAINCHAIN_NODE")
	if mainchainNode != "" {
		cfg.MainChain.Node = mainchainNode
	}

	mainchainDoorAddress := os.Getenv("MAINCHAIN_BRIDGE_ADDRESS")
	if mainchainDoorAddress != "" {
		cfg.MainChain.DoorAddress = mainchainDoorAddress
	}

	mainchainStartingBlock := os.Getenv("MAINCHAIN_STARTING_BLOCK")
	if mainchainStartingBlock != "" {
		block, err := strconv.Atoi(mainchainStartingBlock)
		if err != nil {
			cfg.MainChain.StartingBlock = uint64(block)
		}
	}

	mainchainSignerListSeconds := os.Getenv("MAINCHAIN_SIGNER_LIST_SECONDS")
	if mainchainSignerListSeconds != "" {
		seconds, err := strconv.Atoi(mainchainSignerListSeconds)
		if err != nil {
			cfg.MainChain.SignerListSeconds = int64(seconds)
		}
	}

	mainchainMaxGasFactor := os.Getenv("MAINCHAIN_MAX_GAS_FACTOR")
	if mainchainMaxGasFactor != "" {
		seconds, err := strconv.Atoi(mainchainMaxGasFactor)
		if err != nil {
			cfg.MainChain.MaxGasFactor = int64(seconds)
		}
	}

	sidechainType := os.Getenv("SIDECHAIN_TYPE")
	if sidechainType != "" {
		if sidechainType == "xrp" {
			cfg.SideChain.Type = Xrp
		} else if sidechainType == "evm" {
			cfg.SideChain.Type = Evm
		}
	}

	sidechainNode := os.Getenv("SIDECHAIN_NODE")
	if sidechainNode != "" {
		cfg.SideChain.Node = sidechainNode
	}

	sidechainDoorAddress := os.Getenv("SIDECHAIN_BRIDGE_ADDRESS")
	if sidechainDoorAddress != "" {
		cfg.SideChain.DoorAddress = sidechainDoorAddress
	}

	sidechainStartingBlock := os.Getenv("SIDECHAIN_STARTING_BLOCK")
	if sidechainStartingBlock != "" {
		block, err := strconv.Atoi(sidechainStartingBlock)
		if err != nil {
			cfg.SideChain.StartingBlock = uint64(block)
		}
	}

	sidechainSignerListSeconds := os.Getenv("SIDECHAIN_SIGNER_LIST_SECONDS")
	if sidechainSignerListSeconds != "" {
		seconds, err := strconv.Atoi(sidechainSignerListSeconds)
		if err != nil {
			cfg.SideChain.SignerListSeconds = int64(seconds)
		}
	}

	sidechainMaxGasFactor := os.Getenv("SIDECHAIN_MAX_GAS_FACTOR")
	if sidechainMaxGasFactor != "" {
		seconds, err := strconv.Atoi(sidechainMaxGasFactor)
		if err != nil {
			cfg.SideChain.MaxGasFactor = int64(seconds)
		}
	}

	readSignerEnv(cfg)
}
