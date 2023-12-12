export default {
    abi: [
        {
          "inputs": [
            {
              "internalType": "address",
              "name": "doorAccount_",
              "type": "address"
            },
            {
              "internalType": "string",
              "name": "currency_",
              "type": "string"
            },
            {
              "internalType": "string",
              "name": "currency2_",
              "type": "string"
            }
          ],
          "stateMutability": "nonpayable",
          "type": "constructor"
        },
        {
          "inputs": [],
          "name": "amount",
          "outputs": [
            {
              "internalType": "uint256",
              "name": "",
              "type": "uint256"
            }
          ],
          "stateMutability": "view",
          "type": "function"
        },
        {
          "inputs": [],
          "name": "amount2",
          "outputs": [
            {
              "internalType": "uint256",
              "name": "",
              "type": "uint256"
            }
          ],
          "stateMutability": "view",
          "type": "function"
        },
        {
          "inputs": [],
          "name": "currency",
          "outputs": [
            {
              "internalType": "string",
              "name": "",
              "type": "string"
            }
          ],
          "stateMutability": "view",
          "type": "function"
        },
        {
          "inputs": [],
          "name": "currency2",
          "outputs": [
            {
              "internalType": "string",
              "name": "",
              "type": "string"
            }
          ],
          "stateMutability": "view",
          "type": "function"
        },
        {
          "inputs": [],
          "name": "doorAccount",
          "outputs": [
            {
              "internalType": "contract DoorAccount",
              "name": "",
              "type": "address"
            }
          ],
          "stateMutability": "view",
          "type": "function"
        },
        {
          "inputs": [
            {
              "internalType": "uint256",
              "name": "amount_",
              "type": "uint256"
            },
            {
              "internalType": "uint256",
              "name": "amount2_",
              "type": "uint256"
            }
          ],
          "name": "updateData",
          "outputs": [],
          "stateMutability": "nonpayable",
          "type": "function"
        }
      ]
}