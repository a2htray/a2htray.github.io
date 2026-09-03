+++
date = '2026-09-02T11:11:02+08:00'
draft = false
title = 'Multi-Agent 架构详解：让 Agent 组成团队分工协作'
categories = ['人工智能', '智能体']
tags = ['AI Agent', 'Multi-Agent', 'LangChain', 'Agent 架构']
toc = true
+++

> 系列文章第 5 篇 · AI Agent 六大主流架构拆解
> 本篇拆解 Multi-Agent 架构，从"一个 Agent 包打天下"走向"一个团队分工协作"，落到**油菜新品种审定申报材料的协同准备**这一典型多角色场景上。

## 什么是 Multi-Agent

核心思想一句话：**多个 Agent 各司其职，由一个 Supervisor 主管负责拆解任务、派发、回收结果、汇总。**

一个最小运行单元：

```
用户任务
  ↓
Supervisor 主管:拆解任务 · 派发 Worker
  ↓
专职 Worker A / B / C（各挂一个领域工具）
  ↓
结果回传 Supervisor
  ↓
Supervisor:决定继续派发 / 收敛汇总
  ↓
汇总输出
```

Multi-Agent 不是简单的多 Agent 组合，关键在于**角色化分工与统一调度**：

- **角色化分工**：每个 Worker 有专属的提示词、工具、记忆，为单一职责优化，效果上限远高于单个通用 Agent。
- **统一调度**：Supervisor 决定“该谁做 / 接下来做什么 / 何时收尾”，让多角色不互相推诿、也不无序循环。
- **过程可观测**：派发记录与中间产物即审计线索，便于复盘。

## 示例背景：油菜新品种审定申报材料准备

油菜新品种要想拿到审定编号，需要向品种审定委员会提交一份完整、规范的申报材料包，典型包括：

- 品种特征特性描述（亲本来源、生育期、株高、抗病性等 DUS 性状档案）。
- DUS 测试报告核查（特异性 Distinctness、一致性 Uniformity、稳定性 Stability 三项是否达标）。
- 区域试验与生产试验总结（一般两年区试 + 一年生试）。
- 品质检测报告（芥酸、硫苷、含油量等是否达双低标准）。
- 栽培技术要点与适宜推广区域。

这是一份**多角色协作才能做好的工作**：

1. **田间调查员**需要从品种档案库拉出 DUS 性状，回答「品种长什么样、怎么长的」。
2. **DUS 审查员**需要核查三项测试是否达标，回答「品种够不够资格」。
3. **材料清单员**需要给出按最新审定规则的申报材料清单与起草要点，回答「该准备哪些材料、怎么写」。

如果用一个全能 Agent做，要么提示词塞得太满、工具冲突严重，要么某个角色做得不够深入。

Multi-Agent 的做法是：**对内分三个专职 Worker，每个只关心自己的领域；对外由 Supervisor 统一调度，按需派发、回收、汇总，最终交付一份完整的申报材料草稿。**

## 架构与运行机制

![Multi-Agent 运行机制](/imgs/learn-agents/multi-agent-mechanism.png)

运行机制三要素：

- **Supervisor 主管**：唯一拥有**调度权**的 Agent，基于当前任务与已完成的派发记录，决定下一步派给谁、按什么指令做、何时结束。
- **专职 Worker**：每个 Worker 只做自己职责范围内的工作，可挂专属工具，可以独立被替换或迭代。
- **派发—回收闭环**：Supervisor 派发任务 → Worker 执行 → 结果回传 Supervisor → Supervisor 根据结果再次决策（继续派发 / 补充 / 收尾），**多轮循环**直至信息足够并汇总。

关键设计点：

- **Supervisor 只调度，不亲自答题**。这样不会和 Worker 抢活，也能让 Worker 持续优化迭代而 Supervisor 不动。
- **派发指令要明确**：Supervisor 给 Worker 的指令必须包含所需参数（如品种名、要查的字段），否则 Worker 容易因信息不足反复追问。
- **结果回传要结构化**：Worker 用 `response_format` 返回 Pydantic 结构化结果（含结论 + 关键字段），Supervisor 与汇总员直接基于字段复用，无需解析自然语言。
- **轮次上限是硬安全网**：Supervisor 自己决定何时结束，但若其陷入循环或漏派，需要 `max_rounds` 上限兜底，强制收尾走汇总。
- **可选汇总员**：多轮派发结束后，可单独由一个「汇总员」角色把派发记录与 Worker 输出合成最终报告，让 Supervisor 保持调度职责单一。

## 基于 LangChain 的代码演示

### 定义三个领域工具

每个 Worker 各挂一个工具，分别负责查询品种特征特性档案、DUS 测试结果、最新审定申报材料清单。

```python
from langchain.tools import tool
from pydantic import BaseModel, Field

# 定义工具

# 模拟品种特征特性档案（DUS 性状）
VARIETY_PROFILE = {
    "秦油2026": {
        "亲本来源": "陕油28 × 华油杂62R",
        "生育期d": 232,
        "株高cm": 168,
        "一次分枝数": 9,
        "单株角果数": 320,
        "每角粒数": 22,
        "抗菌核病": "高抗",
        "抗病毒病": "中抗",
        "抗倒性": "强",
    },
}

# 模拟 DUS 测试结果
DUS_TEST = {
    "秦油2026": {
        "特异性(D)": "通过（与近似品种差异显著）",
        "一致性(U)": "通过（变异系数 2.8%）",
        "稳定性(S)": "通过（两年繁殖性状稳定）",
    },
}

# 模拟最新审定申报材料清单
APPROVAL_RULES = {
    "申报书": "《非主要农作物品种登记申请表》或省级审定申请书",
    "附件清单": [
        "品种选育报告",
        "区域试验与生产试验总结（2 年区试 + 1 年生试）",
        "DUS 测试报告",
        "品质检测报告（双低指标）",
        "抗病性鉴定报告",
        "特征特性描述表",
        "栽培技术要点与适宜区域说明",
        "知识产权与育种者承诺书",
    ],
    "栽培要点": [
        "长江中下游直播 9 月下旬至 10 月中旬,移栽提前 7–10 天。",
        "直播每亩留苗 2.0–2.5 万株。",
        "基肥每亩有机肥 1000–1500 kg + 复合肥 30–40 kg。",
        "蕾薹期与初花期注意补硼防花而不实。",
    ],
}


class VarietyInput(BaseModel):
    variety_name: str = Field(description="品种名称,如 '秦油2026'")


@tool(
    name_or_callable="query_variety_profile",
    description="查询品种特征特性档案（DUS 性状）",
    args_schema=VarietyInput
)
def query_variety_profile(variety_name: str) -> str:
    if variety_name in VARIETY_PROFILE:
        p = VARIETY_PROFILE[variety_name]
        return f"品种 '{variety_name}' 特征特性档案:\n" + "\n".join(
            f"- {k}: {v}" for k, v in p.items()
        )
    return f"未找到品种 '{variety_name}' 的档案。"


@tool(
    name_or_callable="check_dus_test",
    description="核查某品种的 DUS 测试结果（特异性/一致性/稳定性三项）",
    args_schema=VarietyInput
)
def check_dus_test(variety_name: str) -> str:
    if variety_name in DUS_TEST:
        d = DUS_TEST[variety_name]
        lines = [f"品种 '{variety_name}' DUS 测试结果:"]
        for k, v in d.items():
            lines.append(f"- {k}: {v}")
        all_pass = all("通过" in v for v in d.values())
        lines.append(f"\n综合:{'三项全部达标,可申请审定' if all_pass else '存在未达标项,需复核'}")
        return "\n".join(lines)
    return f"未找到品种 '{variety_name}' 的 DUS 测试结果。"


@tool(
    name_or_callable="query_approval_rules",
    description="查询当前品种审定申报材料清单与栽培要点"
)
def query_approval_rules() -> str:
    rules = APPROVAL_RULES
    lines = ["品种审定申报材料清单:", f"- 申报书:{rules['申报书']}", "- 附件清单:"]
    for item in rules["附件清单"]:
        lines.append(f"  · {item}")
    lines.append("- 栽培要点:")
    for item in rules["栽培要点"]:
        lines.append(f"  · {item}")
    return "\n".join(lines)
```

### 定义三个专职 Worker Agent

每个 Worker 只挂一个工具，system_prompt 只关心自己职责；更重要的一点：**每个 Worker 都用 `response_format` 返回结构化结果**，而不是纯文本——这样 Supervisor 与汇总员都能基于字段直接做后续操作，不必再解析自然语言。

```python
from langchain.agents import create_agent
from langchain_ollama import ChatOllama
from langchain.agents.structured_output import ToolStrategy

llm = ChatOllama(model="qwen3.5:9b", temperature=0)


class ProfileResult(BaseModel):
    """田间调查员的结构化产出。"""
    variety_name: str = Field(description="品种名称")
    parents: str = Field(description="亲本来源")
    key_traits: dict = Field(description="关键 DUS 性状,含生育期/株高/分枝数/角果数/每角粒数等")
    disease_resistance: str = Field(description="抗病性结论,含抗菌核病/抗病毒病/抗倒性")
    summary: str = Field(description="一句话特征特性摘要")


class DusResult(BaseModel):
    """DUS 审查员的结构化产出。"""
    distinctness: bool = Field(description="特异性(D)是否通过")
    uniformity: bool = Field(description="一致性(U)是否通过")
    stability: bool = Field(description="稳定性(S)是否通过")
    all_passed: bool = Field(description="三项是否全部达标")
    conclusion: str = Field(description="是否可申请审定的结论")


class DocResult(BaseModel):
    """材料清单员的结构化产出。"""
    application_form: str = Field(description="申报书名称")
    attachments: list = Field(description="附件清单,字符串数组")
    cultivation_points: list = Field(description="栽培要点列表,字符串数组")


profile_worker = create_agent(
    model=llm,
    tools=[query_variety_profile],
    system_prompt=(
        "你是品种特征特性员。接到 Supervisor 的指令后,调用 query_variety_profile 查询品种档案,"
        "把关键 DUS 性状整理成结构化结果返回。"
    ),
    response_format=ToolStrategy(ProfileResult),
)

dus_worker = create_agent(
    model=llm,
    tools=[check_dus_test],
    system_prompt=(
        "你是 DUS 审查员。接到 Supervisor 的指令后,调用 check_dus_test 核查三项 DUS 测试结果,"
        "逐项给出是否通过以及是否可申请审定的结论,以结构化结果返回。"
    ),
    response_format=ToolStrategy(DusResult),
)

doc_worker = create_agent(
    model=llm,
    tools=[query_approval_rules],
    system_prompt=(
        "你是材料清单员。接到 Supervisor 的指令后,调用 query_approval_rules 拉取最新申报材料清单与栽培要点,"
        "按申报书 / 附件清单 / 栽培要点 三段整理成结构化结果返回。"
    ),
    response_format=ToolStrategy(DocResult),
)
```

### 定义 Supervisor Agent 和 Reporter Agent

Supervisor 用 `response_format` + Pydantic 模型做结构化输出，避免手工解析 LLM 输出：

```python
from typing import Literal


class DispatchDecision(BaseModel):
    """Supervisor 的单步调度决策。"""
    next_worker: Literal["田间调查员", "DUS审查员", "材料清单员", "FINISH"] = Field(
        description="下一步交给哪个 Worker;信息已足够汇总时填 FINISH"
    )
    instruction: str = Field(
        description="交给该 Worker 的具体任务描述,必须包含所需参数(如品种名、要看哪些字段)"
    )
    reason: str = Field(description="一句话说明为什么选它 / 为什么可以收尾")


SUPERVISOR_PROMPT = (
    "你是品种审定申报任务的 Supervisor，你手下有三个专职 Worker：\n\n"
    "* 田间调查员(查品种 DUS 性状)\n"
    "* DUS审查员(核查三项 DUS 测试是否达标)\n"
    "* 材料清单员(按最新审定规则给出申报材料清单与栽培要点)\n\n"
    "你的职责只有两件：1) 把总任务拆成派发指令分给合适的 Worker；2) 等信息足够后填 FINISH 结束调度。"
    "不要替 Worker 回答他们的专业问题;不要重复派发已经做完的步骤;"
    "各 Worker 回传给你的都是结构化结果(含结论与关键字段),请基于这些字段判断信息是否足够;"
    "当三个 Worker 都把结果回传且足够支撑一份完整申报材料时,必须填 FINISH。"
)

supervisor = create_agent(
    model=llm,
    tools=[],
    system_prompt=SUPERVISOR_PROMPT,
    response_format=ToolStrategy(DispatchDecision),
)

reporter = create_agent(
    model=llm,
    tools=[],
    system_prompt=(
        "你是申报材料汇总员。根据 Supervisor 给你的「总任务」与派发记录,把每个 Worker 的产出整理成"
        "一份结构化的审定申报材料草稿,分『特征特性摘要』『DUS 核查结论』『申报材料清单与栽培要点』三部分输出。"
    ),
)


WORKERS = {
    "田间调查员": profile_worker,
    "DUS审查员": dus_worker,
    "材料清单员": doc_worker,
}
```

### 高度-执行-汇总

```python
from dataclasses import dataclass


@dataclass
class WorkerRun:
    agent_name: str
    instruction: str
    output: str


def multi_agent_run(question, max_rounds=5) -> str:
    runs: list[WorkerRun] = []
    for i in range(1, max_rounds + 1):
        prompt = (
            f"用户提问：{question}\n\n"
            "已派发任务与结果：\n\n"
            + ("\n".join([f"* {run.agent_name}: {run.instruction} -> {run.output}" for run in runs]) or "（暂无）")
            + "\n\n"
            "请基于上述结构化产出决定是否还需派发，或填 FINISH 结束调度。"
        )
        resp = supervisor.invoke({"messages": [{"role": "user", "content": prompt}]})
        dispatch_decision: DispatchDecision = resp["structured_response"]

        if dispatch_decision.next_worker == "FINISH":
            break
        if dispatch_decision.next_worker not in WORKERS:
            break

        worker = WORKERS[dispatch_decision.next_worker]
        worker_prompt = (
            f"用户提问：{question}\n\n"
            f"Supervisor 给你的任务：{dispatch_decision.instruction}"
        )
        worker_resp = worker.invoke(
            {"messages": [{"role": "user", "content": worker_prompt}]}
        )
        result = worker_resp["structured_response"]
        runs.append(
            WorkerRun(
                agent_name=dispatch_decision.next_worker,
                instruction=dispatch_decision.instruction,
                output=result.model_dump_json(),
            )
        )

        print((
            f"====== 第 {i} 轮结果 ======\n"
            f"Worker：{dispatch_decision.next_worker}\n"
            f"任务：{dispatch_decision.instruction}\n"
            f"理由：{dispatch_decision.reason}\n"
            f"产出：{result.model_dump_json()}"
        ))

    summary_prompt = (
        f"用户提问：{question}\n\n"
        "各 Worker 的产出如下：\n\n"
        + ("\n".join([f"* {run.agent_name}: {run.instruction} -> {run.output}" for run in runs]) or "（未派发任何 Worker）")
    )
    final = reporter.invoke({"messages": [{"role": "user", "content": summary_prompt}]})
    return final["messages"][-1].content
```

### 运行

```python
if __name__ == "__main__":
    question = (
        "为油菜新品种 '秦油2026' 准备一份品种审定申报材料草稿:"
        "核对特征特性档案、确认 DUS 测试是否达标、"
        "理清申报材料清单并给出栽培要点。"
    )
    print(multi_agent_run(question))
```

执行输出：

```text
====== 第 1 轮结果 ======
Worker：田间调查员
任务：请核查油菜新品种'秦油2026'的特征特性档案，包括品种名称、特征特性描述、DUS测试数据等关键字段
理由：首先需要核对特征特性档案，这是申报材料的基础信息
产出：{"variety_name":"秦油2026","parents":"陕油28 × 华油杂62R","key_traits":{"一次分枝数":"9个","单株角果数":"320个","株高":"168cm","每角粒数":"22粒","生育期":"232天"},"disease_resistance":"抗菌核病: 高抗; 抗病毒病: 中抗; 抗倒性: 强","summary":"秦油2026为陕油28与华油杂62R杂交选育的油菜新品种，生育期232天，株高168cm，一次分枝9个，单株角果320个，每角粒22粒，抗菌核病高抗、抗病毒病中抗、抗倒性强"}
====== 第 2 轮结果 ======
Worker：DUS审查员
任务：请核查油菜新品种'秦油2026'的三项DUS测试（特异性、一致性、稳定性）是否达标，需包含品种名称、DUS测试结果、达标结论等关键字段
理由：田间调查员已完成特征特性档案核查，接下来需要确认DUS测试是否达标
产出：{"distinctness":true,"uniformity":true,"stability":true,"all_passed":true,"conclusion":"可申请审定"}
====== 第 3 轮结果 ======
Worker：材料清单员
任务：请为油菜新品种'秦油2026'整理申报材料清单并给出栽培要点，需包含品种名称、申报材料清单（按最新审定规则）、栽培要点等关键字段
理由：田间调查员和DUS审查员已完成工作，但材料清单员尚未参与，需要获取申报材料清单与栽培要点才能完成完整申报材料准备
产出：{"application_form":"《非主要农作物品种登记申请表》或省级审定申请书","attachments":["品种选育报告","区域试验与生产试验总结（2年区试+1年生试）","DUS测试报告","品质检测报告（双低指标）","抗病性鉴定报告","特征特性描述表","栽培技术要点与适宜区域说明","知识产权与育种者承诺书"],"cultivation_points":["长江中下游直播9月下旬至10月中旬,移栽提前7–10天。","直播每亩留苗2.0–2.5万株。","基肥每亩有机肥1000–1500 kg + 复合肥30–40 kg。","蕾薹期与初花期注意补硼防花而不实。"]}
# 油菜新品种 '秦油 2026' 品种审定申报材料草稿

## 一、特征特性摘要

*   **品种名称**：秦油 2026
*   **亲本组合**：陕油 28 × 华油杂 62R
*   **主要农艺性状**：
    *   **生育期**：232 天（中熟）
    *   **株高**：168cm
    *   **一次分枝数**：9 个
    *   **单株角果数**：320 个
    *   **每角粒数**：22 粒
*   **抗病性表现**：
    *   **菌核病**：高抗
    *   **病毒病**：中抗
    *   **抗倒性**：强
*   **品种简介**：秦油 2026 为陕油 28 与华油杂 62R 杂交选育的油菜新品种，具有生育期适中、株型紧凑、抗病性强等特点。

## 二、DUS 核查结论

*   **特异性（Distinctness）**：达标 (True)
*   **一致性（Uniformity）**：达标 (True)
*   **稳定性（Stability）**：达标 (True)
*   **综合判定**：三项测试全部通过，符合品种审定要求。
*   **审查结论**：**可申请审定**。

## 三、申报材料清单与栽培要点

### 1. 申报材料清单（按最新审定规则）
*   **主表**：《非主要农作物品种登记申请表》或省级审定申请书
*   **附件材料**：
    1.  品种选育报告
    2.  区域试验与生产试验总结（含 2 年区试 +1 年生试）
    3.  DUS 测试报告
    4.  品质检测报告（双低指标）
    5.  抗病性鉴定报告
    6.  特征特性描述表
    7.  栽培技术要点与适宜区域说明
    8.  知识产权与育种者承诺书

### 2. 栽培技术要点
*   **播种时间**：长江中下游地区直播宜在 9 月下旬至 10 月中旬；若采用移栽方式，需提前 7–10 天。
*   **种植密度**：直播每亩留苗 2.0–2.5 万株。
*   **施肥管理**：基肥每亩施用有机肥 1000–1500 kg + 复合肥 30–40 kg。
*   **关键农事**：蕾薹期与初花期注意补硼，以防花而不实。
```

所有的 Agent 都是结构化的返回，Supervisor 根据提问和上下文，依次调用了 Worker，并最终将结果汇总由 Reporter Agent 生成申报材料草稿。

## 优缺点

### 优点

- **角色化分工**：每个 Worker 按单一职责定制提示词与工具，深度更足。
- **可并行与复用**：多个 Worker 可并行推进，成熟 Worker 还能跨任务复用。
- **能力可扩展**：新增能力只需加一个 Worker，Supervisor 调度逻辑不变。
- **过程可观测**：派发记录与中间产物即审计线索，便于定位问题环节。

### 缺点

- **协调成本高**：每轮派发都要过 Supervisor，token 与延迟成倍上升。
- **上下文共享难**：Worker 之间默认隔离，结果冲突与信息孤岛需额外机制。
- **主管是单点**：Supervisor 拆错或漏派，整个团队跟着空转。
- **收敛难保证**：多轮循环可能来回打转，需要轮次上限与终止条件。