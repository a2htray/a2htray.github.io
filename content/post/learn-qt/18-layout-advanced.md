+++
date = '2026-11-16T15:57:21+08:00'
draft = false
title = 'PySide6 布局进阶：伸缩因子 stretch、间距、边距、SizePolicy 尺寸策略'
categories = ['后端技术', 'PySide6']
tags = ['Python', 'PySide6', 'GUI', '布局']
toc = true
+++

![](/imgs/learn-qt/cover18.png)

## 前言

初学 PySide6 布局时，只会简单 `addWidget`，界面一缩放就失控：控件挤在一起、空白乱飘、控件不会跟着窗口拉伸、或者拉得过大变形。
布局看起来是 “摆放控件”，本质是**尺寸分配规则**。

核心有四块：

1. `stretch` 伸缩因子：空间怎么分
2. `spacing`：控件之间的间距
3. `contentsMargins`：布局四边的边距
4. `SizePolicy`：控件自身的伸缩策略

## 伸缩因子 stretch

`stretch` 只分配"布局放完所有控件最小尺寸之后剩下的空闲空间"，按比例瓜分这份空白，不是直接平分整个容器宽度/高度。

`stretch` 的基本规则是：

- stretch=0：**不主动抢多余空间**，窗口拉大，它的大小不变
- stretch>0：数值是**相对比例**，数值越大，分配的空闲空间越多

相关方法有：

- layout.addWidget(widget, stretch=0)
- layout.addLayout(sub_layout, stretch=0)
- layout.addStretch(stretch=1)
- QGridLayout.setColumnStretch(col_idx, val) / QGridLayout.setRowStretch(row_idx, val)

下面演示一个二行三列的布局（QVBoxLayout 套两个 QHBoxLayout），形如：

![](/imgs/learn-qt/img-74.png)

```python
import sys

from PySide6.QtWidgets import QApplication, QWidget, QVBoxLayout, QHBoxLayout, QLabel, QLineEdit


class MyWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("stretch 演示")
        self.resize(450, 260)

        self.initUI()

    def initUI(self):
        """初始化UI"""
        layout_v = QVBoxLayout(self)
        self.setLayout(layout_v)

        # 第一行
        layout_h0 = QHBoxLayout()
        layout_v.addLayout(layout_h0)

        layout_h0.addWidget(QLabel("服务器"))
        layout_h0.addWidget(QLineEdit(), stretch=1)
        layout_h0.addWidget(QLineEdit())

        # 第二行
        layout_h1 = QHBoxLayout()
        layout_v.addLayout(layout_h1)

        layout_h1.addWidget(QLabel("DSN"))
        layout_h1.addWidget(QLineEdit(), stretch=1)
        layout_h1.addStretch(stretch=1)
        

app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

演示效果如下：

![](/imgs/learn-qt/img-75.png)

当窗口拉大时，控件会根据 `stretch` 来分配空间：

- layout_h0 的第 2 个控件会拉伸，占用所有空闲空间
- layout_h0 的第 3 个控件大小不变
- layout_h1 的第 2 个控件会拉伸，占用所有空闲空间
- layout_h1 的第 3 个控件用 `addStretch(stretch=1)` 进行占位

## 间距 spacing

`spacing` 是控件与控件之间的间隔，不包含布局四边和父容器的距离。常用布局 `QVBoxLayout`、`QHBoxLayout`、`QGridLayout` 都有 `setSpacing` 方法来设置间距，单位是像素。

- 水平布局：左右控件之间的横向空隙
- 垂直布局：上下控件之间的纵向空隙

```python
import sys

from PySide6.QtWidgets import (
    QApplication, QWidget, QVBoxLayout, QHBoxLayout, QComboBox, QPushButton, QLabel,
)


class MyWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("不同 spacing 演示")
        self.resize(450, 260)

        root = QVBoxLayout(self)

        top = QHBoxLayout()
        self.combo = QComboBox()
        self.combo.addItems(["0 px", "10 px", "20 px", "40 px", "80 px"])
        self.combo.setCurrentText("20 px")
        top.addWidget(QLabel("控件间距 spacing:"))
        top.addWidget(self.combo)
        top.addStretch()
        root.addLayout(top)

        self.btn_layout = QHBoxLayout()
        self.btn_layout.setSpacing(20)
        for text in ("Btn1", "Btn2", "Btn3"):
            self.btn_layout.addWidget(QPushButton(text), stretch=1)
        root.addLayout(self.btn_layout)

        self.combo.currentTextChanged.connect(self.change_spacing)

    def change_spacing(self, text):
        self.btn_layout.setSpacing(int(text.split()[0]))


app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

演示效果如下：

![](/imgs/learn-qt/img-76.png)

![](/imgs/learn-qt/img-77.png)


## 边距 contentsMargins

`contentsMargins` 是布局四边的边距，单位是像素，方法 `layout.setContentsMargins(left, top, right, bottom)` 的参数分别是：左、上、右、下。

```python
import sys

from PySide6.QtWidgets import (
    QApplication, QWidget, QVBoxLayout, QHBoxLayout, QComboBox, QPushButton, QLabel,
)


class MyWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("不同 contentMargins 演示")
        self.resize(450, 260)

        root = QVBoxLayout(self)

        top = QHBoxLayout()
        self.combo = QComboBox()
        self.combo.addItems(["0 px", "10 px", "20 px", "40 px", "80 px"])
        self.combo.setCurrentText("20 px")
        top.addWidget(QLabel("四周边距 contentMargins:"))
        top.addWidget(self.combo)
        top.addStretch()
        root.addLayout(top)

        self.btn_layout = QHBoxLayout()
        self.change_margins(self.combo.currentText())
        for text in ("Btn1", "Btn2", "Btn3"):
            self.btn_layout.addWidget(QPushButton(text), stretch=1)
        root.addLayout(self.btn_layout)

        self.combo.currentTextChanged.connect(self.change_margins)

    def change_margins(self, text):
        margin = int(text.split()[0])
        self.btn_layout.setContentsMargins(margin, margin, margin, margin)


app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

演示效果如下：

![](/imgs/learn-qt/img-78.png)

![](/imgs/learn-qt/img-79.png)

## 尺寸策略 SizePolicy

`SizePolicy` 是控件在空间多余或不足的情况下，自身如何应用伸缩行为规则，常用的枚举值有：

- `Policy.Fixed`：尺寸固定，不缩放
- `Policy.Minimum`：最小尺寸是 sizeHint，能放大、不能缩小
- `Policy.Maximum`：最大尺寸是 sizeHint，能缩小、不能放大
- `Policy.Preferred`：优先用 sizeHint，可放大可缩小
- `Policy.Expanding`：同 Preferred，但抢剩余空间优先级更高
- `Policy.MinimumExpanding`：下限是 sizeHint，并且积极抢占多余空间
- `Policy.Ignored`：无视 sizeHint，大小完全交给布局说了算

```python
import sys

from PySide6.QtWidgets import (
    QApplication, QWidget, QVBoxLayout, QHBoxLayout, QComboBox, QPushButton, QLabel,
    QSizePolicy,
)
from PySide6.QtCore import QTimer


class MyWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("sizePolicy 演示")
        self.resize(480, 260)

        root = QVBoxLayout(self)

        top = QHBoxLayout()
        self.combo = QComboBox()
        self.policies = {
            "Fixed": QSizePolicy.Policy.Fixed,
            "Minimum": QSizePolicy.Policy.Minimum,
            "Maximum": QSizePolicy.Policy.Maximum,
            "Preferred": QSizePolicy.Policy.Preferred,
            "Expanding": QSizePolicy.Policy.Expanding,
            "MinimumExpanding": QSizePolicy.Policy.MinimumExpanding,
            "Ignored": QSizePolicy.Policy.Ignored,
        }
        self.combo.addItems(self.policies)
        self.combo.setCurrentText("Preferred")
        top.addWidget(QLabel("水平 sizePolicy:"))
        top.addWidget(self.combo)
        top.addStretch()
        root.addLayout(top)

        self.btn_layout = QHBoxLayout()
        self.buttons = [QPushButton(t) for t in ("Btn1", "Btn2", "Btn3")]
        for btn in self.buttons:
            self.btn_layout.addWidget(btn)
        root.addLayout(self.btn_layout)

        self.info = QLabel()
        root.addWidget(self.info)

        self.combo.currentTextChanged.connect(self.change_policy)
        self.change_policy(self.combo.currentText())

    def change_policy(self, text):
        policy = QSizePolicy(self.policies[text], QSizePolicy.Fixed)
        for btn in self.buttons:
            btn.setSizePolicy(policy)
        QTimer.singleShot(0, self.show_widths)

    def show_widths(self):
        widths = [btn.width() for btn in self.buttons]
        self.info.setText(f"当前 {self.combo.currentText()} 下三个按钮宽度: {widths}")


app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

演示效果如下：

![](/imgs/learn-qt/img-80.png)

![](/imgs/learn-qt/img-81.png)

![](/imgs/learn-qt/img-82.png)

## 小结

- `stretch`：布局分配多余空间的比例规则
- `spacing`：控件之间的距离
- `contentsMargins`：布局四边空白
- `SizePolicy`：控件尺寸拉伸规则
