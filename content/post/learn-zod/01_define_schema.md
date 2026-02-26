+++
date = '2026-02-25T15:09:54+08:00'
draft = false
title = 'Zod 的学习与使用：定义模式（二）'
categories = ['前端技术', 'Zod']
tags = ['Zod', 'TypeScript']
toc = true
+++

- 学习资料：[https://zod.dev/api](https://zod.dev/api)
- 辅助：豆包
- Zod 版本：4

模式，也就是数据结构，使用 Zod 时，我们要先定义模式，即定义数据结构的规则，然后调用 `.parse` 或 `safeParse` 进行校验。

## 开发环境的搭建

本节要搭建一个能运行 Zod 库的 TypeScript 开发环境，步骤如下：

```bash
# 项目目录
$ mkdir learn_zod
$ cd learn_zod
$ npm init -y
# 安装依赖
$ npm install zod
$ npm install -D typescript @types/node ts-node
```

- typescript：TS 编译器
- @types/node：Node.js 的 TS 类型声明
- ts-node：无需手动编译，直接运行 TS 文件

创建 TS 配置文件：

```bash
$ touch tsconfig.json
```

内容如下：

```json
{
    "compilerOptions": {
        "target": "ES2020",
        "module": "CommonJS",
        "outDir": "./dist",
        "rootDir": "./src",
        "strict": true,
        "esModuleInterop": true,
        "skipLibCheck": true,
        "forceConsistentCasingInFileNames": true
    },
    "include": [
        "src/**/*"
    ],
    "exclude": [
        "node_modules"
    ]
}
```

```bash
# 创建测试 ts 文件目录
$ mkdir -p src/tests
```

### 跑个示例

验证开发环境是否最小完备：

```bash
$ touch src/tests/zod_env.ts
```

内容如下：

```ts
import * as z from 'zod'

const Action = z.string().startsWith('Zod')

const action = Action.parse('Zod learning')

console.log(action)
```

然后，在命令行中执行：

```bash
$ npx ts-node src/tests/zod_env.ts
Zod learning
```

上述过程没有报错，即开发环境搭建完毕。

## 用户注册

现在假设用户注册这一场景中，用户需要提供 3 个字段：username, password, confirmPassword，规则如下：

* username 字段仅支持小写字母（a-z）、大写字母（A-Z）、数字（0-9）和下划线（_），长度至少 6 位，最多 12 位
* password 字段支持小写字母（a-z）、大写字母（A-Z）、数字（0-9）和特殊字符（!@#¥%&*），其中小写字母、大写字母、数字、特殊字符必须同时包含，长度至少 8 位，最多 20 位
* confirmPassword 字段的值必须与 password 字段值一致辞

上述需求接近于实际，可以很好地诠释 Zod 的“数据结构校验”的使用过程。

```bash
$ touch src/tests/user_register.ts
```

内容如下：

```ts
import * as z from 'zod'

const UsernameSchema = z.string()
  .min(6, { message: '用户名长度不能少于 6 位' })
  .max(12, { message: '用户名长度不能超过 12 位' })
  .regex(/^[a-zA-Z0-9_]+$/, { message: '用户名只能包含字母、数字和下划线' })

const PasswordSchema = z.string()
  .min(8, { message: '密码长度不能少于 8 位' })
  .max(20, { message: '密码长度不能超过 20 位' })
  .regex(/^[a-zA-Z0-9!@#¥%&*]+$/, {
    message: "密码仅支持大小写字母、数字和特殊字符（!@#¥%&*）"
  })
  .refine(val => /[a-z]/.test(val), {
    message: "密码必须包含小写字母"
  })
  .refine(val => /[A-Z]/.test(val), {
    message: "密码必须包含大写字母"
  })
  .refine(val => /[0-9]/.test(val), {
    message: "密码必须包含数字"
  })
  .refine(val => /[!@#¥%&*]/.test(val), {
    message: "密码必须包含特殊字符（!@#¥%&*）"
  })

const ConfirmPasswordSchema = z.string()
  .nonempty({ message: '确认密码不能为空' })
  
const RegisterFormSchema = z.object({
  username: UsernameSchema,
  password: PasswordSchema,
  confirmPassword: ConfirmPasswordSchema,
})
  .refine(({ password, confirmPassword }) => password === confirmPassword, {
    message: "两次输入的密码不一致",
    path: ['confirmPassword'],
  })

const validForm = {
  username: 'a2htray',
  password: 'P@ssw0rd!',
  confirmPassword: 'P@ssw0rd!',
}

const validResult = RegisterFormSchema.safeParse(validForm)

console.log(validResult.success) // true

const invalidForm = {
  username: 'a2htray',
  password: 'P@ssw0rd!',
  confirmPassword: 'P@ssw0r',
}

const invalidResult = RegisterFormSchema.safeParse(invalidForm)

console.log(invalidResult.success) // false

const issue = invalidResult.error?.issues[0]
console.log(`字段 ${issue?.path}：${issue?.message}`) // 字段 confirmPassword：两次输入的密码不一致
```

* 定义 3 个 Schema（UsernameSchema、PasswordSchema、ConfirmPasswordSchema），再组合成 1 个 RegisterFormSchema
  * 调用基础类型（`z.string`）、字符中格式（`z.min`、`z.max`）和自定义校验（`z.refine`），声明一个 Schema
* 调用 Schema 的 `.safeParse` 方法进行值校验
* 通过校验结果中的 `success` 字段判断是否校验成功，若校验失败，则可以通过 `error` 字段进行捕获

用户注册这一需求，从代码实现来看，做了 3 件事情：

1. 定义模式
2. 模式校验值
3. 获取校验结果

## 特别点

在学习资源中（Zod 的官方文档中），有着许多内置的数据类型、字符串格式，绝大部分都只有在用到时才会看，但有一些特别的点，是一定要有点印象。

**Coercion（胁迫类型）**

<span style="color: blue;">The coerced variant of these schemas attempts to convert the input value to the appropriate type.</span>

定义的胁迫类型，会尝试将输入值转换成合适的类型。

```bash
$ touch src/tests/coerced_variant.ts
```

内容如下：

```ts
import * as z from 'zod'

const BooleanSchema = z.boolean()
const CoercedBooleanSchema = z.coerce.boolean()

// console.log(BooleanSchema.parse('true')) // 报错
// console.log(BooleanSchema.parse('false')) // 报错
// console.log(BooleanSchema.parse(null)) // 报错
// console.log(BooleanSchema.parse('123')) // 报错
console.log(BooleanSchema.parse(true)) // true
console.log(BooleanSchema.parse(false)) // false

console.log(CoercedBooleanSchema.parse('true')) // true
console.log(CoercedBooleanSchema.parse('false')) // true
console.log(CoercedBooleanSchema.parse(null)) // false
console.log(CoercedBooleanSchema.parse('123')) // true
console.log(CoercedBooleanSchema.parse(true)) // true
console.log(CoercedBooleanSchema.parse(false)) // false
```

**字面量结构定义**

如何要约束一个值是几个值中的某一个，要使用 `z.literal` 方法。

```ts
import * as z from 'zod'

const MyLiteralSchema = z.literal(['red', 'green', 'yellow', 1, 2, 3])

MyLiteralSchema.parse('green') // 'green'
MyLiteralSchema.parse(1) // 1
// MyLiteralSchema.parse('blue') // 报错
// MyLiteralSchema.parse(4) // 报错
```

**字符串格式中，传递一个参数修改校验的行为**

不同字符串格式支持的参数字段各有不同，用到时需要仔细看文档。

```ts
import * as z from 'zod'

const UrlSchema = z.url()
 
UrlSchema.parse("https://example.com")
UrlSchema.parse("http://localhost")
UrlSchema.parse("mailto:noreply@zod.dev")

const CustomUrlSchema = z.url({
  hostname: /^example\.com$/,
  message: "主机名必须是 example.com",
})
 
CustomUrlSchema.parse("https://example.com")
// CustomUrlSchema.parse("https://zombo.com") // 报错
```

**z.enum 限定输入只能字符串集合中的其中一个**

```ts
import * as z from 'zod'

const ColorSchema = z.enum(["red", "green", "blue"])

ColorSchema.parse("red")
ColorSchema.parse("green")
ColorSchema.parse("blue")
// ColorSchema.parse("yellow") // 报错
```

与 `z.literal` 的区别在于：`z.enum` 集合中的值类型只能是字符串。

**🌱Zod 4 新增 z.stringbool**

字符串的逻辑表示，校验并转换成 boolean 类型。

```ts
import * as z from 'zod'
const strbool = z.stringbool();

strbool.parse("true")         
strbool.parse("1")            
strbool.parse("yes")          
strbool.parse("on")           
strbool.parse("y")            
strbool.parse("enabled")

strbool.parse("false") 
strbool.parse("0")
strbool.parse("no")
strbool.parse("off")
strbool.parse("n")
strbool.parse("disabled")
```

支持的字符串有：

```json
{
  "values": [
    "true",
    "1",
    "yes",
    "on",
    "y",
    "enabled",
    "false",
    "0",
    "no",
    "off",
    "n",
    "disabled"
  ]
}
```

由此可推测，其内部实现可能用了 `z.enum`。

**z.optional, z.nullable, z.nullish 允许输入的不同**

* z.optional 允许输入为 undefined
* z.nullable 允许输入为 null
* z.nullish 允许输入为 undefined 和 null

**z.strictObject 严格模式**

如果出现未知的 key 会报错。

```ts
import * as z from 'zod'

const ObjectSchema = z.strictObject({
  a: z.string(),
  b: z.number(),
})

let object: any = {
  a: "123",
  b: 456,
}

let result = ObjectSchema.safeParse(object)
console.log(result.success)

object = {
  a: "123",
  b: 456,
  c: "789",
}
result = ObjectSchema.safeParse(object)
console.log(result.success)
```



