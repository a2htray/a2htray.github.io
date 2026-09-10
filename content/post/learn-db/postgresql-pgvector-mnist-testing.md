+++
date = '2026-09-09T09:12:48+08:00'
draft = false
title = 'PostgreSQL：pgvector 扩展存储与检索手写数字图像（MNIST 数据集）'
categories = ['后端技术', 'PostgreSQL']
tags = ['PostgreSQL', 'pgvector', '数据库', '向量数据库', 'hafvec', 'MNIST 数据集']
toc = true
+++

![](/imgs/learn-db/ScreenShot_2026-09-10_164450_760.png)

本文是关于 PostgreSQL pgvector 扩展的一个示例，使用 hafvec 向量类型存储和检索手写数字图像（MNIST 数据集）。

目标：向量化存储 60,000 张训练图像，并用 10,000 张测试图像进行检索，统计检索准确率。

## MNIST 数据集

MNIST 数据集是一个用于训练和测试机器学习模型的数字图像数据集。它由 60,000 张训练图像和 10,000 张测试图像组成，每个图像都是 28x28 像素的灰度图像。

一张图本质上就是一个 28×28 的二维矩阵，把每行依次拼接（展平 / flatten）后，就得到一个 784 维的向量 —— 这正是它能直接套用向量检索的天然前提。

下载数据集并预览前 3 张：

```python
from torchvision import datasets
import matplotlib.pyplot as plt

plt.rcParams["font.sans-serif"] = ["Arial Unicode MS"]
plt.rcParams["axes.unicode_minus"] = False

mnist = datasets.MNIST(root="./data", train=True, download=True)

# 可视化训练集前 3 张,一行三列子图
fig, axes = plt.subplots(1, 3, figsize=(9, 3))
for i, ax in enumerate(axes):
    image, label = mnist[i]
    ax.imshow(image, cmap="gray")
    ax.set_title(f"index={i+1}, label={label}")
    ax.axis("off")

fig.suptitle("MNIST 训练集前 3 张")
fig.tight_layout()
fig.savefig("mnist_preview.png")
```

![](/imgs/learn-db/mnist_preview.png)

## 创建数据库表

```sql
CREATE TABLE mnist_images (
    id SERIAL PRIMARY KEY, -- 序号
    label SMALLINT, -- 真实数字
    pixels vector(784), -- 28×28 展平、原始的像素向量
    embedding vector(784) -- 28×28 展平、归一化后的像素向量
);
```

> 对比检索原始像素向量和归一化后的像素向量的准确率。

## 插入数据

```python
import math

import numpy as np
import psycopg2
from pgvector.psycopg2 import register_vector
from psycopg2.extras import execute_values
from torchvision import datasets

# 连接数据库
PG_URL = "postgresql://postgres:yourpassword@localhost:5432/testdb"
conn = psycopg2.connect(PG_URL)
# 让 psycopg2 能识别并写入 vector 类型的字段
register_vector(conn)

mnist = datasets.MNIST(root="./data", train=True, download=True)

INSERT_SQL = """
INSERT INTO mnist_images (label, pixels, embedding)
VALUES %s
"""


def l2_normalize(vec):
    """对展平向量做 L2 归一化;零向量直接返回避免除零。"""
    norm = math.sqrt(sum(v * v for v in vec))
    if norm == 0.0:
        return vec
    return [v / norm for v in vec]


def main():
    rows = []
    for idx, (image, label) in enumerate(mnist):
        arr = np.array(image)
        pixels = arr.flatten() # 原始像素向量
        embedding = l2_normalize(pixels) # 归一化后的像素向量
        rows.append((int(label), pixels, embedding))

    with conn.cursor() as cur:
        execute_values(cur, INSERT_SQL, rows, page_size=1000)
    conn.commit()
    print(f"已插入 {len(rows)} 条记录")


if __name__ == "__main__":
    try:
        main()
    finally:
        conn.close()
```

![](/imgs/learn-db/ScreenShot_2026-09-09_141657_111.png)

## 检索数据

距离公式：采用欧式距离。

统计指标：

* Top-1 命中率：第 1 个直接命中的图像准确度率
* Top-5 命中率：前 5 个直接命中的图像准确度率
* 平均最近邻距离：测试图像与检索出来的第 1 个图像的向量距离的平均值，用于评估检索结果的相似度

```python
import math
import os
import threading
from concurrent.futures import ThreadPoolExecutor, as_completed

import numpy as np
import psycopg2
from pgvector.psycopg2 import register_vector
from torchvision import datasets

PG_URL = "postgresql://postgres:yourpassword@localhost:5432/testdb"

# 每个线程一份独立的数据库连接 (psycopg2 连接/游标线程不安全)
_local = threading.local()
_conn_lock = threading.Lock()
_all_conns = []


def get_conn():
    conn = getattr(_local, "conn", None)
    if conn is None or conn.closed:
        conn = psycopg2.connect(PG_URL)
        # 让 psycopg2 能识别 vector 类型
        register_vector(conn)
        _local.conn = conn
        with _conn_lock:
            _all_conns.append(conn)
    return conn


# 使用手写数字的测试集
mnist_test = datasets.MNIST(root="./data", train=False, download=True)


def l2_normalize(vec):
    norm = math.sqrt(sum(v * v for v in vec))
    return vec if norm == 0.0 else [v / norm for v in vec]


def topk_by_column(cur, query_vec, column, k=5):
    """在指定列上做欧氏距离最近邻检索,返回 [(label, distance), ...]。"""
    # <-> 欧氏距离运算符
    cur.execute(
        f"""
        SELECT label, {column} <-> %s::vector AS distance
        FROM mnist_images
        ORDER BY distance
        LIMIT %s
        """,
        (query_vec, k),
    )
    return cur.fetchall()


def process_one(idx, k):
    """处理单条样本,返回 6 个统计量: (raw_top1, raw_top5, raw_dist,
    emb_top1, emb_top5, emb_dist)。每个线程使用各自的连接与游标。"""
    cur = get_conn().cursor()

    try:
        image, true_label = mnist_test[idx]
        # 0~255 原始像素
        arr = np.array(image).flatten().tolist()
        # 归一化向量
        emb = l2_normalize(arr)

        # 1) 基于原始向量 pixels
        raw_res = topk_by_column(cur, arr, "pixels", k)
        raw_top1 = int(raw_res[0][0] == true_label)
        raw_top5 = int(any(lbl == true_label for lbl, _ in raw_res))
        raw_dist = raw_res[0][1]

        # 2) 基于归一化向量 embedding
        emb_res = topk_by_column(cur, emb, "embedding", k)
        emb_top1 = int(emb_res[0][0] == true_label)
        emb_top5 = int(any(lbl == true_label for lbl, _ in emb_res))
        emb_dist = emb_res[0][1]
    finally:
        cur.close()

    return raw_top1, raw_top5, raw_dist, emb_top1, emb_top5, emb_dist


def main():
    k = 5
    n = len(mnist_test)
    # 控制并发线程数,避免压垮数据库 (也可用环境变量 OMP_THREADS 指定)
    max_workers = int(os.environ.get("QUERY_WORKERS", os.cpu_count() or 4))

    # 统计指标
    raw_top1_hit = raw_top5_hit = 0
    emb_top1_hit = emb_top5_hit = 0
    raw_dist_sum = emb_dist_sum = 0.0

    with ThreadPoolExecutor(max_workers=max_workers) as executor:
        futures = [executor.submit(process_one, idx, k) for idx in range(n)]
        for fut in as_completed(futures):
            r1, r5, rd, e1, e5, ed = fut.result()
            raw_top1_hit += r1
            raw_top5_hit += r5
            raw_dist_sum += rd
            emb_top1_hit += e1
            emb_top5_hit += e5
            emb_dist_sum += ed

    print(f"测试集样本数: {n}, Top-k = {k}")
    print("-" * 50)
    print("基于原始向量 pixels (欧氏距离 <->):")
    print(f"  Top-1 命中率: {raw_top1_hit / n:.4f} ({raw_top1_hit}/{n})")
    print(f"  Top-5 命中率: {raw_top5_hit / n:.4f} ({raw_top5_hit}/{n})")
    print(f"  平均最近邻距离: {raw_dist_sum / n:.4f}")
    print("-" * 50)
    print("基于归一化向量 embedding (欧氏距离 <->):")
    print(f"  Top-1 命中率: {emb_top1_hit / n:.4f} ({emb_top1_hit}/{n})")
    print(f"  Top-5 命中率: {emb_top5_hit / n:.4f} ({emb_top5_hit}/{n})")
    print(f"  平均最近邻距离: {emb_dist_sum / n:.4f}")


def _close_all_conns():
    with _conn_lock:
        conns = list(_all_conns)
        _all_conns.clear()
    for c in conns:
        try:
            c.close()
        except Exception:
            pass


if __name__ == "__main__":
    try:
        main()
    finally:
        _close_all_conns()
```

> 不启用多线程处理的话，要耗时 36 分钟左右，故当前版本采用多线程的方式进行检索与统计，代码编写相对复杂一些。

执行输出如下：

```bash
测试集样本数: 10000, Top-k = 5
--------------------------------------------------
基于原始向量 pixels (欧氏距离 <->):
  Top-1 命中率: 0.9691 (9691/10000)
  Top-5 命中率: 0.9914 (9914/10000)
  平均最近邻距离: 1075.5468
--------------------------------------------------
基于归一化向量 embedding (欧氏距离 <->):
  Top-1 命中率: 0.9723 (9723/10000)
  Top-5 命中率: 0.9920 (9920/10000)
  平均最近邻距离: 0.4436
```

## 小结

从检索的统计指标可以看出：

* 在 Top-1 命中率上，基于归一化向量检索的准确度略高于基于原始向量检索，相差 0.0032
* 在 Top-5 命中率上，两者都能达到 0.99，侧面说明采用采用原始向量与归一化向量差别不大
* 在平均最近邻距离上，显示基于归一化向量检索的平均距离能更好的量化检索结果的相似度