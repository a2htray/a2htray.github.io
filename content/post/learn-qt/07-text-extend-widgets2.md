+++
date = '2026-11-05T17:06:54+08:00'
title = 'PySide6 扩展文本控件（二）：QKeySequenceEdit、QFontComboBox、QGraphicsTextItem'
categories = ['后端技术', 'PySide6']
tags = ['Python', 'PySide6', 'GUI', '文本控件', '文本扩展控件']
toc = true
+++

![](/imgs/learn-qt/cover7.png)

## 另外 3 个文本扩展控件

承接上文 [PySide6 扩展文本控件（一）：QTextBrowser、QLCDNumber](/post/learn-qt/06-text-extend-widgets/)，本文介绍另外 3 个文本扩展控件。

`QKeySequenceEdit`、`QFontComboBox` 两个控件本身并不能处理文本，属于文本操作中需要的辅助控件。

- `QKeySequenceEdit`：不处理文字内容；它捕获按键快捷键，绑定到文本编辑器动作（复制、粘贴、加粗、保存）
- `QFontComboBox`：用于选择字体，支持设置字体、字体大小、字体样式等

`QGraphicsTextItem` **内部直接复用 QTextDocument**，和 QTextEdit 是同一套富文本渲染内核；支持 HTML、段落、文字选中、文本编辑等功能，但只能跑在 GraphicsView 画布坐标系里，可以**旋转缩放**。

## QKeySequenceEdit

专门用来让用户录制快捷键（Ctrl+C、Alt+S 这类），捕获键盘序列，不是普通文字输入。

```python
from PySide6.QtWidgets import QApplication, QWidget, QVBoxLayout, QKeySequenceEdit, QPlainTextEdit, QLabel
import sys

class KeyShortcutDemo(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QKeySequenceEdit 快捷键录入")
        self.resize(350, 150)
        layout = QVBoxLayout(self)


        self.ks_edit = QKeySequenceEdit()
        self.edit_process = QPlainTextEdit("快捷键录入过程：")
        self.lbl_final = QLabel("最终快捷键：")
        self.ks_edit.keySequenceChanged.connect(self.on_key_changed)
        self.ks_edit.editingFinished.connect(self.on_editing_finished)
        layout.addWidget(self.ks_edit)
        layout.addWidget(self.edit_process)
        layout.addWidget(self.lbl_final)

    def on_key_changed(self, ks):
        self.edit_process.appendPlainText(ks.toString())

    def on_editing_finished(self):
        self.lbl_final.setText(f"录制完成，快捷键：{self.ks_edit.keySequence().toString()}")


app = QApplication(sys.argv)
window = KeyShortcutDemo()
window.show()
sys.exit(app.exec())
```

`KeyShortcutDemo` 控件中包含 `QKeySequenceEdit`、`QPlainTextEdit`、`QLabel` 三个控件：

- `QKeySequenceEdit` 用于录入快捷键
- `QPlainTextEdit` 用于显示快捷键录入过程
- `QLabel` 用于显示最终快捷键

运行效果如下：

![](/imgs/learn-qt/img-23.png)

`keySequenceChanged 与 editingFinished 信号的区别：

- 参数
    - `keySequenceChanged`：最新的 `QKeySequence` 对象
    - `editingFinished`：无参数，需要手动调用 `.keySequence()` 取值
- 触发时机
    - `keySequenceChanged`：序列一变化就触发
    - `editingFinished`：一轮录制结束后，只触发一次
- 场景
    - `keySequenceChanged`：记录过程
    - `editingFinished`：保存最终结果

## QFontComboBox

内置系统字体列表，属于特殊文本相关控件，选中即触发字体变化。

```python
import sys

from PySide6.QtWidgets import QApplication, QWidget, QVBoxLayout, QFontComboBox, QLabel


class FontDemo(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QFontComboBox 字体选择")
        self.resize(350, 120)
        layout = QVBoxLayout(self)

        font_cb = QFontComboBox()
        self.tip = QLabel("当前字体：默认")
        font_cb.currentFontChanged.connect(self.on_font_changed)
        layout.addWidget(font_cb)
        layout.addWidget(self.tip)

    def on_font_changed(self, font):
        self.tip.setFont(font)
        self.tip.setText(f"当前字体：{font.family()}")


app = QApplication(sys.argv)
window = FontDemo()
window.show()
sys.exit(app.exec())
```

选择字体后，文本框以该字段显示字体名称，主要方法如下：

- currentFontChanged(QFont)：选择字体时触发，参数是 QFont
- font.family()：取字体族名

运行效果如下：

![](/imgs/learn-qt/img-24.png)

![](/imgs/learn-qt/img-25.png)


## QGraphicsTextItem

`QGraphicsTextItem` 属于 `GraphicsView（QtGraphics）`体系，不是普通窗口控件 QWidget。在 QGraphicsScene 里的富文本图元，可缩放、旋转、移动。

```python
import sys

from PySide6.QtWidgets import (
    QApplication, QGraphicsView, QGraphicsScene, QGraphicsTextItem
)


class GraphicsTextDemo(QGraphicsView):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QGraphicsTextItem 画布文字")
        self.resize(500, 300)

        scene = QGraphicsScene(self)
        text_item = QGraphicsTextItem("<h2>静夜思</h2><p>床前明月光</p>")
        text_item.setPos(50, 50)
        text_item.setRotation(15)
        text_item.setScale(1.5)
        scene.addItem(text_item)
        self.setScene(scene)


app = QApplication(sys.argv)
window = GraphicsTextDemo()
window.show()
sys.exit(app.exec())
```

- `setPos`：设置文本位置
- `setRotation`：设置文本角度
- `setScale`：设置文本缩放比例

运行效果如下：

![](/imgs/learn-qt/img-26.png)

## 小结

`QKeySequenceEdit`、`QFontComboBox`、`QGraphicsTextItem` 3 个控件是“文本生态”的一员，或控制或另一种形式的文本展示，搭配核心的文本控件以满足不同的需求。
