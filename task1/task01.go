package main

import (
	"fmt"
	"sort"
)

func main() {
	// 1.只出现一次的数字
	singleNumberResult := singleNumber([]int{1, 2, 2, 3, 3})
	fmt.Printf("只出现一次的数字是: %v\n", singleNumberResult)

	// 2.回文数
	num := 121
	fmt.Printf("%d是否为回文数: %v\n", num, isPalindrome(num))

	// 3.有效的括号
	s1, s2 := "()[]{}", "()[{}]}"
	fmt.Printf("%s是否是有效的括号: %t , %s是否是有效的括号: %t\n", s1, isValidBrackets(s1), s2, isValidBrackets(s2))

	// 4.最长公共前缀
	s3 := []string{"flower", "flow", "flight", "flag"}
	fmt.Printf("%v最长公共前缀是: %s\n", s3, longestCommonPrefix([]string{"flower", "flow", "flight"}))

	// 5.删除排序数组中的重复项
	fmt.Printf("去重后的数组长度是: %d\n", removeDuplicate([]int{1, 2, 3, 3}))

	// 6.加一
	intputArray := []int{1, 2, 3}
	fmt.Printf("加一后输出: %v\n", plusOne(intputArray))

	// 7.合并区间
	intervals := [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
	fmt.Printf("合并区间输入: %v 输出: %v\n", intervals, merge(intervals))

	// 8.两数之和
	nums := []int{1, 2, 3}
	target := 5
	result := twoSum(nums, target)
	fmt.Printf("两数之和: nums = %v, target = %d → %v\n", nums, target, result)
}

func singleNumber(nums []int) int {
	result := 0
	for _, num := range nums {
		result ^= num
	}
	return result
}

func isPalindrome(x int) bool {
	if x < 0 {
		return false
	}
	if x != 0 && x%10 == 0 {
		return false
	}

	reversed := 0
	original := x

	for original > reversed {
		reversed = reversed*10 + original%10
		original /= 10
	}

	return original == reversed || original == reversed/10
}

func isValidBrackets(s string) bool {
	stack := []rune{}

	bracketMap := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
	}

	for _, char := range s {
		switch char {
		case '(', '[', '{':
			stack = append(stack, char)
		case ')', ']', '}':
			if len(stack) == 0 {
				return false
			}

			top := stack[len(stack)-1]

			if top != bracketMap[char] {
				return false
			}

			stack = stack[:len(stack)-1]
		}
	}

	return len(stack) == 0
}

func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}

	base := strs[0]

	for i := 0; i < len(base); i++ {
		for j := 1; j < len(strs); j++ {
			if i >= len(strs[j]) {
				return base[:i]
			}

			if strs[j][i] != base[i] {
				return base[:i]
			}
		}
	}

	return base
}

func removeDuplicate(nums []int) int {
	n := len(nums)
	if n < 2 {
		return n
	}

	slow := 0

	for fast := 1; fast < n; fast++ {
		if nums[fast] != nums[slow] {
			slow++
			nums[slow] = nums[fast]
		}
	}

	return slow + 1
}

func plusOne(digits []int) []int {
	n := len(digits)

	for i := n - 1; i >= 0; i-- {
		digits[i]++

		if digits[i] < 10 {
			return digits
		}

		digits[i] = 0
	}

	return append([]int{1}, digits...)
}

func merge(intervals [][]int) [][]int {
	if len(intervals) <= 1 {
		return intervals
	}

	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	result := [][]int{intervals[0]}

	for i := 1; i < len(intervals); i++ {
		last := result[len(result)-1]
		current := intervals[i]

		if current[0] <= last[1] {
			if current[1] > last[1] {
				last[1] = current[1]
			}
		} else {
			result = append(result, current)
		}
	}

	return result
}

func twoSum(nums []int, target int) []int {
	numMap := make(map[int]int)

	for i, num := range nums {
		complement := target - num

		if idx, found := numMap[complement]; found {
			return []int{idx, i}
		}

		numMap[num] = i
	}

	return nil
}
