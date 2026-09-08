+++
date = '2026-09-07T14:36:48+08:00'
draft = false
title = '给 PostgreSQL 装上"向量引擎"：pgvector'
categories = ['后端技术', 'PostgreSQL']
tags = ['PostgreSQL', 'pgvector', '数据库', '向量数据库']
toc = true
+++

![](/imgs/learn-db/ScreenShot_2026-09-08_112250_865.png)

在大模型时代，应用须提供语义搜索或知识检索功能，如果后端数据库本身使用的就是 PG，则可以直接使用 pgvector 来实现，无需额外的数据库组件。

## 什么是 pgvector

**pgvector 是 PostgreSQL 的一个开源扩展，它让 PG 能够原生存储、索引和检索高维向量（embedding）。** 它把"向量相似度搜索"这一原本需要专门向量数据库的能力，直接下沉到了你熟悉的关系型数据库里。**常用于语义搜索、RAG 检索、推荐系统等 AI 场景**，让应用用一条 SQL 就能做"找最相似的 N 条数据"。

## 如何使用 pgvector

下面用 Docker 起一个 PG，再启用扩展，这是最省心的落地方式：

直接用社区维护好的预装镜像 `ankane/pgvector`：

```bash
$ docker run -d \
  --name pgvector \
  -e POSTGRES_PASSWORD=yourpassword \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_DB=testdb \
  -p 5432:5432 \
  ankane/pgvector:latest
```

**只需一条 SQL 启用扩展**（每个数据库执行一次），在宿主机上执行：

```bash
$ docker exec -t pgvector psql -U postgres -d testdb -c "CREATE EXTENSION vector;" 
CREATE EXTENSION
```

验证是否启用成功：

```bash
$ docker exec -t pgvector psql -U postgres -d testdb -c "SELECT extname, extversion FROM pg_extension WHERE extname = 'vector';"
extname | extversion 
---------+------------
 vector  | 0.5.1
(1 row)
```

> 如果是使用原生 PostgreSQL 镜像或主机安装 PG 服务，则需要先编译安装扩展源码，或在宿主机用 `apt`/`yum` 安装 `postgresql-<版本>-pgvector` 包，再执行上面的 `CREATE EXTENSION`。

## pgvector 的优势

## 3. pgvector 的优势

pgvector 优势在于它**最省事、最契合已有架构**，如果现有应用中已使用到了 PG，则只需要安装、启用 pgvector 扩展即可，没有技术负担。

| 工程维度       | pgvector 的优势        | 说明                                                    |
| ---------- | ------------------- | ----------------------------------------------------- |
| **基础设施复用** | 零新增组件               | 备份、监控、权限、连接池全部沿用现有 PG 体系                   |
| **混合检索**   | 极强                  | 向量 + 结构化字段可同一条 SQL 做 `JOIN` / `WHERE` 过滤，无需应用层二次查询再合并 |
| **事务一致性**  | ACID                | 向量写入与业务数据落同一事务，不会出现"业务成功但向量丢了"                        |
| **渐进式采用**  | 低门槛                 | SQL 开发者即可上手，无需学新查询语言                                  |
| **成本**     | 少一套系统               | 少一份独立系统的人力与运维成本                              |
| **生态成熟**   | 即插即用                | ORM、BI、ETL、CDC 工具全部围绕 PG 生态                           |
| **索引多样**   | HNSW + IVFFlat + 量化 | HNSW 高召回、IVFFlat 省内存、halfvec/bit 量化省空间                |
| **多租户**    | 天然支持                | 配合 PG Schema 或行级安全（RLS）做租户隔离，向量自然跟随                   |

> **客观短板**：当向量规模进入**数千万乃至亿级**，且要求极高并发 QPS / 极低延迟时，pgvector 的检索延迟会明显上升，且 PG 本身水平扩展（分片）较麻烦。这种场景才需要 Milvus / Qdrant 等专用库。

## 插入与检索

**建表**

```bash
$ docker exec -t pgvector psql -U postgres -d testdb -c "CREATE TABLE IF NOT EXISTS documents (id SERIAL PRIMARY KEY, content TEXT, embedding VECTOR(4096));"
```

**插入数据**

```python
import psycopg2

from langchain_ollama.embeddings import OllamaEmbeddings

PG_URL = "postgresql://postgres:yourpassword@localhost:5432/testdb"
conn = psycopg2.connect(PG_URL)

embeddings = OllamaEmbeddings(model="qwen3-embedding:8b")

knowledges = (
    '甘蓝型油菜是国内主栽油菜类型，育种重点兼顾丰产性、抗逆性与菜籽品质，低芥酸低硫苷是优质油菜核心品质指标。',
    '油菜分子标记辅助育种可快速筛选目标基因材料，缩短育种周期，减少大规模田间表型鉴定的人力成本。',
    '油菜 DUS 测试包含特异性、一致性、稳定性三项考核，是油菜新品种申请品种权必须完成的鉴定流程。',
    '油菜菌核病属于高发真菌病害，育种上主要通过种质资源挖掘抗性亲本，聚合多个抗性位点提升品种耐病水平。',
    '油菜杂种优势利用主要依靠细胞质雄性不育、细胞核雄性不育系统，筛选优良恢复系与保持系是杂交育种关键工作。',
    '油菜种质资源库保存野生、地方品种与突变材料，通过表型和基因型鉴定，发掘耐冻、耐渍、抗病等优良等位基因。',
    '油菜早熟品种育种需要协调开花期、终花期与成熟期，解决早熟条件下产量容易下降、千粒重降低的矛盾。',
    '油菜抗倒伏性状受茎秆强度、根系发育、株高分枝等多基因控制，表型鉴定需要结合田间大风环境开展多年多点试验。',
    '油菜品质育种除芥酸、硫苷之外，同时关注含油量提升，高油育种需要平衡油脂合成通路与植株整体生长发育。',
)

with conn.cursor() as cursor:
    for idx, knowledge in enumerate(knowledges, start=1):
        embedding = embeddings.embed_query(knowledge)
        cursor.execute("INSERT INTO documents (content, embedding) VALUES (%s, %s::vector)", (knowledge, str(embedding)))
        conn.commit()
        print(f"已插入第 {idx} 条知识")
```

执行输出：

```bash
已插入第 1 条知识
已插入第 2 条知识
已插入第 3 条知识
已插入第 4 条知识
已插入第 5 条知识
已插入第 6 条知识
已插入第 7 条知识
已插入第 8 条知识
已插入第 9 条知识
```

**知识检索**

```python
import psycopg2

from langchain_ollama.embeddings import OllamaEmbeddings

PG_URL = "postgresql://postgres:yourpassword@localhost:5432/testdb"
conn = psycopg2.connect(PG_URL)

embeddings = OllamaEmbeddings(model="qwen3-embedding:8b")

queries = (
    '甘蓝型油菜优质品种主要关注哪些品质指标？',
    '油菜 DUS 测试具体考核哪些内容，有什么用途？',
    '油菜菌核病在育种工作中一般如何开展抗性改良？',
)

with conn.cursor() as cursor:
    for idx, query in enumerate(queries, start=1):
        embedding = embeddings.embed_query(query)
        cursor.execute("""SELECT
id, content, embedding <=> %s::vector AS distance
FROM documents
ORDER BY embedding <=> %s::vector
LIMIT 3;""", (str(embedding), str(embedding)))
        
        results = cursor.fetchall()
        print(f"查询 {idx}：{query}，知识排序结果：")
        for result in results:
            print(f" - 序号 {result[0]}，距离 {result[2]:.4f}，内容 {result[1]}")
```

执行输出：

```bash
查询 1：甘蓝型油菜优质品种主要关注哪些品质指标？，知识排序结果：
 - 序号 1，距离 0.1529，内容 甘蓝型油菜是国内主栽油菜类型，育种重点兼顾丰产性、抗逆性与菜籽品质，低芥酸低硫苷是优质油菜核心品质指标。
 - 序号 9，距离 0.2892，内容 油菜品质育种除芥酸、硫苷之外，同时关注含油量提升，高油育种需要平衡油脂合成通路与植株整体生长发育。
 - 序号 7，距离 0.3519，内容 油菜早熟品种育种需要协调开花期、终花期与成熟期，解决早熟条件下产量容易下降、千粒重降低的矛盾。
查询 2：油菜 DUS 测试具体考核哪些内容，有什么用途？，知识排序结果：
 - 序号 3，距离 0.0730，内容 油菜 DUS 测试包含特异性、一致性、稳定性三项考核，是油菜新品种申请品种权必须完成的鉴定流程。
 - 序号 5，距离 0.4137，内容 油菜杂种优势利用主要依靠细胞质雄性不育、细胞核雄性不育系统，筛选优良恢复系与保持系是杂交育种关键工作。
 - 序号 8，距离 0.4327，内容 油菜抗倒伏性状受茎秆强度、根系发育、株高分枝等多基因控制，表型鉴定需要结合田间大风环境开展多年多点试验。
查询 3：油菜菌核病在育种工作中一般如何开展抗性改良？，知识排序结果：
 - 序号 4，距离 0.0986，内容 油菜菌核病属于高发真菌病害，育种上主要通过种质资源挖掘抗性亲本，聚合多个抗性位点提升品种耐病水平。
 - 序号 5，距离 0.3349，内容 油菜杂种优势利用主要依靠细胞质雄性不育、细胞核雄性不育系统，筛选优良恢复系与保持系是杂交育种关键工作。
 - 序号 8，距离 0.3563，内容 油菜抗倒伏性状受茎秆强度、根系发育、株高分枝等多基因控制，表型鉴定需要结合田间大风环境开展多年多点试验。
```

### 距离运算符

* `<->`：L2 欧氏距离，`ORDER BY col <-> vec`，值越小越相似
* `<#>`：负内积，`ORDER BY col <#> vec`，值越小越相似
* `<=>`：余弦距离，`ORDER BY col <=> vec`，值越小越相似，范围[0,2]

## 对比向量数据库

把 pgvector 和主流专用向量库放在一起看，结论很清楚：**没有"最强"，只有"最合适"。**

| 维度        | **pgvector** | Milvus  | Qdrant        | Chroma | Weaviate |
| --------- | ------------ | ------- | ------------- | ------ | -------- |
| **形态**    | PG 扩展        | 分布式向量库  | Rust 向量服务     | 嵌入式/轻量 | 向量 + 图   |
| **适合规模**  | 百万 ~ 数千万     | **十亿级** | 百万 ~ 数亿       | < 百万   | 百万 ~ 数亿  |
| **混合查询**  | **极强（SQL）**  | 中       | 强（payload 过滤） | 弱      | 极强       |
| **事务一致性** | **ACID**     | 弱       | 弱             | 弱      | 弱        |
| **运维复杂度** | **低**（复用 PG） | 高（K8s）  | 中             | 极低     | 中        |
| **学习曲线**  | **低（SQL）**   | 陡       | 中             | 极低     | 中        |
| **水平扩展**  | 难            | 原生      | 分片            | 无      | 支持       |

**关键结论**：

- **pgvector** 赢在"不引入新系统 + 强一致性 + 混合检索"，适合已有 PG 的中小规模场景。
- **Milvus** 赢在"十亿级规模 + GPU 加速"，但运维重。
- **Qdrant** 赢在"中大规模 + 云原生 + 实时增量更新 + 过滤检索"。
- **Chroma** 赢在"本地开发一把梭"，但基本不上生产。
- **Weaviate** 赢在"结构化 + 向量混合搜索 / 多模态"。

> 一个常被引用的基准（50M 向量、768 维）：配合 [pgvectorscale](https://github.com/timescale/pgvectorscale) 扩展后，pgvector 在 QPS 与延迟上已能和专用库正面对抗。**但规模一旦上亿，专用库仍是首选。**

