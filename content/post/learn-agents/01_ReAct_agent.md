+++
date = '2026-08-19T22:30:00+08:00'
draft = false
title = 'ReAct 架构详解：让大模型学会“边想边做”'
categories = ['人工智能', '智能体']
tags = ['AI Agent', 'ReAct', 'LangChain', 'Agent 架构']
toc = true
+++

> 系列文章第 1 篇 · AI Agent 六大主流架构拆解
> 本篇拆解最经典的 ReAct(Reasoning + Acting)，并把它落到一个真实场景——**油菜育种种质助手**。

## 什么是 ReAct

ReAct 由姚顺雨等人在 2022 年论文《ReAct: Synergizing Reasoning and Acting in Language Models》中提出的 AI Agent 框架，其核心思想是:**让大语言模型把"推理(Reasoning)"和"行动(Acting)"交错成循环** —— 每一步先想清楚下一步该做什么(Thought)，再去真正调用一个外部工具(Action)，拿到工具返回的结果(Observation)，再回想，如此往复，直到信息足够给出最终答案。

一个最小循环单元长这样:

```
Thought 1: 我需要先查 X
Action 1: query(X)
Observation 1: X 的结果是 ...
Thought 2: 根据 X 的结果，我还需要查 Y
Action 2: query(Y)
Observation 2: ...
...
Thought N: 信息已充分，可以回答了
Final Answer: ...
```

### 与相邻范式的关键区别

- **相比纯 Chain-of-Thought(CoT)**：CoT 只在模型脑子里推，无法获取新信息，容易凭空编造；ReAct 把*推理与外部工具结合，结果可验证*。
- **相比纯工具调用(Function Calling)**：纯工具调用是“一次调用拿结果”，模型不展示中间推理；ReAct 把推理过程显式化，模型可以多轮规划、自我调整方向。

一句话总结 ReAct 的价值：**LLM 负责想，工具负责查和算**，两者交错，既避免了纯推理的幻觉，也避免了纯工具调用的僵化。

## 示例背景：油菜育种

油菜(Brassica napus)是我国重要的油料作物，育种是一个典型的**多步推理 + 外部数据依赖**场景:

1. **数据分散在多个系统**：种质资源库、性状表型数据库、分子标记数据库、田间试验记录、育种文献，任何一步决策都要查多个来源。
2. **决策分步且依赖前序结果**：推荐一个杂交组合，要先*确定育种目标* → *查候选亲本性状* → *评估亲本配合力 / 遗传距离* → *结合文献佐证* → *给出组合建议*，步骤有序，存在依赖。
3. **需要可解释性**：育种决策要能向专家讲清楚“为什么这么配”，ReAct 把思考链条留痕，天然可审计。

这种<font color="blue" style="font-weight: bolder;">查-想-再查-再想-...-结果</font>的工作流，正是 ReAct 循环的用武之地。下面用 LangChain 实现一个油菜育种 ReAct Agent。

## 架构与运行机制

ReAct 的运行机制由四要素构成:

![ReAct 运行机制四要素](/imgs/learn-agents/react-four-elements.png)

- **Thought**：模型对当前已知信息与缺失信息的自然语言推理，决定下一步动作。
- **Action**：工具名 + 入参，真正去外部拿数据。
- **Observation**：工具执行后返回的文本结果，拼回上下文。
- **循环终止**：模型在 Thought 中判定<font color="blue" style="font-weight: bolder;">信息已充分</font>，输出 `Final Answer`，跳出循环。

关键设计点:

- 工具描述(prompt 里的工具清单)必须足够清晰，模型才知道「有什么能调用」;
- 每次循环都要把历史 Thought / Action / Observation 累积进上下文，因此**步数越多 token 消耗越大** —— 这是 ReAct 的主要成本来源。

### 循环终止：信息判定

经验上，模型在做这个判断时，隐性对照的是：

1. **覆盖度**：原问题拆出的子问题/所需事实，是否都已有 Observation 支撑?还缺不缺？
2. **可推导性**：已有信息是否足以逻辑推出结论(而不是还卡在某个未知)？
3. **成功标准**：问题里若有硬约束(如"亩产 > 200kg、含油量 > 45%、抗病")，这些条件是否都有数据可以判断？
4. **一致性**：多次 Observation 之间有没有矛盾需要再查？

如何使用判定更可靠？

* 在 prompt 里写明成功标准：把“什么算回答完成”显式化
* 让 Thought 显式自检：提示模型在判定前先复述“已覆盖哪些子问题、还缺什么”，再决定是否收尾。
* 控制 Observation 信息密度：返回精炼、结构化(关键字段优先)，判定更准。
* 设合理 max_iterations：最大循环次数兜底。

## 基于 LangChain 的代码演示

* 模型使用本地搭建的 Ollama 中的 `qwen3.5:9b`

### 可观测性：日志配置

打印 Agent 的调用过程。

```python
import logging

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s | %(levelname)-7s | %(name)s | %(message)s",
    datefmt="%H:%M:%S",
)
logger = logging.getLogger("ReAct Agent")
```

### 定义油菜育种工具

用 `@tool` 装饰器定义两个工具：

* query_germplasm：查询种质信息
* search_literature：查询种质文献信息

```python
from langchain.tools import tool
from pydantic import BaseModel, Field


# 模拟种质资源数据源
GERMPLASM_DB = {
    "中双11号": {"产量kg/亩": 195, "含油量%": 43.2, "菌核病抗性": "中抗", "生育期d": 230},
    "华油杂62R": {"产量kg/亩": 208, "含油量%": 45.6, "菌核病抗性": "抗", "生育期d": 225},
    "秦油33": {"产量kg/亩": 202, "含油量%": 44.8, "菌核病抗性": "中抗", "生育期d": 235},
    "宁杂1818": {"产量kg/亩": 188, "含油量%": 46.1, "菌核病抗性": "抗", "生育期d": 228},
    "陕油28": {"产量kg/亩": 210, "含油量%": 45.0, "菌核病抗性": "高抗", "生育期d": 232},
}

# 模拟文献数据源
LITERATURE = [
    {"title": "甘蓝型油菜菌核病抗性 QTL 定位研究", "abstract": "陕油28、宁杂1818 对菌核病表现稳定抗性，含抗病位点。"},
    {"title": "高油酸油菜亲本配合力分析", "abstract": "华油杂62R × 陕油28 组合一般配合力高，后代含油量与产量兼优。"},
]


# 种质查询工具
class GermplasmQueryArgs(BaseModel):
    germplasm_name: str = Field(description="The name of the germplasm to query")


@tool(
    name_or_callable="query_germplasm",
    description="Query information about a specific germplasm.",
    args_schema=GermplasmQueryArgs,
)
def query_germplasm(germplasm_name: str) -> str:
    logger.info(f"TOOL query_germplasm | args={germplasm_name}")
    if germplasm_name in GERMPLASM_DB:
        result = f"Information about germplasm '{germplasm_name}': {GERMPLASM_DB[germplasm_name]}"
    else:
        result = f"Germplasm '{germplasm_name}' not found."
    logger.info(f"TOOL query_germplasm | result={result}")
    return result


# 文献查询工具
class LiteratureSearchArgs(BaseModel):
    keywords: list[str] = Field(description="The keywords to search for in literature titles and abstracts")


@tool(
    name_or_callable="search_literature",
    description="Search literature for papers related to a keyword.",
    args_schema=LiteratureSearchArgs,
)
def search_literature(keywords: list[str]) -> str:
    logger.info(f"TOOL search_literature | args={keywords}")
    hits = [d for d in LITERATURE if any(keyword in d["title"] or keyword in d["abstract"] for keyword in keywords)]
    if not hits:
        result = f"No literature found for keywords '{keywords}'"
    else:
        result = "\n".join(f"- {d['title']}:{d['abstract']}" for d in hits)
    logger.info(f"TOOL search_literature | hits={len(hits)} result={result}")
    return result


tools = [query_germplasm, search_literature]
```

### 构建 ReAct Agent

```python
from langchain.agents import create_agent
from langchain_ollama import ChatOllama

llm = ChatOllama(model="qwen3.5:9b", temperature=0.3)
agents = create_agent(llm, tools=tools, system_prompt="You are a Rapeseed Germplasm and Literature Query Assistant.")
```

### 查询函数

```python
import time
from langchain.messages import HumanMessage


def user_query(query: str):
    start = time.perf_counter()
    logger.info(f"RUN start | query={query}")
    messages = [HumanMessage(content=query)]

    final_chunk = None
    for step_idx, chunk in enumerate(agents.stream({"messages": messages}, stream_mode="updates"), start=1):
        logger.info(f"Chunk={chunk}")
        for node_name, node_update in chunk.items():
                logger.info(f"STEP {step_idx:02d} | node={node_name}")
                messages = node_update.get("messages", []) if isinstance(node_update, dict) else []
                for msg in messages:
                    if getattr(msg, "tool_calls", None):
                        for tc in msg.tool_calls:
                            logger.info(f"  └─ LLM → TOOL CALL | name={tc.get("name")} args={tc.get("args")}")
                    elif getattr(msg, "tool_call_id", None):
                        logger.info(f"  └─ TOOL → RESULT | id={msg.tool_call_id} content={str(msg.content)[:200]}")
                    else:
                        content = getattr(msg, "content", "")
                        if content:
                            logger.info(f"  └─ LLM → CONTENT | {str(content)[:200]}", )
        final_chunk = chunk

    elapsed = time.perf_counter() - start
    answer = final_chunk[list(final_chunk)[-1]]["messages"][-1].content
    logger.info(f"RUN end | elapsed={elapsed:.2f}s")
    print("\n========== FINAL ANSWER ==========\n")
    print(answer)
```

### query1：查询抗菌核病的两个亲本的种质信息

```python
user_query("查询抗菌核病的两个亲本的种质信息")
```

输出：

```bash
20:04:21 | INFO    | ReAct Agent | RUN start | query=查询抗菌核病的两个亲本的种质信息
20:04:21 | INFO    | httpx | HTTP Request: POST http://127.0.0.1:11434/api/chat "HTTP/1.1 200 OK"
20:04:30 | INFO    | ReAct Agent | Chunk={'model': {'messages': [AIMessage(content='', additional_kwargs={}, response_metadata={'model': 'qwen3.5:9b', 'created_at': '2026-08-19T09:04:30.148677Z', 'done': True, 'done_reason': 'stop', 'total_duration': 9068554958, 'load_duration': 318049041, 'prompt_eval_count': 402, 'prompt_eval_duration': 161936000, 'eval_count': 205, 'eval_duration': 8426183000, 'logprobs': None, 'model_name': 'qwen3.5:9b', 'model_provider': 'ollama'}, id='lc_run--01a01943-ae3e-76e1-825f-88dd8a9adb1e-0', tool_calls=[{'name': 'search_literature', 'args': {'keywords': ['抗菌核病', '亲本']}, 'id': 'cf14c80b-a294-435a-8ee2-7de6d927d9a0', 'type': 'tool_call'}], invalid_tool_calls=[], usage_metadata={'input_tokens': 402, 'output_tokens': 205, 'total_tokens': 607})]}}
20:04:30 | INFO    | ReAct Agent | STEP 01 | node=model
20:04:30 | INFO    | ReAct Agent |   └─ LLM → TOOL CALL | name=search_literature args={'keywords': ['抗菌核病', '亲本']}
20:04:30 | INFO    | ReAct Agent | TOOL search_literature | args=['抗菌核病', '亲本']
20:04:30 | INFO    | ReAct Agent | TOOL search_literature | hits=1 result=- 高油酸油菜亲本配合力分析:华油杂62R × 陕油28 组合一般配合力高，后代含油量与产量兼优。
20:04:30 | INFO    | ReAct Agent | Chunk={'tools': {'messages': [ToolMessage(content='- 高油酸油菜亲本配合力分析:华油杂62R × 陕油28 组合一般配合力高，后代含油量与产量兼优。', name='search_literature', id='5ad4a94e-ba4a-40a0-9685-c82b9b51ae27', tool_call_id='cf14c80b-a294-435a-8ee2-7de6d927d9a0')]}}
20:04:30 | INFO    | ReAct Agent | STEP 02 | node=tools
20:04:30 | INFO    | ReAct Agent |   └─ TOOL → RESULT | id=cf14c80b-a294-435a-8ee2-7de6d927d9a0 content=- 高油酸油菜亲本配合力分析:华油杂62R × 陕油28 组合一般配合力高，后代含油量与产量兼优。
20:04:31 | INFO    | httpx | HTTP Request: POST http://127.0.0.1:11434/api/chat "HTTP/1.1 200 OK"
20:04:39 | INFO    | ReAct Agent | Chunk={'model': {'messages': [AIMessage(content='', additional_kwargs={}, response_metadata={'model': 'qwen3.5:9b', 'created_at': '2026-08-19T09:04:39.319848Z', 'done': True, 'done_reason': 'stop', 'total_duration': 9130557834, 'load_duration': 294287459, 'prompt_eval_count': 493, 'prompt_eval_duration': 542766000, 'eval_count': 204, 'eval_duration': 8272120999, 'logprobs': None, 'model_name': 'qwen3.5:9b', 'model_provider': 'ollama'}, id='lc_run--01a01943-d1e8-7a22-8bb0-441fdbdfe3fe-0', tool_calls=[{'name': 'query_germplasm', 'args': {'germplasm_name': '华油杂62R'}, 'id': 'fcfb7ff2-4331-47c1-afe1-b226320061dc', 'type': 'tool_call'}, {'name': 'query_germplasm', 'args': {'germplasm_name': '陕油28'}, 'id': 'a7ed1f80-e0d6-48e3-b49b-0b347aa3dec7', 'type': 'tool_call'}], invalid_tool_calls=[], usage_metadata={'input_tokens': 493, 'output_tokens': 204, 'total_tokens': 697})]}}
20:04:39 | INFO    | ReAct Agent | STEP 03 | node=model
20:04:39 | INFO    | ReAct Agent |   └─ LLM → TOOL CALL | name=query_germplasm args={'germplasm_name': '华油杂62R'}
20:04:39 | INFO    | ReAct Agent |   └─ LLM → TOOL CALL | name=query_germplasm args={'germplasm_name': '陕油28'}
20:04:39 | INFO    | ReAct Agent | TOOL query_germplasm | args=陕油28
20:04:39 | INFO    | ReAct Agent | TOOL query_germplasm | result=Information about germplasm '陕油28': {'产量kg/亩': 210, '含油量%': 45.0, '菌核病抗性': '高抗', '生育期d': 232}
20:04:39 | INFO    | ReAct Agent | TOOL query_germplasm | args=华油杂62R
20:04:39 | INFO    | ReAct Agent | TOOL query_germplasm | result=Information about germplasm '华油杂62R': {'产量kg/亩': 208, '含油量%': 45.6, '菌核病抗性': '抗', '生育期d': 225}
20:04:39 | INFO    | ReAct Agent | Chunk={'tools': {'messages': [ToolMessage(content="Information about germplasm '陕油28': {'产量kg/亩': 210, '含油量%': 45.0, '菌核病抗性': '高抗', '生育期d': 232}", name='query_germplasm', id='d6e1d556-6ca1-490f-9472-2ad6f9d91383', tool_call_id='a7ed1f80-e0d6-48e3-b49b-0b347aa3dec7')]}}
20:04:39 | INFO    | ReAct Agent | STEP 04 | node=tools
20:04:39 | INFO    | ReAct Agent |   └─ TOOL → RESULT | id=a7ed1f80-e0d6-48e3-b49b-0b347aa3dec7 content=Information about germplasm '陕油28': {'产量kg/亩': 210, '含油量%': 45.0, '菌核病抗性': '高抗', '生育期d': 232}
20:04:39 | INFO    | ReAct Agent | Chunk={'tools': {'messages': [ToolMessage(content="Information about germplasm '华油杂62R': {'产量kg/亩': 208, '含油量%': 45.6, '菌核病抗性': '抗', '生育期d': 225}", name='query_germplasm', id='d9bc3d4d-42f4-4deb-8ff5-6ff9a8c418aa', tool_call_id='fcfb7ff2-4331-47c1-afe1-b226320061dc')]}}
20:04:39 | INFO    | ReAct Agent | STEP 05 | node=tools
20:04:39 | INFO    | ReAct Agent |   └─ TOOL → RESULT | id=fcfb7ff2-4331-47c1-afe1-b226320061dc content=Information about germplasm '华油杂62R': {'产量kg/亩': 208, '含油量%': 45.6, '菌核病抗性': '抗', '生育期d': 225}
20:04:40 | INFO    | httpx | HTTP Request: POST http://127.0.0.1:11434/api/chat "HTTP/1.1 200 OK"
20:04:55 | INFO    | ReAct Agent | Chunk={'model': {'messages': [AIMessage(content='根据查询结果，抗菌核病的两个亲本的种质信息如下：\n\n### 1. **华油杂62R**  \n- **产量**: 208 kg/亩  \n- **含油量**: 45.6%  \n- **菌核病抗性**: 抗  \n- **生育期**: 225天  \n\n### 2. **陕油28**  \n- **产量**: 210 kg/亩  \n- **含油量**: 45.0%  \n- **菌核病抗性**: 高抗  \n- **生育期**: 232天  \n\n这两个亲本在抗病育种中表现出良好的配合力，后代兼具高产、优质和抗病的特性。', additional_kwargs={}, response_metadata={'model': 'qwen3.5:9b', 'created_at': '2026-08-19T09:04:55.889843Z', 'done': True, 'done_reason': 'stop', 'total_duration': 16523086292, 'load_duration': 259025750, 'prompt_eval_count': 694, 'prompt_eval_duration': 966987000, 'eval_count': 377, 'eval_duration': 15277772000, 'logprobs': None, 'model_name': 'qwen3.5:9b', 'model_provider': 'ollama'}, id='lc_run--01a01943-f5c0-7bd0-81bd-54679b35a2b0-0', tool_calls=[], invalid_tool_calls=[], usage_metadata={'input_tokens': 694, 'output_tokens': 377, 'total_tokens': 1071})]}}
20:04:55 | INFO    | ReAct Agent | STEP 06 | node=model
20:04:55 | INFO    | ReAct Agent |   └─ LLM → CONTENT | 根据查询结果，抗菌核病的两个亲本的种质信息如下：

### 1. **华油杂62R**  
- **产量**: 208 kg/亩  
- **含油量**: 45.6%  
- **菌核病抗性**: 抗  
- **生育期**: 225天  

### 2. **陕油28**  
- **产量**: 210 kg/亩  
- **含油量**: 45.0%  
- **菌核病抗性**: 高抗  
- 
20:04:55 | INFO    | ReAct Agent | RUN end | elapsed=34.85s

========== FINAL ANSWER ==========

根据查询结果，抗菌核病的两个亲本的种质信息如下：

### 1. **华油杂62R**  
- **产量**: 208 kg/亩  
- **含油量**: 45.6%  
- **菌核病抗性**: 抗  
- **生育期**: 225天  

### 2. **陕油28**  
- **产量**: 210 kg/亩  
- **含油量**: 45.0%  
- **菌核病抗性**: 高抗  
- **生育期**: 232天  

这两个亲本在抗病育种中表现出良好的配合力，后代兼具高产、优质和抗病的特性。
```

通过日志，可得到大致的过程：

1. 用户的 query 先给到模型
2. 模型返回需要 `TOOL CALL`，调用 `search_literature`，参数为 `['抗菌核病', '亲本']`
3. 工具调用的结果转成 `ToolMessage` 再给回模型
4. 模型返回需要两次 `TOOL CALL`，调用 `query_germplasm`，两次参数分别为 `陕油28`、`华油杂62R`
5. 两次调用的结果都转成 `ToolMessage` 再给回模型
6. 模型总结并生成结果

### query2：中双11号的种质信息

```python
user_query("中双11号的种质信息")
```

输出：

```bash
21:18:59 | INFO    | ReAct Agent | RUN start | query=中双11号的种质信息
21:19:10 | INFO    | httpx | HTTP Request: POST http://127.0.0.1:11434/api/chat "HTTP/1.1 200 OK"
21:19:15 | INFO    | ReAct Agent | Chunk={'model': {'messages': [AIMessage(content='', additional_kwargs={}, response_metadata={'model': 'qwen3.5:9b', 'created_at': '2026-08-19T10:19:15.331972Z', 'done': True, 'done_reason': 'stop', 'total_duration': 15747650750, 'load_duration': 9306504750, 'prompt_eval_count': 400, 'prompt_eval_duration': 1714101000, 'eval_count': 110, 'eval_duration': 4699443000, 'logprobs': None, 'model_name': 'qwen3.5:9b', 'model_provider': 'ollama'}, id='lc_run--01a01988-0479-7351-91c1-ddf90f79b316-0', tool_calls=[{'name': 'query_germplasm', 'args': {'germplasm_name': '中双11号'}, 'id': '447d1e8e-1b00-44c6-8546-a55b752ff4d8', 'type': 'tool_call'}], invalid_tool_calls=[], usage_metadata={'input_tokens': 400, 'output_tokens': 110, 'total_tokens': 510})]}}
21:19:15 | INFO    | ReAct Agent | STEP 01 | node=model
21:19:15 | INFO    | ReAct Agent |   └─ LLM → TOOL CALL | name=query_germplasm args={'germplasm_name': '中双11号'}
21:19:15 | INFO    | ReAct Agent | TOOL query_germplasm | args=中双11号
21:19:15 | INFO    | ReAct Agent | TOOL query_germplasm | result=Information about germplasm '中双11号': {'产量kg/亩': 195, '含油量%': 43.2, '菌核病抗性': '中抗', '生育期d': 230}
21:19:15 | INFO    | ReAct Agent | Chunk={'tools': {'messages': [ToolMessage(content="Information about germplasm '中双11号': {'产量kg/亩': 195, '含油量%': 43.2, '菌核病抗性': '中抗', '生育期d': 230}", name='query_germplasm', id='eabf0fa9-a72e-48b9-8946-ace536f51f08', tool_call_id='447d1e8e-1b00-44c6-8546-a55b752ff4d8')]}}
21:19:15 | INFO    | ReAct Agent | STEP 02 | node=tools
21:19:15 | INFO    | ReAct Agent |   └─ TOOL → RESULT | id=447d1e8e-1b00-44c6-8546-a55b752ff4d8 content=Information about germplasm '中双11号': {'产量kg/亩': 195, '含油量%': 43.2, '菌核病抗性': '中抗', '生育期d': 230}
21:19:16 | INFO    | httpx | HTTP Request: POST http://127.0.0.1:11434/api/chat "HTTP/1.1 200 OK"
21:19:24 | INFO    | ReAct Agent | Chunk={'model': {'messages': [AIMessage(content='以下是“中双11号”种质的详细信息：\n\n| 性状 | 数值/描述 |\n| :--- | :--- |\n| **产量** | 195 kg/亩 |\n| **含油量** | 43.2% |\n| **菌核病抗性** | 中抗 |\n| **生育期** | 230天 |\n\n该品种具有较好的高产潜力和中等以上的抗病性，适合在相应区域种植。', additional_kwargs={}, response_metadata={'model': 'qwen3.5:9b', 'created_at': '2026-08-19T10:19:24.675368Z', 'done': True, 'done_reason': 'stop', 'total_duration': 9240134292, 'load_duration': 282904375, 'prompt_eval_count': 507, 'prompt_eval_duration': 709154000, 'eval_count': 202, 'eval_duration': 8223772000, 'logprobs': None, 'model_name': 'qwen3.5:9b', 'model_provider': 'ollama'}, id='lc_run--01a01988-4268-72c0-bbd2-febae65d56e5-0', tool_calls=[], invalid_tool_calls=[], usage_metadata={'input_tokens': 507, 'output_tokens': 202, 'total_tokens': 709})]}}
21:19:24 | INFO    | ReAct Agent | STEP 03 | node=model
21:19:24 | INFO    | ReAct Agent |   └─ LLM → CONTENT | 以下是“中双11号”种质的详细信息：

| 性状 | 数值/描述 |
| :--- | :--- |
| **产量** | 195 kg/亩 |
| **含油量** | 43.2% |
| **菌核病抗性** | 中抗 |
| **生育期** | 230天 |

该品种具有较好的高产潜力和中等以上的抗病性，适合在相应区域种植。
21:19:24 | INFO    | ReAct Agent | RUN end | elapsed=25.13s

========== FINAL ANSWER ==========

以下是“中双11号”种质的详细信息：

| 性状 | 数值/描述 |
| :--- | :--- |
| **产量** | 195 kg/亩 |
| **含油量** | 43.2% |
| **菌核病抗性** | 中抗 |
| **生育期** | 230天 |

该品种具有较好的高产潜力和中等以上的抗病性，适合在相应区域种植。
```

通过日志，可得到大致的过程：

1. 模型识别到了“中双11号”
2. 模型返回需要 `TOOL CALL`，调用 `query_germplasm`，参数为 `中双11号`
3. 结果转换成 `ToolMessage`
4. 模型总结并生成结果


## 优缺点

### 优点

- **避免幻觉**：结论来自工具返回的真实数据，而非模型脑补。
- **轻量易落地**：单个 LLM + 一组工具即可起步，不需要复杂编排框架。
- **工具可插拔**：新增一个种质库接口，只需加一个 `@tool`，模型自动学会调用。

### 缺点

- **步数线性增长 token 与延迟**：每步都要把历史轨迹拼回上下文，长任务成本高。
- **缺乏全局纠偏**：某一步 Thought 走偏，没有机制纠正，会越错越远。
- **工具失败会传导**：某个工具报错或返回脏数据，后续推理全部受影响，需要做好异常处理与数据清洗。
- **上下文窗口压力**：复杂任务步数多时，历史轨迹可能撑爆窗口，需要做摘要 / 裁剪。
- **依赖工具描述质量**：工具描述写得不清楚，模型就不知道何时该用，误调用率高。

## 小结

ReAct 是几乎所有 Agent 架构的起点 —— 它解决了“一个 LLM 怎么一步步把任务做完”的问题。但当任务变长、变复杂，ReAct 的短板(缺乏全局规划、无纠偏机制、token 爆炸)就暴露出来。