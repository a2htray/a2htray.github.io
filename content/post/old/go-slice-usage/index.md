---
title: "Go Slice 的使用"
date: 2022-09-19T23:33:58+08:00
draft: false
reward: false
categories:
  - 后端技术
  - Go
tags:
 - slice
 - 翻译
---

原文：[Slices/arrays explained: create, index, slice, iterate](https://yourbasic.org/golang/slices-explained/)

<!--more-->

## 概述

slice 是 Go 中的一种数据结构，描述了底层实现 array 的部分信息，并不会真实地存储任何数据。

* 修改 slice 中的元素，实际上就是修改底层实现 array 中的元素，引用相同 array 的 slice 也会相应被修改；

```go
func changeSliceElement() {
    s := []int{1, 2, 3}
    s1 := s
    s2 := s
    s[0] = 0

    fmt.Println(s1)
    fmt.Println(s2)
}
```

```bash
[0 2 3]
[0 2 3]
```

* slice 会进行扩容和缩容；

```go
func growSlice() {
    s := []int{0, 1, 2, 3}

    fmt.Printf("slice 的长度为 %d，容量为 %d\n", len(s), cap(s))
    for i := 5; i < 10; i++ {
        s = append(s, i)
        fmt.Printf("append %d 次，slice 的长度为 %d，容量为 %d\n", i-4, len(s), cap(s))
    }
}
```

```bash
slice 的长度为 4，容量为 4
append 1 次，slice 的长度为 5，容量为 8
append 2 次，slice 的长度为 6，容量为 8
append 3 次，slice 的长度为 7，容量为 8
append 4 次，slice 的长度为 8，容量为 8
append 5 次，slice 的长度为 9，容量为 16
```

在第 5 次 append 后，slice 进行了扩容（以 2 为倍数）。

不同于 C 中的 realloc 方法，Go 中的 slice 没有直接的缩容方法，参考【1】中的给出了间接的方法示例。

```go
func shrinkSlice() {
    s := []int{0, 1, 2, 3}
    fmt.Printf("slice 的长度为 %d，容量为 %d\n", len(s), cap(s))

    s1 := append([]int{}, s[:2]...)
    fmt.Printf("slice 的长度为 %d，容量为 %d\n", len(s1), cap(s1))
}
```

```bash
slice 的长度为 4，容量为 4
slice 的长度为 2，容量为 2
```

* slice 的元素访问；

## 定义

```go
var s []int                   // a nil slice
s1 := []string{"foo", "bar"}
s2 := make([]int, 2)          // same as []int{0, 0}
s3 := make([]int, 2, 4)       // same as new([4]int)[:2]
fmt.Println(len(s3), cap(s3)) // 2 4
```

* slice 的零值为 nil，内置方法 len、cap 和 append 都会把 nil 视为一个容量为 0 的 slice；

```go
func nilSlice() {
    var s []int
    if s == nil {
        fmt.Printf("var s []int 声明了一个 []int 的零值 nil\n")
    }
    fmt.Printf("len(s) = %d, cap(s) = %d\n", len(s), cap(s))
}
```

```bash
var s []int 声明了一个 []int 的零值 nil
len(s) = 0, cap(s) = 0
```

* 内置方法 len 和 cap 分别可以得到 slice 的长度和容量；

## 切片

```go
a := [...]int{0, 1, 2, 3} // an array
s := a[1:3]               // s == []int{1, 2}        cap(s) == 3
s = a[:2]                 // s == []int{0, 1}        cap(s) == 4
s = a[2:]                 // s == []int{2, 3}        cap(s) == 2
s = a[:]                  // s == []int{0, 1, 2, 3}  cap(s) == 4
```

支持从一个 slice 创建出新的 slice：

* 通过指定索引值区间 s[low:high] 创建一个 slice；

```go
func Index() {
   s := []int{1, 2, 3, 4, 5}
   fmt.Printf("s 的长度为 %d，容量为 %d\n", len(s), cap(s))

   s1 := s[:2]
   fmt.Printf("s1 = s[:2] 的长度为 %d，容量为 %d\n", len(s1), cap(s1))

   // 打印出 s 和 s1 的地址
   fmt.Printf("s 的地址 %p，s1 的地址 %p\n", &s, &s1)
   // s 和 s1 的第 1 个元素是否在同一个地址上？
   if &s == &s1 {
      fmt.Println("s 和 s1 的第 1 个元素在同一个地址上")
   } else {
      fmt.Println("s 和 s1 的第 1 个元素不在同一个地址上")
   }

   // low 和 high 的取值范围
   //  创建 1 个长度为 4，容量为 8 的 slice
   s2 := make([]int, 4, 8)
   s2[0] = 1
   s2[1] = 2
   s2[2] = 3
   s2[3] = 4
   fmt.Printf("s2 的长度为 %d，容量为 %d\n", len(s2), cap(s2))

   s3 := s2[0:]
   fmt.Printf("s3 的长度为 %d，容量为 %d\n", len(s3), cap(s3))
   s4 := s2[0:4]
   fmt.Printf("s4 的长度为 %d，容量为 %d\n", len(s4), cap(s4))
   s5 := s2[0:5]
   fmt.Printf("s5 的长度为 %d，容量为 %d\n", len(s5), cap(s5))
   s6 := s2[0:8]
   fmt.Printf("s6 的长度为 %d，容量为 %d\n", len(s6), cap(s6))
   fmt.Printf("第 8 个元素为 %d\n", s6[7])
}
```

```bash
s 的长度为 5，容量为 5
s1 = s[:2] 的长度为 2，容量为 5
s 的地址 0xc0000b4090，s1 的地址 0xc0000b40a8
s 和 s1 的第 1 个元素不在同一个地址上
s2 的长度为 4，容量为 8
s3 的长度为 4，容量为 8
s4 的长度为 4，容量为 8
s5 的长度为 5，容量为 8
s6 的长度为 8，容量为 8
第 8 个元素为 0
```

从上述示例中，可以得出：

1. 通过 s[low:high] 创建的 slice 的内存地址与 s 的内存地址不同；
2. high 的最大值为 slice 的容量，若容量大于长度时l，high 的值可能会大于长度，此时通过切片会得到一个以元素零值填充的 slice；

* slice 中的元素为引用类型时，通过切片得到的新的 slice 中的元素保持相同的引用；

```go
func IndexRefObjs() {
   s := make([]*strings.Reader, 4, 8)
   for i, char := range []string{"A", "B", "C", "D"} {
      s[i] = strings.NewReader(char)
   }

   s1 := s[0:4]
   if s[0] == s1[0] {
      fmt.Println("s 的第 1 个元素和 s1 的第 1 个元素为相同引用")
   }
}
```

```bash
s 的第 1 个元素和 s1 的第 1 个元素为相同引用
```

## 迭代

```go
s := []string{"Foo", "Bar"}
for i, v := range s {
    fmt.Println(i, v)
}
```

```bash
0 Foo
1 Bar
```

* range 表达式用于迭代 slice；
* 两个迭代值 i 和 v 分别为索引值和元素值；
* 第 2 个迭代值是可选的；
* 如果 slice 为 nil，则其可迭代值为 0；

```go
func IterNilSlice() {
   var s []int
   if s == nil {
      fmt.Println("s 值为 nil")
   }

   for i, v := range s {
      fmt.Println(i, v)
   }
}
```

```bash
s 值为 nil
```

上述示例可知，若一个 slice 为 nil 时，仍然可以作为 range 表达式的操作值。

## 添加和复制

* append 函数可以将一个元素添加到 slice 的尾部，如果超过了 slice 的容量会进行自动扩容；
* copy 函数可以将源 slice 中的元素复制到目标 slice 中，可复制的元素个数为源 slice 和目标 slice 长度中较小值；

## 参考

1. [Does Go have no real way to shrink a slice? Is that an issue?](https://stackoverflow.com/questions/16748330/does-go-have-no-real-way-to-shrink-a-slice-is-that-an-issue)
