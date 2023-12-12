const { EXPLORER } = require("./constant")

async function create_trustline(client, wallet, currency_code, issuer_address){
    const trust_result = await client.submitAndWait({
        "TransactionType": "TrustSet",
        "Account": wallet.address,
        "LimitAmount": {
          "currency": currency_code,
          "issuer": issuer_address,
          "value": "10000000000" // Large limit, arbitrarily chosen
        }
      }, {autofill: true, wallet: wallet})
      if (trust_result.result.meta.TransactionResult == "tesSUCCESS") {
        console.log(`Trust line created: ${EXPLORER}/${trust_result.result.hash}`)
      } else {
        throw `Error sending transaction: ${trust_result}`
      }
}

module.exports = create_trustline;