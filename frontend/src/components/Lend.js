import React, { useEffect, useState } from "react";
import { Button } from "@mui/material";
import { amountLent, amountToClaim, lend } from "../utils/eth";
import CircularProgress from "@mui/material/CircularProgress";

export default function Lend({ evmAddr, onUpdate }) {
  const [loading, setLoading] = useState(false);
  const [lentInput, setLentInput] = useState("");
  const [claim, setClaim] = useState("");
  const [lent, setLent] = useState("");

  const updateData = async () => {
    const lent_ = await amountLent(evmAddr);
    const claim_ = await amountToClaim(evmAddr);
    setClaim(claim_);
    setLent(lent_);
    setLoading(false);
  };

  useEffect(() => {
    setLoading(true);
    updateData();
  }, []);
  return (
    <div className=" w-full flex flex-col items-center  my-8">
      {!loading && (
        <div className=" w-full flex flex-col items-center ">
          <div className="text-white">
            <span className="text-white  text-xl font-semibold">
              Current Lent Amount:{" "}
            </span>
            <span className="ml-2 text-lg">{lent} TXT</span>
          </div>
          <div className="mt-4 flex justify-center items-center">
            <input
              onChange={(e) => {
                setLentInput(e.target.value);
              }}
              value={lentInput}
              className="py-1 px-1 text-lg w-52"
              type="text"
            />
            <Button
              onClick={async () => {
                await lend(lentInput, evmAddr);
                updateData()
                onUpdate()
              }}
              variant="contained"
            >
              Lend Amount
            </Button>
          </div>
          <div className="mt-6 text-white flex flex-col items-center">
            <span className="font-bold text-xl mb-2">Actions:</span>
            <Button variant="contained">Claim {claim} TXT</Button>
            <div  className="mt-4">
            <Button variant="contained">withdraw Stake</Button>

            </div>

          </div>
        </div>
      )}
      {loading && <CircularProgress />}
    </div>
  );
}
