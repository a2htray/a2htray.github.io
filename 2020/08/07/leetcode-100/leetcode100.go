// LeetCode 100
// https://leetcode-cn.com/problems/same-tree/
package main

import "fmt"

// 晚上的时间刷一道 "每日一题"

// TreeNode ...
type TreeNode struct {
	Val int
	Left *TreeNode
	Right *TreeNode
}

// bfs 通过广度搜索输出各结点的值
func bfs(t *TreeNode) []int {
	values := make([]int, 0)
	queue := make([]*TreeNode, 0)
	queue = append(queue, t)

	for ; len(queue) != 0; {
		p := queue[0]
		values = append(values, p.Val)
		queue = queue[1:]

		if p.Left != nil {
			queue = append(queue, p.Left)
		} else {
			if p.Right != nil {
				values = append(values, 10000000)
			}
		}

		if p.Right != nil {
			queue = append(queue, p.Right)
		}
	}

	return values
}

func isSameTree(p *TreeNode, q *TreeNode) bool {
	pvs := bfs(p)
	qvs := bfs(q)
	fmt.Println(pvs, qvs)

	if len(pvs) != len(qvs) {
		return false
	}

	for i, v := range pvs {
		if v != qvs[i] {
			return false
		}
	}
	return true
}

func main() {
	p := &TreeNode{
			Val: 1,
			Left: &TreeNode{
				Val: 2,
			},
		}
	q := &TreeNode{
		Val: 1,
		Left: &TreeNode{
			Val: 2,
		},
	}

	if isSameTree(p, q) {
		fmt.Println("same")
	} else {
		fmt.Println("not same")
	}
}