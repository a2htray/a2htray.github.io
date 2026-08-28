+++
date = '2026-08-28T14:43:23+08:00'
draft = false
title = 'Reflexion 架构详解：让 Agent 学会自我纠错'
categories = ['人工智能', '智能体']
tags = ['AI Agent', 'Reflexion', 'LangChain', 'Agent 架构']
toc = true
+++

> 系列文章第 3 篇 · AI Agent 六大主流架构拆解
> 本篇拆解 Reflexion，把“执行 → 评估 → 反思 → 记忆”的自我纠错循环装进 Agent，并把它落到**油菜杂交组合配置方案的迭代优化**这一典型“试错改进”任务上。

## 什么是 Reflexion

核心思想一句话：**在执行循环外面再套一层 —— Actor 产出尝试，Evaluator 给评分与反馈，Self-Reflector 根据反馈生成反思经验存入 Memory；下一轮重试时把累积的经验注入上下文，逐轮逼近目标。**

> Note：反思经验。

一个最小运行单元：

```
任务
  → Actor 执行器: 产出尝试结果(可调用工具)
  → Evaluator 评估: 评分 + 反馈
  → Self-Reflector 反思: 据反馈生成经验教训
  → Memory 长期记忆: 存储反思经验
  →(未达标)经验注入下一轮 → Actor 重试
  → 达标 → 最终输出
```

与 ReAct / Plan-and-Execute 的关键区别：

- ReAct 解决“怎么做”，Plan-and-Execute 解决“怎么规划长任务”，两者都没有纠偏机制，执行错了一步就一路错下去，没有改进的退路。
- Reflexion 补齐退路：把“做完复盘-下次改进”的认知过程工程化，用**自然语言反思**作为自我反馈信号，**跨轮**累积经验。

一句话总结：ReAct 是“思考驱动执行”，Plan-and-Execute 是“规划驱动执行”，**Reflexion 是“反思驱动改进”**。

## 示例背景：油菜杂交组合配置

油菜杂交组合配置是育种核心决策之一——选两个亲本组成杂交组合，期望后代综合目标性状(抗菌核病、亩产>200kg、含油量>45%、油酸>75%等)。这是一项典型的**单次配置往往不理想、需多轮试错改进**的任务：

1. **一次配置很难同时满足所有目标**：油酸高的亲本往往抗病或产量不够，反之亦然，需要权衡与重选。
2. **每次失败都有可学习的经验**："上次选的两个亲本油酸都偏低"——这种经验下一轮可避免。
3. **经验可跨轮累积**：多轮反思形成的经验链(Actor 每次重试都带着前几轮教训)，让方案逐步逼近达标。
4. **可审计**：每轮的尝试、评分、反思都留痕，育种方案可追溯。

相比之下，ReAct 适合“即席查询”、Plan-and-Execute 适合“多阶段长任务”；“单次不达预期、需多轮试错改进”的场景，Reflexion 明显更合适。工具上也跟前两篇不重叠——本篇两个工具是 `evaluate_combination`(组合评估)/ `query_breeding_target`(育种目标查询)，与种质查询、区试数据零重叠。

---

## 架构与运行机制

![Reflexion 运行机制](/imgs/learn-agents/reflexion-mechanism.png)

运行机制四要素：

- **Actor 执行器**：接收任务与历史经验，产出尝试结果(可调用工具)；用一个 `create_agent` 子 Agent，内含工具调用循环。
- **Evaluator 评估器**：对尝试打分(0–100)并给出具体反馈(哪里没达标、怎么改)；用一个 `create_agent` Agent(不挂工具)，prompt 里写清目标与评分标准，返回 JSON 评分。
- **Self-Reflector 反思器**：用一个 `create_agent` Agent(不挂工具)，据 Evaluator 反馈生成一条简洁的反思经验(“下一轮该怎么调整”)，存入 Memory。
- **Memory 长期记忆**：跨轮累积反思经验，下一轮 Actor 重试时把经验注入上下文；本例用 JSON 文件显式持久化(见下文代码)。

关键设计点：

- **Evaluator 必须可靠**：它给 Self-Reflector 喂反馈，反馈不准则反思跑偏、整条链都空转。生产中常用“具体可验证的指标”而非纯主观打分。
- **反思经验要简洁可执行**：一条经验一句话，聚焦“下一轮该做什么/不做什么”，否则会撑爆上下文。
- **Memory 可跨任务复用**：同一类任务的反思经验(例如“油菜杂交优先看油酸 + 抗病双达标”)可在多次任务间共享，不只是单次会话的轮内累积。
- **多轮成本是主要代价**：每轮都要跑一次 Evaluator + Reflector，token 与延迟线性增长；需设 `max_rounds` 兜底。

## 基于 LangChain 的代码演示

### 定义杂交组合配置工具

用 `@tool` + Pydantic `args_schema` 定义两个工具：种质资源查询和杂交组合计算得分。

```python
import math

from langchain.tools import tool
from pydantic import BaseModel, Field

PARENTS = {
    "陕油28":    {"抗病": "高抗", "产量kg/亩": 210, "含油量%": 45.0, "油酸%": 68},
    "华油杂62R":  {"抗病": "抗",   "产量kg/亩": 208, "含油量%": 45.6, "油酸%": 70},
    "高油酸1号":  {"抗病": "中抗", "产量kg/亩": 195, "含油量%": 47.2, "油酸%": 82},
    "抗病9号":    {"抗病": "高抗", "产量kg/亩": 205, "含油量%": 44.5, "油酸%": 66},
}


@tool(
    "get_germplasms",
    description="获取所有油菜种质名称",
)
def get_germplasms() -> str:
    """获取所有种质名称"""
    return "、".join(PARENTS.keys())


def target_minmax(target: str) -> tuple[float, float]:
    """计算目标性状的最小值和最大值"""
    values = [PARENTS[key][target] for key in PARENTS]
    return min(values), max(values)


def minmax_normalize100(value: float, min_value: float, max_value: float) -> float:
    """将值归一化到 0-100 之间"""
    return (value - min_value) / (max_value - min_value) * 100


class ParentPair(BaseModel):
    parent_a: str = Field(description="第一个亲本")
    parent_b: str = Field(description="第二个亲本")


@tool(
    "evaluate_combination",
    description="评估杂交组合性状，按目标性状（抗病、产量、含油量、油酸）进行逐项打分（0-100）。",
    args_schema=ParentPair,
)
def evaluate_combination(parent_a: str, parent_b: str) -> str:
    if parent_a not in PARENTS or parent_b not in PARENTS:
        return f"未找到亲本，可选亲本有{'、'.join(PARENTS.keys())}"
    a, b = PARENTS[parent_a], PARENTS[parent_b]

    anti_scores = {"抗": 75, "中抗": 85, "高抗": 95}
    anti_score = math.sqrt(anti_scores[a["抗病"]] * anti_scores[b["抗病"]])

    yield_minmax = target_minmax("产量kg/亩")
    yield_score = math.sqrt(minmax_normalize100(a["产量kg/亩"], *yield_minmax) * minmax_normalize100(b["产量kg/亩"], *yield_minmax))

    oil_minmax = target_minmax("含油量%")
    oil_score = math.sqrt(minmax_normalize100(a["含油量%"], *oil_minmax) * minmax_normalize100(b["含油量%"], *oil_minmax))

    oil_content_minmax = target_minmax("油酸%")
    oil_content_score = math.sqrt(minmax_normalize100(a["油酸%"], *oil_content_minmax) * minmax_normalize100(b["油酸%"], *oil_content_minmax))

    return f"""{parent_a} x {parent_b} 组合各性状得分（0-100）：
* 抗病：{anti_score:.2f}
* 产量：{yield_score:.2f}
* 含油量：{oil_score:.2f}
* 油酸：{oil_content_score:.2f}"""

tools = [get_germplasms,evaluate_combination]
```

### 构建 Actor / Evaluator / Reflector 三个 Agent

```python
from langchain.agents import create_agent
from langchain.messages import SystemMessage, HumanMessage
from langchain_ollama import ChatOllama

llm = ChatOllama(model="qwen3.5:9b", temperature=0.2)

actor = create_agent(
    llm,
    tools=tools,
    system_prompt=SystemMessage(
        "你是一个专业的油菜杂交组合模拟助手，选择两个亲本组成一个杂交组合，能够计算杂交组合的各性状得分。"
    ),
)


class EvaluationResult(BaseModel):
    score: float = Field(description="组合的综合评分（0-100）")
    suggestion: str = Field(description="改进建议")
    finish: bool = Field(description="是否完成")


evaluator = create_agent(
    llm,
    tools=[],
    system_prompt=SystemMessage(
        "你是一个专业的油菜杂交组合评估助手，能够基于杂交组合的性状得分，给出组合的综合评分（0-100）和改进建议， 若评估结果符合目标性状，则返回 finish=True，否则返回 finish=False。"
    ),
    response_format=EvaluationResult,
)


class ReflectionResult(BaseModel):
    reflection: str = Field(description="上一轮的尝试和反馈的简洁经验")


reflector = create_agent(
    llm,
    tools=[],
    system_prompt=SystemMessage(
        "你是一个专业的油菜杂交组合反思器，根据上一轮的尝试和反馈，生成一条简洁的反思经验，聚集于下一轮的调整。"
    ),
    response_format=ReflectionResult,
)
```

* actor Agent：根据用户问题，选择两个亲本组成一个杂交组合，能够计算杂交组合的各性状得分。
* evaluator Agent：根据杂交组合的性状得分，给出组合的综合评分（0-100）和改进建议。
* reflector Agent：根据上一轮的尝试和反馈，生成一条简洁的反思经验，聚集于下一轮的调整。

### 迭代 Agent

```python
question = "从油菜种质资源库中选择两个亲本，组配出一个高产和高油酸的杂交组合。"

memory: list[str] = []


def query(question, max_rounds: int = 3) -> str:
    final_content = ""
    for i in range(max_rounds):
        print(f"========= 第{i+1}轮的经验数：{len(memory)} ==========")
        print(f"""历史经验：
{'\n'.join(memory)}""")

        actor_messages = [HumanMessage(content=question)] + [SystemMessage(content=f"历史经验{i+1}: {m}") for i, m in enumerate(memory, start=1)]
        actor_content = actor.invoke({
            "messages": actor_messages,
        })["messages"][-1].content

        final_content = actor_content

        evaluator_resp = evaluator.invoke({
            "messages": [HumanMessage(content=question), SystemMessage(content=actor_content)],
        })
        evaluator_result: EvaluationResult = evaluator_resp["structured_response"]

        if evaluator_result.finish:
            break
        else:
            reflector_resp = reflector.invoke({
                "messages": [HumanMessage(content=actor_content), SystemMessage(content=evaluator_result.suggestion)],
            })
            reflector_result: ReflectionResult = reflector_resp["structured_response"]
            memory.append(reflector_result.reflection)
    return final_content


if __name__ == "__main__":
    print(query(question))
```

* memory 列表：存储每一轮的反思经验，调用 actor Agent 时作为系统提示词传入。
* 根据 evaluator Agent 返回的 finish 字段，判断是否完成任务。
* 若未完成任务，调用 reflector Agent 生成反思经验，添加到 memory 列表中。

执行结果如下：

```text
========= 第1轮的经验数：0 ==========
历史经验：

========= 第2轮的经验数：1 ==========
历史经验：
上一轮筛选显示单一亲本难以同时满足高产与高油酸双重目标，需转向多性状协同育种策略：1）优先选择产量≥80分且含油量≥35分的亲本进行杂交；2）引入抗病性强的背景材料作为基础，通过回交或分子标记辅助稳定优良基因型；3）扩大亲本库筛选范围，探索不同遗传背景的组配潜力。

1. **陕油28 × 高油酸1号**  
   - 抗病性得分：**89.86分**（表现优异）  
   - 产量得分：**0.00分**（未检测到显著增产潜力）  
   - 含油量得分：**43.03分**（中等水平）  
   - 油酸得分：**35.36分**（符合高油酸目标，但仍有提升空间）

2. **华油杂62R × 高油酸1号**  
   - 抗病性得分：**79.84分**（表现良好）  
   - 产量得分：**0.00分**（同样未检测到显著增产潜力）  
   - 含油量得分：**63.83分**（显著提升，接近高产目标）  
   - 油酸得分：**50.00分**（达到高油酸标准且表现更优）

### 结论与建议：
- **华油杂62R × 高油酸1号**组合在含油量（63.83分）和油酸含量（50.00分）上均优于另一组配，更接近“高产和高油酸”的双重目标。  
- 尽管两个组合的产量得分均为0.00分，但结合历史经验，这可能是因为当前亲本库中缺乏高产品种或需进一步优化育种策略。建议后续引入更多高产型亲本（如通过分子标记辅助选择）进行杂交试验，以突破产量瓶颈。

如需进一步分析其他组配潜力，可继续探索不同遗传背景的亲本组合！
```

* 第 1 次反思，给出了反思经验，提示后续筛选高产量且含油量高的亲本。
* 在第 2 轮时，evaluator Agent 认为结果符合，并返回 finish=True。
* 得出的杂交组合是：华油杂62R × 高油酸1号。

## 优缺点

### 优点

- **自我纠错**：让 Agent 在“试错”中逐步改进，显著提升高准确度任务的成功率。
- **经验复用**：Memory 里的反思经验可跨轮、跨任务累积，同一类任务越来越快收敛。

### 缺点

- **依赖基线 Actor 能力**：若 Actor Agent 本身就弱，反思再多次也难救。
- **Evaluator 可靠性依赖**：依赖 Evaluator Agent 的可靠性，若 E Evaluator 给不出有信息量的反馈，整个调用链条将“走偏”。
- **多轮成本是主要代价**：每轮 Actor + Evaluator + Reflector 三次 LLM 调用，token 与延迟线性增长；`max_rounds` 兜底但不能根治。
- **反思可能误导**：Reflector Agent 可能“想错”，把一次偶然失败归因为错误原因，导致后续轮次走偏。需要 Evaluator Agent 反馈“具体可验证”来约束。

## 小结

Reflexion 解决了“怎么越做越好”，给 ReAct / Plan-and-Execute 补上“反思-改进”的闭环。它让油菜杂交组合配置这种“单次难达标、需多轮试错”的任务有了可承载的架构。