const xrpl = require("xrpl");
const createToken = require("./create_token");
const getAmmInfo = require("./amm_info");
const createAmm = require("./create_amm");
const createTrustLine = require("./create_trustline");
const mintToken = require("./mint_token");
var ellipticcurve = require("starkbank-ecdsa");
var Ecdsa = ellipticcurve.Ecdsa;
var PrivateKey = ellipticcurve.PrivateKey;

const { Command } = require("commander");
const fs = require("fs");

const program = new Command();

const WS_URL = "wss://s.devnet.rippletest.net:51233/";

async function main() {
  process.stdin.resume();

  console.log("Connecting xrpl client!");

  const client = new xrpl.Client(WS_URL);
  await client.connect();

  process.on("exit", async () => {
    console.log("Exiting cli :)");
    await client.disconnect();
  });

  console.log("xrpl client connected!");

  let wallet = null;
  if (fs.existsSync("./wallet.json")) {
    console.log("Importing exsisting wallet");
    const seedJson = fs.readFileSync("./wallet.json");
    wallet = xrpl.Wallet.fromSeed(JSON.parse(seedJson).seed);
  } else {
    console.log("Creating new wallet!");
    wallet = (await client.fundWallet()).wallet;
    fs.writeFileSync(
      "./wallet.json",
      JSON.stringify({
        seed: wallet.seed,
      })
    );
  }

  program
    .command("create-token")
    .argument("<currency_code>", "Currency Code")
    .argument("<issue_quantity>", "Quantity")
    //.argument("<issuer>", "Issuer")
    .action(async (a, b) => {
      //const seedJson = fs.readFileSync(c)
      //const issuer = xrpl.Wallet.fromSeed(JSON.parse(seedJson).seed)
      const newToken = await createToken(client, wallet, a, b);
      fs.writeFileSync(`./tokens/${a}_token.json`, JSON.stringify(newToken));
      process.exit();
    });

  program
    .command("create-trustline")
    .argument("<token-file>", "Token file")
    .action(async (a) => {
      const token = JSON.parse(fs.readFileSync(a));
      await createTrustLine(client, wallet, token.currency, token.issuer);
      process.exit();
    });

  program
    .command("create-amm")
    .argument("<token-file-1>", "Token file 1")
    .argument("<token-file-2>", "Token file 2")
    .argument("<amount-1>", "Amount 1")
    .argument("<amount-2>", "Amount 2")
    .argument("<wallet>", "wallet")
    .action(async (a, b, c, d, e) => {
      const token1 = JSON.parse(fs.readFileSync(a));
      const token2 = JSON.parse(fs.readFileSync(b));
      const amount1 = c;
      const amount2 = d;
      const wallet = xrpl.Wallet.fromSeed(JSON.parse(fs.readFileSync(e)).seed);
      let tokenObj1 = {
        currency: token1.currency,
      };
      if (token1.currency == "XRP") {
        tokenObj1 = xrpl.xrpToDrops(amount1);
      } else {
        tokenObj1.issuer = token1.issuer;
        tokenObj1.value = amount1;
      }
      let tokenObj2 = {
        currency: token2.currency,
      };
      if (token2.currency == "XRP") {
        tokenObj2 = xrpl.xrpToDrops(amount2);
      } else {
        tokenObj2.issuer = token2.issuer;
        tokenObj2.value = amount2;
      }
      await createAmm(client, wallet, tokenObj1, tokenObj2);
      process.exit();
    });
  program
    .command("mint-token")
    .argument("<wallet-dest>", "Wallet dest")
    .argument("<token-file>", "Token file")
    .argument("<amount>", "Amount")
    .action(async (c, a, b) => {
      const token = JSON.parse(fs.readFileSync(a));
      const issuer = xrpl.Wallet.fromSeed(token.seed);
      const destWallet = JSON.parse(fs.readFileSync(c));
      const wallet = xrpl.Wallet.fromSeed(destWallet.seed);
      await mintToken(client, wallet, b, issuer, token.currency);
      process.exit(0);
    });

  program
    .command("get-pk")
    .argument("<wallet-file>", "Wallet file")
    .action(async (a) => {
      const wallet_json = JSON.parse(fs.readFileSync(a));
      const wallet = xrpl.Wallet.fromSeed(wallet_json.seed);
      console.log(wallet.address);
      xrpl.Wallet.
      process.exit(0);
    });

  program
    .command("wallet-from-pk")
    .argument("<pub-key>", "public key")
    .argument("<pk>", "private key")
    .action(async (a, b) => {
      const wallet = new xrpl.Wallet(a, b);
      //await client.fundWallet(wallet)
      console.log(wallet.address);
      process.exit(0);
    });

  program
    .command("create-wallet")
    .argument("<wallet-name>", "Wallet name")
    .action(async (a) => {
      const wallet = (await client.fundWallet()).wallet
      fs.writeFileSync(`./${a}_wallet.json`, JSON.stringify({
        seed: wallet.seed,
        address: wallet.address
      }))
      process.exit(0);
    });

  program
    .command("amm-info")
    .argument("<token-file-1>", "Token file 1")
    .argument("<token-file-2>", "Token file 2")
    .action(async (a, b) => {
      const token1 = JSON.parse(fs.readFileSync(a));
      const token2 = JSON.parse(fs.readFileSync(b));
      let tokenObj1 = {
        currency: token1.currency,
      };
      if (token1.currency != "XRP") {
        tokenObj1.issuer = token1.issuer;
      }
      let tokenObj2 = {
        currency: token2.currency,
      };
      if (token2.currency != "XRP") {
        tokenObj2.issuer = token2.issuer;
      }
      await getAmmInfo(client, tokenObj1, tokenObj2);
      process.exit();
    });

  program.parse();
}

main();
