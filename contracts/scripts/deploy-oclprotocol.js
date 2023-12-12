const ethers = require("hardhat").ethers;

async function main() {
    const [deployer] = await ethers.getSigners();
  
    console.log("Deploying contracts with the account:", deployer.address);
  
    const oclProtocolFactory = await ethers.getContractFactory("oclProtocol");
    const oclProtocol = await oclProtocolFactory.deploy("0x133EEf561F068511bFc2740A1Fd33192E246637E", "0xaf09826Fab7224678CfAE325C0FB2340f3454608", "0x3fec2f39e32ad951e0113f573736ff500126755a", "0x0FCCFB556B4aA1B44F31220AcDC8007D46514f31")
    console.log("oclProtocol address:", await oclProtocol.getAddress());
  }
  
  main()
    .then(() => process.exit(0))
    .catch((error) => {
      console.error(error);
      process.exit(1);
    });