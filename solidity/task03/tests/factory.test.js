const { expect } = require("chai");
const { ethers, upgrades } = require("hardhat");

describe("AuctionFactory", function () {
  let factory;
  let nft;
  let owner, seller, bidder;

  beforeEach(async function () {
    [owner, seller, bidder] = await ethers.getSigners();
    
    // 部署 NFT 合约
    const MyNFT = await ethers.getContractFactory("MyNFT");
    nft = await MyNFT.deploy();
    await nft.deployed();
    
    // 铸造测试 NFT
    await nft.mint(seller.address);
    const tokenId = 1;
    
    // 部署工厂合约
    const AuctionFactory = await ethers.getContractFactory("AuctionFactory");
    factory = await upgrades.deployProxy(AuctionFactory, ["0x694AA1769357215DE4FAC081bf1f309aDC325306"]); 
    await factory.deployed();
    
    // 设置价格源
    await factory.setPriceFeed(
      ethers.constants.AddressZero, 
      "0x694AA1769357215DE4FAC081bf1f309aDC325306"
    );
  });

  it("should create a new auction", async function () {
    const tokenId = 1;
    await nft.connect(seller).approve(factory.address, tokenId);
    
    const tx = await factory.connect(seller).createAuction(
      nft.address,
      tokenId,
      3600, 
      ethers.constants.AddressZero 
    );
    
    const receipt = await tx.wait();
    const auctionCreatedEvent = receipt.events.find(e => e.event === "AuctionCreated");
    const auctionAddress = auctionCreatedEvent.args.auction;
    
    expect(auctionAddress).to.not.equal(ethers.constants.AddressZero);
    expect(await factory.isAuction(auctionAddress)).to.be.true;
  });

  it("should upgrade factory contract", async function () {
    const AuctionFactoryV2 = await ethers.getContractFactory("AuctionFactoryV2");
    const factoryV2 = await upgrades.upgradeProxy(factory.address, AuctionFactoryV2);
    
    // 检查新功能
    expect(await factoryV2.version()).to.equal("V2");
    
    // 创建拍卖测试新功能
    const tokenId = 1;
    await nft.connect(seller).approve(factoryV2.address, tokenId);
    
    await factoryV2.connect(seller).createAuction(
      nft.address,
      tokenId,
      3600,
      ethers.constants.AddressZero
    );
    
    expect(await factoryV2.auctionsCreated(seller.address)).to.equal(1);
  });
});