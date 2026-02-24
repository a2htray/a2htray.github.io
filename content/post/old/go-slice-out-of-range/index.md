---
title: "Slice 什么时候报 out of range"
date: 2022-04-26T09:24:01+08:00
draft: false
reward: false
categories:
  - 后端技术
  - Go
tags:
 - slice
 - stackoverflow
 - 面试经
images:
 - images/go-slice.png
---

面试的时候问到了一个关于 go Slice 的问题，即为什么在 `a[i:]` 中 `i` 的取值可以是 `a` 的长度。平时开发中也是这么用的，但没太深入的了解，所以在这篇文章中对其进行一些探讨。

<!--more-->

## slice 删除元素

Go 标准内置包没有提供太多操作 slice 的方法，所以如果要删除 slice 的元素通常都能找到以下的实现。

```go
func remove(slice []int, s int) []int {
    return append(slice[:s], slice[s+1:]...)
}
```

这里就引发出一个疑问：当要删除最后一个元素时，`s+1` 不就等于 slice 的长度了，但程序为什么没有报 `index out of range` 错误。

## 两种表达式的官方解释

官方对于 `a[i]` 和 `a[low:high]`有不同的定义，分别为索引表达式 和 slice 表示式，所以两者并非一个东西。

**a[i] 表达式**

> 下标的取值范围在 $0 \le i \lt len(a)$，如果超出了会报 `out of range` 运行时错误。

**a[low:high] 表示式**

>  对于数组或字符串，下标的取值范围在 $0 \le low \le high \le len(a)$，即下标可以取到数组的长度或字符串的长度。对于 slice 来说，下标的取值上限是 slice 的容量，显然 slice 的容量会大于等于其长度。

```go
package main

import (
	"fmt"
)

func main() {
	ints := make([]int, 0)
	for i := 0; i < 3; i++ {
		ints = append(ints, i)
	}

	fmt.Printf("length: %d, capacity: %d\n", len(ints), cap(ints))
	fmt.Println(ints[3:4])
}
```

```bash
length: 3, capacity: 4
[0]
```

`ints` 的容量为 4，所以 `ints[3:4]` 符合定义。另外，扩容操作只有在使用 `append` 方法后才会执行，正常初始化的 slice 的长度与容量相同，如下：

```go
package main

import "fmt"

func main() {
	ints := []int{1, 2, 3}
	fmt.Printf("length: %d, capacity: %d\n", len(ints), cap(ints))
}
```

```go
length: 3, capacity: 3
```

值得注意的是，如果 `a[low:high]` 中的 `high` 缺省了会默认为 `a` 的长度，则 `a[len(a):]` 会变成 `a[len(a):len(a)]` ，而 `len(a) - len(a) = 0`，所以会取到一个空 `[]` 的 slice。

## 猜想一

猜想：`a[len(a):]` 是不是读到了 slice 相邻内存上的数据，因为相邻内存上没有数据，所以才会返回 `[]`。所以要先了解下 slice 容量的扩展方式，例子如下：

```go
package main

import (
	"fmt"
)

func main() {
	ints := make([]int, 0)
	for i := 0; i < 9; i++ {
		// sh := (*reflect.SliceHeader)(unsafe.Pointer(&ints))
		fmt.Printf("before adding a element [%d]\n", i)
		fmt.Printf("length: %d, capacity: %d\n", len(ints), cap(ints))

		ints = append(ints, i)
		fmt.Printf("after adding a element [%d]\n", i)
		fmt.Printf("length: %d, capacity: %d\n", len(ints), cap(ints))
		fmt.Println()
		// slice 容量从 4 开始以 2 的倍数增加
	}
}
```

```bash
before adding a element [0]
length: 0, capacity: 0
after adding a element [0]
length: 1, capacity: 1

before adding a element [1]
length: 1, capacity: 1
after adding a element [1]
length: 2, capacity: 2

before adding a element [2]
length: 2, capacity: 2
after adding a element [2]
length: 3, capacity: 4

before adding a element [3]
length: 3, capacity: 4
after adding a element [3]
length: 4, capacity: 4

before adding a element [4]
length: 4, capacity: 4
after adding a element [4]
length: 5, capacity: 8

before adding a element [5]
length: 5, capacity: 8
after adding a element [5]
length: 6, capacity: 8

before adding a element [6]
length: 6, capacity: 8
after adding a element [6]
length: 7, capacity: 8

before adding a element [7]
length: 7, capacity: 8
after adding a element [7]
length: 8, capacity: 8

before adding a element [8]
length: 8, capacity: 8
after adding a element [8]
length: 9, capacity: 16
```

从输出可以得知：

1. 添加 `int` 0 后，长度为 1，容量为 1；
2. 添加 `int` 1 后，长度为 2，容量为 2；
3. 添加 `int` 2 后，长度为 3，容量为 4；
4. 添加 `int` 4 后，长度为 5，容量为 8；
5. 添加 `int` 8 后，长度为 9，容量为 16；

综上，slice 底层数组的扩展规则为**容量以 2 的倍数增长**。

```go
package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	ints := make([]int, 0)
	ints = append(ints, 1, 2, 3)
	// slice 的长度为 3，容量为 3，取下标 3 的元素会越界
	// fmt.Println("try to get the 4th element", ints[3])

	// 取 slice 从下标 3 之后的元素，返回的是个 []
	fmt.Println("try to get the rest elements from the 4th element: ", ints[3:])

	// sh := (*reflect.SliceHeader)(unsafe.Pointer(&ints))
	// 打印 ints 的起始地址
	fmt.Printf("the start address of ints: %p\n", ints)
	// 第一次：打印 ints[3:] 的起始地址
	fmt.Printf("first the start address of ints[3:]: %p\n", ints[3:])
	// 第二次：打印 ints[3:] 的起始地址
	fmt.Printf("second the start address of ints[3:]: %p\n", ints[3:])
	// 打印 ints[2] 的地址
	fmt.Printf("the address of ints[2]: %p\n", &(ints[2]))

	fmt.Println()
	rest := ints[3:]
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&rest))
	fmt.Printf("length: %d, capacity: %d\n", sh.Len, sh.Cap)
}
```

```bash
try to get the rest elements from the 4th element:  []
the start address of ints: 0xc0000be090
first the start address of ints[3:]: 0xc0000be090
second the start address of ints[3:]: 0xc0000be090
the address of ints[2]: 0xc0000be0a0
```

`ints` 初始化为长度为 3、容量为 4 的 slice，对其进行取值操作，得到：

1. `ints[3:]` 返回一个空的 slice `[]`，同时其起始地址与 `ints` 起始地址相同；
2. 多次调用 `ints[3:]` 始终返回相同的结果，且长度和容量均为 0；

显然，`ints[3:]` 并非取到了 `ints` 相邻内存中的值，所以猜想不成立。

## 猜想二

猜想：当 `len(a) < cap(a)` 时，`a[i:j] (i <= j, len(a) < j)` 取到了已分配内存中的零值。

```go
package main

import (
	"fmt"
)

func main() {
	ints := make([]int, 0)
	for i := 0; i < 5; i++ {
		ints = append(ints, i)
	}

	fmt.Printf("length: %d, capacity: %d\n", len(ints), cap(ints))
	fmt.Println(ints[4:8])
}
```

 ```bash
length: 5, capacity: 8
[4 0 0 0]
 ```

在上述代码中，`ints` 的长度为 5、容量为 8，`ints[4:8]` 中的下标值 8 符合小于等于容量的规定，语法有效。同时看到输出，也确实取到了已分配（未使用）内存中的值。

## 参考

1. [stackoverflow: Why go doesn't report "slice bounds out of range" in this case?](https://stackoverflow.com/questions/52250700/why-go-doesnt-report-slice-bounds-out-of-range-in-this-case)

2. [stackoverflow: Why does go allow slicing from len(slice)?](https://stackoverflow.com/questions/30208182/why-does-go-allow-slicing-from-lenslice)
3. [官方 Index 表示式](https://go.dev/ref/spec#Index_expressions)
4. [官方 Slice 表达式](https://go.dev/ref/spec#Slice_expressions)
5. [stackoverflow: Why does go allow slicing from len(slice)?](https://stackoverflow.com/questions/30208182/why-does-go-allow-slicing-from-lenslice)

