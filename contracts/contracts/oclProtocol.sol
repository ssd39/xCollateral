//SPDX-License-Identifier: UNLICENSED

pragma solidity ^0.8.20;

import './PriceOracle.sol';
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

interface BridgeDoorNative{
    function commit(address receiver, uint256 claimId, uint256 amount) external payable;
}


contract oclProtocol{
    PriceOracle priceOracle;
    ERC20 txtToken;

    struct borrowData {
        address borrower;
        uint startTimestamp;
        uint256 xrpAmount;
        uint256 txtAmount;
        uint256 xrpReward;
        bool isLiquidated;
    }

    uint256 public txtLocked;

    uint256 public txtReward;
    uint256 public rewardsToLiquidate;
    uint256 public curTxTReward;
    mapping(address => uint) public lastClaim;
    mapping(uint => uint) public rewardClaimAmounts;

    mapping(address => uint) public lendBalances;
    mapping(uint => borrowData) public borrowedEntries;
    mapping(address => uint[]) public borrowDataByAddress;

    uint public lastBorrowId;
    uint public lastRewardBatchId;
    uint public lastRewardTimeStamp;
    uint public loanDurationBlocks = 1000;
    uint public collateralRatio = 7;
    uint public yieldPT = 15;
    address public xChainDoor;
    BridgeDoorNative nativeBridgeDoor;
    event AmmLiquidate(uint256 amount);

    constructor(address priceOracledAddress, address tokenAddress, address xChainDoor_, address nativeBridgeDoor_) {  
        priceOracle = PriceOracle(priceOracledAddress);
        txtToken = ERC20(tokenAddress);
        txtLocked = 0;
        txtReward = 0;
        lastBorrowId = 0;
        lastRewardBatchId = 0;
        curTxTReward = 0;
        rewardsToLiquidate = 0;
        lastRewardTimeStamp = block.timestamp;
        xChainDoor = xChainDoor_;
        lastClaim[address(this)] = 0;
        nativeBridgeDoor = BridgeDoorNative(nativeBridgeDoor_);
    }

    function lend(uint256 lendingAmount) public {
        require(txtToken.allowance(msg.sender, address(this)) >= lendingAmount, "Allowance required");
        require(txtToken.balanceOf(msg.sender) >= lendingAmount, "Low balance");
        if(txtToken.transferFrom(msg.sender, address(this), lendingAmount)){
            lendBalances[msg.sender] += lendingAmount;
        }
    }

    function withdraw(uint withdrawAmount) public {
        require(lendBalances[msg.sender]>=withdrawAmount, "Withdrawing more then your balance!");
        require(lastClaim[msg.sender]<lastRewardBatchId, "Claim your rewards before withdrawing!");
        require(lendBalances[msg.sender] * (txtToken.balanceOf(address(this))-txtReward) > withdrawAmount * (txtToken.balanceOf(address(this)) + txtLocked - txtReward), "Withdrawing more then current possible your stake");
        txtToken.transfer(msg.sender, withdrawAmount);
        lendBalances[msg.sender] -= withdrawAmount;
    }

    function claimReward() public {
        require(lendBalances[msg.sender]<=0, "You are not lender!");
        if(block.timestamp - lastRewardTimeStamp >= 3600) {
            lastRewardBatchId += 1;
            rewardClaimAmounts[lastRewardBatchId] = curTxTReward;
            curTxTReward = 0;
            lastRewardTimeStamp = block.timestamp;
        }
        uint pendingClaims = lastRewardBatchId-lastClaim[msg.sender]; 
        if(pendingClaims > 30) {
            pendingClaims = 30; // in one go can only claim first 30 days rewards
        }
        uint totalClaimAmount = 0;
        for(uint i=lastClaim[msg.sender]+1; i<=lastClaim[msg.sender]+pendingClaims;i++){
            totalClaimAmount += (lendBalances[msg.sender] / txtToken.balanceOf(address(this)) + txtLocked - txtReward) * rewardClaimAmounts[i];
        }
        require(txtToken.balanceOf(address(this)) >= totalClaimAmount, "Not enough balance in vault");
        lastClaim[msg.sender] = lastClaim[msg.sender] + pendingClaims;
        txtReward -= totalClaimAmount;
        txtToken.transfer(msg.sender, totalClaimAmount);
    }

    function closeBorrow(uint borrowId) public {
        require(!borrowedEntries[borrowId].isLiquidated, "loan already closed!");
        uint256 txtReward_ = (borrowedEntries[borrowId].xrpReward * priceOracle.amount2()) / priceOracle.amount();
        require(txtToken.allowance(msg.sender, address(this)) >= borrowedEntries[borrowId].txtAmount, "Allowance required");
        require(txtToken.balanceOf(msg.sender) >= borrowedEntries[borrowId].txtAmount, "Low balance");
        if(txtToken.transferFrom(msg.sender, address(this), borrowedEntries[borrowId].txtAmount)){
            txtReward += txtReward_;
            curTxTReward += txtReward_;
            borrowedEntries[borrowId].isLiquidated = true;
            //borrowDataByAddress[borrowedEntries[borrowId].borrower][borrowId] = false;
            rewardsToLiquidate += borrowedEntries[borrowId].xrpReward;
            payable(msg.sender).transfer(borrowedEntries[borrowId].xrpAmount -borrowedEntries[borrowId].xrpReward );
        }
    }

    function liquidateRewards(uint claimId) public {
        address[] memory witnessList = priceOracle.doorAccount().getWitnesses();
        bool isWintesss = false;
        for(uint i = 0; i<witnessList.length;i++){
            if(witnessList[i] == msg.sender){
                isWintesss = true;
            }
        }
        require(isWintesss, "Only witness can call this");
        require(rewardsToLiquidate > 0, "no rewards to liqiuidate");
        nativeBridgeDoor.commit{value: rewardsToLiquidate}(xChainDoor , claimId, rewardsToLiquidate);
        emit AmmLiquidate(rewardsToLiquidate);
        rewardsToLiquidate = 0;
    }

    function liquidate(uint borrowId, uint claimId) public {
        require(!borrowedEntries[borrowId].isLiquidated, "loan already closed!");
        uint256 xrpToTxt = (borrowedEntries[borrowId].xrpAmount * priceOracle.amount2()) / priceOracle.amount();
        uint256 txtReward_ = (borrowedEntries[borrowId].xrpReward * priceOracle.amount2()) / priceOracle.amount();
        if(txtReward_ * 100 < collateralRatio * xrpToTxt || block.timestamp - borrowedEntries[borrowId].startTimestamp > 30 * 3600 ){
            txtReward += txtReward_;
            curTxTReward += txtReward_;
            borrowedEntries[borrowId].isLiquidated = true;
            //borrowDataByAddress[borrowedEntries[borrowId].borrower][borrowId] = false;
            nativeBridgeDoor.commit{value: borrowedEntries[borrowId].xrpAmount}(xChainDoor , claimId, borrowedEntries[borrowId].xrpAmount);
            emit AmmLiquidate(borrowedEntries[borrowId].xrpAmount);
        }
    }

    function borrow() public payable {
       uint256 xrpToTxt = (msg.value * priceOracle.amount2()) / priceOracle.amount();
       uint256 collatral = (yieldPT * xrpToTxt) / 100;
       uint256 xrpReward = (yieldPT * msg.value) / 100;
       require(xrpToTxt-collatral <= (txtToken.balanceOf(address(this))-txtReward), "Not enough balnce txt to borrow loan");
       txtToken.transfer(msg.sender, xrpToTxt-collatral);
       lastBorrowId += 1;
       borrowedEntries[lastBorrowId] = borrowData({ borrower: msg.sender, startTimestamp: block.timestamp, xrpReward: xrpReward, xrpAmount: msg.value, txtAmount:  xrpToTxt-collatral, isLiquidated: false });
       borrowDataByAddress[msg.sender].push(lastBorrowId);
       txtLocked += xrpToTxt-collatral;
    }

    function getClaimAmount(address owner) public view returns(uint){
        uint pendingClaims = lastRewardBatchId-lastClaim[msg.sender]; 
        if(pendingClaims > 30) {
            pendingClaims = 30; // in one go can only claim first 30 days rewards
        }
        uint totalClaimAmount = 0;
        for(uint i=lastClaim[owner]+1; i<=lastClaim[owner]+pendingClaims;i++){
            totalClaimAmount += (lendBalances[owner] / txtToken.balanceOf(address(this)) + txtLocked - txtReward) * rewardClaimAmounts[i];
        }
        return totalClaimAmount;
    }

    function getLentAmountByAddress(address owner) public view returns(uint){
        return lendBalances[owner];
    }

    function getBorrowData(address borrower) public view returns(uint[] memory){
        return borrowDataByAddress[borrower];
    }


    function getBorrowId(uint id) public view returns(borrowData memory){
        return borrowedEntries[id];
    }
    
}