pragma solidity ^0.8.0;

contract BinarySearch {
    function binarySearch(int256[] memory arr, int256 target) public pure returns (int256) {
        if (arr.length == 0) {
            return -1;
        }
        
        uint256 left = 0;
        uint256 right = arr.length - 1;
        
        while (left <= right) {
            uint256 mid = left + (right - left) / 2;
            int256 midValue = arr[mid];
            
            if (midValue == target) {
                return int256(mid);
            } else if (midValue < target) {
                left = mid + 1;
            } else {
                if (mid == 0) break;
                right = mid - 1;
            }
        }
        
        return -1;
    }
}