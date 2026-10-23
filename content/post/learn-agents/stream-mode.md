+++
date = '2026-10-23T09:35:50+08:00'
draft = false
title = 'LangChain/LangGraph 流式输出：stream_mode 不同值的效果 1'
categories = ['人工智能', '智能体']
tags = ['AI Agent', 'LangGraph', 'LangChain', 'Stream Mode']
toc = true
+++

> 演示 stream_mode 为 values、updates、messages 的效果。

![](/imgs/learn-agents/20261023.png)

## vs. invoke

`invoke` 方法直接返回模型的最终结果，屏蔽了模型内部的工作流程，`stream` 方法则会返回一个迭代器，逐步 yield 中间状态。

什么时候用 invoke，什么时候用 stream：

- 用 invoke：只关心最终结果、任务短、或写离线批处理脚本
- 用 stream：你要看进度、要流式输出 LLM 文本、或构建交互式 UI

## 构建一个带查询当前时间的智能体

```python
from datetime import datetime

from langchain.agents import create_agent
from langchain.tools import tool
from langchain_ollama import ChatOllama

llm = ChatOllama(model="qwen3.5:9b")


@tool
def current_time():
    """返回当前时间"""
    now = datetime.now()
    return now.strftime("%Y-%m-%d %H:%M:%S")


agent = create_agent(llm, tools=[current_time])
```

在调用智能体时，期望步骤是：用户输入 -> 智能决定调用工具 -> 工具调用 -> 智能体总结工具输出。

## 演示 stream_mode=values 的效果

指定 `stream_mode=values` 后，模型会流式输出模型的中间结果：

```python
def get_types(messages: list) -> list:
    """返回消息类型"""
    return [message.__class__.__name__ for message in messages]


for i, chunk in enumerate(agent.stream({"messages": [{"role": "user", "content": "现在是什么时候"}]}, stream_mode="values")):
    print(f"================= Chunk {i+1} =================")
    message_types = get_types(chunk["messages"])
    last_message_type = message_types[-1]
    last_message = chunk["messages"][-1]
    print(f"消息类型：{message_types}，最新一条消息类型为 {last_message_type}")

    if last_message_type == "HumanMessage":
        print(f"用户输入：{last_message.content}")
    elif last_message_type == "AIMessage":
        print(f"输入 token：{last_message.usage_metadata['input_tokens']}，输出 token：{last_message.usage_metadata['output_tokens']}")
        if last_message.tool_calls:
            for tool_call in last_message.tool_calls:
                print(f"模型决定调用工具 {tool_call['name']}")
        else:
            print(f"模型回复：{last_message.content}")
    elif last_message_type == "ToolMessage":
        print(f"调用工具 {last_message.name}，返回 {last_message.content}")
    else:
        print(f"未知消息类型 {last_message_type}")
    print()
```

* `get_types` 返回消息列表中的消息类型，针对不同的消息类型，打印不同的信息，消息类型包括 `HumanMessage`、`AIMessage`、`ToolMessage`。
* `chunk["messages"][-1]` 为最新一条消息
    * 如果是 `HumanMessage`，打印用户输入
    * 如果是 `AIMessage`，打印模型回复工具调用决策
    * 如果是 `ToolMessage`，打印调用工具和返回结果
* 消息带 `tool_calls`，表示模型要求调用工具

执行结果如下：

```bash
$ python stream_mode_values.py
================= Chunk 1 =================
消息类型：['HumanMessage']，最新一条消息类型为 HumanMessage
用户输入：现在是什么时候

================= Chunk 2 =================
消息类型：['HumanMessage', 'AIMessage']，最新一条消息类型为 AIMessage
输入 token：253，输出 token：39
模型决定调用工具 current_time

================= Chunk 3 =================
消息类型：['HumanMessage', 'AIMessage', 'ToolMessage']，最新一条消息类型为 ToolMessage
调用工具 current_time，返回 2026-10-22 20:26:47

================= Chunk 4 =================
消息类型：['HumanMessage', 'AIMessage', 'ToolMessage', 'AIMessage']，最新一条消息类型为 AIMessage
输入 token：303，输出 token：73
模型回复：现在是 2026 年 10 月 22 日 20:26:47。
```

直观的感觉就是，每一步的模型状态是一个快照，包含了 `messages` 的所有信息，最新的模型消息在列表的最后。

## 演示 stream_mode=updates 的效果

指定 `stream_mode=updates` 后，模型会流式输出模型的中间结果：

```python
def get_types(messages: list) -> list:
    """返回消息类型"""
    return [message.__class__.__name__ for message in messages]


for i, chunk in enumerate(agent.stream({"messages": [{"role": "user", "content": "现在是什么时候"}]}, stream_mode="updates")):
    print(f"================= Chunk {i+1} =================")

    messages = chunk["model"]["messages"] if "model" in chunk else chunk["tools"]["messages"]
    message_types = get_types(messages)
    last_message_type = message_types[-1]
    last_message = messages[-1]

    print(f"消息类型：{message_types}，最新一条消息类型为 {last_message_type}")

    if "model" in chunk:
        if last_message.tool_calls:
            print(f"模型决定调用工具 {last_message.tool_calls[0]['name']}")
        else:
            print(f"模型回复：{last_message.content}")
    if "tools" in chunk:
        print(f"调用工具 {last_message.name}，返回 {last_message.content}")

    print()
```

```bash
$ python stream_mode_updates.py
================= Chunk 1 =================
消息类型：['AIMessage']，最新一条消息类型为 AIMessage
模型决定调用工具 current_time

================= Chunk 2 =================
消息类型：['ToolMessage']，最新一条消息类型为 ToolMessage
调用工具 current_time，返回 2026-10-23 10:36:53

================= Chunk 3 =================
消息类型：['AIMessage']，最新一条消息类型为 AIMessage
模型回复：现在是 **2026年10月23日 10:36:53**。
```

每一步的输出，都是最新的一个消息（增量更新）。

## 演示 stream_mode=messages 的效果

```python
for i, (message, metadata) in enumerate(agent.stream({"messages": [{"role": "user", "content": "现在是什么时候"}]}, stream_mode="messages")):
    print(f"================= Chunk {i+1} =================")
    message_class_name = message.__class__.__name__
    print(f"消息类型：{message_class_name}, 元数据 Node：{metadata['langgraph_node']}")

    if message_class_name == "AIMessageChunk":
        print(f"模型回复：{message.content}")
    elif message_class_name == "ToolMessage":
        print(f"调用工具 {message.name}，返回 {message.content}")

    print()
```

```bash
$ python stream_mode_messages.py
================= Chunk 1 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 2 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 3 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 4 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 5 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 6 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 7 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 8 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 9 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 10 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 11 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 12 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 13 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 14 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 15 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 16 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 17 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 18 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 19 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 20 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 21 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 22 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 23 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 24 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 25 =================
消息类型：ToolMessage, 元数据 Node：tools
调用工具 current_time，返回 2026-10-23 10:50:32

================= Chunk 26 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 27 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 28 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 29 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 30 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 31 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 32 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 33 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 34 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 35 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 36 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 37 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 38 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 39 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 40 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 41 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 42 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 43 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 44 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 45 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 46 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 47 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 48 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 49 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 50 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 51 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 52 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 53 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 54 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 55 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 56 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 57 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 58 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 59 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 60 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 61 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 62 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 63 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 64 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 65 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 66 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 67 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 68 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：现在是

================= Chunk 69 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复： 2

================= Chunk 70 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：0

================= Chunk 71 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：2

================= Chunk 72 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：6

================= Chunk 73 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：年

================= Chunk 74 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：1

================= Chunk 75 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：0

================= Chunk 76 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：月

================= Chunk 77 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：2

================= Chunk 78 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：3

================= Chunk 79 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：日

================= Chunk 80 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复： 1

================= Chunk 81 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：0

================= Chunk 82 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：:

================= Chunk 83 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：5

================= Chunk 84 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：0

================= Chunk 85 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：:

================= Chunk 86 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：3

================= Chunk 87 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：2

================= Chunk 88 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：。

================= Chunk 89 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：

================= Chunk 90 =================
消息类型：AIMessageChunk, 元数据 Node：model
模型回复：
```

输出最为详细，逐 token 输出，用于打字机效果。