+++
date = '2026-09-01T10:01:19+08:00'
draft = false
title = 'Router 架构详解：让 Agent 学会按需分发'
categories = ['人工智能', '智能体']
tags = ['AI Agent', 'Router', 'LangChain', 'Agent 架构']
toc = true
+++

> 系列文章第 4 篇 · AI Agent 六大主流架构拆解  
> 本篇拆解 Router，解决**一个入口、多个专长**的分发问题，并把它落到油菜种植技术咨询这一典型多领域场景上。

## 什么是 Router

核心思想一句话：**前置一个路由器，先把输入按意图/领域分类，再分发给对应的领域子 Agent，各自处理后再以统一格式返回。**

一个最小运行单元：

```
用户问题
  ↓
Router 路由器：判定领域
  ↓
领域子 Agent A / B / C
  ↓
汇聚输出
```

## 示例背景：油菜种植技术咨询

油菜种植涉及多个专业领域：栽培(播期、密度、灌溉)、植保(菌核病、蚜虫等病虫害防治)、土肥(施肥方案、硼肥等)。农户或农技人员问一个问题，往往只对应其中一个领域。

如果用一个通用 Agent 回答所有问题，会面临三个麻烦：

1. **提示词臃肿**：必须把栽培、植保、土肥三套知识、三套工具都塞进同一份 system prompt，互相干扰。
2. **工具冲突**：一个通用 Agent 可能同时拿到栽培工具和植保工具，在不确定领域时乱调用。
3. **效果天花板低**：每个领域都有专业术语与最佳实践，通用 Agent 很难在每个方向上都做深。

Router 的做法是：**对外只暴露一个问答入口，对内按问题领域分发给专门的栽培 Agent、植保 Agent、土肥 Agent。** 每个子 Agent 只配自己领域的提示词和工具，Router 则只做一件事——判定这个问题属于哪个领域。

## 架构与运行机制

![Router 运行机制](/imgs/learn-agents/router-mechanism.png)

运行机制三要素：

- **统一入口**：用户只面对一个入口，不需要提前选择领域。
- **Router 路由器**：根据输入语义判定领域/意图，返回路由结果(领域 + 理由)。
- **领域子 Agent + 结果汇聚**：被命中的子 Agent 用专属提示词和工具处理问题，未命中时走兜底策略(通用 Agent 或追问澄清)，最后统一格式返回给用户。

关键设计点：

- **Router 只做分类，不做业务逻辑**：它不应该自己生成最终答案，只决定交给哪个子 Agent。
- **子 Agent 之间尽量解耦**：每个子 Agent 的系统提示、工具集、输出格式只关心自己的领域，降低相互干扰。
- **路由判定的准确率是瓶颈**：一旦分错，后续子 Agent 再强也救不回来。因此 Router 的 prompt 要给出清晰的领域定义与示例，必要时加兜底或二次确认。
- **输出格式最好统一**：各子 Agent 返回的格式应一致(例如都用「结论 + 依据」)，否则前端拼接体验差。

## 基于 LangChain 的代码演示

### 定义三个领域工具

用 `@tool` + Pydantic `args_schema` 定义三个互不重迭的领域查询工具：

```python
from langchain.tools import tool
from pydantic import BaseModel, Field

CULTIVATION_DB = {
    "播种期": "长江中下游直播油菜适宜播种期为 9 月下旬至 10 月中旬，移栽油菜可提前 7–10 天。",
    "种植密度": "直播每亩留苗 2.0–2.5 万株，移栽每亩 0.6–0.8 万株。",
    "灌溉": "蕾薹期与花期是需水关键期，保持田间持水量 70% 左右。",
}

PEST_DB = {
    "菌核病": "油菜菌核病防治：花期用 50% 腐霉利可湿性粉剂 1000–1500 倍液喷雾，连防 2–3 次。",
    "蚜虫": "蚜虫可用 10% 吡虫啉可湿性粉剂 2000 倍液或 25% 噻虫嗪水分散粒剂 5000 倍液防治。",
    "菜青虫": "菜青虫可用苏云金杆菌(Bt)或氯虫苯甲酰胺防治。",
}

FERTILIZER_DB = {
    "基肥": "每亩施腐熟有机肥 1000–1500 kg + 复合肥(N-P₂O₅-K₂O 15-15-15)30–40 kg。",
    "追肥": "蕾薹期亩追尿素 5–8 kg，初花期喷施 0.2% 硼砂溶液防花而不实。",
    "硼肥": "缺硼土壤在苗期和初花期各喷一次 0.2% 硼砂溶液。",
}


class CultivationInput(BaseModel):
    keyword: str = Field(description="栽培关键词，如 '播种期'、'种植密度'、'灌溉'")

class PestInput(BaseModel):
    keyword: str = Field(description="病虫害关键词，如 '菌核病'、'蚜虫'、'菜青虫'")

class FertilizerInput(BaseModel):
    keyword: str = Field(description="施肥关键词，如 '基肥'、'追肥'、'硼肥'")


@tool(name_or_callable="query_cultivation_tech", args_schema=CultivationInput)
def query_cultivation_tech(keyword: str) -> str:
    """查询油菜栽培技术要点。"""
    if keyword in CULTIVATION_DB:
        return f"栽培技术：{keyword} — {CULTIVATION_DB[keyword]}"
    return f"未找到 '{keyword}' 的栽培技术，请尝试：{list(CULTIVATION_DB.keys())}"


@tool(name_or_callable="query_pest_control", args_schema=PestInput)
def query_pest_control(keyword: str) -> str:
    """查询油菜病虫害防治方案。"""
    if keyword in PEST_DB:
        return f"病虫害防治：{keyword} — {PEST_DB[keyword]}"
    return f"未找到 '{keyword}' 的防治方案，请尝试：{list(PEST_DB.keys())}"


@tool(name_or_callable="query_fertilizer_plan", args_schema=FertilizerInput)
def query_fertilizer_plan(keyword: str) -> str:
    """查询油菜施肥方案。"""
    if keyword in FERTILIZER_DB:
        return f"施肥方案：{keyword} — {FERTILIZER_DB[keyword]}"
    return f"未找到 '{keyword}' 的施肥方案，请尝试：{list(FERTILIZER_DB.keys())}"
```

### 领域子 Agent 挂载工具

三个子 Agent 各挂一个领域工具：

```python
from langchain.agents import create_agent

llm = ChatOllama(model="qwen3.5:9b", temperature=0)

cultivation_agent = create_agent(
    model=llm,
    tools=[query_cultivation_tech],
    system_prompt="你是油菜栽培技术专家。用户问栽培/农事/播种/密度/灌溉等问题，调用 query_cultivation_tech 查询后给出简明建议。",
)

pest_agent = create_agent(
    model=llm,
    tools=[query_pest_control],
    system_prompt="你是油菜植保专家。用户问病虫害识别与防治，调用 query_pest_control 查询后给出防治建议。",
)

fertilizer_agent = create_agent(
    model=llm,
    tools=[query_fertilizer_plan],
    system_prompt="你是油菜土肥专家。用户问施肥方案/肥料/追肥/硼肥，调用 query_fertilizer_plan 查询后给出建议。",
)
```

### 构建 Router + 三个领域子 Agent

Router 本身也是一个 Agent，只负责输出结构化的路由决策：

```python
from langchain.agents.structured_output import ToolStrategy
from langchain_ollama import ChatOllama
from pydantic import BaseModel, Field
from typing import Literal


class RouteDecision(BaseModel):
    """路由决策结果。"""
    route: Literal["栽培", "植保", "土肥", "其他"] = Field(
        description="问题所属领域：栽培(播种/密度/灌溉等农事)、植保(病虫害防治)、土肥(施肥方案)；与三者均无关时填 '其他'"
    )
    reason: str = Field(description="一句话说明判定理由")


ROUTER_PROMPT = (
    "你是油菜种植技术咨询的领域路由器。"
    "判断用户问题属于栽培、植保、土肥中的哪一个领域，无法归类时选 '其他'。"
    "只做分类，不要回答用户的问题本身。"
)

router = create_agent(
    model=llm,
    tools=[],
    system_prompt=ROUTER_PROMPT,
    response_format=ToolStrategy(RouteDecision),
)

AGENTS = {
    "栽培": cultivation_agent,
    "植保": pest_agent,
    "土肥": fertilizer_agent,
}
```

### 定义问答入口函数

```python
def route_and_answer(question: str) -> str:
    result = router.invoke({"messages": [{"role": "user", "content": question}]})
    decision: RouteDecision = result["structured_response"]

    if decision.route not in AGENTS:
        return f"无法识别问题领域(识别为：{decision.route})，请补充描述或换个说法。理由：{decision.reason}"

    answer = AGENTS[decision.route].invoke({
        "messages": [{"role": "user", "content": question}]
    })["messages"][-1].content

    return f"【路由到：{decision.route}】依据：{decision.reason}\n\n{answer}"
```

函数返回路由决策结果及对应领域子 Agent 的回答。

### 运行：多领域问题自动分发

```python
if __name__ == "__main__":
    print(route_and_answer("油菜播种期"))
    print(route_and_answer("油菜菌核病怎么防治"))
    print(route_and_answer("油菜施肥方案"))
```

执行输出：

```text
【路由到：栽培】依据：播种期属于油菜种植农事操作中的播种管理范畴

关于油菜播种期的建议：

**长江中下游地区：**
- **直播油菜**：适宜在 **9月下旬至10月中旬** 进行播种
- **移栽油菜**：可提前 **7–10天**（约9月中下旬）完成移栽

**注意事项：**
- 播种期应结合当地气候条件灵活调整，避免过早或过晚影响产量
- 直播油菜需确保土壤墒情适宜，移栽油菜则需预留缓苗时间

如需了解其他栽培技术要点（如种植密度、灌溉管理等），欢迎继续提问！
【路由到：植保】依据：油菜菌核病属于病虫害防治范畴

油菜菌核病是油菜生产中的常见病害，主要危害花器和茎秆，严重时会导致“烂腰”和减产。以下是综合防治建议：

### 一、化学防治（关键时期）
- **最佳时机**：**花期**是防治的关键窗口期，此时病菌孢子活跃，易侵染。
- **推荐药剂**：50% 腐霉利可湿性粉剂 1000–1500 倍液喷雾。
- **施药次数**：连防 2–3 次，间隔 7–10 天一次。
- **注意事项**：喷药时注意均匀覆盖花部，避免重喷漏喷。

### 二、农业防治（基础措施）
- **合理轮作**：与非十字花科作物（如玉米、小麦等）轮作，减少土壤中病菌积累。
- **清除病残体**：收获后及时清理田间病株残体，集中销毁，降低菌源。
- **合理密植**：避免种植过密，保持田间通风透光，降低湿度。
- **排水降湿**：雨季注意排水，避免田间积水，减少高湿环境利于发病的条件。

### 三、其他建议
- 若病情较重，可结合使用其他杀菌剂（如多菌灵、咪鲜胺等）轮换使用，延缓抗药性产生。
- 选用抗病品种也是预防病害的重要措施。

通过“农业防治 + 化学防治”相结合，可有效控制油菜菌核病的发生与蔓延。
【路由到：土肥】依据：问题涉及油菜施肥方案，属于土肥领域

油菜基肥施肥建议如下：

- **有机肥**：每亩施用腐熟有机肥 1000–1500 kg；
- **复合肥**：每亩施用 N-P₂O₅-K₂O 比例为 15-15-15 的复合肥 30–40 kg。

如果您还想了解追肥、硼肥或其他施肥阶段的方案，请随时告诉我！
```

## 优缺点

### 优点

- **一个入口，多个专长**：产品对外统一，对内可挂载任意多个垂直 Agent，扩展性强。
- **提示词与工具解耦**：每个子 Agent 只需承载一个领域的知识，避免提示臃肿和工具冲突。
- **效果上限更高**：子 Agent 可以为单一职责深度优化(提示词、工具、示例)。
- **可独立迭代**：新增领域只需新增一个子 Agent + 一条路由规则，不影响已有领域。

### 缺点

- **路由准确率是瓶颈**：一旦分错，后续回答必然跑偏；模糊问题或跨领域问题尤其容易误判。
- **子 Agent 之间默认不通信**：一个复杂问题可能需要栽培 + 植保 + 土肥联合回答，Router 只做单选，需要更高层编排(这就是 Multi-Agent 要解决的问题)。
- **领域切分是隐性成本**：哪些领域独立、边界在哪、是否允许重叠，需要人工设计与持续维护。
- **输出格式需治理**：子 Agent 返回格式不一致时，前端拼接体验差。