+++
date = '2026-11-05T09:47:33+08:00'
title = 'PySide6 扩展文本控件（一）：QTextBrowser、QLCDNumber'
categories = ['后端技术', 'PySide6']
tags = ['Python', 'PySide6', 'GUI', '文本控件', '文本扩展控件']
toc = true
+++

![](/imgs/learn-qt/cover6.png)

## 扩展控件

基础的四个控件（`QLabel / QLineEdit / QTextEdit / QPlainTextEdit`）能覆盖绝大多数场景，但 PySide6 还提供了一批"特定用途"的文本相关控件。这些扩展控件专用于某一特定场景的需求：

- `QTextBrowser`：只读富文本展示 + 超链接导航
- `QLCDNumber`：数码管风格数字展示，来展示数字、少量符号
- `QKeySequenceEdit`：让用户录制快捷键
- `QFontComboBox`：选择系统字体
- `QGraphicsTextItem`：在画布（GraphicsView）上放文字

## QTextBrowser

`QTextBrowser` 继承自 `QTextEdit`，默认只读，专门用来做富文本查看器、帮助文档、日志说明页、简易浏览器。

下面代码演示`帮助文档`与`简易浏览器`的实现：

```python
import sys

from PySide6.QtWidgets import (
    QApplication, QMainWindow, QWidget, QHBoxLayout, QVBoxLayout,
    QTextBrowser, QLineEdit, QPushButton, QLabel,
)
from PySide6.QtCore import QUrl


class MainWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QTextBrowser 演示：左帮助文档 / 右简易浏览器")
        self.resize(900, 520)

        central = QWidget()
        self.setCentralWidget(central)
        root = QHBoxLayout(central)

        # ============ 左列：帮助文档 ============
        self.help = QTextBrowser()
        self.help.setOpenLinks(True)
        self.help.setHtml(self.help_html())
        self.help.anchorClicked.connect(
            lambda url: self.status(f"帮助文档：点击了链接 {url.toString()}")
        )
        root.addWidget(self.help)

        # ============ 右列：简易浏览器 ============
        right = QWidget()
        rlayout = QVBoxLayout(right)

        nav = QHBoxLayout()
        self.url_input = QLineEdit("doc://home")
        btn_go = QPushButton("Go")
        btn_back = QPushButton("Back")
        nav.addWidget(QLabel("地址:"))
        nav.addWidget(self.url_input)
        nav.addWidget(btn_go)
        nav.addWidget(btn_back)
        rlayout.addLayout(nav)

        self.browser = QTextBrowser()
        rlayout.addWidget(self.browser)

        self.history = []
        self.pages = self.site_pages()

        self.browser.anchorClicked.connect(self.on_link)
        btn_go.clicked.connect(self.go)
        btn_back.clicked.connect(self.back)
        root.addWidget(right)

        self.status_label = QLabel()
        rlayout.addWidget(self.status_label)

        self.go()  # 载入首页
        self.status("就绪")

    # ---------- 左：帮助文档内容 ----------
    def help_html(self):
        return """
        <h2>使用帮助</h2>
        <p>这是用 <b>QTextBrowser</b> 渲染的帮助文档，支持 HTML 与页内跳转。</p>
        <ul>
          <li><a href="#install">安装说明</a></li>
          <li><a href="#usage">基本用法</a></li>
          <li><a href="#faq">常见问题</a></li>
        </ul>
        <h3 id="install">安装说明</h3>
        <p>使用 pip 安装依赖即可：<code>pip install PySide6</code>。</p>
        <h3 id="usage">基本用法</h3>
        <p>把 QTextBrowser 当作只读、支持超链接的富文本视图使用。</p>
        <h3 id="faq">常见问题</h3>
        <p>QTextBrowser 与 QTextEdit 的区别：前者只读且内置链接/历史导航。</p>
        """

    # ---------- 右：简易站点（内存中的几页 HTML） ----------
    def site_pages(self):
        home = """
        <h2>首页 (doc://home)</h2>
        <p>这是一个用 QTextBrowser 实现的极简“浏览器”演示。</p>
        <ul>
          <li><a href="doc://about">关于页面</a></li>
          <li><a href="doc://contact">联系页面</a></li>
        </ul>
        """
        about = """
        <h2>关于 (doc://about)</h2>
        <p>本页面演示 QTextBrowser 渲染 HTML 并处理 <a href="doc://home">返回首页</a> 链接。</p>
        """
        contact = """
        <h2>联系 (doc://contact)</h2>
        <p>邮箱：<a href="doc://about">回到关于</a></p>
        """
        return {"doc://home": home, "doc://about": about, "doc://contact": contact}

    def go(self):
        url = self.url_input.text().strip()
        self.load(url)

    def on_link(self, url: QUrl):
        self.url_input.setText(url.toString())
        self.load(url.toString())

    def load(self, url):
        if url in self.pages:
            self.history.append(url)
            self.browser.setHtml(self.pages[url])
            self.status(f"已加载：{url}")
        else:
            self.browser.setHtml(f"<h2>404</h2><p>未找到页面：{url}</p>")
            self.status(f"未找到：{url}")

    def back(self):
        if len(self.history) > 1:
            self.history.pop()
            prev = self.history[-1]
            self.url_input.setText(prev)
            self.browser.setHtml(self.pages[prev])
            self.status(f"后退到：{prev}")
        else:
            self.status("已在首页，无法后退")

    def status(self, msg):
        self.status_label.setText(msg)


if __name__ == "__main__":
    app = QApplication(sys.argv)
    window = MainWindow()
    window.show()
    app.exec()
```

采用水平布局，左侧是帮助文档，右侧是简易浏览器：

- `setOpenLinks` 方法设置是否打开超链接，默认不打开
- `anchorClicked` 信号用于处理超链接点击事件，传入 `QUrl` 对象
- `setHtml` 方法设置 HTML 内容，支持页内跳转
- `self.history` 维护浏览历史

运行效果如下：

![](/imgs/learn-qt/img-20.png)

QTextBrowser 内置历史栈：`back()` / `forward()` / `home()`；支持 `setSource()`、`searchPaths` 本地文档路径查找，示例如下：

```python
import sys

from PySide6.QtWidgets import (
    QApplication, QMainWindow, QWidget, QVBoxLayout, QHBoxLayout,
    QTextBrowser, QLineEdit, QPushButton, QLabel,
)
from PySide6.QtCore import QUrl
from PySide6.QtGui import QTextDocument


class DemoBrowser(QTextBrowser):
    """重写 loadResource，让 QTextBrowser 能从内存“站点”加载 HTML，
    从而真正使用它自带的 source + 历史栈机制（setSource / backward / forward）。"""

    def __init__(self, pages):
        super().__init__()
        self.pages = pages
        self.setOpenLinks(True)

    def loadResource(self, resource_type, url):
        if resource_type == QTextDocument.HtmlResource and url.toString() in self.pages:
            return self.pages[url.toString()]
        return super().loadResource(resource_type, url)


class MainWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QTextBrowser 自带历史栈演示")
        self.resize(720, 560)

        central = QWidget()
        self.setCentralWidget(central)
        root = QVBoxLayout(central)

        # ---- 工具栏：地址栏 + Go / Back / Forward ----
        nav = QHBoxLayout()
        self.url_input = QLineEdit("doc://home")
        btn_go = QPushButton("Go")
        self.btn_back = QPushButton("Back")
        self.btn_forward = QPushButton("Forward")
        nav.addWidget(QLabel("地址:"))
        nav.addWidget(self.url_input)
        nav.addWidget(btn_go)
        nav.addWidget(self.btn_back)
        nav.addWidget(self.btn_forward)
        root.addLayout(nav)

        # ---- 浏览器 ----
        self.browser = DemoBrowser(self.site_pages())
        root.addWidget(self.browser)

        # ---- 状态栏：展示历史栈信息 ----
        self.status_label = QLabel()
        root.addWidget(self.status_label)

        # 信号连接：导航后刷新状态；点击链接时同步地址栏
        btn_go.clicked.connect(self.go)
        self.btn_back.clicked.connect(self.browser.backward)
        self.btn_forward.clicked.connect(self.browser.forward)
        self.browser.sourceChanged.connect(self.update_status)
        self.browser.anchorClicked.connect(
            lambda url: self.url_input.setText(url.toString())
        )

        self.go()
        self.update_status()

    def site_pages(self):
        home = """
        <h2>首页 (doc://home)</h2>
        <p>这是 QTextBrowser 自带历史栈演示。点击链接或用 Back/Forward 导航。</p>
        <ul>
          <li><a href="doc://about">关于页面</a></li>
          <li><a href="doc://contact">联系页面</a></li>
        </ul>
        """
        about = """
        <h2>关于 (doc://about)</h2>
        <p>QTextBrowser 的 backward()/forward() 基于它内部的 source 历史栈。</p>
        <p><a href="doc://contact">去联系页面</a> | <a href="doc://home">回首页</a></p>
        """
        contact = """
        <h2>联系 (doc://contact)</h2>
        <p>状态栏会显示 backwardHistoryCount / forwardHistoryCount。</p>
        <p><a href="doc://about">回关于页面</a> | <a href="doc://home">回首页</a></p>
        """
        return {"doc://home": home, "doc://about": about, "doc://contact": contact}

    # ---------- 导航 ----------
    def go(self):
        # setSource 会把当前 URL 压入 QTextBrowser 自带的 source 历史栈
        self.browser.setSource(QUrl(self.url_input.text().strip()))

    def update_status(self):
        src = self.browser.source().toString()
        back_n = self.browser.backwardHistoryCount()
        fwd_n = self.browser.forwardHistoryCount()
        self.btn_back.setEnabled(self.browser.isBackwardAvailable())
        self.btn_forward.setEnabled(self.browser.isForwardAvailable())
        self.status_label.setText(
            f"当前: {src}  |  历史栈 后退={back_n}  前进={fwd_n}"
        )


app = QApplication(sys.argv)
window = MainWindow()
window.show()
app.exec()
```

- `QTextBrowser` 的 `sourceChanged` 信号触发更新内容
- `QTextBrowser` 的 `anchorClicked` 信号触发点击链接

运行效果如下：

![](/imgs/learn-qt/img-21.png)

## QLCDNumber

QLCDNumber 用于展示数字、少量符号，如计数器、时间显示等，其核心能力包括：

- 支持十进制 / 十六进制 / 八进制 / 二进制
- 支持 `int` / `double` / `str` 三种方式赋值
- 三种数码段风格：Flat / Filled / Outline
- 支持小时钟样式 `hh:mm:ss`

```python
import sys

from PySide6.QtWidgets import (
    QApplication, QMainWindow, QWidget, QVBoxLayout, QHBoxLayout,
    QLCDNumber, QComboBox, QLineEdit, QPushButton, QLabel,
)
from PySide6.QtCore import QTimer, QTime


class MainWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QLCDNumber 演示")
        self.resize(480, 420)

        central = QWidget()
        self.setCentralWidget(central)
        root = QVBoxLayout(central)

        # ===== 显示屏：8 位，足够显示 hh:mm:ss =====
        self.lcd = QLCDNumber()
        self.lcd.setDigitCount(8)
        self.lcd.setSegmentStyle(QLCDNumber.SegmentStyle.Outline)
        self.lcd.display(0)
        root.addWidget(self.lcd)

        # ===== 进制（Dec / Hex / Oct / Bin） =====
        mode_row = QHBoxLayout()
        self.mode_combo = QComboBox()
        self.mode_combo.addItems(["十进制 Dec", "十六进制 Hex", "八进制 Oct", "二进制 Bin"])
        self.mode_combo.setCurrentText("十进制 Dec")
        mode_row.addWidget(QLabel("进制:"))
        mode_row.addWidget(self.mode_combo)
        root.addLayout(mode_row)

        # ===== 数码段风格（Flat / Filled / Outline） =====
        style_row = QHBoxLayout()
        self.style_combo = QComboBox()
        self.style_combo.addItems(["Flat", "Filled", "Outline"])
        self.style_combo.setCurrentText("Filled")
        style_row.addWidget(QLabel("段风格:"))
        style_row.addWidget(self.style_combo)
        root.addLayout(style_row)

        # ===== 三种赋值方式：int / double / str =====
        val_row = QHBoxLayout()
        self.type_combo = QComboBox()
        self.type_combo.addItems(["int", "double", "str"])
        self.val_input = QLineEdit("123")
        btn_display = QPushButton("显示")
        val_row.addWidget(QLabel("类型:"))
        val_row.addWidget(self.type_combo)
        val_row.addWidget(QLabel("值:"))
        val_row.addWidget(self.val_input)
        val_row.addWidget(btn_display)
        root.addLayout(val_row)

        # ===== 小时钟样式 hh:mm:ss =====
        self.btn_clock = QPushButton("启动时钟 hh:mm:ss")
        root.addWidget(self.btn_clock)

        # ---- 信号连接 ----
        self.mode_combo.currentTextChanged.connect(self.apply_mode)
        self.style_combo.currentTextChanged.connect(self.apply_style)
        btn_display.clicked.connect(self.display_value)
        self.btn_clock.clicked.connect(self.toggle_clock)

        # ---- 时钟用的定时器 ----
        self.timer = QTimer()
        self.timer.setInterval(1000)
        self.timer.timeout.connect(self.show_clock)

        self.apply_mode(self.mode_combo.currentText())
        self.apply_style(self.style_combo.currentText())
        self.display_value()

    # ---------- 进制 ----------
    def apply_mode(self, text):
        mapping = {
            "十进制 Dec": QLCDNumber.Mode.Dec,
            "十六进制 Hex": QLCDNumber.Mode.Hex,
            "八进制 Oct": QLCDNumber.Mode.Oct,
            "二进制 Bin": QLCDNumber.Mode.Bin,
        }
        self.lcd.setMode(mapping[text])

    # ---------- 数码段风格 ----------
    def apply_style(self, text):
        mapping = {
            "Flat": QLCDNumber.SegmentStyle.Flat,
            "Filled": QLCDNumber.SegmentStyle.Filled,
            "Outline": QLCDNumber.SegmentStyle.Outline,
        }
        self.lcd.setSegmentStyle(mapping[text])

    # ---------- 三种赋值方式 ----------
    def display_value(self):
        t = self.type_combo.currentText()
        text = self.val_input.text().strip()
        if t == "int":
            self.lcd.display(int(text))        # 以当前进制显示整数
        elif t == "double":
            self.lcd.display(float(text))      # 显示小数（含小数点）
        else:  # str
            self.lcd.display(text)             # 直接显示字符串字符

    # ---------- 小时钟 ----------
    def toggle_clock(self):
        if self.timer.isActive():
            self.timer.stop()
            self.btn_clock.setText("启动时钟 hh:mm:ss")
        else:
            self.show_clock()
            self.timer.start()
            self.btn_clock.setText("停止时钟")

    def show_clock(self):
        # 用字符串方式赋值，按 hh:mm:ss 显示当前时间
        self.lcd.display(QTime.currentTime().toString("hh:mm:ss"))


app = QApplication(sys.argv)
window = MainWindow()
window.show()
app.exec()
```

运行效果如下：

![](/imgs/learn-qt/img-22.png)

## 小结

本文主要介绍了文本扩展控件 `QTextBrowser` 和 `QLCDNumber`，并给出代码示例。

下一篇：QKeySequenceEdit、QFontComboBox、QGraphicsTextItem。

持续学习 ing -^_^-
