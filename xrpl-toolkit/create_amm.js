const xrpl = require("xrpl");
const { EXPLORER } = require("./constant")

async function createAmm(client, wallet, amount1, amount2) {
    const ss = await client.request({"command": "server_state"})
    const amm_fee_drops = ss.result.state.validated_ledger.reserve_inc.toString()
    console.log(`Current AMMCreate transaction cost: 
                 ${xrpl.dropsToXrp(amm_fee_drops)} XRP`)
    const ammcreate_result = await client.submitAndWait({
        "TransactionType": "AMMCreate",
        "Account": wallet.address,
        "Amount": amount1, /*{
            currency: "TST",
            issuer: "rP9jPyP5kyvFRb6ZiRghAGw5u8SGAmU4bd",
            value: "15"
        },*/
        "Amount2": amount2, /*{
            "currency": foo_amount.currency,
            "issuer": foo_amount.issuer,
            "value": "100"
        },*/
        "TradingFee": 500, // 0.5%
        "Fee": amm_fee_drops
        }, {autofill: true, wallet: wallet, fail_hard: true})
    // Use fail_hard so you don't waste the tx cost if you mess up
    if (ammcreate_result.result.meta.TransactionResult == "tesSUCCESS") {
        console.log(`AMM created: ${EXPLORER}/transactions/${ammcreate_result.result.hash}`)
    } else {
        console.error(`Error sending transaction: ${ammcreate_result}`)
    }
}

module.exports = createAmm