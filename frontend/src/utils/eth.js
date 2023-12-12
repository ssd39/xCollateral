import Web3 from "web3";
import ERC20 from "../abis/ERC20";
import PriceOracle from "../abis/PriceOracle";
import MultiTokenDoor from "../abis/MultiTokenDoor";
import { xrplAccountToEvmAddress } from "./address";
import OclProtocol from "../abis/OclProtocol";
const TOKEN_CONTRACT = "0xaf09826Fab7224678CfAE325C0FB2340f3454608";
const PRICE_ORACLE_CONTRACT = "0x133EEf561F068511bFc2740A1Fd33192E246637E";
const ISSUING_DOOR_CONTRACT = "0x337BE5e12E59298a3384F3d8d95AaCE89465A62c";
const OCL_PROTOCOL_CONTRACT = "0x024e87F9b38070d092aE1568Fa872691e5b068dC";

async function getEvmBalance(provider, addr) {
  const web3 = new Web3(provider);
  const balance = await web3.eth.getBalance(addr);
  return (Number(balance.toString()) / Math.pow(10, 18)).toFixed(2);
}

async function getEvmTxTBalance(provider, addr) {
  const web3 = new Web3(provider);
  const token = new web3.eth.Contract(ERC20.abi, TOKEN_CONTRACT);
  return (
    Number((await token.methods.balanceOf(addr).call()).toString()) /
    Math.pow(10, 18)
  ).toFixed(2);
}

async function amountToClaim(evmAddr) {
  const web3 = new Web3(window.web3.currentProvider);
  const oclProtocol = new web3.eth.Contract(
    OclProtocol.abi,
    OCL_PROTOCOL_CONTRACT
  );
  return (
    Number(
      (await oclProtocol.methods.getClaimAmount(evmAddr).call()).toString()
    ) / Math.pow(10, 18)
  ).toFixed(2);
}

async function amountLent(evmAddr) {
  const web3 = new Web3(window.web3.currentProvider);
  const oclProtocol = new web3.eth.Contract(
    OclProtocol.abi,
    OCL_PROTOCOL_CONTRACT
  );
  return (
    Number(
      (
        await oclProtocol.methods.getLentAmountByAddress(evmAddr).call()
      ).toString()
    ) / Math.pow(10, 18)
  ).toFixed(2);
}

async function vaultBalance() {
  const web3 = new Web3(window.web3.currentProvider);
  const token = new web3.eth.Contract(ERC20.abi, TOKEN_CONTRACT);
  return (
    Number(
      (await token.methods.balanceOf(OCL_PROTOCOL_CONTRACT).call()).toString()
    ) / Math.pow(10, 18)
  ).toFixed(2);
}

async function lend(amount_, evmAddr) {
  const web3 = new Web3(window.web3.currentProvider);
  const token = new web3.eth.Contract(ERC20.abi, TOKEN_CONTRACT);
  const oclProtocol = new web3.eth.Contract(
    OclProtocol.abi,
    OCL_PROTOCOL_CONTRACT
  );
  const amount = parseFloat(amount_) * Math.pow(10, 18);
  if (
    Number(
      await token.methods.allowance(evmAddr, OCL_PROTOCOL_CONTRACT).call()
    ) < amount
  ) {
    console.log("roar it!!");
    await token.methods
      .approve(OCL_PROTOCOL_CONTRACT, amount)
      .send({ from: evmAddr });
  }
  await oclProtocol.methods.lend(amount).send({ from: evmAddr });
  alert("Lent sucessfully!");
}

async function borrowData(evmAddr) {
  const web3 = new Web3(window.web3.currentProvider);
  const oclProtocol = new web3.eth.Contract(
    OclProtocol.abi,
    OCL_PROTOCOL_CONTRACT
  );
  const borrowData = await oclProtocol.methods.getBorrowData(evmAddr).call();
  const out = [];
  for (let data of borrowData) {
    const bdId = await oclProtocol.methods.getBorrowId(data).call();
    if (!bdId.isLiquidated) {
      bdId.id = parseInt(data);
      bdId.txtAmountView = (
        Number(bdId.txtAmount.toString()) / Math.pow(10, 18)
      ).toFixed(2);
      bdId.startTimestampView =   parseInt((Number(bdId.startTimestamp.toString()) + 3600 * 24 * 30 - (Date.now()/1000)) / (3600 * 24))
      out.push(bdId);
    }
  }
  return out;
}

async function closeBorrow(evmAddr, amount, bid){
    const web3 = new Web3(window.web3.currentProvider);
    const token = new web3.eth.Contract(ERC20.abi, TOKEN_CONTRACT);
    const oclProtocol = new web3.eth.Contract(
      OclProtocol.abi,
      OCL_PROTOCOL_CONTRACT
    );
  
    if (
      Number(
        await token.methods.allowance(evmAddr, OCL_PROTOCOL_CONTRACT).call()
      ) < amount
    ) {
      console.log("roar it!!");
      await token.methods
        .approve(OCL_PROTOCOL_CONTRACT, amount)
        .send({ from: evmAddr });
    }
    await oclProtocol.methods.closeBorrow(bid).send({ from: evmAddr })
    alert("Loan closed sucessfully!")

}

async function borrow(evmAddr, amount_) {
  const web3 = new Web3(window.web3.currentProvider);
  const oclProtocol = new web3.eth.Contract(
    OclProtocol.abi,
    OCL_PROTOCOL_CONTRACT
  );
  const amount = parseFloat(amount_) * Math.pow(10, 18);
  await oclProtocol.methods.borrow().send({ from: evmAddr, value: amount });
}

async function getXrpTxtRate(provider) {
  const web3 = new Web3(provider);
  const oracle = new web3.eth.Contract(PriceOracle.abi, PRICE_ORACLE_CONTRACT);
  const amout = await oracle.methods.amount().call();
  const amout2 = await oracle.methods.amount2().call();
  return (Number(amout2.toString()) / Number(amout.toString())).toFixed(2);
}

async function createEvmClaim(provider, addr1, addr2) {
  const web3 = new Web3(provider);
  const doorContract = new web3.eth.Contract(
    MultiTokenDoor.abi,
    ISSUING_DOOR_CONTRACT
  );
  const tx = await doorContract.methods
    .createClaimId(
      [
        "0x3fec2f39e32ad951e0113f573736ff500126755a",
        ["0xb11e527233c77590da06056753dfe24f1cfc6b75", "TXT"],
        "0x7ff8622aee4d28f7848a64be82c99c30cbac4d9b",
        ["0x7ff8622aee4d28f7848a64be82c99c30cbac4d9b", "TXT"],
      ],
      xrplAccountToEvmAddress(addr1)
    )
    .send({ from: addr2, value: Math.pow(10, 18) });
  return tx.events.CreateClaim.returnValues.claimId;
}

export {
  getEvmBalance,
  getEvmTxTBalance,
  getXrpTxtRate,
  createEvmClaim,
  amountLent,
  amountToClaim,
  lend,
  borrow,
  vaultBalance,
  borrowData,
  closeBorrow
};
