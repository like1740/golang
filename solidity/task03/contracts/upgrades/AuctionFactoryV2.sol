// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../AuctionFactory.sol";

contract AuctionFactoryV2 is AuctionFactory {
    // 拍卖统计
    mapping(address => uint256) public auctionsCreated;
    
    // 黑名单
    mapping(address => bool) public blacklisted;
    
    function createAuction(
        address _nftContract,
        uint256 _nftId,
        uint256 _duration,
        address _paymentToken
    ) external override returns (address) {
        require(!blacklisted[msg.sender], "Blacklisted address");
        address auction = super.createAuction(_nftContract, _nftId, _duration, _paymentToken);
        auctionsCreated[msg.sender]++;
        return auction;
    }
    
    function blacklistAddress(address user, bool status) external onlyOwner {
        blacklisted[user] = status;
    }
    
    // 版本标识
    function version() external pure returns (string memory) {
        return "V2";
    }
}