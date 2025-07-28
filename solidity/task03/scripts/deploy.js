// scripts/deploy.js
const { ethers, upgrades } = require("hardhat");

async function main() {
  // 部署 NFT 合约
  const MyNFT = await ethers.getContractFactory("MyNFT");
  const nft = await MyNFT.deploy();
  await nft.deployed();
  console.log("MyNFT deployed to:", nft.address);

  // 部署工厂合约
  const AuctionFactory = await ethers.getContractFactory("AuctionFactory");
  const factory = await upgrades.deployProxy(AuctionFactory, [
    "0x694AA1769357215DE4FAC081bf1f309aDC325306"
  ]);
  await factory.deployed();
  console.log("AuctionFactory deployed to:", factory.address);
  
  // 设置 ETH 价格源
  await factory.setPriceFeed(
    ethers.constants.AddressZero, 
    "0x694AA1769357215DE4FAC081bf1f309aDC325306"
  );
  console.log("ETH price feed set");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });