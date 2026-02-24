---
title: "Note: Evaluating the Go Programming Language with Design Patterns"
date: 2023-04-18T20:43:16+08:00
draft: true
reward: false
categories:
 - 后端技术
 - Go
tags:
 - Design Pattern
 - Reading Note
---

《Evaluating the Go Programming Language with Design Patterns》是 Frank Schmager 编写的一本关于 Go 设计模式的教程书，本人无意间在网上得知，正逢想多了解一些 Go 设计模式的相关知识，故将书中有用的句子摘录到本文。本文的书写规范如下：



<!--more-->

## Background

### Go

Go is an object-oriented programming language with C-like syntax.

> Go 是一种类 C 语法的面向对象编程语言。

Go supports **a mixture of static and dynamic typing**, and is designed to be safe and efficient.

> Go 同时支持静态类型检查（强类型语言）和动态类型编程（interface），旨在编写出安全和高效的程序。

Communicating Sequential Processes (CSP) was developed by C. A. R. Hoare in 1978. Hoare **introduces the idea of using channels for interprocess communication**.

> CSP：引入了使用通道进行通信的概念。

Go's interprocess communication is largely influenced by CSP and its successors as outlined above. Go uses typed channels which can be buffered or unbuffered; Go's `select` statement implements the process synchronization mechanism described by Hoare and implemented in Newsqueak.

> Go 的进程间通信机制深受 CSP 的影响，如 Go 使用限定类型的缓冲或非缓冲的通道、Go 的 select 语句实现了 Hoare 所描述的进程同步机制。

Go is systems programming language with focus on concurrency.

> Go 语言更关注于并发的处理。

Goroutines are lightweight and **multiplexed into threads**.

> 协程非常轻量，并且可以在线程中被调度。

Go distinguishes between methods and functions. Methods are the functions that have a receiver, and the receiver can be any valid variable name.

> Go 的方法和函数是不等同的，Go 的方法是一类包含接收者的函数，同时接收者可以是任何有效的变量名。

Following C++, rather than Java, methods in Go are defined outside struct declarations.

> Go 的方法定义在结构之外。如：
>
> ```go
> type Car struct {}
>
> func (c *Car) Drive() {}
> ```

Go has no class-based inheritance, code reuse is supported by embedding.

> Go 没有类继承机制，而是通过内嵌的方式实现代码复用。如：
>
> ```go
> type Truck struct {
> 	*Car
> }
>
> func (t *Truck) Attack() {}
> ```

If a method or field can not be found in an object's type definition, then embedded types are searched, and method calls are forward to the embedded object.

> 如果一个方法或字段不能在对象的结构定义中被找到，则会继续在内嵌结构中查询，并直接在内嵌对象上进行方法调用。

Interfaces in Go are **abstract representations of behavior** - sets of method signatures.

> 接口是行为的抽象表示，包含一组方法声明。

Interfaces are allowed to embed other interfaces.

> 支持内嵌接口。
>
> ```go
> type Writer interface {
> 	Write()
> }
>
> type Closer interface {
> 	Writer
> 	Close()
> }
>
> type DiskCloser struct{}
>
> func (d *DiskCloser) Write() {}
>
> func (d *DiskCloser) Close() {}
> ```

Every objects implements the empty interface `interface{}`.

> 所有对象都隐式实现了空接口 `interface {}`。

The dispatch of method calls with embedding is different from that of class-based inheritance.

> 内嵌类型方法调用的分配与类继承的方法调用分配不同。
>
> ```go
> package main
>
> import "fmt"
>
> type T struct{}
>
> func (t T) Foo() {
> 	fmt.Println("T")
> }
>
> func (t T) Baz() {
> 	t.Foo()
> }
>
> type S struct {
> 	T
> }
>
> func (s S) Foo() {
> 	fmt.Println("S")
> }
>
> func main() {
> 	s := new(S)
> 	s.Baz() // T
> }
> ```

Wildcards cannot be used for importing packages, every package has to be listed individually.

> 不能通过通配符一次引入多个包，每一个包的引入需要独立列出。

An object of type `T` can also be created with the built-in function `new()`. `new(T)` allocates zeroed storage for a new object of type `T` and returns a pointer to it.

> 可以使用内置函数 `new()` 创建 `T` 类型的对象。

Note that, the first return parameter could be `nil` even though the second is `true`, since values stored in a map can be `nil`.

> 因为 `map` 的值可以是 `nil`，所以即使第二个返回值是` true`，第一个返回值也可能是 `nil`。
>
> ```go
> func main() {
> 	m := map[string]interface{}{
> 		"id1": "123",
> 		"id2": nil,
> 	}
>
> 	id1, ok := m["id1"]
> 	fmt.Println(id1, ok) // 123 true
>
> 	id2, ok := m["id2"]
> 	fmt.Println(id2, ok) // <nil> true
> }
> ```

### Design Patterns

Patterns are a description of reoccurring design problems and general solutions for the problem.

> 设计模式描述的是常见的设计问题并给出通用的解决方案。

A pattern consists usually of a name; a concise, abstract description of the problem; an example; detail of implementation; advantages and disadvantages; and known uses.

> 一个设计模式通常包含：
>
> * 名称
> * 简洁且抽象的问题描述
> * 具体的实现
> * 优缺点
> * 已知的用途