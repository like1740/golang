const { expect } = require("chai");
const { ethers, upgrades } = require("hardhat");
const { time } = require("@nomicfoundation/hardhat-network-helpers");

describe("NFT拍卖市场", function () {
  let auction;
  let nft;
  let factory;
  let owner, seller, bidder1, bidder2;
  let tokenId;
  
  const DURATION = 3600; // 1小时
  const ETH = ethers.constants.AddressZero;
  const ETH_USD_PRICE_FEED = "0x694AA1769357215DE4FAC081bf1f309aDC325306";

  beforeEach(async function () {
    [owner, seller, bidder1, bidder2] = await ethers.getSigners();
    
    // 部署NFT合约
    const MyNFT = await ethers.getContractFactory("MyNFT");
    nft = await MyNFT.deploy();
    await nft.deployed();
    
    // 铸造NFT
    await nft.mint(seller.address);
    tokenId = 1;
    
    // 部署工厂合约
    const AuctionFactory = await ethers.getContractFactory("AuctionFactory");
    factory = await upgrades.deployProxy(AuctionFactory, [ETH_USD_PRICE_FEED], {
      initializer: "initialize",
    });
    await factory.deployed();
    
    // 设置ETH价格源
    await factory.setPriceFeed(ETH, ETH_USD_PRICE_FEED);
    
    // 卖家批准工厂合约
    await nft.connect(seller).approve(factory.address, tokenId);
    
    // 创建拍卖
    await factory.connect(seller).createAuction(
      nft.address,
      tokenId,
      DURATION,
      ETH
    );
    
    // 获取拍卖合约
    const auctions = await factory.getAllAuctions();
    auction = await ethers.getContractAt("Auction", auctions[0]);
  });

  describe("拍卖创建", function () {
    it("应正确创建拍卖", async function () {
      expect(await auction.seller()).to.equal(seller.address);
      expect(await auction.nftContract()).to.equal(nft.address);
      expect(await auction.nftId()).to.equal(tokenId);
      expect(await auction.startTime()).to.be.gt(0);
      expect(await auction.endTime()).to.equal(await auction.startTime() + DURATION);
      expect(await auction.ended()).to.be.false;
    });
    
    it("应正确记录拍卖地址", async function () {
      const auctions = await factory.getAllAuctions();
      expect(auctions.length).to.equal(1);
      expect(await factory.isAuction(auctions[0])).to.be.true;
    });
  });

  describe("出价功能", function () {
    it("应接受ETH出价", async function () {
      const bidAmount = ethers.utils.parseEther("0.1");
      
      await expect(
        auction.connect(bidder1).placeBid(0, { value: bidAmount })
        .to.emit(auction, "BidPlaced")
        .withArgs(bidder1.address, bidAmount);
      
      expect(await auction.highestBid()).to.equal(bidAmount);
      expect(await auction.highestBidder()).to.equal(bidder1.address);
    });
    
    it("应拒绝低于当前最高出价的出价", async function () {
      const firstBid = ethers.utils.parseEther("0.1");
      const secondBid = ethers.utils.parseEther("0.09");
      
      await auction.connect(bidder1).placeBid(0, { value: firstBid });
      
      await expect(
        auction.connect(bidder2).placeBid(0, { value: secondBid })
        .to.be.revertedWith("Bid too low");
    });
    
    it("应返还前一个最高出价", async function () {
      const firstBid = ethers.utils.parseEther("0.1");
      const secondBid = ethers.utils.parseEther("0.15");
      
      // 第一次出价
      await auction.connect(bidder1).placeBid(0, { value: firstBid });
      
      // 检查余额变化
      await expect(
        auction.connect(bidder2).placeBid(0, { value: secondBid })
        .to.changeEtherBalance(bidder1, firstBid);
      
      expect(await auction.highestBidder()).to.equal(bidder2.address);
    });
    
    it("应拒绝卖家出价", async function () {
      await expect(
        auction.connect(seller).placeBid(0, { value: ethers.utils.parseEther("0.1") })
        .to.be.revertedWith("Seller cannot bid");
    });
  });

  describe("结束拍卖", function () {
    it("应正确结束拍卖并转移资产", async function () {
      const bidAmount = ethers.utils.parseEther("0.5");
      
      // 出价
      await auction.connect(bidder1).placeBid(0, { value: bidAmount });
      
      // 快进到拍卖结束
      await time.increase(DURATION + 1);
      
      // 结束拍卖
      await expect(auction.connect(seller).endAuction())
        .to.emit(auction, "AuctionEnded")
        .withArgs(bidder1.address, bidAmount);
      
      // 检查NFT所有权
      expect(await nft.ownerOf(tokenId)).to.equal(bidder1.address);
      
      // 检查资金转移
      const sellerBalanceBefore = await ethers.provider.getBalance(seller.address);
      await auction.connect(seller).endAuction();
      const sellerBalanceAfter = await ethers.provider.getBalance(seller.address);
      expect(sellerBalanceAfter.sub(sellerBalanceBefore)).to.be.closeTo(
        bidAmount,
        ethers.utils.parseEther("0.001")
      );
    });
    
    it("应拒绝提前结束拍卖", async function () {
      await expect(
        auction.connect(seller).endAuction())
        .to.be.revertedWith("Auction not ended");
    });
    
    it("应处理无人出价的情况", async function () {
      // 快进到拍卖结束
      await time.increase(DURATION + 1);
      
      // 结束拍卖
      await auction.connect(seller).endAuction();
      
      // 检查NFT所有权（应返回卖家）
      expect(await nft.ownerOf(tokenId)).to.equal(seller.address);
    });
  });

  describe("预言机集成", function () {
    it("应获取美元价格", async function () {
      const bidAmount = ethers.utils.parseEther("1");
      await auction.connect(bidder1).placeBid(0, { value: bidAmount });
      
      const usdValue = await auction.getPriceInUSD();
      
      expect(usdValue).to.be.gt(0);
      
      console.log(`1 ETH ≈ $${usdValue / 100000000} (基于Chainlink预言机)`);
    });
  });
});