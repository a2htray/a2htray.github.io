+++
date = '2026-05-13T14:32:50+08:00'
draft = false
title = 'ref、reactive 响应式 API'
categories = ['前端技术', 'Vue']
tags = ['Vue', '响应式 API', 'ref', 'reactive']
toc = true
+++

`本文在 LLM 对话中学习、总结、记录`

## 响应式 API

响应式 API 是 Vue 官方提供的一组函数，专门用来把**普通变量 / 对象变成自动监听变化、自动更新视图**的接口。Vue3 响应式 API共分成 4 类：

* 创建响应式：ref、reactive
* 转换 / 解构响应式：toRef、toRefs
* 只读响应式：readonly、shallowReadonly
* 浅层响应式：shallowRef、shallowReactive

本文主要介绍创建响应式 API - **ref**、**reactive**。

## ref 和 reactive

 ref，全称 reference，**能为任意类型创建响应式变量**，它的原理是：把普通值包一层对象，内部有 value 属性，实现响应式。

 > ref 用于基本数据类型。

reactive，属于代理式响应式，**只能代理对象 / 数组（引用类型）**，它的原理是：给原始对象创建 Proxy 代理。

> reactive 用于对象 / 数组等引用类型。

## 差异

ref 和 reactive 的差异如下表：

|特性|ref|reactive|
|--|--|--|
|支持的数据类型|所有类型（基本类型 + 引用类型）|仅引用类型（对象、数组、Map、Set 等）|
|访问 / 修改方式|JS 中必须用 .value	|直接访问，无需 .value|
|响应式丢失|不会丢失|解构 / 展开会丢失响应式|
|模板中使用|自动解包，不用写 .value	|直接使用|
|重新赋值|可以直接重新赋值|不建议直接重新赋值（会破坏响应式）|

## 开发实践

* 优先用 ref：通用、安全、不会踩坑，能处理所有数据
* reactive：只适合完整对象 / 表单数据这种整体结构
* 基本类型（数字、字符串、布尔）：只能用 ref

