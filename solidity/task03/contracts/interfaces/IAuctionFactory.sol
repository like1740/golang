// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IAuctionFactory {
    event AuctionCreated(
        address indexed auction,
        address indexed seller,
        address nftContract,
        uint256 nftId
    );
    
    function createAuction(
        address _nftContract,
        uint256 _nftId,
        uint256 _duration,
        address _paymentToken
    ) external returns (address);
    
    function getAllAuctions() external view returns (address[] memory);
    function setPriceFeed(address token, address priceFeed) external;
    function setDefaultPaymentToken(address token) external;
    
    // 状态变量视图函数
    function auctions(uint256 index) external view returns (address);
    function isAuction(address auction) external view returns (bool);
    function defaultPaymentToken() external view returns (address);
    function ethUsdPriceFeed() external view returns (address);
    function tokenPriceFeeds(address token) external view returns (address);
}