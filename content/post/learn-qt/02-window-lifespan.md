+++
date = '2026-10-30T16:12:58+08:00'
draft = false
title = 'PySide6 窗口生命周期：show、close、hide、事件、窗口标志 WindowFlags'
categories = ['后端技术', 'PySide6']
tags = ['Python', 'PySide6', 'GUI']
toc = true
+++

![](/imgs/learn-qt/cover2.png)

## 前言

窗口有一套完整的**生命周期**：对象实例化 → 准备显示 → 显示 → 隐藏 → 关闭 → 对象销毁。同时配套 `show()` / `hide()` / `close()` 三个核心方法，以及窗口标志 `Qt.WindowFlags` 控制窗口样式（有无标题栏、是否置顶、弹窗样式等）。

> show() 显示、hide() 隐藏、close() 关闭。

本文会区分**逻辑隐藏**和**真正关闭销毁**，讲解关键事件，附可直接运行的代码。

> 重点区分（高频混淆）
> - `hide()`：隐藏窗口，**对象还在内存**，随时可以 `show()` 再次打开
> - `close()`：尝试关闭窗口，会触发关闭事件，**不一定销毁对象**
> - `show()`：显示窗口，不会新建对象

## 窗口生命周期

![](/imgs/learn-qt/img-07.png)

1. `__init__`：构造函数，只是创建窗口对象，**窗口不会显示**
2. `show()`：绘制并展示窗口，触发 `showEvent`
3. `hide()`：窗口从屏幕消失，内存对象保留，触发 `hideEvent`
4. `close()`：发送关闭请求，进入 `closeEvent`，这里可以拦截关闭动作
5. `Qt.WA_DeleteOnClose` 属性：**关键开关**
   - 默认：False，close 后只是 hide，对象还在，可以再次 show
   - 设置 True：close 被接受之后，销毁窗口对象

通过一个示例完整说明窗口的生命周期：以 Dialog 为例，主窗口左侧为文本区域，显示事件过程，右侧为按钮区域，演示事件操作。

```python
import sys
from datetime import datetime

from PySide6.QtWidgets import (
    QApplication, QMainWindow, QDialog, QPushButton, QWidget,
    QHBoxLayout, QVBoxLayout, QTextEdit, QDialogButtonBox, QLabel,
)


class DemoDialog(QDialog):
    """用于演示生命周期的对话框：把自身的事件转发到主窗口的日志框。"""

    def __init__(self, log_fn, parent=None):
        super().__init__(parent)
        self._log = log_fn
        self.setWindowTitle("Dialog")
        self.resize(300, 160)

        layout = QVBoxLayout(self)

        tip = QLabel("点击下方按钮，观察主窗口文本框中的事件日志")
        tip.setWordWrap(True)
        layout.addWidget(tip)

        box = QDialogButtonBox()
        btn_accept = box.addButton("接受(accept)", QDialogButtonBox.ButtonRole.AcceptRole)
        btn_reject = box.addButton("拒绝(reject)", QDialogButtonBox.ButtonRole.RejectRole)
        btn_hide = box.addButton("隐藏(hide)", QDialogButtonBox.ButtonRole.ActionRole)
        btn_close = box.addButton("关闭(close)", QDialogButtonBox.ButtonRole.DestructiveRole)
        layout.addWidget(box)

        # accept()/reject() 在非模态 show() 下等价于 hide()（仅隐藏，不销毁）
        btn_accept.clicked.connect(self.accept)
        btn_reject.clicked.connect(self.reject)
        btn_hide.clicked.connect(self.hide)
        btn_close.clicked.connect(self.close)

    # ---- 生命周期事件：重写以便记录到主窗口日志 ----
    def showEvent(self, event):
        self._log("Dialog.showEvent -> 对话框显示")
        super().showEvent(event)

    def hideEvent(self, event):
        self._log("Dialog.hideEvent -> 对话框隐藏")
        super().hideEvent(event)

    def closeEvent(self, event):
        self._log("Dialog.closeEvent -> 对话框关闭（默认仅隐藏，不销毁）")
        super().closeEvent(event)


class MainWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("主窗口")
        self.resize(400, 400)

        self.dialog = None  # 复用的对话框实例

        central = QWidget(self)
        self.setCentralWidget(central)
        layout = QVBoxLayout(central)

        # 顶部：两个按钮
        btn_row = QWidget()
        btn_layout = QHBoxLayout(btn_row)
        self.btn_new = QPushButton("show() 打开新 Dialog")
        self.btn_reuse = QPushButton("show() 重新打开隐藏的 Dialog")
        btn_layout.addWidget(self.btn_new)
        btn_layout.addWidget(self.btn_reuse)
        layout.addWidget(btn_row)

        # 底部：文本框（只读，用于展示 Dialog 的事件操作）
        self.log_edit = QTextEdit()
        self.log_edit.setReadOnly(True)
        layout.addWidget(self.log_edit, stretch=1)

        self.btn_new.clicked.connect(self.open_new_dialog)
        self.btn_reuse.clicked.connect(self.reopen_dialog)

        self.log("主窗口已启动")

    def log(self, msg):
        ts = datetime.now().strftime("%H:%M:%S")
        self.log_edit.append(f"[{ts}] {msg}")

    def open_new_dialog(self):
        # 每次新建一个 Dialog 实例，并用非模态的 show() 打开
        dlg = DemoDialog(self.log, self)
        self.dialog = dlg
        self.log("调用 dlg.show() 打开新的 Dialog（非模态）")
        dlg.show()

    def reopen_dialog(self):
        # 复用上一次创建的实例，演示对已隐藏/已关闭 Dialog 再次 show()
        if self.dialog is None:
            self.log("还没有创建过 Dialog，请先点击「show() 打开新 Dialog」")
            return
        if self.dialog.isVisible():
            self.log("Dialog 当前已可见，show() 不会重复触发 showEvent")
        self.log("调用 dlg.show() 重新打开之前隐藏/关闭的 Dialog")
        self.dialog.show()


app = QApplication(sys.argv)
window = MainWindow()
window.show()
app.exec()
```

![](/imgs/learn-qt/img-08.png)

## WindowFlags 窗口标志（Qt.WindowFlags）

`setWindowFlags(flags)` 用来控制窗口样式：标题栏、最大化最小化按钮、弹窗样式、置顶、无边框等。

> 结构说明：`WindowFlags` 本质是一个 `IntFlag`，由两部分按位或组成：
> - **低位（0x0 ~ 0xff）**：`WindowType` 窗口类型，决定窗口"属于哪一种"
> - **高位（0x100 及以上）**：`WindowHint` 窗口提示，逐位控制具体外观/行为 

**窗口类型 WindowType（低位，决定窗口类别）**

| 常量 | 值 | 含义 |
|---|---|---|
| `Qt.Widget` | `0x0` | 普通控件/子控件，不是顶层窗口（默认） |
| `Qt.Window` | `0x1` | 独立顶层窗口，带标题栏与边框，可移动 |
| `Qt.Dialog` | `0x3` | 对话框窗口（= Window \| Dialog），一般无最大化/最小化按钮 |
| `Qt.Sheet` | `0x5` | macOS 工作表（Sheet），从父窗口顶沿滑出的模态窗 |
| `Qt.Drawer` | `0x7` | macOS 抽屉窗口，从父窗口侧边滑出 |
| `Qt.Popup` | `0x9` | 弹出式窗口（无标题栏、立即夺焦点），如菜单/下拉 |
| `Qt.Tool` | `0xb` | 工具窗口，小标题栏，任务栏不显示图标 |
| `Qt.ToolTip` | `0xd` | 工具提示窗口，无标题栏、不接收焦点，一闪即过 |
| `Qt.SplashScreen` | `0xf` | 启动闪屏窗口（无边框，通常位于屏幕中央） |
| `Qt.Desktop` | `0x11` | 桌面窗口（覆盖整个桌面，很少直接用） |
| `Qt.SubWindow` | `0x12` | 子窗口（用于 MDI 多文档界面的内部子窗口） |
| `Qt.ForeignWindow` | `0x21` | 包装外部系统窗口的句柄（嵌入其他进程窗口） |
| `Qt.CoverWindow` | `0x41` | 覆盖窗口，被其他窗口遮挡时也保留绘制缓冲 |
| `Qt.WindowType_Mask` | `0xff` | 类型掩码，用于按位提取窗口类型部分 |

**窗口提示标志 WindowHint（高位，逐位控制外观/行为）**

| 常量 | 值 | 含义 |
|---|---|---|
| `Qt.MSWindowsFixedSizeDialogHint` | `0x100` | Windows：对话框固定大小，禁止拖动边框调整 |
| `Qt.MSWindowsOwnDC` | `0x200` | Windows：为控件分配独立的设备上下文（DC） |
| `Qt.BypassWindowManagerHint` | `0x400` | 绕过窗口管理器（WM）直接管理 |
| `Qt.X11BypassWindowManagerHint` | `0x400` | X11：绕过窗口管理器（与上值相同） |
| `Qt.FramelessWindowHint` | `0x800` | **无边框窗口**，没有标题栏、边框、按钮 |
| `Qt.WindowTitleHint` | `0x1000` | 显示标题栏 |
| `Qt.WindowSystemMenuHint` | `0x2000` | 显示系统菜单（标题栏左侧菜单） |
| `Qt.WindowMinimizeButtonHint` | `0x4000` | 显示最小化按钮 |
| `Qt.WindowMaximizeButtonHint` | `0x8000` | 显示最大化按钮 |
| `Qt.WindowMinMaxButtonsHint` | `0xc000` | 同时显示最小化 + 最大化按钮 |
| `Qt.WindowContextHelpButtonHint` | `0x10000` | 显示上下文帮助按钮（"?"） |
| `Qt.WindowShadeButtonHint` | `0x20000` | 显示卷起/隐藏（Shade）按钮 |
| `Qt.WindowStaysOnTopHint` | `0x40000` | **窗口置顶**，永远在最上层 |
| `Qt.WindowTransparentForInput` | `0x80000` | 窗口对鼠标输入透明（穿透，不接收鼠标事件） |
| `Qt.WindowOverridesSystemGestures` | `0x100000` | 覆盖系统手势（如触控板/多点触控） |
| `Qt.WindowDoesNotAcceptFocus` | `0x200000` | 窗口不接收键盘焦点 |
| `Qt.ExpandedClientAreaHint` | `0x400000` | Windows：客户区扩展到系统按钮周围 |
| `Qt.MaximizeUsingFullscreenGeometryHint` | `0x400000` | macOS：最大化时使用全屏几何（与上同值） |
| `Qt.NoTitleBarBackgroundHint` | `0x800000` | macOS：标题栏不绘制背景（透明） |
| `Qt.CustomizeWindowHint` | `0x2000000` | 自定义窗口，关闭标准标题栏装饰 |
| `Qt.WindowStaysOnBottomHint` | `0x4000000` | **窗口置底**，永远在最下层 |
| `Qt.WindowCloseButtonHint` | `0x8000000` | 显示关闭按钮 |
| `Qt.MacWindowToolBarButtonHint` | `0x10000000` | macOS：标题栏工具栏按钮（统一标题栏） |
| `Qt.BypassGraphicsProxyWidget` | `0x20000000` | 绕过图形代理控件（Graphics Proxy） |
| `Qt.NoDropShadowWindowHint` | `0x40000000` | 禁用窗口投影（阴影） |
| `Qt.WindowFullscreenButtonHint` | `0x80000000` | 显示全屏按钮 |

低位标识与高位标识搭配使用，如下代码：

```python
from PySide6.QtWidgets import QApplication, QMainWindow
from PySide6.QtCore import Qt


app = QApplication([])

window = QMainWindow()
# window.setWindowFlags(Qt.WindowType.Window)
# window.setWindowFlags(Qt.WindowType.Window | Qt.WindowType.FramelessWindowHint)
# window.setWindowFlags(Qt.WindowType.Dialog | Qt.WindowType.WindowStaysOnTopHint)
window.setWindowFlags(Qt.WindowType.Dialog | Qt.WindowType.MacWindowToolBarButtonHint)
window.show()

app.exec()
```