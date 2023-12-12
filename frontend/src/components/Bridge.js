import React, { useRef, useState } from "react";
import { Button } from "@mui/material";
import HeightIcon from "@mui/icons-material/Height";
import { createEvmClaim } from "../utils/eth";
import { commit } from "../utils/xrp";
export default function Bridge({ xumm, evmAddr, xrpAddr }) {
  const t1 = useRef(0);
  const t2 = useRef(0);
  const [direction, setDirection] = useState(0);

  const XrpInput = ({ onUpdate, val }) => {
    const [t, setT] = useState(val);
    return (
      <div className="flex flex-col">
        <span className="text-white text-lg font-bold">XRPL (TXT):</span>
        <input
          onChange={(e) => {
            setT(e.target.value);
            onUpdate(e.target.value);
          }}
          value={t}
          className="py-1 px-1 text-lg"
          type="text"
        />
      </div>
    );
  };

  const TxtInput = ({ onUpdate, val }) => {
    const [t, setT] = useState(val);
    return (
      <div className="flex flex-col">
        <span className="text-white text-lg font-bold">EVM Chain (TXT):</span>
        <input
          onChange={(e) => {
            setT(e.target.value);
            onUpdate(e.target.value);
          }}
          value={t}
          className="py-1 px-1 text-lg"
          type="text"
        />
      </div>
    );
  };

  const evmToXrp = async () => {};

  const xrpToEvm = async () => {
    const claimId = await createEvmClaim(
      window.web3.currentProvider,
      xrpAddr,
      evmAddr
    );
    await commit(xumm, claimId, t1.current, xrpAddr, evmAddr);
    alert("Bridged succeed!");
  };

  return (
    <div className="w-full flex-1  flex  flex-col my-8 items-center">
      <span className="text-white text-xl font-semibold">TXT Bridge</span>
      <div className="my-5 flex flex-col justify-center items-center rounded-lg p-4" style={{ borderWidth: 2, borderColor:'white'}}>
        {direction == 0 ? (
          <XrpInput val={t1.current} onUpdate={(e) => (t1.current = e)} />
        ) : (
          <TxtInput val={t2.current} onUpdate={(e) => (t1.current = e)} />
        )}
        <div className="my-4">
          <Button
            onClick={() => setDirection(direction == 0 ? 1 : 0)}
            variant="contained"
            startIcon={<HeightIcon />}
          >
            Switch
          </Button>
        </div>
        {direction == 0 ? (
          <span className="text-white text-lg font-bold">To EVM Chain (TXT)</span>
        ) : (
            <span className="text-white text-lg font-bold">To XRPL (TXT)</span>

        )}

      </div>
      <div className="mt-4">
          <Button
            onClick={direction == 0 ? xrpToEvm : evmToXrp}
            variant="contained"
          >
            Transfer
          </Button>
        </div>
    </div>
  );
}
