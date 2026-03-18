+++
date = '2026-03-18T13:41:40+08:00'
draft = false
title = 'Ollama：部署与使用'
categories = ['人工智能', '大语言模型']
tags = ['Ollama', 'LLM']
toc = true
+++

Ollama 是大模型的运行平台，本文主要介绍 Ollama 的安装与使用。

## 安装

听说 Ollama 类似于 Docker 一样，出了 Desktop 版本，所以这次我选择下载 `.dmg` 文件方式安装。

访问官网下载页面 [https://ollama.com/download](https://ollama.com/download)

下载 `.dmg` 安装文件

![](/imgs/learn-ai/deploy_ollama_01.png)

双击安装，拖动图标至 Applications

![](/imgs/learn-ai/deploy_ollama_02.png)

打开程序，如下图：

**Desktop**

![](/imgs/learn-ai/deploy_ollama_03.png)

**Terminal**

```bash
$ ollama -v
ollama version is 0.17.4
```

## 下载模型

考虑到磁盘的现有空间和后续使用，先下两个小模型使用，分别是 `qwen3.5:2b` 和 `qwen3.5:9b`。

```bash
$ ollama pull qwen3.5:2b
pulling manifest 
pulling b709d81508a0: 100% ▕████████████████████████████████████████████████████▏ 2.7 GB                         
pulling 9be69ef46306: 100% ▕████████████████████████████████████████████████████▏  11 KB                         
pulling 9371364b27a5: 100% ▕████████████████████████████████████████████████████▏   65 B                         
pulling ee043a99abe5: 100% ▕████████████████████████████████████████████████████▏  473 B
verifying sha256 digest 
writing manifest 
success
$ ollama pull qwen3.5:9b
pulling manifest 
pulling dec52a44569a: 100% ▕████████████████████████████████████████████████████▏ 6.6 GB                         
pulling 7339fa418c9a: 100% ▕████████████████████████████████████████████████████▏  11 KB                         
pulling 9371364b27a5: 100% ▕████████████████████████████████████████████████████▏   65 B                         
pulling be595b49fe22: 100% ▕████████████████████████████████████████████████████▏  475 B                         
verifying sha256 digest 
writing manifest 
success
```

| 模型名称       | 大小   | 上下文大小 | 输入    |
|------------|------|-------|-------|
| qwen3.5:2b | 2.7G | 256K  | 文本/图片 |
| qwen3.5:9b | 6.6G | 256K  | 文本/图片 |


> ollama 的子命令规则与 docker 相似。

## 使用

**查看当前环境的模型**

```bash
$ ollama list
NAME          ID              SIZE      MODIFIED          
qwen3.5:9b    6488c96fa5fa    6.6 GB    12 minutes ago       
qwen3.5:2b    324d162be6ca    2.7 GB    About an hour ago
```

**启动 Ollama**

1. 点击桌面图标启动
2. 命令行启动

```bash
$ nohup ollama serve &
```

**运行模型**

```bash
# 进入交互式问答
$ ollama run qwen3.5:2b
```

按 Ctrl + D 退出。

**查看已运行模型**

```bash
$ ollama ps  
NAME          ID              SIZE      PROCESSOR    CONTEXT    UNTIL              
qwen3.5:2b    324d162be6ca    4.1 GB    100% GPU     4096       3 minutes from now
```

**停止运行模型**

```bash
$ ollama stop qwen3.5:2b
```

## API 文档

Ollama 默认运行在 11434 端口，所以 Base URL 为：

```text
http://localhost:11434
```

使用 cURL 调用模型列表接口，如下：

```bash
$ curl -s http://localhost:11434/api/tags | jq                       
{
  "models": [
    {
      "name": "qwen3.5:9b",
      "model": "qwen3.5:9b",
      "modified_at": "2026-03-18T15:39:21.748402333+08:00",
      "size": 6594474711,
      "digest": "6488c96fa5faab64bb65cbd30d4289e20e6130ef535a93ef9a49f42eda893ea7",
      "details": {
        "parent_model": "",
        "format": "gguf",
        "family": "qwen35",
        "families": [
          "qwen35"
        ],
        "parameter_size": "9.7B",
        "quantization_level": "Q4_K_M"
      }
    },
    {
      "name": "qwen3.5:2b",
      "model": "qwen3.5:2b",
      "modified_at": "2026-03-18T14:48:28.343656637+08:00",
      "size": 2741192820,
      "digest": "324d162be6ca5629ae4517c8710434d0bd2d665bc94dbad46e9af8fbf8a2f0df",
      "details": {
        "parent_model": "",
        "format": "gguf",
        "family": "qwen35",
        "families": [
          "qwen35"
        ],
        "parameter_size": "2.3B",
        "quantization_level": "Q8_0"
      }
    }
  ]
}
```

**!注意！**

API 文档中有两个功能相似的接口：

1. POST /api/generate
2. POST /api/chat

两者的区别在于：<font color="blue">**/api/generate 用于单次文本的生成，/api/chat 用于多轮对话**</font>，如下两个调用示例：

```bash
$ curl -s http://localhost:11434/api/generate -d '{
  "model": "qwen3.5:2b",
  "prompt": "你是谁，10 个字内回复我。",
  "stream": false,
  "think": false
}' | jq
{
  "model": "qwen3.5:2b",
  "created_at": "2026-03-18T08:50:05.767219Z",
  "response": "我是通义千问，人工智能。",
  "done": true,
  "done_reason": "stop",
  "context": [
    248045,
    846,
    198,
    144810,
    3709,
    16,
    15,
    220,
    121493,
    95988,
    100535,
    95815,
    1710,
    248046,
    198,
    248045,
    74455,
    198,
    248068,
    271,
    248069,
    271,
    103724,
    95974,
    96359,
    96617,
    96094,
    3709,
    109015,
    1710
  ],
  "total_duration": 1595669583,
  "load_duration": 168393750,
  "prompt_eval_count": 22,
  "prompt_eval_duration": 1219554417,
  "eval_count": 9,
  "eval_duration": 203112833
}
```


```bassh
$ curl -s http://localhost:11434/api/chat -d '{
  "model": "qwen3.5:2b",
  "messages": [
    {
      "role": "user",
      "content": "你是谁，10 个字内回复我。"
    }
  ],
  "stream": false,
  "think": false
}' | jq
{
  "model": "qwen3.5:2b",
  "created_at": "2026-03-18T08:51:02.73661Z",
  "message": {
    "role": "assistant",
    "content": "我是通义实验室研发的 AI 助手 Qwen3.5。"
  },
  "done": true,
  "done_reason": "stop",
  "total_duration": 1748286292,
  "load_duration": 192718167,
  "prompt_eval_count": 22,
  "prompt_eval_duration": 1185704375,
  "eval_count": 15,
  "eval_duration": 356360626
}
```

支持多轮对话，则需要在 `messages` 数组中加入对话内容。

## 资源

* [官方 API 文档](https://docs.ollama.com/api/introduction)