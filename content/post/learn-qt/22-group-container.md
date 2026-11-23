+++
date = '2026-11-23T17:18:59+08:00'
draft = false
title = 'PySide6 QGroupBox、QFrame 分组容器，界面分区'
categories = ['后端技术', 'PySide6']
tags = ['Python', 'PySide6', 'GUI', 'QGroupBox', 'QFrame', '界面分区']
toc = true
+++

![](/imgs/learn-qt/cover22.png)

本文讲解两个最常用的分区容器：**QGroupBox（带标题分组框）**、**QFrame（纯边框分区面板）**，两个控件的适用场景分别是：

- **QGroupBox**：需要分组显示控件，且分组有标题时使用，如功能模块分组、参数面板、设置项分组等。
- **QFrame**：需要分组显示控件，且分组无标题时使用，卡片、分割面板、视觉区块、自定义边框面板等。

## QGroupBox

QGroupBox，是一个带标题的分组框，最核心特色：`checkable` 开关，勾选时自动启用 / 禁用内部所有子控件。

常用的属性包括：

- `title`：标题文字，支持快捷键 `&`
- `alignment`：标题对齐
- `checkable`：是否开启复选框
- `checked`：勾选状态；`checkable=True` 时生效，**取消勾选 → 内部控件全部 disable**
- `flat`：扁平化，隐藏大部分边框，只保留标题分隔感

```python
import sys

from PySide6.QtCore import Qt
from PySide6.QtWidgets import (
    QApplication,
    QCheckBox,
    QFormLayout,
    QGroupBox,
    QLineEdit,
    QPushButton,
    QVBoxLayout,
    QWidget,
)

class MyWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QGroupBox 演示")
        self.resize(360, 220)

        group = QGroupBox("账户设置(&A)")
        group.setAlignment(Qt.AlignmentFlag.AlignRight | Qt.AlignmentFlag.AlignVCenter)
        group.setCheckable(True)
        group.setFlat(True)

        form = QFormLayout(group)
        form.addRow("昵称：", QLineEdit("张三"))
        form.addRow(QCheckBox("接收通知"))
        form.addRow(QPushButton("保存设置"))

        group.setChecked(False)

        root = QVBoxLayout(self)
        root.addWidget(group)


app = QApplication(sys.argv)
app.setStyle("Fusion")
window = MyWindow()
window.show()
sys.exit(app.exec())
```

- `QGroupBox("账户设置(&A)")` 创建分组框
- `setAlignment(...)` 设置标题右对齐并垂直居中，对应 alignment 属性
- `setCheckable(True)` 让标题左侧出现复选框，对应 checkable 属性
- `setFlat(True)` 开启扁平化，隐藏大部分边框，只保留标题分隔感，对应 flat 属性
- `QFormLayout(group)` 把表单布局设置给分组框，内部放入输入框、复选框和按钮
- `group.setChecked(False)` 取消勾选分组框，对应 checked 属性。因为 checkable=True，取消勾选后，QGroupBox 会自动禁用内部所有控件；再次勾选会重新启用。

演示效果如下：

![](/imgs/learn-qt/img-90.png)

## QFrame

`QFrame` 是一个**纯边框分区面板**，常用于卡片、分割面板、视觉区块等场景，它有两个组合属性控制外观：

- `QFrame.Shape`：边框形状
- `QFrame.Shadow`：边框立体效果

**QFrame.Shape**

- `QFrame.Shape.NoFrame`：无边框（默认）
- `QFrame.Shape.Box`：矩形方框
- `QFrame.Shape.Panel`：面板
- `QFrame.Shape.StyledPanel`：跟随系统样式面板
- `QFrame.Shape.HLine`：水平分割线
- `QFrame.Shape.VLine`：垂直分割线

**QFrame.Shadow**

- `QFrame.Shadow.Plain`：平面，无阴影
- `QFrame.Shadow.Raised`：凸起
- `QFrame.Shadow.Sunken`：凹陷

```python
import sys

from PySide6.QtWidgets import (
    QApplication,
    QComboBox,
    QFormLayout,
    QFrame,
    QLabel,
    QVBoxLayout,
    QWidget,
)


class MyWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QFrame 演示")
        self.resize(420, 220)

        self.frame = QFrame()
        self.frame.setMinimumSize(360, 100)
        self.frame.setFrameShape(QFrame.Shape.Box)
        self.frame.setFrameShadow(QFrame.Shadow.Raised)

        self.shape = QComboBox()
        self.shape.addItems(["NoFrame", "Box", "Panel", "StyledPanel", "HLine", "VLine"])
        self.shape.currentTextChanged.connect(self.update_frame)

        self.shadow = QComboBox()
        self.shadow.addItems(["Plain", "Raised", "Sunken"])
        self.shadow.currentTextChanged.connect(self.update_frame)

        form = QFormLayout()
        form.addRow("frameShape：", self.shape)
        form.addRow("frameShadow：", self.shadow)

        root = QVBoxLayout(self)
        root.addWidget(self.frame)
        root.addWidget(QLabel("切换下拉框，观察边框形状与立体效果"))
        root.addLayout(form)

    def update_frame(self):
        shape = getattr(QFrame.Shape, self.shape.currentText())
        shadow = getattr(QFrame.Shadow, self.shadow.currentText())
        self.frame.setFrameShape(shape)
        self.frame.setFrameShadow(shadow)


app = QApplication(sys.argv)
app.setStyle("Fusion")
window = MyWindow()
window.show()
sys.exit(app.exec())
```

关键代码是：

```python
self.frame = QFrame()
self.frame.setMinimumSize(360, 100)
self.frame.setFrameShape(QFrame.Shape.Box)
self.frame.setFrameShadow(QFrame.Shadow.Raised)
```

并且用 `.setFrameShape()` 和 `.setFrameShadow()` 方法设置**边框形状**和**立体效果**。

演示效果如下：

![](/imgs/learn-qt/img-91.png)

![](/imgs/learn-qt/img-92.png)

![](/imgs/learn-qt/img-93.png)

## 小结

总的来说，QGroupBox 是一个**有标题的分组框**，可以控制一组控件是否可用，而 QFrame 是一个**纯边框分区面板**，偏向于视觉工具，用于绘制边框、分割线与卡片面板。
