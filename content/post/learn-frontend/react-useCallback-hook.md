+++
date = '2026-10-27T10:02:45+08:00'
draft = false
title = 'React useCallback Hook - 缓存函数引用'
categories = ['前端技术', 'React']
tags = ['React', 'React Hooks', 'useCallback']
toc = true
+++

## 什么是 useCallback

区别于 `useMemo`，`useCallback` **缓存函数引用，依赖不变时，多次渲染返回同一个函数对象**，不是缓存函数执行结果。

格式如下：

```javascript
const memoFn = useCallback(
  (...args) => {
    // 你的逻辑
  },
  [deps]
)
```

* 参数 1：要缓存的函数
* 参数 2：依赖数组，依赖不变时，函数引用不变
* 返回值：缓存后的函数引用

与 `useMemo` 的等价写法：

```javascript
useCallback(fn, deps)

useMemo(() => fn, deps)
```

## 使用场景

**场景 1：搭配 React.memo - 减少子组件渲染**

```javascript
const Child = React.memo(({ onClick }) => <button onClick={onClick} />)

// 父组件
const handleClick = useCallback(() => {
  console.log('click')
}, [])
<Child onClick={handleClick} />
```

用 `useCallback` 固定函数引用，依赖不变时，传给子组件的回调不变，React.memo 生效。

**场景 2：函数作为 `useEffect` 的依赖项**

```javascript
const loadList = useCallback(async () => {
  const res = await fetch(`/api/list?id=${id}`)
}, [id])

useEffect(() => {
  loadList()
}, [loadList])
```

用 `useCallback` 固定函数引用，依赖不变时，`useEffect` 依赖项不变，避免重复执行。适用于：

- 接口请求函数
- 定时器回调
- 事件订阅解绑函数

