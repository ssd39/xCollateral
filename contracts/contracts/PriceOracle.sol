//SPDX-License-Identifier: UNLICENSED

pragma solidity ^0.8.20;

interface DoorAccount{
    function getWitnesses() external view returns (address[] memory);
}

contract PriceOracle{
    
    DoorAccount public doorAccount;
    
    string public currency;
    string public currency2;
    uint256 public amount;
    uint256 public amount2;

    constructor(address doorAccount_, string memory currency_, string memory currency2_) {
        doorAccount = DoorAccount(doorAccount_);
        currency = currency_;
        currency2 = currency2_;
    }

    function updateData(uint256 amount_, uint256 amount2_) public {
        address[] memory witnessList = doorAccount.getWitnesses();
        for(uint i = 0; i<witnessList.length;i++){
            if(witnessList[i] == msg.sender){
                amount = amount_;
                amount2 = amount2_;
            }
        }
    }

}