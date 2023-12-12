import { decodeAccountID, encodeAccountID } from "ripple-address-codec";

export const xrplAccountToEvmAddress = (account) => {
    const accountId = decodeAccountID(account);
    return `0x${accountId.toString("hex")}`;
};

export const evmAddressToXrplAccount = (address) => {
    const accountId = Buffer.from(address.slice(2), "hex");
    return encodeAccountID(accountId);
};
