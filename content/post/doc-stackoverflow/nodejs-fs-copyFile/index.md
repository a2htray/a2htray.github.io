---
title: "静态资源的管理-文件复制"
date: 2022-05-12T11:09:53+08:00
draft: false
toc: true
reward: false
categories:
 - 后端技术
 - Nodejs
tags:
 - stackoverflow
 - Node.js
---

这是一个简单的 Node.js 的 API 的使用示例：如何基于 Node.js 的完成文件的复制。

<!--more-->

**问题起源**

该问题源起于一个前后端不分离项目的前端静态资源如 JS、CSS 文件的管理。考虑到要对第三方静态资源的修改非常少且需要对其进行版本管理，所以想到使用 `npm` 对其进行管理。

**步骤**

实施的步骤如下：

1. 使用 `npm install PACKAGE_NAME` 下载到本地；
2. 在 `package.json` 中配置发布脚本；
3. 发布脚本的作用是将 `PACKAGE_NAME` 中的静态资源复制到指定目录；

**stack overflow**

stack overflow 的问题：在 Node.js 中，如何快速地复制文件。高赞答案提供了两种实现方式：

1. 使用 `fs.copyFile`；

```js
const fs = require('fs');

// File destination.txt will be created or overwritten by default.
fs.copyFile('source.txt', 'destination.txt', (err) => {
  if (err) throw err;
  console.log('source.txt was copied to destination.txt');
});
```

2. 流的读取与写入；

```js
const fs = require('fs');
fs.createReadStream('test.log').pipe(fs.createWriteStream('newLog.log'));
```

**在项目中使用**

在项目中，则是预先定义一系列的资源文件的源地址与目标地址，通过 `fs.copyFile` 完成复制操作。

```js
const fs = require('fs')

const universalFromPrefix = (path) => './node_modules/' + path
const universalToPrefix = (path) => './static/' + path

const resources = {
    bootstrap: [
        {
            from: universalFromPrefix('bootstrap/dist/js/bootstrap.min.js'),
            to: universalToPrefix('js/bootstrap.min.js'),
        },
        {
            from: universalFromPrefix('bootstrap/dist/js/bootstrap.min.js.map'),
            to: universalToPrefix('js/bootstrap.min.js.map'),
        },
        {
            from: universalFromPrefix('bootstrap/dist/css/bootstrap.min.css'),
            to: universalToPrefix('css/bootstrap.min.css'),
        },
        {
            from: universalFromPrefix('bootstrap/dist/css/bootstrap.min.css.map'),
            to: universalToPrefix('css/bootstrap.min.css.map'),
        },
    ]
}

for (let resKey in resources) {
    let resource = resources[resKey]
    console.log(`copy resource ${resKey}`)
    for (let i = 0; i < resource.length; i++) {
        let {from, to} = resource[i]
        fs.copyFile(from, to, (err) => {
            if (err) throw err
            console.log(`=> copy ${from} to ${to}`)
        })
    }
}
```

**好处**

包的下载依赖于 `npm`，无须手动地下载特定的静态资源。