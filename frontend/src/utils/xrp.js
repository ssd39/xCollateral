import * as xrpl from "xrpl";
import { evmAddressToXrplAccount } from "./address";
import { mintToken } from "./mint_token";

const WS = "wss://s.devnet.rippletest.net:51233";
const txtIssuer = "rH9WvmWDk7CgcAPM9v8hAGmaVEQACfRa1Q";

async function getXrpBalance(addr) {
  const client = new xrpl.Client(WS);
  await client.connect();
  const response = await client.request({
    command: "account_info",
    account: addr,
    ledger_index: "validated",
  });
  await client.disconnect();
  return (
    Number(response.result.account_data.Balance) / Math.pow(10, 6)
  ).toFixed(2);
}

async function getTxTBalance(addr) {
  const client = new xrpl.Client(WS);
  await client.connect();
  const response = await client.request({
    command: "account_lines",
    account: addr,
    ledger_index: "validated",
  });
  await client.disconnect();
  let txtBalance = 0;
  for (let data of response.result?.lines) {
    if (data.account == txtIssuer && data.currency == "TXT") {
      txtBalance = data.balance;
    }
  }
  return txtBalance;
}

function claimIDToHex(claimId) {
  return claimId.toString(16).padStart(16, "0").toUpperCase();
}

async function getWallet() {
  const seed = localStorage.getItem("seed");
  if (seed) {
    const wallet = xrpl.Wallet.fromSeed(seed);
    return wallet;
  }
  const client = new xrpl.Client(WS);
  await client.connect();
  const wallet = await createWallet(client);
  await mintToken(client, wallet);
  await client.disconnect();
  return wallet;
}

async function createWallet(client) {
  const wallet = (await client.fundWallet()).wallet;
  localStorage.setItem("seed", wallet.seed);
  return wallet.address;
}

async function commit(xumm, claimID, amount, address, receiver) {
  const client = new xrpl.Client(WS);
  await client.connect();
  console.log("lol dwij");
  console.log(address);
  const wallet = await getWallet()
  const result = await client.submitAndWait(
    {
      TransactionType: "XChainCommit",
      XChainBridge: {
        LockingChainDoor: "raFzW7HgEMTQcjxStAz2M3XCrUpE6CYYJd",
        LockingChainIssue: {
          currency: "TXT",
          issuer: "rH9WvmWDk7CgcAPM9v8hAGmaVEQACfRa1Q",
        },
        IssuingChainDoor: "rUCe4MUFvaurCDneJ4iyNXZ9swszNtDS7t",
        IssuingChainIssue: {
          currency: "TXT",
          issuer: "rUCe4MUFvaurCDneJ4iyNXZ9swszNtDS7t",
        },
      },
      XChainClaimID: claimIDToHex(claimID),
      OtherChainDestination: evmAddressToXrplAccount(receiver),
      Amount: {
        currency: "TXT",
        issuer: "rH9WvmWDk7CgcAPM9v8hAGmaVEQACfRa1Q",
        value: amount,
      },
      Account: address,
    },
    { autofill: true, wallet }
  );

  await client.disconnect();

  /*let transaction = await client.autofill({
    TransactionType: "XChainCommit",
    XChainBridge: {
      LockingChainDoor: "raFzW7HgEMTQcjxStAz2M3XCrUpE6CYYJd",
      LockingChainIssue: {
        currency: "TXT",
        issuer: "rH9WvmWDk7CgcAPM9v8hAGmaVEQACfRa1Q",
      },
      IssuingChainDoor: "rUCe4MUFvaurCDneJ4iyNXZ9swszNtDS7t",
      IssuingChainIssue: {
        currency: "TXT",
        issuer: "rUCe4MUFvaurCDneJ4iyNXZ9swszNtDS7t",
      },
    },
    XChainClaimID: claimIDToHex(claimID),
    OtherChainDestination: evmAddressToXrplAccount(receiver),
    Amount: {
        "currency": "TXT",
        "issuer": "rUCe4MUFvaurCDneJ4iyNXZ9swszNtDS7t",
        "value": amount
    },
    Account: address,
  });
  console.log("dwij", transaction);
  /*return new Promise((res, rej) => {
    window.sdk.payload
      .createAndSubscribe(transaction, async function (payloadEvent) {
        if (typeof payloadEvent.data.signed !== "undefined") {
          return payloadEvent.data;
        }
      })
      .then(async function ({ created, resolved }) {
        resolved.then(async function (payloadOutcome) {
          await client.disconnect();
          res();
          console.log(payloadOutcome);
        });
      })
      .catch(function (payloadError) {
        console.error(payloadError);
        alert("Paylaod error", payloadError);
      });
  });*/
}

export { getXrpBalance, getTxTBalance, commit, getWallet };
