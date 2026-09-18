+++
date = '2026-09-18T10:52:03+08:00'
draft = false
title = '可视化 - 线图：matplotlib vs. ggplot2'
categories = ['后端技术', '可视化']
tags = ['Python', '可视化', 'matplotlib', 'ggplot2', "线图", "R"]
toc = true
+++

![](/imgs/learn-vis/ScreenShot_2026-09-18_105617_844.png)

> 本文为可视化系列第 01 篇，将介绍线图的基础概念、适用场景与数据集选取，并分别使用 Python Matplotlib 和 R ggplot2 两套工具，演示绘制同一幅线图。

## 线图

线图通过**线段连接离散的数据点**，用来呈现指标沿着连续变量（最常见为时间）的变化情况，能够直观反映数值的增减波动、变化快慢以及拐点与周期特征，多系列线图还可以同时对比多组数据的走势差异。

适用场景：

* 展示数据随时间的连续变化（如每日股价、月度营收）
* 观察指标变化趋势、周期性波动、增长 / 衰减速率
* 多组同维度数据趋势对比

不适合：分类项目大小对比、占比展示。

## 数据集选取

**数据集 1 - 油菜生育期动态数据集：rape_growth_dynamics.csv**

展示油菜播种后随生育天数推进，株高、叶面积、干物质积累的连续变化趋势。

|字段|释义|单位|
|--|--|--|
|day_after_sowing|播种后天数|d|
|plant_height_min_cm|株高最小值|cm|
|plant_height_max_cm|株高最大值|cm|
|plant_height_mean_cm|株高平均值|cm|
|leaf_area_min_cm2|单株叶面积最小值|cm²|
|leaf_area_max_cm2|单株叶面积最大值|cm²|
|leaf_area_mean_cm2|单株叶面积平均值|cm²|
|dry_matter_g|单株干物质积累量|g|

## 可视化 1 - 油菜株高随播种后天数变化

**Python Matplotlib**

```python
import matplotlib.pyplot as plt
import pandas as pd

df = pd.read_csv("../data/rape_growth_dynamics.csv")

plt.rcParams["font.sans-serif"] = ["Heiti TC"]
plt.rcParams["axes.unicode_minus"] = False

fig, ax = plt.subplots(figsize=(10, 6))

ax.plot(df["day_after_sowing"], df["plant_height_max_cm"], label="最大值", color="#D62728")
ax.plot(df["day_after_sowing"], df["plant_height_mean_cm"], label="平均值", color="#0072B2")
ax.plot(df["day_after_sowing"], df["plant_height_min_cm"], label="最小值", color="#FFB600")

ax.set_xlabel("播种后天数 (d)")
ax.set_ylabel("株高 (cm)")
ax.set_title("油菜株高随播种后天数变化")

ax.grid(True, axis="x", linestyle="--", color="#CBC8C8")
ax.legend(loc="lower right")

fig.savefig("plant_height_by_matplotlib.png")
```

* `ax.grid` 方法设置网格线，`True` 表示显示网格线，`axis="x"` 表示仅在x 轴显示网格线，`linestyle="--"` 表示虚线，`color="#CBC8C8"` 表示网格线颜色为 #CBC8C8；
* `ax.legend` 方法添加图例，`loc="lower right"` 表示图例位置为右下角，其它值如 `upper left`、`upper right` 等。

![](/imgs/learn-vis/plant_height_by_matplotlib.png)

```r
library(ggplot2)

df <- read.csv("../data/rape_growth_dynamics.csv")

font_family <- "Heiti TC"

p <- ggplot(df, aes(x = day_after_sowing)) +
  geom_line(aes(y = plant_height_max_cm, color = "最大值")) +
  geom_line(aes(y = plant_height_mean_cm, color = "平均值")) +
  geom_line(aes(y = plant_height_min_cm, color = "最小值")) +
  scale_color_manual(
    name = NULL,
    values = c("最大值" = "#D62728", "平均值" = "#0072B2", "最小值" = "#FFB600")
  ) +
  labs(
    x = "播种后天数 (d)",
    y = "株高 (cm)",
    title = "油菜株高随播种后天数变化"
  ) +
  theme_minimal() +
  theme(
    text = element_text(family = font_family),
    plot.title = element_text(
      hjust = 0.5,
      vjust = 1,
    ),
    legend.position = c(0.99, 0.01),
    legend.justification = c(1, 0),
    legend.box.background = element_rect(
      fill = NA, colour = "black", linewidth = 0.2
    ),
    legend.background = element_blank(),
    panel.border = element_rect(fill = NA, colour = "black", linewidth = 0.8),
    panel.grid.major.x = element_line(
      color = "#C8C8C8", linewidth = 0.3
    ),
    panel.grid.major.y = element_blank(),
    panel.grid.minor = element_blank()
  )

ggsave("plant_height_by_r.png", plot = p, width = 12, height = 6, dpi = 300, units = "cm")
```

* `geom_line` 函数用于绘制线图，`aes` 函数用于指定数据映射；
* `labs` 函数用于添加轴标签、标题；
* `theme` 函数进一步细分绘制要求。

![](/imgs/learn-vis/plant_height_by_r.png)

## 可视化 2 - 叶面积、干物质积累随播种后天数变化

双轴线图，左轴为叶面积，右轴为干物质积累量。

```python
import matplotlib.pyplot as plt
import pandas as pd

df = pd.read_csv("../data/rape_growth_dynamics.csv")

plt.rcParams["font.sans-serif"] = ["Heiti TC"]
plt.rcParams["axes.unicode_minus"] = False

fig, ax = plt.subplots(figsize=(10, 6))
line_color = "#0072B2"
line, = ax.plot(
    df["day_after_sowing"], df["leaf_area_mean_cm2"],
    marker="o", label="叶面积", color=line_color, linewidth=1,
)
ax.set_xlabel("播种后天数 (d)")
ax.set_ylabel("叶面积 (cm²)", color=line_color)
ax.tick_params(axis="y", labelcolor=line_color)
ax.set_ylim(0, max(df["leaf_area_mean_cm2"]) * 1.15)

ax2 = ax.twinx()
line_color2 = "#FFB600"
line2, = ax2.plot(
    df["day_after_sowing"], df["dry_matter_g"], 
    marker="s", label="干物质积累量", color=line_color2, linewidth=1,
)
ax2.set_ylabel("干物质积累量 (g)", color=line_color2)
ax2.tick_params(axis="y", labelcolor=line_color2)
ax2.set_ylim(0, max(df["dry_matter_g"]) * 1.15)

lines = [line, line2]
labels = [ln.get_label() for ln in lines]
ax.legend(lines, labels, loc="lower right", framealpha=0.9)

ax.set_title("叶面积与干物质积累量随生育期的变化")
ax.grid(True, linestyle="--", alpha=0.4, color="#CBC8C8")

fig.savefig("leaf_area_dry_matter_by_matplotlib.png")
```

* `ax.twinx()` 方法创建右侧双轴，共用 x 轴；

![](/imgs/learn-vis/leaf_area_dry_matter_by_matplotlib.png)

```r
library(ggplot2)

df <- read.csv("../data/rape_growth_dynamics.csv")

font_family <- "Heiti TC"

scale <- max(df$leaf_area_mean_cm2) / max(df$dry_matter_g)
df$dry_matter_g_scaled <- df$dry_matter_g * scale

p <- ggplot(df, aes(x = day_after_sowing)) +
  geom_line(aes(y = leaf_area_mean_cm2, color = "叶面积"), linewidth = 1) +
  geom_point(aes(y = leaf_area_mean_cm2, color = "叶面积", shape = "叶面积")) +
  geom_line(aes(y = dry_matter_g_scaled, color = "干物质积累量"), linewidth = 1) +
  geom_point(aes(y = dry_matter_g_scaled, color = "干物质积累量", shape = "干物质积累量")) +
  scale_color_manual(
    name = NULL,
    values = c("叶面积" = "#0072B2", "干物质积累量" = "#FFB600")
  ) +
  scale_shape_manual(
    name = NULL,
    values = c("叶面积" = 19, "干物质积累量" = 15)
  ) +
  scale_y_continuous(
    name = "叶面积 (cm²)",
    sec.axis = sec_axis(~ . / scale, name = "干物质积累量 (g)")
  ) +
  labs(
    title = "叶面积与干物质积累量随生育期的变化",
    x = "播种后天数 (d)"
  ) +
  theme_minimal() +
  theme(
    text = element_text(family = font_family),
    plot.title = element_text(
      hjust = 0.5,
      vjust = 1,
    ),
    axis.title.y = element_text(color = "#0072B2", size = 11),
    axis.text.y = element_text(color = "#0072B2", size = 11),
    axis.title.y.right = element_text(color = "#FFB600", size = 11),
    axis.text.y.right = element_text(color = "#FFB600", size = 11),
    legend.position = c(0.98, 0.04),
    legend.justification = c(1, 0),
    legend.box.background = element_rect(
      fill = NA, colour = "black", linewidth = 0.2
    ),
    legend.background = element_blank(),
    panel.border = element_rect(fill = NA, colour = "black", linewidth = 0.8),
    panel.grid.major = element_line(linetype = 2, color = "#C8C8C8"),
    panel.grid.minor = element_blank()
  )

ggsave("leaf_area_dry_matter_by_r.png", p, width = 9, height = 5.5, dpi = 150)
```

* `scale_shape_manual` 函数用于手动指定形状映射；
* `scale_y_continuous` 函数用于指定 y 轴的刻度标签，`sec.axis` 参数用于指定右侧轴的刻度标签；
* `axis.title.y`、`axis.text.y`、`axis.title.y.right`、`axis.text.y.right` 函数用于指定轴标签、刻度标签的颜色和大小；
* `panel.grid.major` 应用 `element_line` 函数指定主网格线的样式，其中 `linetype = 2` 表示虚线。

![](/imgs/learn-vis/leaf_area_dry_matter_by_r.png)