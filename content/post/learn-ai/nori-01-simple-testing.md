+++
date = '2026-09-14T14:04:20+08:00'
draft = false
title = 'Nori：面向表格回归的上下文学习模型'
categories = ['人工智能', '结构化大模型']
tags = ['Nori', '结构化大模型', '表格数据']
toc = true
+++

Nori 是一系列面向**表格数据回归**的轻量级模型，由 Synthefy 推出。它的核心思路是**上下文学习（In-Context Learning）**：无需传统意义上的参数训练，只要把一小段“已知样本（上下文）”喂给模型，它就能对新的查询行直接给出预测。

## 安装依赖

```bash
$ pip install synthefy-nori torch pandas numpy
```

## 快速开始 - 房价预测

```python
import os
os.environ["HF_ENDPOINT"] = "https://hf-mirror.com"
os.environ["HF_XET_HIGH_PERFORMANCE"] = "1"

import pandas as pd
from synthefy_nori import NoriRegressor

context = pd.DataFrame({
    "sqft": [1200, 1500, 1800, 2100],
    "bedrooms": [2, 3, 3, 4],
    "house_age": [20, 15, 10, 5],
    "price": [240000, 300000, 360000, 420000]
})

query = pd.DataFrame({
    "sqft": [1600],
    "bedrooms": [3],
    "house_age": [12]
})

X_context = context.drop("price", axis=1)
y_context = context["price"]

model = NoriRegressor(model="nori-30m")
model.fit(X_context, y_context)

preds = model.predict(query)
print("预测房价：", preds)
```

运行结果：

```bash
预测房价： [329112.99150825]
```

## 工作原理

Nori 在推理时大致经历以下步骤：

1. **注入上下文**：`fit(X_context, y_context)` 把若干已知 `(特征, 标签)` 样本存入模型
2. **行注意力检索**：对每条查询行，模型通过**行注意力（row attention）**在上下文样本中
   找到与之最相似的若干行，作为预测的参考依据
3. **插值 / 推理预测**：`predict(query)` 基于检索到的相似行做加权插值，或直接前向输出预测值，
   全程为**单轮前向**，无反向传播

## Nori-30M（普通版） vs Nori-30M-thinking（Thinking 增强版）

两者**参数量都是 30M，基础主干架构完全一致**，只是预训练目标与推理行为不同。

| 维度 | Nori-30M（普通版） | Nori-30M-thinking（增强版） |
| --- | --- | --- |
| 设计定位 | 基础版本，快速直接预测 | 增强数值推理版本 |
| 擅长场景 | 线性/弱非线性关系的一般回归 | 非线性、多步数值关系 |
| 推理逻辑 | 单轮前向，基于行注意力做插值预测 | 主干外增加**数值推理预训练分支**，建模特征间数学关系 |
| 推理速度 | 更快 | 因额外分支略慢 |
| 显存开销 | 更低 | 因额外分支略高 |
| 预测精度 | 一般场景足够 | 复杂数值关系上更准 |

简单总结三点差异：

1. **设计定位**：Nori-30M 是基础版本，可用于快速直接预测；Nori-30M-thinking 是增强数值推理版本，专门优化非线性、多步数值关系。
2. **推理逻辑**：Nori-30M 单轮前向，直接基于行注意力做插值预测；Nori-30M-thinking 在主干之外增加数值推理预训练分支，建模特征间数学关系。
3. **速度与资源**：Nori-30M 性能更优（更快、更省显存）；Nori-30M-thinking 由于额外的数值推理分支，显存开销略高、推理略慢，但换来复杂数值关系上更高的精度。

> **本地推理支持说明**：Thinking 增强版仅在 hosted Synthefy API 上运行，本地推理不支持。

## 适用场景与选型建议

* **选 Nori-30M（普通版）**：样本量小、特征关系近似线性、对延迟/显存敏感、追求“即用即走”
* **选 Nori-30M-thinking（增强版）**：特征间存在明显非线性或需要多步数值推导、可接受稍高延迟与显存、追求更高精度

## 常见问题（FAQ）

**Q1：模型下载到哪里了？**

默认位于 Hugging Face 缓存目录：`~/.cache/huggingface/hub/models--Synthefy--Nori-30M`，首次运行自动下载，后续复用。

**Q2：fit 会训练模型参数吗？**

不会。Nori 是上下文学习模型，`fit` 仅把上下文样本注入模型，真正的计算发生在 `predict` 的单次前向。

**Q3：上下文样本越多越好吗？**

通常更多上下文能提升参考质量，但也会增大输入规模、增加推理开销。建议结合实际数据做小范围调参。

**Q4：没有 GPU 能跑吗？**

可以。模型会在 CPU 上推理。