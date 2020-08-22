package main

import "fmt"

type Node struct {
	Val string
	Left *Node
	Right *Node
}

type Tree *Node

func PreOrderRecursive(tree Tree) {
	if tree == nil {
		return
	}

	fmt.Println(tree.Val)
	PreOrderRecursive(tree.Left)
	PreOrderRecursive(tree.Right)
}

func PreOrder(tree Tree) {
	if tree == nil {
		return
	}

	stack := []*Node{tree}

	for {
		size := len(stack)
		if size == 0 {
			break
		}
		node := stack[size - 1]
		stack = stack[:size - 1]

		fmt.Println(node.Val)

		if node.Right != nil {
			stack = append(stack, node.Right)
		}

		if node.Left != nil {
			stack = append(stack, node.Left)
		}
	}
}

func PostOrderRecursive(tree Tree)  {
	if tree == nil {
		return
	}

	if tree.Left != nil {
		PostOrderRecursive(tree.Left)
	}

	if tree.Right != nil {
		PostOrderRecursive(tree.Right)
	}

	fmt.Println(tree.Val)
}

func PostOrder(tree Tree) {

}

func InOrderRecursive(tree Tree) {
	if tree == nil {
		return
	}

	if tree.Left != nil {
		InOrderRecursive(tree.Left)
	}

	fmt.Println(tree.Val)

	if tree.Right != nil {
		InOrderRecursive(tree.Right)
	}
}

func InOrder(tree Tree) {
	if tree == nil {
		return
	}

	stack := make([]*Node, 0)

	for {
		size := len(stack)
		if size == 0 && tree == nil{
			break
		}

		for {
			if tree == nil {
				break
			}
			stack = append(stack, tree)
			tree = tree.Left
		}

		size = len(stack)
		node := stack[size - 1]
		stack = stack[:size - 1]

		fmt.Println(node.Val)

		tree = node.Right
	}
}

func main() {
	tree := &Node{
		Val: "A",
		Left: &Node{
			Val: "B",
			Left: &Node{
				Val: "D",
			},
			Right: &Node{
				Val: "F",
				Left: &Node{
					Val: "E",
				},
			},
		},
		Right: &Node{
			Val: "C",
			Left: &Node{
				Val: "G",
				Right: &Node{
					Val: "H",
				},
			},
			Right: &Node{
				Val: "I",
			},
		},
	}

	fmt.Println("Pre")
	PreOrderRecursive(tree)

	fmt.Println("Post")
	PostOrderRecursive(tree)

	fmt.Println("In")
	InOrderRecursive(tree)

	fmt.Println("Pre without recursion")
	PreOrder(tree)

	fmt.Println("Post without recursion")
	PostOrder(tree)

	fmt.Println("In without recursion")
	InOrder(tree)
}

