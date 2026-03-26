+++
date = '2026-03-25T10:54:47+08:00'
draft = false
title = '5 天学习 Playwright（Day3）：视频教程学习'
categories = ['后端技术', 'Playwright']
tags = ['自动化测试', '爬虫', 'Playwright', '视频教程']
toc = true
+++

文档看得差不多，在 B 站上找了个机构的培训视频，播放量挺多，有 4.4w，作为总结。

视频讲的内容相较于文档，内空是少了很多，但也补充了一些和 Selenium 的差异及其它。

文章的最后是 3 个视频。

## 和 Selenium 的核心差异

|              | Selenium                                                  | Playwright                                                           |
|:-------------|:----------------------------------------------------------|:---------------------------------------------------------------------|
| **通信原理**     | WebDriver 协议：基于HTTP，单向同步                                  | DevTools 协议，基于WS，双向异步                                                |
| **浏览器支持**    | 最广泛：通过第三方厂商提供的 Driver 支持几乎所有浏览器                           | 原生跨浏览器：项目原生支持 Chromium (Chrome, Edge), Firefox 和 WebKit (Safari 内核)。 |
| **执行速度与稳定性** | 相对较慢：多层协议转换带来延迟。稳定性依赖于开发者手动编写的显式等待                        | 极快且稳定：WS协议通信延迟极低。并且实现了自动等待机制，消除了绝大部分稳定性问题。                           |
| **核心功能与生态**  | 功能基础：专注于浏览器驱动，许多高级功能（如网络拦截、截图录制）需借助第三方库或复杂配置实现。生态成熟，社区庞大。 | 功能全面且强大：几乎支持Chromium全部功能，生态由微软公司维护                                   |

## Playwright 优势速览

基于与 Selenium 的对比，Playwright 的核心优势可总结为以下四点：

- **现代化的架构**：采用 DevTools Protocol，绕开了 WebDriver 的性能瓶颈，提供了更快、更可靠的浏览器控制能力。
- **与生俱来的稳定性**：强大的自动等待机制，是 Playwright 解决业界普遍存在的测试“脆弱性”问题的关键所在。
- **强大的原生工具链**：集成代码生成器（Codegen）、调试器（Inspector）和追踪查看器（Trace Viewer），极大提升了测试脚本的编写与调试效率。
- **全面的测试能力**：原生支持网络请求拦截、移动端设备模拟等复杂场景，无需依赖繁杂的第三方库，即可满足现代 Web 应用的全方位测试需求。

## 视频

<iframe src="//player.bilibili.com/player.html?bvid=BV1BBvAz2EB4&page=1&autoplay=0" scrolling="no" frameborder="0" style="width:100%; height:600px;"></iframe>

<iframe src="//player.bilibili.com/player.html?bvid=BV1BBvAz2EB4&page=2&autoplay=0" scrolling="no" frameborder="0" style="width:100%; height:600px;"></iframe>

<iframe src="//player.bilibili.com/player.html?bvid=BV1BBvAz2EB4&page=3&autoplay=0" scrolling="no" frameborder="0" style="width:100%; height:600px;"></iframe>
