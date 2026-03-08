+++
date = '2026-03-08T20:23:44+08:00'
draft = true
title = '5 天学习 Nuxt.js（Day 1）：介绍、安装、配置、视图'
categories = ['前端技术', 'Nuxt.js']
tags = ['Nuxt.js', 'Vue']
toc = true
+++

5 天学习 Nuxt.js，主要以官方文档为主，对于不明白的地方则问下豆包，加深理解。

资源：

1. [官方 - 介绍](https://nuxt.com/docs/4.x/getting-started/introduction)
2. [官方 - 安装](https://nuxt.com/docs/4.x/getting-started/installation)
3. [官方 - 配置](https://nuxt.com/docs/4.x/getting-started/configuration)
4. [官方 - 视图](https://nuxt.com/docs/4.x/getting-started/views)

## 介绍

_Nuxt is a free and open-source framework with an intuitive and extendable way to create type-safe, performant and
production-grade full-stack web applications and websites with Vue.js._

Nuxt 是基于 Vue.js 的全栈 Web 框架，其特点包括：

* 基于文件的路由：根据应用的 app/pages/ 目录结构来定义路由
* 代码分割：Nuxt 会自动将代码拆分为更小的代码块，有助于缩短应用的初始加载时间
* 服务端渲染：Nuxt 内置了服务端渲染（SSR）能力
* 自动导入：在对应的目录中编写 Vue 组合式函数和组件后，可直接使用无需手动导入
* 数据获取工具：Nuxt 提供了组合式函数，可处理兼容 SSR 的数据获取逻辑
* 零配置 TypeScript 支持：自动生成的类型文件和 tsconfig.json 配置
* 预置的构建工具：Vite

作为一款支持 SSR 能力的框架，SSR 的优势则必须一提，包括：

* 更快的初始页面加载速度
* 搜索引擎优化（SEO）
* 更适用于低性能设备
* 更好的可访问性
* 更易实现缓存

> 优势的本质在于：服务端把完整的 HTML 代码返回给浏览器，浏览器无须执行复杂的脚本逻辑，同时可以缓存整个页面。

## 安装

以 `pnpm` 的方式安装 Nuxt，所以先安装 `pnpm`：

```bash
$ sudo corepack enable
$ corepack prepare pnpm@latest --activate
$ pnpm -v
10.30.3
```

创建项目

```bash
$ pnpm create nuxt@latest llm_agent
```

交互式创建过程中，选择的选项分别是：

1. minimal – Minimal setup for Nuxt 4 (recommended)
2. pnpm
3. Initialize git repository? - Yes
4. browse and install modules? - No

运行程序：

```bash
$ cd llm_agent
$ pnpm dev -o
```

## 配置

_The nuxt.config.ts file can override or extend this default configuration._

nuxt.config.ts 可用于：

1. 添加自定义脚本
2. 注册模块
3. 修改渲染模式

配置项大全。见 https://nuxt.com/docs/4.x/api/nuxt-config。

### 环境配置重写

在 nuxt.config.ts 文件定义不同环境的配置，如：

```typescript
export default defineNuxtConfig({
  $production: {
    routeRules: {
      '/**': { isr: true },
    },
  },
  $development: {
    //
  },
  $env: {
    staging: {
      //
    },
  },
})
```

在构建应用时，使用 --envName 选项指定环境，如：

```bash
$ nuxt build --envName staging
```

### 环境变量

### 应用配置

### runtimeConfig vs. app.config

## 视图

### 组件

### 页面

### 布局

## 总结




