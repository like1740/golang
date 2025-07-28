// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IAuction {
    event BidPlaced(address indexed bidder, uint256 amount);
    event AuctionEnded(address indexed winner, uint256 amount);
    
    function placeBid(uint256 amount) external payable;
    function endAuction() external;
    function getPriceInUSD() external view returns (uint256);
    
    // 状态变量视图函数
    function seller() external view returns (address);
    function nftContract() external view returns (address);
    function nftId() external view returns (uint256);
    function startTime() external view returns (uint256);
    function endTime() external view returns (uint256);
    function highestBid() external view returns (uint256);
    function highestBidder() external view returns (address);
    function ended() external view returns (bool);
    function isEthAuction() external view returns (bool);
    function paymentToken() external view returns (address);
}