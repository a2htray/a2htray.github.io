+++
date = '2026-08-18T09:53:32+08:00'
draft = false
title = 'Airflow 是什么，以及为什么你大概需要它'
categories = ['后端技术', 'Airflow']
tags = ['Airflow', '批处理', '定时任务']
toc = true
+++

> 这是 Airflow 学习系列的第 1 篇。本篇不写代码，先讲清楚一件事：Airflow 到底解决什么问题，它的核心抽象是什么，以及它适合什么样的场景。

## 凌晨三点的告警

先讲一个场景。

你负责一套数据管道，每天凌晨两点从业务库抽数到数仓，做几张宽表，跑几个聚合，最后刷新 BI 报表。整个链条七八个步骤，有的串行，有的并行。你用 cron 加一堆 shell 脚本串起来，每个脚本之间靠文件标记或者上一条命令的退出码做衔接。

某个凌晨三点，你被告警吵醒。某一步失败了。

问题立刻堆上来：是哪一步失败的？失败前跑了什么？现在数据停在哪个中间状态？要不要重跑，重跑会不会重复写入？下游报表显示的，是今天的旧数据还是错误数据？

你爬起来，登服务器，翻日志，发现是第三步的某个 SQL 超时了。你手动补跑第三步，再接着跑第四五六七步。处理完天都快亮了。

这套流程的问题不在"它会挂"——任何系统都会挂——而在"**它挂了之后你很难干净地恢复**"。没有依赖管理，没有状态记录，没有统一的重跑入口，没有全局视角的可观测性。你在黑暗里摸黑接线。

Airflow 就是来处理这件事的，提供：

1. 管理依赖
2. 记录状态
3. 统一入口
4. 可观测性

## 数据编排到底难在哪

把"按时跑一系列任务"这件事拆开看，难点其实集中在四个地方。

1. **依赖管理**。任务之间有先后，有的还能并行。用脚本的话你得自己写 wait 逻辑、自己判断前置是否成功。任务一多，拓扑就乱。一个典型的 ETL 是这样的：取数有三路并行，取完汇总，汇总完分两路做特征和报表，特征做完再训模型。这种拓扑你手写 if/else 维护，迟早会出错。
2. **失败与重跑**。任务可能因为网络、数据质量、资源任何原因挂掉。挂了之后要能重跑，而且重跑不能搞坏数据——比如已经写入了一半的表，直接重跑会重复或冲突。理想情况是任务幂等，且编排系统能从失败点继续，而不是从头来。
3. **可观测性**。当前跑到哪一步了？每步花了多久？哪步慢了？失败原因是什么？历史趋势怎样？这些信息散在各个脚本日志里，基本等于没有。出了问题你只能挨个去翻。
4. **回填**。今天修了个 bug，想把它应用到过去三十天的历史数据上。你需要按日期逐天重跑那段管道。cron 做不到这件事，你得自己写循环、自己管日期边界、自己处理哪天已经跑过哪天没有。

这四个点，就是工作流编排器（workflow orchestrator）要解决的核心问题。Airflow 是其中最知名的一个。

## Airflow 是什么

Apache Airflow 是一个开源的工作流编排平台，用 Python 写的，2014 年由 Airbnb 内部开发，2016 年进入 Apache 孵化，2019 年成为顶级项目。到现在它几乎是数据工程领域的事实标准之一。

用一句话说：你用 Python 代码定义一批任务和它们的依赖关系，Airflow 负责按时调度、分配资源执行、记录状态、提供界面让你看和操作。

关键点在于"用 Python 代码定义"。Airflow 的核心抽象是 DAG（有向无环图），一个 DAG 就是一个 Python 文件，里面声明了若干 Task 以及它们之间的依赖。比如：

```python
from airflow import DAG
from airflow.operators.bash import BashOperator
from datetime import datetime

with DAG(
    dag_id="daily_etl",
    schedule="0 2 * * *",   # 每天凌晨两点
    start_date=datetime(2024, 1, 1),
    catchup=False,
) as dag:
    extract   = BashOperator(task_id="extract",   bash_command="sh extract.sh")
    transform = BashOperator(task_id="transform", bash_command="sh transform.sh")
    load      = BashOperator(task_id="load",      bash_command="sh load.sh")

    extract >> transform >> load
```

因为是代码，所以可以版本管理、可以 review、可以测试、可以复用、可以动态生成。这一套叫 **DAG as code**，是 Airflow 区别于很多老式调度工具的根本。后面你会看到，几乎所有 Airflow 的设计取舍，都从这个前提推出来。

## 核心概念，一次性讲清

Airflow 的术语不少，初学容易被绕晕。这里把最核心的几个一次讲清，后面文章会反复用到。

**DAG**。有向无环图，一组带方向的依赖关系，不能成环。在 Airflow 里，一个 DAG 就是一个工作流的骨架——它定义了"有哪些任务、谁先谁后"，但不定义"任务具体干什么"。

**Task**。DAG 里的一个执行单元，是 DAG 在某次运行中实例化出来的节点。一个 DAG 每天跑一次，每天就产生一组 Task 实例。

**Operator**。Task 的类型模板，决定了这个任务怎么跑。常用的有：

- `BashOperator`：跑一条 shell 命令
- `PythonOperator`：跑一个 Python 函数
- `EmailOperator`：发邮件
- `HttpOperator`：调一个 HTTP 接口
- `SqlOperator`：执行 SQL
- 各种特定系统的 Operator，比如 `SparkSubmitOperator`、`KubernetesPodOperator`

Operator 是 Airflow 生态最丰富的地方，几乎所有常见外部系统都有对应 Operator，这是它最大的护城河之一。

**Sensor**。一种特殊的 Operator，它不"做事"，而是"等"——等某个条件成立再放行。比如等一个文件出现在 S3、等一个数据库表有新分区、等一个外部任务完成。**Sensor 让管道能基于事件推进，而不只是基于时间**。

**Hook**。与外部系统交互的底层封装，一个 Hook 对应一种连接（`S3Hook`、`PostgresHook`、`SlackHook`）。Operator 内部通常用 Hook 来干活。你也可以在自己的 `PythonOperator` 里直接用 Hook，跳过 Operator 这一层。

**XCom**。cross-communication 的缩写，task 之间传递小量数据的机制。一个 task 把结果 push 出去，另一个 task pull 进来。注意 XCom 设计上是给小数据用的（默认序列化后存进 Metadata DB），大数据应该走外部存储（S3、数据库表），XCom 只传引用——比如文件路径、表名。

**Scheduler**。调度器，整个系统的大脑。它持续扫描所有 DAG，决定哪个 DAG 该在什么时候触发、哪个 Task 的上游跑完了可以入队。Scheduler 是 Airflow 的核心循环。

**Executor**。执行器，决定 Task 实际在哪里跑。不同 Executor 对应不同规模：

- `SequentialExecutor`：单进程串行，只适合测试
- `LocalExecutor`：本机多进程
- `CeleryExecutor`：通过 Celery 分布式到多台 worker，生产标配
- `KubernetesExecutor`：每个 Task 起一个 K8s Pod，弹性最好

**Metadata DB**。元数据库（通常 PostgreSQL 或 MySQL），存所有状态——DAG 定义、运行历史、Task 实例状态、XCom、变量、连接配置。这是 Airflow 的"真相之源"，挂了等于全挂。

**Webserver**。Web 界面，让你看 DAG 列表、运行历史、Graph 视图、日志，可以手动触发和重跑。它是给人用的窗口，本身不参与调度。

![](/imgs/learn-airflow/airflow-concept-map.png)

这几个概念的关系一句话概括：你写 DAG（骨架）→ DAG 里是 Operator（任务类型）→ Operator 用 Hook 连外部系统 → Scheduler 决定何时跑 → Executor 决定在哪跑 → 状态写进 Metadata DB → 你在 Webserver 上看。

## 设计哲学

理解 Airflow 的几个核心取舍，能帮你判断它适不适合你的场景。

**DAG as code，而不是 DAG as config**。很多调度系统让你填表单、拖拽、配 YAML 来定义工作流。Airflow 坚持用完整的 Python。代价是学习曲线陡一点，收益是你可以用代码做任何事——循环生成任务、按配置动态拼接、复用函数、写辅助逻辑。复杂工作流用配置几乎一定会撞墙，用代码不会。

**显式依赖，调度不管数据流**。Airflow 管的是"任务执行顺序和触发时机"，不管数据怎么流。数据要么通过 XCom 传小量，要么走外部存储由 task 自己读写。这是 Airflow 和 Dagster 的根本分歧——Dagster 把"数据资产"作为一等公民，Airflow 把"任务"作为一等公民。两种思路各有适用场景，没有绝对优劣。

**幂等与可重跑是一等公民**。Airflow 假设任务可能失败、可能被重跑很多次。所以它鼓励你把*任务写成幂等的*，并提供了清晰的重跑入口（从某次运行、某个任务开始重跑）。这个假设深刻影响了它的状态模型设计——一个 Task 实例不是简单的"成功/失败"，而是有 queued、running、success、failed、up_for_retry、upstream_failed、skipped 一整套状态机。

**拉取而非推送**。Airflow 主动去扫描、决定、触发，而不是被动等外部系统通知它（虽然 Sensor 能模拟等待，但整体是 pull 模型）。这意味着它特别擅长"按时间表做事"，对"事件即时触发"的支持相对弱。这也解释了为什么它叫批处理编排器，不叫事件引擎。

## 跟其他工具比

编排器赛道不小，几个常被拿来对比的：

**Luigi**。Spotify 2012 年开源，Airflow 的精神前辈。也是 Python DAG，思路很像。但 Luigi 基本停止积极演进，社区和新功能都在 Airflow 这边。新项目没必要选 Luigi。

**Prefect**。2018 年出现，定位更现代。强调动态工作流（运行时决定跑什么）、Pythonic API、原生云。代码写起来比 Airflow 顺手，本地开发体验好。但生态深度、企业部署成熟度、Operator 数量还不及 Airflow。

**Dagster**。把"数据资产"（asset）作为核心抽象，强调数据血缘、软件工程化的资产定义和测试。如果你的关注点是"我有哪些数据资产、它们怎么衍生、怎么测试"，Dagster 的模型比 Airflow 顺。但对纯任务编排场景，它的概念更重，学习成本不低。

**Mage**。比较新的工具，主打可视化 + 代码混合，AI 辅助生成。上手快，但生态和成熟度还在早期，复杂场景的工程化能力有待验证。

一句话总结对比：Airflow 是"成熟、生态厚、配置略繁"的稳妥选择；Prefect/Dagster 是"更现代、更软件工程化"的新选择；Mage 是"快速上手"的尝试。新项目如果团队已经在重度用 dbt、看重数据资产视角，可以认真考虑 Dagster；其余大多数批处理编排场景，Airflow 仍然是默认选项。

## 什么时候用 Airflow，什么时候别用

**适合的场景**：

- 批处理 ETL/ELT 管道，每天或每小时跑一次
- 定时报表刷新、数据同步
- 机器学习训练管道：取数 → 特征 → 训练 → 评估 → 上线
- 多步骤的运维任务编排：备份、清理、对账
- 需要回填历史数据的场景

**不适合的场景**：

- 实时流处理。Airflow 是批处理的，毫秒级延迟的流式要用 Flink、Kafka Streams 这类
- 低延迟事件驱动。一个事件来了立刻触发一个动作，这是消息队列 + 函数的活
- 极简单的单步定时任务。一个 cron 表达式能搞定的事，上 Airflow 是过度设计
- 对数据资产血缘有强需求。这种场景 Dagster 更对路

判断标准很简单：如果你有"多步骤、有依赖、需要可观测和可重跑"的批处理管道，Airflow 基本是对的；如果是"事件来了立刻响应"或"一条命令搞定"，别用 Airflow。
