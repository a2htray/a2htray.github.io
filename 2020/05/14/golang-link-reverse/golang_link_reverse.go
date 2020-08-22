package main

import "fmt"

type LNode struct {
	Val int
	Next *LNode
}

func CreateLink(ints []int) *LNode {
	header := &LNode{}
	cur := header
	for _, val := range ints {
		node := &LNode{}
		node.Val = val
		cur.Next = node
		cur = cur.Next
	}

	return header
}

func PrintLink(node *LNode) {
	if node == nil {
		return
	}
	cur := node.Next
	for {
		if cur == nil {
			break
		}
		fmt.Println(cur.Val)
		cur = cur.Next
	}
}

// 实现一，就地逆序
// node *LNode 链表的头结点
func ReverseLinkLocal(node *LNode) {
	if node == nil || node.Next == nil {
		return
	}
	var pre *LNode
	var cur *LNode
	next := node.Next
	for next != nil {
		cur = next.Next
		next.Next = pre
		pre = next
		next = cur
	}
	node.Next = pre
}

// 实现二，递归逆序
// node *LNode 链表的头结点
func reverseLinkRecursive(node *LNode) *LNode {
	if node == nil || node.Next == nil {
		return node
	}
	head := reverseLinkRecursive(node.Next)
	node.Next.Next = node
	node.Next = nil
	return head
}

func ReverseLinkRecursive(node *LNode) {
	first := node.Next
	head := reverseLinkRecursive(first.Next)
	node.Next = head
}

// 实现三，插入逆序
func ReverseLinkInsert(node *LNode) {
	if node == nil || node.Next == nil {
		return
	}

	var cur *LNode
	var next *LNode
	cur = node.Next.Next
	node.Next.Next = nil

	for cur != nil {
		next = cur.Next
		cur.Next = node.Next
		node.Next = cur
		cur = next
	}
}

func main() {
	link := CreateLink([]int{1, 2, 3, 4, 5, 6, 7})
	fmt.Println("before reversing")
	PrintLink(link)

	ReverseLinkLocal(link)
	fmt.Println("after [ReverseLinkLocal] reversing")
	PrintLink(link)

	link = CreateLink([]int{1, 2, 3, 4, 5, 6, 7})
	ReverseLinkRecursive(link)
	fmt.Println("after [ReverseLinkRecursive] reversing")
	PrintLink(link)

	link = CreateLink([]int{1, 2, 3, 4, 5, 6, 7})
	ReverseLinkInsert(link)
	fmt.Println("after [ReverseLinkInsert] reversing")
	PrintLink(link)
}