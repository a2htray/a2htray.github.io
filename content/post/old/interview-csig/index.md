---
title: "CSIG 线上面试"
date: 2022-04-15T11:24:13+08:00
draft: false
reward: false
categories:
 - 生活感想
tags:
 - 面试经
 - 大数除法
images:
 - images/logo-interview-artist-magazine-interview.jpg
---

有幸搞了个 CSIG 的线上面试，感觉是“没什么感觉”，一般般吧，没过。

前面介绍什么就不说了，我这边没突出什么工作亮点，然后就直接共享桌面写代码了。题目是编程实现一个由字符串数组表示的大数的除以 9 的计算，后面又追问了小数点后值如何保存，所以索性在线下实现也写了写。

<!--more-->

其实，对于这种手撕算法题还是挺反感的，有点类似于“形而上”的学习态度，”结伴编程“多少会是有些紧张，没写出来也很正常。但是换位思考一下，问题确实来源于实际，而且看别人码代码总是能看出一些面试者的风格或问题，多少可以作为出题人考查的标准。所以没对没错吧，自己也确实没有准备过算法题，一般般吧。

## 自己的实现

回到这个问题，大数是指那些无法用固定长度类型保存的数值，所以需要用可变长的数组来模拟计算和存储结果。下方代码的实现逻辑比较简单，就是按位对数值进行除以 9 取商取模的操作：

1. 计算第 1 位数值除以 9，取商取模；
2. 计算后续的数值，保存到一个新的 `[]string` 作为结果返回；

```go
package main

import (
	"fmt"
	"strconv"
)

func div(ns []string, prec int) []string {
	result := make([]string, 0)
	// 被除数
	dividend, _ := strconv.ParseInt(ns[0], 10, 64)
	// 模
	var remainder int64
	if dividend < 9 {
		 remainder = dividend
	} else {
		result = append(result, "1")
	}
    
	for _, v := range ns[1:] {
		addition, _ := strconv.ParseInt(v, 10, 64)
		dividend = remainder * 10 + addition
		result = append(result, fmt.Sprintf("%d", dividend / 9))
		remainder = dividend % 9
	}
	if prec > 0 {
		result = append(result, ".")
	}
	for prec > 0 {
		dividend = remainder * 10
		result = append(result, fmt.Sprintf("%d", dividend / 9))
		remainder = dividend % 9
		prec--
	}

	return result
}

func main() {
	op1 := []string{"1", "2", "3", "4", "5"}
	prec := 6
	result := div(op1, prec)
	for _, v := range result {
		fmt.Print(v)
	}
	fmt.Println()

	op2 := []string{"5", "4", "3", "2", "1"}
	prec = 7
	result = div(op2, prec)
	for _, v := range result {
		fmt.Print(v)
	}
	fmt.Println()
}
```

```bash
1371.666666
6035.6666666
```

## math/big 包实现

Go 的 `math/big` 对于大数运算是有实现的，顺带也来看一看。`math/big` 的 `Int` 类型的结构如下：

```go
type Int struct {
	neg bool // 符号
	abs nat  // 整数的绝对值
}

type nat []Word

type Word uint
```

所以 `Int` 类型的底层实现是一个 `uint` 切片，除法运算如下：

```go
func (z *Int) Div(x, y *Int) *Int {
	y_neg := y.neg // z may be an alias for y
	var r Int
	z.QuoRem(x, y, &r) // Step 1
	if r.neg {
		if y_neg {
			z.Add(z, intOne) 
		} else {
			z.Sub(z, intOne)
		}
	}
	return z
}

func (z *Int) QuoRem(x, y, r *Int) (*Int, *Int) {
	z.abs, r.abs = z.abs.div(r.abs, x.abs, y.abs) // Step 2
	z.neg, r.neg = len(z.abs) > 0 && x.neg != y.neg, len(r.abs) > 0 && x.neg // 0 has no sign
	return z, r
}

func (z nat) div(z2, u, v nat) (q, r nat) {
	if len(v) == 0 {
		panic("division by zero")
	}

	if u.cmp(v) < 0 {
		q = z[:0]
		r = z2.set(u)
		return
	}

	if len(v) == 1 {
		var r2 Word
		q, r2 = z.divW(u, v[0])
		r = z2.setWord(r2)
		return
	}

	q, r = z.divLarge(z2, u, v)
	return
}
```

也是带了个余数 `r *Int` 参与计算，除法的计算使用了 [Knuth's Algorithm D](https://surface.syr.edu/cgi/viewcontent.cgi?article=1162&context=eecs_techreports)。

## 总结

后面查资料发现，整数除以 9 是有一定规律的，所以出题人才会这么出，这个确实没接触过，具体计算看这个吧[《任意多位数除以 9：计算规律让你一辈子不忘》](http://www.360doc.com/content/17/0731/20/9930982_675674336.shtml)。

1. 会就会，不会就是不会；
2. 除以 9 的规律，可以知道，但不用去记，这种 tricky 的技巧并不具备普适性；
3. 以后尽量多看看算法题；