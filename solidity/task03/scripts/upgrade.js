const { ethers, upgrades } = require("hardhat");

async function main() {
  // 获取代理合约地址
  const proxyAddress = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512";
  
  // 部署新版本的工厂合约
  const AuctionFactoryV2 = await ethers.getContractFactory("AuctionFactoryV2");
  console.log("准备升级 AuctionFactory 合约...");
  
  // 执行升级
  const upgraded = await upgrades.upgradeProxy(proxyAddress, AuctionFactoryV2);
  await upgraded.deployed();
  
  console.log("AuctionFactory 已成功升级到 V2 版本");
  console.log("新实现合约地址:", await upgrades.erc1967.getImplementationAddress(proxyAddress));
  
  // 验证新功能
  console.log("新版本标识:", await upgraded.version());
  
  // 测试新功能
  const [owner] = await ethers.getSigners();
  await upgraded.blacklistAddress(owner.address, false);
  console.log("黑名单功能测试通过");
}

main()
  .then(() => process.exit(0))
  .catch(error => {
    console.error("升级过程中发生错误:", error);
    process.exit(1);
  });