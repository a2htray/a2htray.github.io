+++
date = '2026-11-16T14:05:20+08:00'
draft = false
title = 'PySide6 QFormLayout 表单布局'
categories = ['后端技术', 'PySide6']
tags = ['Python', 'PySide6', 'GUI', '布局', 'QFormLayout']
toc = true
+++

![](/imgs/learn-qt/cover17.png)

在先前的文章中，学习了 QVBoxLayout、QHBoxLayout、QGridLayout 3 个布局类，今天来学习 QFormLayout，这是一个专门用于表单布局的布局类。

## 前言

`QFormLayout` 就是 Qt 专门为“**标签 + 输入控件**”这种两列表单场景封装的布局管理器，不用手动管理行列，自带标签关联、跨平台对齐规则、自适应换行能力，非常适合做各类工具的参数配置面板。

`QFormLayout` 固定是**多行两列结构**：

- 左列：标签（QLabel），用来描述参数名称
- 右列：字段控件（QLineEdit / QSpinBox / QComboBox / QCheckBox 等输入组件）

以下是一个用户登录表单示例：

```python
import sys

from PySide6.QtWidgets import QApplication, QWidget, QFormLayout, QLineEdit, QHBoxLayout, QPushButton


class MyWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("登录表单")

        form_layout = QFormLayout()
        self.edt_name = QLineEdit()
        form_layout.addRow("用户名", self.edt_name)
        self.edt_password = QLineEdit()
        self.edt_password.setEchoMode(QLineEdit.EchoMode.Password)
        form_layout.addRow("密码", self.edt_password)

        btn_layout = QHBoxLayout()
        btn_login = QPushButton("登录")
        btn_reset = QPushButton("重置")
        btn_layout.addWidget(btn_login)
        btn_layout.addWidget(btn_reset)
        form_layout.addRow(btn_layout)

        self.setLayout(form_layout)

        btn_reset.clicked.connect(self.on_reset)
        btn_login.clicked.connect(lambda: print(f"用户名：{self.edt_name.text()}，密码：{self.edt_password.text()}"))

    def on_reset(self):
        self.edt_name.clear()
        self.edt_password.clear()
        self.edt_name.setFocus()


app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

演示效果如下：

![](/imgs/learn-qt/img-72.png)

- `addRow` 方法：添加一行标签 + 输入控件

## 核心 API

```python
import sys

from PySide6.QtCore import Qt
from PySide6.QtWidgets import (
    QApplication, QWidget, QFormLayout, QLineEdit,
    QSpinBox, QComboBox
)


class FormApiDemo(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QFormLayout 核心API演示")
        self.resize(450, 260)

        form = QFormLayout(self)

        form.setLabelAlignment(Qt.AlignmentFlag.AlignRight)
        form.setFormAlignment(Qt.AlignmentFlag.AlignTop | Qt.AlignmentFlag.AlignLeft)
        form.setHorizontalSpacing(12)                   
        form.setVerticalSpacing(24)                         
        form.setFieldGrowthPolicy(QFormLayout.FieldGrowthPolicy.ExpandingFieldsGrow)
        # form.setRowWrapPolicy(QFormLayout.RowWrapPolicy.WrapAllRows)

        form.addRow("名称", QLineEdit("demo"))
        form.addRow("数量", QSpinBox())
        cb = QComboBox()
        cb.addItems(["选项A", "选项B", "选项C"])
        form.addRow("类型", cb)


app = QApplication(sys.argv)
window = FormApiDemo()
window.show()
sys.exit(app.exec())
```

演示效果如下：

![](/imgs/learn-qt/img-73.png)

`QFormLayout` 除了 `addRow` 方法，还有其他核心 API 方法，用于设置标签对齐、字段宽度、换行策略等：

- `setLabelAlignment` 方法：设置标签对齐方式
- `setFormAlignment` 方法：设置表单对齐方式
- `setHorizontalSpacing` 方法：设置水平间距
- `setVerticalSpacing` 方法：设置垂直间距
- `setFieldGrowthPolicy` 方法：设置字段宽度增长策略
- `setRowWrapPolicy` 方法：设置换行策略
