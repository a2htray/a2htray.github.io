+++
date = '2026-11-06T16:20:42+08:00'
draft = false
title = 'PySide6 按钮类控件（一）：QPushButton、QRadioButton'
categories = ['后端技术', 'PySide6']
tags = ['Python', 'PySide6', 'GUI', '按钮控件']
toc = true
+++

![](/imgs/learn-qt/cover8.png)

## 前言

按钮承担着 GUI 程序里最基础的交互职责：接收用户的点击，触发对应的事件，完成相应的业务。PySide6 中所有按钮都继承自同一个基类，

```bash
QWidget
    - QAbstractButton
        - QPushButton
        - QRadioButton
        - QCheckBox
        - QToolButton
```

如下图：

![](/imgs/learn-qt/img-27.png)

基类提供的通用接口如下：

| 方法 | 描述 |
| --- | --- |
| `setText / text` | 设置和获取按钮的文本内容 |
| `setIcon / icon` | 设置和获取按钮的图标 |
| `setIconSize / iconSize` | 设置和获取按钮的图标大小 |
| `setShortcut / shortcut` | 设置和获取按钮的快捷键 |
| `setChecked / isChecked` | 设置和获取按钮的选中状态 |
| `setCheckable / isCheckable` | 设置按钮是否可选中 |

当然，还有很多，可进入基类源码查看。

本文主要介绍 `QPushButton` 和 `QRadioButton` 两个按钮控件：

- `QPushButton`：普通按钮，执行动作，触发事件
- `QRadioButton`：单选按钮，选择一个选项，维护单选状态

## QPushButton

QPushButton 点击后执行一次操作，不维持选中状态。它是“提交、打开、保存、删除”这类动作的交互入口。

```python
import sys

from PySide6.QtWidgets import QApplication, QMainWindow, QWidget, QVBoxLayout, QPushButton, QLabel


class MyWindow(QMainWindow):
    def __init__(self):
            super().__init__()
            self.setWindowTitle("QPushButton 演示")
            self.resize(400, 300)

            widget = QWidget()
            self.setCentralWidget(widget)
            layout = QVBoxLayout(widget)

            self.btn = QPushButton()
            self.btn.setText("点击我")
            self.btn.clicked.connect(self.on_btn_clicked)
            layout.addWidget(self.btn)

            self.btn2 = QPushButton("无边框与背景")
            self.btn2.setFlat(True)
            self.btn2.clicked.connect(self.on_btn2_clicked)
            layout.addWidget(self.btn2)

    def on_btn_clicked(self):
        self.btn.setText("按钮被点击了")

    def on_btn2_clicked(self):
         flat = self.btn2.isFlat()
         self.btn2.setFlat(not flat)
         self.btn2.setText("无边框与背景" if not flat else "有边框与背景")


app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

演示效果如下：

![](/imgs/learn-qt/img-28.png)

`QPushButton` 提供了 `clicked` 信号，用于响应按钮点击事件，其它信号还有：

- `pressed`：按钮被按下时触发
- `released`：按钮被松开时触发
- `toggleed`：按钮被切换时触发，仅在 checkable 为 True 时有效

> clicked 与 pressed 的差别不只是时机——按下后把鼠标移出按钮再松开，clicked 不会触发，pressed 会触发。因此“一旦按下就不可撤销”的操作，必须用 clicked，避免误触。

更复杂些的示例，有图标设置，快捷键设置、下拉菜单等，如下：

```python
import sys

from PySide6.QtWidgets import QApplication, QWidget, QPushButton, QVBoxLayout, QMenu
from PySide6.QtGui import QIcon, QKeySequence
from PySide6.QtCore import Slot, QSize


class ToolWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("图标与菜单按钮")
        self.resize(320, 220)
        layout = QVBoxLayout(self)

        self.save_btn = QPushButton("保存")
        self.save_btn.setIcon(QIcon.fromTheme("document-save"))
        self.save_btn.setIconSize(QSize(18, 18))
        self.save_btn.setShortcut(QKeySequence("Ctrl+S"))
        self.save_btn.clicked.connect(self.on_save)
        layout.addWidget(self.save_btn)

        menu_btn = QPushButton("更多")
        menu = QMenu(menu_btn)
        menu.addAction("导出为 CSV", self.on_export)
        menu.addAction("打开设置", self.on_settings)
        menu_btn.setMenu(menu)
        layout.addWidget(menu_btn)

    @Slot()
    def on_save(self):
        print("保存")

    @Slot()
    def on_export(self):
        print("导出")

    @Slot()
    def on_settings(self):
        print("设置")


app = QApplication(sys.argv)
window = ToolWindow()
window.show()
sys.exit(app.exec())
```
演示效果如下：

![](/imgs/learn-qt/img-29.png)

## QRadioButton

QRadioButton 是单元框，用于表达“一组互斥选项中的一个”，其核心特性是**互斥**：同组内同一时间只能有一个被选中。

> 普通按钮产生一个事件，单选框维持一个状态。

QRadioButton 的互斥规则是：

- 父控件规则：autoExclusive 默认为 True 的单选框，会与同一父控件下的同级单选框互斥。
- 分组规则：加入同一个 QButtonGroup 的单选框彼此互斥，与父控件无关；即使分布在不同的布局、不同的容器里，也依然互斥。

下一个示例演示：同一父控件下的单选框互斥与默认选中，如下：

```python
import sys

from PySide6.QtWidgets import QApplication, QWidget, QRadioButton, QLabel, QVBoxLayout, QHBoxLayout
from PySide6.QtCore import Slot


class DemoWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QRadioButton 演示")
        self.resize(300, 220)
        layout = QVBoxLayout(self)

        widget = QWidget()
        btn_layout = QHBoxLayout()
        widget.setLayout(btn_layout)
        layout.addWidget(widget)

        self.r1 = QRadioButton("选项 A")
        self.r2 = QRadioButton("选项 B")
        self.r3 = QRadioButton("选项 C")
        self.r1.setChecked(True)

        self.r1.toggled.connect(self.on_toggle)
        self.r2.toggled.connect(self.on_toggle)
        self.r3.toggled.connect(self.on_toggle)

        btn_layout.addWidget(self.r1)
        btn_layout.addWidget(self.r2)
        btn_layout.addWidget(self.r3)

        self.tip = QLabel("当前选中：选项 A")
        layout.addWidget(self.tip)

    @Slot(bool)
    def on_toggle(self, checked):
        if checked:
            self.tip.setText(f"当前选中：{self.sender().text()}")


app = QApplication(sys.argv)
window = DemoWindow()
window.show()
sys.exit(app.exec())
```

- 3 个 QRadioButton 互斥，只能选择一个，它们的父控件是 widget，布局是 QHBoxLayout
- 默认选中的是选项 A：`self.r1.setChecked(True)`

演示效果如下：

![](/imgs/learn-qt/img-30.png)

QButtonGroup 是**纯逻辑容器**：不渲染、不可见、不能放进任何布局，只负责把若干按钮绑定成一组并管理它们的 id。

如：

```python
import sys

from PySide6.QtWidgets import (QApplication, QWidget, QRadioButton,
                               QVBoxLayout, QHBoxLayout, QButtonGroup)
from PySide6.QtCore import Slot


class DemoWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QRadioButton 多分组")
        self.resize(350, 260)
        root_layout = QVBoxLayout(self)

        g1_layout = QHBoxLayout()
        self.gender_group = QButtonGroup(self)
        rb_male = QRadioButton("男")
        rb_female = QRadioButton("女")
        self.gender_group.addButton(rb_male, 1)
        self.gender_group.addButton(rb_female, 2)
        rb_male.setChecked(True)
        g1_layout.addWidget(rb_male)
        g1_layout.addWidget(rb_female)
        root_layout.addLayout(g1_layout)

        g2_layout = QHBoxLayout()
        self.level_group = QButtonGroup(self)
        rb_l1 = QRadioButton("初级")
        rb_l2 = QRadioButton("中级")
        rb_l3 = QRadioButton("高级")
        self.level_group.addButton(rb_l1, 10)
        self.level_group.addButton(rb_l2, 20)
        self.level_group.addButton(rb_l3, 30)
        rb_l2.setChecked(True)
        g2_layout.addWidget(rb_l1)
        g2_layout.addWidget(rb_l2)
        g2_layout.addWidget(rb_l3)
        root_layout.addLayout(g2_layout)

        self.gender_group.idToggled.connect(self.on_gender_toggle)
        self.level_group.idToggled.connect(self.on_level_toggle)

    @Slot(int, bool)
    def on_gender_toggle(self, bid, checked):
        if checked:
            print(f"性别ID：{bid}")

    @Slot(int, bool)
    def on_level_toggle(self, bid, checked):
        if checked:
            print(f"等级ID：{bid}")


app = QApplication(sys.argv)
window = DemoWindow()
window.show()
sys.exit(app.exec())
```

演示效果如下：

![](/imgs/learn-qt/img-31.png)

## 小结

演示了 QPushButton 与 QRadioButton 的用法，比较基础。
