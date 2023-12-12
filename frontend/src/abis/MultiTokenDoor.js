export default {
  abi: [
    {
      inputs: [
        {
          internalType: "contract GnosisSafeL2",
          name: "safe",
          type: "address",
        },
      ],
      stateMutability: "nonpayable",
      type: "constructor",
    },
    {
      anonymous: false,
      inputs: [
        {
          indexed: true,
          internalType: "bytes32",
          name: "bridgeKey",
          type: "bytes32",
        },
        {
          indexed: true,
          internalType: "uint256",
          name: "claimId",
          type: "uint256",
        },
        {
          indexed: true,
          internalType: "address",
          name: "witness",
          type: "address",
        },
        {
          indexed: false,
          internalType: "uint256",
          name: "value",
          type: "uint256",
        },
        {
          indexed: false,
          internalType: "address",
          name: "receiver",
          type: "address",
        },
      ],
      name: "AddClaimAttestation",
      type: "event",
    },
    {
      anonymous: false,
      inputs: [
        {
          indexed: true,
          internalType: "bytes32",
          name: "bridgeKey",
          type: "bytes32",
        },
        {
          indexed: true,
          internalType: "address",
          name: "witness",
          type: "address",
        },
        {
          indexed: true,
          internalType: "address",
          name: "receiver",
          type: "address",
        },
        {
          indexed: false,
          internalType: "uint256",
          name: "value",
          type: "uint256",
        },
      ],
      name: "AddCreateAccountAttestation",
      type: "event",
    },
    {
      anonymous: false,
      inputs: [
        {
          indexed: true,
          internalType: "bytes32",
          name: "bridgeKey",
          type: "bytes32",
        },
        {
          indexed: true,
          internalType: "uint256",
          name: "claimId",
          type: "uint256",
        },
        {
          indexed: true,
          internalType: "address",
          name: "sender",
          type: "address",
        },
        {
          indexed: false,
          internalType: "uint256",
          name: "value",
          type: "uint256",
        },
        {
          indexed: false,
          internalType: "address",
          name: "destination",
          type: "address",
        },
      ],
      name: "Claim",
      type: "event",
    },
    {
      anonymous: false,
      inputs: [
        {
          indexed: true,
          internalType: "bytes32",
          name: "bridgeKey",
          type: "bytes32",
        },
        {
          indexed: true,
          internalType: "uint256",
          name: "claimId",
          type: "uint256",
        },
        {
          indexed: true,
          internalType: "address",
          name: "sender",
          type: "address",
        },
        {
          indexed: false,
          internalType: "uint256",
          name: "value",
          type: "uint256",
        },
        {
          indexed: false,
          internalType: "address",
          name: "receiver",
          type: "address",
        },
      ],
      name: "Commit",
      type: "event",
    },
    {
      anonymous: false,
      inputs: [
        {
          indexed: true,
          internalType: "bytes32",
          name: "bridgeKey",
          type: "bytes32",
        },
        {
          indexed: true,
          internalType: "uint256",
          name: "claimId",
          type: "uint256",
        },
        {
          indexed: true,
          internalType: "address",
          name: "sender",
          type: "address",
        },
        {
          indexed: false,
          internalType: "uint256",
          name: "value",
          type: "uint256",
        },
      ],
      name: "CommitWithoutAddress",
      type: "event",
    },
    {
      anonymous: false,
      inputs: [
        {
          indexed: true,
          internalType: "address",
          name: "receiver",
          type: "address",
        },
        {
          indexed: false,
          internalType: "uint256",
          name: "value",
          type: "uint256",
        },
      ],
      name: "CreateAccount",
      type: "event",
    },
    {
      anonymous: false,
      inputs: [
        {
          indexed: true,
          internalType: "bytes32",
          name: "bridgeKey",
          type: "bytes32",
        },
        {
          indexed: true,
          internalType: "address",
          name: "creator",
          type: "address",
        },
        {
          indexed: true,
          internalType: "address",
          name: "destination",
          type: "address",
        },
        {
          indexed: false,
          internalType: "uint256",
          name: "value",
          type: "uint256",
        },
        {
          indexed: false,
          internalType: "uint256",
          name: "signatureReward",
          type: "uint256",
        },
      ],
      name: "CreateAccountCommit",
      type: "event",
    },
    {
      anonymous: false,
      inputs: [
        {
          indexed: true,
          internalType: "bytes32",
          name: "bridgeKey",
          type: "bytes32",
        },
        {
          indexed: false,
          internalType: "address",
          name: "lockingChainDoor",
          type: "address",
        },
        {
          indexed: false,
          internalType: "address",
          name: "lockingChainIssueIssuer",
          type: "address",
        },
        {
          indexed: false,
          internalType: "string",
          name: "lockingChainIssueCurency",
          type: "string",
        },
        {
          indexed: false,
          internalType: "address",
          name: "issuingChainDoor",
          type: "address",
        },
        {
          indexed: false,
          internalType: "address",
          name: "issuingChainIssueIssuer",
          type: "address",
        },
        {
          indexed: false,
          internalType: "string",
          name: "issuingChainIssueCurency",
          type: "string",
        },
      ],
      name: "CreateBridge",
      type: "event",
    },
    {
      anonymous: false,
      inputs: [
        {
          indexed: false,
          internalType: "address",
          name: "tokenAddress",
          type: "address",
        },
      ],
      name: "CreateBridgeRequest",
      type: "event",
    },
    {
      anonymous: false,
      inputs: [
        {
          indexed: true,
          internalType: "bytes32",
          name: "bridgeKey",
          type: "bytes32",
        },
        {
          indexed: true,
          internalType: "uint256",
          name: "claimId",
          type: "uint256",
        },
        {
          indexed: true,
          internalType: "address",
          name: "creator",
          type: "address",
        },
        {
          indexed: false,
          internalType: "address",
          name: "sender",
          type: "address",
        },
      ],
      name: "CreateClaim",
      type: "event",
    },
    {
      anonymous: false,
      inputs: [
        {
          indexed: true,
          internalType: "bytes32",
          name: "bridgeKey",
          type: "bytes32",
        },
        {
          indexed: true,
          internalType: "uint256",
          name: "claimId",
          type: "uint256",
        },
        {
          indexed: true,
          internalType: "address",
          name: "receiver",
          type: "address",
        },
        {
          indexed: false,
          internalType: "uint256",
          name: "value",
          type: "uint256",
        },
      ],
      name: "Credit",
      type: "event",
    },
    {
      anonymous: false,
      inputs: [
        {
          indexed: true,
          internalType: "address",
          name: "previousOwner",
          type: "address",
        },
        {
          indexed: true,
          internalType: "address",
          name: "newOwner",
          type: "address",
        },
      ],
      name: "OwnershipTransferred",
      type: "event",
    },
    {
      anonymous: false,
      inputs: [
        {
          indexed: false,
          internalType: "address",
          name: "account",
          type: "address",
        },
      ],
      name: "Paused",
      type: "event",
    },
    {
      anonymous: false,
      inputs: [
        {
          indexed: false,
          internalType: "address",
          name: "account",
          type: "address",
        },
      ],
      name: "Unpaused",
      type: "event",
    },
    {
      inputs: [],
      name: "MIN_CREATE_BRIDGE_REWARD",
      outputs: [
        {
          internalType: "uint256",
          name: "",
          type: "uint256",
        },
      ],
      stateMutability: "view",
      type: "function",
    },
    {
      inputs: [],
      name: "_safe",
      outputs: [
        {
          internalType: "contract GnosisSafeL2",
          name: "",
          type: "address",
        },
      ],
      stateMutability: "view",
      type: "function",
    },
    {
      inputs: [
        {
          components: [
            {
              internalType: "address",
              name: "lockingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "lockingChainIssue",
              type: "tuple",
            },
            {
              internalType: "address",
              name: "issuingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "issuingChainIssue",
              type: "tuple",
            },
          ],
          internalType: "struct XChainTypes.BridgeConfig",
          name: "bridgeConfig",
          type: "tuple",
        },
        {
          internalType: "uint256",
          name: "claimId",
          type: "uint256",
        },
        {
          internalType: "uint256",
          name: "amount",
          type: "uint256",
        },
        {
          internalType: "address",
          name: "sender",
          type: "address",
        },
        {
          internalType: "address",
          name: "destination",
          type: "address",
        },
      ],
      name: "addClaimAttestation",
      outputs: [],
      stateMutability: "nonpayable",
      type: "function",
    },
    {
      inputs: [
        {
          components: [
            {
              internalType: "address",
              name: "lockingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "lockingChainIssue",
              type: "tuple",
            },
            {
              internalType: "address",
              name: "issuingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "issuingChainIssue",
              type: "tuple",
            },
          ],
          internalType: "struct XChainTypes.BridgeConfig",
          name: "bridgeConfig",
          type: "tuple",
        },
        {
          internalType: "address",
          name: "destination",
          type: "address",
        },
        {
          internalType: "uint256",
          name: "amount",
          type: "uint256",
        },
        {
          internalType: "uint256",
          name: "signatureReward",
          type: "uint256",
        },
      ],
      name: "addCreateAccountAttestation",
      outputs: [],
      stateMutability: "nonpayable",
      type: "function",
    },
    {
      inputs: [
        {
          components: [
            {
              internalType: "address",
              name: "lockingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "lockingChainIssue",
              type: "tuple",
            },
            {
              internalType: "address",
              name: "issuingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "issuingChainIssue",
              type: "tuple",
            },
          ],
          internalType: "struct XChainTypes.BridgeConfig",
          name: "bridgeConfig",
          type: "tuple",
        },
        {
          internalType: "uint256",
          name: "claimId",
          type: "uint256",
        },
        {
          internalType: "uint256",
          name: "amount",
          type: "uint256",
        },
        {
          internalType: "address",
          name: "destination",
          type: "address",
        },
      ],
      name: "claim",
      outputs: [],
      stateMutability: "nonpayable",
      type: "function",
    },
    {
      inputs: [
        {
          components: [
            {
              internalType: "address",
              name: "lockingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "lockingChainIssue",
              type: "tuple",
            },
            {
              internalType: "address",
              name: "issuingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "issuingChainIssue",
              type: "tuple",
            },
          ],
          internalType: "struct XChainTypes.BridgeConfig",
          name: "bridgeConfig",
          type: "tuple",
        },
        {
          internalType: "address",
          name: "receiver",
          type: "address",
        },
        {
          internalType: "uint256",
          name: "claimId",
          type: "uint256",
        },
        {
          internalType: "uint256",
          name: "amount",
          type: "uint256",
        },
      ],
      name: "commit",
      outputs: [],
      stateMutability: "payable",
      type: "function",
    },
    {
      inputs: [
        {
          components: [
            {
              internalType: "address",
              name: "lockingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "lockingChainIssue",
              type: "tuple",
            },
            {
              internalType: "address",
              name: "issuingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "issuingChainIssue",
              type: "tuple",
            },
          ],
          internalType: "struct XChainTypes.BridgeConfig",
          name: "bridgeConfig",
          type: "tuple",
        },
        {
          internalType: "uint256",
          name: "claimId",
          type: "uint256",
        },
        {
          internalType: "uint256",
          name: "amount",
          type: "uint256",
        },
      ],
      name: "commitWithoutAddress",
      outputs: [],
      stateMutability: "payable",
      type: "function",
    },
    {
      inputs: [
        {
          components: [
            {
              internalType: "address",
              name: "lockingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "lockingChainIssue",
              type: "tuple",
            },
            {
              internalType: "address",
              name: "issuingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "issuingChainIssue",
              type: "tuple",
            },
          ],
          internalType: "struct XChainTypes.BridgeConfig",
          name: "bridgeConfig",
          type: "tuple",
        },
        {
          internalType: "address",
          name: "destination",
          type: "address",
        },
        {
          internalType: "uint256",
          name: "amount",
          type: "uint256",
        },
        {
          internalType: "uint256",
          name: "signatureReward",
          type: "uint256",
        },
      ],
      name: "createAccountCommit",
      outputs: [],
      stateMutability: "payable",
      type: "function",
    },
    {
      inputs: [
        {
          components: [
            {
              internalType: "address",
              name: "lockingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "lockingChainIssue",
              type: "tuple",
            },
            {
              internalType: "address",
              name: "issuingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "issuingChainIssue",
              type: "tuple",
            },
          ],
          internalType: "struct XChainTypes.BridgeConfig",
          name: "config",
          type: "tuple",
        },
        {
          components: [
            {
              internalType: "uint256",
              name: "minCreateAmount",
              type: "uint256",
            },
            {
              internalType: "uint256",
              name: "signatureReward",
              type: "uint256",
            },
          ],
          internalType: "struct XChainTypes.BridgeParams",
          name: "params",
          type: "tuple",
        },
      ],
      name: "createBridge",
      outputs: [],
      stateMutability: "nonpayable",
      type: "function",
    },
    {
      inputs: [
        {
          internalType: "address",
          name: "tokenAddress",
          type: "address",
        },
      ],
      name: "createBridgeRequest",
      outputs: [],
      stateMutability: "payable",
      type: "function",
    },
    {
      inputs: [
        {
          components: [
            {
              internalType: "address",
              name: "lockingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "lockingChainIssue",
              type: "tuple",
            },
            {
              internalType: "address",
              name: "issuingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "issuingChainIssue",
              type: "tuple",
            },
          ],
          internalType: "struct XChainTypes.BridgeConfig",
          name: "bridgeConfig",
          type: "tuple",
        },
        {
          internalType: "address",
          name: "sender",
          type: "address",
        },
      ],
      name: "createClaimId",
      outputs: [
        {
          internalType: "uint256",
          name: "",
          type: "uint256",
        },
      ],
      stateMutability: "payable",
      type: "function",
    },
    {
      inputs: [
        {
          internalType: "address",
          name: "to",
          type: "address",
        },
        {
          internalType: "uint256",
          name: "value",
          type: "uint256",
        },
        {
          internalType: "bytes",
          name: "data",
          type: "bytes",
        },
        {
          internalType: "enum Enum.Operation",
          name: "operation",
          type: "uint8",
        },
      ],
      name: "execute",
      outputs: [
        {
          internalType: "bool",
          name: "success",
          type: "bool",
        },
      ],
      stateMutability: "nonpayable",
      type: "function",
    },
    {
      inputs: [
        {
          components: [
            {
              internalType: "address",
              name: "lockingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "lockingChainIssue",
              type: "tuple",
            },
            {
              internalType: "address",
              name: "issuingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "issuingChainIssue",
              type: "tuple",
            },
          ],
          internalType: "struct XChainTypes.BridgeConfig",
          name: "bridgeConfig",
          type: "tuple",
        },
        {
          internalType: "uint256",
          name: "claimId",
          type: "uint256",
        },
      ],
      name: "getBridgeClaim",
      outputs: [
        {
          internalType: "address",
          name: "",
          type: "address",
        },
        {
          internalType: "address",
          name: "",
          type: "address",
        },
        {
          internalType: "bool",
          name: "",
          type: "bool",
        },
      ],
      stateMutability: "view",
      type: "function",
    },
    {
      inputs: [
        {
          components: [
            {
              internalType: "address",
              name: "lockingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "lockingChainIssue",
              type: "tuple",
            },
            {
              internalType: "address",
              name: "issuingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "issuingChainIssue",
              type: "tuple",
            },
          ],
          internalType: "struct XChainTypes.BridgeConfig",
          name: "bridgeConfig",
          type: "tuple",
        },
      ],
      name: "getBridgeConfig",
      outputs: [
        {
          internalType: "address",
          name: "",
          type: "address",
        },
        {
          internalType: "address",
          name: "",
          type: "address",
        },
        {
          internalType: "string",
          name: "",
          type: "string",
        },
        {
          internalType: "address",
          name: "",
          type: "address",
        },
        {
          internalType: "address",
          name: "",
          type: "address",
        },
        {
          internalType: "string",
          name: "",
          type: "string",
        },
      ],
      stateMutability: "view",
      type: "function",
    },
    {
      inputs: [
        {
          components: [
            {
              internalType: "address",
              name: "lockingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "lockingChainIssue",
              type: "tuple",
            },
            {
              internalType: "address",
              name: "issuingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "issuingChainIssue",
              type: "tuple",
            },
          ],
          internalType: "struct XChainTypes.BridgeConfig",
          name: "bridgeConfig",
          type: "tuple",
        },
        {
          internalType: "address",
          name: "account",
          type: "address",
        },
      ],
      name: "getBridgeCreateAccount",
      outputs: [
        {
          internalType: "uint256",
          name: "",
          type: "uint256",
        },
        {
          internalType: "bool",
          name: "",
          type: "bool",
        },
        {
          internalType: "bool",
          name: "",
          type: "bool",
        },
      ],
      stateMutability: "view",
      type: "function",
    },
    {
      inputs: [
        {
          components: [
            {
              internalType: "address",
              name: "lockingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "lockingChainIssue",
              type: "tuple",
            },
            {
              internalType: "address",
              name: "issuingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "issuingChainIssue",
              type: "tuple",
            },
          ],
          internalType: "struct XChainTypes.BridgeConfig",
          name: "bridgeConfig",
          type: "tuple",
        },
      ],
      name: "getBridgeKey",
      outputs: [
        {
          internalType: "bytes32",
          name: "",
          type: "bytes32",
        },
      ],
      stateMutability: "pure",
      type: "function",
    },
    {
      inputs: [
        {
          components: [
            {
              internalType: "address",
              name: "lockingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "lockingChainIssue",
              type: "tuple",
            },
            {
              internalType: "address",
              name: "issuingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "issuingChainIssue",
              type: "tuple",
            },
          ],
          internalType: "struct XChainTypes.BridgeConfig",
          name: "bridgeConfig",
          type: "tuple",
        },
      ],
      name: "getBridgeParams",
      outputs: [
        {
          internalType: "uint256",
          name: "",
          type: "uint256",
        },
        {
          internalType: "uint256",
          name: "",
          type: "uint256",
        },
      ],
      stateMutability: "view",
      type: "function",
    },
    {
      inputs: [
        {
          components: [
            {
              internalType: "address",
              name: "lockingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "lockingChainIssue",
              type: "tuple",
            },
            {
              internalType: "address",
              name: "issuingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "issuingChainIssue",
              type: "tuple",
            },
          ],
          internalType: "struct XChainTypes.BridgeConfig",
          name: "bridgeConfig",
          type: "tuple",
        },
      ],
      name: "getBridgeToken",
      outputs: [
        {
          internalType: "address",
          name: "",
          type: "address",
        },
      ],
      stateMutability: "view",
      type: "function",
    },
    {
      inputs: [
        {
          internalType: "uint256",
          name: "page",
          type: "uint256",
        },
      ],
      name: "getBridgesPaginated",
      outputs: [
        {
          components: [
            {
              internalType: "address",
              name: "lockingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "lockingChainIssue",
              type: "tuple",
            },
            {
              internalType: "address",
              name: "issuingChainDoor",
              type: "address",
            },
            {
              components: [
                {
                  internalType: "address",
                  name: "issuer",
                  type: "address",
                },
                {
                  internalType: "string",
                  name: "currency",
                  type: "string",
                },
              ],
              internalType: "struct XChainTypes.BridgeChainIssue",
              name: "issuingChainIssue",
              type: "tuple",
            },
          ],
          internalType: "struct XChainTypes.BridgeConfig[]",
          name: "configs",
          type: "tuple[]",
        },
        {
          components: [
            {
              internalType: "uint256",
              name: "minCreateAmount",
              type: "uint256",
            },
            {
              internalType: "uint256",
              name: "signatureReward",
              type: "uint256",
            },
          ],
          internalType: "struct XChainTypes.BridgeParams[]",
          name: "params",
          type: "tuple[]",
        },
      ],
      stateMutability: "view",
      type: "function",
    },
    {
      inputs: [],
      name: "getWitnesses",
      outputs: [
        {
          internalType: "address[]",
          name: "",
          type: "address[]",
        },
      ],
      stateMutability: "view",
      type: "function",
    },
    {
      inputs: [
        {
          internalType: "address",
          name: "token",
          type: "address",
        },
      ],
      name: "isTokenRegistered",
      outputs: [
        {
          internalType: "bool",
          name: "",
          type: "bool",
        },
      ],
      stateMutability: "view",
      type: "function",
    },
    {
      inputs: [],
      name: "owner",
      outputs: [
        {
          internalType: "address",
          name: "",
          type: "address",
        },
      ],
      stateMutability: "view",
      type: "function",
    },
    {
      inputs: [],
      name: "pause",
      outputs: [],
      stateMutability: "nonpayable",
      type: "function",
    },
    {
      inputs: [],
      name: "paused",
      outputs: [
        {
          internalType: "bool",
          name: "",
          type: "bool",
        },
      ],
      stateMutability: "view",
      type: "function",
    },
    {
      inputs: [],
      name: "renounceOwnership",
      outputs: [],
      stateMutability: "nonpayable",
      type: "function",
    },
    {
      inputs: [
        {
          internalType: "address",
          name: "newOwner",
          type: "address",
        },
      ],
      name: "transferOwnership",
      outputs: [],
      stateMutability: "nonpayable",
      type: "function",
    },
    {
      inputs: [],
      name: "unpause",
      outputs: [],
      stateMutability: "nonpayable",
      type: "function",
    },
    {
      stateMutability: "payable",
      type: "receive",
    },
  ],
};
