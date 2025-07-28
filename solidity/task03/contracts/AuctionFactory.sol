// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "./Auction.sol";

contract AuctionFactory is UUPSUpgradeable, OwnableUpgradeable {
    address[] public auctions;
    mapping(address => bool) public isAuction;
    
    address public defaultPaymentToken;
    address public ethUsdPriceFeed;
    mapping(address => address) public tokenPriceFeeds;
    
    event AuctionCreated(address indexed auction, address indexed seller, address nftContract, uint256 nftId);

    function initialize(address _ethUsdPriceFeed) public initializer {
        __Ownable_init();
        ethUsdPriceFeed = _ethUsdPriceFeed;
    }

    function createAuction(
        address _nftContract,
        uint256 _nftId,
        uint256 _duration,
        address _paymentToken
    ) external returns (address) {
        require(_nftContract != address(0), "Invalid NFT contract");
        
        address priceFeed = _paymentToken == address(0) 
            ? ethUsdPriceFeed 
            : tokenPriceFeeds[_paymentToken];
        
        Auction auction = new Auction(
            msg.sender,
            _nftContract,
            _nftId,
            _duration,
            _paymentToken,
            priceFeed
        );
        
        auctions.push(address(auction));
        isAuction[address(auction)] = true;
        
        emit AuctionCreated(address(auction), msg.sender, _nftContract, _nftId);
        return address(auction);
    }

    function getAllAuctions() external view returns (address[] memory) {
        return auctions;
    }
    
    function setPriceFeed(address token, address priceFeed) external onlyOwner {
        tokenPriceFeeds[token] = priceFeed;
    }
    
    function setDefaultPaymentToken(address token) external onlyOwner {
        defaultPaymentToken = token;
    }
    
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}