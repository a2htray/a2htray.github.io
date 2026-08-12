+++
date = '2026-08-11T09:23:27+08:00'
draft = false
title = 'Neo4j 教程：掌握变长关系与路径算法'
categories = ['后端技术', '数据库']
tags = ['Neo4j', 'Cypher', '图数据库', 'T 系列', '变长', '最短路径', 'Neo4j 教程']
toc = true
+++

## 引言

欢迎来到这篇深度教程，我们将探讨 Neo4j 最强大的特性之一：**处理变长路径和关系**。本指南将带你学习如何遍历复杂网络、查找实体之间的联系，以及在图数据中发现最优路线。

图数据库在分析关联数据方面表现卓越，Neo4j 提供了专门的语法和算法来高效探索这些联系。无论你是在构建社交网络、知识图谱、物流系统还是推荐引擎，本教程中介绍的技术都将帮助你从互联数据中挖掘出有价值的信息。

通过本教程的学习，你将掌握以下能力：

- 查询变长关系
- 查找节点之间的最短路径
- 发现节点之间的所有最短路径
- 将这些技术应用于解决实际问题

## 搭建示例数据库

在本教程中，我们将使用一个表示城市间公路连接的交通网络作为示例。这个例子将帮助我们理解如何在图中**查找路线**——这是图数据中路径查找的一个典型应用场景。

首先创建示例数据库：

```bash
// 创建城市节点
CREATE (nyc:City {name: "New York", population: 8400000})
CREATE (la:City {name: "Los Angeles", population: 3900000})
CREATE (chicago:City {name: "Chicago", population: 2700000})
CREATE (houston:City {name: "Houston", population: 2300000})
CREATE (phoenix:City {name: "Phoenix", population: 1600000})
CREATE (philly:City {name: "Philadelphia", population: 1500000})
CREATE (dallas:City {name: "Dallas", population: 1300000})
CREATE (sf:City {name: "San Francisco", population: 880000})
CREATE (denver:City {name: "Denver", population: 710000})
CREATE (boston:City {name: "Boston", population: 690000})
CREATE (seattle:City {name: "Seattle", population: 730000})
CREATE (atlanta:City {name: "Atlanta", population: 500000})
CREATE (miami:City {name: "Miami", population: 450000})
CREATE (austin:City {name: "Austin", population: 950000})
CREATE (dc:City {name: "Washington DC", population: 690000})

// 创建 ROAD_TO 关系并附带距离（英里）
MATCH (nyc:City {name: "New York"}), (philly:City {name: "Philadelphia"})
CREATE (nyc)-[:ROAD_TO {distance: 95}]->(philly)
CREATE (philly)-[:ROAD_TO {distance: 95}]->(nyc)

MATCH (nyc:City {name: "New York"}), (boston:City {name: "Boston"})
CREATE (nyc)-[:ROAD_TO {distance: 215}]->(boston)
CREATE (boston)-[:ROAD_TO {distance: 215}]->(nyc)

MATCH (philly:City {name: "Philadelphia"}), (dc:City {name: "Washington DC"})
CREATE (philly)-[:ROAD_TO {distance: 140}]->(dc)
CREATE (dc)-[:ROAD_TO {distance: 140}]->(philly)

MATCH (dc:City {name: "Washington DC"}), (atlanta:City {name: "Atlanta"})
CREATE (dc)-[:ROAD_TO {distance: 600}]->(atlanta)
CREATE (atlanta)-[:ROAD_TO {distance: 600}]->(dc)

MATCH (atlanta:City {name: "Atlanta"}), (miami:City {name: "Miami"})
CREATE (atlanta)-[:ROAD_TO {distance: 660}]->(miami)
CREATE (miami)-[:ROAD_TO {distance: 660}]->(atlanta)

MATCH (atlanta:City {name: "Atlanta"}), (dallas:City {name: "Dallas"})
CREATE (atlanta)-[:ROAD_TO {distance: 780}]->(dallas)
CREATE (dallas)-[:ROAD_TO {distance: 780}]->(atlanta)

MATCH (dallas:City {name: "Dallas"}), (houston:City {name: "Houston"})
CREATE (dallas)-[:ROAD_TO {distance: 240}]->(houston)
CREATE (houston)-[:ROAD_TO {distance: 240}]->(dallas)

MATCH (houston:City {name: "Houston"}), (austin:City {name: "Austin"})
CREATE (houston)-[:ROAD_TO {distance: 165}]->(austin)
CREATE (austin)-[:ROAD_TO {distance: 165}]->(houston)

MATCH (dallas:City {name: "Dallas"}), (austin:City {name: "Austin"})
CREATE (dallas)-[:ROAD_TO {distance: 195}]->(austin)
CREATE (austin)-[:ROAD_TO {distance: 195}]->(dallas)

MATCH (dallas:City {name: "Dallas"}), (phoenix:City {name: "Phoenix"})
CREATE (dallas)-[:ROAD_TO {distance: 1065}]->(phoenix)
CREATE (phoenix)-[:ROAD_TO {distance: 1065}]->(dallas)

MATCH (phoenix:City {name: "Phoenix"}), (la:City {name: "Los Angeles"})
CREATE (phoenix)-[:ROAD_TO {distance: 370}]->(la)
CREATE (la)-[:ROAD_TO {distance: 370}]->(phoenix)

MATCH (la:City {name: "Los Angeles"}), (sf:City {name: "San Francisco"})
CREATE (la)-[:ROAD_TO {distance: 380}]->(sf)
CREATE (sf)-[:ROAD_TO {distance: 380}]->(la)

MATCH (sf:City {name: "San Francisco"}), (seattle:City {name: "Seattle"})
CREATE (sf)-[:ROAD_TO {distance: 810}]->(seattle)
CREATE (seattle)-[:ROAD_TO {distance: 810}]->(sf)

MATCH (seattle:City {name: "Seattle"}), (denver:City {name: "Denver"})
CREATE (seattle)-[:ROAD_TO {distance: 1300}]->(denver)
CREATE (denver)-[:ROAD_TO {distance: 1300}]->(seattle)

MATCH (denver:City {name: "Denver"}), (dallas:City {name: "Dallas"})
CREATE (denver)-[:ROAD_TO {distance: 790}]->(dallas)
CREATE (dallas)-[:ROAD_TO {distance: 790}]->(denver)

MATCH (denver:City {name: "Denver"}), (chicago:City {name: "Chicago"})
CREATE (denver)-[:ROAD_TO {distance: 1000}]->(chicago)
CREATE (chicago)-[:ROAD_TO {distance: 1000}]->(denver)

MATCH (chicago:City {name: "Chicago"}), (nyc:City {name: "New York"})
CREATE (chicago)-[:ROAD_TO {distance: 790}]->(nyc)
CREATE (nyc)-[:ROAD_TO {distance: 790}]->(chicago)
```

![](/imgs/learn-ai/Neo4j_06.png)


现在我们拥有了一个以美国主要城市为节点、城市间公路为关系的连通网络，边的属性包含了以英里为单位的距离。让我们返回所有城市和公路来可视化这个网络：

```bash
MATCH (c1:City)-[r:ROAD_TO]->(c2:City)
WHERE id(c1) < id(c2)  // 避免重复路径
RETURN c1, r, c2
```

这将展示我们的交通网络，后续我们将用它来演示变长路径查询和最短路径算法。

## 变长关系

Neo4j 允许我们使用带有星号（`*`）的特殊语法来查询变长路径。当你不确定两个节点之间隔了多少跳时，这一功能极为强大。

### 变长关系基础语法

变长关系的基本语法如下：

```bash
-[关系类型*最小跳数..最大跳数]->
```

其中：

- `关系类型` 指定要沿何种关系类型进行遍历
- `最小跳数` 是最少跳数（可选，默认为 1）
- `最大跳数` 是最多跳数（可选）

**如果同时省略最小和最大跳数，语法变为 `*`，表示"任意跳数"**。

让我们通过一些示例来理解：

### 查找 2 条道路连接以内的城市

要查找从纽约出发，最多经过 2 条公路可以到达的所有城市：

```bash
MATCH (nyc:City {name: "New York"})-[:ROAD_TO*1..2]->(destination:City)
RETURN destination.name AS City, 
       CASE 
         WHEN destination.name = "New York" THEN "Starting Point"
         ELSE "Destination" 
       END AS Type
```

该查询将返回与纽约直接相连的城市，以及需要两条公路才能到达的城市。

### 查找城市之间所有可能的路线

要查找纽约和迈阿密之间的所有可能路线（无论长度）：

```bash
MATCH path = (:City {name: "New York"})-[:ROAD_TO*]->(:City {name: "Miami"})
RETURN path, length(path) AS hops
ORDER BY hops
LIMIT 5
```

**警告**：在大规模图上，此查询可能会非常消耗计算资源，因为它会探索所有可能的路径。通常，更好的做法是限制最大路径长度：

```bash
MATCH path = (:City {name: "New York"})-[:ROAD_TO*1..6]->(:City {name: "Miami"})
RETURN path, length(path) AS hops
ORDER BY hops
LIMIT 5
```

### 从路径中提取节点信息

在处理路径时，你可以使用 `nodes()` 和 `relationships()` 等函数来提取路径上的节点和关系信息：

```bash
MATCH path = (:City {name: "New York"})-[:ROAD_TO*1..4]->(:City {name: "Miami"})
WITH path, length(path) AS hops
ORDER BY hops
LIMIT 1
RETURN [node IN nodes(path) | node.name] AS cities,
       [rel IN relationships(path) | rel.distance] AS distances,
       reduce(total = 0, distance IN [rel IN relationships(path) | rel.distance] | total + distance) AS totalDistance
```

这个查询找到了一条从纽约到迈阿密的路线，提取了沿途城市名称、每段路的距离，以及总距离。

## 最短路径算法

虽然变长关系允许我们找到节点之间的所有路径，但在实际应用中，我们通常只想找到最高效的路线。Neo4j 为此提供了内置函数：

- `shortestPath()`：查找跳数最少的路径
- `allShortestPaths()`：查找所有跳数最少的路径

### 查找城市之间的最短路径

要查找纽约和洛杉矶之间最短路径（最少公路数）：

```bash
MATCH path = shortestPath((start:City {name: "New York"})-[:ROAD_TO*]-(end:City {name: "Los Angeles"}))
RETURN [node IN nodes(path) | node.name] AS route,
       length(path) AS hops,
       reduce(total = 0, rel IN relationships(path) | total + rel.distance) AS totalDistance
```

注意我们使用了无向关系语法 `()-[]-()`，它允许路径沿任意方向遍历关系。

### 查找所有最短路径

有时会有多条相同（最小）跳数的路径。要查找所有这样的路径：

```bash
MATCH paths = allShortestPaths((start:City {name: "Dallas"})-[:ROAD_TO*]-(end:City {name: "Chicago"}))
RETURN [node IN nodes(paths) | node.name] AS route,
       length(paths) AS hops,
       reduce(total = 0, rel IN relationships(paths) | total + rel.distance) AS totalDistance
ORDER BY totalDistance
```

此查询返回从达拉斯到芝加哥、跳数最少的所有路线，按总距离排序。

## 高级路径查找技术

现在让我们探索一些在 Neo4j 中处理路径的更高级技术。

### 按距离查找真正的最短路径

`shortestPath()` 函数查找的是跳数最少的路径，但这未必是距离最短的路径。要找到总距离最短的路径：

```bash
MATCH path = (:City {name: "New York"})-[:ROAD_TO*]->(:City {name: "Los Angeles"})
WITH path, reduce(total = 0, rel IN relationships(path) | total + rel.distance) AS totalDistance
ORDER BY totalDistance
LIMIT 1
RETURN [node IN nodes(path) | node.name] AS route,
       length(path) AS hops,
       totalDistance
```

此查询评估从纽约到洛杉矶的所有路径，返回总距离最短的那条。

### 查找一定距离范围内的城市

要查找芝加哥 500 英里范围内的所有城市（基于累积路程距离）：

```bash
MATCH path = (:City {name: "Chicago"})-[:ROAD_TO*1..10]-(destination:City)
WHERE destination.name <> "Chicago"
WITH destination, path, 
     reduce(total = 0, rel IN relationships(path) | total + rel.distance) AS totalDistance
WHERE totalDistance <= 500
RETURN DISTINCT destination.name AS city, totalDistance
ORDER BY totalDistance
```

此查询查找所有与芝加哥相连的城市，计算每条路径的总距离，并过滤出距离在 500 英里以内的。

### 查找最中心的城市

路径查找的一个有趣应用是判断网络中"最中心"的节点。我们可以使用平均最短路径长度来衡量中心性：

```bash
MATCH (origin:City)
MATCH (destination:City)
WHERE origin <> destination
MATCH path = shortestPath((origin)-[:ROAD_TO*]-(destination))
WITH origin, avg(length(path)) AS avgPathLength
RETURN origin.name AS city, avgPathLength
ORDER BY avgPathLength
LIMIT 5
```

此查询计算每个城市到所有其他城市的平均最短路径长度，从而衡量哪些城市在公路网络中处于最中心的位置。

## 实际应用

让我们探索变长关系和路径查找的一些实际应用场景。

### 旅行规划

假设我们计划一次从波士顿到旧金山的公路旅行，并希望途经特定城市：

```bash
MATCH (start:City {name: "Boston"}),
      (end:City {name: "San Francisco"}),
      (denver:City {name: "Denver"})
MATCH p1 = shortestPath((start)-[:ROAD_TO*]-(denver))
MATCH p2 = shortestPath((denver)-[:ROAD_TO*]-(end))
WITH nodes(p1) + tail(nodes(p2)) as route
RETURN [node IN route | node.name] AS roadTripCities,
       reduce(total = 0, i IN range(0, size(route)-2) |
           total + size([
               (route[i])-[r:ROAD_TO]-(route[i+1]) | r.distance
           ])
       ) AS totalDistance
```

此查询查找从波士顿到旧金山、途经丹佛的最短路线。

### 路线优化

对于配送业务来说，找到途经多个城市的最优路线至关重要：

```bash
MATCH (origin:City {name: "New York"})
WITH origin
MATCH (stop1:City {name: "Washington DC"})
MATCH (stop2:City {name: "Atlanta"})
MATCH (stop3:City {name: "Dallas"})
MATCH (destination:City {name: "Houston"})

MATCH path1 = shortestPath((origin)-[:ROAD_TO*]-(stop1))
MATCH path2 = shortestPath((stop1)-[:ROAD_TO*]-(stop2))
MATCH path3 = shortestPath((stop2)-[:ROAD_TO*]-(stop3))
MATCH path4 = shortestPath((stop3)-[:ROAD_TO*]-(destination))

WITH 
  [node IN nodes(path1) | node.name] +
  tail([node IN nodes(path2) | node.name]) +
  tail([node IN nodes(path3) | node.name]) +
  tail([node IN nodes(path4) | node.name]) AS route,

  reduce(total = 0, rel IN relationships(path1) | total + rel.distance) +
  reduce(total = 0, rel IN relationships(path2) | total + rel.distance) +
  reduce(total = 0, rel IN relationships(path3) | total + rel.distance) +
  reduce(total = 0, rel IN relationships(path4) | total + rel.distance) AS totalDistance

RETURN route, totalDistance
```

此查询为需要按指定顺序途经多个城市的配送卡车计算最高效的路线。

### 检测聚类

我们可以使用变长路径来检测网络中的聚类或社区：

```bash
// 查找彼此之间最多 2 条道路连接的城市
MATCH (c1:City)
MATCH (c2:City)
WHERE id(c1) < id(c2)
MATCH path = shortestPath((c1)-[:ROAD_TO*1..2]-(c2))
WHERE length(path) <= 2
WITH c1, collect(c2) AS nearby
RETURN c1.name AS city, [city IN nearby | city.name] AS nearbyCities, size(nearby) AS clusterSize
ORDER BY clusterSize DESC
```

此查询识别紧密相连的城市，形成自然的地理聚类。

## 最佳实践与优化技巧

在 Neo4j 中处理路径查询时，请牢记以下最佳实践：

1. **始终限制路径长度**：无边界路径长度（`*`）的查询可能会非常昂贵。尽可能使用合理的上限（`*1..n`）：

```bash
// 比直接使用 * 更好
MATCH path = (a)-[:ROAD_TO*1..10]-(b)
```

2. **尽可能使用有向关系**：如果你知道关系的方向，使用有向语法（`->` 或 `<-`）而非无向语法（`-`），因为这会缩小搜索空间：

```bash
// 若知道方向，此写法更高效
MATCH path = (a)-[:ROAD_TO*1..5]->(b)
```

3. **尽早过滤**：在查询中尽早应用过滤条件，以减少处理的数据量：

```bash
// 更好的做法
MATCH (a:City {name: "New York"}), (b:City {name: "Los Angeles"})
MATCH path = shortestPath((a)-[:ROAD_TO*]-(b))

// 效率较低的写法
MATCH path = shortestPath((a)-[:ROAD_TO*]-(b))
WHERE a.name = "New York" AND b.name = "Los Angeles"
```

4. **使用 PROFILE 分析查询性能**：在查询前加上 `PROFILE`，查看 Neo4j 如何执行查询以及瓶颈所在：

```bash
PROFILE MATCH path = shortestPath((a:City {name: "New York"})-[:ROAD_TO*]-(b:City {name: "Los Angeles"}))
RETURN path
```

5. **考虑使用 Graph Data Science 库**：对于高级路径查找算法，建议使用 Neo4j 的图数据科学库（Graph Data Science Library），它提供了许多算法的优化实现：

```bash
// 使用 GDS 的示例（需要安装 GDS 库）
CALL gds.alpha.shortestPath.stream({
    nodeQuery: 'MATCH (n:City) RETURN id(n) as id',
    relationshipQuery: 'MATCH (n)-[r:ROAD_TO]->(m) RETURN id(n) as source, id(m) as target, r.distance as weight',
    startNode: $startNodeId,
    endNode: $endNodeId,
    relationshipWeightProperty: 'weight'
})
YIELD nodeId, cost
MATCH (n) WHERE id(n) = nodeId
RETURN n.name AS city, cost
```

## 常见陷阱及避免方法

在处理变长路径和路径查找算法时，请注意以下常见陷阱：

1. **内存密集型操作**：

在密集连接的图中查找两个节点之间的所有路径可能消耗大量内存。务必先在小数据集上测试，并始终使用路径长度限制。

2. **方向性混淆**：

明确你是希望沿特定方向（`->`、`<-`）还是沿任意方向（`-`）遍历关系。

3. **忽略关系属性**：

`shortestPath()` 查找的是跳数最少的路径，而不一定是权重（如距离）最小的路径。使用自定义聚合函数来查找加权最短路径。

4. **性能问题**：

如果查询速度慢，考虑为过滤条件中使用的属性添加索引，并使用更具体的模式来缩小搜索空间。

## 总结

在本教程中，我们探索了 Neo4j 处理变长关系和在图中查找路径的强大能力。这些特性让图数据库在关联数据分析方面真正脱颖而出。

我们学习了：

- 使用 `*` 语法查询变长关系
- 使用 `shortestPath()` 查找节点之间的最短路径
- 使用 `allShortestPaths()` 发现所有最短路径
- 使用 `nodes()` 和 `relationships()` 等函数从路径中提取信息
- 使用聚合函数计算加权路径
- 将这些技术应用于解决实际问题

无论你是在构建社交网络、推荐引擎、物流系统，还是任何处理关联数据的应用程序，这些路径查找能力都将成为你 Neo4j 工具箱中不可或缺的利器。