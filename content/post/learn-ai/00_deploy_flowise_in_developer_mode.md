+++
date = '2026-03-17T21:12:33+08:00'
draft = false
title = 'Flowise：本地部署、系统架构、技术实现和应用构建'
categories = ['人工智能', '应用平台']
tags = ['Flowise', '大模型应用开发平台', 'TypeScript']
toc = true
+++

Flowise 是开源的生成式人工智能开发平台，可用于构建**智能体**和**大语言模型工作流**。

Flowise 提供以下 6 个完整的解决方案：

1. 可视化的构建过程
2. 追踪和分析功能
3. 应用评价体系
4. 人为干预、修正
5. 对外 API、命令行工具、SDK
6. 团队和工作空间

## 本地部署

环境说明：

* macOS Tahoe 26.2
* Node v18.17.0
* pnpm 10.30.3

> Node 版本太高，会导致后续 faiss-node 运行脚本失败和构建报错。

出于以下两个原因：

1. 学习 Flowise 的开发技术
2. 二次开发 Flowise

所以我采用下载源码部署的方式。

**下载源码**

```bash
$ git clone git@github.com:FlowiseAI/Flowise.git
```

如果下载过慢，可以使用 `gitee` 的同步仓库：

```bash
$ git clone git@gitee.com:trusted-list/Flowise.git
```

**安装依赖**

```bash
$ pnpm install --registry=https://registry.npmmirror.com
```

> 上述命令在我电脑上报 `CMake executable is not found`，解决过程见 [CMake executable is not found](#cmake-executable-is-not-found)。

**构建项目**

```bash
$ pnpm build
```

> 构建时，又报了一堆 TypeScript 语法上的错误，实在难以改得动，好在有 TRAE 帮助，过程见 [TypeScript 编译错误](#typescript-编译错误)

**复制配置文件**

```bash
$ cp packages/server/.env.example packages/server/.env
$ cp packages/ui/.env.example packages/ui/.env
```

先都采用默认配置。

**开启服务**

```bash
$ pnpm dev
```

访问 [http://localhost:8080](http://localhost:8080)，注册一个账号，登录后如下图：

![](/imgs/learn-ai/deploy_flowise_01.png)

## 系统架构

## 技术实现

## 应用构建

Flowise 支持构建两种智能体，分别是 Chatflow 和 Agentflow，两者的区别在于：

| /    | Chatflow       | Agentflow                 |
|------|----------------|---------------------------|
| 定位   | 单智能体，对话优先的工作流  | 多智能体协作，复杂工作流              |
| 本质   | 一个 LLM 走完整个工作流 | 多个具有独立 LLM/提示词/工具/记忆 的智能体 |
| 适用场景 | 聊天机器人，RAG 问答   | 多任务自动化，多角色协同，长时状态工作流      |


### Chatflow

Flowise 本地化部署已经完成，试着搭建一个 Chatflow，示例相对简单，使用到了 3 个节点：

* Chat Models 下的 ChatOllama 节点
* Prompts 下的 Chat Prompt Template 节点
* Chains 下的 LLM Chain 节点

连接方式如下：

![](/imgs/learn-ai/deploy_flowise_03.png)

> Ollama 的本地部署可以参看的我 [Ollama：部署与使用](/post/learn-ai/01_deploy_ollama_on_local/) 一文。

进行简单的问答，体验效果不是很好，主要是回复得非常慢，可能是本地版本或 Ollama 的问题，但起码有一次完整的问答过程，如下图：

![](/imgs/learn-ai/deploy_flowise_04.png)

#### 接口实现

问答时，请求的接口是 `/api/v1/internal-prediction`，路由文件和控制器文件分别是：

* 路由文件：packages/server/src/routes/internal-predictions/index.ts

```typescript
router.post(['/', '/:id'], internalPredictionsController.createInternalPrediction)
```

* 控制器文件：packages/server/src/controllers/internal-predictions/index.ts

查看请求传参：

```json
{"question":"你好","chatId":"7360cfc0-6f7f-4faa-9562-6bec84db74c2","streaming":true}
```

服务器端最终调用的是 `createAndStreamInternalPrediction` 函数，产生流式的响应。

代码逻辑如下：

* 获取 App 实例的 sseStreamer 属性，SSEStreamer 用于管理 SSE 的连接信息、发送流式消息等功能
* 调用 sseStreamer.addClient 方法建立 chatId 和 response 的对应关系
* 设置 response 流式响应头
* 调用 utilBuildChatflow 函数，主要的逻辑在这一块
* 最终，解除 chatId 和 response 的关系，并消除相关信息

### Agentflow

## 遇到问题

本小节罗列在使用 Flowise 过程中遇到的问题，因非重点，所以不做过多描述，仅记录解决问题的过程。

### CMake executable is not found

```bash
$ curl -OL https://github.com/Kitware/CMake/releases/download/v4.3.0/cmake-4.3.0-macos-universal.tar.gz
$ tar -zxvf cmake-4.3.0-macos-universal.tar.gz
$ sudo mv cmake-4.3.0-macos-universal /usr/local/cmake-4.3.0
```

编辑 `~/.zshrc`，新增一行：

```shell
export PATH="/usr/local/cmake-4.3.0/CMake.app/Contents/bin:$PATH"
```

```bash
$ source ~/.zshrc
$ cmake --version
cmake version 4.3.0

CMake suite maintained and supported by Kitware (kitware.com/cmake).
```

### TypeScript 编译错误

把终端的错误信息给到 TRAE，让其解决，主要添加了 Polyfill 和修改了代码书写方式。

![](/imgs/learn-ai/deploy_flowise_01.png)

## 资源

* [官方文档 - Introduction](https://docs.flowiseai.com/)
* [精讲 AI 教程: 免费使用 Flowise 搭建 LLM 工作流应用](https://zhuanlan.zhihu.com/p/688195105)
