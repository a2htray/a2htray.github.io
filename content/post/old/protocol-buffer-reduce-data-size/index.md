---
title: "Protocol Buffer 减少传输数据的大小"
date: 2022-04-09T15:01:17+08:00
draft: false
reward: false
categories:
 - 后端技术
 - Go
tags:
 - Go
 - protobuf
---

Protocol Buffer 的介绍与语法已在文章[《Protocol Buffer 语法》](/posts/protocol-buffer-syntax/)给出，本文则演示了 Protocol Buffer 如何减少了传输数据的大小。

<!--more-->

## 使用 protoc 命令生成 pb 代码文件

首先新建 `proto/user.proto` 文件来定义数据结构，其内容如下：

```protobuf
syntax = "proto3";

option go_package = "/model";

message User {
    string name = 1;
    enum Gender {
        MALE = 0;
        FEMALE = 1;
    };
    Gender gender = 2;
}
```

然后，执行如下命令生成源代码文件：

```bash
protoc --go_out=. proto\user.proto
```

## 主程序 main

主程序的作用是比较不同 `User` 结构序列化后的字节数据的大小。首先，新建 `main.go` 文件并定义一个 `User` 结构，如下：

```go
type User struct {
	Name   string
	Gender int32
}
```

 主程序的逻辑如下：

```go
func main() {
    user := &User{
		Name:   "a2htray",
		Gender: 1,
	}
	data, err = json.Marshal(user)
	fmt.Println(data, err)
}
```

执行输出如下：

```bash
[123 34 78 97 109 101 34 58 34 97 50 104 116 114 97 121 34 44 34 71 101 110 100 101 114 34 58 49 125] 29 <nil>
```

从输出可知，自定义的 `User` 的值序列化后的字节长度为 29。

接着，使用 Protocol Buffer 生成的 `User` 结构并使用 `proto.Marshal` 方法对值进行序列化，代码如下：

```go
func main() {
    userPB := &model.User{
		Name:   "a2htray",
		Gender: 1,
	}
	data, err = proto.Marshal(userPB)
	fmt.Println(data, len(data), err)
}
```

执行输出如下：

```bash
[10 7 97 50 104 116 114 97 121 16 1] 11 <nil>
```

从输出可知，Protocol Buffer 生成的 `User` 类型的值序列化后的字节长度为 11。

综上，分别使用 JSON 和 Protocol Buffer 序列化相同的数据信息，使用 Protocol Buffer 得到的字节长度要更小，更有得于在网络中的传输。

## 完整代码

演示完成后，当前项目的目录结构如下：

```bash
model
  - user.pb.go # 通过 protoc 命令生成
proto
  - user.proto # 定义数据结构
main.go
```

`main.go` 的完整内容如下：

```go
package main

import (
	"encoding/json"
	"fmt"

	"./model"
	"google.golang.org/protobuf/proto"
)

type User struct {
	Name   string
	Gender int32
}

func main() {
	user := &User{
		Name:   "a2htray",
		Gender: 1,
	}
	data, err := json.Marshal(user)
	fmt.Println(data, len(data), err)

	userPB := &model.User{
		Name:   "a2htray",
		Gender: 1,
	}
	data, err = proto.Marshal(userPB)
	fmt.Println(data, len(data), err)
}
```

## 总结

1. Protocol Buffer 序列化的数据量更小；