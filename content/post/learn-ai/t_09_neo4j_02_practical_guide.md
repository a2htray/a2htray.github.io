+++
date = '2026-08-05T16:33:05+08:00'
draft = false
title = 'Neo4j 教程：图形查询语言的实用指南'
categories = ['后端技术', '数据库']
tags = ['Neo4j', 'Cypher', '图数据库', 'T 系列', '实用指南']
toc = true
+++

## 引言

Cypher 对于 Neo4j 就像 SQL 对于关系型数据库一样—，它是一种优雅的**声明式查询语言**，让你无需精确指定如何检索数据，即可描述想要从数据库中获取的内容。Cypher 的独特之处在于其基于 ASCII 艺术的可视化语法，能够模仿并匹配你图中所需的图形结构模式。

> ASCII 艺术（ASCII Art，简称 AA）是一种艺术形式，通过使用字符和符号来模仿图形、形状、图像或其他视觉效果。

让我们通过实际例子来探索 Cypher，结合真实场景展示其强大功能和灵活性。

## 开始使用 Cypher 语法

Cypher 的语法设计得直观且易于理解，节点用括号 `()` 表示，关系用箭头 `-[]->` 表示。这使得查询语句易于阅读，类似于在白板上绘制想要查找的模式。

**节点与关系模式基础**

```bash
MATCH (p:Person {name: "John"})
RETURN p

MATCH (p:Person {name: "John"})-[:FRIENDS_WITH]->(friend)
RETURN friend.name
```

在这些示例中：
* `()` 表示一个节点  
* `:Person` 是用于分类该节点的标签  
* `{name: "John"}` 是属性约束  
* `-[:FRIENDS_WITH]->` 表示具有类型为 "FRIENDS_WITH" 的**有向关系**

## 使用 Cypher 创建数据

从创建一个小型社交网络数据集开始：

**创建节点**

```bash
CREATE (alice:Person {name: "Alice", age: 32, occupation: "Data Scientist"})
CREATE (bob:Person {name: "Bob", age: 35, occupation: "Software Engineer"})
CREATE (charlie:Person {name: "Charlie", age: 28, occupation: "UX Designer"})
CREATE (diana:Person {name: "Diana", age: 41, occupation: "Project Manager"})
CREATE (edward:Person {name: "Edward", age: 25, occupation: "Data Analyst"})

CREATE (graphdb:Interest {name: "Graph Databases", category: "Technology"})
CREATE (cycling:Interest {name: "Cycling", category: "Sports"})
CREATE (cooking:Interest {name: "Cooking", category: "Hobby"})
CREATE (photography:Interest {name: "Photography", category: "Arts"})
CREATE (travel:Interest {name: "Travel", category: "Lifestyle"})
```

**创建关系**

```bash
MATCH (alice:Person {name: "Alice"}), (bob:Person {name: "Bob"})
CREATE (alice)-[:FRIENDS_WITH {since: 2018}]->(bob)

MATCH (alice:Person {name: "Alice"}), (charlie:Person {name: "Charlie"})
CREATE (alice)-[:FRIENDS_WITH {since: 2020}]->(charlie)

MATCH (bob:Person {name: "Bob"}), (diana:Person {name: "Diana"})
CREATE (bob)-[:FRIENDS_WITH {since: 2015}]->(diana)

MATCH (charlie:Person {name: "Charlie"}), (diana:Person {name: "Diana"})
CREATE (charlie)-[:FRIENDS_WITH {since: 2019}]->(diana)

MATCH (diana:Person {name: "Diana"}), (edward:Person {name: "Edward"})
CREATE (diana)-[:FRIENDS_WITH {since: 2021}]->(edward)

// Create interest relationships
MATCH (alice:Person {name: "Alice"}), (graphdb:Interest {name: "Graph Databases"})
CREATE (alice)-[:INTERESTED_IN {level: "Expert"}]->(graphdb)

MATCH (alice:Person {name: "Alice"}), (cycling:Interest {name: "Cycling"})
CREATE (alice)-[:INTERESTED_IN {level: "Intermediate"}]->(cycling)

MATCH (bob:Person {name: "Bob"}), (graphdb:Interest {name: "Graph Databases"})
CREATE (bob)-[:INTERESTED_IN {level: "Beginner"}]->(graphdb)

MATCH (bob:Person {name: "Bob"}), (cooking:Interest {name: "Cooking"})
CREATE (bob)-[:INTERESTED_IN {level: "Advanced"}]->(cooking)

MATCH (charlie:Person {name: "Charlie"}), (photography:Interest {name: "Photography"})
CREATE (charlie)-[:INTERESTED_IN {level: "Expert"}]->(photography)

MATCH (diana:Person {name: "Diana"}), (travel:Interest {name: "Travel"})
CREATE (diana)-[:INTERESTED_IN {level: "Advanced"}]->(travel)

MATCH (diana:Person {name: "Diana"}), (cooking:Interest {name: "Cooking"})
CREATE (diana)-[:INTERESTED_IN {level: "Intermediate"}]->(cooking)

MATCH (edward:Person {name: "Edward"}), (graphdb:Interest {name: "Graph Databases"})
CREATE (edward)-[:INTERESTED_IN {level: "Beginner"}]->(graphdb)

MATCH (edward:Person {name: "Edward"}), (cycling:Interest {name: "Cycling"})
CREATE (edward)-[:INTERESTED_IN {level: "Advanced"}]->(cycling)
```

## 使用 Cypher 查询数据

探索不同类型的查询，来提取数据中有价值的信息。

**通过标签和属性查找节点**

```bash
MATCH (p:Person)
RETURN p

MATCH (p:Person)
WHERE p.age > 30
RETURN p.name, p.age, p.occupation
```

**检索关系**

```bash
MATCH (p1:Person)-[r:FRIENDS_WITH]->(p2:Person)
RETURN p1.name, p2.name, r.since

MATCH (p:Person {name: "Alice"})-[:FRIENDS_WITH]->(friend)
RETURN friend.name
```

**检索朋友的朋友**

```bash
MATCH (alice:Person {name: "Alice"})-[:FRIENDS_WITH]->(friend)-[:FRIENDS_WITH]->(fof)
WHERE NOT (alice)-[:FRIENDS_WITH]->(fof) AND alice <> fof
RETURN DISTINCT fof.name as FriendOfFriend
```

> 寻找 Alice 的朋友，并再扩展一层人际关系、排除 Alice。

**检索共同兴趣**

```bash
MATCH (alice:Person {name: "Alice"})-[:INTERESTED_IN]->(interest)<-[:INTERESTED_IN]-(other)
WHERE alice <> other
RETURN other.name as Person, interest.name as SharedInterest
```

**聚合和排序**

```bash
MATCH (p:Person)-[:FRIENDS_WITH]->(friend)
RETURN p.name as Person, COUNT(friend) as NumberOfFriends
ORDER BY NumberOfFriends DESC

MATCH (i:Interest)<-[:INTERESTED_IN]-(p:Person)
RETURN i.name as Interest, COUNT(p) as Popularity
ORDER BY Popularity DESC
```

## 进阶技巧

**路径变量和函数**

```bash
MATCH path = shortestPath((alice:Person {name: "Alice"})-[:FRIENDS_WITH*]-(edward:Person {name: "Edward"}))
RETURN path

MATCH path = shortestPath((alice:Person {name: "Alice"})-[:FRIENDS_WITH*]-(edward:Person {name: "Edward"}))
RETURN [node in nodes(path) | node.name] as People, length(path) as PathLength
```

> 找出 Alice 和 Edward 之间最短的路径。

> 返回路径的长度。

**集合操作**

```bash
MATCH (p:Person)-[:INTERESTED_IN]->(i:Interest)
RETURN p.name as Person, COLLECT(i.name) as Interests

MATCH (p:Person)
WHERE ALL(interest IN ["Graph Databases", "Cycling"] 
          WHERE (p)-[:INTERESTED_IN]->(:Interest {name: interest}))
RETURN p.name
```

> `COLLECT` 收集每个人的所有兴趣。 

> 找到对 `["Graph Databases", "Cycling"]` 感兴趣的人。

**CASE 表达式**

```bash
MATCH (p:Person)
RETURN p.name, 
       CASE
         WHEN p.age < 30 THEN "Young Professional"
         WHEN p.age >= 30 AND p.age < 40 THEN "Mid-career"
         ELSE "Senior Professional"
       END AS AgeCategory
```

## 查询优化

随着图的增长，优化查询变得非常重要。这里有一些技巧：

**使用索引**

```bash
CREATE INDEX person_name FOR (p:Person) ON (p.name)

CREATE INDEX interest_name FOR (i:Interest) ON (i.name)
```

> `CREATE INDEX` 创建索引。

**查询分析**

```bash
PROFILE MATCH (p:Person {name: "Alice"})-[:FRIENDS_WITH*1..3]-(other)
RETURN other.name
```

> 类似于 SQL 中的 EXPLAIN。

## 更贴合实际的例子

运用所学的知识来解决一些常见的图形问题：

**推荐系统**

```bash
MATCH (alice:Person {name: "Alice"})-[:FRIENDS_WITH]->(friend)-[:INTERESTED_IN]->(interest)
WHERE NOT (alice)-[:INTERESTED_IN]->(interest)
RETURN interest.name as RecommendedInterest, COUNT(friend) as CommonFriends
ORDER BY CommonFriends DESC
```

> 基于 Alice 的朋友的兴趣，推荐 Alice 可能感兴趣的领域。

**网络分析**

```bash
MATCH (alice:Person {name: "Alice"})-[:FRIENDS_WITH*1..2]-(person)-[:FRIENDS_WITH]-(connection)
RETURN person.name, COUNT(DISTINCT connection) as Connections
ORDER BY Connections DESC
LIMIT 1
```

> 找到 Alice 长度在 2 内的朋友及朋友数量。

**模式识别**

```bash
MATCH (p1:Person)-[:FRIENDS_WITH]-(p2:Person)-[:FRIENDS_WITH]-(p3:Person)-[:FRIENDS_WITH]-(p1)
WHERE p1.name < p2.name AND p2.name < p3.name  // To avoid duplicate results
RETURN p1.name, p2.name, p3.name
```

> 找到互为朋友的三元组。

## 建议

编写 Cypher 查询的最佳实践
* 从简单的模式开始，逐渐增加复杂性
* 选择描述性强的变量名，使查询更易读
* 开发时使用 LIMIT 限制结果集大小进行测试
* 为复杂的逻辑添加注释
