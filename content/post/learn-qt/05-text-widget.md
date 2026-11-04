+++
date = '2026-11-04T09:39:38+08:00'
draft = false
title = 'PySide6 文本控件：QLabel、QLineEdit、QTextEdit、QPlainTextEdit 使用全解'
categories = ['后端技术', 'PySide6']
tags = ['Python', 'PySide6', 'GUI', '文本控件']
toc = true
+++

![](/imgs/learn-qt/cover5.png)

## 文本控件

PySide6 提供了多种文本控件，用于显示和编辑文本，包括：

- QLabel：用于显示文本，不支持编辑
- QLineEdit：用于编辑文本，支持输入
- QTextEdit：用于编辑多行文本，支持格式化
- QPlainTextEdit：用于编辑多行文本，不支持格式化

**一图全览**

![](/imgs/learn-qt/img-15.png)

文本控件选型：

- 只展示不用输入 → QLabel
- 只能一行输入 → QLineEdit
- 排版富文本 → QTextEdit
- 日志、代码、纯文本 → QPlainTextEdit

## QLabel

窗口界面中的 QLabel，分别显示文本、图片、HTML。

```python
import sys

from PySide6.QtWidgets import QApplication, QMainWindow, QWidget, QVBoxLayout, QLabel
from PySide6.QtGui import QPixmap
from PySide6.QtCore import QSize


class MainWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QLabel 演示")
        self.resize(400, 400)

        widget = QWidget(self)
        layout = QVBoxLayout(widget)
        self.setCentralWidget(widget)

        # 纯文本
        label_text = QLabel("纯文本 QLabel")
        layout.addWidget(label_text)

        # 图片
        pixmap = QPixmap("cover5.png").scaled(QSize(220, 100))
        label_img = QLabel()
        label_img.setPixmap(pixmap)
        layout.addWidget(label_img)

        # HTML 格式
        label_html = QLabel("<b>HTML 格式 QLabel</b>")
        layout.addWidget(label_html)

        # 超链接
        label_link = QLabel("<a href='https://www.pythonguis.com'>这是一个链接</a>")
        # 开启外部链接
        label_link.setOpenExternalLinks(True)
        layout.addWidget(label_link)


app = QApplication(sys.argv)
window = MainWindow()
window.show()
sys.exit(app.exec())
```

运行效果如下：

![](/imgs/learn-qt/img-16.png)

使用 QLabel 需要注意以下几点：

1. QLabel 不支持编辑，无法获取焦点
2. setText 方法设置文本内容或实例化时传入，setPixmap 方法设置图片
3. setPixmap 会覆盖 setText 方法设置的文本内容，只能显示文本或图片之一
4. 图片显示支持 png、jpg、jpeg、gif 等格式，支持缩放、居中显示等操作

## QLineEdit

QLineEdit 用于编辑文本，支持单行输入。

```python
import sys

from PySide6.QtWidgets import QApplication, QMainWindow, QWidget, QVBoxLayout, QLineEdit, QPlainTextEdit
from PySide6.QtCore import Slot


class MainWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QLabel 演示")
        self.resize(400, 200)

        widget = QWidget(self)
        layout = QVBoxLayout(widget)
        self.setCentralWidget(widget)

        # 用户名
        self.line_username = QLineEdit()
        self.line_username.setPlaceholderText("请输入用户名")
        self.line_username.returnPressed.connect(self.username_finish)
        layout.addWidget(self.line_username)

        # 密码
        self.line_password = QLineEdit()
        self.line_password.setPlaceholderText("请输入密码")
        self.line_password.setEchoMode(QLineEdit.EchoMode.Password)
        self.line_password.returnPressed.connect(self.password_finish)
        layout.addWidget(self.line_password)

        # 加个日志，响应文本变化
        self.edit_log = QPlainTextEdit()
        layout.addWidget(self.edit_log, stretch=1)

    @Slot()
    def username_finish(self):
        username = self.line_username.text()
        self.edit_log.appendPlainText(f"用户名: {username}")

    @Slot()
    def password_finish(self):
        password = self.line_password.text()
        self.edit_log.appendPlainText(f"密码: {password}")


app = QApplication(sys.argv)
window = MainWindow()
window.show()
sys.exit(app.exec())
```

使用 QLineEdit 实现了用户名、密码输入框：

- 使用 `setPlaceholderText` 方法设置占位符文本
- `setEchoMode` 方法设置密码输入模式
- `returnPressed` 信号触发，获取输入文本

运行效果如下：

![](/imgs/learn-qt/img-17.png)

使用 QLineEdit 需要注意以下几点：

1. 单行输入，不能换行：回车会触发 returnPressed 信号
2. 密码框必须设 setEchoMode(QLineEdit.EchoMode.Password)，否则明文显示
3. 空输入判断：`xxx.text().strip() == ""`，注意去掉首尾空白
4. 限制长度：`setMaxLength(n)`，n 为最大输入长度

## QTextEdit

QTextEdit 用于编辑多行文本，支持富文本格式。

```python
import sys

from PySide6.QtWidgets import QApplication, QMainWindow, QWidget, QGridLayout, QHBoxLayout, QTextEdit, QPushButton


class MainWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QTextEdit 演示")
        self.resize(800, 400)

        widget = QWidget(self)
        grid_layout = QGridLayout(widget)
        self.setCentralWidget(widget)

        self.edit_text_pure = QTextEdit()
        self.edit_text_pure.append("""《水浒传》又名《忠义水浒传》、《江湖豪客传》，是我国第一部歌颂农民起义的长篇章回体小说，写于元末明初, 故事取材于北宋末年的宋江起义。 它生动地描写了梁山好汉们从起义到兴盛再到最终失败的全过程， 特别是通过写众多草莽英雄不同的人生经历和反抗道路，鲜明地表现了“官逼民反”的主题，是一部反抗封建暴政的英雄传奇。 小说真实地叙写了起义斗争的发展过程：从个人反抗到集体反抗，从分散斗争到有组织斗争，规模从小到大。""")
        grid_layout.addWidget(self.edit_text_pure, 0, 0)

        self.edit_text_html = QTextEdit()
        self.edit_text_html.append("""<h3 style="text-align:center;">静夜思</h3>
<p style="text-align:center; color: blue;">床前明月光，</p>
<p style="text-align:center; font-size: 20px;">疑是地上霜。</p>
<p style="text-align:center; text-decoration: underline;">举头望明月，</p>
<p style="text-align:center;">低头思故乡。</p>""")
        grid_layout.addWidget(self.edit_text_html, 0, 1)

        btns_widget = QWidget()
        grid_layout.addWidget(btns_widget, 1, 0, 1, 2)
        btns_layout = QHBoxLayout(btns_widget)

        self.btn_can_edit = QPushButton("禁用编辑")
        self.btn_can_edit.clicked.connect(self.toggle_readonly)
        btns_layout.addWidget(self.btn_can_edit)

        self.btn_clear = QPushButton("清空内容")
        self.btn_clear.clicked.connect(self.clear_content)
        btns_layout.addWidget(self.btn_clear)

        self.btn_get_content = QPushButton("获取内容")
        self.btn_get_content.clicked.connect(self.get_content)
        btns_layout.addWidget(self.btn_get_content)

    def toggle_readonly(self):
        self.edit_text_pure.setReadOnly(not self.edit_text_pure.isReadOnly())
        self.edit_text_html.setReadOnly(not self.edit_text_html.isReadOnly())
        self.btn_can_edit.setText("启用编辑" if self.edit_text_pure.isReadOnly() else "禁用编辑")

    def clear_content(self):
        self.edit_text_pure.clear()
        self.edit_text_html.clear()

    def get_content(self):
        print(self.edit_text_pure.toPlainText())
        print(self.edit_text_html.toHtml())

app = QApplication(sys.argv)
window = MainWindow()
window.show()
sys.exit(app.exec())
```

- 使用 `QGridLayout` 布局，将 QTextEdit 放置在不同的列中
- 调用 `QTextEdit.append` 方法添加文本内容，支持纯文本、HTML
- `QTextEdit.setReadOnly` 方法设置文本内容为只读模式，通过 `QTextEdit.isReadOnly` 方法判断是否只读
- `QTextEdit.clear` 方法清空文本内容
- `QTextEdit.toPlainText` 方法获取纯文本内容
- `QTextEdit.toHtml` 方法获取 HTML 内容

运行效果如下：

![](/imgs/learn-qt/img-18.png)

使用 QTextEdit 需要注意以下几点：

1. 获取内容的两种方式：toHtml() 返回带标签的完整 HTML；toPlainText() 返回纯文本；
2. `append` 方法自动换行追加到末尾
3. `textChanged` 每次编辑便会触发

## QPlainTextEdit

QPlainTextEdit 用于编辑多行文本，不支持富文本格式。

```python
import sys

from PySide6.QtWidgets import QApplication, QMainWindow, QWidget, QGridLayout, QHBoxLayout, QPlainTextEdit, QPushButton


class MainWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QPlainTextEdit 演示")
        self.resize(800, 400)

        widget = QWidget(self)
        grid_layout = QGridLayout(widget)
        self.setCentralWidget(widget)

        self.edit_text_pure = QPlainTextEdit()
        self.edit_text_pure.appendPlainText("""《水浒传》又名《忠义水浒传》、《江湖豪客传》，是我国第一部歌颂农民起义的长篇章回体小说，写于元末明初, 故事取材于北宋末年的宋江起义。 它生动地描写了梁山好汉们从起义到兴盛再到最终失败的全过程， 特别是通过写众多草莽英雄不同的人生经历和反抗道路，鲜明地表现了“官逼民反”的主题，是一部反抗封建暴政的英雄传奇。 小说真实地叙写了起义斗争的发展过程：从个人反抗到集体反抗，从分散斗争到有组织斗争，规模从小到大。""")
        grid_layout.addWidget(self.edit_text_pure, 0, 0)

        self.edit_text_pure2 = QPlainTextEdit()
        self.edit_text_pure2.appendPlainText("""<h3 style="text-align:center;">静夜思</h3>
<p style="text-align:center; color: blue;">床前明月光，</p>
<p style="text-align:center; font-size: 20px;">疑是地上霜。</p>
<p style="text-align:center; text-decoration: underline;">举头望明月，</p>
<p style="text-align:center;">低头思故乡。</p>""")
        grid_layout.addWidget(self.edit_text_pure2, 0, 1)

        btns_widget = QWidget()
        grid_layout.addWidget(btns_widget, 1, 0, 1, 2)
        btns_layout = QHBoxLayout(btns_widget)

        self.btn_can_edit = QPushButton("禁用编辑")
        self.btn_can_edit.clicked.connect(self.toggle_readonly)
        btns_layout.addWidget(self.btn_can_edit)

        self.btn_clear = QPushButton("清空内容")
        self.btn_clear.clicked.connect(self.clear_content)
        btns_layout.addWidget(self.btn_clear)

        self.btn_get_content = QPushButton("获取内容")
        self.btn_get_content.clicked.connect(self.get_content)
        btns_layout.addWidget(self.btn_get_content)

    def toggle_readonly(self):
        self.edit_text_pure.setReadOnly(not self.edit_text_pure.isReadOnly())
        self.edit_text_pure2.setReadOnly(not self.edit_text_pure2.isReadOnly())
        self.btn_can_edit.setText("启用编辑" if self.edit_text_pure.isReadOnly() else "禁用编辑")

    def clear_content(self):
        self.edit_text_pure.clear()
        self.edit_text_pure2.clear()

    def get_content(self):
        print(self.edit_text_pure.toPlainText())
        print(self.edit_text_pure2.toPlainText())

app = QApplication(sys.argv)
window = MainWindow()
window.show()
sys.exit(app.exec())
```

与 QTextEdit 相比，QPlainTextEdit 不支持格式化，只能显示纯文本；不支持 append 方法，但支持 appendPlainText 方法。

运行效果如下：

![](/imgs/learn-qt/img-19.png)

使用 QPlainTextEdit 需要注意以下几点：

1. 只支持纯文本，HTML 不会渲染，HTML 会原样显示
2. 追加用 appendPlainText()
3. 读取全文 toPlainText()

## 小结

本文介绍了 PySide6 文本控件 QLabel、QLineEdit、QTextEdit、QPlainTextEdit 的使用，包括设置占位符、密码输入、追加文本、获取文本内容等，并提供相应示例代码进行说明