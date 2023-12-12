import * as xrpl from "xrpl";

async function mintToken(client, wallet){
 
    const issuer  = xrpl.Wallet;

    const trust_result = await client.submitAndWait({
        "TransactionType": "TrustSet",
        "Account": wallet.address,
        "LimitAmount": {
          "currency": "TXT",
          "issuer": issuer.address,
          "value": "10000000000" // Large limit, arbitrarily chosen
        }
      }, {autofill: true, wallet: wallet})
      if (trust_result.result.meta.TransactionResult == "tesSUCCESS") {
      } else {
        throw `Error sending transaction: ${trust_result}`
      }
    
      // Issue tokens -------------------------------------------------------------
      const issue_result = await client.submitAndWait({
        "TransactionType": "Payment",
        "Account": issuer.address,
        "Amount": {
          "currency": "TXT",
          "value": "1000",
          "issuer": issuer.address
        },
        "Destination": wallet.address
      }, {autofill: true, wallet: issuer})
      if (issue_result.result.meta.TransactionResult == "tesSUCCESS") {
      } else {
        throw `Error sending transaction: ${issue_result}`
      }
}

export {mintToken}