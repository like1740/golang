pragma solidity ^0.8.0;

contract RomanToInteger {
    mapping(bytes1 => uint256) private romanValues;
    
    constructor() {
        // 初始化罗马数字对应值
        romanValues['I'] = 1;
        romanValues['V'] = 5;
        romanValues['X'] = 10;
        romanValues['L'] = 50;
        romanValues['C'] = 100;
        romanValues['D'] = 500;
        romanValues['M'] = 1000;
    }
    
    // 罗马数字转整数函数
    function romanToInt(string memory s) public view returns (uint256) {
        bytes memory roman = bytes(s);
        uint256 length = roman.length;
        uint256 total = 0;
        
        for (int256 i = int256(length) - 1; i >= 0; i--) {
            uint256 currentValue = romanValues[roman[uint256(i)]];
            
            if (i < int256(length) - 1 && currentValue < romanValues[roman[uint256(i + 1)]]) {
                total -= currentValue;
            } else {
                total += currentValue;
            }
        }
        
        return total;
    }
}