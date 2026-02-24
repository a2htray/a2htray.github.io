---
title: "Python PAT 甲级 1003"
date: 2022-08-03T00:33:33+08:00
draft: false
reward: false
categories:
  - 算法
  - PAT
tags:
 - PAT
---

PAT 甲级 1003 。

```python
# -*- coding:utf-8 -*-

if __name__ == '__main__':
    # n 树中节点个数
    # m 非叶子节点个数
    n, m = list(map(int, input().strip().split(' ')))
    # 二维数组
    # 元素的下标表示节点的 ID
    # 第 1 个元素不使用
    tree = [[] for _ in range(n+1)]
    while m != 0:
        tokens = list(map(int, input().split(' ')))
        _id = tokens[0]
        children = tokens[2:]
        tree[_id] = children
        m = m - 1

    # 记录每一层叶子节点的数量
    counts = [0 for _ in range(n+1)]
    # 记录最大深度
    max_depth = 0

    def dfs(i, depth):
        global max_depth, counts, tree

        if len(tree[i]) == 0:  # 节点 i 为叶子节点
            counts[depth] = counts[depth] + 1
            max_depth = max(depth, max_depth)
            return
        for j in tree[i]:
            dfs(j, depth + 1)

    dfs(1, 0)

    print(' '.join([str(count) for count in counts[:max_depth+1]]))

```

测试点 4 没过，应该是没有理解 *The input ends with N being 0. That case must NOT be processed.* 的原因，因为不知道要输出什么。

![](WX20220803-003514@2x.png)