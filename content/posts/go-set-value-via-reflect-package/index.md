---
title: "使用 reflect 包设置值"
date: 2022-04-21T19:59:01+08:00
draft: true
comment: false
reward: false
categories:
 - Go
tags:
 - reflect
 - 反射
images:
 - images/golang-reflect.png
---

首先贴上 Go 开发团队对 reflect 包的描述：

> Package reflect implements run-time reflection, allowing a program to manipulate objects with arbitrary types. The typical use is to take a value with static type interface{} and extract its dynamic type information by calling TypeOf, which returns a Type.
>
> A call to ValueOf returns a Value representing the run-time data. Zero takes a Type and returns a Value representing a zero value for that type.

从描述中，我们得到以下几点：

1. reflect 包实现了运行时的反射机制，允许程序操作任意类型的对象；
2. TypeOf 可以得到一个 interface{} 的具体类型，ValueOf 可以得到一个 interface{} 的具体值；

<!--more-->

## 重要的类型

reflect 包中定义了几种重要的、常用的类型，分别为：

1. Kind

### Kind

Kind 类型用于修饰类型 `Type`，用于表示 `Type` 的种类，底层是用一个**无符号整数**表示：

```go
type Kind uint;
```

其中修饰了的基础数据类型包括`Bool`、`Int`、`Int32`、`Float32`、`Float64` 等，引用数据类型包括`Slice`、`Map`、`Chan`、`Interface`、`Func` 等。

```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	var v1 int = 1
	var v2 float32 = 1.0
	v3 := struct {
	}{}
	v4 := func() {}
	PrintKind(v1)
	PrintKind(v2)
	PrintKind(v3)
	PrintKind(v4)
}

func PrintKind(v interface{}) {
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Int:
		fmt.Println("int")
	case reflect.Float32:
		fmt.Println("float32")
	case reflect.Struct:
		fmt.Println("struct")
	case reflect.Func:
		fmt.Println("function")
    // 其它的情况 ...
	}
}
```

```bash
int
float32
struct
function
```

如果只是单纯地打印类型 `Type` 的种类，直接调用 `.String` 方法即可。

```go
func PrintKind(v interface{}) {
	t := reflect.TypeOf(v)
	fmt.Println(t.String())
}
```

```bash
int
float32
struct {}
func()
```

### Value

Value 类型表示一个接口具体的值，类型的定义如下：

```go
type Value struct {
	// 值的真实类型
	typ *rtype
	// 指向值的指针
	ptr unsafe.Pointer
	// 值的元数据信息
	// flag 可以是多种情况的组合
	flag
}
```

在知道值类型的情况下，可以调用相应的函数获取其本身的值，即值的类型是 `int`，经过 `ValueOf` 返回 `Value` 类型的值就可以调用 `Int` 方法得到其值。但是，如果调用了不符合本身类型的方法会报错，这一类方法的开头都对值的类型做了类型判断。

```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	var v1 int = 1
	fmt.Println(reflect.ValueOf(v1).Int())
	var v2 float32 = 1.0
	// fmt.Println(reflect.ValueOf(v2).Int()) 会报错
}
```

```go
func (v Value) Int() int64 {
	k := v.kind()
	p := v.ptr
	switch k {
	case Int:
		return int64(*(*int)(p))
	case Int8:
		return int64(*(*int8)(p))
	case Int16:
		return int64(*(*int16)(p))
	case Int32:
		return int64(*(*int32)(p))
	case Int64:
		return *(*int64)(p)
	}
	panic(&ValueError{"reflect.Value.Int", v.kind()})
}

func (f flag) mustBe(expected Kind) {
	// TODO(mvdan): use f.kind() again once mid-stack inlining gets better
	if Kind(f&flagKindMask) != expected {
		panic(&ValueError{methodName(), f.kind()})
	}
}
```

### SliceHeader

SliceHeader 是 slice 的底层实现，其结构如下：

```go
type SliceHeader struct {
    // 指向底层数据的指针
	Data uintptr
    // slice 的长度
	Len  int
    // slice 的容量
	Cap  int
}
```





## 参考

1. [Go 标准库 reflect](https://pkg.go.dev/reflect)