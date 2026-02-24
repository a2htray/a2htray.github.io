---
title: "Python PAT 甲级 1003"
date: 2022-08-01T23:24:15+08:00
draft: false
reward: false
categories:
 - 算法
 - PAT
tags:
 - PAT
---

PAT 甲级 1003 。

<!-- more -->

```python
# -*- coding:utf-8 -*-
import sys

MAX_INT = sys.maxsize

if __name__ == '__main__':
    # m 城市个数
    # n 路径条数
    # start 起始城市下标
    # end 结束城市下标
    m, n, start, end = map(int, input().strip().split(' '))
    # nums_of_teams 各城市救援队的数量
    nums_of_teams = list(map(int, input().strip().split(' ')))

    assert m == len(nums_of_teams)

    # 城市间的路径矩阵
    roadmap = [
        [MAX_INT for _ in range(m)] for _ in range(m)
    ]
    # 根据输入初始化 roadmap
    while n != 0:
        i, j, d = map(int, input().strip().split(' '))
        roadmap[i][j] = d
        roadmap[j][i] = d  # 对角矩阵表示无向图
        n = n - 1

    roadmap[start][start] = 0  # 测试点中的起始城市和结束城市可能相同

    # 起始城市到其余城市最短路径条数
    nums_of_short_paths = [0 for _ in range(m)]
    # 保存访问记录, 0 表示未访问, 1 表示已访问
    visited = [0 for _ in range(m)]
    # 起始城市到其余城市最短路径
    dists = [MAX_INT for _ in range(m)]
    # 起始城市到其余城市在最短路径的基础上, 各路径上城市救援队个数的和
    weights = [0 for _ in range(m)]

    # 初始状态
    nums_of_short_paths[start] = 1  # 起始城市到自己, 表示有 1 条
    dists[start] = 0
    weights[start] = nums_of_teams[start]

    for i in range(m):
        u = -1
        min_d = MAX_INT
        # 在剩余未访问的城市中, 找到距离起始城市最近的城市
        for j in range(m):
            if visited[j] == 0 and dists[j] < min_d:
                min_d = dists[j]
                u = j
        if u == -1:
            break

        visited[u] = 1  # 第 1 次找到的必然是起始城市
        for k in range(m):
            # for 循环用于计算：
            # > 已访问城市到未访问城市的最短路径之和
            if visited[k] == 0 and roadmap[u][k] != MAX_INT:
                # 已访问城市与未访问城市之间必须具有连通性
                if dists[u] + roadmap[u][k] < dists[k]:
                    dists[k] = dists[u] + roadmap[u][k]
                    nums_of_short_paths[k] = nums_of_short_paths[u]
                    weights[k] = weights[u] + nums_of_teams[k]
                elif dists[u] + roadmap[u][k] == dists[k]:
                    nums_of_short_paths[k] = nums_of_short_paths[k] + nums_of_short_paths[u]
                    if nums_of_teams[k] + weights[u] > weights[k]:
                        weights[k] = weights[u] + nums_of_teams[k]

    print(f'{nums_of_short_paths[end]} {weights[end]}', end='')
```

