+++
date = '2026-03-09T17:27:25+08:00'
draft = false
title = '工具方法：获取响应的 cookie'
categories = ['前端技术', 'Nuxt.js']
tags = ['Nuxt.js', 'Vue', 'TypeScript']
toc = true
+++

起因：应用在服务端要请求第三方的接口，获取数据的同时，还保存响应头中的 Cookie 信息。

本来是想命名为 `getCookie`，但该方法名与 Nuxt 内置的一个方法重名，两种方式可以解决：

1. 另取别名，如 `getCustomCookie` 之类，但我觉得不雅，所以没有采用
2. 命名空间隔离，**采用这种**

> 内置的 `getCookie` 方法用于在服务端获取请求中携带的 Cookie 信息。

工具方法定义在 `/server/utils/cookie.ts` 中，如下：

```typescript
export const responseCookie = {
  get(response: Response, cookieName: string): string | undefined {
    const cookie = response.headers.get('set-cookie')
    if (!cookie) {
        return undefined
    }

    const cookieValue = cookie.split('; ').find((item) => item.startsWith(`${cookieName}=`))
    if (!cookieValue) {
        return undefined
    }

    return cookieValue.split('=')[1]
  }
}
```

## 总结

1. 服务端的工具方法放置至 `/server/utils` 目录
2. api 代码中无须显示引入，直接使用 `responseCookie.get`