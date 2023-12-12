const xrpl = require("xrpl");
const { EXPLORER } = require("./constant")
async function get_new_token(client, wallet, currency_code, issue_quantity) {
    // Get credentials from the Testnet Faucet -----------------------------------
    //console.log("Funding an issuer address with the faucet...")
    const issuer = (await client.fundWallet()).wallet
    //console.log(`Got issuer address ${issuer.address}.`)
  
    // Enable issuer DefaultRipple ----------------------------------------------
    const issuer_setup_result = await client.submitAndWait({
      "TransactionType": "AccountSet",
      "Account": issuer.address,
      "SetFlag": xrpl.AccountSetAsfFlags.asfDefaultRipple
    }, {autofill: true, wallet: issuer} )
    if (issuer_setup_result.result.meta.TransactionResult == "tesSUCCESS") {
      console.log(`Issuer DefaultRipple enabled: ${EXPLORER}/${issuer_setup_result.result.hash}`)
    } else {
      throw `Error sending transaction: ${issuer_setup_result}`
    }
  
    // Create trust line to issuer ----------------------------------------------
    const trust_result = await client.submitAndWait({
      "TransactionType": "TrustSet",
      "Account": wallet.address,
      "LimitAmount": {
        "currency": currency_code,
        "issuer": issuer.address,
        "value": "10000000000" // Large limit, arbitrarily chosen
      }
    }, {autofill: true, wallet: wallet})
    if (trust_result.result.meta.TransactionResult == "tesSUCCESS") {
      console.log(`Trust line created: ${EXPLORER}/${trust_result.result.hash}`)
    } else {
      throw `Error sending transaction: ${trust_result}`
    }
  
    // Issue tokens -------------------------------------------------------------
    const issue_result = await client.submitAndWait({
      "TransactionType": "Payment",
      "Account": issuer.address,
      "Amount": {
        "currency": currency_code,
        "value": issue_quantity,
        "issuer": issuer.address
      },
      "Destination": wallet.address
    }, {autofill: true, wallet: issuer})
    if (issue_result.result.meta.TransactionResult == "tesSUCCESS") {
      console.log(`Tokens issued: ${EXPLORER}/${issue_result.result.hash}`)
    } else {
      throw `Error sending transaction: ${issue_result}`
    }
  
    return {
      "currency": currency_code,
      "value": issue_quantity,
      "issuer": issuer.address,
      "seed": issuer.seed
    }
  }

module.exports = get_new_token;