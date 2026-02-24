---
title: "Go flag 自定义选项类型"
date: 2023-03-09T20:38:36+08:00
draft: false
reward: false
categories:
  - 后端技术
  - Go
tags:
 - go
 - flags
---

`flag` 包定义了一系列函数，可用于定义命令行参数，支持的参数类型如下：

* string：flag.StringVar 函数
* bool：flag.BoolVar 函数
* time.Duration: flag.DurationVar 函数
* int: flag.IntVar 函数
* uint: flag.UintVar 函数
* float64: flag.Float64Var 函数
* int64: flag.Int64Var 函数
* uint64: flag.Uint64Var 函数

<!--more-->

各函数声明见下图：

![](define_flag_functions.png)

直接调用上述函数的话，会给默认的 `CommandLine` 定义一个命令行选项，内部则是调用了 `FlagSet.Var` 函数，如下：

```go
// StringVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a string variable in which to store the value of the flag.
func StringVar(p *string, name string, value string, usage string) {
	CommandLine.Var(newStringValue(value, p), name, usage)
}

// Var defines a flag with the specified name and usage string. The type and
// value of the flag are represented by the first argument, of type Value, which
// typically holds a user-defined implementation of Value. For instance, the
// caller could create a flag that turns a comma-separated string into a slice
// of strings by giving the slice the methods of Value; in particular, Set would
// decompose the comma-separated string into the slice.
func (f *FlagSet) Var(value Value, name string, usage string) {
  	// ...
}
```

`FlagSet.Var` 函数的第一个参数类型为 `Value`，所以自定义选项类型需要从 `Value` 接口入手。

## 代码

`Value` 接口的声明如下（★ ★ ★ ★ ★）：

```go
type Value interface {
	String() string
	Set(string) error
}
```

自定义的选项类型要实现 `Value` 接口，完整代码如下：

```go
// user_defined_flag_type.go
package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

// Gender 性别
type Gender int

const (
	// genderMale 男性
	genderMale Gender = iota
	// genderFemale 女性
	genderFemale
)

var (
	// ErrUserParse 无法解析用户信息
	ErrUserParse = errors.New("parse user error")
)

// User 自定义选项结构
// 支持以命令行选项的方式实例化，使用 : 作为分隔，选项值形如：`xiaoming:18:0`
type User struct {
	Name   string
	age    int
	gender Gender
}

func (u *User) String() string {
	return "user"
}

// Set 解析传递来的字符串，实例化 User 结构
func (u *User) Set(s string) error {
	tokens := strings.Split(s, ":")
	if len(tokens) != 3 {
		return ErrUserParse
	}
	u.Name = tokens[0]

	v, err := strconv.ParseInt(tokens[1], 0, strconv.IntSize)
	if err != nil {
		return err
	}
	u.age = int(v)

	v, err = strconv.ParseInt(tokens[2], 0, strconv.IntSize)
	if err != nil {
		return err
	}
	gender := Gender(v)
	if !(gender == genderMale || gender == genderFemale) {
		return ErrUserParse
	}
	u.gender = gender

	return nil
}

func main() {
	var user = User{}
	flag.Var(&user, "user", "user information")
	flag.Parse()

	fmt.Printf("name: %s\n", user.Name)
	fmt.Printf("age: %d\n", user.age)
	fmt.Printf("gender: %d\n", user.gender)
}
```

## 运行

**正常解析**

```bash
$ go run user_defined_flag_type.go --user=xiaoming:18:0
name: xiaoming
age: 18
gender: 0
```

**user 解析失败：信息不足**

```bash
$ go run user_defined_flag_type.go --user=xiaoming:18
invalid value "xiaoming:18" for flag -user: parse user error
```

**user 解析失败：性别不符**

```bash
$ go run user_defined_flag_type.go --user=xiaoming:18:2
invalid value "xiaoming:18:2" for flag -user: parse user error
```

## 小结

如果要实现自定义的选项类型，则该选项需要实现 `flag.Value` 接口，关键在于 `Set(s string) error` 如何实现字符串的解析并将信息转换成类型实例信息。[代码](https://github.com/a2htray/code-notebook/blob/main/BlogCode/Go/user_defined_flag_type/user_defined_flag_type.go)。