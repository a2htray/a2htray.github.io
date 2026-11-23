+++
date = '2026-11-23T10:41:06+08:00'
draft = false
title = 'PySide6 QSS 进阶：伪状态 hover/pressed/disabled、子控件 ::indicator'
categories = ['后端技术', 'PySide6']
tags = ['Python', 'PySide6', 'GUI', 'QSS', '样式']
toc = true
+++

![](/imgs/learn-qt/cover20.png)

## 伪状态

伪状态（Pseudo-State）是指在控件的正常状态下，通过鼠标悬停、点击等交互操作，改变控件的外观。

> 伪状态 = 给控件样式加一条 “什么时候生效” 的条件。

常用伪状态清单：

|伪状态|说明|
|---|---|
|:hover|鼠标悬浮在控件上方|
|:pressed|鼠标左键按下还没松开|
|:disabled|setEnabled(False)，控件禁用|
|:enabled|控件可用|
|:checked|QCheckBox / QRadioButton 被勾选|
|:focus|控件拿到输入焦点|
|:!xxx|取反，:!hover = 不在悬浮状态|

```python
import sys

from PySide6.QtWidgets import (
    QApplication, QWidget, QVBoxLayout, QLabel,
    QPushButton, QCheckBox, QRadioButton, QLineEdit,
)


class MyWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QSS 伪状态演示")
        self.resize(520, 480)

        root = QVBoxLayout(self)
        root.addWidget(QLabel("把鼠标移到按钮/按住、点击勾选、切换焦点，观察颜色变化："))

        # :hover —— 鼠标悬浮
        btn_hover = QPushButton("鼠标悬浮到这里 (:hover)")
        btn_hover.setObjectName("btn_hover")

        # :pressed —— 鼠标左键按下未松开
        btn_pressed = QPushButton("按住我不松 (:pressed)")
        btn_pressed.setObjectName("btn_pressed")

        # :disabled / :enabled —— 禁用 / 可用
        btn_disabled = QPushButton("禁用状态 (:disabled)")
        btn_disabled.setObjectName("btn_disabled")
        btn_disabled.setEnabled(False)
        btn_enabled = QPushButton("可用状态 (:enabled)")
        btn_enabled.setObjectName("btn_enabled")

        # :checked —— 勾选状态
        chk = QCheckBox("勾选我 (:checked)")
        chk.setObjectName("chk")
        radio = QRadioButton("选中我 (:checked)")
        radio.setObjectName("radio")

        # :focus —— 获得输入焦点
        line = QLineEdit()
        line.setObjectName("line")
        line.setPlaceholderText("点一下获得焦点 (:focus)")

        # :!hover —— 取反：不在悬浮状态
        btn_not_hover = QPushButton("鼠标移开时变色 (:!hover)")
        btn_not_hover.setObjectName("btn_not_hover")

        for w in (btn_hover, btn_pressed, btn_disabled, btn_enabled,
                  chk, radio, line, btn_not_hover):
            root.addWidget(w)

        self.setStyleSheet("""
            QPushButton#btn_hover:hover        { background: #7ec8ff; }
            QPushButton#btn_pressed:pressed    { background: #ff9f9f; }
            QPushButton#btn_disabled:disabled  { background: #dddddd; color: #999999; }
            QPushButton#btn_enabled:enabled    { background: #b7eb8f; }
            QCheckBox#chk:checked,
            QRadioButton#radio:checked         { color: #d4380d; font-weight: bold; }
            QLineEdit#line:focus               { background: #fffbe6; border: 2px solid #faad14; }
            QPushButton#btn_not_hover:!hover   { background: #ffe7ba; }
        """)


app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

* setObjectName() 方法用于为控件设置一个唯一的名称，用于在 QSS 中引用该控件
* 伪状态的格式为 `控件类型#控件名称:伪状态`，例如 `QPushButton#btn_hover:hover`
* setStyleSheet() 方法用于为控件设置样式表，样式表中可以使用 `setObjectName()` 方法设置的名称

```css
QPushButton#btn_hover:hover        { background: #7ec8ff; }
QPushButton#btn_pressed:pressed    { background: #ff9f9f; }
QPushButton#btn_disabled:disabled  { background: #dddddd; color: #999999; }
QPushButton#btn_enabled:enabled    { background: #b7eb8f; }
QCheckBox#chk:checked,
QRadioButton#radio:checked         { color: #d4380d; font-weight: bold; }
QLineEdit#line:focus               { background: #fffbe6; border: 2px solid #faad14; }
QPushButton#btn_not_hover:!hover   { background: #ffe7ba; }
```

演示效果如下：

![](/imgs/learn-qt/img-86.png)

## 子控件 ::indicator

`::indicator` 是 Qt 预定义的**子控件（sub-control，双冒号 `::`）**，不是伪状态，它指代复选框、单选框内部那个勾选方框 / 圆形单选圆点。

只有这几个常用控件支持 `::indicator`：

- `QCheckBox`：前面的方形勾选框
- `QRadioButton`：前面的圆形单选圆点
- `QGroupBox`（可选勾选模式 `setCheckable(True)`）：标题前面的勾选框

```python
import sys

from PySide6.QtWidgets import (
    QApplication, QWidget, QVBoxLayout, QLabel,
    QCheckBox, QRadioButton, QGroupBox,
)


class MyWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QSS ::indicator 子控件演示")
        self.resize(420, 420)

        root = QVBoxLayout(self)
        root.addWidget(QLabel("用 ::indicator 自定义勾选指示器："))

        # 1) QCheckBox：方形勾选框
        chk = QCheckBox("QCheckBox 方形勾选框")
        chk.setChecked(True)

        # 2) QRadioButton：圆形单选圆点
        radio1 = QRadioButton("QRadioButton 圆形单选圆点（选中）")
        radio2 = QRadioButton("QRadioButton 圆形单选圆点（未选）")
        radio1.setChecked(True)

        # 3) QGroupBox：标题前面的勾选框（需 setCheckable(True)）
        group = QGroupBox("QGroupBox 可勾选标题")
        group.setCheckable(True)
        group.setChecked(True)
        group_layout = QVBoxLayout(group)
        group_layout.addWidget(QLabel("组内的内容"))

        root.addWidget(chk)
        root.addWidget(radio1)
        root.addWidget(radio2)
        root.addWidget(group)

        self.setStyleSheet("""
            /* QCheckBox：方形勾选框 */
            QCheckBox::indicator {
                width: 18px; height: 18px;
                border: 2px solid #888; border-radius: 3px; background: white;
            }
            QCheckBox::indicator:checked {
                background: #52c41a; border-color: #389e0d;
            }

            /* QRadioButton：圆形单选圆点 */
            QRadioButton::indicator {
                width: 18px; height: 18px;
                border: 2px solid #888; border-radius: 9px; background: white;
            }
            QRadioButton::indicator:checked {
                background: #1677ff; border-color: #0958d9;
            }

            /* QGroupBox：标题前的勾选框 */
            QGroupBox::indicator {
                width: 16px; height: 16px;
                border: 2px solid #888; border-radius: 2px; background: white;
            }
            QGroupBox::indicator:checked {
                background: #fa8c16; border-color: #d46b08;
            }
        """)


app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

演示效果如下：

![](/imgs/learn-qt/img-87.png)

## 小结

1. **伪状态 `:`**：给样式加生效条件，`:hover`、`:pressed`、`:disabled`、`:checked`、`:focus` 等覆盖了几乎全部交互反馈场景，`!` 取反可处理“非某状态”的情况。
2. **子控件 `::`**：`::indicator` 用于定制 QCheckBox、QRadioButton、QGroupBox 内部的勾选指示器，配合 `:checked` 可分别定义选中/未选中外观。
