+++
date = '2026-10-22T15:45:48+08:00'
draft = false
title = '学习 Langchain Langgraph：代码也给你'
categories = ['人工智能', '智能体']
tags = ['AI Agent', 'LangGraph', 'LangChain']
toc = true
+++

本文演示跑通 LangChain、LangGraph 相关概念的最小代码示例。

## C1 创建智能体

```python
from langchain.agents import create_agent
from langchain_ollama import OllamaLLM

llm = OllamaLLM(model="qwen3.5:9b")

agent = create_agent(llm, tools=[])

result = agent.invoke({"messages": [{"role": "user", "content": "你好"}]})
print(result["messages"][-1].content)
```

* LLM 可以换成其它模型，因为我本地部署了 Ollama，所以这里用 qwen3.5:9b 模型
* 步骤：
    * `create_agent` 创建智能体
    * `agent.invoke` 调用智能体

## C2 带工具调用的智能体

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

result = agent.invoke({"messages": [{"role": "user", "content": "现在是什么时候"}]})
print(result["messages"][-1].content)
```

* @tool 装饰器定义一个工具
* 在 create_agent 中传入工具列表

## C3 带工具调用的 LangGraph 智能体

```python
from langgraph.graph import StateGraph, MessagesState, START, END
from langgraph.prebuilt import ToolNode, tools_condition
from langchain.tools import tool
from langchain_ollama import ChatOllama

llm = ChatOllama(model="qwen3.5:9b")


@tool
def add(a: int, b: int) -> int:
    """对两个整型数值做加法运算并返回结果。"""
    return a + b


@tool
def multiply(a: int, b: int) -> int:
    """对两个整型数值做乘法运算并返回结果。"""
    return a * b


tools = [add, multiply]
llm_with_tools = llm.bind_tools(tools)


def agent(state: MessagesState) -> dict:
    """调用 LLM（已绑定工具），返回新的 AIMessage。"""
    response = llm_with_tools.invoke(state["messages"])
    return {"messages": [response]}


builder = StateGraph(MessagesState)

builder.add_node("agent", agent)
builder.add_node("tools", ToolNode(tools))
builder.add_edge(START, "agent")
builder.add_conditional_edges("agent", tools_condition, {"tools": "tools", END: END})
builder.add_edge("tools", "agent")

graph = builder.compile()


if __name__ == "__main__":
    user_input = "请先算 3 乘以 4，再加 5，最终等于多少？请逐步使用工具计算。"
    result = graph.invoke({"messages": [{"role": "user", "content": user_input}]})

    for msg in result["messages"]:
        msg.pretty_print()

    print(graph.get_graph().draw_mermaid())
```

* 构建步骤如下：
    1. 定义工具：@tool
    2. 定义智能体：agent 函数
    3. 定义工具节点：ToolNode
    4. 定义边：add_edge, add_conditional_edges
    5. 编译图：builder.compile()
    6. 调用图：graph.invoke()

## C4 短期记忆 checkpointer

```python
from langchain_ollama import ChatOllama
from langgraph.graph import StateGraph, MessagesState, START, END
from langgraph.checkpoint.memory import InMemorySaver
from langgraph.config import RunnableConfig

llm = ChatOllama(model="qwen3.5:9b")

checkpointer = InMemorySaver()


def agent(state: MessagesState) -> dict:
    response = llm.invoke(state["messages"])
    return {"messages": [response]}


builder = StateGraph(MessagesState)

builder.add_node("agent", agent)
builder.add_edge(START, "agent")
builder.add_edge("agent", END)

graph = builder.compile(checkpointer=checkpointer)

if __name__ == "__main__":
    cfg = RunnableConfig(configurable={"thread_id": "user-001"})
    result = graph.invoke({"messages": [{"role": "user", "content": "我是小明，你是谁"}]}, config=cfg)
    print("------ A1 ------")
    print(result["messages"][-1].content)
    result = graph.invoke({"messages": [{"role": "user", "content": "我是谁"}]}, config=cfg)
    print("------ A2 ------")
    print(result["messages"][-1].content)
```

* checkpointer，适用于单个会话（thread）的短期记忆
* 在 `graph.invoke` 中传入 `config=cfg`，即可实现短期记忆

## C5 长期记忆 store

```python
from langchain.messages import SystemMessage
from langchain.tools import tool
from langchain_ollama import ChatOllama
from langgraph.checkpoint.memory import InMemorySaver
from langgraph.graph import StateGraph, MessagesState, START, END
from langgraph.prebuilt import ToolNode, tools_condition
from langgraph.store.memory import InMemoryStore

llm = ChatOllama(model="qwen3.5:9b")

checkpointer = InMemorySaver()   # 短期记忆：每个 thread_id 一份对话历史
store = InMemoryStore()          # 长期记忆：跨 thread_id 共享

NAMESPACE = ("user_profile",)


@tool
def save_memory(key: str, value: str) -> str:
    """保存一条长期记忆。key 为唯一标识（如 name/job/hobby），value 为具体内容。
    当用户透露姓名、职业、爱好、偏好等个人信息时调用本工具。"""
    store.put(NAMESPACE, key, {"value": value})
    return f"已保存长期记忆: {key} = {value}"


@tool
def get_memory(key: str) -> str:
    """根据 key 读取之前保存的长期记忆；若没有记录则返回未找到。
    当用户询问关于他本人的信息时调用本工具回忆。"""
    item = store.get(NAMESPACE, key)
    if item is None:
        return f"未找到关于 {key} 的长期记忆"
    return item.value.get("value", "")


tools = [save_memory, get_memory]
llm_with_tools = llm.bind_tools(tools)

SYSTEM_PROMPT = (
    "你是一个有长期记忆的助手。\n"
    "规则：\n"
    "1. 当用户告诉你他的姓名、职业、爱好、偏好等个人信息时，务必调用 save_memory 保存。\n"
    "2. 当用户询问关于他本人的信息（如'我是谁''我做什么''我喜欢什么'）时，"
    "务必调用 get_memory 回忆后再回答，不要编造。\n"
)


def agent(state: MessagesState) -> dict:
    messages = [SystemMessage(content=SYSTEM_PROMPT)] + state["messages"]
    response = llm_with_tools.invoke(messages)
    return {"messages": [response]}


builder = StateGraph(MessagesState)
builder.add_node("agent", agent)
builder.add_node("tools", ToolNode(tools))

builder.add_edge(START, "agent")
builder.add_conditional_edges("agent", tools_condition, {"tools": "tools", END: END})
builder.add_edge("tools", "agent")                       

graph = builder.compile(checkpointer=checkpointer, store=store)


if __name__ == "__main__":
    config_a = {"configurable": {"thread_id": "thread-A"}}
    graph.invoke(
        {"messages": [{
            "role": "user",
            "content": "我叫王明，是一名数据分析师，平时喜欢爬山，最爱用的语言是 Python。",
        }]},
        config_a,
    )

    print("=== 长期记忆 store（跨 thread 可见）===")
    for item in store.search(NAMESPACE):
        print(f"  {item.key}: {item.value.get('value')}")

    config_b = {"configurable": {"thread_id": "thread-B"}}
    result = graph.invoke(
        {"messages": [{
            "role": "user",
            "content": "你还记得我是谁吗？我是做什么工作的，又喜欢什么？",
        }]},
        config_b,
    )

    print("\n=== 会话 B 的回答（依赖长期记忆）===")
    print(result["messages"][-1].content)

    print("\n=== 对比说明 ===")
    print("· checkpointer(短期): 按 thread_id 隔离，thread-B 看不到 thread-A 的对话历史")
    print("· store(长期)       : 跨 thread_id 共享，thread-B 仍能通过 get_memory 回忆 thread-A 保存的信息")
```

* NAMESPACE：用于区分不同的长期记忆空间，避免冲突
* @tool 定义两个工具：save_memory 和 get_memory，用于保存和回忆长期记忆
* builder.compile(checkpointer=checkpointer, store=store)：编译图时传入 checkpointer 和 store，即可实现短期记忆和长期记忆

## C6 Human-in-the-Loop

```python
from langchain.messages import SystemMessage
from langchain_core.messages import ToolMessage
from langchain.tools import tool
from langchain_ollama import ChatOllama
from langgraph.checkpoint.memory import InMemorySaver
from langgraph.graph import StateGraph, MessagesState, START, END
from langgraph.prebuilt import ToolNode
from langgraph.types import interrupt, Command

llm = ChatOllama(model="qwen3.5:9b")

checkpointer = InMemorySaver()


@tool
def send_email(to: str, subject: str) -> str:
    """发送一封邮件（敏感操作，需要人工确认）。to=收件人，subject=主题。"""
    return f"邮件已发送给 {to}，主题：{subject}"


@tool
def search_web(query: str) -> str:
    """搜索网页资料（普通操作，直接执行，无需人工确认）。"""
    return f"关于 '{query}' 的搜索结果摘要……"


tools = [send_email, search_web]
llm_with_tools = llm.bind_tools(tools)

SENSITIVE_TOOLS = {"send_email"}

SYSTEM_PROMPT = (
    "你是一个助手。当用户要发邮件时调用 send_email；当用户要查资料时调用 search_web。"
)


def agent(state: MessagesState) -> dict:
    messages = [SystemMessage(content=SYSTEM_PROMPT)] + state["messages"]
    response = llm_with_tools.invoke(messages)
    return {"messages": [response]}


def human_approval(state: MessagesState) -> dict:
    last = state["messages"][-1]
    tool_call = last.tool_calls[0]

    decision = interrupt({
        "action": "请人工确认是否执行该敏感操作",
        "tool": tool_call["name"],
        "args": tool_call["args"],
    })

    if decision == "approve":
        return ToolNode(tools).invoke(state)
    else:
        return {"messages": [ToolMessage(
            content="用户拒绝了该敏感操作，未执行。",
            tool_call_id=tool_call["id"],
        )]}


def route_after_agent(state: MessagesState) -> str:
    last = state["messages"][-1]
    if getattr(last, "tool_calls", None):
        name = last.tool_calls[0]["name"]
        return "human_approval" if name in SENSITIVE_TOOLS else "tools"
    return END


builder = StateGraph(MessagesState)
builder.add_node("agent", agent)
builder.add_node("human_approval", human_approval)
builder.add_node("tools", ToolNode(tools))

builder.add_edge(START, "agent")
builder.add_conditional_edges("agent", route_after_agent)
builder.add_edge("tools", "agent")
builder.add_edge("human_approval", "agent")

graph = builder.compile(checkpointer=checkpointer)


def run_with_human(thread_id: str, user_req: str, human_decision: str) -> None:
    config = {"configurable": {"thread_id": thread_id}}

    result = graph.invoke({"messages": [{"role": "user", "content": user_req}]}, config)

    if "__interrupt__" in result:
        req = result["__interrupt__"][0].value
        print(f"[线程 {thread_id}] ⏸ 人工审核请求: {req}")
        print(f"[线程 {thread_id}]   人工决定: {human_decision}")
        result = graph.invoke(Command(resume=human_decision), config)
    else:
        print(f"[线程 {thread_id}] (未触发人工审核)")

    print(f"[线程 {thread_id}] 最终回答: {result['messages'][-1].content}\n")


if __name__ == "__main__":
    print(graph.get_graph().draw_mermaid())
    run_with_human(
        "demo-approve",
        "帮我给 lisi@company.com 发一封邮件，主题是：项目进度同步",
        "approve",
    )

    run_with_human(
        "demo-reject",
        "帮我给 boss@company.com 发邮件，主题是：请假申请",
        "reject",
    )
```

* interrupt 获取人工审核请求的结果
* Command 用于继续执行图

## C7 终端交互 Human-in-the-Loop

```python
from langchain.messages import SystemMessage
from langchain_core.messages import ToolMessage
from langchain.tools import tool
from langchain_ollama import ChatOllama
from langgraph.checkpoint.memory import InMemorySaver
from langgraph.graph import StateGraph, MessagesState, START, END
from langgraph.prebuilt import ToolNode
from langgraph.types import interrupt, Command

llm = ChatOllama(model="qwen3.5:9b")

checkpointer = InMemorySaver()


@tool
def send_email(to: str, subject: str) -> str:
    """发送一封邮件（敏感操作，需要人工确认）。to=收件人，subject=主题。"""
    return f"邮件已发送给 {to}，主题：{subject}"


@tool
def search_web(query: str) -> str:
    """搜索网页资料（普通操作，直接执行，无需人工确认）。"""
    return f"关于 '{query}' 的搜索结果摘要……"


tools = [send_email, search_web]
llm_with_tools = llm.bind_tools(tools)

SENSITIVE_TOOLS = {"send_email"}

SYSTEM_PROMPT = (
    "你是一个助手。当用户要发邮件时调用 send_email；当用户要查资料时调用 search_web。"
)


def agent(state: MessagesState) -> dict:
    messages = [SystemMessage(content=SYSTEM_PROMPT)] + state["messages"]
    response = llm_with_tools.invoke(messages)
    return {"messages": [response]}


def human_approval(state: MessagesState) -> dict:
    last = state["messages"][-1]
    tool_call = last.tool_calls[0]

    decision = interrupt({
        "action": "请人工确认是否执行该敏感操作",
        "tool": tool_call["name"],
        "args": tool_call["args"],
    })

    if decision == "approve":
        return ToolNode(tools).invoke(state)
    else:
        return {"messages": [ToolMessage(
            content="用户拒绝了该敏感操作，未执行。",
            tool_call_id=tool_call["id"],
        )]}


def route_after_agent(state: MessagesState) -> str:
    last = state["messages"][-1]
    if getattr(last, "tool_calls", None):
        name = last.tool_calls[0]["name"]
        return "human_approval" if name in SENSITIVE_TOOLS else "tools"
    return END


builder = StateGraph(MessagesState)
builder.add_node("agent", agent)
builder.add_node("human_approval", human_approval)
builder.add_node("tools", ToolNode(tools))

builder.add_edge(START, "agent")
builder.add_conditional_edges("agent", route_after_agent)
builder.add_edge("tools", "agent")
builder.add_edge("human_approval", "agent")

graph = builder.compile(checkpointer=checkpointer)


def run_interactive() -> None:
    config = {"configurable": {"thread_id": "interactive"}}
    print("=== 人机协作会话（输入 exit / quit 退出）===")
    print("提示：涉及发邮件(send_email)会触发人工审批，请在终端输入 approve / reject\n")

    while True:
        try:
            user_req = input("你: ").strip()
        except (EOFError, KeyboardInterrupt):
            print("\n再见。")
            break
        if user_req.lower() in {"exit", "quit", "q"}:
            print("再见。")
            break
        if not user_req:
            continue

        result = graph.invoke({"messages": [{"role": "user", "content": user_req}]}, config)

        while "__interrupt__" in result:
            req = result["__interrupt__"][0].value
            print(f"[待审批] 即将执行工具 {req['tool']}，参数: {req['args']}")
            decision = input("是否批准该操作？(approve/reject): ").strip().lower()
            result = graph.invoke(Command(resume=decision), config)

        print("助手:", result["messages"][-1].content)


if __name__ == "__main__":
    run_interactive()
```

* 从终端获取用户审核结果并按需执行图

## C8 多智能体

```python
from langchain.tools import tool
from langchain.agents import create_agent
from langchain_ollama import ChatOllama
from langgraph_supervisor import create_supervisor

llm = ChatOllama(model="qwen3.5:9b")


@tool
def add(a: int, b: int) -> int:
    """对两个整数做加法运算并返回结果。"""
    return a + b


@tool
def multiply(a: int, b: int) -> int:
    """对两个整数做乘法运算并返回结果。"""
    return a * b


@tool
def search_doc(query: str) -> str:
    """根据关键词检索内部文档资料（此处为模拟检索）。"""
    return f"关于 '{query}' 的文档摘要：这是一份示例检索结果。"


math_agent = create_agent(
    llm,
    tools=[add, multiply],
    system_prompt="你是一个数学计算助手，只能使用提供的工具完成加减乘除等计算，不要编造结果。",
    name="math_agent",
)

research_agent = create_agent(
    llm,
    tools=[search_doc],
    system_prompt="你是一个资料检索助手，只能使用提供的检索工具回答资料相关问题。",
    name="research_agent",
)


supervisor = create_supervisor(
    agents=[math_agent, research_agent],
    model=llm,
    prompt=(
        "你是一个主管，负责把用户的问题分配给合适的子智能体处理：\n"
        "- math_agent：负责所有数学计算类问题（如加减乘除）。\n"
        "- research_agent：负责资料检索 / 查询类问题。\n"
        "请使用对应的 transfer_to_* 工具把任务交给子智能体，"
        "等其返回结果后再汇总成最终回答，不要自己直接计算或检索。"
    ),
).compile()


if __name__ == "__main__":
    user_req = "帮我算一下 (3 + 5) * 2，并查一下 LangGraph 是什么。"
    result = supervisor.invoke({"messages": [{"role": "user", "content": user_req}]})

    for msg in result["messages"]:
        msg.pretty_print()
```

* create_supervisor 创建的是 `StateGraph` 实例，langgraph_supervisor 包帮助我们定义了图的节点、边、状态等。

## 小结

示例简单，一跑就会，真实业务，就难得很。