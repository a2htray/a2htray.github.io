+++
date = '2026-07-16T10:58:49+08:00'
draft = false
title = '和 CodeBuddy 聊了 7 天，开发了一款财务系统'
categories = ['人工智能', '应用平台']
tags = ['CodeBuddy', '财务系统', '生产力工具', 'AI 辅助编码']
toc = true
+++

放在以前，单个这个标题，我一般就直接关闭了。因为我会想：

> “7 天开发一个系统，哪怕二次开发、CURD 也要半个月吧，做出来也是个没用的系统。”

现在，对上述的观点，我保留一半，能不能要在后续的开发中持续验证。

## 背景

朋友最近接手了一家小公司，有一些财务上的代帐和报税的需求，我爱人就跟我说：

![](/imgs/learn-ai/codebuddy-01.png)

她是看了我发的关于 Multica 的公众号，所以才这么跟我说的。如果在以前，我都会直接拒绝，原因很简单：没个财务相关的领域知识，不好做。

当时想法是，先找些一些现有的开源项目改改，再让 AI 改改，指不定能做出来。**首先，我锁定了技术栈（Go 或 Python），然后，问了豆包，最后放弃了**。原因有三：

1. 成熟财务系统以 Java 为主，看不懂
2. Go 或 Python 项目，功能不全，或不适配国内，或二次开发工作量大
3. 没个背景，看懂源码是个问题

所以决定让 AI 写，让 AI 教我写。

## CodeBuddy

在用 CodeBuddy 之前，我用的是 Claude Code，平时用用 OpenRouter 的免费模型，但免费额度有限，也不好用。这个时候也正好看到 Hy3 的新闻，新用户免费试用两周，就在 VS Code 上安装了 CodeBuddy 插件并开始深度使用。

## 开发过程

### 熟悉财务知识

了解新领域的知识，我一般都是从术语开始，对应用到系统开发中，它可能是一个实体模型、操作、流程等。刚开始还没有着急让 CodeBuddy 代码，就是和它聊天：

* 账务系统都有哪些术语？
* 个别术语都是什么意思？
* 账务系统一般都有哪些功能？

慢慢地，我缩小了问题的范围，有针对性地问：

1. 什么是**帐套**？
2. 会议科目都有哪些，怎么分类，借代方向又是什么意思？
3. 凭证有什么作用，如果做凭证，其数据来源是什么
4. 资产怎么做帐，工资又怎么做帐？

问了几轮后，一个账务系统的基本轮廓就已经出来了，如下：

```bash
帐套 Book ──< BookUser >── 用户 User ──< UserRole >── 角色 Role ──< 权限 Permission
```

```bash
科目 Account（树形层级，余额方向借/贷）
   └─ 期初 OpeningBalance（期初/累计/年初，带符号）
        └─ 凭证 Voucher → JournalEntry（借贷分录）
             └─ 保存即更新 Account.Balance（事务内）
```

![](/imgs/learn-ai/codebuddy-02.png)

上面的领域知识也是在和 CodeBuddy 聊天、编程过程中逐步建立起来的，一些关键性的结论，我也会让他记录到特定的 markdown 文件中，以便后续的参考。

### 确定技术栈

- 后端：Go 1.25、Gin、GORM（PostgreSQL 生产）、JWT 鉴权、CORS
- 前端：React 18、TypeScript、Vite、React Router、Axios
- 领域模型：复式记账（科目 / 凭证 / 分录），保存凭证时自动校验借贷平衡并更新科目余额
- 权限模型：基于 RBAC 的用户角色管理 + 帐套（Book）数据隔离

### 搭建基础框架

后端系统该有的东西都得有，如认证授权、用户权限，生成的第一个版本，过于简单。

用户直接绑定了角色，角色对应的权限也是通过硬编码的方式进行映射，然后我让 CodeBuddy 基于 RBAC 的表结构设计重新设计了数据库，目前勉强符合要求，现阶段的 RBAC 的设计如下：

表结构

![](/imgs/learn-ai/codebuddy-03.png)

权限设计

| 分组 | 前缀示例 | 说明 |
|---|---|---|
| 系统设置 | `user:view/manage`、`role:*` | 用户/角色管理 |
| 基础设置 | `base:book:*`、`base:department:*`、`base:employee:*`、`base:project:*`、`base:audit:view` | 帐套/部门/职员/项目/日志 |
| 账务设置 | `account:*`、`opening-balance:*`、`currency:*`、`unit:*`、`salary:config` | 科目/期初/币别/单位 |
| 业务核算 | `voucher:*`、`fund-journal:*`、`category:*`、`fund-account:*`、`invoice:*`、`fixed-asset:*`、`asset-category:*`、`salary:*`、`report:view` | 日常核算与报表 |

系统默认内置了 3 个角色：超级管理员、帐套管理员、会计，每个角色对应不同的权限，后台必然要支持权限分配功能：

![](/imgs/learn-ai/codebuddy-04.png)

权限分配的弹出层的第一个版本，很丑，后面让 CodeBuddy 美化了一个，但也没有说明我具体想要什么样，但第二个版本就是上面的样子。

### 实现业务逻辑

系统的基础功能实现之后，难点在于业务功能的实现，毕竟不懂财务知识。在了解一些基本的概念之后，我尝试着去使用开源项目的 demo 系统，“照虎画猫”，倒也是一步步实现了基本的业务流程。

以下是我和 CodeBuddy 的对话，截取了部分：

**明细表和汇总表的实现**

![](/imgs/learn-ai/codebuddy-05.png)

**生成的代码字段冗余，让 CodeBuddy 做删减**

![](/imgs/learn-ai/codebuddy-06.png)

**确认新增发票的需求**

![](/imgs/learn-ai/codebuddy-07.png)

**连续发需求**

![](/imgs/learn-ai/codebuddy-08.png)

**最初实现还能使用帐套进行数据隔离，但实现多轮后，CodeBuddy 好像忘记了这个设计，现在让其加上**

![](/imgs/learn-ai/codebuddy-09.png)

### 现阶段效果

Dashboard

![](/imgs/learn-ai/codebuddy-10.png)

凭证录入

![](/imgs/learn-ai/codebuddy-11.png)

日记帐

![](/imgs/learn-ai/codebuddy-12.png)

![](/imgs/learn-ai/codebuddy-13.png)

核对总帐

![](/imgs/learn-ai/codebuddy-14.png)

发票管理

![](/imgs/learn-ai/codebuddy-15.png)

![](/imgs/learn-ai/codebuddy-16.png)

会计科目

![](/imgs/learn-ai/codebuddy-17.png)

帐套管理

![](/imgs/learn-ai/codebuddy-18.png)

目前已经在线上发布了测试环境，有需要体验的可以后台联系我。

## 真心话

真心话放在最后，我是认真的。

代码百分百都是 CodeBuddy 写的，我在其中的作用就是跟他对话：确定技术栈、说明需求、调整实现方向。必要时，我会阶段性地 review 代码，并让其优化，当然优化方向是我给的。

代码中有我看不明白的，也是 CodeBuddy 帮我生成文档，教会了我绝大部分财务知识，当然也不多。

现在这个系统能不能用，肯定是不能用的，因为我还没有完全看懂里面的计算逻辑，需要在后面慢慢打磨，边用、边学、边改。

如果给这个系统打分，给到 6 分（10 分制），给 CodeBuddy 打分的话，给到 8.5 分，绝大多数情况下，能在单次对话中完成，少数情况也不超过 3 次对话。

用过其他的一些 Code Agent，他们的底层能力还是取决于大模型，所以要夸一下 Hy3，至少让我感觉用起来和 GPT 5.5、opus 4.7 差不多。

**主要还是场景相对简单**
