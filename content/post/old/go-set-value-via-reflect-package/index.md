---
title: "Go 的反射包 reflect"
date: 2022-04-21T19:59:01+08:00
draft: false
reward: false
categories:
  - 后端技术
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

## 重要类型

reflect 包中定义了几种重要的、常用的类型，分别为：

1. Kind；
2. Value；
3. SliceHeader；
4. StringHeader；
5. Method；
6. StructField；

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

SliceHeader 是 slice 运行时的底层实现，其结构如下：

```go
type SliceHeader struct {
    // 指向底层数据的指针
    // 无符号的整数表示内存中的地址
	Data uintptr
    // slice 的长度
	Len  int
    // slice 的容量
	Cap  int
}
```

```go
func main() {
	data := []int{1, 2, 3, 4, 5}
	// unsafe.Pointer 表示任意类型的指针
	dataPointer := unsafe.Pointer(&data)
	sh1 := (*reflect.SliceHeader)(dataPointer)
	fmt.Printf(`SliceHeader {
    Data: %d
    Len: %d
    Cap: %d
}
`, sh1.Data, sh1.Len, sh1.Cap)
	for i := 0; i < sh1.Len; i++ {
		// slice 中各元素的地址
		addr := uint(sh1.Data) + uint(unsafe.Sizeof(1)) * uint(i)
		// slice 中各元素的值
		value := *(*int)(unsafe.Pointer(uintptr(addr)))
		fmt.Printf("%dth element addr: 0x%x, value: %d\n", i, addr, value)
	}
}
```

```bash
SliceHeader {
    Data: 824634670792
    Len: 5
    Cap: 5
}
0th element addr: 0xc0000e7ec8, value: 1
1th element addr: 0xc0000e7ed0, value: 2
2th element addr: 0xc0000e7ed8, value: 3
3th element addr: 0xc0000e7ee0, value: 4
4th element addr: 0xc0000e7ee8, value: 5
```

上述代码构造了一个 SliceHeader 类型的变量 `sh1`，并逐一打印出各元素的地址和值：

1. `&data` 得到了 slice 的指针，使用 ` unsafe.Pointer(&data)` 强制转换变成 `*ArbitraryType`；
2. 使用 `(*reflect.SliceHeader)(dataPointer)` 得到了一个指向 SliceHeader 的指针；
3. `fmt.Printf` 打印变量中 3 个字段的信息；
4. for 循环：先计算每一个元素的地址，再通过地址得到地址上的值，最后打印输出；
5. 地址和值都需要进行一系列转换；
   1. 如 uintptr 不支持算法运算符，所以通过转换成 uint 类型进行计算；
   2. 值的获取则是先将 uint 转换成 uintptr，再通过 `unsafe.Pointer` 转换成可表示任意类型的指针，接着转成 `*int` 类型的指针，最后使用 `*` 取指针的值；

### StringHeader

StringHeader 是 string 运行时的底层实现，其结构如下：

```go
type StringHeader struct {
    // 指向底层数据的指针
    // 无符号的整数表示内存中的地址
	Data uintptr
    // 字符串的长度
	Len  int
}
```

```go
func main() {
	data := []string{"1", "22", "333", "4444", "55555"}
	sh1 := (*reflect.StringHeader)(unsafe.Pointer(&data))
	fmt.Printf("string slice length: %d\n", sh1.Len)

	for i := 0; i < sh1.Len; i++ {
		// 计算第 i 个元素的地址
		addr := uint(sh1.Data) + uint(i) * uint(unsafe.Sizeof(data[i]))
		// 转换成任意类型的指针
		arbitraryPointer := unsafe.Pointer(uintptr(addr))
		// 转换成字符串指针
		stringPointer := (*string)(arbitraryPointer)
		value := *stringPointer
		fmt.Printf("%dth value: %s\n", i, value)
	}
}
```

```bash
string slice length: 5
0th value: 1
1th value: 22
2th value: 333
3th value: 4444
4th value: 55555
```

StringHeader 的例子与 SliceHeader 类似，区别在于在指针转换时 StringHeader 使用了 `*string`。

### Method

Method 表示一个方法，可以使用 reflect 包动态的调用方法，传入的参数是一个 `reflect.Value` 的 slice，其返回值也是使用一个 `reflect.Value` 的 slice 表示。

```go
type anyFunc func(int) int

type A struct {}

func (a A) AnyFunc() string { return "A" }

func main() {
	var oneFunc anyFunc = func(i int) int {
		return i + 1
	}

	// oneFunc 是一个函数指针
	fmt.Printf("1. pointer %p\n", oneFunc)
	// 函数类型的字符串表示
	fmt.Printf("2. %s\n", reflect.TypeOf(oneFunc).String())
	// 函数指针
	fmt.Printf("3. %s\n",  (reflect.ValueOf(&oneFunc)).Kind().String())

	fmt.Println("== call method by reflect package")
	i := 1000
	// reflect.Value 得到函数值
	// .Call 调用函数
	// []reflect.Value{reflect.ValueOf(i)} 是传入函数的参数
	result := reflect.ValueOf(oneFunc).Call([]reflect.Value{reflect.ValueOf(i)})
	fmt.Println(result[0])

	// 类型中的方法
	a := A{}
	method := reflect.TypeOf(a).Method(0)
	fmt.Printf(`== basic information of A.AnyFunc
Type: %s
Name: %s
PkgPath： %s
Index: %d
`, method.Type, method.Name, method.PkgPath, method.Index)

	// 调用类型中的方法
	fmt.Println("== call a.AnyFunc")
	fmt.Println(reflect.ValueOf(a).Method(0).Call([]reflect.Value{}))
}
```

### StructField

StructField 表示一个结构的字段，StructField 还会保存字段的 Tag 信息（StructTag）。StructTag 的底层实现是 string：

```go
type StructTag string
```

```go

type A struct {
	StringField string `json:"jsonStringField"`
	IntField    int    `json:"jsonIntField"`
}

func main() {
	a := A{
		StringField: "simple string",
		IntField:    100,
	}
	fmt.Printf("number of struct fields: %d\n", reflect.TypeOf(a).NumField())

	fmt.Println("== get struct field value")
	numField := reflect.TypeOf(a).NumField()
	for i := 0; i < numField; i++ {
		fmt.Printf("%s: %v\n", reflect.TypeOf(a).Field(i).Name, reflect.ValueOf(a).Field(i))
	}

	fmt.Println("== get struct field tag")
	for i := 0; i < numField; i++ {
		tag := reflect.TypeOf(a).Field(i).Tag
		fmt.Printf("%s's tag string: `%s`\n", reflect.TypeOf(a).Field(i).Name, tag)
		fmt.Printf("  json: %v\n", tag.Get("json"))
	}
}
```

## 总结

1. reflect.ValueOf(值) 和 reflect.ValueOf(指针).Elem() 等价；
2. unsafe.Pointer 提供了任意类型指针的访问；

## 参考

1. [Go 标准库 reflect](https://pkg.go.dev/reflect)
2. [stackoverflow 动态调用对象的方法](https://stackoverflow.com/questions/8103617/call-a-struct-and-its-method-by-name-in-go)

