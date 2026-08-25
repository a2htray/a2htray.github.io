+++
date = '2025-12-23T14:35:28+08:00'
draft = false
title = '使用 Excel 做线性回归分析'
categories = ['算法', '机器学习']
tags = ['Excel', '数据分析', '线性回归', '解决问题']
+++

起因是我在刷短视频的时候，看到了一个关于 **Excel 比赛**的视频，那些选手真是把 Excel 玩出了花。所以也想着，以后的工作中，尽可能使用 Excel 完成一些简单的数据分析工作。

这不今天就有那么一个任务。

## 背景

对接了一个文件上传的服务，需要测试文件上传速率。因为这个文件上传接口对文件大小有限制，即不能上传大文件，所以我想着能不能在小文件上多采几个样，通过预测的方法把大文件的上传给估算出来。

## 收集数据

数据的方法很简单，也很笨。

用 `dd` 命令生成 1M 到 20M 的文件，对每个文件执行 5 次上传操作。

然后，复制到 Excel 中，最终于得到的结果如下：

![](/imgs/do-data-analysis-and-linear-regression-with-excel/data.png)

## 绘制散点图

选择 A 列和 I 列，插入散点图：

![](/imgs/do-data-analysis-and-linear-regression-with-excel/insert-scatter.png)

得到下图：

![](/imgs/do-data-analysis-and-linear-regression-with-excel/scatter.png)

## 添加线性回归线

选中一个散点，右键，选择 `Add Trendline...`

![](/imgs/do-data-analysis-and-linear-regression-with-excel/add-trendline.png)

在 Trendline Options 中选择如下：

![](/imgs/do-data-analysis-and-linear-regression-with-excel/trendline-options.png)

* linear 表示线性回归
* Display Equation on chart 显示公式
* Display R-squared value on chart 显示 R 方

得到：

![](/imgs/do-data-analysis-and-linear-regression-with-excel/basic-plot.png)

再加一点点样式，最终结果如下：

![](/imgs/do-data-analysis-and-linear-regression-with-excel/final-plot.png)

## 小结

以后，有什么数据图 📊 📈 📉 都优先使用 Excel 画一下，做个验证。