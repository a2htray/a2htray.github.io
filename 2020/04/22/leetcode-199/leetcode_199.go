package main

type TreeNode struct {
	Val int
	Left *TreeNode
	Right *TreeNode
}

func rightSideView(root *TreeNode) []int {
	result := make([]int, 0)
	if root == nil {
		return result
	}
	depthMap := map[int][]int{0: []int{root.Val}}

	stack := []*TreeNode{root}

	for {
		size := len(stack)
		if size == 0 {
			break
		}



	}


	return result
}