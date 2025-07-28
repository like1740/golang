// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@chainlink/contracts/src/v0.8/interfaces/AggregatorV3Interface.sol";

contract Auction is ReentrancyGuard {
    address public seller;
    address public factory;
    IERC721 public nftContract;
    uint256 public nftId;
    
    uint256 public startTime;
    uint256 public endTime;
    uint256 public highestBid;
    address public highestBidder;
    bool public ended;
    
    // 代币出价
    IERC20 public paymentToken;
    bool public isEthAuction;
    
    // Chainlink
    AggregatorV3Interface internal priceFeed;
    address public chainlinkPriceFeed;

    event BidPlaced(address bidder, uint256 amount);
    event AuctionEnded(address winner, uint256 amount);

    modifier onlySeller() {
        require(msg.sender == seller, "Not seller");
        _;
    }

    modifier onlyFactory() {
        require(msg.sender == factory, "Not factory");
        _;
    }

    constructor(
        address _seller,
        address _nftContract,
        uint256 _nftId,
        uint256 _duration,
        address _paymentToken,
        address _priceFeed
    ) {
        factory = msg.sender;
        seller = _seller;
        nftContract = IERC721(_nftContract);
        nftId = _nftId;
        startTime = block.timestamp;
        endTime = startTime + _duration;
        
        if (_paymentToken == address(0)) {
            isEthAuction = true;
        } else {
            paymentToken = IERC20(_paymentToken);
            isEthAuction = false;
        }
        
        chainlinkPriceFeed = _priceFeed;
        if (_priceFeed != address(0)) {
            priceFeed = AggregatorV3Interface(_priceFeed);
        }
    }

    function placeBid(uint256 amount) external payable nonReentrant {
        require(block.timestamp >= startTime && block.timestamp <= endTime, "Auction not active");
        require(!ended, "Auction ended");
        require(msg.sender != seller, "Seller cannot bid");
        
        uint256 bidAmount = isEthAuction ? msg.value : amount;
        require(bidAmount > highestBid, "Bid too low");
        
        // 返还前一个最高出价
        if (highestBidder != address(0)) {
            if (isEthAuction) {
                payable(highestBidder).transfer(highestBid);
            } else {
                paymentToken.transfer(highestBidder, highestBid);
            }
        }
        
        // 接受新出价
        if (isEthAuction) {
            require(msg.value == bidAmount, "ETH amount mismatch");
        } else {
            require(paymentToken.transferFrom(msg.sender, address(this), bidAmount), "Token transfer failed");
        }
        
        highestBid = bidAmount;
        highestBidder = msg.sender;
        
        emit BidPlaced(msg.sender, bidAmount);
    }

    function endAuction() external nonReentrant {
        require(block.timestamp > endTime || msg.sender == factory, "Auction not ended");
        require(!ended, "Auction already ended");
        
        ended = true;
        
        if (highestBidder != address(0)) {
            // 转移 NFT
            nftContract.transferFrom(seller, highestBidder, nftId);
            
            // 转移资金
            if (isEthAuction) {
                payable(seller).transfer(highestBid);
            } else {
                paymentToken.transfer(seller, highestBid);
            }
        } else {
            // 如果没有出价，NFT 返回给卖家
            nftContract.transferFrom(address(this), seller, nftId);
        }
        
        emit AuctionEnded(highestBidder, highestBid);
    }
    
    function getPriceInUSD() public view returns (uint256) {
        if (chainlinkPriceFeed == address(0)) return 0;
        
        (, int256 price, , , ) = priceFeed.latestRoundData();
        uint8 decimals = priceFeed.decimals();
        
        // 计算美元价值 = (出价金额 * 价格) / 10^decimals
        return (highestBid * uint256(price)) / (10 ** decimals);
    }
    
    // 工厂合约可以提前结束拍卖
    function forceEnd() external onlyFactory {
        endAuction();
    }
}