+++
date = '2026-10-23T14:22:46+08:00'
title = 'LangChain/LangGraph 流式输出：stream_mode 不同值的效果 2'
categories = ['人工智能', '智能体']
tags = ['AI Agent', 'LangGraph', 'LangChain', 'Stream Mode']
toc = true
+++

## stream_mode=checkpoints

* 触发时机：每一次 super-step（超步）生成 checkpoint 之后，就产出一个 chunk
* 适用场景：
    * 拿到 `checkpoint_id` 就可以从这个快照恢复执行
    * 长时间运行的任务、故障恢复

## stream_mode=tasks

按单个 Task 的生命周期推送事件：Task Start、Task Finish，一般返回的 chunk 两个为一个 Task。

## 演示 stream_mode=checkpoints 的效果

要使用 `stream_mode="checkpoints"`，必须要传入 `checkpointer` 参数，相应地，在 `stream()` 方法中，必须要传入 `RunnableConfig` 参数。

1. 构建智能体

```python
from langchain.agents import create_agent
from langchain.tools import tool
from langchain_ollama import ChatOllama
from langgraph.config import RunnableConfig
from langgraph.checkpoint.memory import InMemorySaver

llm = ChatOllama(model="qwen3.5:9b")


@tool
def current_time():
    """返回当前时间"""
    now = datetime.now()
    return now.strftime("%Y-%m-%d %H:%M:%S")

checkpointer = InMemorySaver()
config = RunnableConfig(configurable={"thread_id": "user-001"})
agent = create_agent(llm, tools=[current_time], checkpointer=checkpointer)
```

2. 定义帮助方法 `get_types`

```python
def get_types(messages: list) -> list:
    """返回消息类型"""
    return [message.__class__.__name__ for message in messages]
```

3. 流式执行

```python
for i, chunk in enumerate(agent.stream({"messages": [{"role": "user", "content": "现在是什么时候"}]}, config=config, stream_mode="checkpoints")):
    print(f"================= Chunk {i+1} =================")

    chunk_config = chunk["config"]
    print(f"检查点：{chunk_config['configurable']['checkpoint_id']}")

    chunk_values = chunk["values"]
    messages = chunk_values["messages"]
    message_types = get_types(messages)

    print(f"消息数：{len(chunk_values['messages'])}，消息类型：{message_types}")
    print(f"下一步执行节点 {chunk['next']}")

    if len(chunk_values['messages']) > 0:
        last_message = messages[-1]
        last_message_type = message_types[-1]
        print(f"最新一条消息类型为 {last_message_type}")

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

* 每一步迭代返回一个 chunk，主要属性有 `config`，`values`、`next` 等
* chunk 中包含当前检测点信息及下一步要执行的节点
* `values` 中包含当前状态的所有的累积消息

```bash
$ python stream_mode_checkpoints.py
================= Chunk 1 =================
检查点：1f1ceacd-e1ec-61aa-bfff-8f84c2228379
消息数：0，消息类型：[]
下一步执行节点 ['__start__']

================= Chunk 2 =================
检查点：1f1ceacd-e1ee-64e6-8000-e7e51134d964
消息数：1，消息类型：['HumanMessage']
下一步执行节点 ['model']
最新一条消息类型为 HumanMessage
用户输入：现在是什么时候

================= Chunk 3 =================
检查点：1f1ceacd-f43a-6442-8001-ff52b609ce28
消息数：2，消息类型：['HumanMessage', 'AIMessage']
下一步执行节点 ['tools']
最新一条消息类型为 AIMessage
输入 token：253，输出 token：38
模型决定调用工具 current_time

================= Chunk 4 =================
检查点：1f1ceacd-f43f-62b2-8002-7b05ff39ecad
消息数：3，消息类型：['HumanMessage', 'AIMessage', 'ToolMessage']
下一步执行节点 ['model']
最新一条消息类型为 ToolMessage
调用工具 current_time，返回 2026-10-23 14:42:10

================= Chunk 5 =================
检查点：1f1ceace-1090-625e-8003-83ba61768fda
消息数：4，消息类型：['HumanMessage', 'AIMessage', 'ToolMessage', 'AIMessage']
下一步执行节点 []
最新一条消息类型为 AIMessage
输入 token：303，输出 token：71
模型回复：现在是 **2026年10月23日 14:42:10**。
```

## 演示 stream_mode=tasks 的效果

```python
for i, chunk in enumerate(agent.stream({"messages": [{"role": "user", "content": "现在是什么时候"}]}, config=config, stream_mode="tasks")):
    print(f"================= Chunk {i+1} =================")

    print(f"任务 ID：{chunk['id']}，名称：{chunk['name']}")
    if "input" in chunk:
        print(f"输入：{chunk['input']}")
    elif "result" in chunk:
        print(f"结果：{chunk['result']}")

    print()
```

```python
python stream_mode_tasks.py
================= Chunk 1 =================
任务 ID：904b8460-6e41-b46c-a7cb-908e4a60a128，名称：model
输入：{'messages': [HumanMessage(content='现在是什么时候', additional_kwargs={}, response_metadata={}, id='492773d7-bdd7-4a76-9196-b6a565c3bf08')]}

================= Chunk 2 =================
任务 ID：904b8460-6e41-b46c-a7cb-908e4a60a128，名称：model
结果：{'messages': [AIMessage(content='', additional_kwargs={}, response_metadata={'model': 'qwen3.5:9b', 'created_at': '2026-10-23T07:12:21.087478Z', 'done': True, 'done_reason': 'stop', 'total_duration': 3643127000, 'load_duration': 55083834, 'prompt_eval_count': 253, 'prompt_eval_duration': 2150367000, 'eval_count': 38, 'eval_duration': 1356073000, 'logprobs': None, 'model_name': 'qwen3.5:9b', 'model_provider': 'ollama'}, id='lc_run--01a1679a-7207-7292-bdc4-0b09c7502962-0', tool_calls=[{'name': 'current_time', 'args': {}, 'id': 'c56f8ced-b443-4668-94a0-b311b33dd55f', 'type': 'tool_call'}], invalid_tool_calls=[], usage_metadata={'input_tokens': 253, 'output_tokens': 38, 'total_tokens': 291})]}

================= Chunk 3 =================
任务 ID：eb46f426-8c63-962b-2d0c-2c1e61b4b701，名称：tools
输入：{'__type': 'tool_call_with_context', 'tool_call': {'name': 'current_time', 'args': {}, 'id': 'c56f8ced-b443-4668-94a0-b311b33dd55f', 'type': 'tool_call'}, 'state': {'messages': [HumanMessage(content='现在是什么时候', additional_kwargs={}, response_metadata={}, id='492773d7-bdd7-4a76-9196-b6a565c3bf08'), AIMessage(content='', additional_kwargs={}, response_metadata={'model': 'qwen3.5:9b', 'created_at': '2026-10-23T07:12:21.087478Z', 'done': True,'done_reason': 'stop', 'total_duration': 3643127000, 'load_duration': 55083834, 'prompt_eval_count': 253, 'prompt_eval_duration': 2150367000, 'eval_count': 38, 'eval_duration': 1356073000, 'logprobs': None, 'model_name': 'qwen3.5:9b', 'model_provider': 'ollama'}, id='lc_run--01a1679a-7207-7292-bdc4-0b09c7502962-0', tool_calls=[{'name': 'current_time', 'args': {}, 'id': 'c56f8ced-b443-4668-94a0-b311b33dd55f', 'type': 'tool_call'}], invalid_tool_calls=[], usage_metadata={'input_tokens': 253, 'output_tokens': 38, 'total_tokens': 291})]}}

================= Chunk 4 =================
任务 ID：eb46f426-8c63-962b-2d0c-2c1e61b4b701，名称：tools
结果：{'messages': [ToolMessage(content='2026-10-23 15:12:21', name='current_time', id='6b80f732-c317-43f3-87bc-a6fe709e3363', tool_call_id='c56f8ced-b443-4668-94a0-b311b33dd55f')]}

================= Chunk 5 =================
任务 ID：f6894b2f-c9d8-ec71-d8c3-074caf1ca207，名称：model
输入：{'messages': [HumanMessage(content='现在是什么时候', additional_kwargs={}, response_metadata={}, id='492773d7-bdd7-4a76-9196-b6a565c3bf08'), AIMessage(content='', additional_kwargs={}, response_metadata={'model': 'qwen3.5:9b', 'created_at': '2026-10-23T07:12:21.087478Z', 'done': True, 'done_reason': 'stop', 'total_duration': 3643127000, 'load_duration': 55083834, 'prompt_eval_count': 253, 'prompt_eval_duration': 2150367000, 'eval_count': 38, 'eval_duration': 1356073000, 'logprobs': None, 'model_name': 'qwen3.5:9b', 'model_provider': 'ollama'}, id='lc_run--01a1679a-7207-7292-bdc4-0b09c7502962-0', tool_calls=[{'name': 'current_time', 'args': {}, 'id': 'c56f8ced-b443-4668-94a0-b311b33dd55f', 'type': 'tool_call'}], invalid_tool_calls=[], usage_metadata={'input_tokens': 253, 'output_tokens': 38, 'total_tokens': 291}), ToolMessage(content='2026-10-23 15:12:21', name='current_time', id='6b80f732-c317-43f3-87bc-a6fe709e3363', tool_call_id='c56f8ced-b443-4668-94a0-b311b33dd55f')]}

================= Chunk 6 =================
任务 ID：f6894b2f-c9d8-ec71-d8c3-074caf1ca207，名称：model
结果：{'messages': [AIMessage(content='现在是 2026年10月23日 星期五 下午 15时12分21秒。', additional_kwargs={}, response_metadata={'model': 'qwen3.5:9b', 'created_at': '2026-10-23T07:12:24.126197Z', 'done': True, 'done_reason':'stop', 'total_duration': 3031338584, 'load_duration': 2495000, 'prompt_eval_count': 303, 'prompt_eval_duration': 413064000, 'eval_count': 72, 'eval_duration': 2610336000, 'logprobs': None, 'model_name': 'qwen3.5:9b', 'model_provider': 'ollama'}, id='lc_run--01a1679a-8065-7251-bc97-c2dca7afbd1b-0', tool_calls=[], invalid_tool_calls=[], usage_metadata={'input_tokens': 303, 'output_tokens': 72, 'total_tokens': 375})]}
```

* 从上述的输出可知，chunk1 和 chunk2 是同一个任务的开始与结束，同理 chunk3 和 chunk4、chunk5 和 chunk6
* TaskStart 是输入信息，TaskFinish 是结果信息