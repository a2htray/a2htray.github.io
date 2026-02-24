---
title: "Go 1.18 特性 - 泛型"
date: 2022-04-10T12:11:38+08:00
draft: false
reward: false
categories:
  - 后端技术
  - Go
tags:
 - Go 1.18
images:
 - "images/golang-1.18-release.png"
---

Go 1.18 在 2022 年 3 月 15 日发布，根据团队的[博文](https://go.dev/blog/go1.18)介绍，1.18 版本包含 4 个重要特性：

1. **泛型**；
2. fuzzing；
3. 工作空间；
4. 20% 的性能提升；

<!--more-->

## 泛型

泛型是一种无须关心具体操作类型的编码方式，它将逻辑实现与具体类型解耦，体现在程序中的 3 个地方：

1. 函数和类型的类型参数；
2. 用于指定类型的集合；
3. 类型推断，不需要显式指定类型；

> 本节是官方泛型教程的截取或修改内容，详细请查看[此处](https://go.dev/doc/tutorial/generics)。

### 不同类型求和函数

```go
package main

import "fmt"

// sumInts 计算 int slice 的和
func sumInts(values []int) int {
	var total int
	for _, value := range values {
		total += value
	}
	return total
}

// sumFloat32s 计算 float32 slice 的和
func sumFloat32s(values []float32) float32 {
	var total float32
	for _, value := range values {
		total += value
	}
	return total
}

// sum 计算的 slice 元素可以是 int 类型或 float32 类型
func sum[Element int | float32](values []Element) Element {
	var total Element
	for _, value := range values {
		total += value
	}
	return total
}

func main() {
	intValues := []int{1, 2, 3}
	float32Values := []float32{4, 5, 6}
	fmt.Println(sumInts(intValues))
	fmt.Println(sumFloat32s(float32Values))
    // sum(intValues) 等价于 sum[int](intValues)
	fmt.Println(sum(intValues))
	fmt.Println(sum[int](intValues))
    // sum(float32Values) 等价于 sum[float32](float32Values)
	fmt.Println(sum(float32Values))
	fmt.Println(sum[float32](float32Values))
}
```

### 类型约束

示例 1：

```go
package main

import "fmt"

// Number 类型约束，限制 Number 可以是 int 或 float32
type Number interface {
	int | float32
}

// sum 求和
// [Element Number] 限定 Element 需要符合 Number 类型约束
// 即 Element 只能是 int 或 float32
func sum[Element Number](values []Element) Element {
	var total Element
	for _, value := range values {
		total += value
	}
	return total

}

func main() {
	intValues := []int{1, 2, 3}
	float32Values := []float32{4, 5, 6}
	fmt.Println(sum(intValues))
	fmt.Println(sum(float32Values))
}
```

示例 2：

```go
package main

import "fmt"

// Flag 底层类型为 int
type Flag int

const (
	Flag_A Flag = iota
	Flag_B
	Flag_C
)

// Number 类型约束
// ~ 操作符表示底层类型为 int 的类型也符合
type Number interface {
	float32 | ~int
}

func sum[Element Number](values []Element) Element {
	var total Element
	for _, value := range values {
		total += value
	}
	return total
}

func main() {
	flagValues := []Flag{
		Flag_A, Flag_B, Flag_C,
	}
	float32Values := []float32{4, 5, 6}
	fmt.Println(sum(flagValues))
	fmt.Println(sum(float32Values))
}
```

示例 3：

```go
package main

import (
	"fmt"
	"strconv"
)

type Flag int

func (f Flag) String() string {
	return strconv.Itoa(int(f))
}

const Flag_A Flag = iota

// Number 包含了方法的类型约束
// 指定了 Number 类型的底层类型是 int 并实现了 String() string 方法
type Number interface {
	~int
	String() string
}

func PrintNumber[V Number](number V) {
	fmt.Println(number.String())
}

func main() {
	PrintNumber(Flag_A)
}
```

示例 4：

```go
package main

import "fmt"

type Flag int

const (
	Flag_A Flag = iota
)

// PrintNumber 使用 ~ 操作符，定义类型底层实现
func PrintNumber[Number ~int](value Number) {
	fmt.Println(value)
}

func main() {
	PrintNumber(Flag_A)
    PrintNumber(1)
}
```

## 总结

1. 使用泛型可以减少相同逻辑（不同具体类型）的代码量；
2. 使用 `~` 操作符指定底层类型；
3. 类型约束中可以声明方法；

