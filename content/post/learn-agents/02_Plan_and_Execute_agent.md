+++
date = '2026-08-21T10:52:54+08:00'
draft = false
title = 'Plan-and-Execute 架构详解：先规划,再执行'
categories = ['人工智能', '智能体']
tags = ['AI Agent', 'Plan-and-Execute', 'LangChain', 'Agent 架构']
toc = true
+++

> 系列文章第 2 篇 · AI Agent 六大主流架构拆解  
> 本篇拆解 Plan-and-Execute，把“规划”与“执行”分离，并把它落到**油菜新品种区域试验评价**这一长周期多阶段任务上(场景与第 1 篇的「种质性状查询」互补)。

## 什么是 Plan-and-Execute

Plan-and-Execute 核心思想一句话：**先用 Planner 一次性把任务拆成有序子任务列表，再用 Executor 逐项执行，中间结果可回灌，必要时触发 Re-Plan 修订后续计划。**

一个最小运行单元：

```
用户问题
  → Planner：生成 [Step 1, Step 2, Step 3, ...]
  → Executor：逐项执行 Step，调用工具，产出中间结果
  → (可选) Re-Plan：根据中间结果修订剩余计划
  → 汇聚 → 最终答案
```

与 ReAct 的关键区别：

- ReAct 是“边走边看” —— 每一步都在 Thought 里决定下一步，**缺乏全局视野**，长任务容易绕远路，且每步全量推理，token 随步数爆炸。
- Plan-and-Execute 是“先有蓝图再施工” —— Planner 先给出全局计划，Executor 只管按图施工，必要时回头改图。

一句话总结：ReAct 是**思考驱动执行**，Plan-and-Execute 是**规划驱动执行**。

## 品种区域试验评价

油菜新品种从育成到审定，必须经过**区域试验**这一关键长周期环节：在多个生态区试点、多年度重复，系统评价新品种的丰产性、稳定性、适应性，作为能否推广的决策依据。整个评价流程是一个典型的**多阶段、有明确工序、跨季度**的工作流：

**“汇总各试点多年试验数据 → 稳定性参数计算 → 生态适应性区划查询 → 综合丰产/稳产/适配评价 → 给出推广建议”**

这种特性正好契合 Plan-and-Execute：

1. **多阶段有依赖**：必须先汇总数据，才能算稳定性;先有稳定性，才能做适应性判断——前一阶段输出是后一阶段输入。
2. **可中途修订**：某试点数据缺失或某品种稳定性不达标，可 Re-Plan 调整后续步骤(剔除该品种、补充试点查询)。
3. **可审计**：区试评价报告要能存档、评审，Plan 本身就是可交付物。

## 架构与运行机制

![Plan-and-Execute 运行机制](/imgs/learn-agents/plan-execute-mechanism.png)

- **Planner 规划器**：接收用户目标，一次性输出一份有序子任务列表(可含步骤间依赖)。
- **Plan 任务列表**：Planner 产出的有序步骤，既是执行的蓝图，也是可审计、可回放的中间产物。
- **Executor 执行器**：逐项执行每个 Step，每步可调用外部工具(检索、计算、查库等)获取新信息并产出中间结果;工程上常用一个独立子 Agent 作为执行内核。
- **Re-Plan 重规划**：每完成一步，根据已完成步骤与中间结果，决定是继续按原计划、插入新步骤、修订剩余计划，还是收尾——而非死守初始计划。

关键设计点：

- **Planner 与 Executor 解耦**，可分别选型——Planner 用强模型保证规划质量，Executor 用快模型降低成本。
- **每步结果回灌**到上下文，供后续步骤消费，也供 Re-Plan 判断计划是否需要修订。
- **Re-Plan 是兜底**：没有它，Plan-and-Execute 就退化为“一次死计划”;有了它，即便 Planner 第一次拆得不够准，执行中也能动态修订。

## 基于 LangChain 的代码演示

### 日志配置

```python
import logging

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s | %(levelname)-7s | %(name)s | %(message)s",
    datefmt="%H:%M:%S",
)
logger = logging.getLogger("Plan-and-Execute Agent")
```

### 定义区试评价工具

用 `@tool` 装饰器加 Pydantic `args_schema` 定义三个工具：
1. 多地点试验数据查询
2. 稳定性参数计算
3. 适应性区划查询。

```python
from langchain.tools import tool
from pydantic import BaseModel, Field

# 模型区域试验数据：多试点 x 多年份 x 多表型
TRIAL_DB = {
    "陕油28": {
        "武汉2023": {"产量kg/亩": 212, "含油量%": 45.2, "芥酸%": 0.8, "硫苷μmol/g": 22, "千粒重g": 4.2, "株高cm": 168},
        "武汉2024": {"产量kg/亩": 208, "含油量%": 45.5, "芥酸%": 0.7, "硫苷μmol/g": 21, "千粒重g": 4.1, "株高cm": 170},
        "长沙2023": {"产量kg/亩": 205, "含油量%": 44.8, "芥酸%": 0.9, "硫苷μmol/g": 23, "千粒重g": 4.0, "株高cm": 172},
        "长沙2024": {"产量kg/亩": 210, "含油量%": 45.1, "芥酸%": 0.8, "硫苷μmol/g": 22, "千粒重g": 4.2, "株高cm": 169},
        "南昌2023": {"产量kg/亩": 198, "含油量%": 44.5, "芥酸%": 1.0, "硫苷μmol/g": 24, "千粒重g": 3.9, "株高cm": 174},
        "南昌2024": {"产量kg/亩": 204, "含油量%": 44.9, "芥酸%": 0.9, "硫苷μmol/g": 23, "千粒重g": 4.1, "株高cm": 171},
    },
    "华油杂62R": {
        "武汉2023": {"产量kg/亩": 206, "含油量%": 46.0, "芥酸%": 1.2, "硫苷μmol/g": 25, "千粒重g": 3.8, "株高cm": 162},
        "武汉2024": {"产量kg/亩": 210, "含油量%": 46.2, "芥酸%": 1.1, "硫苷μmol/g": 24, "千粒重g": 3.9, "株高cm": 160},
        "长沙2023": {"产量kg/亩": 200, "含油量%": 45.8, "芥酸%": 1.3, "硫苷μmol/g": 26, "千粒重g": 3.7, "株高cm": 164},
        "长沙2024": {"产量kg/亩": 203, "含油量%": 46.0, "芥酸%": 1.2, "硫苷μmol/g": 25, "千粒重g": 3.8, "株高cm": 163},
        "南昌2023": {"产量kg/亩": 195, "含油量%": 45.5, "芥酸%": 1.4, "硫苷μmol/g": 27, "千粒重g": 3.6, "株高cm": 166},
        "南昌2024": {"产量kg/亩": 199, "含油量%": 45.7, "芥酸%": 1.3, "硫苷μmol/g": 26, "千粒重g": 3.7, "株高cm": 165},
    },
}

# 模型品种适应性数据
ADAPTABILITY_DB = {
    "陕油28": {"长江中游": "适宜", "长江下游": "适宜", "黄淮": "一般"},
    "华油杂62R": {"长江中游": "适宜", "长江下游": "适宜", "黄淮": "不适宜"},
}


class TrialDataQueryArgs(BaseModel):
    name: str = Field(description="油菜品种")


@tool(
    "query_trial_data",
    description="查询某油菜品种在多个区域试验点的多年试验数据,返回产量、含油量、芥酸、硫苷、千粒重、株高等表型记录。",
    args_schema=TrialDataQueryArgs
)
def query_trial_data(name: str) -> str:
    data = TRIAL_DB.get(name)
    if not data:
        lines = [f"未找到品种 {name} 的区试数据，可选品种："]
        lines.extend([f" - {k}" for k in TRIAL_DB.keys()])
        return "\n".join(lines)

    lines = [f"{name} 区试数据："]
    for site, rec in data.items():
        lines.append(
            f"- {site}：产量 {rec['产量kg/亩']} kg/亩，含油量 {rec['含油量%']}%，"
            f"芥酸 {rec['芥酸%']}%，硫苷 {rec['硫苷μmol/g']} μmol/g，"
            f"千粒重 {rec['千粒重g']} g，株高 {rec['株高cm']} cm"
        )
    return "\n".join(lines)


class CalcStabilityArgs(BaseModel):
    name: str = Field(description="油菜品种")


@tool(
    "calc_stability",
    description="根据某品种多地点多年产量数据，计算丰产性(均值)与稳定性参数(变异系数 CV)，并给出稳定性评价。",
    args_schema=CalcStabilityArgs,
)
def calc_stability(name: str) -> str:
    data = TRIAL_DB.get(name)
    if not data:
        lines = [f"未找到品种 {name} 的区试数据，可选品种："]
        lines.extend([f" - {k}" for k in TRIAL_DB.keys()])
        return "\n".join(lines)

    yields = [rec["产量kg/亩"] for rec in data.values()]
    mean = sum(yields) / len(yields)
    var = sum((y - mean) ** 2 for y in yields) / len(yields)
    cv = (var ** 0.5) / mean
    stability = "稳定" if cv < 0.05 else "较稳定" if cv < 0.10 else "欠稳定"
    return f"{name} 多点平均产量：{mean:.1f} kg/亩，变异系数 CV：{cv:.3f}，稳定性：{stability}"


class AdaptabilityZoneQueryArgs(BaseModel):
    name: str = Field(description="油菜品种")


@tool(
    name_or_callable="query_adaptability_zone",
    description="查询某品种在各主要油菜生态区的适应性评价(适宜 / 一般 / 不适宜)。",
    args_schema=AdaptabilityZoneQueryArgs,
)
def query_adaptability_zone(name: str) -> str:
    """查询某品种的生态适应性区划。"""
    zones = ADAPTABILITY_DB.get(name)
    if not zones:
        lines = [f"未找到品种 {name} 的适应性数据，可选品种："]
        lines.extend([f" - {k}" for k in ADAPTABILITY_DB.keys()])
        return "\n".join(lines)

    lines = [f"{name} 适应性区划:"]
    for zone, adaptability in zones.items():
        lines.append(f"- {zone}:{adaptability}")
    return "\n".join(lines)


tools = [query_trial_data, calc_stability, query_adaptability_zone]
```

### Executor 实现

用 `create_agent` 把三个工具装进一个子 Agent，作为 Executor。

```python
from langchain.chat_models import BaseChatModel
from langchain.agents import create_agent
from langchain.messages import SystemMessage, HumanMessage
from langchain_ollama import ChatOllama

llm = ChatOllama(model="qwen3.5:9b", temperature=0)

class Executor:
    def __init__(self, llm: BaseChatModel, *, tools: list, system_prompt: str):
        self.agent = create_agent(
            llm,
            tools=tools,
            system_prompt=SystemMessage(system_prompt),
        )

    def invoke(self, question: str, step: str) -> str:
        logger.info(f"调用执行器，步骤：{step}")
        content = f"用户提问：{question}，当前执行步骤为：{step}"
        resp = self.agent.invoke({"messages": [HumanMessage(content)]})
        content = resp["messages"][-1].content
        logger.info(f"步骤（{step}）的执行结果：{content[:200]}")
        return content
```

Executor 的作用是对于特定的一个步骤，模型返回结果，必要时调用工具。

### Planner 实现

对于 `Pydantic` 的类型提示功能，让模型返回结构化的数据。

```python
class Step(BaseModel):
    description: str = Field(description="步骤描述")
    done: bool = Field(default=False, description="是否完成，True 完成，False 未完成")

    def __str__(self):
        return self.description


class PlannerResult(BaseModel):
    steps: list[Step] = Field(description="步骤数据")


class ReplanResult(BaseModel):
    finished: bool = Field(description="是否已掌握足够信息可作答，True 表示结束，False 表示还有下一步")
    next_step: str = Field(default="", description="下一个步骤描述；若 finished 为 True 则为空字符串")


class FinalAnswerResult(BaseModel):
    answer: str = Field(description="基于所有步骤执行结果，对原始问题的最终回答")


class Planner:
    """规划师 + 重规划师 + 综合器：制定计划、判断下一步、汇总最终回答。"""

    def __init__(self, llm, *, planner_prompt: str, replanner_prompt: str,
                 synthesizer_prompt: str, executor: Executor):
        # 规划
        self.planner_agent = create_agent(
            llm,
            system_prompt=SystemMessage(planner_prompt),
            response_format=PlannerResult,
        )
        # 重新规划
        self.replanner_agent = create_agent(
            llm,
            system_prompt=SystemMessage(replanner_prompt),
            response_format=ReplanResult,
        )
        # 总结
        self.synthesizer_agent = create_agent(
            llm,
            system_prompt=SystemMessage(synthesizer_prompt),
            response_format=FinalAnswerResult,
        )
        self.executor = executor
        # 存储用户提问
        self.question = None
        # 存储步骤
        self.steps: list[Step] = []
        # 存储每一个步骤由模型返回的内容
        self.contents: list[str] = []
        # 存储最终答案
        self.final_answer: str = ""

    def make_plan(self, question: str):
        self.question = question
        resp = self.planner_agent.invoke({"messages": [HumanMessage(self.question)]})
        result: PlannerResult = resp["structured_response"]

        logger.info(f"规划师生成 {len(result.steps)} 个步骤：")
        for idx, step in enumerate(result.steps, start=1):
            logger.info(f"步骤 {idx} ：{step}")
        # 保存模型返回的步骤
        self.steps = result.steps
        return self

    def call_executor(self):
        for idx, step in enumerate(self.steps):
            content = self.executor.invoke(self.question, step.description)
            self.contents.append(content)
            step.done = True
            replan_result = self.replan()
            if replan_result.finished:
                self.steps = self.steps[:idx]
                break
            else:
                self.steps.insert(idx + 1, Step(description=replan_result.next_step))

        self.final_answer = self.synthesize()
        return self

    def replan(self) -> ReplanResult:
        """根据原始目标与已执行步骤的结果，判断是否可以结束，或给出下一步。"""
        # 已完成的步骤及其结果：steps 与 contents 下标一一对应
        log_lines = [f"原始目标：{self.question}", "", "已完成步骤及结果："]
        for i, (step, content) in enumerate(zip(self.steps, self.contents), start=1):
            log_lines.append(f"{i}. {step.description}\n   结果：{content}")
        log_lines.append(
            "\n请判断：是否已掌握足够信息来回答原始目标？"
            "若已足够请设置 finished=True（next_step 留空）；"
            "否则给出唯一一个最有价值的下一步（finished=False）。不要重复已完成的步骤。"
        )
        prompt = "\n".join(log_lines)

        resp = self.replanner_agent.invoke({"messages": [HumanMessage(prompt)]})
        result: ReplanResult = resp["structured_response"]

        if result.finished:
            logger.info("重规划：已结束（无需更多步骤）")
        else:
            logger.info(f"重规划：下一步 = {result.next_step}")
        return result

    def synthesize(self) -> str:
        step_results = "\n".join([f"""{idx}. {step}

{self.contents[idx]}""" for idx, step in enumerate(self.steps, start=1)])
        content = f"""用户问题：{self.question}，以下是各步骤及对应的结果：

{step_results}

基于以上内容，总结并回答用户问题。"""
        resp = self.synthesizer_agent.invoke({"messages": HumanMessage(content)})
        final_answer: FinalAnswerResult = resp["structured_response"]
        return final_answer.answer

    def get_final_answer(self) -> str:
        return self.final_answer
```

类与重要方法说明：
* `Step`：模型返回单个步骤的结构化表示
* `PlannerResult`：Planer Agent 返回的结构化表示，包含多个 `Step`
* `ReplanResult`：Replanner Agent 返回的结构化表示，提示信息已完备，或给出下一步骤
* `FinalAnswerResult`：Synthesizer Agent 返回的结构化表示，为最终回答
* `Planner`：包含三个 Agent，分别是 Planer Agent、Replanner Agent 和 Synthesizer Agent，还包含 Executor 实例
    * `make_plan` 方法：调用 Planer Agent，返回步骤
    * `call_executor` 方法：调用 Executor，实际行为是为当前步骤请求模型并返回结果
    * `replan` 方法：调用 Replanner Agent，返回规划结束或执行下一步骤
    * `synthesize` 方法：调用 Synthesizer Agent，返回最终的回答

### query 方法

定义 `query` 方法，接收用户提问。

```python
def query(question: str) -> str:
    executor = Executor(
        llm,
        tools=tools,
        system_prompt="你是油菜区试评价的步骤执行器，职责是调用工具执行给定的单个步骤，并给出该步骤的简明结果。",
    )
    planner = Planner(
        llm,
        planner_prompt=(
            "你是油菜区试评价规划师。请将用户的区试评价目标拆解为有序、可独立执行的子任务步骤，"
            "每步应尽量对应一次工具调用（如查询区试数据、计算稳定性、查询适应性区划等）。"
        ),
        replanner_prompt=(
            "你是油菜区试评价重规划师。根据用户原始目标与已完成步骤的结果，判断是否已经掌握足够信息来回答问题。"
            "若已足够，设置 finished=True；否则给出唯一一个最有价值的下一步（finished=False）。不要重复已完成的步骤。"
        ),
        synthesizer_prompt="你是油菜区试评价综合师，请基于各步骤的执行结果，凝练出对原始问题的最终专业回答。",
        executor=executor,
    )

    logger.info(f"用户提问：{question}")

    return planner.make_plan(question).call_executor().get_final_answer()
```

* 实例化 Executor 和 Planner
* 链式调用返回最终回答

### 调用

```python
question = "比较“陕油28”与“华油杂62R”的丰产性、稳定性与生态适应性，给出区域布局建议"
answer = query(question)
print("\n========== 最终回答 ==========\n")
print(answer)
```

执行的日志如下，较长：

```bash
23:15:15 | INFO    | Plan-and-Execute Agent | 用户提问：比较“陕油28”与“华油杂62R”的丰产性、稳定性与生态适应性，给出区域布局建议
23:15:17 | INFO    | httpx | HTTP Request: POST http://127.0.0.1:11434/api/chat "HTTP/1.1 200 OK"
23:15:41 | INFO    | Plan-and-Execute Agent | 规划师生成 6 个步骤：
23:15:41 | INFO    | Plan-and-Execute Agent | 步骤 1 ：获取'陕油28'品种多点位区试基础数据（产量、生育期等）
23:15:41 | INFO    | Plan-and-Execute Agent | 步骤 2 ：获取'华油杂62R'品种多点位区试基础数据（产量、生育期等）
23:15:41 | INFO    | Plan-and-Execute Agent | 步骤 3 ：分析两品种的丰产性：计算平均产量及变异系数，进行显著性检验
23:15:41 | INFO    | Plan-and-Execute Agent | 步骤 4 ：评估稳定性指标：计算回归系数、环境方差贡献率等稳定性参数
23:15:41 | INFO    | Plan-and-Execute Agent | 步骤 5 ：查询油菜生态适应性区划信息，明确不同区域的气候与土壤条件差异
23:15:41 | INFO    | Plan-and-Execute Agent | 步骤 6 ：综合比较结果并给出品种的区域布局建议
23:15:41 | INFO    | Plan-and-Execute Agent | 调用执行器，步骤：获取'陕油28'品种多点位区试基础数据（产量、生育期等）
23:15:44 | INFO    | httpx | HTTP Request: POST http://127.0.0.1:11434/api/chat "HTTP/1.1 200 OK"
23:15:52 | INFO    | httpx | HTTP Request: POST http://127.0.0.1:11434/api/chat "HTTP/1.1 200 OK"
23:16:21 | INFO    | Plan-and-Execute Agent | 步骤（获取'陕油28'品种多点位区试基础数据（产量、生育期等））的执行结果：**步骤结果：**  
已成功获取“陕油28”品种在武汉、长沙、南昌三个区域试验点（2023-2024年）的多年表型记录：  

| 地点 | 年份 | 产量 (kg/亩) | 含油量 (%) | 芥酸 (%) | 硫苷 (μmol/g) | 千粒重 (g) | 株高 (cm) |
|------|------|--------------|------------|----------|----
23:16:24 | INFO    | httpx | HTTP Request: POST http://127.0.0.1:11434/api/chat "HTTP/1.1 200 OK"
23:16:42 | INFO    | Plan-and-Execute Agent | 重规划：下一步 = 获取“华油杂62R"品种在武汉、长沙、南昌三个区域试验点的多年表型数据（产量、生育期等），以便与陕油28进行对比分析
23:16:42 | INFO    | Plan-and-Execute Agent | 调用执行器，步骤：获取“华油杂62R"品种在武汉、长沙、南昌三个区域试验点的多年表型数据（产量、生育期等），以便与陕油28进行对比分析
23:16:45 | INFO    | httpx | HTTP Request: POST http://127.0.0.1:11434/api/chat "HTTP/1.1 200 OK"
23:16:57 | INFO    | httpx | HTTP Request: POST http://127.0.0.1:11434/api/chat "HTTP/1.1 200 OK"
23:17:06 | INFO    | Plan-and-Execute Agent | 步骤（获取“华油杂62R"品种在武汉、长沙、南昌三个区域试验点的多年表型数据（产量、生育期等），以便与陕油28进行对比分析）的执行结果：已获取“华油杂62R"在武汉、长沙、南昌三个区域试验点2023-2024年的表型数据（含产量、含油量等），数据显示该品种在不同区域的平均产量约为198-207 kg/亩，性状表现稳定。下一步将获取“陕油28”的对应数据进行对比分析。
23:17:10 | INFO    | httpx | HTTP Request: POST http://127.0.0.1:11434/api/chat "HTTP/1.1 200 OK"
23:17:37 | INFO    | Plan-and-Execute Agent | 重规划：下一步 = 对陕油28与华油杂62R在三个区域的产量数据进行统计分析（计算平均产量、变异系数等），评估丰产性与稳定性差异，并分析不同区域的环境适应性特征以给出布局建议
23:17:37 | INFO    | Plan-and-Execute Agent | 调用执行器，步骤：对陕油28与华油杂62R在三个区域的产量数据进行统计分析（计算平均产量、变异系数等），评估丰产性与稳定性差异，并分析不同区域的环境适应性特征以给出布局建议
23:17:40 | INFO    | httpx | HTTP Request: POST http://127.0.0.1:11434/api/chat "HTTP/1.1 200 OK"
23:18:05 | INFO    | httpx | HTTP Request: POST http://127.0.0.1:11434/api/chat "HTTP/1.1 200 OK"
23:18:43 | INFO    | Plan-and-Execute Agent | 步骤（对陕油28与华油杂62R在三个区域的产量数据进行统计分析（计算平均产量、变异系数等），评估丰产性与稳定性差异，并分析不同区域的环境适应性特征以给出布局建议）的执行结果：**统计分析结果：**

1.  **丰产性对比**：
    *   “陕油28”多点平均产量为 **206.2 kg/亩**。
    *   “华油杂62R”多点平均产量为 **202.2 kg/亩**。
    *   **结论**：“陕油28”丰产性略高于“华油杂62R”。

2.  **稳定性对比**：
    *   “陕油28”变异系数（CV）为 **0.022**，评价为稳定。

23:18:48 | INFO    | httpx | HTTP Request: POST http://127.0.0.1:11434/api/chat "HTTP/1.1 200 OK"
23:19:05 | INFO    | Plan-and-Execute Agent | 重规划：已结束（无需更多步骤）
23:19:10 | INFO    | httpx | HTTP Request: POST http://127.0.0.1:11434/api/chat "HTTP/1.1 200 OK"

========== 最终回答 ==========

基于多点位区试数据分析，'陕油28'与'华油杂62R'的综合评价及区域布局建议如下：

1. **丰产性对比**：“陕油28”多点平均产量为 206.2 kg/亩，“华油杂62R"为 202.2 kg/亩，前者略高于后者。
2. **稳定性对比**：两者变异系数均较低（分别为 0.022 和 0.024），评价均为稳定；但“陕油28”产量波动更小，稳定性表现更优。
3. **生态适应性对比**：在长江中游及下游地区，两品种适应等级均为“适宜”；但在黄淮地区，“陕油28"为“一般”，而“华油杂62R"评价为“不适宜”。

**区域布局建议：**
*   **长江中下游生态区**：可主推这两个高产品种。
*   **黄淮地区**：建议选择适应性较好的“陕油28"，避免使用“华油杂62R"。
```

根据模型返回的步骤列表，执行步骤时可能会调用工具，并且在适当时会重新规划。

## 优缺点

### 优点

- **全局规划**：先有蓝图再施工，长任务不绕远路。
- **计划可审计**：Plan 是显式中间产物，区试评价报告可存档、评审、追溯。
- **可重规划**：Re-Plan 让计划能适应执行中暴露的新信息(如某品种稳定性不达标需剔除)，不死守原计划。

### 缺点

- **Planner 拆错则全链跑偏**：早期计划质量决定全局，对 Planner(通常用强模型)依赖高。
- **早期计划难应对新信息**：即便有 Re-Plan，频繁重规划也会推高成本与延迟。
- **计划粒度难把握**：拆太粗则执行不精确(每步仍要走长推理)，拆太细则等于退化成 ReAct。
