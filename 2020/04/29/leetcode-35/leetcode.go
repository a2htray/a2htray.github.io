package main

import "fmt"

func searchInsert(nums []int, target int) int {
	j := 0
	for i, v := range(nums) {
		if v == target {
			return i
		}
		if target > nums[i] {
			j++
		}
	}
	return j
}

func main() {
	fmt.Println(searchInsert([]int{1,3,5,6}, 5))
	fmt.Println(searchInsert([]int{1,3,5,6}, 2))
	fmt.Println(searchInsert([]int{1,3,5,6}, 7))
	fmt.Println(searchInsert([]int{1,3,5,6}, 0))
	
}