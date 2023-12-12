const { EXPLORER } = require("./constant")
async function mintToken(client, wallet, amount, issuer, currency){
    const trust_result = await client.submitAndWait({
        "TransactionType": "TrustSet",
        "Account": wallet.address,
        "LimitAmount": {
          "currency": currency,
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
          "currency": currency,
          "value": amount,
          "issuer": issuer.address
        },
        "Destination": wallet.address
      }, {autofill: true, wallet: issuer})
      if (issue_result.result.meta.TransactionResult == "tesSUCCESS") {
        console.log(`Tokens issued: ${EXPLORER}/${issue_result.result.hash}`)
      } else {
        throw `Error sending transaction: ${issue_result}`
      }
}

module.exports = mintToken;