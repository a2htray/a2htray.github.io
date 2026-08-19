+++
date = '2026-08-19T09:58:18+08:00'
draft = false
title = 'Airflow 是怎么把 DAG 跑起来的'
categories = ['后端技术', 'Airflow']
tags = ['Airflow', '批处理', '定时任务', 'DAG']
toc = true
+++

> 这是 Airflow 学习系列的第 2 篇。上一篇讲了 Airflow 是什么、为什么需要它、核心概念词汇表。这一篇钻进内部，把"齿轮怎么转"讲透——你看完之后，第 3 篇上手写 DAG 时，每一步操作都不会是黑盒。

## Airflow 组件

第 1 篇列了五大组件：Scheduler、Executor、Metadata DB、Webserver，外加你写的 DAG。那次是定义，这次讲它们怎么咬合。

一句话总览：你写好 DAG 文件放进目录，**Scheduler** 发现它、决定何时触发、决定哪个 task 可以跑，把 task 交给 **Executor** 去实际执行，**Executor** 跑完把状态写回 **Metadata DB**，**Webserver** 读 Metadata DB 把这一切展示给你看。

注意一个关键点：**Metadata DB 是所有组件的共享内存**。Scheduler、Executor、Webserver 之间不直接通信，它们都通过读写同一个数据库来协作。这也是为什么生产环境必须用 PostgreSQL 或 MySQL，不能用 sqlite——并发写入压力全压在这一个库上。

![](.//imgs/learn-airflow/airflow_5_core_components.png)

这种"以数据库为总线"的架构是 Airflow 的根本设计，理解了它，后面很多行为就好懂了。

## DAG 发现：DagBag

你把 DAG 写成 .py 文件，丢进 Airflow 的 DAG 文件夹（默认 `~/airflow/dags`）。Scheduler 不会立刻就知道它来了。

Scheduler 内部有一个叫 **DagBag** 的东西，它会周期性地扫描 DAG 文件夹，把每个 .py 文件当作 Python 模块导入执行，提取里面定义的 DAG 对象。这个过程叫"收集"（collect）。

这里有一个 Airflow 最重要的认知，也是无数初学者踩坑的地方：**DAG 文件每次被扫描时，整个文件作为 Python 代码从头到尾执行一遍**。这意味着你写在文件顶层的所有代码——import、函数定义、变量赋值、甚至网络请求——每次调度循环都会跑一次（实际频率是几秒到几十秒一次，取决于配置）。

所以 Airflow 的铁律之一是：**DAG 文件的顶层只放声明，不写实际做事的代码**。不要在顶层发 HTTP 请求、连数据库、读大文件。这些动作应该放进 task 内部（Operator 的 `execute` 方法或 PythonOperator 的函数体）。顶层只做"声明有哪些 task、依赖是什么"。

扫描完之后，每个 DAG 会做一次**拓扑排序**（topological sort）——把任务按依赖关系排成线性执行顺序。这是后续调度的基础。Airflow 内部用 networkx 来做这个排序。

## 调度最容易误解的概念：logical_date 与 data interval

在讲 Scheduler 的循环之前，必须先把 Airflow 里最容易让新手懵的概念讲清：**logical_date**（旧名 execution_date）和 **data interval**。

很多人初学 Airflow 的第一反应是："我设了 `schedule="0 2 * * *"`（每天凌晨2点），那 2 点触发的这次运行，处理的是哪天的数据？"

直觉答案是"今天的数据"，但 Airflow 的答案是"昨天的数据"。

原因在于 Airflow 的调度模型是"事后处理"：一次运行被触发时，它处理的是**刚刚过去那个完整周期**的数据。

举例：schedule 是每天一次，start_date 是 2024-01-01。

- 2024-01-02 凌晨 2 点，触发第一次 run
- 这次 run 的 logical_date = 2024-01-01
- data_interval = [2024-01-01 00:00, 2024-01-02 00:00)
- 它处理的"逻辑时间"是 1 月 1 日，数据区间是 1 月 1 日整天

为什么叫 logical_date？因为它不是"真实触发时间"（那是 1 月 2 日凌晨 2 点），而是"这次运行逻辑上代表的时间点"（1 月 1 日）。这个命名在 Airflow 2.2 改的，旧名 execution_date 误导性更强，让很多人以为"任务会在 1 月 1 日执行"。

理解了这一层，你就能解释 Airflow 几个常见困惑：

- 为什么 start_date 设了今天，但 DAG 今天不跑（今天还没结束，数据区间不完整，要等明天才触发处理今天的区间）
- `catchup=True` 时为什么会有一堆历史 run（从 start_date 到现在，每个周期补一个 run）
- 为什么 task 里 `{{ ds }}` 模板变量是 logical_date 的日期格式，而不是今天

这个概念是理解 Airflow 调度的钥匙。记住一句：**Airflow 在某个周期结束时触发运行，处理的是这个周期内的数据**。

## Scheduler 的循环：一个 tick 干了什么

Scheduler 是一个**常驻进程**，它的核心是一个无限循环，每个循环叫一个 **tick**。一个 tick 大致做这几件事：

1. **收集 DAG**：让 DagBag 扫描 DAG 文件夹，发现新增或变更的 DAG 文件，更新内部维护的 DAG 列表。
2. **处理 DAG 文件**：对每个 DAG 文件，交给一个 `DagFileProcessor`，解析 DAG、做拓扑排序、判断是否到了触发时间，如果是就创建 DagRun。
3. **调度 Task**：对每个 DagRun，检查每个 task 的上游状态，把"上游全 success 且未被调度过"的 task 标记为可执行，发送给 Executor。

Airflow 2.x 里，Scheduler 会启动多个 DagFileProcessor 并行处理不同的 DAG 文件，提高扫描吞吐。文件处理器是独立进程，因为导入 Python 文件可能慢或出错，独立进程能隔离故障——一个文件 import 报错不会拖垮整个 Scheduler。

几个关键配置参数：

- `scheduler heartbeat`：调度器心跳，默认 5 秒，是 Scheduler 主循环的节奏
- `min_file_process_interval`：同一个 DAG 文件两次处理之间的最小间隔，避免太频繁地 import
- `parsing_workers`：并行解析 DAG 文件的进程数
- `max_dagruns_to_create_per_loop`：每个循环最多创建多少个 DagRun，避免突发时雪崩

这些参数在调优时是杠杆。但初学不用纠结，先把行为搞懂。

## 一个 Task 的完整生命周期

现在聚焦到最微观的层次：单个 Task 的状态机。这是排查问题的核心——你看到 UI 上某个状态，得知道它意味着什么、下一步会去哪。

Task 实例的可能状态：

- **none**：初始，刚创建，啥也没干
- **scheduled**：Scheduler 看完上游，决定可以跑，已标记
- **queued**：已发给 Executor，在队列里排队等执行
- **running**：Worker 拿到，正在执行
- **success**：执行成功
- **failed**：执行失败
- **up_for_retry**：失败，但还有重试次数，会等一会儿重试
- **upstream_failed**：上游 task 失败了，所以这个 task 不用跑了
- **skipped**：被 BranchOperator 等机制跳过
- **deferred**：deferrable operator 主动让出 worker，等外部事件恢复（Airflow 2.2+）
- **removed**：DAG 从文件里删了，旧的 task 实例残留

![](.//imgs/learn-airflow/airflow_task_status.png)

排查问题时记住几个常见的"卡住"状态：

- 一直 `queued`：Executor 没消费，可能 worker 挂了或队列堵了
- 一直 `running`：task 真在跑（慢），还是僵尸（进程死了但状态没回写）
- 一直 `up_for_retry`：在等重试间隔，看 retry_delay 配置
- `deferred`：sensor 让出 slot 了，看有没有对应的 trigger 在恢复它

Airflow 有个机制叫 **zombie detection**（僵尸检测）：一个 task 状态是 running 但对应进程实际已经死了，Scheduler 会定期检查，把这种 task 标记为 failed。否则它会永远卡在 running。

## Executor：四种模型的取舍

Executor 决定 task 在哪里跑。选错 executor，要么撑不住规模，要么过度复杂。

**SequentialExecutor**。单线程串行，一次只跑一个 task。只能配合 sqlite 用，**仅限开发测试**。生产绝不能用，因为 sqlite 本身就不能并发。

**LocalExecutor**。在本机起多进程跑 task，并发数受 `parallelism` 参数控制。需要真正的数据库（Postgres/MySQL），不能用 sqlite。适合**中小规模、单机部署**——比如一个团队几十个 DAG、并发不超过一二十个 task。优点是部署简单，缺点是单机容量上限明显，机器挂了全挂。

**CeleryExecutor**。通过 Celery 框架 + 消息队列（Redis 或 RabbitMQ）把 task 分发到多台 worker 节点。这是**生产分布式标配**。Scheduler 把 task 推到队列，多台 worker 各自拉取执行。优点是水平扩展、可做高可用（worker 多台冗余）；代价是引入了 Celery + 队列中间件，运维复杂度上升。

**KubernetesExecutor**。每个 task 起一个独立的 K8s Pod 跑，跑完销毁。优点是**隔离性极强**（每个 task 独立环境、独立资源、独立故障域）和**弹性**（按需起 pod，不长期占用资源）。代价是 Pod 启动开销大（每个 task 要等几秒到几十秒拉镜像、起容器），不适合大量短任务。适合对隔离 / 弹性要求高、task 重量大的场景，比如 ML 训练、Spark 作业。

还有混合方案 **CeleryKubernetesExecutor**：轻任务走 Celery，重任务走 K8s，按队列分流。

选型决策很直接：

- 测试：Sequential
- 单机生产：Local
- 分布式稳态生产：Celery
- 强隔离 / 弹性 / 已有 K8s 体系：Kubernetes
- 既有轻任务又有重任务：CeleryKubernetes

## 端到端：一次 DagRun 的完整旅程

把前面的零件串起来，看一次 DagRun 从无到有的完整流程：

![DagRun 完整旅程流程图](/imgs/learn-airflow/airflow-dagrun-flow.png)

逐步拆解如下：

1. 你写好 DAG 文件放进 dags 目录
2. Scheduler 的 DagBag 扫到这个文件，import 执行，提取出 DAG 对象
3. 对 DAG 做拓扑排序，得到任务的执行顺序
4. 时间到了（满足 schedule，且当前时间已过该 data_interval 的结束 + 调度延迟）
5. Scheduler 创建一个 DagRun，logical_date 设为对应周期的起点
6. DagRun 创建后，所有 task 实例初始为 none
7. Scheduler 进入下一个 tick，检查这个 DagRun 的 task：
   - 入口 task（没有上游）的依赖已满足 → 标记 scheduled
8. scheduled 的 task 发给 Executor，状态变 queued
9. Worker 拿到 task，状态变 running，执行 Operator 的 `execute` 逻辑
10. execute 返回，状态变 success（或 failed）
11. 如果 success，Scheduler 下一个 tick 发现该 task 的下游上游已满足，标记下游为 scheduled，回到第 8 步
12. 如果 failed 且 retries > 0，状态变 up_for_retry，等 retry_delay 后重新走 8
13. 如果 failed 且 retries 用尽，下游变 upstream_failed
14. 所有 task success → DagRun 状态 success
15. 期间所有状态变化都实时写回 Metadata DB，Webserver 读取展示给你

整个链条里，**Metadata DB 是唯一的真相之源**。Scheduler、Executor、Worker 之间没有直接握手，全靠读写数据库的状态来协调。这也是为什么 Airflow 对数据库性能敏感——状态读写在热路径上。

## 几个影响理解的细节

**Scheduler 不是越频繁越好**。heartbeat 太短会让 DAG 文件被频繁 import，拖慢系统。生产默认 5 秒通常够。

**task 不是一到点就跑**。从 DagRun 创建到入口 task 真正 running，会有几秒到几十秒延迟——取决于 Scheduler tick 间隔、DAG 数量、Executor 队列深度。如果你的管道对启动延迟敏感，要把这个预期考虑进去。

**Webserver 不参与调度**。它只读数据库。所以 Webserver 挂了你的管道照跑，只是你看不到 UI。反过来，UI 上能看到的"卡住"状态，反映的是数据库里的真实状态。

**回填和正常调度走的是同一套机制**。`airflow dags backfill` 创建的是额外的 DagRun，走同样的 DagRun → task 调度流程。区别只是这些 DagRun 的 logical_date 是过去的时间点。
