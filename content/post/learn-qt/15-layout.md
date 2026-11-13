+++
date = '2026-11-12T16:38:35+08:00'
draft = false
title = 'PySide6 布局入门：QVBoxLayout、QHBoxLayout 水平垂直布局'
categories = ['后端技术', 'PySide6']
tags = ['Python', 'PySide6', 'GUI', '布局', 'QVBoxLayout', 'QHBoxLayout']
toc = true
+++

![](/imgs/learn-qt/cover15.png)

## 前言

PySide6 里，如果直接用 `move()`、`resize()` 写死控件坐标，窗口缩放时控件不会跟着自适应，布局类（Layout）就是用来自动管理控件位置、大小、间距的工具。

两个最基础、最常用布局：

- `QHBoxLayout`：水平布局，控件从左往右依次排列
- `QVBoxLayout`：垂直布局，控件从上往下依次排列

> 核心特点：自动自适应窗口大小、支持嵌套、统一设置边距和控件间距。

二者都继承自 **QBoxLayout**，所以下面全部方法通用；区别只在于默认排列方向不同。

![](/imgs/learn-qt/img-63.png)

按功能进行分类，简要介绍如下：

**末尾追加类（add）**

- addItem：向布局末尾添加底层 QLayoutItem 布局项
- addLayout：向布局末尾嵌套添加子布局，可设置拉伸权重
- addSpacerItem：向布局末尾添加自定义 QSpacerItem 空白占位项
- addSpacing：向布局末尾添加固定像素大小的空白间隔，不可拉伸
- addStretch：向布局末尾添加可弹性伸缩空白占位，窗口放大自动分配多余空间
- addStrut：强制限制布局在垂直方向的最小尺寸（水平布局限制最小高度，垂直布局限制最小宽度）
- addWidget：向布局末尾添加控件，可设置拉伸权重、对齐方式

**指定位置插入类（insert）**

- insertItem：在指定索引位置插入底层 QLayoutItem
- insertLayout：在指定索引位置插入子布局
- insertSpacerItem：在指定索引位置插入自定义空白占位项
- insertSpacing：在指定索引位置插入固定像素间隔
- insertStretch：在指定索引位置插入弹性伸缩占位
- insertWidget：在指定索引位置插入控件，可设置拉伸权重、对齐方式

**读取、查询类**

- count：返回布局内所有子项总数（控件、布局、占位都算一项）
- direction：获取当前布局的排列方向
- expandingDirections：返回布局支持拉伸扩展的方向
- hasHeightForWidth：判断布局是否支持 “高度随宽度变化” 的计算逻辑
- heightForWidth：给定宽度，计算布局对应的推荐高度（文本自动换行场景使用）
- itemAt：根据索引获取布局项，**不会从布局移除该项**
- maximumSize：获取布局的最大推荐尺寸
- minimumHeightForWidth：给定宽度，计算布局可接受的最小高度
- minimumSize：获取布局的最小推荐尺寸
- sizeHint：获取布局的推荐默认尺寸
- spacing：获取当前设置的控件之间默认间距
- stretch：获取指定索引子项的拉伸权重

**设置、修改类**

- setDirection：修改布局排列方向
- setGeometry：设置布局所在矩形区域的位置与大小
- setSpacing：设置布局内相邻子项之间的默认间距
- setStretch：根据索引直接设置子项的拉伸权重
- setStretchFactor：给控件 / 子布局设置拉伸权重

**刷新与移除类**

- invalidate：标记布局失效，触发布局重新计算、刷新
- takeAt：按索引取出布局项，**同时从布局中移除该子项**，返回取出的布局项

## QVBoxLayout

`QVBoxLayout` 简单示例：

```python
import sys

from PySide6.QtWidgets import QApplication, QWidget, QVBoxLayout, QPushButton, QLineEdit, QLabel


class MyWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QVBoxLayout 垂直布局")
        self.resize(300, 200)

        v_layout = QVBoxLayout()
        v_layout.setSpacing(10)
        v_layout.setContentsMargins(20, 20, 20, 20)

        v_layout.addWidget(QLabel("用户名"))
        v_layout.addWidget(QLineEdit())
        v_layout.addWidget(QPushButton("登录"))

        self.setLayout(v_layout)


app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

控件按垂直方向排列，间距为 10 像素，边距为 20 像素。

![](/imgs/learn-qt/img-64.png)

## QHBoxLayout

`QHBoxLayout` 简单示例：

```python
import sys

from PySide6.QtWidgets import QApplication, QWidget, QHBoxLayout, QPushButton

class MyWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QHBoxLayout 水平布局")
        self.resize(400, 120)

        h_layout = QHBoxLayout()
        h_layout.setSpacing(12)
        h_layout.setContentsMargins(20, 20, 20, 20)

        h_layout.addWidget(QPushButton("上一步"))
        h_layout.addWidget(QPushButton("下一步"))
        h_layout.addWidget(QPushButton("完成"))

        self.setLayout(h_layout)


app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

控件按水平方向排列，间距为 12 像素，边距为 20 像素。

![](/imgs/learn-qt/img-65.png)

## 复合布局

复合布局是指将多个布局组合起来使用，比如一个垂直布局里嵌套一个水平布局，或一个水平布局里嵌套一个垂直布局，或一个垂直布局里嵌套多个水平布局，或一个水平布局里嵌套多个垂直布局。

```python
import sys

from PySide6.QtWidgets import (
    QApplication, QWidget, QVBoxLayout, QHBoxLayout,
    QLabel, QLineEdit, QPushButton, QTextEdit
)


class CompoundLayoutDemo(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("复合布局：QVBoxLayout + QHBoxLayout 嵌套")
        self.resize(500, 320)

        root_v_layout = QVBoxLayout()
        root_v_layout.setContentsMargins(24, 24, 24, 24)
        root_v_layout.setSpacing(16)

        h1 = QHBoxLayout()
        h1.setSpacing(10)
        label_user = QLabel("用户名：")
        edit_user = QLineEdit()
        h1.addWidget(label_user)
        h1.addWidget(edit_user)
        h1.setStretchFactor(edit_user, 1)
        root_v_layout.addLayout(h1)

        h2 = QHBoxLayout()
        h2.setSpacing(10)
        label_pwd = QLabel("密　码：")
        edit_pwd = QLineEdit()
        edit_pwd.setEchoMode(QLineEdit.EchoMode.Password)
        h2.addWidget(label_pwd)
        h2.addWidget(edit_pwd)
        h2.setStretchFactor(edit_pwd, 1)
        root_v_layout.addLayout(h2)

        label_note = QLabel("备注信息：")
        root_v_layout.addWidget(label_note)
        text_note = QTextEdit()
        root_v_layout.addWidget(text_note, stretch=1)

        root_v_layout.addStretch(1)

        h_btn = QHBoxLayout()
        h_btn.addStretch(1)
        btn_reset = QPushButton("重置")
        btn_submit = QPushButton("提交")
        h_btn.addWidget(btn_reset)
        h_btn.addSpacing(12)
        h_btn.addWidget(btn_submit)
        root_v_layout.addLayout(h_btn)

        self.setLayout(root_v_layout)


app = QApplication(sys.argv)
window = CompoundLayoutDemo()
window.show()
sys.exit(app.exec())
```

![](/imgs/learn-qt/img-66.png)

整体是 **1 个外层垂直布局（QVBoxLayout）作为根布局，里面嵌套多个水平布局（QHBoxLayout）**：

1. **root_v_layout 根布局（垂直）**

所有元素按从上到下顺序摆放，是整个窗口的基础骨架。

- 设置了统一的四周边距 `setContentsMargins`，让内容不和窗口边框贴死；
- 设置统一 `setSpacing`，控制每一行之间的距离。

2. **h1：用户名行（水平布局）**

同一行内，标签和输入框左右并排。

- `addWidget` 依次放入标签、输入框；
- `setStretchFactor(edit_user,1)`：让输入框抢占这一行多余的宽度，窗口拉大时标签宽度不变、输入框自动拉长。

3. **h2：密码行（水平布局）**

和 h1 结构完全一样，标签 + 密码输入框横向排列，输入框自适应拉伸。

4. **备注标签 + 备注文本框（直接放在根垂直布局里，不嵌套水平布局）**

- 标签单独一行；
- 文本框 `addWidget(text_note, stretch=1)`：给文本框设置纵向拉伸权重，窗口拉高时，文本框优先变高。

5. **root_v_layout.addStretch (1) 纵向弹性占位**

这个是纵向空白弹性块，作用是瓜分根布局里垂直方向多余空间，把下面的按钮栏整体**压到窗口底部**。

6. **h_btn 底部按钮行（水平布局）**

同一行横向排布按钮：

- 先 `addStretch(1)`：水平弹性占位，抢占这一行左边所有空余宽度，把两个按钮推到这一行最右侧；
- `addSpacing(12)`：在两个按钮中间插入**固定像素间隔**，和布局的 spacing 互相独立；
- 依次添加重置、提交按钮。

## 小结

QHBoxLayout、QVBoxLayout 本质是同一套 QBoxLayout，仅默认排列方向不同。在复合界面开发中，二者通过互相嵌套组合，就能搭建出绝大多数窗体骨架；配合 `addStretch`、`setStretchFactor`、边距与间距，轻松实现控件自适应拉伸、对齐，避免硬编码坐标，窗口缩放时界面自动适配。