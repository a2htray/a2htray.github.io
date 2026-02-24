+++
date = '2026-01-15T15:01:00+08:00'
draft = false
title = '自然语言处理介绍'
categories = ['人工智能', '自然语言处理']
tags = ['技术笔记', 'NLP', '《自然语言处理基础》']
+++

目标：

1. 能够介绍“什么是 NLP”
2. 学习 NLP 的历史
3. NLP 与文本分析的区别
4. 实现特定处理工作
5. NLP 工程的不同阶段

## NLP

自然语言处理（NLP）是一个与多个学科交叉的领域，它脱胎于人工智能、语言学、形式语言及编译器等学科。

多学科：

* 人工智能（Artificial Intelligence）：NLP 的核心母体学科，为其提供算法与算力基础
* 语言学（Linguistics）：提供人类语言的语法、语义、语用规则体系
* 形式语言（Formal Languages）：研究抽象符号系统的语法规则
* 编译器（Compilers）：负责将高级语言转换为机器语言

<span style="color: blue;">Previously, a traditional **rule-based system** was used for computations.</span>

> 此前，计算任务采用的是传统的基于规则的系统。

<span style="color: blue;">Today, computations on natural language are being done using machine learning and deep learning techniques.</span>

> 如今，自然语言相关的计算任务正通过机器学习与深度学习技术来完成。

基于机器学习的自然语言处理的主要工作开始于 1980 年。

## Vs. 文本分析

<span style="color: blue;">Text analytics is the method of extracting meaningful insights and answering questions from text data.</span>

文本分析是从文本数据中提取有价值的洞察并解答相关问题的方法。

文本分析：从文本中提取价值。

NLP 的处理对象除了文本，还可以是语音，NLP 大体可分成 Natural Language Understanding (NLU) and Natural Language Generation (NLG) 两个核心子任务。

核心子任务：

* NLU：让计算机理解语言的过程；
* NLG：让计算机表达出人类理解的内容的过程；

## NLP 步骤

* Tokenization
* PoS Tagging（词性标注）

Tokenization（分词）：将一个完整句子拆解为其组成基本单位（token）的过程。

n-grams（n 元语法）：基于滑动窗口思想，将文本序列（字符 / 单词）切分为连续的 n 个元素组成的子序列，是 NLP 中表示文本局部特征的基础方法。如：

```text
I am reading a book.
```

* 一元语法（unigram, n=1）：`["I", "am", "reading", "a", "book"]`
* 二元语法（bigram, n=2）：`["I am", "am reeding", "reading a", "a book"]`
* 三元语法（trigram, n=3）：`["I am reading", "am reeding a", "reading a book"]`

> n 的选择原则：n 越小，计算成本越低，但丢失的上下文信息越多；n 越大，上下文信息越丰富，但易出现数据稀疏问题。

PoS Tagging：对分词后的文本中每个 token（词元） 标注其语法词性（如名词、动词、形容词等）的过程，是 NLP 文本预处理与句法分析的关键步骤。

核心作用：

1. 文本预处理：为后续的句法分析、命名实体识别（NER）、情感分析等任务提供语法特征
2. 歧义消解：区分多义词的词性（如 book 作名词是 “书籍”，作动词是 “预订”）
3. 提高模型精度：在文本分类、机器翻译等任务中，**词性特征可增强模型对文本语义的理解**

> nltk 包中的词性标签。

停用词（stop word）：高频出现但无实质语义贡献的通用词汇，仅起语法支撑作用，是 NLP 文本清洗环节的核心处理对象。

- 英文停用词：the, a, an, is, are, in, on, and, of
- 中文停用词：的, 是, 在, 和, 就, 但, 了

停用词移除的目的：

1. 降低特征维度，减少计算资源消耗
2. 突出核心词汇（如名词、动词），提升模型对文本主题的捕捉能力
3. 避免高频无义词干扰后续任务（如文本分类、主题建模）


