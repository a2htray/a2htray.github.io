+++
date = '2026-10-27T16:32:28+08:00'
draft = false
title = '声明式访问控制引擎 Casbin'
categories = ['后端技术', 'Casbin']
tags = ['Python', 'Casbin', '权限控制']
toc = true
+++

![](/imgs/learn-backend/2026-10-28.png)

Casbin 是一个基于声明式访问控制引擎，用于在应用中实现**访问控制策略**。

> Casbin 只负责授权，不负责认证。

## 基础概念

Casbin 包含 5 个核心组件：

1. **模型（Model）**：定义一套权限“语法框架”，规定请求长什么样、策略字段有哪些、怎么匹配、命中后怎么算结果。
2. **策略（Policy）**：具体的权限记录，是填充到 Model 模板里的实例数据。
3. **执行器（Enforcer）**：调度中心，绑定 Model + Adapter，对外提供鉴权、策略管理 API。
4. **适配器（Adapter）**：隔离 Enforcer 和存储介质，负责 Policy 的读。
5. **匹配器（Matcher）**：定义了主体和资源的匹配规则，用于判断是否命中策略。

> Enforcer **不是线程安全**。生产一般全局单例；写策略的时候建议加锁。

5 个组件之间的关系如下图：

![](/imgs/learn-backend/2026-10-27.png)

## 配置文件

Casbin 配置文件（`model.conf`）是一个包含权限“语法框架”的文件，用于定义主体、资源和操作的类型，包括：

```conf
[request_definition]   # r：请求结构
[policy_definition]    # p：单条策略结构
[role_definition]      # g：角色/分组关系（RBAC专属）
[policy_effect]        # e：多条策略冲突时的判定规则
[matchers]             # m：匹配表达式（核心逻辑）
```

## 权限模型

### ACL - 控制访问列表

ACL 控制访问列表是最简单的权限模型，直接定义每一个主体对某一资源的访问权限。

配置文件 `model.conf` 如下：

```conf
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
```

策略文件 `policy.csv` 如下：

```csv
p, alice, article, read
p, alice, article, write
p, bob, article, read
```

> 策略的持久化存储可以使用文件、数据库等。

在程序中使用 Enforcer 配置模型和策略，即可实现 ACL 权限模型控制访问列表。

```python
from casbin import Enforcer

enforcer = Enforcer("model.conf", "policy.csv")

print(enforcer.enforce("alice", "article", "read")) # True
print(enforcer.enforce("bob", "article", "write")) # False
```

* `Enforcer` 加载模型和策略文件
* `enforce` 方法用于判断主体是否有权限访问资源

### RBAC - 基于角色访问控制

RBAC 基于角色访问控制是一种基于角色的权限模型，将主体和资源的访问权限分配给角色，而不是直接分配给主体。

配置文件 `model.conf` 如下：

```conf
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && r.act == p.act
```

`g = _, _` 中 `g` 是 **grouping 函数**，专门用来描述“主体 ↔ 角色”的映射关系，也就是 RBAC 的用户 - 角色绑定。`_` 是占位符，分别表示主体和角色。

`[matchers]` 节中 `g(r.sub, p.sub)` 是 **角色匹配器**，用于判断主体 r.sub 是否属于角色 p.sub。

`keyMatch(r.obj, p.obj)` 是**匹配器**函数，用于判断主体要访问的资源是否匹配策略中的资源。

策略文件 `policy.csv` 如下：

```csv
p, admin, /api/*, write
p, admin, /api/*, read
p, viewer, /api/info, read
g, alice, admin
g, bob, viewer
```

验证程序如下：

```python
from casbin import Enforcer

enforcer = Enforcer("model.conf", "policy.csv")
# alice 是 admin 角色，对 `/api/*` 有 write 权限
print(enforcer.enforce("alice", "/api/user", "write")) # True
# bob 是 viewer 角色，对 `/api/*` 无 write 权限
print(enforcer.enforce("bob", "/api/user", "write")) # False
# bob 是 viewer 角色，对 `/api/info` 有 read 权限
print(enforcer.enforce("bob", "/api/info", "read")) # True

enforcer.add_role_for_user("bob", "admin")
# 给 bob 加了 admin 角色，所以对 `/api/*` 有 write 权限
print(enforcer.enforce("bob", "/api/user", "write")) # True
```

### ABAC - 基于属性访问控制

ABAC 基于属性访问控制是一种基于属性的权限模型，通过对主体属性判断是否对资源有访问权限。

配置文件 `model.conf` 如下：

```conf
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub_rule, obj_rule, act_rule

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = eval(p.sub_rule) && eval(p.obj_rule) && eval(p.act_rule)
```

验证程序如下：

```python
from casbin import Enforcer

enforcer = Enforcer("model.conf")
# 增加一条规则：用户 id == 资源 owner，动作是 read 则放行
enforcer.add_policy(
    "r.sub['user_id'] == r.obj['owner_id']",
    "True",
    "r.act == 'read'"
)

# 用户
user_a = {"user_id": 1001, "dept": "dev"}
# 两条资源
article_own = {"owner_id": 1001, "title": "我的笔记"}
article_other = {"owner_id": 2002, "title": "别人的笔记"}

print(enforcer.enforce(user_a, article_own, "read"))    # True
print(enforcer.enforce(user_a, article_other, "read"))  # False
print(enforcer.enforce(user_a, article_own, "write"))  # False，规则只允许 read
```

通常，在企业级应用中，采用 `RBAC + ABAC` 来实现权限控制：

- RBAC 控制用户对某一模块的访问权限
- ABAC 控制用户对行级数据的访问权限


### ReBAC - 基于关系的访问控制

ReBAC 基于关系的访问控制是一种基于关系的权限模型，通过对主体和资源的关系判断是否对资源有访问权限。如用户与文档的关系可以是 "owner" 或 "reader"。

配置文件 `model.conf` 如下：

```conf
[request_definition]
r = sub, rel, obj, act

[policy_definition]
p = rel, act, eft

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.rel == p.rel && r.act == p.act
```

策略文件 `policy.csv` 如下：

```csv
p, owner, read, allow
p, owner, write, allow
p, owner, delete, allow
p, editor, read, allow
p, editor, write, allow
p, viewer, read, allow
```

验证程序如下：

```python
from casbin import Enforcer

enforcer = Enforcer("model.conf", "policy.csv")

user_resource_rel = {
    ("alice", "doc1"): "owner",
    ("bob", "doc1"): "viewer",
    ("bob", "doc2"): "editor"
}

def check_perm(user: str, doc: str, act: str) -> bool:
    rel = user_resource_rel.get((user, doc))
    if rel is None:
        return False
    return enforcer.enforce(user, rel, doc, act)

print("alice 是 doc1 owner，write：", check_perm("alice", "doc1", "write"))   # True
print("bob 是 doc1 viewer，write：", check_perm("bob", "doc1", "write"))     # False
print("bob 是 doc2 editor，write：", check_perm("bob", "doc2", "write"))     # True
print("bob 是 doc2 editor，delete：", check_perm("bob", "doc2", "delete"))   # False
```

- user_resource_rel 是一个字典，键为 (user, doc)，值为用户与文档的关系
- enforcer.enforce 的参数对应 model.conf 中的 `r` 的四个字段

## 不只 csv

上述的示例中，策略都是使用 csv 文件存储的。在实际的应用中，数据库存储才是首选。

所以下面的示例，使用 sqlite 数据库存储策略，而不是 csv 文件。

```python
import csv

from casbin import Enforcer
from casbin_sqlalchemy_adapter import Adapter

adapter = Adapter(f"sqlite:///./rbac.db")
enforcer = Enforcer("./model.conf", adapter)

if not enforcer.get_grouping_policy() or not enforcer.get_policy():
    with open("./policy.csv", newline="", encoding="utf-8") as f:
        for row in csv.reader(f):
            if not row:
                continue
            if row[0] == "p":
                enforcer.add_policy(row[1].strip(), row[2].strip(), row[3].strip())
            elif row[0] == "g":
                enforcer.add_grouping_policy(row[1].strip(), row[2].strip())

enforcer.build_role_links()

print(enforcer.enforce("alice", "/api/user", "write"))
print(enforcer.enforce("bob", "/api/user", "write"))
print(enforcer.enforce("bob", "/api/info", "read"))

enforcer.add_role_for_user("bob", "admin")
print(enforcer.enforce("bob", "/api/user", "read"))
```

安装 `casbin_sqlalchemy_adapter` 依赖：

```bash
pip install casbin_sqlalchemy_adapter
```

先前是传递 `policy.csv` 文件，现在是传递 `Adapter`，在 `Adapter` 实例化后，会创建数据库表，表创建 SQL 如下：

```sql
create table casbin_rule
(
    id    INTEGER not null
        primary key,
    ptype VARCHAR(255),
    v0    VARCHAR(255),
    v1    VARCHAR(255),
    v2    VARCHAR(255),
    v3    VARCHAR(255),
    v4    VARCHAR(255),
    v5    VARCHAR(255)
);
```

`casbin_rule` 表中存储策略规则和组策略规则。

当表里没有策略时（`not enforcer.get_grouping_policy() or not enforcer.get_policy()` 为 True），读取 `policy.csv` 文件，将策略写入数据库表。

载入后显式调用 `enforcer.build_role_links()`，丢弃初始化时建立的陈旧角色链接。