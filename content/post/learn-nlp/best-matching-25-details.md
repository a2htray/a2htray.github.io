+++
date = '2026-09-14T17:20:12+08:00'
draft = false
title = '基于词袋的概率检索模型 BM25 详解'
categories = ['人工智能', '自然语言处理']
tags = ['自然语言处理', '文本检索', 'BM25', '概率检索模型']
toc = true
+++

![](/imgs/learn-nlp/ScreenShot_2026-09-15_163226_605.png)

BM25（Best Matching 25）是**基于词袋的概率检索模型**，目前是搜索引擎、ES/OpenSearch 全文检索默认的相关性打分算法，用来计算**查询词 Q 和文档 D 的相关性分数**。

> 核心思想：文档 D 与查询词 Q 的相关性分数越高，则表示越相似（相关）。

## 前置基础 - TF 和 IDF

**TF**（Term Frequency）：词频，即词 t 在文档 D 中出现的次数。

$$
tf(t, D) = count(t, D)
$$

* $count(t, D)$：词 t 在文档 D 中出现的次数。

**IDF**（Inverse Document Frequency）：逆文档频率，衡量词稀有度。词越罕见，IDF 越大。


$$
idf(t) = ln(\frac{N - n_t + 0.5}{n_t + 0.5} + 1)
$$

* $N$：全部文档总数
* $n_t$：包含 t 的文档数
* $0.5$：平滑系数，防止分母为 0

## BM25 公式

$$
bm25(Q, D) = \sum_{t \in Q} (idf(t) \cdot \frac{tf(t, D) \cdot (k_1 + 1)}{tf(t, D) + k_1 \cdot (1 - b + b \cdot \frac{\left | D \right | }{avgdl})})
$$

* $Q$：查询词集，包含多个查询词 t
* $\left | D \right | $：文档 D 的词数
* $avgdl$：全部文档的平均词数
* $k_1$：词频饱和系数，用于调整词频对相关性的影响，$k_1$ 越大，词频提升对分数影响越强
* $b$：文档长度因子，用于调整文档长度对相关性的影响，$b$ 越大，文档长度对分数影响越强

> $k_1$ 在 Elasticsearch 中默认值为 1.2，$b$ 默认值为 0.75。

## 代码示例

**实现 tf 和 idf 函数**

```python
import math


def tf(t, d):
    """计算词频 (term frequency)：词 t 在文档 d 中出现的次数。

    d 可为分词后的列表/元组（逐词计数），或 {词: 计数} 的映射（直接取值）。
    """
    if isinstance(d, dict):
        return d.get(t, 0)
    return sum(1 for w in d if w == t)


def idf(nt, n):
    """计算 BM25 的逆文档频率 (inverse document frequency)。

    nt: 包含词 t 的文档数
    n: 文档总数。
    """
    return math.log((n - nt + 0.5) / (nt + 0.5) + 1.0)
```

**实现简单分词并处理 5 个模拟文档**

```python
import string

import jieba as jb

punctuation = set(string.punctuation + "，。、；：？！‘’“”（）《》【】……——")

# 模拟文档库：5 条关于 AI 的简短描述
documents = (
    "人工智能（AI）是让机器模拟人类智能的技术，涵盖机器学习、自然语言处理与计算机视觉等领域。它正逐步渗透到医疗、金融与交通等行业，成为推动数字化转型的重要引擎。",
    "大语言模型通过海量文本预训练，能够理解并生成自然语言，是近年来人工智能最热门的方向之一。以 GPT 为代表的生成式模型，已能胜任写作、编程与问答等多种任务。",
    "机器学习是人工智能的子集，它让计算机从数据中自动学习规律，而无需人工显式编写每一条规则。常见算法包括决策树、支持向量机与神经网络等，各有其适用场景。",
    "计算机视觉致力于让机器看懂图像与视频，已广泛应用于人脸识别、自动驾驶和医疗影像诊断。借助深度学习，模型在目标检测与图像生成上的表现已接近甚至超越人类。",
    "强化学习通过与环境交互、以奖励信号指导决策，曾在围棋等复杂游戏中战胜人类顶尖选手。AlphaGo 正是其代表作，展示了智能体在庞大状态空间中自主规划的能力。",
)

def simple_cut(text) -> list:
    """简单分词，使用jieba分词器"""
    STOP_WORDS = {
        "的", "了", "是", "在", "我", "有", "和", "就", "都", "而", "及", "与",
        "也", "很", "还", "个", "之", "对", "这", "那", "该", "些", "等", "可以",
        "应该", "能够", "会", "要", "把", "被", "让", "从", "向", "给", "到", "上", 
        "下", "既", "也"
    }
    return [w for w in jb.cut(text) if w not in STOP_WORDS and w.strip() and w not in punctuation]


class TokenizedDocument:
    def __init__(self, text: str):
        self.text = text
        self.tokens = simple_cut(text)
        self.token_freq = self._compute_token_freq()

    def _compute_token_freq(self):
        return {token: self.tokens.count(token) for token in self.tokens}

    def token_count(self):
        return len(self.tokens)

    def tf(self, token):
        return tf(token, self.token_freq)

    def contains(self, token):
        return token in self.token_freq


tokenized_documents = [TokenizedDocument(doc) for doc in documents]
```

* `simple_cut`：简单分词，使用 jieba 分词器
* `punctuation`：标点符号集合，用于过滤掉标点符号
* `TokenizedDocument`：分词后的文档，包含原始文本和分词后的词列表
* `_compute_token_freq`：计算每个词的文档频率，用于计算 IDF
* `token_count`：返回文档中词的数量
* `tf`：计算词频 (term frequency)
* `contains`：检查文档是否包含指定词


以下是分词结果：

```text
原文：人工智能（AI）是让机器模拟人类智能的技术，涵盖机器学习、自然语言处理与计算机视觉等领域。它正逐步渗透到医疗、金融与交通等行业，成为推动数字化转型的重要引擎。
分词：['人工智能', 'AI', '机器', '模拟', '人类', '智能', '技术', '涵盖', '机器', '学习', '自然语言', '处理', '计算机', '视觉', '领域', '它', '正', '逐步', '渗透到', '医疗', '金融', '交通', '行业', '成为', '推动', '数字化', '转型', '重要', '引擎']
==================================================
原文：大语言模型通过海量文本预训练，能够理解并生成自然语言，是近年来人工智能最热门的方向之一。以 GPT 为代表的生成式模型，已能胜任写作、编程与问答等多种任务。
分词：['大', '语言', '模型', '通过', '海量', '文本', '预', '训练', '理解', '并', '生成', '自然语言', '近年来', '人工智能', '最', '热门', '方向', '之一', '以', 'GPT', '为', '代表', '生成式', '模型', '已能', '胜任', '写作', '编程', '问答', '多种', '任务']
==================================================
原文：机器学习是人工智能的子集，它让计算机从数据中自动学习规律，而无需人工显式编写每一条规则。常见算法包括决策树、支持向量机与神经网络等，各有其适用场景。
分词：['机器', '学习', '人工智能', '子集', '它', '计算机', '数据', '中', '自动', '学习', '规律', '无需', '人工', '显式', '编写', '每', '一条', '规则', '常见', '算法', '包括', '决策树', '支持', '向量', '机与', '神经网络', '各有', '其', '适用', '场景']
==================================================
原文：计算机视觉致力于让机器看懂图像与视频，已广泛应用于人脸识别、自动驾驶和医疗影像诊断。借助深度学习，模型在目标检测与图像生成上的表现已接近甚至超越人类。
分词：['计算机', '视觉', '致力于', '机器', '看', '懂', '图像', '视频', '已', '广泛应用', '于', '人脸识别', '自动', '驾驶', '医疗', '影像', '诊断', '借助', '深度', '学习', '模型', '目标', '检测', '图像', '生成', '表现', '已', '接近', '甚至', '超越', '人类']
==================================================
原文：强化学习通过与环境交互、以奖励信号指导决策，曾在围棋等复杂游戏中战胜人类顶尖选手。AlphaGo 正是其代表作，展示了智能体在庞大状态空间中自主规划的能力。
分词：['强化', '学习', '通过', '环境', '交互', '以', '奖励', '信号', '指导', '决策', '曾', '围棋', '复杂', '游戏', '中', '战胜', '人类', '顶尖', '选手', 'AlphaGo', '正是', '其', '代表作', '展示', '智能', '体在', '庞大', '状态', '空间', '中', '自主', '规划', '能力']
==================================================
```

**bm25 计算、文档检索**

```python
def bm25(tokens: list, doc: TokenizedDocument, k1: float = 1.5, b: float = 0.75):
    score = 0.0
    n = len(tokenized_documents)
    avgdl = 1.0 * sum(doc.token_count() for doc in tokenized_documents) / n
    for token in tokens:
        nt = sum(1 for d in tokenized_documents if d.contains(token))
        idf_val = idf(nt, n)
        tf_val = doc.tf(token)

        score += idf_val * tf_val * (k1 + 1) / (tf_val + k1 * (1 - b + b * doc.token_count() / avgdl))
    return score


def search_documents(query_tokens: list, topk: int = 2):
    scores = {doc: bm25(query_tokens, doc) for doc in tokenized_documents}
    return sorted(scores.items(), key=lambda x: x[1], reverse=True)[:topk]
```

* `bm25`：计算 BM25 分数
* `search_documents`：检索文档，返回分数最高的 topk 个文档

**示例调用**

```python
if __name__ == "__main__":
    query = "我要学习人工智能"
    query_tokens = simple_cut(query)

    print(f"查询：{query}")
    print(f"查询分词：{query_tokens}")
    print(f"搜索结果：")
    for doc, score in search_documents(query_tokens):
        print(f" - {score:.4f} / {doc.text}")
```

结果输出：

```text
查询：我要学习人工智能
查询分词：['我要', '学习', '人工智能']
搜索结果：
 - 0.9598 / 机器学习是人工智能的子集，它让计算机从数据中自动学习规律，而无需人工显式编写每一条规则。常见算法包括决策树、支持向量机与神经网络等，各有其适用场景。
 - 0.8490 / 人工智能（AI）是让机器模拟人类智能的技术，涵盖机器学习、自然语言处理与计算机视觉等领域。它正逐步渗透到医疗、金融与交通等行业，成为推动数字化转型的重要引擎。
```

肉眼评价的话，对比其它三个文档，结果中的两个文档与查询词（人工智能、学习）更契合。

## 小结

BM25 能成为 ES/OpenSearch 全文检索的默认相关性算法，在于它用三个简单而有效的机制，把**词频**、**罕见度**和**文档长度**组合成了一个可控的打分公式：

1. **IDF 稀有度加权**：罕见词比常见词更能代表文档主题，因此获得更高权重
2. **词频饱和（$k_1$）**：词频对分数的贡献边际递减，避免了靠堆砌关键词刷分
3. **文档长度归一化（$b$）**：对长文档做惩罚，避免长文因词多而天然占优

BM25 是**词袋模型**，只统计词是否出现、出现多少次，**不理解词序与语义**——"苹果手机"与"手机苹果"对它等价，也无法召回同义词。这决定了它更适合作为**召回层**（快速从海量文档中筛出候选集），再交由向量检索或 rerank 模型做语义层面的精排，这也是当前 RAG 系统中**BM25 + 向量检索**混合召回成为主流的原因。