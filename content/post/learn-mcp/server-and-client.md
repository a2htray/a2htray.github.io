+++
date = '2026-02-05T16:06:49+08:00'
draft = false
title = 'MCP：服务器和客户端'
categories = ['人工智能', '大语言模型']
tags = ['MCP', '摘录']
toc = true
+++

本文继续做官方文档的学习和摘录，在服务器和客户端这一块，理解有效，有些翻译甚至有些勉强。

学习来源：

1. [https://modelcontextprotocol.io/docs/learn/server-concepts](https://modelcontextprotocol.io/docs/learn/server-concepts)
2. [https://modelcontextprotocol.io/docs/learn/client-concepts](https://modelcontextprotocol.io/docs/learn/client-concepts)

## 服务端

> MCP servers are programs that expose specific capabilities to AI applications through standardized protocol interfaces.

重点：<span style="color: blue; font-weight: bold;">对外提供特殊能力</span>。

服务端关键的特性：

1. 工具（Tools）：LLM 可调用的功能（<span style="color: blue; font-weight: bold;">由模型控制</span>）
2. 资源（Resources）：只读访问的数据源（<span style="color: blue; font-weight: bold;">由应用控制</span>）
3. 提示词（Prompts）：预设的指令模板（<span style="color: blue; font-weight: bold;">由用户控制</span>）

### 工具

Tools enable AI models to perform actions. Each tool defines a specific operation with typed inputs and outputs. The model requests tool execution based on context.

<span style="color: blue; font-weight: bold;">有定义输入和输出的功能</span>，可以认为是工具。MCP 约定的与工具相关的操作有：

| 方法         | 目的      | 返回        |
|------------|---------|-----------|
| tools/list | 发现可用的工具 | 模式固定的工具列表 |
| tools/call | 执行特定工具  | 工具的执行结果   |

### 资源

Resources provide <span style="color: blue; font-weight: bold;">structured access to information</span> that <span style="color: blue; font-weight: bold;">the AI application can retrieve and provide to models as context.</span>

两点：

1. 结构化的信息
2. 构成模型的上下文

应用程序可以直接访问资源并决定如何使用资源，使用方式可以是：

1. 选择部分相关内容
2. 向量化后检索
3. 一次性全部交给模型

资源支持两种发现模式：

1. 直接资源：特定数据的固定 URI，如 calendar://events/2024
2. 资源模板：动态 URI，如 travel://activities/{city}/{category}

资源模板包含的元信息，如：

1. 标题
2. 描述
3. MIME

MCP 约定的与资源相关的操作有：

| 方法                       | 目的                       | 返回        |
|--------------------------|--------------------------|-----------|
| resources/list           | 查询可用的直接资源                | 资源数组      |
| resources/templates/list | resources/templates/list | 资源数组      |
| resources/read           | 获取资源内容                   | 带元数据的资源内容 |
| resources/subscribe      | 监控资源变化                   | 确认消息      |

### 提示词

Prompts provide reusable templates. They allow MCP server authors to provide parameterized prompts for a domain, or showcase how to best use the MCP server.

可利用的模板。

MCP 约定的与提示词相关的操作有：

| 方法           | 目的       | 返回      |
|--------------|----------|---------|
| prompts/list | 查询可用的提示词 | 提示词数组   |
| prompts/get  | 获取提示词内容  | 带传参的提示词 |

## 客户端

MCP clients are <span style="color: blue; font-weight: bold;">instantiated by host applications</span> to communicate with particular MCP servers.

MCP 客户端由主机应用程序实例化，用于与特定 MCP 服务器通信。

The host is the application users interact with, while clients are <span style="color: blue; font-weight: bold;">the protocol-level components that enable server connections</span>.

客户端是实现服务器连接的协议级组件。

客户端还可以为服务器提供若干功能：

1. 探询（Elicitation）：在交互中获取用户特定信息，供服务器收集、使用
2. 根目录（Roots）：允许客户端指定<span style="color: blue; font-weight: bold;">服务器应关注的目录</span>，通过协调机制传达预期范围
3. 采样（Sampling）：服务器通过客户端请求 LLM 完成，从而实现代理式工作流程

