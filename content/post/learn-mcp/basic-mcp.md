+++
date = '2026-01-29T15:01:42+08:00'
draft = false
title = 'MCP：基础概念'
categories = ['人工智能', '大语言模型']
tags = ['MCP', '摘录']
toc = true
+++

做项目要用到 MCP，所以对官方的文档进行学习及摘录，不然做出的东西一知半解，不深入。

学习来源：[https://modelcontextprotocol.io](https://modelcontextprotocol.io)

## 什么是 MCP

> MCP (Model Context Protocol) is an open-source standard for connecting AI applications to external systems.

MCP Server 是**连接 AI 应用与外部系统**的开源标准。

通过 MCP，AI 应用可以连接：

1. 数据源（本地文件、数据库）
2. 工具（Shell 命令、搜索引擎）
3. 工作流（文中举例是“特殊提示词”，我理解应该是“<span style="color: blue;">具体任务</span>”）

![](/imgs/learn-mcp/1769670526298.jpg)

> MCP can have a range of benefits.
> - **Developers**: MCP reduces development time and complexity when building, or integrating with, an AI application or agent.
> - **AI applications or agents**: MCP provides access to an ecosystem of data sources, tools and apps which will enhance capabilities and improve the end-user experience.
> - **End-users**: MCP results in more capable AI applications or agents which can access your data and take actions on your behalf when necessary.

- 为开发者降低开发成本、复杂性
- 赋能 AI 应用或智能体
- 终端用户受益，更好地使用 AI 应用或智能体

MCP 提供了一套“连接”的标准，无疑是规范了开发者的开发逻辑，成本和复杂性必然是会降低的。同时，终端用户使用了赋予更多功能的 AI 应用或智能体，自然是提高了生产效率。

## MCP 规范

<span style="color: blue;">JSON-RPC 消息格式</span>、<span style="color: blue;">有状态连接</span>

MCP 为 AI 应用提供了一套标准化的方式：
- 与语言模型共享上下文信息
- 向 AI 系统展示工具和功能
- 构建可组合的集成和工作流

MCP 为 AI 应用提供了一套标准化的方式：
- 与语言模型共享上下文信息
- 向 AI 系统展示工具和功能
- 构建可组合的集成和工作流


<span style="color: blue;">不同于常规网络模型，MCP 客户端的关键特性：</span>
- <span style="color: blue;">归属于主机应用，非单独的应用或进程</span>
- <span style="color: blue;">具备协议适配、能力桥接的功能</span>
- <span style="color: blue;">是主机应用发起连接的具体执行者</span>

规范约定服务端提供以下 3 种特性：
1. 资源（Resources）：用户或 AI 模型的上下文、数据
2. 提示词（Prompts）：用户的模板化的消息和工作流
3. 工具（Tools）：AI 模型可调用的方法

同时，规范约定客户端向服务端提供以下 3 种特性：

1. 采样（Sampling）：支持服务端主动发起智能体式行为，以及递归的大语言模型交互能力
2. 根目录（Roots）：允许服务端主动发起查询，获取可操作的 URI（统一资源标识符）或文件系统边界，明确自身的操作范围
3. 探询（Elicitation）：支持服务端主动发起请求，向终端用户收集 / 获取额外的信息

<span style="color: blue;">允许服务端的请求的发起、调用，及资源的操作</span>

### MCP 关键组件

MCP 包含的关键组件如下：

1. 基础协议：JSON-RPC 消息类型
2. 生命周期管理：连接建立、能力协商、会话控制
3. 授权机制：基于 HTTP 传输的鉴权和授权
4. 服务端特性：资源、提示词、工具
5. 客户端特性：采样、根目录、探询
6. 实用工具：日志功能等

消息类型包括：请求（Request）、响应（Response）和通知（Notification）。

生命周期的各个阶段：

1. 初始化：
   - 确认协议
   - 能力协商
   - 信息共享
2. 运行：
   - 遵循协商确定的协议版本
   - 仅使用经协商成功敲定的能力特性
3. 关闭：
   - 对于 stdio 类型：
     - 关闭子进程的输入流
     - 等待服务端程序退出，或发送 SIGTERM
     - 若上一操作无响应，则发送 SIGKILL
   - 对于 HTTP 类型
     - 关闭 HTTP 连接

### MCP 层次结构

MCP 包含两个层级：

1. 数据层：基于 JSON-RPC 的客户端-服务端通信
2. 传输层：通信的传输机制及通道

MCP 支持两种传输机制：

1. Stdio 传输：使用进程间通信的标准输入输出流，没有网络开销
2. 流式 HTTP 传输

## 核心基本部件

服务端：工具、资源、提示词。

Each primitive type has associated methods for discovery (*/list), retrieval (*/get), and in some cases, execution (tools/call). MCP clients will use the */list methods to discover available primitives. For example, a client can first list all available tools (tools/list) and then execute them.

通过接口暴露服务端特性。

<span style="color: blue;">（计划内）</span>

<span style="color: blue;">任务（实验）：持久的执行包装器，支持对MCP请求的延迟结果检索和状态跟踪（例如，昂贵的计算、工作流自动化、批处理、多步骤操作）</span>

客户端与服务端的交互

1. 客户端发送 method 为 “initialize” 的请求
2. 客户端通过 method 为 “tools/list” 的请求，获取服务端可用的的工具
3. 服务端返回工具的基本信息
   - name
   - title
   - description
   - inputSchema
4. 客户端通过 method 为 "tools/call" 的请求调用服务端的工具
5. 服务端响应结果中的 "result.content" 字段为执行结果
   - result.content 是一个数组
   - 数组中元素的 type 字段标识内容的类型，有 text、json、code、table、list
