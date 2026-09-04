+++
date = '2026-09-04T10:39:44+08:00'
draft = false
title = '文本相似度计算：Jaccard 和 Dice 系数介绍与示例'
categories = ['人工智能', '自然语言处理']
tags = ['自然语言处理', '文本相似度', 'Jaccard 系数', 'Dice 系数']
toc = true
+++

## 文本相似度

文本相似度，简单说就是量化两段文本的**内容重合**、**语义相近程度**，是 NLP 最基础、最高频的任务之一。

日常业务中随处可见：

- 文章查重、内容去重
- 问答系统：用户提问与历史问题匹配
- 搜索推荐：相似内容召回
- 评论舆情：重复言论过滤
- 短文本匹配、指令对齐

文本相似度算法主要分两类：

1. 基于集合的离散匹配：Jaccard、Dice（只看词语「是否存在」，不看出现次数）
2. 基于向量的语义匹配：余弦相似度、Embedding 相似度（考虑词频、语义空间）

对于短文本、口语化文本、去重场景，不需要复杂向量计算，Jaccard 和 Dice 足够好用，且速度更快、结果可解释。

## Jaccard 和 Dice 系数公式

将两段文本分词、去重后得到两个词语集合 $A$、$B$。

- 交集 $A \cap B = \{a \in A \land a \in B\}$：两段文本共同包含的词语
- 并集 $A \cup B = \{a \in A \lor a \in B\}$：两段文本所有不重复词语的总和

**Jaccard 系数**

$$
J(A, B) = \frac{|A \cap B|}{|A \cup B|}
$$

即共同词语个数除以所有不重复词语个数，范围 $[0, 1]$。

特点：对文本差异更敏感，适合查重、去重、精准区分相似文本。

**Dice 系数**

$$
Dice(A,B) = \frac{2\cdot|A \cap B|}{|A|+|B|}
$$

即 2 倍共同词语个数除以所有词语总数，范围 $[0, 1]$。

特点：放大共同特征权重，弱化独有词语干扰，适合模糊匹配、相似召回、内容对齐。

## 文本特征提取

Jaccard、Dice 是集合算法，只能对**离散、去重、标准化的词语集合**计算，所以要先对文本进行分词、去重、标准化处理，步骤如下：

1. 文本清洗：移除、替换特殊字符等，保留有效词语
2. 分词处理：将文本转换为词语序列，去除停用词、标点符号等
3. 语义归一：将词语转换为统一的表示，统一同义表述
4. 集合构建：将分词后的词语转换为集合，去重
5. 计算相似度：根据 Jaccard 或 Dice 系数公式，计算两段文本的相似度

![](/imgs/learn-nlp/text-feature-extract-pipeline.png)

## 代码演示

**文本清洗**

* 全角字符统一转为半角：消除肉眼相同但编码不同的字符变体，保障分词、过滤、同义词匹配与特征提取的准确性
* 移除特殊字符：如标点符号、特殊字符等，保留有效词语

```python
def fullwidth_to_halfwidth(text: str) -> str:
    """全角转半角"""
    result = []
    for char in text:
        code = ord(char)
        if 0xFF01 <= code <= 0xFF5E:
            result.append(chr(code - 0xFEE0))
        elif code == 0x3000:
            result.append(chr(0x0020))
        else:
            result.append(char)
    return "".join(result)


def cleaned_text(text: str) -> str:
    """文本清洗"""
    text = fullwidth_to_halfwidth(text).lower()
    text = re.sub(r"[^\u4e00-\u9fa5a-zA-Z0-9\s]", " ", text)
    text = re.sub(r"\s+", " ", text).strip()
    return text
```

**分词处理**

* 调用 jieba 分词
* 移除停用词和空格

```python
STOP_WORDS = {
    "的", "了", "是", "在", "我", "有", "和", "就", "都", "而", "及", "与",
    "也", "很", "还", "个", "之", "对", "这", "那", "该", "些", "等", "可以",
    "应该", "能够", "会", "要", "把", "被", "让", "从", "向", "给", "到", "上", "下", "既",
}


def tokenized_text(text: str) -> list[str]:
    """文本分词"""
    
    tokens = [token for token in jb.lcut(text) if token not in STOP_WORDS and token.strip()]
    return tokens
```

**语义归一**

```python
SYNONYM_MAP = {
    "人工智能": ["AI", "ai", "人工智慧"],
    "大数据": ["海量数据"],
    "机器学习": ["ML"],
    "生活": ["日常生活"],
    "快速": ["飞速"],
    "用户": ["使用者", "客户", "个人"],
    "性能": ["效能", "效率"],
}

SYNONYM_REVERT_MAP = {syn: word for word, syns in SYNONYM_MAP.items() for syn in syns}

# 语义归一
def semantic_normalization(tokens: list[str]) -> list[str]:
    """语义归一"""
    normalized_tokens = []
    for token in tokens:
        if token in SYNONYM_REVERT_MAP:
            normalized_tokens.append(SYNONYM_REVERT_MAP[token])
        else:
            normalized_tokens.append(token)
    return normalized_tokens
```

**系数计算**

```python
def jaccard(l1: list[str], l2: list[str]) -> float:
    """jaccard 系数"""
    set1 = set(l1)
    set2 = set(l2)
    intersection = set1 & set2
    union = set1 | set2
    return len(intersection) / len(union)


def dice(l1: list[str], l2: list[str]) -> float:
    """dice 系数"""
    set1 = set(l1)
    set2 = set(l2)
    intersection = set1 & set2
    return 2 * len(intersection) / (len(set1) + len(set2))
```

**调用与结果**

```python
if __name__ == "__main__":
    text1 = "AI 的快速发展，给我们的生活带来了许多方便，也提升了个人的工作效率。"
    text2 = "人工智能技术飞速迭代，既为日常生活带来诸多便利，也有效提升了人们的工作效能"
    l1 = list(set(semantic_normalization(tokenized_text(cleaned_text(text1)))))
    l2 = list(set(semantic_normalization(tokenized_text(cleaned_text(text2)))))

    print("文本 1 分词:", l1)
    print("文本 2 分词:", l2)

    print(f"jaccard 系数: {jaccard(l1, l2):.4f}, dice 系数: {dice(l1, l2):.4f}")
```

输出：

```bash
文本 1 分词: ['人工智能', '快速', '发展', '我们', '生活', '带来', '许多', '方便', '提升', '用户', '工作效率']
文本 2 分词: ['人工智能', '技术', '快速', '迭代', '为', '生活', '带来', '诸多', '便利', '有效', '提升', '人们', '工作', '性能']
jaccard 系数: 0.2500, dice 系数: 0.4000
```

**结果分析**

从语义的表达来说，文本 1 和文本 2 是相近的，但 Jaccard 系数和 Dice 系数相对较低，表明它们的共同词语较少。所以除了分词技术外，还需要考虑语义归一化，将同义词转换为统一的表示，以提高相似度计算的准确性。

## 小结

Jaccard 和 Dice 作为轻量化文本相似度算法，不依赖大模型、无需向量训练，凭借低算力、高可解释、落地简单的优势，适用于短文本相似度计算的场景。

其缺点也足够明显，如不考虑词语频率、顺序等因素，导致无法准确评估文本的相似度。
