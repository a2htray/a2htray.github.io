---
title: "Protocol Buffer 语法"
date: 2022-04-08T16:33:52+08:00
draft: false
reward: false
categories:
  - 后端技术
  - Go
tags:
  - Go
  - protobuf
---

Protocol Buffer（Protobuf） 是一种高效的数据结构序列化的机制，同时也是一种结构化数据的存储格式。

> 序列化与反序列化
>
> * 序列化：将数据结构或对象转换成二进制串的过程；
> * 反序列化：将序列化后的二进制串转换成数据结构或对象的过程；

<!-- more -->

## 语法

```protobuf
/*
* 语法
*/

/*
* 指定 Protobuf 解析使用的版本，可以是 proto3 或 proto2
*/
syntax = "proto3";

/*
* message 定义中的每一个字段都有一个唯一标识，该标识用于在二进制格式中识别字段
* 字段的标识一旦使用就不要进行修改
* 当标识为 1 到 15 时，使用一个字节进行编码，字节信息中包含字段的标识以及类型
* 当标识为 16 到 2047 时，使用两个字节进行编码
* Field numbers in the range 16 through 2047 take two bytes. So you should reserve the numbers 1 through 15 for very frequently occurring message elements.
* 在编码的过程中，标识使用应当留有余地，便于将来扩展
* 标识最小的值是 1，最大的值为 2^29-1
* 不能使用的标识为 19000 到 19999
* 不能再使用已经被 reserved 的标识
*/

/*
定义 message 的语法：
    message ${MessageName} {
        ${Scalar Value Type} ${FieldName1} = ${Tag Number1};
                .
                .
                .
        ${Scalar Value Type} ${FieldNameN} = ${Tag NumberN};
    }
*/

message MessageTypes {
    /*
    * 标量值类型
    */
    string stringType = 1; // 字符串可以是 UTF-8 编码，也可以是一个 7 比特的 ASCII 字符，默认为“”
    // 数值类型，默认为 0
    int32 int32Type = 2; // 使用变量长度进行编码，如果是负数，请使用 sint32
    int64 int64Type = 3; // 使用变量长度进行编码，如果是负数，请使用 sint64
    uint32 uInt32Type = 4; // 使用变量长度进行编码
    uint64 uInt64Type = 5; // 使用变量长度进行编码
    sint32 sInt32Type = 6; // 使用变量长度进行编码，处理负数更高效
    sint64 sInt64Type = 7; // 使用变量长度进行编码，处理负数更高效

    fixed32 fixed32Type = 8; // 变量总是占 4 个字节，当值大于 2^28 时，比使用 uint32 更有效率
    fixed64 fixed64Type = 9; // 变量总是占 8 个字节，当值大于 2^56 时，比使用 uint64 更有效率

    sfixed32 sfixed32Type = 10; // 变量总是占 4 个字节
    sfixed64 sfixed64Type = 11; // 变量总是占 8 个字节

    bool boolType = 12; // 布尔类型，默认为 false

    bytes bytesType = 13; // 可包含任意长度的字节数组，默认为长度为 0 的字节数组

    double doubleType = 14;
    float floatType = 15;

    enum Week {
        UNDEFINED = 0; // 第 1 个值
        SUNDAY = 1;
        MONDAY = 2;
        TUESDAY = 3;
        WEDNESDAY = 4;
        THURSDAY = 5;
        FRIDAY = 6;
        SATURDAY = 7;
    }
    Week wkDayType = 16;

    /*
    * 定义标量值类型的集合
    * Syntax: repeated ${ScalarType} ${name} = TagValue
    */
    repeated string listOfString = 17; // List[String]
}

/*
* 在其它 message 中使用已定义的 message
*/
message Person {
    string fname = 1;
    string sname = 2;
}

message City {
    Person p = 1;
}

/*
* 嵌套的 message 定义
*/
message NestedMessages {
    message FirstLevelNestedMessage {
        string firstString = 1;
        message SecondLevelNestedMessage {
            string secondString = 2;
        }
    }
    FirstLevelNestedMessage msg = 1;
    FirstLevelNestedMessage.SecondLevelNestedMessage msg2 = 2;
}

/*
* .proto 文件的引入
*/

// one.proto
// message One {
//     string oneMsg = 1;
// }

// two.proto
//  import "myproject/one.proto"
//  message Two {
//       string twoMsg = 2;
//  }


/*
* 高级知识点
*/

/*
* message 发生改变时，永远不要修改或使用已经删除字段的标识
*/

/*
* 使用 reserved 保留已删除的标识或字段名
*/
message ReservedMessage {
    reserved 0, 1, 2, 3 to 10; // 这里的标识不可再使用
    reserved "firstMsg", "secondMsg", "thirdMsg"; // 这里的字段名不可再使用
}

/*
* 引用其它文件中定义的 message
*/
import "google/protobuf/any.proto";
message AnySampleMessage {
    repeated google.protobuf.Any.details = 1;
}

/*
*  OneOf
* 相同于 union，只能是其中一个
* 使用 oneof 的 message 不能被 repeated
*/
message OneOfMessage {
    oneof msg {
        string fname = 1;
        string sname = 2;
    };
}

/*
* Maps
* map 字段不能被 repeated
*/
message MessageWithMaps {
    map<string, string> mapOfMessages = 1;
}


/*
* Packages
* 声明一个包名，防止同名的 message
* 语法:
    package ${packageName};

    访问方式
    ${packageName}.${messageName} = ${tagNumber};
*/

/*
* 在 RPC 系统中使用，其中可以定义方法
*/
message SearchRequest {
    string queryString = 1;
}

message SearchResponse {
    string queryResponse = 1;
}
service SearchService {
    rpc Search (SearchRequest) returns (SearchResponse);
}
```

### 数据类型

Protobuf 内置的数据类型以及在 Go 中对应的数据类型：

| Protobuf | Go      |
| -------- | ------- |
| double   | float64 |
| float    | float32 |
| int32    | int32   |
| int64    | int64   |
| uint32   | uint32  |
| uint64   | uint64  |
| sint32   | int32   |
| sint64   | int64   |
| fixed32  | uint32  |
| fixed64  | uint64  |
| sfixed64 | int64   |
| bool     | bool    |
| string   | string  |
| bytes    | []byte  |

## 书写规范

代码编写要遵循某种特定的规则，如 Python 的 PEP8，`.proto` 文件的内容也应该按照统一的格式，这样才能规范团队编码风格，易于他人理解。

* message 采用驼峰命名法且首字母大写；

```protobuf
message UserConfig {} // 符合
message user_config {} // 不符合
```

* 字段名采用下划线分隔命名法且均小写；

```protobuf
message Product {
    string fasta_origin = 1;
} // 符合
message Product {
    string FastaOrigin = 1
} // 不符合
```

* 枚举型命名格式与 message 相同，枚举值全部大写，并且用下划线分隔命名法；

```protobuf
enum Week {
    MY_MONDAY = 0
} // 符合
enum Week {
    MyMonday = 0
} // 不符合
```

* service 和方法名都采用驼峰命名法，并且首字母大写；

```protobuf
service Greeter {
    rpc SayHello(HelloRequest) returns (HelloReply);
} // 符合
service Greeter {
    rpc say_hello(HelloRequest) returns (HelloReply);
} // 不符合
```

## 与 JSON 的区别

| Protobuf                                                     | JSON                         |
| ------------------------------------------------------------ | ---------------------------- |
| 由 Google 开发、用于序列化和反序列化结构化数据的高效编码方式 | 轻量型的数据交换格式         |
| 可自定义规则、方法的消息格式                                 | 仅仅只是一种消息格式         |
| 二进制格式                                                   | 文本格式                     |
| 支持的语言有限                                               | 绝大部分语言都支持           |
| 微服务间数据传输的格式                                       | WEB 应用与服务器间的传输格式 |
| 相同数据序列化后的数据量较小                                 | 相同数据序列化后的数据量较大 |

## 参考

1. [https://learnxinyminutes.com/docs/protocol-buffer-3/](https://learnxinyminutes.com/docs/protocol-buffer-3/)
2. [https://zhuanlan.zhihu.com/p/95869546](https://zhuanlan.zhihu.com/p/95869546)

