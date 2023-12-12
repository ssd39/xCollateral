import React, { useEffect, useState } from "react";
import { Button } from "@mui/material";
import { borrow, vaultBalance, borrowData, closeBorrow } from "../utils/eth";
import CircularProgress from "@mui/material/CircularProgress";
import ArrowDownwardIcon from "@mui/icons-material/ArrowDownward";

export default function Borrow({ evmAddr, onUpdate, rate }) {
  const [loading, setLoading] = useState(false);
  const [xrpInput, setXrpInput] = useState("");
  const [vaultB, setVaultB] = useState("");
  const [borrowList, setBorrowList] = useState([]);

  const updateData = async () => {
    vaultBalance().then((d) => {
      setVaultB(d);
    });
    borrowData(evmAddr).then((d) => {
      setBorrowList(d);
    });
  };
  useEffect(() => {
    updateData();
  }, []);
  return (
    <div className=" w-full flex flex-col items-center  my-8">
      <div className="w-full flex-1  flex  flex-col items-center">
        <div
          className="my-5 flex flex-col justify-center items-center rounded-lg p-4"
          style={{ borderWidth: 2, borderColor: "white" }}
        >
          <div className="my-4 flex flex-col items-center justify-center">
            <span className="text-white font-semibold mb-2">
              Vault Balance: {vaultB}
            </span>
            <div className="flex justify-center">
              <input
                value={xrpInput}
                onChange={(e) => {
                  setXrpInput(e.target.value);
                }}
                className="py-1 px-1 text-lg mb-3"
                type="text"
              />
              <span className="ml-2 text-white">XRP</span>
            </div>
            <ArrowDownwardIcon className="text-white" color="white" />
            <span className="text-white font-semibold">
              {parseFloat(isNaN(parseFloat(xrpInput)) ? 0 : xrpInput) *
                rate *
                0.85}{" "}
              TXT
            </span>
          </div>
        </div>
        <div className="mt-4">
          <Button
            variant="contained"
            onClick={async () => {
              await borrow(evmAddr, xrpInput);
              onUpdate();
              updateData();
              alert("TXT borrowed sucessfully!");
            }}
          >
            Borrow
          </Button>
        </div>
        <div className="w-full mx-2 bg-white my-6" style={{ height: 2 }}></div>
        <div className="w-full flex flex-col items-center">
        <span className="font-bold text-xl text-white mb-4">Active Loans:</span>
          {borrowList.map((val, ind) => {
            return (
              <div style={{width:'550px'}} className=" bg-white flex justify-between items-center rounded-lg px-4 py-2">
                <div>
                  <span className="font-semibold">Id: </span>
                  <span className="ml-2">{val.id}</span>
                </div>
                <div>
                  <span className="font-semibold">Expires In:</span>
                  <span className="ml-2">{ parseInt(parseInt(val.startTimestampView))} Days</span>
                </div>
                <div>
                  <span className="font-semibold">Amount:</span>
                  <span className="ml-2">{val.txtAmountView} TXT</span>
                </div>
                <div>
                <Button variant="contained" onClick={async ()=>{
                    await closeBorrow(evmAddr, val.txtAmount, val.id)
                    updateData()
                }}>
                  Close
                </Button>
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
