---
title: "channel 的方向"
date: 2022-05-05T17:17:29+08:00
draft: false
reward: false
categories:
  - 后端技术
  - Go
tags:
 - channel
---

在回忆管道方向的语法上时不时地会出错，所以搜罗一些资料以加强自身的记忆。

<!--more-->

## 基础

channel 是 Go 提供**同步的**、**强类型的**消息传输功能的一种数据结构，搭配上 goroutine，构建了 Go 的 Communicating Sequential Processes（CSP）并发模型。

> channel 和 goroutine 是 Go 的 CSP 并发模型的基石。

定义一个 channel 只需要一个 `chan` 关键字及 channel 中传输的元素的类型，如：

```go
var c chan int
// 或
c := make(chan int)
```

默认情况下，channel 的双向的，即 channel 的一端可读、另一端可写。因为没有定义 channel 的方向，在编写程序的过程中，很有可以向一个已关闭的 channel 发送信息。比较优的编码实践应该在代码中指定 channel 的方向。

## channel 的方向

操作符 `<-` 用于指定 channel 是方向，表示读或写的操作。由此可知，声明一个 channel 有 3 种方式：

```go
var c chan int // 双向
var c <-chan int // 只能读
var c chan<- int // 只能写
```

记忆方式：

```go
<-chan // 从 channel 中来 
chan<- // 到 channel 中去
```

值得注意的是，一开始的 channel 声明可以是双向的，并且可以作为参数的过程中指定 channel 的方法，这样就可以避免在读取 channel 的过程中关闭 channel（因为只有写 channel 的代码才有权限关闭 channel）。

> 只读的 channel 不能被关闭。

在函数签名中，也应该使用带有方向的 channel 作为参数类型或返回值类型，除非在某一函数中完全使用双向的 channel。

```go
func ReturnReadOnly() <-chan int {
    c := make(chan int)
    // ...
    return c
}

// 因为 ReturnReadOnly 返回的是只读的 channel，所以 readOnly 是只读的
// 但在 ReturnReadOnly 函数中声明的是一个双向的 channel `c := make(chan int)`
readOnly := ReturnReadOnly() 
```

使用带方向的 channel 的好处：

1. 事先声明 channel 的读写性质；
2. 编写的 channel 相关代码语义明显；

## 参考

1. [Directional Channels in Go](https://blog.gopheracademy.com/advent-2019/directional-channels/#:~:text=Directional%20Channels%20in%20Go%20Go%E2%80%99s%20channels%20provide%20a,the%20backbone%20of%20Go%E2%80%99s%20CSP%20-inspired%20concurrency%20model.)