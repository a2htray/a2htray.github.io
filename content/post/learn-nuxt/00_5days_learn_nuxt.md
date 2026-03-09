+++
date = '2026-03-08T20:23:44+08:00'
draft = false
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

字段 runtimeConfig 定义了应用程序可访问的值（类似于环境变量），默认情况下，这些值在服务端侧可读，而字段 runtimeConfig.public 和 
runtimeConfig.app 定义的值则也可以在客户端侧。

```typescript
export default defineNuxtConfig({
  runtimeConfig: {
    // 只能在服务端侧访问
    apiSecret: '123',
    // 客户端测同样可以访问
    public: {
      apiBase: '/api',
    },
  },
})
```

> Composable 组合式函数：抽离出来的、可复用的逻辑函数。

这些变量会通过 `useRuntimeConfig()` 组合式函数暴露给应用的其他部分使用。

```html
<script setup lang="ts">
  const runtimeConfig = useRuntimeConfig()
</script>
```

### 应用配置

app.config.ts 文件位于源码目录（默认是 app/ 目录），用于暴露构建阶段即可确定的公共变量。与 runtimeConfig 配置项不同，这些变量无法通过环境变量覆盖。

这些变量会通过 `useAppConfig()` 组合式函数暴露给应用的其他部分使用。

```html
<script setup lang="ts">
  const appConfig = useAppConfig()
</script>
```

### runtimeConfig vs. app.config

runtimeConfig 是运行时生效的配置（可通过环境变量动态修改，支持私有 / 公有），app.config 是**构建时固化的配置**（打包后无法修改，仅存公有信息）。

关键差异：

- 生效时机：
  - runtimeConfig：配置在应用运行阶段生效，即使项目打包部署后，你依然可以通过修改环境变量（比如 .env 文件、服务器环境变量）来更新配置，无需重新打包。
  - app.config：配置在项目构建（build）阶段就被编译进代码里，打包完成后就 “固定死了”，想要修改必须重新打包部署。
- 数据安全性
  - runtimeConfig：支持私有（private） 和公有（public） 配置：
    - 私有配置（如 API 密钥、数据库密码）：仅在服务端可用，不会暴露到客户端代码中。
    - 公有配置（如公开的 API 基础地址）：可通过 useRuntimeConfig().public 暴露到客户端。
  - app.config：所有配置都是公开的，会被编译到客户端代码中（能在浏览器 DevTools 里看到），绝对不能放敏感信息（比如密钥、token）。
- 适用场景
  - runtimeConfig：存放需要动态调整、包含敏感信息的配置，比如：
    - 私有：API 密钥、数据库连接字符串、第三方服务的 secret；
    - 公有：可动态修改的 API 基础地址、环境标识（dev/prod）。
  - app.config：存放固定不变、非敏感的项目配置，比如：
    - 网站标题、logo 地址、主题色；
    - 静态的功能开关（如是否显示某个模块）、语言列表；
    - 项目版本号、作者信息。
- 访问方式
  - runtimeConfig：必须通过 `useRuntimeConfig()` 组合式函数访问，且服务端 / 客户端访问规则不同
  - app.config：通过 `useAppConfig()` 组合式函数访问，服务端 / 客户端均可访问全部内容

## 视图

### 组件

在 Nuxt 中，在 `app/components/` 目录下创建这些组件，它们会自动在整个应用中可用，**无需手动显式导入**。

在 `app/components` 目录下创建文件 `Topbar.vue`，代码如下：

```vue
<template>
  <div class="topbar">
    <div class="logo">
      <router-link to="/">
        <h1><slot name="logo" /></h1>
      </router-link>
    </div>
    <div class="nav-links">
      <router-link to="/tasks">任务列表</router-link>
      <router-link to="/settings">系统配置</router-link>
    </div>
  </div>
</template>

<style scoped>
.topbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
  background-color: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
}

.logo h1 {
  margin: 0;
  font-size: 1.5rem;
  font-weight: bold;
}

.logo a {
  text-decoration: none;
  color: #333;
}

.nav-links {
  display: flex;
  gap: 2rem;
}

.nav-links a {
  text-decoration: none;
  color: #333;
  font-weight: 500;
  transition: color 0.3s ease;
}

.nav-links a:hover {
  color: #007bff;
}
</style>
```

在 `app/app.vue` 使用 `Topbar` 组件：

```vue
<template>
  <div>
    <Topbar>
      <template #logo>
        LLM Agent
      </template>
    </Topbar>
    Index
  </div>
</template>
```

### 页面

多个页面跳转，采用 `app/pages/index.vue`：

1. 删除 `app/page.vue` 文件
2. 创建 `app/pages/index.vue` 文件

```bash
$ rm app/page.vue
$ rm app/pages/index.vue
```

**index.vue**

```vue
<template>
  <div>
    <Topbar>
      <template #logo>
        LLM Agent
      </template>
    </Topbar>
    Index
  </div>
</template>
```

在 `app/pages` 目录下创建页面，如：

```bash
$ mkdir app/pages
$ touch app/pages/tasks.vue
$ touch app/pages/settings.vue
```

**tasks.vue**

```vue
<template>
  <div>
    <Topbar>
      <template #logo>
        LLM Agent
      </template>
    </Topbar>
    Tasks
  </div>
</template>
```

**settings.vue**

```vue
<template>
  <div>
    <Topbar>
      <template #logo>
        LLM Agent
      </template>
    </Topbar>
    Settings
  </div>
</template>
```

浏览器效果

![](/imgs/learn-nuxt/1773021309698.gif)

### 布局

布局是使用 `<slot />` 插槽组件来展示页面内容的 Vue 文件。`app/layouts/default.vue` 文件会被默认使用，自定义布局可通过
**页面元数据**进行配置。

创建 layout 文件

```bash
$ mkdir app/layouts
$ touch app/layouts/default.vue
```

**layouts/default.vue**

```vue
<template>
  <div>
    <Topbar>
      <template #logo>
        LLM Agent
      </template>
    </Topbar>
    <slot />
  </div>
</template>
```

调整 `index.vue` 、`tasks.vue` 、`settings.vue`：

**index.vue**

```vue
<template>
  <div>
    Index
  </div>
</template>
```

**tasks.vue**

```vue
<template>
  <div>
    Tasks
  </div>
</template>
```

**settings.vue**

```vue
<template>
  <div>
    Settings
  </div>
</template>
```

浏览器效果不变。

## 总结

了解了 Nuxt 的特点和 SSR 的优势，同时掌握了 Nuxt 项目的搭建过程。

学习了环境变量（`nuxt.config.ts`）、应用配置（`app.config.ts`）的定义、使用，以及它们的区别，如：

1. 定义的位置不同
2. 是否可以被覆盖
3. 使用场景
4. 数据安全性
5. 生效时机

在编码实践中，遵循 Nuxt 约定的规范，如：

1. 组件放置在目录 `/app/components`，组件无须显示导入
2. 页面放置在目录 `/app/pages`
3. 布局放置在目录 `/app/layouts`，默认使用 `/app/layouts/default.vue`，也可通过在页面中设置元数据修改引用布局


