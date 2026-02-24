---
title: "gRPC 的基本概述"
date: 2022-04-08T15:39:54+08:00
draft: true
reward: false
categories:
  - 后端技术
  - Go
tags:
 - gRPC
 - RPC
 - protocol buffer
---

[gRPC](https://grpc.io/) 是开源的高性能的远程过程调用（Remote Procedure Call，RPC）框架，其特点包括：

* 服务定义简单易用；
* 启动快，高性能；
* 双向的数据流传输；
* 插件支持；
* 跨平台；
* 跨语言，支持众多主流编程语言，包括 C#、C++、Dart、Go 等；

<!--more-->

<div style="display: flex;">
    <img src="gRPC-request-response.svg" style="margin: 0 auto;" />
</div>

<p style="text-align: center;"><i>图 1 gRPC 通信过程</i></p>

如图 1，服务器端部署了C++版本 gRPC 服务，两个客户端分别为 Ruby 和 Android-Java 实现，不同编程语言的实现可通过 gRPC 发生远程调用。

## 相关概念

### 服务定义

在 gRPC 框架中，服务器端定义了一个服务，该服务指定了一个可接收参数且可返回的函数，客户端调用函数实现远程调用。默认情况下，gRPC 使用 [protocol buffers](https://developers.google.com/protocol-buffers) 作为接口定义语言（Interface Definition Language，IDL），用于约定服务的接口以及传值的结构。gRPC 定义了 4 种服务方法的类型：

1. `rpc SayHello(HelloRequest) returns (HelloResponse)`：客户端发送请求，服务器端返回响应；
2. `rpc LotsOfReplies(HelloRequest) returns (stream HelloResponse)`：客户端发送请求，服务器返回流式数据，客户端流式读取，直至读写所有响应；
3. `rpc LotsOfGreetings(stream HelloRequest) returns (HelloResponse)`：客户端流式请求，服务端返回响应；
4. `rpc BidiHello(stream HelloRequest) returns (stream HelloResponse)`：客户端与服务器端均以流式数据进行传输；

> gRPC 提供了流式传输和读取数据的功能

### .proto 文件

.proto 文件定义了 RPC 服务器可对外提供的功能。下面分别就上述的 4 种服务方法类型编写如下文件：

```protobuf
syntax = "proto3";


```

### 接口定义

gRPC 可以解析 .proto 文件来生成服务器端和客户端的代码，



## 与 JSON 的区别



