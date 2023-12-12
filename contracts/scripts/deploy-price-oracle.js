const ethers = require("hardhat").ethers;

async function main() {
    const [deployer] = await ethers.getSigners();
  
    console.log("Deploying contracts with the account:", deployer.address);
  
    const PriceOracleFactory = await ethers.getContractFactory("PriceOracle");
    const PriceOracle = await PriceOracleFactory.deploy("0x337BE5e12E59298a3384F3d8d95AaCE89465A62c", "XRP", "TXT")
    console.log("PriceOracle address:", await PriceOracle.getAddress());
    
  }
  
  main()
    .then(() => process.exit(0))
    .catch((error) => {
      console.error(error);
      process.exit(1);
    });