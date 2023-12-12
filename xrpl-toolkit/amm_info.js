async function getAmmInfo(client, asset1, asset2) {
  const amm_info_request = {
    "command": "amm_info", 
    "asset":asset1, /*{
      "currency": "TST",
      "issuer": "rP9jPyP5kyvFRb6ZiRghAGw5u8SGAmU4bd",
    },*/
    "asset2": asset2, /* {
       "currency": "XRP"
    },*/
    "ledger_index": "validated"
  }
  try {
    const amm_info_result = await client.request(amm_info_request)
    console.log(JSON.stringify(amm_info_result))
    return amm_info_result
  } catch(err) {
    if (err.data.error === 'actNotFound') {
      console.log(`No AMM exists yet for the pair 
                   XRP / 
                   TST.rP9jPyP5kyvFRb6ZiRghAGw5u8SGAmU4bd.
                   (This is probably as expected.)`)
    } else {
      console.error(err)
    }
  }
  return null
}

module.exports = getAmmInfo;