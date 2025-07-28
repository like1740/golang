const { ethers } = require("hardhat");
const { parseEther } = ethers.utils;

async function main() {
  // 获取测试账户
  const [owner, seller, bidder1, bidder2] = await ethers.getSigners();
  
  console.log("使用账户:");
  console.log("所有者:", owner.address);
  console.log("卖家:", seller.address);
  console.log("出价者1:", bidder1.address);
  console.log("出价者2:", bidder2.address);
  
  // 获取已部署的合约
  const MyNFT = await ethers.getContractFactory("MyNFT");
  const nft = await MyNFT.attach("0x5FbDB2315678afecb367f032d93F642f64180aa3");
  
  const AuctionFactory = await ethers.getContractFactory("AuctionFactory");
  const factory = await AuctionFactory.attach("0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512");
  
  // 1. 铸造NFT
  console.log("\n步骤1: 卖家铸造NFT...");
  await nft.connect(seller).mint(seller.address);
  const tokenId = 1;
  console.log(`NFT #${tokenId} 已铸造给 ${seller.address}`);
  
  // 2. 卖家批准工厂合约转移NFT
  console.log("\n步骤2: 卖家批准工厂合约...");
  await nft.connect(seller).approve(factory.address, tokenId);
  console.log("工厂合约已获得NFT转移权限");
  
  // 3. 创建拍卖
  console.log("\n步骤3: 创建拍卖...");
  const duration = 300;
  const tx = await factory.connect(seller).createAuction(
    nft.address,
    tokenId,
    duration,
    ethers.constants.AddressZero 
  );
  
  const receipt = await tx.wait();
  const auctionCreatedEvent = receipt.events.find(e => e.event === "AuctionCreated");
  const auctionAddress = auctionCreatedEvent.args.auction;
  
  console.log(`拍卖已创建: ${auctionAddress}`);
  console.log(`NFT合约: ${nft.address}, Token ID: ${tokenId}`);
  console.log(`持续时间: ${duration} 秒`);
  
  // 4. 出价
  const Auction = await ethers.getContractFactory("Auction");
  const auction = await Auction.attach(auctionAddress);
  
  console.log("\n步骤4: 出价者1出价...");
  const bid1 = parseEther("0.1");
  await auction.connect(bidder1).placeBid(0, { value: bid1 });
  console.log(`出价者1出价: ${ethers.utils.formatEther(bid1)} ETH`);
  
  console.log("\n步骤5: 出价者2出价...");
  const bid2 = parseEther("0.15");
  await auction.connect(bidder2).placeBid(0, { value: bid2 });
  console.log(`出价者2出价: ${ethers.utils.formatEther(bid2)} ETH`);
  
  // 5. 等待拍卖结束
  console.log("\n等待拍卖结束...");
  await new Promise(resolve => setTimeout(resolve, (duration + 10) * 1000));
  
  // 6. 结束拍卖
  console.log("\n步骤6: 结束拍卖...");
  await auction.connect(seller).endAuction();
  
  // 7. 验证结果
  console.log("\n拍卖结果:");
  console.log("获胜者:", await auction.highestBidder());
  console.log("最高出价:", ethers.utils.formatEther(await auction.highestBid()), "ETH");
  
  const usdValue = await auction.getPriceInUSD();
  console.log("美元价值:", usdValue.toString(), "USD (基于Chainlink预言机)");
  
  // 验证NFT所有权
  const newOwner = await nft.ownerOf(tokenId);
  console.log(`NFT新所有者: ${newOwner} (应该是获胜者: ${await auction.highestBidder()})`);
}

main()
  .then(() => process.exit(0))
  .catch(error => {
    console.error("测试过程中发生错误:", error);
    process.exit(1);
  });