pragma solidity ^0.8.0;

contract MergeSortedArrays {
    function merge(int256[] memory arr1, int256[] memory arr2) public pure returns (int256[] memory) {
        int256[] memory merged = new int256[](arr1.length + arr2.length);
        
        uint256 i = 0; // arr1的指针
        uint256 j = 0; // arr2的指针
        uint256 k = 0; // merged数组的指针
        
        while (i < arr1.length && j < arr2.length) {
            if (arr1[i] <= arr2[j]) {
                merged[k] = arr1[i];
                i++;
            } else {
                merged[k] = arr2[j];
                j++;
            }
            k++;
        }
        
        while (i < arr1.length) {
            merged[k] = arr1[i];
            i++;
            k++;
        }
        
        while (j < arr2.length) {
            merged[k] = arr2[j];
            j++;
            k++;
        }
        
        return merged;
    }
}