---
title: "Go 泛型函数 - 解构类型参数"
date: 2024-01-07T14:41:11+08:00
draft: false
reward: false
categories:
 - 后端技术
 - Go
tags:
 - 类型推断
 - 泛型函数
 - 翻译
---

原文地址：[https://go.dev/blog/deconstructing-type-parameters](https://go.dev/blog/deconstructing-type-parameters)

**译者评论**

本文通过 slices.Clone 泛型函数介绍了 Go 是如何使用<font color="#0d6efd" style="font-weight: bold;">类型推断</font>完成参数类型的解构。简单来说，如果第一个类型参数是一个复合类型，则可以通过第二、第三或更多的类型参数约束复杂类型中的类型参数，而类型推断则可以通过第一个参数推断出后续类型参数的实际类型。另外本文还说明为消除歧义而引入 <font color="#0d6efd" style="font-weight: bold;">~ 符号</font>，即用于指定类型的底层类型。

<!--more-->

## slices 包函数签名

slices.Clone 函数非常简单，该函数可以克隆任意类型元素的 slice，函数签名如下：

```go
func Clone[S ~[]E, E any](s S) S {
    return append(s[:0:0], s...)
}
```

slices.Clone 函数可以正常运行，因为<font color="#0d6efd" style="font-weight: bold;">向一个空容量的 slice 添加元素会在底层分配一个新的数组</font>。我们看到，函数体内容要比函数签要短，为什么会这样？在这篇文章中，我们会解释如此设计函数签名的原因。

## 简单 Clone 实现

首先，我们在 slices 包之外写一个简单的泛型 Clone1 函数。我们希望函数接收一个任意元素类型的 slice，并返回一个 slice。

```go
func Clone1[E any](s []E) []E {
    // body omitted
}
```

Clone1 泛型函数只有一个类型参数 E，且函数参数是一个类型为 E 的 slice s，并返回一个相同类型的 slice。对于了解 Go 泛型的开发者来说，函数签名直接了当，很好理解。

然而，这里隐含着一个问题。尽管 slice 命名类型在 Go 中不常见，但开发者偶尔也会使用到。

```go
// MySlice 是一个包含特殊方法的 string slice
type MySlice []string

// String 返回 slice 拼接的结果
func (s MySlice) String() string {
    return strings.Join(s, "+")
}
```

现在我们想要实现一个 MySlice 拷贝再打印出里面的内容，只不过其中字符串是经过排序的。

```go
func PrintSorted(ms MySlice) string {
    c := Clone1(ms)
    slices.Sort(c)
    return c.String() // FAILS TO COMPILE
}
```

但不幸的是，上述代码不能正常工作，编译器会报一个 error：

```bash
c.String undefined (type []string has no field or method String)
```

如果用传参类型代替类型参数的类型，手动实例化 Clone1 函数（如 InstantiatedClone1），就会看到这个问题。

```go
func InstantiatedClone1(s []string) []string
```

Go 的赋值规则允许我们为类型为 []string 的参数传入一个类型为 MySlice 的值，<font color="#0d6efd" style="font-weight: bold;">但是 Clone1 函数返回值的类型为 []string，并不是 MySlice</font>。由于 []string 类型没有 String 方法，所以编译器报了这个错。

## 更灵活的 Clone

为了解决这个问题，我们不得不再写一个 Clone2 函数，该函数返回值的类型与其参数类型一致。如果这么做了，然后向 Clone2 函数传入类型为 MySlice 的值，函数也会返回一个类型为 MySlice 的结果。

新函数签名就应该像下面一样。

```go
func Clone2[S ?](s S) S // INVALID
```

Clone2 函数返回值的类型与其参数类型一致。

这里我们将约束定义为 ?，但这只是一个占位符。在 Clone1 函数中，我们使用了一个 any 类型参数，但在 Clone2 函数中，并不适合用 any 类型参数，因为我们需要一个 slice 类型。

因为我们知道需要一个 slice，S 的约束就必须是一个 slice。我们先不管 slice 的元素类型，先姑且像 Clone1 函数中一样，声明一个类型参数 E。

```go
func Clone3[S []E](s S) S // INVALID
```

上面的代码依然是无效的，因为我们并没有定义类型参数 E。类型参数 E 可以是任意类型，这也就意味着 E 可以是其本身，因为可以是任意类型，所以 E 的约束应该 any。

```go
func Clone4[S []E, E any](s S) S
```

我们已经慢慢接近正确答案了，至少现在是可以编译通过的。如果我们按上述函数声明运行，依然会报一个 error。

```go
MySlice does not satisfy []string (possibly missing ~ for []string in []string)
```

编译器告诉我们，我们不能用 MySlice 作为 S 的类型参数，因为 MySlice 并不满足约束 []E。<font color="#0d6efd" style="font-weight: bold;">这是因为 []E 只允许 slice 字面量，不支持 MySlice</font>。

## 底层类型约束

因为上面的错误提示，我们应该加上 ~。

```go
func Clone5[S ~[]E, E any](s S) S
```

再次声明，约束 [S []E, E any] 表示<font color="#0d6efd" style="font-weight: bold;">类型参数 S 可以是一个 slice 类型，但不能是一个 slice 的命名类型</font>。将其改成 [S ~[]E, E any]，<font color="#0d6efd" style="font-weight: bold;">则表示 S 可以是任意底层类型是 slice 的任意类型</font>。

因为 MySlice 的底层类型是一个 string slice，所以我们可以向 Clone5 函数传递类型为 MySlice 的参数。最终，Clone5 函数的签名已经和 slices.Clone 函数的签名一样了。

在继续之前，我们讨论下为什么 Go 语法需要 ~。就好像我们只需要总是允许传递 MySlice 类型就好了，为什么不使用其默认类型。或者，我们只需要实现精确匹配就好，即用 []E 只匹配 slice 类型字面量。

为解释这个原因，我们首先<font color="#0d6efd" style="font-weight: bold;">观察到使用 [T ~MySlice] 没有多大意义，因为 MySlice 并不是其它类型的底层类型</font>。比如，我们定义一个类型 MySlice2 MySlice

```go
type MySlice2 MySlice
```

但 MySlice2 的底层类型依然是 []string。所以 [T ~MySlice] 的作用与 [T MySlice] 一样，仅仅只约束到了 MySlice 类型的传参。为了消除困惑，所以编译器禁用了这一类使用，不然会报以下的错误

```bash
invalid use of ~ (underlying type of MySlice is []string)
```

如果 Go 语法中不使用 ~，那么 <font color="#0d6efd" style="font-weight: bold;">[S []E] 将会精确匹配到任意以 []E 作为底层类型的类型，这样我们就不得不定义 [S MySlice] 作为约束</font>。

Go 语法禁止 [S MySlice]，或者说 [S MySlice] 只能匹配到 MySlice，但是对语言预定义的类型会造成困惑。作为预定义类型的 int，其底层类型依然是 int。我们希望 Go 语言能够能开发者提供精确匹配和定义约束底层类型为 int 的方式，如在程序中使用 [T ~int]。如果我们不使用 ~，<font color="#0d6efd" style="font-weight: bold;">[T int] 不能很好表明要使用底层类型为 int 语义。如果这么做了，那么 [T MySlice] 和 [T int] 的约束行为就会有歧义</font>。

我们可能会认为 [S MySlice] 匹配任意底层类型为 MySlice 的底层类型的类型，但这样会很困惑。

所以我们觉得使用 ~ 表明其底层类型会更好一些。

## 类型推断

我们已经解释了为什么 slices.Clone 的签名，现在再看看如果在实际中使用 slice.Clone。slices.Clone 的签名如下：

```go
func Clone[S ~[]E, E any](s S) S
```

在调用 slices.Clone 时，需要将一个 slice 传递给参数 s。在类型推断时，<font color="#0d6efd" style="font-weight: bold;">编译器在 slice 传递给 Clone 之前就会推断出类型参数 S 的实际参数类型</font>。类型推断足够强大，即可以确定类型参数 E 是传递给 S 中的元素类型。

所以我们可以像下面这么写：

```go
c := Clone(ms)
```

而不需要像下面这么写：

```go
c := Clone[MySlice, string](ms)
```

如果我们要引用 Clone，而不是直接调用，因为编译器没有任何信息作为推断，所以我们需要指定 S 的类型参数。幸运的是在这个示例中，类型抢断可以从 S 的参数中推断出 E 的类型参数，所以我们不需要依次指定类型参数。

这就是为什么我们可以像下面这么写：

```go
myClone := Clone[MySlice]
```

而不需要像下面这么写：

```go
myClone := Clone[MySlice, string]
```

## 解构类型参数

我们使用到的技术，即定义一个使用类型参数 E 的类型参数 S，是<font color="#0d6efd" style="font-weight: bold;">一种在泛型函数签名中解构类型的方式</font>。通过解构类型，我们可以命名、约束类型的各个方面。

比如，maps.Clone 的签名如下：

```go
func Clone[M ~map[K]V, K comparable, V any](m M) M
```

和 slices.Clone 一样，我们使用了类型参数 M 来约束参数 m，然后定义类型参数 K 和 V 用于解构类型。

在 maps.Clone 中，我们约束 K 必须是可比较型的，这与 map 的 key 的约束一致。也正因为这一特性，我们可以在开发过程中实现对复合类型的解构。

```go
func WithStrings[S ~[]E, E interface { String() string }](s S) (S, []string)
```

上述示例中，我们要求 WithStrings 的参数类型必须是一个元素类型为带 String 方法的 slice。

因此，我们可以 Go 语言中在复合类型中使用类型推断来推断出其实际类型。
