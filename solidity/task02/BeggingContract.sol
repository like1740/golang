// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract BeggingContract {
    // 合约所有者地址
    address public owner;
    
    // 记录每个地址的捐赠金额
    mapping(address => uint256) public donations;
    
    // 合约总捐赠金额
    uint256 public totalDonations;
    
    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }
    
    constructor() {
        owner = msg.sender;
    }
    
    // 捐赠函数
    function donate() external payable {
        require(msg.value > 0, "Donation amount must be greater than 0");
        
        // 更新捐赠记录
        donations[msg.sender] += msg.value;
        totalDonations += msg.value;
    }
    
    // 提取合约所有资金
    function withdraw() external onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No funds available to withdraw");
        
        // 转账给所有者
        (bool success, ) = payable(owner).call{value: balance}("");
        require(success, "Transfer failed");
    }
    
    // 查询指定地址的捐赠金额
    function getDonation(address donor) external view returns (uint256) {
        return donations[donor];
    }
    
    // 回退函数
    receive() external payable {
        donations[msg.sender] += msg.value;
        totalDonations += msg.value;
    }
}