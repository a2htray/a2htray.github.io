+++
date = '2026-04-13T14:47:38+08:00'
draft = false
title = 'Python 中的类装饰器'
categories = ['后端技术', 'Python']
tags = ['Python', '类装饰器', '@dataclass', '@property', '@staticmethod', '@classmethod', '@abstractmethod']
toc = true
+++

本文用于记录在使用 Python 过程中不知道的类装饰器的用法与功能。

## @dataclass

```python
"""Add dunder methods based on the fields defined in the class.

Examines PEP 526 __annotations__ to determine fields.

If init is true, an __init__() method is added to the class. If repr
is true, a __repr__() method is added. If order is true, rich
comparison dunder methods are added. If unsafe_hash is true, a
__hash__() method is added. If frozen is true, fields may not be
assigned to after instance creation. If match_args is true, the
__match_args__ tuple is added. If kw_only is true, then by default
all fields are keyword-only. If slots is true, a new class with a
__slots__ attribute is returned.
"""
```

根据类中字段的定义，给类添加魔法方法：

* 如果 init 为 True，则会为类添加一个 `__init__()` 方法
* 如果 repr 为 True，则会添加一个 `__repr__()` 方法
* 如果 order 为 True，则会添加用于大小比较的魔法方法
* 如果 unsafe_hash 为 True，则会添加一个 `__hash__()` 方法
* 如果 frozen 为 True，则在实例创建后，不允许修改字段的值
* 如果 match_args 为 True，则会添加 `__match_args__` 元组
* 如果 kw_only 为 True，则默认所有字段都是仅限关键字参数
* 如果 slots 为 True，则会返回一个带有 `__slots__` 属性的新类

@dataclass 的声明如下：

![](/imgs/learn-python/01-dataclass-decalare.png)

在不传入参数的情况下，默认会增加 3 个魔法方法：

1. `__init__()`
2. `__repr__()`
3. `__eq__()`

示例如下：

```python
from dataclasses import dataclass

@dataclass
class Student:
    name: str
    age: int

student = Student(name="Joe", age=28)
print(student.__repr__()) # Student(name='Joe', age=28)

student2 = Student(name="Joe", age=28)
print("equal:", student == student2) # equal: True
```

## @property

@property 把类的方法当作属性调用，可以实现**只读**、**属性校验**、**动态计算属性**。

```python
class User:
    def __init__(self, age: int):
        self._age = age

    @property
    def age(self):
        return self._age

    @age.setter
    def age(self, value: int):
        if value < 0:
            raise ValueError("年龄不能小于0")
        if value > 120:
            raise ValueError("年龄不能大于120")
        self._age = value


user = User(18)
print(user.age)
user.age = 120
print(user.age)
```

## @staticmethod

@staticmethod 声明方法为静态方法，无需实例即可调用。

```python
class Math:
    """数学工具类"""
    @staticmethod
    def add(a: int, b: int) -> int:
        return a + b

print(Math.add(1, 2))
```

## @classmethod

@classmethod 把方法变成类方法，常用于工厂方法和多构造函数。

* 第一个参数为类本身，cls
* 直接`类名.方法名()` 调用

```python
class Product:
    """产品类"""

    def __init__(self, name: str, price: float):
        self.name = name
        self.price = price

    @classmethod
    def default(cls) -> "Product":
        return cls("default", 0.0)

    @classmethod
    def create(cls, name: str, price: float) -> "Product":
        return cls(name, price)


default_product = Product.default()
print(default_product.name)
print(default_product.price)

product = Product.create("test", 100.0)
print(product.name)
print(product.price)
```

## @abstractmethod

`多态设计`

@abstractmethod 用于定义抽象方法，必须配合抽象基类 ABC 使用，作用如下：

1. 子类必须实现父类定义的抽象方法
2. 定义接口约束

```python
from abc import ABC, abstractmethod


class Animal(ABC):
    """动物类"""
    @abstractmethod
    def sound(self) -> str:
        """动物的叫声"""
        pass


class Dog(Animal):
    """狗类"""
    def sound(self) -> str:
        return "汪汪汪"


dog = Dog()
print(dog.sound())
```


