package main

import "fmt"

type ListNode struct {
	Val int
	Next *ListNode
}

// 遍历法
func mergeTwoLists(l1 *ListNode, l2 *ListNode) *ListNode {
	if l1 == nil { return l2 }
	if l2 == nil { return l1 }

	p := &ListNode{}
	q := p
	for {
		if l1 == nil || l2 == nil {
			break
		}

		if l1.Val < l2.Val {
			q.Next = l1
			l1 = l1.Next
		} else {
			q.Next = l2
			l2 = l2.Next
		}
		q = q.Next
	}

	if l1 == nil {
		q.Next = l2
	} else {
		q.Next = l1
	}

	return p.Next
}

// 递归
func recursiveMerge(l1 *ListNode, l2 *ListNode) *ListNode {
	if l1 == nil { return l2 }
	if l2 == nil { return l2 }

	if l1.Val < l2.Val {
		l1.Next = recursiveMerge(l1.Next, l2)
		return l1
	} else {
		l2.Next = recursiveMerge(l2.Next, l1)
		return l2
	}
}

func print(l1 *ListNode) {
	p := l1
	for {
		if p == nil {
			break
		}
		fmt.Println(p.Val)
		p = p.Next
	}
}

func main() {
	fmt.Println("round 1")
	l1 := &ListNode{
		Val: 1,
		Next: &ListNode{
			Val: 2,
			Next: &ListNode{
				Val: 4,
				Next: nil,
			},
		},
	}

	l2 := &ListNode{
		Val: 1,
		Next: &ListNode{
			Val: 3,
			Next: &ListNode{
				Val: 4,
				Next: nil,
			},
		},
	}

	print(recursiveMerge(l1, l2))

	fmt.Println("round 2")
	l3 := &ListNode{
		Val: 2,
	}

	l4 := &ListNode{
		Val: 1,
	}

	print(recursiveMerge(l3, l4))
}