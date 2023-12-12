import "./App.css";
import DialogTitle from "@mui/material/DialogTitle";
import Dialog from "@mui/material/Dialog";
import DialogContent from "@mui/material/DialogContent";
import { useEffect, useState } from "react";
import Button from "@mui/material/Button";
import CircularProgress from "@mui/material/CircularProgress";
import { XummPkce } from "xumm-oauth2-pkce";
import { getXrpBalance, getTxTBalance, getWallet } from "./utils/xrp";
import { getEvmBalance, getEvmTxTBalance, getXrpTxtRate } from "./utils/eth";
import Web3 from "web3";
import Bridge from "./components/Bridge";
import Lend from "./components/Lend";
import Borrow from "./components/Borrow";
function App() {
  const [isShowApp, setShowApp] = useState(false);

  const [evmAddr, setEvmAddr] = useState("");
  const [isXummLoading, setXummLoading] = useState(false);
  const [xrpAddr, setXrpAddr] = useState("");
  const [xauth, setXauth] = useState(null);
  const [balnceLoader, setBalanceLoader] = useState(true);
  const [xrp1, setXrp1] = useState(0);
  const [txt1, setTxt1] = useState(0);
  const [xrp2, setXrp2] = useState(0);
  const [txt2, setTxt2] = useState(0);
  const [priceR, setPriceR] = useState(0);
  const [curMenu, setCurentMenu] = useState(0);
  const [connecting, setConnecting] = useState(false);

  const signedInHandler = (authorized) => {
    window.sdk = authorized.sdk;
  };

  const ConnectXumm = async () => {
    setXummLoading(true);

    /*var auth = new XummPkce("dd0b803a-6c67-4f97-98c3-2d728da39ba6",
    {
      implicit: true,
    });

    auth.on("error", (error) => {
      console.log("error", error);
    });

    auth.on("success", async () => {
      console.log("success");
      auth.state().then((state) => {
        if (state.me) {
          console.log("success, me", JSON.stringify(state.me));
          setXrpAddr(state.me.account);
          setXummLoading(false);

          signedInHandler(state);
        }
      });
    });

    auth.on("retrieved", async () => {
      // Redirect, e.g. mobile. Mobile may return to new tab, this
      // must retrieve the state and process it like normally the authorize method
      // would do
      console.log("Results are in, mobile flow, process sign in");

      auth.state().then((state) => {
        console.log(state);
        if (state) {
          console.log("retrieved, me:", JSON.stringify(state.me));
          setXrpAddr(state.me.account);
          setXummLoading(false);
          signedInHandler(state);
        }
      });
    });

    await auth.authorize();*/
    const wallet = await getWallet();
    setXrpAddr(wallet.address);
    setXummLoading(false);
  };

  const ConenctMetamask = async () => {
    try {
      setConnecting(true);
      if (!window.ethereum) {
        setConnecting(false);
        alert("Metamask not found!");
        return;
      }
      const accounts = await window.ethereum.request({
        method: "eth_requestAccounts",
      });
      const account = accounts[0];
      setEvmAddr(account);
      setConnecting(false);
    } catch (err) {
      console.warn(`failed to connect..`, err);
    }
  };

  const updateBalance = async () => {
    setBalanceLoader(true);
    let t = await getXrpBalance(xrpAddr);
    let t2 = await getTxTBalance(xrpAddr);
    let t3 = await getEvmBalance(window.web3.currentProvider, evmAddr);
    let t4 = await getEvmTxTBalance(window.web3.currentProvider, evmAddr);
    let t5 = await getXrpTxtRate(window.web3.currentProvider);
    setTxt1(t2);
    setXrp1(t);
    setXrp2(t3);
    setTxt2(t4);
    setPriceR(t5);
    setBalanceLoader(false);
  };

  useEffect(() => {
    if (xrpAddr != "" && evmAddr !== "") {
      setShowApp(true);
      updateBalance();
    }
  }, [xrpAddr, evmAddr]);

  return (
    <div className="flex flex-col w-full min-h-screen bg-black">
      {isShowApp && (
        <>
          <div className="flex justify-between px-4 items-center py-2">
            <div>
              <span className="text-white font-extrabold text-4xl">
                xCollateral
              </span>
            </div>
            <div className="text-white flex flex-col">
              <div className="flex">
                <span className="font-bold">Connected Wallets</span>
              </div>
              <div className="flex flex-col">
                <div className="rounded-full bg-purple-500 px-2">
                  <span className="font-semibold">XRP:</span>
                  <span className="ml-2">{xrpAddr}</span>
                </div>
                <div className="rounded-full bg-blue-500 px-2 mt-1">
                  <span className="font-semibold ">Evm:</span>
                  <span className="ml-2">{evmAddr}</span>
                </div>
              </div>
            </div>
          </div>
          <div className="w-full bg-white" style={{ height: 2 }}></div>
          <div className="mt-4 w-full">
            <div className="w-full">
              {!balnceLoader && (
                <div
                  className="flex rounded-lg flex-col w-full justify-center items-center py-4"
                  style={{ borderColor: "white", borderWidth: 2 }}
                >
                  <div className="text-white rounded-full bg-red-500 px-2">
                    <span className="font-semibold">Rate:</span>
                    <span className="ml-2">{priceR} TXT/XRP</span>
                  </div>
                  <div className="flex text-white text-2xl font-semibold">
                    Balances:
                  </div>
                  <div className="text-white flex mt-2">
                    <div className="flex flex-col items-center justify-center mx-2">
                      <span className="font-semibold">XRPL</span>
                      <div className="flex flex-col mt-1">
                        <div className="rounded-full bg-green-500 px-2">
                          <span className="font-semibold">XRP:</span>
                          <span className="ml-2 text-lg">{xrp1}</span>
                        </div>
                        <div className="rounded-full bg-orange-500 px-2 mt-1">
                          <span className="font-semibold ">TXT:</span>
                          <span className="ml-2 text-lg">{txt1}</span>
                        </div>
                      </div>
                    </div>
                    <div className="flex flex-col items-center justify-center mx-2">
                      <span className="font-semibold">EVM Chain</span>
                      <div className="flex flex-col mt-1">
                        <div className="rounded-full bg-green-500 px-2">
                          <span className="font-semibold">XRP:</span>
                          <span className="ml-2 text-lg">{xrp2}</span>
                        </div>
                        <div className="rounded-full bg-orange-500 px-2 mt-1">
                          <span className="font-semibold ">TXT:</span>
                          <span className="ml-2 text-lg">{txt2}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              )}
              {balnceLoader && (
                <div
                  className="flex rounded-lg flex-col w-full justify-center items-center py-4"
                  style={{ borderColor: "white", borderWidth: 2 }}
                >
                  <CircularProgress />
                </div>
              )}
            </div>
            <div className="w-full  mt-3 py-2 flex justify-center items-center">
              <div
                onClick={() => setCurentMenu(0)}
                className={`rounded-full mx-5 py-1 hover:bg-sky-700  px-3 text-white cursor-pointer	 font-semibold ${
                  curMenu == 0 ? "bg-blue-800" : ""
                }`}
              >
                Borrow
              </div>
              <div
                onClick={() => setCurentMenu(1)}
                className={`rounded-full mx-5 py-1 hover:bg-sky-700 px-3 text-white  cursor-pointer	 font-semibold  ${
                  curMenu == 1 ? "bg-blue-800" : ""
                }`}
              >
                Bridge
              </div>
              <div
                onClick={() => setCurentMenu(2)}
                className={`rounded-full mx-5 py-1 hover:bg-sky-700 px-3 text-white  cursor-pointer	 font-semibold ${
                  curMenu == 2 ? "bg-blue-800" : ""
                }`}
              >
                Lend
              </div>
            </div>
          </div>
          {curMenu == 0 && (
            <div className="w-full flex flex-1">
              <Borrow
                onUpdate={() => {
                  updateBalance();
                }}
                evmAddr={evmAddr}
                rate={priceR}
              />
            </div>
          )}
          {curMenu == 1 && (
            <div className="w-full flex flex-1">
              <Bridge xumm={xauth} evmAddr={evmAddr} xrpAddr={xrpAddr} />
            </div>
          )}
          {curMenu == 2 && (
            <div className="w-full flex flex-1">
              <Lend
                onUpdate={() => {
                  updateBalance();
                }}
                evmAddr={evmAddr}
              />
            </div>
          )}
        </>
      )}
      <Dialog
        sx={{
          "& .MuiPaper-root": {
            background: "#000",
            borderColor: "white",
            borderWidth: "1px",
          },
          "& .MuiBackdrop-root": {
            backgroundColor: "transparent", // Try to remove this to see the result
          },
        }}
        open={!isShowApp}
      >
        <DialogTitle sx={{ fontWeight: "bold", color: "white" }}>
          <div className="flex flex-col">
            <span>Welcome to the xCollateral</span>
            <span className="text-sm">Connect wallets to continue</span>
            <div className="mt-2 w-full bg-white" style={{ height: 1 }}></div>
          </div>
        </DialogTitle>
        <DialogContent>
          <div className="flex flex-col  w-full">
            <div className="w-full flex  justify-center items-center">
              {!isXummLoading && xrpAddr == "" && (
                <Button variant="contained" onClick={ConnectXumm}>
                  Import/Create XRPL Wallet
                </Button>
              )}
              {isXummLoading && <CircularProgress color="secondary" />}
            </div>
            <div className="w-full flex justify-center items-center mt-2">
              {!connecting && evmAddr == "" && (
                <Button variant="contained" onClick={ConenctMetamask}>
                  Connect Metamask
                </Button>
              )}
              {connecting && <CircularProgress color="secondary" />}
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}

export default App;
