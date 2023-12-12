require('dotenv').config()
require("@nomicfoundation/hardhat-toolbox");

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  networks: {
    xrpl: {
      url: "https://rpc-evm-sidechain.xrpl.org",
      chainId: 1440002,
      accounts: {
        mnemonic: process.env.MNEMONIC,
      },
    },
  },
  solidity: "0.8.20",
};
