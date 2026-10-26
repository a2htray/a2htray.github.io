+++
date = '2026-10-26T16:38:33+08:00'
draft = false
title = 'React useMemo Hook - 缓存函数调用结果'
categories = ['前端技术', 'React']
tags = ['React', 'React Hooks', 'useMemo']
toc = true
+++

## 什么是 useMemo

`useMemo` 是 React 提供的一个 Hook，用来将计算结果缓存起来，只有当依赖项发生变化时才会才会重新计算。

> useMemo 不保证一定缓存命中，在内存压力、并发渲染等情况下，React 可能丢弃缓存并在下次渲染时重新计算。

## 使用场景

**场景 1：大量的同步计算**

比如大列表过滤、排序、复杂数学运算、大数据遍历，渲染频繁触发时。

```javascript
const list = [/* 上万条数据 */]

const filteredList = list.filter(item => item.name.includes(keyword))
```

组件每次渲染，函数 / 表达式就从头到尾执行一遍。


```javascript
const list = [/* 上万条数据 */]

const filteredList = useMemo(() => {
  return list.filter(item => item.name.includes(keyword))
}, [list, keyword])
```

**场景 2：配合 React.memo 减少子组件渲染**

`useMemo`：依赖不变时，**返回上次的同一个引用**。

`React.memo(Child)`：子组件会对 props 做**浅层比较**，props 没变就跳过渲染。

```javascript
import { useState, useMemo } from 'react'

// 子组件，用 React.memo 包裹
const Child = React.memo(({ user }) => {
  return <div>名字：{user.name}</div>
})

export default function Parent() {
  const [name, setName] = useState('小明')

  // useMemo 缓存对象引用
  const user = useMemo(() => ({ name }), [name])

  console.log('父组件渲染')
  return (
    <div>
      <p>count: {count}</p>
      <button onClick={() => setName('小红')}>改名字</button>
      <Child user={user} />
    </div>
  )
}
```

当 `name` 变化时，`user` 对象引用也会变化，子组件会重新渲染。
