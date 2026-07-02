+++
date = '2026-07-02T15:18:40+08:00'
draft = false
title = 'COF - Connectivity-based Outlier Factor 连通离群因子算法'
categories = ['算法', '异常检测']
tags = ['异常检测', 'Connectivity-based Outlier Factor', '算法', 'COF']
toc = true
+++

Connectivity-based Outlier Factor，连通离群因子算法是 2002 年提出的改进 LOF 算法的无监督异常检测算法。LOF 算法的局限在于：它是基于“**球形局部密度**”假设，对带状、曲线、链式、流形连续分布的数据缺陷极大，会把长条簇两端的正常样本误判为离群点。

COF 算法引入**连通最短路径**衡量<font color="blue">样本融入局部簇的难易程度</font>，适配连续带状分布。如下图：

![](/imgs/learn-algo/cof/01_lof_vs_cof.png)

## 数据集介绍

构建的数据集是月牙形分布，其散点图如下：

![](/imgs/learn-algo/cof/00_dataset_preview.png)

对该数据集应用 LOF 算法的结果如下图所示：

![](/imgs/learn-algo/cof/02_lof_run.png)

LOF 算法成功检测出所有的异常样本，但误报了弧端点的样本。

## 概念

1. SBN Path - Set Based Neighbor Path 集合近邻连通路径

SBN Path，集合近邻连通路径是指样本 $p$ 在其邻域内，连接所有邻居的最短路径。SBN Path 是单条有序链式路径，类似于 $p\ \text{->}\ v_1\ \text{->}\ v_2\ \text{->}\ v_3$。如下图：

![](/imgs/learn-algo/cof/03_sbn_path.png)

SBN Path 为 $A\ \text{->}\ B\ \text{->}\ C\ \text{->}\ D$。

2. $\operatorname{ac-}dist_k(p)$  平均连通距离

$\operatorname{ac-}dist_k(p)$ 是对 SBN Path 做递减加权求和，公式如下：

$$\operatorname{ac-}dist_k(p)=\sum_{i=1}^n\frac{2(n+1-i)}{n(n+1)}\cdot{e_i}$$

其中 $n = {\left|N_k(p)\right|}$，权重 $w_i=\frac{2(n+1-i)}{n(n+1)}$ 满足 $\sum_{i=1}^nw_i=1$。从权重公式可知：

* 越近邻居的距离权重越大
* 最远邻居的距离权重最小
* 越往后的距离权重线性衰减

$\operatorname{ac-}dist_k(p)$ 越大，表示样本 $p$ 与邻居之间越紧密，反之则越疏远。

3. COF 打分公式

COF 是相对比值指标，将样本 $p$ 自身的连通距离，与它 k 个近邻的平均连通距离做对比，公式如下：

$$\text{COF}_k(p) = \frac{{\left|N_k(p)\right|} \cdot \text{ac-dist}_k(p)}{\displaystyle\sum_{o \in N_k(p)} \text{ac-dist}_k(o)}$$

其数据用于样本判断：

* $\text{COF}_k(p)\approx 1$，正常样本
* $\text{COF}_k(p)\gg 1$，异常样本
* $\text{COF}_k(p)\lt 1$，样本更贴近簇心

## 计算过程

1. 对每个样本 $p$，计算距离，找出 $N_k(p)$
2. 构造样本 $p$ 和 $N_k(p)$ 的 SBN Path
3. 代入加权公式，计算 $\operatorname{ac-}dist_k(p)$
4. 遍历所有邻居 $o\in N_k(p)$，求和邻域 $\operatorname{ac-}dist$
5. 计算 $\text{COF}_k(p)$

运用 COF 算法的异常检测结果如下图：

![](/imgs/learn-algo/cof/04_cof_run.png)

COF 算法同样检出了所有异常样本，但也对样本 S020 和 S015 进行了误报。

## 心得

* 在示例数据集中，LOF 和 COF 的性能差别不大，LOF 倾向将弧端的样本报为异常，而 COF 倾向将弧边的样本报为异常
* 如果按异常分数阈值过滤（取前 5），在示例数据值上体现，LOF 还更好，因为 COF 将样本 S015 的异常分数排到了前 5