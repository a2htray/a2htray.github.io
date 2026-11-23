+++
date = '2026-11-18T11:08:12+08:00'
draft = false
title = 'PySide6 QSS 入门：基础语法、选择器、颜色、边框、圆角'
categories = ['后端技术', 'PySide6']
tags = ['Python', 'PySide6', 'GUI', 'QSS', '样式']
toc = true
+++

![](/imgs/learn-qt/cover19.png)

> 如果你有 CSS 的基础，掌握 QSS 会更快速。

## 前言

QSS 全称 Qt Style Sheets，语法借鉴 CSS，用于设置 Qt 应用的样式，包括给 `Qt/PySide6` 控件上色、改边框、做圆角、定义交互悬浮效果等。

本文从基础语法讲起，介绍选择器、颜色、边框、圆角等知识点，同时介绍项目里最常用的**独立 .qss 文件**方案，实现样式与业务代码解耦，方便后续维护、换肤和迭代。

## 基础语法

用规则定义样式，规则的格式如下：

```css
选择器 {
    属性: 值;
}
```

- 选择器：选择要应用样式的控件
- 属性：要应用的样式属性
- 值：属性的具体值

规则使用 `/*`、`*/` 注释，注释内容不会被应用，并严格区分大小写。

应用样式的方式有三：

- 内联样式：直接在控件上设置样式，格式为 `控件.setStyleSheet(样式)`
- 全局样式：在应用启动时设置全局样式，格式为 `window.setStyleSheet(样式)`
- 独立 .qss 文件：将样式写入独立的 .qss 文件，格式为 `window.setStyleSheet("qss 文件内容")`，等同于方式二

## 选择器

QSS 的选择器包括类型选择器、ID 选择器、动态属性选择器、后代选择器和伪类选择器等。

**类型选择器**

```css
/* 所有 QPushButton */
QPushButton { }
/* 所有 QLabel */
QLabel { }
```

**ID 选择器**

```python
btn.setObjectName("submit_btn")
```

```css
#submit_btn {
    background-color: green;
}
```

**动态属性选择器**

```python
btn.setProperty("btn-type", "danger")
```

```css
QPushButton[btn-type="danger"] {
    background-color: #dc2626;
}
```

**后代选择器**

```css
/* QWidget 里面所有 QLabel */
QWidget QLabel { color: #666; }
```

**伪类选择器**

```css
/* 鼠标悬浮 */
QPushButton:hover {}
/* 按下 */
QPushButton:pressed {}
/* 禁用 */
QPushButton:disabled {}
/* 选中（复选框/单选框） */
QCheckBox:checked {}
```

## 颜色

颜色支持 4 种写法：

- 颜色关键字
- 十六进制
- rgb
- rgba 带透明度

```python
import sys

from PySide6.QtWidgets import QApplication, QWidget, QVBoxLayout, QLabel


class MyWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QSS 颜色设置的 4 种写法")
        self.resize(420, 320)

        root = QVBoxLayout(self)
        root.addWidget(QLabel("QSS 中的 color / background-color 支持以下 4 种写法："))

        # 1) 颜色关键字（如 red、green、blue、white...）
        w1 = QLabel("1. 颜色关键字: background-color: lightblue; color: darkblue;")
        w1.setStyleSheet("background-color: lightblue; color: darkblue; padding: 8px;")

        # 2) 十六进制：#RRGGBB（也可 #RGB 简写）
        w2 = QLabel("2. 十六进制: background-color: #FFD700; color: #333333;")
        w2.setStyleSheet("background-color: #FFD700; color: #333333; padding: 8px;")

        # 3) rgb(r, g, b)：每通道 0~255
        w3 = QLabel("3. rgb: background-color: rgb(144, 238, 144); color: rgb(0, 100, 0);")
        w3.setStyleSheet("background-color: rgb(144, 238, 144); color: rgb(0, 100, 0); padding: 8px;")

        # 4) rgba(r, g, b, a)：a 为透明度，0 全透明 ~ 1 不透明
        w4 = QLabel("4. rgba: background-color: rgba(255, 99, 71, 0.4); color: rgba(0, 0, 0, 1);")
        w4.setStyleSheet(
            "background-color: rgba(255, 99, 71, 0.4); color: rgba(0, 0, 0, 1); padding: 8px;"
        )

        for w in (w1, w2, w3, w4):
            root.addWidget(w)


app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

演示效果如下：

![](/imgs/learn-qt/img-83.png)

## 边框

边框设置的格式为 `border: 边框宽度 边框样式 边框颜色;`，也可以单独设置每个边的边框。

```css
PushButton {
    /* 统一边框 */
    border: 2px solid #2563eb;
    /* 单独四边 */
    border-top: 1px solid #ccc;
    border-bottom: 1px solid #ccc;
    border-left: 1px solid #ccc;
    border-right: 1px solid #ccc;
}
```

边框样式支持：

- `solid` 实线（最常用）
- `dashed` 虚线
- `dotted` 点线

```python
import sys

from PySide6.QtWidgets import QApplication, QWidget, QVBoxLayout, QLabel


class MyWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QSS 边框样式演示")
        self.resize(420, 320)

        root = QVBoxLayout(self)
        root.addWidget(QLabel("QSS 的 border-style 支持以下常见样式："))

        # 通用边框属性：宽度 + 样式 + 颜色
        common = "border-width: 2px; border-color: steelblue; padding: 8px;"

        # 1) solid 实线（最常用）
        w1 = QLabel("1. solid 实线: border-style: solid;")
        w1.setStyleSheet(f"{common} border-style: solid;")

        # 2) dashed 虚线
        w2 = QLabel("2. dashed 虚线: border-style: dashed;")
        w2.setStyleSheet(f"{common} border-style: dashed;")

        # 3) dotted 点线
        w3 = QLabel("3. dotted 点线: border-style: dotted;")
        w3.setStyleSheet(f"{common} border-style: dotted;")

        for w in (w1, w2, w3):
            root.addWidget(w)


app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

演示效果如下：

![](/imgs/learn-qt/img-84.png)

## 圆角

圆角设置的格式为 `border-radius: 圆角半径;`，支持设置单个、两个、四个角的圆角半径。

```css
QPushButton {
    /* 一个值统一所有角 */
    border-radius: 8px;
    /* 四个值设置四个角的圆角半径，顺序为：左、上、右、下 */
    border-radius: 4px 8px 12px 16px;
}
```

```python
import sys

from PySide6.QtWidgets import QApplication, QWidget, QVBoxLayout, QLabel


class MyWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QSS border-radius 圆角演示")
        self.resize(420, 360)

        root = QVBoxLayout(self)
        root.addWidget(QLabel("QSS 的 border-radius 常用写法："))

        # 通用底色，便于看清圆角效果
        common = "background-color: #FFD700; padding: 10px; border: 1px solid #999;"

        # 1) 四角统一圆角
        w1 = QLabel("1. 四角统一: border-radius: 10px;")
        w1.setStyleSheet(common + "border-radius: 10px;")

        # 2) 分别指定四角（顺序：左上 右上 右下 左下）
        w2 = QLabel("2. 分别指定四角: border-radius: 2px 10px 2px 10px;")
        w2.setStyleSheet(common + "border-radius: 2px 10px 2px 10px;")

        for w in (w1, w2):
            root.addWidget(w)


app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

演示效果如下：

![](/imgs/learn-qt/img-85.png)

## QSS 文件

QSS 文件是样式的剥离，在应用时通过读取 QSS 文件内容并应用到控件上。

```python
def load_qss(path):
    with open(path, "r", encoding="utf-8") as f:
        return f.read()

qss = load_qss("base.qss") + load_qss("button.qss") + load_qss("input.qss")
app.setStyleSheet(qss)
```

## 小结

QSS 语法借鉴 CSS，上手成本低，核心就三件事：

1. **规则**：`选择器 { 属性: 值; }`，配合类型、ID、动态属性、后代、伪类五类选择器精确定位控件。
2. **常用样式**：颜色（关键字/十六进制/rgb/rgba）、边框（`border: 宽度 样式 颜色`）、圆角（`border-radius`）覆盖了日常 UI 美化的大部分需求。
3. **工程化**：把样式剥离到独立 .qss 文件中统一加载，实现样式与业务代码解耦，便于维护和换肤。

掌握这些基础后，就可以应付绝大多数控件美化场景，更高级的属性（如渐变、子控件定制）可在需要时查阅 Qt 文档按图索骥。

