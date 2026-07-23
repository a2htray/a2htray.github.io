+++
date = '2026-07-21T16:00:35+08:00'
draft = false
title = 'FastMCP 上下文状态管理机制'
categories = ['人工智能', '大语言模型']
tags = ['MCP', 'FastMCP', 'AI', 'Agent', '智能体']
toc = true
+++

FastMCP: 3.4.4

FastMCP 的状态不是靠"把上一步结果塞回参数"来维持的，而是服务端按会话持有 state。

## 基础示例

先创建一个的 MCP Server 并启动：

```python
# server.py
mcp = FastMCP("上下文状态管理")


@mcp.tool(
    title="FastMCP 介绍",
    description="介绍 FastMCP 是什么，以及如何使用 FastMCP 进行上下文状态管理。",
)
async def intro():
    return "FastMCP 是一个用于构建多轮对话的工具，支持上下文状态管理。"
```

```bash
$ fastmcp run --no-banner --transport streamable-http server.py
[07/21/26 19:02:21] INFO     Starting MCP server '上下文状态管理' with transport 'streamable-http' on                transport.py:361
                             http://127.0.0.1:8000/mcp                                                                               
INFO:     Started server process [29291]
INFO:     Waiting for application startup.
INFO:     Application startup complete.
INFO:     Uvicorn running on http://127.0.0.1:8000 (Press CTRL+C to quit)
```

* `--no-banner` 终端不显示 banner
* `--transport streamable-http` HTTP 长连接双向流式、自动携带 SessionID，完整支持 Session State 会话持久化

### 获取工具信息

获取 MCP Server 工具信息的方式有两种：

1. 使用 FastMCP 的 Client 获取
2. 使用 langchain_mcp_adapters 的 MultiServerMCPClient 获取

方式一：**使用 FastMCP 的 Client 获取**

```python
# agent.py
import asyncio

from fastmcp.client import Client


async def tools_info_with_fastmcp():
    client = Client("http://localhost:8000/mcp")
    async with client:
        tools = await client.list_tools()
        print(f"共 {len(tools)} 个工具：")
        for tool in tools:
            print(f" * Tool: {tool.name}, Description: {tool.description}")


if __name__ == "__main__":
    asyncio.run(tools_info_with_fastmcp())
```

* 实例化 Client 类
* 调用实例的 `list_tools()` 方法获取工具列表

运行效果：

```bash
$ python agent.py
共 1 个工具：
 * Tool: intro, Description: 介绍 FastMCP 是什么，以及如何使用 FastMCP 进行上下文状态管理。
```

方式二：**使用 langchain_mcp_adapters 的 MultiServerMCPClient 获取**




## 源码学习

在源文件 `fastmcp/server/context.py` 中定义了 `Context` 类，用于管理对话上下文状态。

第 83 行的：

```python
_current_context: ContextVar[Context | None] = ContextVar("context", default=None)
```

用来在异步/多协程环境下，安全存储当前请求/会话的 `Context` 对象；每个协程、每个请求拥有独立隔离的 `_current_context`，互不干扰。

第 127 到 133 行：

```python
@contextmanager
def set_context(context: Context) -> Generator[Context, None, None]:
    token = _current_context.set(context)
    try:
        yield context
    finally:
        _current_context.reset(token)
```

* `@contextmanager` 装饰器用于创建上下文管理器，确保在进入和退出上下文时执行相应的逻辑。

