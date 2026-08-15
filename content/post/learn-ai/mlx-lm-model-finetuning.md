+++
date = '2026-08-15T21:56:27+08:00'
draft = false
title = '模型微调示例：油菜品种信息查询助手'
categories = ['人工智能', '大模型模型']
tags = ['模型微调', '大语言模型', 'LoRA']
toc = true
+++

> 在 Apple Silicon（M 系列芯片）上，用 **MLX-LM** 对一个开源大模型做 LoRA 微调，
> 把它变成一个"油菜品种信息查询"垂域助手，并对比微调前后的效果。
> 本文记录了从环境准备、数据构造、微调到评测的完整过程，以及踩过的几个版本坑。

## 概述与背景

**目标**：构建一个能回答"某油菜品种生育期多长、含油量多少、适合种在哪里、抗什么病"等问题的垂域大模型。

**为什么用 MLX-LM**：Apple Silicon 的 CPU / GPU / 神经网络引擎共享同一块物理内存（统一内存），MLX 利用这一点做到"零拷贝"，模型可以大到吃满整机内存，而不受独立显存限制。因此可以在没有独立显卡的 Mac 上，用接近 NumPy 的写法做 GPU/神经网络引擎加速的训练与推理。MLX-LM 是其在语言模型方向的开源实现，提供 `generate` / `lora`（微调）/ `fuse`（合并）等命令。

**基座模型选择**：`mlx-community/Qwen2.5-3B-Instruct-4bit`（约 2GB、中文能力强、16GB 内存即可跑，且是 MLX 格式，开箱即用）。

**微调方式**：LoRA（低秩适配）。只训练注入的一小部分低秩参数，基座权重冻结，显存占用小、训练快、训完权重只有几十 MB，且能直接合并回基座分发。

## 环境准备

需要 macOS + Apple Silicon（M1 及以上），Python ≥ 3.9。本文所用环境为 conda `mldl` 环境。

安装必要的依赖包：

```bash
$ pip install -U mlx mlx-lm "huggingface_hub[cli]" hf_transfer rouge-score lm_eval
```

验证版本：

```bash
$ pip list | grep -E "mlx-lm|huggingface_hub|hf_transfer|rouge_score"
hf_transfer               0.1.9
huggingface_hub           1.27.0
mlx-lm                    0.31.3
rouge_score               0.1.2
```

## 下载模型

新版的 Hugging Face 客户端把 `huggingface-cli` 整合进了 `hf` 命令。用 `hf download` 把模型拉到本地目录：

```bash
$ HF_HUB_DISABLE_XET=1 \
    hf download mlx-community/Qwen2.5-3B-Instruct-4bit \
    --local-dir ./models/Qwen2.5-3B-Instruct-4bit
```

- `--local-dir` 把权重直接落到 `./models/Qwen2.5-3B-Instruct-4bit`，之后所有命令的 `--model` 都指这个本地路径，不用每次从 Hugging Face 联网拉取。
- `HF_HUB_DISABLE_XET=1` 关闭 HF 的 Xet/CAS 大文件传输层，规避偶发的 `401 Unauthorized` 下载错误。

下载完成后验证目录里有 `config.json`、`*.safetensors`、`tokenizer*` 即可。

先跑一次推理，确认链路通：

```bash
$ mlx_lm generate \
    --model ./models/Qwen2.5-3B-Instruct-4bit \
    --prompt "你好，介绍一下你自己"
==========
你好！我是Qwen，由阿里云研发的超大规模语言模型。我被训练来生成各种类型的文本，包括但不限于故事、诗歌、新闻文章、科学论文等。我的目标是帮助用户获得信息、娱乐、解决问题或完成特定任务。我可以回答问题、提供信息、创作文字内容，甚至参与对话。我能够处理多种语言，并且在理解上下文和生成连贯内容方面表现出色。如果你有任何问题或需要帮助的地方，欢迎随时提问
==========
Prompt: 33 tokens, 62.431 tokens-per-sec
Generation: 100 tokens, 84.512 tokens-per-sec
Peak memory: 1.872 GB
```

- 提示词含 33 个 token，模型生成了 100 个 token 作为回复
- token 输入速度 62.431 tokens/sec，输出速度 84.512 tokens/sec
- 峰值内存使用 1.872 GB

> 注：`python -m mlx_lm.generate`（点号写法）已被官方弃用，控制台命令用 `mlx_lm generate` 即可。

## 准备数据集

**内容构成**：

- 覆盖 **14 个真实油菜品种**：中油杂19、华油杂62、沣油737、川油36、秦优10号、中油杂2号、秦优7号、浙油50、油研10号、中双11号、华油杂9号、中油杂39、青杂2号、青杂（春油菜）等。
- 每品种生成 **8 类定向查询**：生育期、含油量、适种区域、育种单位/审定、抗性、产量、双低品质、综合介绍。
- 跨品种 **推荐/对比** 类：高油、早熟，长江中游/下游，春油菜区，机收等。
- **概念解释** 类：双低、甘蓝型、三系杂交、CMS、芥酸、硫苷、菌核病、硼肥等。
- 1 条 **"知识边界如实拒答"** 样本，训练模型"不知道就说不知道"。

**数据量**：`train` 86 条 / `valid` 18 条 / `test` 18 条（小数据集常见约 7:1.5:1.5 划分）。

`data` 目录下包含 `train.jsonl`、`valid.jsonl` 和 `test.jsonl` 3 个文件，单条记录如下：

```json
{
    "messages": [
        {
            "role": "system", 
            "content": "你是油菜品种信息查询助手。仅依据官方品种审定与登记资料回答油菜品种的类型、育种单位、审定/登记编号、全生育期、含油量、芥酸与硫苷含量、适种区域、抗性（菌核病/病毒病/倒伏/抗寒）与产量等问题；不确定或超出知识范围时如实说明，不编造数据。回答简洁、准确，必要时给出数值范围。"
        },
        {
            "role": "user",
            "content": "沣油737含油率大概多少？"
        },
        {
            "role": "assistant",
            "content": "沣油737种子含油率约44.86%，属双低优质油菜（低芥酸、低硫苷）。"
        }
    ]
}
```

**三个文件的角色与用法**

| 文件 | 条数 | 角色 | 何时用 | 作用 |
|---|---|---|---|---|
| `train.jsonl` | 86 | 训练集 | 全程 | 模型从这里面"学"参数（更新 LoRA 权重） |
| `valid.jsonl` | 18 | 验证集 | 训练**中途**反复用 | 监控泛化、决定早停、调超参（如 rank / 学习率） |
| `test.jsonl` | 18 | 测试集 | 训练**结束后只用一次** | 报告最终真实效果，模拟"考试" |

> 🔑 **铁律**：`test.jsonl` 在训练和调参阶段**绝对不能碰**。一旦你拿 test 结果去改 rank、降学习率、
> 或决定多训几轮，test 就"偷看了答案"，不再是泛化能力的代表，退化为第二个 valid。
> 正确链路：`train → 用 valid 反复调/早停 → 定稿后，才用 test 评一次出报告`。

![](/imgs/learn-ai/relationship.png)

## 模型微调

创建 LoRA 配置文件 `lora_config.yaml`：

```yaml
lora_parameters:
  rank: 16       # LoRA 秩，越大容量越高（8~32 常用）
  scale: 20.0    # 适配 4bit 基座，通常 20 左右
  dropout: 0.1   # 防过拟合
```

执行微调：

```bash
$ mlx_lm lora \
    --model ./models/Qwen2.5-3B-Instruct-4bit \
    --data ./data \
    --train \
    --adapter-path ./adapters \
    --num-layers 8 \
    --batch-size 4 \
    --iters 110 \
    --learning-rate 1e-4 \
    --mask-prompt \
    -c lora_config.yaml
```

**关键参数说明**：

- `--data ./data`：MLX-LM 会自动读取该目录下 `train.jsonl` / `valid.jsonl` / `test.jsonl`。
- `--num-layers 8`：只对模型**顶层**若干层加 LoRA。层数越多拟合越强但越易过拟合；3B 模型常取 8–16。
- `--mask-prompt`：只对"回答"部分算 loss，避免模型去模仿提问。
- `--iters 110`：按步数训练。`iters = 轮数 × (训练样本 / batch)` = `5 × (86/4)` ≈ 110 步。
- 训练过程会打印 loss，loss 不再下降即可停；想看曲线可加 `--report-to wandb`。

训练产出的 LoRA 权重保存在 `./adapters/`（通常仅几十 MB）。

> 可选增强（按需加进命令）：`--steps-per-eval 20`（每 20 步在 valid 上验证，做早停监控）、
> `--save-every 50`（每 50 步存一次）、`--grad-checkpoint`（内存吃紧时开，省内存但略慢）。

---

## 评估微调后的模型

### 为什么不用 `mlx_lm evaluate`

在 0.31.3 中，`mlx_lm evaluate` 已被重写为 **lm-evaluation-harness** 驱动，**只认 `--tasks`**（hellaswag / mmlu / gsm8k 等公开学术基准），**既不支持 `--data`（喂自定义 jsonl），也不支持 `--adapter-path`（评 LoRA）**。因此：

```bash
# 下面这种"自定义数据评测"在 0.31.3 已不存在，会报 "the following arguments are required: --tasks"
mlx_lm evaluate --model ./models/Qwen2.5-3B-Instruct-4bit --data ./data/test.jsonl   # ❌
```

它适合做的是**通用能力体检**（如 `mlx_lm evaluate --tasks hellaswag arc_challenge --limit 100`），
而不是评测你的领域 `test.jsonl`。要评"基座 vs 微调"在自定义测试集上的表现，得用程序化方式。

### 评测脚本 `eval.py`

`eval.py` 直接调用 MLX-LM 的程序化 API：

- `load(model, adapter_path=...)` 加载基座或带 LoRA 的模型；
- 用 `tokenizer.apply_chat_template(..., add_generation_prompt=True)` 构建提示，**与 `mlx_lm generate` 的模板完全一致**，
  避免评测 / 推理不一致；
- 指标在脚本内**纯 Python 自实现**（字符级 ROUGE-L / tokenF1 / ExactMatch），**零额外依赖**，无需 `rouge_score`；
- 读取 `messages` 格式（兼容 `prompt`/`completion` 简单格式），逐条生成并打分，明细写入 JSON 便于排查。

### 评测结果

**基座模型**

```bash
$ python eval.py \
    --model ./models/Qwen2.5-3B-Instruct-4bit \
    --data ./data/test.jsonl \
    --output ./eval_base.json
========== 汇总 ==========
{
  "model": "./models/Qwen2.5-3B-Instruct-4bit",
  "adapter_path": null,
  "num_samples": 18,
  "ROUGE-L": 0.27695110302209386,
  "tokenF1": 0.38343017108608496,
  "ExactMatch": 0.0
}
```

**微调模型**

```bash
$ python eval.py \
    --model ./models/Qwen2.5-3B-Instruct-4bit \
    --adapter-path ./adapters \
    --data ./data/test.jsonl \
    --output ./eval_lora.json
========== 汇总 ==========
{
  "model": "./models/Qwen2.5-3B-Instruct-4bit",
  "adapter_path": "./adapters",
  "num_samples": 18,
  "ROUGE-L": 0.7296865528423073,
  "tokenF1": 0.8065462200608561,
  "ExactMatch": 0.1111111111111111
}
```

**结果对比表**

| 指标 | 基座 (base) | 微调 (LoRA) | 提升 |
|---|---:|---:|---:|
| ROUGE-L | 0.277 | 0.730 | +163% |
| tokenF1 | 0.383 | 0.807 | +111% |
| ExactMatch | 0.000 | 0.111 | +0.111（18 题命中 1 题） |

![](/imgs/learn-ai/base_tuning_compare.png)

**结论**

1. **LoRA 微调带来显著提升**：ROUGE-L 从 0.277 跃升到 0.730（+163%）、tokenF1 从 0.383 升到 0.807（+111%），
   说明模型对油菜品种问答的"表述结构与关键信息覆盖"已大幅贴合标准答案，微调方向正确有效。
2. **绝对质量仍有空间**：ExactMatch 仅 0.111（18 题仅 1 题完全一致），且基座 tokenF1 已 0.383，
   说明模型"见过字"但离"一字不差"还差一截。后续可增大训练轮数 / 数据量、降低学习率抑制过拟合，
   或在 `eval.py` 中加**关键字段抽取核对**（含油量、适种区域、生育期逐字段比对）来严判领域事实正确性。

## 评测指标详解

ROUGE-L 与 tokenF1 都是**表面重叠度（surface overlap）**指标——只比"预测答案"和"标准答案"在字面上的重合程度，不判断语义对错。区别在"怎么算重合"。

- **ROUGE-L（最长公共子序列，顺序敏感）**
  - 在预测和 gold 里找 LCS，字符可不全连着但相对顺序要对；
  - 召回 R = LCS 长度 / gold 长度，精度 P = LCS 长度 / 预测长度，ROUGE-L = 2·R·P / (R+P)；
  - 把句子换个语序，LCS 变短、分数掉；范围 0~1，越高越好；
  - 适合衡量"表述结构 / 流畅度"像不像标准答案（摘要、长回答等顺序敏感任务）。
- **tokenF1（字符集合交集，顺序无关）**
  - 把预测和 gold 各看成字符集合，算集合 F1；
  - R = 共同字符数 / gold 字符数，P = 共同字符数 / 预测字符数，tokenF1 = 2·R·P / (R+P)；
  - 完全不在乎顺序，只看"字覆盖没覆盖"；范围 0~1，越高越好；
  - 适合衡量"关键信息字覆盖没覆盖"（抽取 / 短答案，如"适种区域""含油量数值"有没有出现）。
- **ExactMatch（精确匹配）**
  - 去空白后预测与 gold 完全一致的比例，最严格；领域抽取 / 短答案类任务的"一字不差"判据。

> ⚠️ 三者都只量"形似"，量不了"事实真"。例如预测把"含油量 45%"写成"含油量 54%"只差一字，
> ROUGE-L / tokenF1 都会很高，但事实错了。领域上线前需配人工或**字段级核查**兜底。

下图用同一组重排示例（`gold=适宜黄淮长江下游`，`pred=适宜长江下游黄淮`）直观对比两个指标：
ROUGE-L 因顺序被打乱而扣分，tokenF1 因字全覆盖而满分。

![](/imgs/learn-ai/rougel-tokenf.png)

## 踩坑记录

**① 下载报 `401 Unauthorized ... cas-server.xethub.hf.co`（Xet/CAS 传输层 bug）**

新版 HF 用 Xet 大文件传输，公开模型也会偶发 401。加 `HF_HUB_DISABLE_XET=1` 关掉它，回退普通 HTTPS 下载即可。
若仍 401，清掉半截缓存：`rm -rf ~/.cache/huggingface/hub/models--mlx-community--Qwen2.5-3B-Instruct-4bit` 再下。
`Qwen2.5-3B-Instruct-4bit` 是公开模型，正常无需 `hf login`。

**② `python -m mlx_lm.generate` 弃用告警**

控制台改用 `mlx_lm generate`（或 `python -m mlx_lm generate` 中间加空格），告警即消失。

**③ `mlx_lm lora` 报 `unrecognized arguments: --num-epochs 5 --lora-parameters {...}`**

0.31.3 已移除这两个旧参数：
- `--num-epochs` → 改用 `--iters`（按"轮数 × 样本数 / batch"换算，本例 5 × 86/4 ≈ 110）；
- `--lora-parameters '{...}'` → 改成 YAML 配置文件 `-c lora_config.yaml`。

**④ `mlx_lm evaluate` 报 `the following arguments are required: --tasks`**

0.31.3 的 `evaluate` 已变为 lm-eval-harness 驱动，只认 `--tasks`，**不支持 `--data` 与 `--adapter-path`**，
无法评自定义 jsonl 或 LoRA。自定义领域评测改用程序化 `eval.py`；通用能力体检才用 `mlx_lm evaluate --tasks ...`。

**⑤ 本地 MLX 模型与标准 HF 模型的区分**

只有 `mlx-community/*` 这类仓库是 MLX 格式（`config.json` 含 `quantization` 块 + `model.safetensors`）。
若 `hf download` 的是普通 HF 模型（如 `Qwen/Qwen2.5-3B-Instruct`），需先转格式：
`mlx_lm convert --hf-path ./models/Qwen2.5-3B-Instruct -q --local-dir ./models/Qwen2.5-3B-Instruct-4bit`。

*本文基于 mlx-lm 0.31.3、Apple Silicon（macOS）环境实测整理。*