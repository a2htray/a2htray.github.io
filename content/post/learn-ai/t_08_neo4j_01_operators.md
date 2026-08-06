+++
date = '2026-08-03T14:21:39+08:00'
draft = false
title = 'Neo4j 教程：逻辑运算符与关系检索'
categories = ['后端技术', '数据库']
tags = ['Neo4j', 'Cypher', '图数据库', 'T 系列', '运算符']
toc = true
+++

## 引言

在本教程中，将学习如何在 Cypher 查询中有效地使用**逻辑运算符**，并掌握在图数据库中创建、检索和修改关系的基本技能。

关系是图数据库强大之处的核心，虽然节点存储着实体，但正是它们之间的关系才真正释放了像 Neo4j 这样的图数据库的价值所在。

完成本教程后，你将能够熟练掌握如何操作关系，并结合**逻辑运算符**创建高效的查询。

## 初始化数据库

在深入操作符和关系之前，先创建一个示例数据库来使用，构建一个简单的社交网络，包含用户、他们的兴趣以及他们之间的各种联系。

**创建 Person 节点**

```bash
CREATE (john:Person {name: "John", age: 28, city: "New York"})
CREATE (sarah:Person {name: "Sarah", age: 26, city: "Boston"})
CREATE (mike:Person {name: "Mike", age: 32, city: "Chicago"})
CREATE (emily:Person {name: "Emily", age: 24, city: "New York"})
CREATE (david:Person {name: "David", age: 35, city: "San Francisco"})
CREATE (laura:Person {name: "Laura", age: 29, city: "Chicago"})
```

**创建 Interest 节点**

```bash
CREATE (music:Interest {name: "Music", type: "Art"})
CREATE (coding:Interest {name: "Coding", type: "Technology"})
CREATE (hiking:Interest {name: "Hiking", type: "Outdoor"})
CREATE (photography:Interest {name: "Photography", type: "Art"})
CREATE (cooking:Interest {name: "Cooking", type: "Culinary"})
```

## Cypher 中的逻辑运算符

我们先来了解如何在 Cypher 查询中使用逻辑运算符 IN、AND 和 OR。这些运算符可以帮助我们根据多个条件筛选结果。

**IN 运算符**

IN 运算符允许您检查属性值是否在指定的值列表中，它是多个 OR 条件的简写形式。

```bash
MATCH (p:Person)
WHERE p.city IN ["New York", "Chicago"]
RETURN p.name, p.city
```

执行结果：

```bash
| p.name  | p.city    |
|---------|-----------|
| "John"  | "New York"|
| "Emily" | "New York"|
| "Mike"  | "Chicago" |
| "Laura" | "Chicago" |
```

**AND 运算符**

AND 运算符用于组合多个条件，只有当所有条件都为真时，记录才会被包含在结果中。

```bash
MATCH (p:Person)
WHERE p.age > 25 AND p.city = "Chicago"
RETURN p.name, p.age, p.city
```

执行结果：

```bash
| p.name  | p.age | p.city    |
|---------|-------|-----------|
| "Mike"  | 32    | "Chicago" |
| "Laura" | 29    | "Chicago" |
```

**OR 运算符**

OR 运算符用于包含至少满足其中一个条件的记录。

```bash
MATCH (p:Person)
WHERE p.age < 25 OR p.age > 30
RETURN p.name, p.age
```

执行结果：

```bash
| p.name  | p.age |
|---------|-------|
| "Mike"  | 32    |
| "Emily" | 24    |
| "David" | 35    |
```

**组合多个运算符**

将这些运算符组合起来创建更复杂的查询，使用括号来控制计算顺序。

```bash
MATCH (p:Person)
WHERE p.city = "New York" AND (p.age < 25 OR p.age > 30)
RETURN p.name, p.age, p.city
```

执行结果：

```bash
| p.name  | p.age | p.city    |
|---------|-------|-----------|
| "Emily" | 24    | "New York"|
```

## 创建关系

现在我们来看看如何在图中创建节点之间的关系。Neo4j 中的关系是单向的，具有类型，并且可以包含属性。

> Neo4j 存储的关系底层是「有向（单向）」，但查询层面可以自由无视方向。

```bash
MATCH (john:Person {name: "John"}), (sarah:Person {name: "Sarah"})
CREATE (john)-[:FRIENDS_WITH {since: 2019}]->(sarah)

MATCH (mike:Person {name: "Mike"}), (john:Person {name: "John"})
CREATE (mike)-[:FRIENDS_WITH {since: 2018}]->(john)

MATCH (emily:Person {name: "Emily"}), (sarah:Person {name: "Sarah"})
CREATE (emily)-[:FRIENDS_WITH {since: 2020}]->(sarah)

MATCH (david:Person {name: "David"}), (mike:Person {name: "Mike"})
CREATE (david)-[:FRIENDS_WITH {since: 2017}]->(mike)

MATCH (laura:Person {name: "Laura"}), (david:Person {name: "David"})
CREATE (laura)-[:FRIENDS_WITH {since: 2021}]->(david)
```

**一次创建多条关系**

你可以在一个查询中创建多个关系，这比逐一创建更高效。

```bash
MATCH (john:Person {name: "John"}), (music:Interest {name: "Music"}),
      (hiking:Interest {name: "Hiking"})
CREATE (john)-[:INTERESTED_IN {level: "high"}]->(music),
       (john)-[:INTERESTED_IN {level: "medium"}]->(hiking)

MATCH (sarah:Person {name: "Sarah"}), (photography:Interest {name: "Photography"}),
      (cooking:Interest {name: "Cooking"})
CREATE (sarah)-[:INTERESTED_IN {level: "high"}]->(photography),
       (sarah)-[:INTERESTED_IN {level: "high"}]->(cooking)

MATCH (mike:Person {name: "Mike"}), (coding:Interest {name: "Coding"}),
      (hiking:Interest {name: "Hiking"})
CREATE (mike)-[:INTERESTED_IN {level: "high"}]->(coding),
       (mike)-[:INTERESTED_IN {level: "medium"}]->(hiking)

MATCH (emily:Person {name: "Emily"}), (music:Interest {name: "Music"}),
      (photography:Interest {name: "Photography"})
CREATE (emily)-[:INTERESTED_IN {level: "medium"}]->(music),
       (emily)-[:INTERESTED_IN {level: "high"}]->(photography)
```

**创建双向关系**

在某些情况下，若希望建立双向关系。这可以通过两条 CREATE 语句来实现。

```bash
MATCH (david:Person {name: "David"}), (laura:Person {name: "Laura"})
CREATE (david)-[:FRIENDS_WITH {since: 2021}]->(laura)
```

## 查询关系

已经建立了关系，接下来就来探索获取和查询这些关系的不同方法。

**基础**

```bash
MATCH (p1:Person)-[r:FRIENDS_WITH]->(p2:Person)
RETURN p1.name, p2.name, r.since
```

**通过属性过滤关系**

```bash
MATCH (p1:Person)-[r:FRIENDS_WITH]->(p2:Person)
WHERE r.since >= 2020
RETURN p1.name, p2.name, r.since
```

**检索特定的关系**

```bash
MATCH (john:Person {name: "John"})-[:FRIENDS_WITH]->(friend)
RETURN friend.name

MATCH (p:Person)-[:INTERESTED_IN]->(:Interest {name: "Hiking"})
RETURN p.name
```

**使用关系的逻辑运算符**

```bash
MATCH (p:Person)
WHERE (p)-[:INTERESTED_IN]->(:Interest {name: "Music"}) 
  AND (p)-[:INTERESTED_IN]->(:Interest {name: "Photography"})
RETURN p.name
```

**使用 IN 与关系类型**

```bash
MATCH (p:Person)-[r:INTERESTED_IN|CREATED|CONTRIBUTED_TO]->(i:Interest)
RETURN p.name, type(r), i.name
```

**检索具有可变长度关系的路径**

```bash
MATCH (john:Person {name: "John"})-[:FRIENDS_WITH*2]->(fof)
RETURN john.name, fof.name
```

## 修改关系

关系创建后，后期可进行修改，让我们看看如何更新关系属性和更改关系结构。

**更新关系属性**

```bash
MATCH (john:Person {name: "John"})-[r:FRIENDS_WITH]->(sarah:Person {name: "Sarah"})
SET r.since = 2018
RETURN john.name, sarah.name, r.since

MATCH (p:Person)-[r:INTERESTED_IN]->(i:Interest)
SET r.updated = date()
RETURN p.name, i.name, r.level, r.updated
```

**删除关系**

```bash
MATCH (mike:Person {name: "Mike"})-[r:FRIENDS_WITH]->(john:Person {name: "John"})
DELETE r

MATCH (:Person)-[r:INTERESTED_IN {level: "high"}]->(:Interest)
DELETE r
```

**替换关系**

有时你可能想用一段不同的关系来取代一段关系。

```bash
MATCH (mike:Person {name: "Mike"})-[r:INTERESTED_IN]->(:Interest {name: "Coding"})
WHERE r.level = "high"
DELETE r
CREATE (mike)-[:EXPERT_IN {since: 2015}]->(:Interest {name: "Coding"})
```

> 先删除，再创建。

**合并关系**

MERGE 命令可以用来匹配现有的关系，或者在关系不存在的情况下创建它们。

```bash
MATCH (emily:Person {name: "Emily"}), (mike:Person {name: "Mike"})
MERGE (emily)-[r:FRIENDS_WITH]->(mike)
ON CREATE SET r.since = 2022, r.through = "Work"
ON MATCH SET r.updated = date()
RETURN emily.name, mike.name, r
```

## 实际用例

让我们来探索一些实际的用例，这些用例结合了我们所学到的关于运算符和关系的所有知识。

**推荐引擎**

```bash
MATCH (john:Person {name: "John"})-[:FRIENDS_WITH]->(friend)-[:INTERESTED_IN]->(interest)
WHERE NOT (john)-[:INTERESTED_IN]->(interest)
RETURN interest.name, count(friend) as frequency
ORDER BY frequency DESC
```

根据 John 朋友的喜好向他推荐新的兴趣爱好。

**寻找相互联系**

```bash
MATCH (sarah:Person {name: "Sarah"})<-[:FRIENDS_WITH]-(mutualFriend)-[:FRIENDS_WITH]->(emily:Person {name: "Emily"})
RETURN mutualFriend.name
```

找到 Sarah 和 Emily 之间的共同朋友。

**复杂模式**

```bash
MATCH (p1:Person)-[:INTERESTED_IN]->(i:Interest)<-[:INTERESTED_IN]-(p2:Person)
WHERE p1.city = p2.city AND p1.name < p2.name
AND NOT (p1)-[:FRIENDS_WITH]-(p2)
RETURN p1.name, p2.name, p1.city, collect(i.name) as sharedInterests
```

找到住在同一个城市，至少有一个共同爱好，但还不是朋友的人。

## 高级技巧

下面是一些在 Neo4j 中处理关系的高级技巧。

**使用 FOREACH 创建条件关系**

```bash
MATCH (p1:Person), (p2:Person)
WHERE p1.city = "Chicago" AND p2.city = "Chicago" AND p1 <> p2
FOREACH (ignoreMe IN CASE WHEN NOT (p1)-[:FRIENDS_WITH]->(p2) THEN [1] ELSE [] END |
  CREATE (p1)-[:FRIENDS_WITH {since: 2022, reason: "Same city"}]->(p2)
)
```

在所有来自芝加哥的人之间建立关系。

**使用 UNWIND 创建批处理关系**

```bash
MATCH (david:Person {name: "David"})
UNWIND ["Coding", "Hiking", "Cooking"] as interestName
MATCH (interest:Interest {name: interestName})
CREATE (david)-[:INTERESTED_IN {level: "medium"}]->(interest)
```

同时创建多个兴趣关系。

**关系去重**

```bash
MATCH (p1:Person)-[r:FRIENDS_WITH]->(p2:Person)
WITH p1, p2, COLLECT(r) as rels
WHERE SIZE(rels) > 1
UNWIND TAIL(rels) as extraRel
DELETE extraRel
```

确保不存在重复的 FRIENDS_WITH 关系。

## 处理关系的最佳实践

* 仔细设计关系：考虑方向性、类型和属性，以最准确地建模你的领域。
* 使用有意义的关系类型：如 FRIENDS_WITH、BELONGS_TO 或 PURCHASED 这样的名称，能清晰地传达关系的性质。
* 优化关系遍历：Neo4j 的强大之处在于高效的图关系遍历。
* 对可变长度路径要谨慎：虽然功能强大，但像 [:FRIENDS_WITH*1..5] 这样的查询在大型数据库上可能会变得代价较大。
* 用于过滤器的索引属性：对于在 WHERE 子句中经常使用的属性，应创建相应的索引。
* 考虑双向关系：对于像友谊这样的关系，如果双向遍历很常见，建议创建双向关系。
* 在查询中使用参数：避免硬编码值，使用参数以提高安全性并缓存查询计划。