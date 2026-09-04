+++
date = '2026-09-03T14:15:35+08:00'
draft = false
title = 'Memory-Augmented 架构详解：让 Agent 记住你是谁'
categories = ['人工智能', '智能体']
tags = ['AI Agent', 'Memory-Augmented', 'LangChain', 'Agent 架构']
toc = true
+++

> 系列文章第 6 篇 · AI Agent 六大主流架构拆解
> 本篇拆解 Memory-Augmented 架构，把"短期工作记忆 + 跨会话长期记忆 + 记忆控制器"装进 Agent，让它"记住用户、记住经验、记住你自己"，示例选择**油菜表型采集**这一典型"长期、连续、按口径"的任务。

## 什么是 Memory-Augmented

核心思想一句话：**在 Agent 外挂一层可读写的记忆层，让 Agent 不仅能调用工具，还能在多次会话之间累积"事实、偏好、经验" —— 下次对话时，Agent 知道你是谁、关心什么、按什么口径回答。**

一个最小运行单元：

```
用户请求
  → ① 记忆召回：从长期记忆库检索与本轮相关条目
  → ② 推理与行动：注入上下文，Agent 照常思考与工具调用
  → 最终应答：把结果回给用户
  → ③ 编码与写入：抽取值得长期保留的信息，结构化落盘
  → 整合与遗忘：去重 · 冲突消解 · 时间衰减(周期任务)
```

Memory-Augmented 关注的是**用户相关的事实与偏好** —— 会写、会更新、会随时间衰减、会冲突消解。

### 记忆的划分

按**时效性**切分，记忆可分成三类：

- **瞬时记忆**（Immediate Memory / Sensory Memory）：单次上下文窗口内，对话结束即丢弃，存储内容包括：当前轮完整 Prompt、工具返回结果、最新的用户输入，原始流式上下文。
- **短期记忆**（Short‑Term Memory / Working Memory）：一次会话周期（用户打开对话到关闭会话），保存会话内上下文，做压缩，规避窗口溢出，会话销毁后丢失。
- **长期记忆**（Long‑Term Memory）：跨会话持久化存储，永久留存，可检索召回，存储内容包括：用户偏好、历史事实、过往经验、反思总结、知识库片段。

按**内容性质**切分，记忆可分成四类：

- **陈述性记忆**（Declarative Memory）：存放可被语言描述的事实、经历与知识，包含情景事件和语义常识。
- **程序性记忆**（Procedural Memory）：存储从过往经验习得的行动策略、处理流程与做事范式，指导 Agent 如何执行动作。
- **反思记忆**（Reflective Memory）：对历史执行结果复盘归因，产出错误教训与改进方案的高阶自省记忆。
- **工作记忆**（Working Memory）：保存当前任务的临时中间状态，用于本轮推理、计划与 Agent 之间的信息流转。

两套划分并不是并列的另一套切片，而是**不同维度对同一件事的描述**：时效性决定**记忆能活多久**，内容性质决定**记忆装的是哪一类东西**。

## 示例背景：油菜表型采集

油菜表型采集是典型的**周期性、强标准、重经验**业务：

1. **分期采集、错过不可补**：从苗期（出苗期、五叶期、苗期生长习性）→ 蕾薹期（抽薹期、薹茎色）→ 花期（初花、盛花、终花期）→ 成熟期（株高、分枝数、角果数、每果粒数），一个生育周期要下地四五趟。花期那趟没记，就只能等下一年。
2. **标准重**：性状名称、描述符编号、计量单位、取样方法全部按《油菜种质资源描述规范和数据标准》执行——该书规定了描述符及其分级标准、字段类型与小数位、以及数据采集全过程的质量控制方法。换个技术员采集，数据必须能对得上。
3. **经验强**：同样的"全株有效角果数 342.6 个"，放在不同密度、不同年份、不同试点，解读不同；异常值是不是仪器误差、要不要现场重测，靠的是长期经验。
4. **跨会话连续**：孙老师上周在 A 会话里说「我负责杨凌试验站的油菜表型采集，报告按《油菜种质资源描述规范和数据标准》的描述符编号和单位输出，重点看全株有效角果数和菌核病」——这周在 B 会话里，他直接问「P-2025-0412 这批表型数据怎么样」，Agent 应当**接着上次的口径回答**，而不是让他再交代一遍。

## 架构与运行机制

![Memory-Augmented 运行机制](/imgs/learn-agents/memory-augmented-mechanism.png)

运行机制四阶段：

- **① 记忆召回（Retrieve）**：新请求到来时，先按 query 的关键词 / 标签 / 语义去检索长期记忆库，把相关条目作为补充上下文注入 system prompt。**这一步骤决定了"Agent 有多了解你"**。
- **② 推理与行动（Reason & Act）**：在记忆增强的上下文中跑常规 Agent 循环：思考、调工具、观察、再思考。Agent 拿到工具返回后，会自动把结果"按记忆里那个用户的口径"重新组织，而不是泛泛地回答。
- **③ 编码与写入（Encode/Write）**：每轮对话结束后，独立运行一个"记忆写入器" Agent，对当轮对话做评估：哪些信息值得长期保留？归到语义 / 情景 / 程序哪一类？打个重要度分（1–5），加几个检索标签，结构化落到 JSON / 向量库。
- **整合与遗忘（Consolidate）**：周期任务，对记忆库做去重、冲突消解、按时间衰减淘汰过期信息。

关键设计点：

| 阶段 | 定位 | 动作 | 设计价值 |
| :--- | :--- | :--- | :--- |
| 记忆召回 | 记忆读取，推理前置增强 | 按语义/关键词/标签检索长期记忆库，相关片段注入系统Prompt | 决定 Agent 对用户理解深度，个性化响应的基础 |
| 推理与行动 | 核心任务执行，记忆落地应用 | 思考‑工具调用‑观察循环，结合记忆适配用户口径输出 | 复用原生 Agent 的能力，完成记忆驱动的任务执行 |
| 编码与写入 | 记忆写入，单轮信息沉淀 | 筛选有效信息、记忆分类、重要度打分、打标签，结构化存储 | 沉淀记忆，业务逻辑与记忆处理职责解耦 |
| 整合与遗忘 | 记忆维护，库生命周期管理 | 去重、冲突消解、时间衰减、淘汰低价值过期记忆 | 防止记忆库膨胀，规避旧信息干扰推理 |


## 基于 LangChain 的代码演示

Memory-Augmented 就是为这种"长期 · 连续 · 按口径"的任务而生。本篇选择油菜表型采集作为示例，定义四个工具：

- 两个业务工具：
    - `query_phenotype_records` - 查询表型采集记录
    - `query_collection_task` - 查询采集任务进度与质控状态
- 两个记忆工具：
    - `search_long_term_memory` - 语义检索长期记忆
    - `write_long_term_memory` - 写入一条结构化记忆

### 定义业务工具：表型采集记录与采集任务

**业务工具**

```python
import json
import re
from datetime import datetime
from pathlib import Path

from langchain.tools import tool
from pydantic import BaseModel, Field

PHENO_RECORDS = {
    "P-2025-0412": {
        "材料编号": "YLY-2025-0412",
        "品种(系)": "秦油2026",
        "试点": "陕西杨凌试验站 3 号田",
        "播期": "2024-09-28",
        "取样株数": 10,
        "采集日期": "2025-05-20",
        "采集人": "孙老师 / 王技术员",
        "数据状态": "已复核",
        "性状": {
            "5.7 盛花期": "2025-03-30",
            "5.9 成熟期": "2025-05-22",
            "5.43 株高(cm)": 168.5,
            "5.44 一次分枝数(个)": 8.2,
            "5.46 有效分枝高度(cm)": 42.3,
            "5.47 主轴有效长度(cm)": 62.7,
            "5.48 主轴有效角果数(个)": 76.4,
            "5.49 全株有效角果数(个)": 342.6,
            "5.54 角果长度(cm)": 6.8,
            "5.59 每果粒数(粒)": 21.3,
            "菌核病发病率(%)": 3.2,
        },
    },
    "P-2025-0508": {
        "材料编号": "YLY-2025-0508",
        "品种(系)": "秦油33",
        "试点": "陕西杨凌试验站 5 号田",
        "播期": "2024-09-30",
        "取样株数": 8,
        "采集日期": "2025-05-24",
        "采集人": "李技术员",
        "数据状态": "待复核(花期缺采 2 项)",
        "性状": {
            "5.9 成熟期": "2025-05-26",
            "5.43 株高(cm)": 175.2,
            "5.44 一次分枝数(个)": 7.1,
            "5.47 主轴有效长度(cm)": 58.4,
            "5.49 全株有效角果数(个)": 268.3,
            "5.54 角果长度(cm)": 6.1,
            "5.59 每果粒数(粒)": 18.6,
            "菌核病发病率(%)": 9.7,
        },
    },
}

COLLECTION_TASKS = {
    "T-2025-03": {
        "任务批次": "T-2025-03",
        "采集时期": "成熟期",
        "试点": "陕西杨凌试验站",
        "计划材料数": 120,
        "已采材料数": 96,
        "计划完成日": "2025-05-25",
        "采集方法": "每材料随机取 10 株,按描述规范逐性状量测,取算术平均",
        "设备": "塔尺 + 游标卡尺 + 手持 PDA 录入",
        "质控规则": "双人复核;异常值(> 均值 ±3SD)现场重测",
        "缺项材料": ["YLY-2025-0508(倒伏,花期缺采)", "YLY-2025-0631(未达成熟期)"],
    },
}


class PhenoInput(BaseModel):
    batch_id: str = Field(description="表型采集批次号,如 'P-2025-0412'")


@tool(
    "query_phenotype_records",
    description="按采集批次号查询油菜表型采集记录(描述符编号、性状、数值、单位、取样株数、采集人、数据状态)。",
    args_schema=PhenoInput,
)
def query_phenotype_records(batch_id: str) -> str:
    rec = PHENO_RECORDS.get(batch_id)
    if not rec:
        return f"未找到采集批次 '{batch_id}'。"
    lines = [
        f"采集批次 {batch_id} 表型记录:",
        f"- 材料编号: {rec['材料编号']}",
        f"- 品种(系): {rec['品种(系)']}",
        f"- 试点: {rec['试点']}",
        f"- 取样株数: {rec['取样株数']} 株",
        f"- 采集日期/采集人: {rec['采集日期']} / {rec['采集人']}",
        f"- 数据状态: {rec['数据状态']}",
        "- 性状:",
    ]
    for k, v in rec["性状"].items():
        lines.append(f"    {k}: {v}")
    return "\n".join(lines)


class TaskInput(BaseModel):
    task_id: str = Field(description="采集任务批次号,如 'T-2025-03'")


@tool(
    "query_collection_task",
    description="按任务批次号查询表型采集任务进度(时期、已采/计划、采集方法、设备、质控规则、缺项材料)。",
    args_schema=TaskInput,
)
def query_collection_task(task_id: str) -> str:
    task = COLLECTION_TASKS.get(task_id)
    if not task:
        return f"未找到采集任务 '{task_id}'。"
    lines = [f"采集任务 {task_id}:"]
    for k, v in task.items():
        if isinstance(v, list):
            lines.append(f"- {k}: {'; '.join(v)}")
        else:
            lines.append(f"- {k}: {v}")
    return "\n".join(lines)
```

### 长期记忆存储与检索层

```python
import json
import os
import time
from typing import Literal
from dataclasses import dataclass, asdict

def calc_jaccard(set_a: set, set_b: set) -> float:
    """计算 Jaccard 系数"""
    if not set_a or not set_b:
        return 0.0
    inter = len(set_a & set_b)
    union = len(set_a | set_b)
    return inter / union


def time_decay(create_ts: float, halflife_days: float = 30.0) -> float:
    """计算时间衰减因子"""
    now = time.time()
    delta_sec = now - create_ts
    delta_day = delta_sec / 86400.0
    return 0.5 ** (delta_day / halflife_days)


@dataclass
class Memory:
    """单条长期记忆实体"""
    id: int
    content: str
    memory_type: Literal["declarative", "procedural", "reflective", "working"]
    tags: list[str]
    importance: int # 1~5
    create_time: float
    update_time: float

    @classmethod
    def from_dict(cls, d: dict[str, any]) -> "Memory":
        """字典转 Memory 对象"""
        return cls(
            id=d["id"],
            content=d["content"],
            memory_type=d["memory_type"],
            tags=d["tags"],
            importance=d["importance"],
            create_time=d["create_time"],
            update_time=d["update_time"]
        )

    def to_dict(self) -> dict[str, any]:
        """Memory对象转字典，用于 JSON 存储"""
        return asdict(self)

    def to_md(self) -> str:
        """Memory 对象转 Markdown 格式"""
        lines = [
            f"记忆 {self.id}",
            f"- 内容: {self.content}",
            f"- 类型: {self.memory_type}",
            f"- 标签: {'; '.join(self.tags)}",
            f"- 重要性: {self.importance}",
        ]
        return "\n".join(lines)

class MemoryStore:
    MEMORY_FILE = os.path.join(os.path.dirname(os.path.abspath(__file__)), "long_term_memory.json")
    """长期记忆存储器"""
    def __init__(self):
        self.init_memory()

    def init_memory(self):
        if not os.path.exists(self.MEMORY_FILE):
            with open(self.MEMORY_FILE, "w", encoding="utf-8") as f:
                json.dump({"memories": []}, f, ensure_ascii=False, indent=2)

    def load_memories(self) -> list[Memory]:
        try:
            with open(self.MEMORY_FILE, "r", encoding="utf-8") as f:
                data = json.load(f)
        except (json.JSONDecodeError, FileNotFoundError):
            return []
        raw_list = data.get("memories", [])
        return [Memory.from_dict(item) for item in raw_list]

    def save_memories(self, memories: list[Memory]):
        dict_list = [m.to_dict() for m in memories]
        with open(self.MEMORY_FILE, "w", encoding="utf-8") as f:
            json.dump({"memories": dict_list}, f, ensure_ascii=False, indent=2)

    def encode_memory(
        self,
        content: str,
        memory_type: str,
        tags: list[str],
        importance: int
    ) -> Memory:
        now_ts = time.time()
        memories = self.load_memories()
        new_id = max((m.id for m in memories), default=0) + 1
        return Memory(
            id=new_id,
            content=content,
            memory_type=memory_type,
            tags=tags,
            importance=max(1, min(5, importance)),
            create_time=now_ts,
            update_time=now_ts
        )

    def write_memory(self, memory_item: Memory):
        memories = self.load_memories()
        memories.append(memory_item)
        self.save_memories(memories)


    def retrieve_memory(self, query_tags: list[str], top_k: int = 3) -> list[Memory]:
        """召回：综合得分 = jaccard * importance * time_decay"""
        memories = self.load_memories()
        q_tag_set = set(query_tags)
        scored = []
        for m in memories:
            jac = calc_jaccard(q_tag_set, set(m.tags))
            decay = time_decay(m.create_time)
            score = jac * m.importance * decay
            scored.append((score, m))
        scored.sort(key=lambda x: x[0], reverse=True)
        return [item[1] for item in scored[:top_k]]


    def consolidate_memory(self, decay_threshold: float = 0.15) -> tuple[int, int]:
        """离线整合遗忘：时间衰减过滤 + 简单内容去重"""
        memories = self.load_memories()
        keep: list[Memory] = []
        seen_content = set()
        for m in memories:
            dec = time_decay(m.create_time)
            score = dec * m.importance
            if score < decay_threshold:
                continue
            if m.content in seen_content:
                continue
            seen_content.add(m.content)
            keep.append(m)
        self.save_memories(keep)
        return len(memories), len(keep)

memory_store = MemoryStore()
```

* 基于 JSON 文件存储长期记忆
* · `Memory` 类定义了单条长期记忆的结构，包括 ID、内容、类型、标签、重要性、创建时间、更新时间
* `MemoryStore` 类负责 JSON 文件的读写，并提供记忆的编码、写入、召回、整合等功能

### 定义记忆工具：查询记忆与写入记忆

```python
class MemorySearchArgs(BaseModel):
    query_tags: list[str] = Field(description="查询标签")
    top_k: int = Field(default=3, description="召回数量")


@tool(
    "search_memories",
    description="在长期记忆库中按关键词/标签检测相关记忆",
    args_schema=MemorySearchArgs,
)
def search_memories(query_tags: list[str], top_k: int = 3) -> str:
    """搜索长期记忆"""
    memories = memory_store.retrieve_memory(query_tags, top_k)
    if not memories:
        return "未找到相关记忆"
    content = (
        f"找到 {len(memories)} 条相关记忆：\n"
        "\n"
        + "\n".join([m.to_md() for m in memories])
    )
    return content


class MemoryWriteArgs(BaseModel):
    content: str = Field(description="记忆内容")
    memory_type: Literal["declarative", "procedural", "reflective", "working"] = Field(
        description="记忆类型，可选值： declarative（陈述性记忆）, procedural（程序性记忆）, reflective（反思性记忆）, working（工作记忆）")
    tags: list[str] = Field(description="记忆标签")
    importance: int = Field(default=3, ge=1, le=5, description="记忆重要性 1-5 之间")


@tool(
    "write_memory",
    description="写入长期记忆",
    args_schema=MemoryWriteArgs,
)
def write_memory(content: str, memory_type: str, tags: list[str], importance: int = 3) -> str:
    memory = memory_store.encode_memory(content, memory_type, tags, importance)
    memory_store.write_memory(memory)
    return f"已写入记忆 {memory.id}"
```

### Agent 定义

```python
from langchain.agents import create_agent
from langchain_ollama import ChatOllama

llm = ChatOllama(model="qwen3.5:9b", temperature=0)

assistant_prompt = (
    "你是'孙老师的油菜表型采集助理',负责杨凌试验站油菜田间表型数据的查询、核对与"
    "归档。回答时:1) 优先调用工具拿到一手采集数据;2) 输出表格时严格按"
    "《油菜种质资源描述规范和数据标准》的描述符编号、性状名、数值、单位四列呈现;"
    "3) 未采集的性状一律标注'未采集',不得插值补数或用经验值填充;"
    "4) 用户的偏好/身份/格式要求已经在'已召回长期记忆'里给出,直接遵守;"
    "5) 如果用户表达了新偏好或新事实,请用 write_memory 工具把它落到长期记忆里。"
)

assistant = create_agent(
    model=llm,
    tools=[
        query_phenotype_records,
        query_collection_task,
        search_memories,
        write_memory,
    ],
    system_prompt=assistant_prompt
)

memory_writer_prompt = (
    "你是'记忆反思器'。给定用户的最新提问与助理的回答,你负责从中抽取"
    "值得长期保留的信息,并用 write_memory 工具写入长期记忆库。\n"
    "写入规则:\n"
    "1) 只写 durable 信息(用户身份、格式偏好、试验站约定、质控规则、品种特性结论等),"
    "不写一次性查询的流水账;\n"
    "2) 每条记忆给 1~3 个标签(便于后续 tag 召回),并打 1~5 的重要性;"
    "3) 若本轮没有值得长期保留的内容,则不要调用任何工具,直接回复'无需写入新记忆';\n"
    "4) 避免写入与已有记忆明显重复的内容。"
)

memory_writer = create_agent(
    model=llm,
    tools=[write_memory],
    system_prompt=memory_writer_prompt
)
```

### 问答与反思

```python
import re


def _extract_tags(text: str) -> list[str]:
    """极简标签抽取:按中英文 token 粗分,供 MemoryStore 的 Jaccard tag 召回使用。
    生产中可替换为 embedding 语义召回,此处保持与存储器的 tag 召回协议一致。"""
    tokens = re.findall(r"[一-龥A-Za-z0-9]+", text)
    return [t for t in tokens if len(t) >= 2][:10]


def run_with_memory(question: str, top_k: int = 3) -> str:
    recalled_raw = search_memories.invoke({
        "query_tags": _extract_tags(question), 
        "top_k": top_k,
    })
    recalled_text = recalled_raw if isinstance(recalled_raw, str) else str(recalled_raw)

    answer_resp = assistant.invoke(
        {"messages": [
            {"role": "user", "content": f"已召回长期记忆:\n{recalled_text}\n\n用户提问: {question}"}
        ]}
    )
    answer = answer_resp["messages"][-1].content

    memory_writer.invoke(
        {"messages": [
            {"role": "user", "content": f"用户提问: {question}\n助理回答: {answer}\n请判断是否需要写入长期记忆。"}
        ]}
    )

    return answer
```

* `_extract_tags` 用于从用户提问中提取标签,用于召回相关记忆
* `run_with_memory` 主流程:召回记忆 -> 助理作答 -> 反思写入

### 调用

```python
if __name__ == "__main__":
    # 演示:多轮对话,观察记忆如何在轮次间被召回与累积
    tasks = [
        "查一下 P-2025-0412 这批的株高和每果粒数,按描述规范四列给我。",
        "我习惯表格里数值带单位,而且我是孙老师团队的人,记住这个偏好。",
        "那 P-2025-0508 这批的菌核病发病率多少?对照我之前说的格式来。",
    ]
    for t in tasks:
        print("=" * 60)
        print(f"用户: {t}")
        print("-" * 60)
        print(run_with_memory(t))

    # 演示离线整合(遗忘)低频记忆
    total, kept = memory_store.consolidate_memory()
    print("=" * 60)
    print(f"记忆整合完成: 整合前 {total} 条 -> 整合后 {kept} 条")
```

执行输出：

```text
============================================================
用户: 查一下 P-2025-0412 这批的株高和每果粒数,按描述规范四列给我。
------------------------------------------------------------
根据采集批次 P-2025-0412 的表型记录，查询结果如下：

| 描述符编号 | 性状名 | 数值 | 单位 |
|------------|--------|------|------|
| 5.43 | 株高 | 168.5 | cm |
| 5.59 | 每果粒数 | 21.3 | 粒 |

以上数据来自杨凌试验站 3 号田秦油2026品种，取样株数10株，采集日期2025-05-20，数据状态已复核。
============================================================
用户: 我习惯表格里数值带单位,而且我是孙老师团队的人,记住这个偏好。
------------------------------------------------------------
好的，孙老师！您的偏好和身份已记录：

- **表格格式**：数值带单位（如株高: 120 cm）
- **团队身份**：孙老师团队成员

后续查询表型数据时，我会严格按此标准输出表格。如有其他需要调整的地方，随时告诉我！
============================================================
用户: 那 P-2025-0508 这批的菌核病发病率多少?对照我之前说的格式来。
------------------------------------------------------------
根据采集批次 **P-2025-0508** 的表型记录，菌核病发病率为：

| 描述符编号 | 性状名 | 数值 | 单位 |
|:---|:---|:---|:---|
| 未指定 | 菌核病发病率 | 9.7 | % |

该批次数据状态为“待复核”，花期缺采2项，但菌核病发病率已采集完成。
============================================================
记忆整合完成: 整合前 2 条 -> 整合后 2 条
```

查看 `long_term_memory.json` 文件，即可看到所有长期记忆，如下：

```json
{
  "memories": [
    {
      "id": 1,
      "content": "用户偏好：表格中数值需带单位；用户身份：孙老师团队成员",
      "memory_type": "declarative",
      "tags": [
        "孙老师团队",
        "数据格式偏好",
        "用户身份"
      ],
      "importance": 5,
      "create_time": 1788404930.074628,
      "update_time": 1788404930.074628
    },
    {
      "id": 2,
      "content": "用户身份：孙老师团队成员；格式偏好：表格中数值需带单位（如株高: 120 cm），后续查询表型数据时严格按此标准输出",
      "memory_type": "declarative",
      "tags": [
        "用户身份",
        "格式偏好"
      ],
      "importance": 3,
      "create_time": 1788404947.558461,
      "update_time": 1788404947.558461
    }
  ]
}
```

## 优缺点

### 优点

- **跨会话连续性**：用户不必重复交代背景，第二次开口就能接着上次
- **个性化与低成本适配**：记住偏好、术语、口径，无需微调即可按用户习惯输出
- **经验可累积**：事实、结论、流程沉淀成可检索资产，越用越准
- **上下文瘦身**：长期信息外置到记忆库，避免每次全量塞进提示词

### 缺点

- **检索质量决定上限**：召回错了等于注入噪声，比没有记忆更糟
- **记忆污染与冲突**：错误信息一旦写入会被反复引用，新旧冲突难自动判定
- **写入时机难把握**：什么都记会撑爆记忆库，什么都不记等于没有记忆
- **工程与合规复杂度**：存储、去重、过期、权限与删除机制成本与隐私风险并存

## 小结

Memory-Augmented 解决的是"Agent 记得你"的问题：把"事实 / 偏好 / 流程 / 经验"从对话里抽出来，结构化、可检索、可遗忘，跨会话复用。