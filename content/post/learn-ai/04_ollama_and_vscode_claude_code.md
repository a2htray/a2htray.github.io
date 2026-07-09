+++
date = '2026-07-09T09:37:05+08:00'
draft = false
title = '生产力工具集成：Ollama + VS Code + Claude Code，效果不太行...'
categories = ['人工智能', '应用平台']
tags = ['Ollama', 'VS Code', 'Claude Code', '工具集成', '工具安装', '生产力工具']
toc = true
+++

最近一直在捣鼓 AI 工具，想着说给我的开发工作提供强有力的帮助，但在使用的过程中，**token 成本太高**是我急需解决的一个问题。

这让我想起了 Ollama，网上也查了相关资料，Claude Code 通过配置是可以和本地的大模型进行通信的，所以有了这篇文章。

> 本文旨在分享将 Ollama、VS Code、Claude Code 生产力工具的集成过程。

## 安装 Ollama

Ollama 是轻量化本地大模型运行框架，可以一键拉取、调度、管理 LLM。安装过程可以参见之前写的文章 [https://a2htray.github.io/post/learn-ai/01_deploy_ollama_on_local/](https://a2htray.github.io/post/learn-ai/01_deploy_ollama_on_local/)。

因为我是个人使用，所以下的是桌面版本，也可以直接在它的 github releases 中下载其二进制版本。

![](/imgs/learn-ai/ovc-01.png)

安装好后，在终端验证是否安装完成：

```bash
$ ollama -v     
ollama version is 0.31.1
```

考虑到这次是为了以后开发，所以咬咬牙就想下个大的，选择了 19G 的 `qwen3-coder:30b`。

![](/imgs/learn-ai/ovc-02.png)

> 笔记本配置是 Apple M2 Pro 16G 内存。

但实际下载后运行，实在是太卡了，CPU 直接飙到了 100%，笔记本也卡的厉害，最终选择了 `qwen3.5:9b`，6.6G。

下载模型的过程如下：

```bash
$ ollama pull qwen3.5:9b
```

## 安装 VS Code

IDE 我用的是 VS Code，相信也是大多数开发者的首选，官网 [https://code.visualstudio.com/](https://code.visualstudio.com/) 下载安装包，点击安装即可。

为了和 Claude Code 搭配起来，需要安装插件 `Claude Code for VS Code`，如下图：

![](/imgs/learn-ai/ovc-03.png)

安装好后，显示右侧边，有 Claude Code 的 tab 的显示，如下图：

![](/imgs/learn-ai/ovc-04.png)

> 在没有安装 Claude Code 或没配置之前，显示的是订阅信息，我是已经配置了，所以可以进行相关操作。

## 安装 Claude Code

我本地的 Node、NPM 版本如下，供参考：

```bash
$ node -v
v24.15.0
$ npm -v
11.12.1
```

安装命令如下：

```bash
$ npm install -g @anthropic-ai/claude-code
```

安装太慢，可以切换成国内的镜像源。安装好后，验证下版本：

```bash
claude -v
2.1.205 (Claude Code)
```

## 安装 LiteLLM

LiteLLM 是大模型统一代理中间件，将 Ollama 的 function calling 格式转换成 Claude Code 认可的 Anthropic 协议，不然模型只能回复你的问题，但无法在项目中编辑。

我新开了一个 conda 环境，用于跑 LiteLLM 的代理服务，命令如下：

```bash
$ conda create -n ai-llm python=3.12
$ conda activate ai-llm
$ pip install 'litellm[proxy]'
```

新建 LiteLLM 的配置文件 `litellm_config.yaml`，内容如下：

```yaml
model_list:
  - model_name: "qwen3.5:9b"
    litellm_params:
      model: "ollama_chat/qwen3.5:9b"
      api_base: "http://localhost:11434"
      api_key: "ollama"
    model_info:
      supports_function_calling: true

litellm_settings:
  drop_params: true
```

* model_name 是给 Claude Code 的模型名
* model 是 Ollama 中的模型名

启动代理服务：

```bash
$ litellm --config litellm_config.yaml --port 4000
INFO:     Started server process [61704]
INFO:     Waiting for application startup.

   ██╗     ██╗████████╗███████╗██╗     ██╗     ███╗   ███╗
   ██║     ██║╚══██╔══╝██╔════╝██║     ██║     ████╗ ████║
   ██║     ██║   ██║   █████╗  ██║     ██║     ██╔████╔██║
   ██║     ██║   ██║   ██╔══╝  ██║     ██║     ██║╚██╔╝██║
   ███████╗██║   ██║   ███████╗███████╗███████╗██║ ╚═╝ ██║
   ╚══════╝╚═╝   ╚═╝   ╚══════╝╚══════╝╚══════╝╚═╝     ╚═╝


#------------------------------------------------------------#
#                                                            #
#              'I don't like how this works...'               #
#        https://github.com/BerriAI/litellm/issues/new        #
#                                                            #
#------------------------------------------------------------#

 Thank you for using LiteLLM! - Krrish & Ishaan
```


## 集成工具

至此，已经完成了 4 件事：

1. 安装了 Ollama，并下载了模型
2. 安装了 VS Code 及插件 `Claude Code for VS Code`
3. 安装了 Claude Code
4. 安装了 LiteLLM

接下来需要配置 Claude Code，我因为全局有配置，所以以项目级配置为例，如下：

```bash
$ mkdir -p test_project/.claude
$ touch test_project/.claude/settings.local.json
```

`settings.local.json` 的内容如下：

```json
{
  "env": {
    "ANTHROPIC_BASE_URL": "http://0.0.0.0:4000",
    "ANTHROPIC_AUTH_TOKEN": "ollama",
    "ANTHROPIC_MODEL": "qwen3.5:9b"
  }
}
```

* `test_project` 是项目名
* `.claude/settings.local.json` 中的配置只对 `test_project` 有效


> 全局配置在 `~/.claude/settings.json` 文件中。

配置好 Claude Code 后，启动 Ollama 和 LiteLLM 代码，如下：

```bash
$ ollama serve
$ litellm --config litellm_config.yaml --port 4000
```

VS Code 打开 `test_project` 项目，在 Claude Code 插件中输入一段需求，如下：

![](/imgs/learn-ai/ovc-05.png)

## 小结

其实本地部模型做开发，是挺尴尬的，大点的模型又跑不起来，小点的模型又不好用，跑一个小任务也卡、耗时，而且效果也不行。当然也可以搞个高配的服务器，部个大点的模型，然后通过 VS Code SSH 开发，起码把 token 的成本转成了固定资产。

最后总结下集成过程，整体比较简单，步骤如下：

1. 下工具 Ollama、VS Code、Claude Code
2. 下模型，装插件 - `ollama pull`、`Claude Code for VS Code`
3. Claude Code 配置 `setting.local.json`、启模型 `ollama serve`
4. LiteLLM 代理服务

还是老老实实用平台提供的 API 方式接入吧，本地跑模型真得能让你怀疑人生。