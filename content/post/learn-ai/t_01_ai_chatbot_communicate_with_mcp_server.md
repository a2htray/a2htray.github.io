+++
date = '2026-07-20T15:15:14+08:00'
draft = false
title = 'Chatbot UI 与 MCP Server 通信'
categories = ['人工智能', '大语言模型']
tags = ['MCP', '翻译', 'Chatbot','T 系列', 'AI']
toc = true
+++

本文是 T 系列的第 01 篇，主要介绍 Chatbot UI 如何与 MCP Server 进行通信。

> T 系列为技术文章翻译系列，加以本人的一些拙见，如有雷同，实属正常，如有不同，欢迎交流。

Chatbot 的用户界面（UI）正从简单的文本显示演变为一个动态客户端，能够与复杂的后端系统进行通信，而该后端通常由模型上下文协议（MCP）Agent 驱动。本文探讨了基于工具的聊天界面架构，**重点分析前端与后端之间的通信流程**。我们将逐步介绍用户请求的路由、调用外部工具以及实时流式响应的全过程。目标是清晰理解支持这些复杂交互的协议和设计模式，从而超越基础的请求-响应模型，构建出更稳健、可扩展且更易用的系统。

<font color="blue">关注点：前端与后端之间的通信 - 用户请求 -> 工具调用 -> 实时响应。</font>

## 前端-后端通信流程

工具支持的 Chatbot 的一个关键架构转变是职责分离。前端负责处理用户输入并渲染响应，后端（MCP Server）则协调 Agent 的逻辑，管理工具上下文，并与外部工具进行交互。这两层之间的通信必须具备高可靠性、高效性，能够支持多轮对话和实时更新。

<font color="blue">现实中的通信流程应该更为复杂，如后端可以进一步拆分：后端可自启 MCP Server 或调用外部 MCP Server。Agent 负责理解用户输入、生成工具调用、调用工具、解释调用结果，后端或 MCP 框架维护用户对话的上下文，并将结果返回给前端。</font>

用户请求到 Agent：

当用户输入消息并按下回车键时，用户界面会向后端发送请求。这通常是一个**异步调用**到服务器端点。后端的任务是：

1. 接收用户的消息
2. 将消息路由至 MCP Agent
3. 启动 MCP 流程，该流程可能涉及一个或多个工具调用

这是一个前端如何通过简单的 fetch API 调用向后端发送请求的基本示例：

```typescript
// src/api/chat.ts
interface ChatMessage {
  role: 'user' | 'assistant' | 'tool';
  content: string;
}

export async function sendMessage(message: string): Promise<ChatMessage[]> {
  try {
    const response = await fetch('/api/chat', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ message }),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data: ChatMessage[] = await response.json();
    return data;
  } catch (error) {
    console.error("Failed to send message:", error);
    throw error;
  }
}
```

在后端，一个无状态的函数调用或专用的 API 端点会接收此请求。该端点内的核心逻辑随后将触发 MCP Agent。

<font color="blue">上述过程与 HTTP 请求无异，无非是前端通过网络发送请求，后端接收并处理请求，再返回响应。</font>

```typescript
// src/server/api/chat.ts
import { mcpAgent } from '../mcp/agent';

export async function handleChatRequest(req: Request): Promise<Response> {
  const { message } = await req.json();
  const agentResponse = await mcpAgent.handleRequest(message);

  return new Response(JSON.stringify(agentResponse), {
    headers: { 'Content-Type': 'application/json' },
  });
}
```

这种简单的请求-响应模型适用于基本的交互，但在处理**多步骤工具调用和长时间运行的操作**时则显得力不从心。用户只能等待，而无法获得代理进程的任何反馈。

<font color="blue">上述前后端的调用过程，对于多步骤工具调用和长耗时操作，用户只能等待响应，从而对用户体验造成影响。</font>

## 流式实时反馈

为了提供更好的用户体验，现代聊天界面采用流式协议。这使得后端能够随着数据可用时发送部分响应，从而让用户实时了解 Agent 的进展情况。实现这一功能的两种常见协议是服务器发送事件（SSE）和 WebSocket。

<font color="blue">流式响应解决了长耗时操作的用户体验问题。</font>

**使用 SSE**

SSE 是一种**轻量级的单向协议**，通过标准 HTTP 连接实现服务器与客户端之间的实时通信。它非常适合需要显示代理进度的聊天界面，例如“正在搜索航班信息……”或“正在为您预订……”。前端订阅一个 SSE 端点，服务器在事件发生时将消息推送到客户端。

在前端，UI订阅事件流：

```typescript
// src/components/ChatStreamComponent.tsx
import { useEffect, useRef } from 'react';

const ChatStreamComponent = ({ userId, message }) => {
  const eventSourceRef = useRef<EventSource | null>(null);

  useEffect(() => {
    eventSourceRef.current = new EventSource(`/api/chat-stream?userId=${userId}&message=${encodeURIComponent(message)}`);

    eventSourceRef.current.onmessage = (event) => {
      const data = JSON.parse(event.data);
      console.log('Received streamed data:', data);
    };

    eventSourceRef.current.onerror = (error) => {
      console.error('SSE Error:', error);
      eventSourceRef.current?.close();
    };

    return () => {
      eventSourceRef.current?.close();
    };
  }, [userId, message]);
};
```

<font color="blue">前端使用 EventSource 实例，并为其添加消息（onmessage）和错误处理（onerror）事件监听器。</font>

在后端，服务器会保持连接打开，并随着 MCP Agent 完成各个步骤而发送事件。

```typescript
// src/server/api/chat-stream.ts
import { mcpAgent } from '../mcp/agent';
import { AgentProgressEvent } from '../mcp/types'; // Define event types

export async function handleStreamRequest(req: Request): Promise<Response> {
  const { userId, message } = req.query;
  // TextEncoder：二进制编码器，把字符串转 UTF-8 字节流，用于写入流式响应
  const encoder = new TextEncoder();

  const readableStream = new ReadableStream({
    async start(controller) {
      // data: [JSON字符串]\n\n（双换行是 SSE 消息分隔符，前端 EventSource 靠它分割每条消息）
      controller.enqueue(encoder.encode(`data: ${JSON.stringify({ type: 'status', content: 'Thinking...' })}\n\n`));

      const progressEvents = mcpAgent.run(userId, message);

      for await (const event of progressEvents) {
        if (event.type === 'tool_call') {
          const data: AgentProgressEvent = { 
            type: 'tool_call_event', 
            tool_name: event.tool_name, 
            status: 'in-progress' 
          };
          controller.enqueue(encoder.encode(`data: ${JSON.stringify(data)}\n\n`));
        } else if (event.type === 'tool_result') {
          const data: AgentProgressEvent = { 
            type: 'tool_result_event', 
            tool_name: event.tool_name, 
            status: 'completed', 
            result: event.result
          };
          controller.enqueue(encoder.encode(`data: ${JSON.stringify(data)}\n\n`));
        }
      }

      const finalResponse = await progressEvents.getFinalResponse();
      controller.enqueue(encoder.encode(`data: ${JSON.stringify({ type: 'final_response', content: finalResponse })}\n\n`));
      controller.close();
    }
  });

  return new Response(readableStream, {
    headers: {
      'Content-Type': 'text/event-stream',
      'Cache-Control': 'no-cache',
      'Connection': 'keep-alive',
    },
  });
}
```

<font color="blue">前端发起请求后，后端持续实时推送 AI Agent 执行过程中的中间状态（工具调用、工具返回结果），最后推送完整最终回答，典型用于 AI 对话实时打字流、Agent 工具调用进度展示。</font>

这种流式显示模式使聊天界面能够展示对话流程中的中间步骤，特别适用于复杂的多轮工具交互场景，用户在这些情况下可能需要等待较长时间而得不到任何反馈。

<font color="blue">使用 SSE 改善用户对话体验。</font>

**幕后：MCP Agent 的角色**

MCP Agent 是系统的核心，它不仅处理用户输入、提示词，还根据明确的状态协调一系列操作。Agent 的流程包括：

1. 请求解析：Agent 接收到用户的消息以及当前的工具上下文
2. 工具选择：根据消息和上下文，Agent 决定需要使用哪个工具（如有）以及执行何种操作
3. 工具调用生成：生成一个结构化的工具调用对象
4. 工具调用执行：**MCP 框架执行该工具调用**，并等待工具结果返回
5. 上下文更新：将工具结果添加到工具上下文中
6. 响应生成：代理利用更新后的上下文决定下一步行动：生成面向用户的响应，或发起新的工具调用

<font color="blue">工具选择的本质是 LLM 驱动的函数调用，MCP Server 工具的结构化 schema（名称、描述、输入参数及类型等）会被拼接成 LLM 的提示词语，让模型知道“哪些工具可用及如何调用”。</font>

```bash
用户消息 + 工具清单 + 历史上下文
        ↓  LLM 推理
判断是否需要工具 → 是 → 选哪个工具 + 填什么参数（生成结构化调用对象）
                   ↓ 否
              直接生成回答
```

多步骤、有状态的过程使得基于工具的复杂工作流成为可能。聊天界面凭借其流式传输功能，充当了这一动态过程的窗口，**实时显示每个步骤的执行情况**。这种透明性对于管理用户期望以及调试问题至关重要。

## 交互性与透明性的设计

MCP Agent 聊天界面设计的一个关键部分是**让用户能够清晰地看到内部流程**。仅仅接收文本片段并显示为一条连续消息是不够的。界面应通过视觉方式呈现代理的操作，以建立信任并提供明确的反馈。

* 操作提示：使用图标或标签等视觉提示来显示代理正在执行特定操作，例如搜索工具的“🔍”图标或日历工具的“📅”图标
* 进度消息：当流式事件到达时，显示简短且信息明确的消息，如“正在搜索航班……”或“从数据库中获取用户数据……”
* 上下文提示：当收到工具结果时，界面可对输出内容进行不同格式化（例如代码块或结构化卡片），以清晰区分其与代理的文本描述。这在展示从外部API或数据库获取的数据时尤为有用

<font color="blue">让用户通过聊天界面知道当前系统正在做什么。</font>

这种设计哲学将聊天体验从简单的对话转变为协作过程，让用户成为 Agent 解决问题过程中的知情者、观察者及参与者。

## 有感而发

转向基于工具的 Chatbot 架构，要求重新评估用户界面设计和通信协议。在所有工具调用完成后仅返回最终响应，对于复杂任务而言是一种糟糕的用户体验。这会让用户感到迷茫，不知系统是否正常运行或已卡住。诸如 SSE（服务器端事件）之类的流式传输协议通过弥合代理内部处理过程与用户感知之间的差距来解决这一问题。这种透明性有助于建立用户信任，并使复杂交互感觉更加自然。

此外，一个设计良好的用户界面不仅应显示最终响应，还应提供使用工具的视觉提示。例如，显示特定于工具调用的加载指示器，可以使代理的操作更加透明。这些界面的未来在于前端与 MCP 协议之间的更紧密集成。用户界面可以变得更加主动，根据用户的输入，在代理决定使用某个工具之前，就向用户推荐该工具。这将使用户界面从被动的展示转变为工具支持对话中的积极参与者，为真正直观且功能强大的对话系统铺平道路。

<font color="blue">1. 请求由异步调用改成实时流式响应；2. MCP Agent 完成工具调用的流程；3. 考虑如何把工具调用流程中的中间状态展示给用户。</font>