---
title: "Go 内置的 RPC 包"
date: 2022-04-10T16:37:30+08:00
draft: false
reward: false
categories:
  - 后端技术
  - Go
tags:
 - RPC
 - jsonrpc
images:
 - "images/operating-system-remote-procedure-call-1.png"
---

<!--Package rpc provides access to the exported methods of an object across a network or other I/O connection. A server registers an object, making it visible as a service with the name of the type of the object. After registration, exported methods of the object will be accessible remotely. A server may register multiple objects (services) of different types but it is an error to register multiple objects of the same type.-->

在网络或 I/O 连接中，可以使用 `net/rpc` 包实现对一个对象的导出方法的调用，即远程过程调用（Remote Procedure Call，RPC）。通过向 RPC 服务注册一个对象，使其可被远程调用，进而实现一些复杂的业务逻辑。

## 项目结构

示例项目的结构如下：

```bash
client
  - client.go
  - json_client.go
models
  - greeting.go
server
  - json_server.go
  - server.go
```

## 注册服务

一个可被远程调用的方法须满足以下条件：

1. 方法所属结构是公开的；
2. 方法是分开的；
3. 方法的参数类型是分开的；
4. 方法带两个参数，第 2 个参数为指针；
5. 方法返回值为 error 类型；

如下，在 `models/greeting.go` 中定义了一个服务：

```go
type GreetingArg struct {
	Name string
}

type GreetingReply struct {
	Message string
}

type Greeting struct {}

// SayHello 方法满足上述条件
func (Greeting) SayHello(arg GreetingArg, reply *GreetingReply) error {
	reply.Message = "hello, " + arg.Name
	return nil
}
```

现在，在 `server/server.go` 中编写服务器端代码：

```go
package main

import (
	"gorpc/models"
	"log"
	"net"
	"net/rpc"
)

func main() {
	server := rpc.NewServer()
	if err := server.Register(&models.Greeting{}); err != nil {
		log.Fatalln(err)
	}

	listener, err := net.Listen("tcp", ":2022")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	server.Accept(listener)
}
```

服务器端注册了 `Greeting` 服务并监听了 2022 端口，等待客户端连接。在客户端 `client/client.go` 的代码如下：

```go
package main

import (
	"fmt"
	"gorpc/models"
	"log"
	"net"
	"net/rpc"
)

func main() {
	conn, err := net.Dial("tcp", ":2022")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := rpc.NewClient(conn)
	greetingArg := models.GreetingArg{Name: "a2htray"}
	greetingReply := models.GreetingReply{}

	if err = client.Call("Greeting.SayHello", greetingArg, &greetingReply); err != nil {
		log.Fatalln(err)
	}
	fmt.Println(greetingReply.Message)
}
```

上述代码完成了以下几件事：

1. 使用 `net.Dial` 连接 2022 端口；
2. 在 TCP 连接之上，使用 `rpc.NewClient` 创建一个 RPC 客户端；
3. 使用 `client.Call` 远程调用 `Greeting` 的 `SayHello` 方法；
4. 返回的值体现在 `greetingReply` 变量中；

### jsonrpc

`net/rpc` 的传输数据使用 `encoding/gob` 进行编码解码，并且不支持跨语言调用，即只能使用 Go 编写的程序进行调用。`encoding/gob` 编码解码在源码中有给出：

```go
// rpc/server.go
func (server *Server) ServeConn(conn io.ReadWriteCloser) {
	buf := bufio.NewWriter(conn)
	srv := &gobServerCodec{
		rwc:    conn,
		dec:    gob.NewDecoder(conn),
		enc:    gob.NewEncoder(buf),
		encBuf: buf,
	}
	server.ServeCodec(srv)
}
```

除了 `net/rpc`，还可以使用 `net/rpc/jsonrpc` 实现 RPC 功能，该方式支持跨语言调用。新建 `server/json_server.go`，代码如下：

```go
package main

import (
	"gorpc/models"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {
	err := rpc.Register(&models.Greeting{})
	if err != nil {
		log.Fatalln(err)
	}
	listener, err := net.Listen("tcp", ":2023")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		go jsonrpc.ServeConn(conn)
	}
}
```

上述代码完成了以下几件事：

1. 在 RPC 服务上注册了 `Greeting`；
2. 监听了 2023 端口，使用 for 循环接受客户端连续；
3. 对每一个连接使用协程进行处理；

新建 `client/json_client.go`，代码如下：

```go
package main

import (
	"fmt"
	"gorpc/models"
	"log"
	"net/rpc/jsonrpc"
)

func main() {
	client, err := jsonrpc.Dial("tcp", ":2023")
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
	
	greetingArg := models.GreetingArg{Name: "a2htray"}
	greetingReply := models.GreetingReply{}

	if err = client.Call("Greeting.SayHello", greetingArg, &greetingReply); err != nil {
		log.Fatalln(err)
	}
	fmt.Println(greetingReply.Message)
}
```

上述代码完成了以下几件事：

1. 使用 `jsonrpc.Dial` 连接到端口 2023；
2. 使用 `client.Call` 调用了 `Greeting.SayHello` 方法；
3. 打印输出返回信息；

## rpc 与 jsonrpc 的区别

Go 内置的 rpc 与 jsonrpc 的区别在于：

1. rpc 使用 `gob` 编码解码，jsonrpc 使用 `json` 编码解码；
2. rpc 不支持跨语言调用，jsonrpc 支持跨语言调用；
3. jsonrpc 在构建在 rpc 之上使用不同数据交换格式的 RPC 服务；

## 参考

1. [golang下的rpc框架jsonrpc理解和使用示例](https://www.pudn.com/news/624c4e61fc37f87c248d01b8.html)
2. [golang实现RPC的几种方式](https://cxymm.net/article/kenkao/99713753)

