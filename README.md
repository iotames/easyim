## 简介

easyim 是一个简单易用，二开友好，方便部署的即时通讯服务器。如数据量大、对性能有要求，请自行扩展，并将`sqlite3` 替换为其他数据库。

在线文档: [https://imdocs.catmes.com](https://imdocs.catmes.com)

代码源于刘丹冰老师视频教程：[8小时转职Golang工程师](https://www.bilibili.com/video/BV1gf4y1r79E/) - 即时通讯系统


## 开发环境

下载并安装Go: https://golang.google.cn/doc/install

设置GO国内代理:

```
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.io,direct
```


## 快速开始

```
# 加载依赖包
go mod tidy

# 首次运行，添加初始化参数--init，初始化数据库
go run . --init
```

客户端调试:

```
# linux下使用nc命令调试
nc 127.0.0.1 8888

# 使用本项目 `examples` 目录中示例文件调试
go run client.go -ip 127.0.0.1 -port 8888
```


## 配置文件

复制 `env.default` 文件为 `.env`, 并更改新配置文件 `.env` 的配置项，以覆盖 `env.default` 配置文件的默认值


## 通讯数据格式

IM数据通讯的长连接，支持数据传输 `json`, `protobuf` 两种格式。默认为 `protobuf`。

如需更改，请在 `.env` 文件，添加配置项 `MSG_FORMAT = "json"`，以覆盖默认值。

数据格式说明:

| 字段名 | 数据类型 | 释义   |
| ------ | --------- | -------- |
| id | string |     消息ID，字符串格式     |         |
| seq   | uint32     | 时序号，整型。确保消息按正确顺序显示。客户端按用户发送顺序，从小到大填写seq值。  |
| from_user_id   | string     | 消息发送方ID |
| to_user_id   | string     | 消息接收方ID |
| chat_type   | ChatType     | 聊天类型(ChatType枚举类型:0单聊,1群聊) |
| msg_type   | MsgType     | 消息类型(MsgType枚举类型:0文本,1图片,2语音,3视频) |
| status   | MsgStatus     | 消息状态(MsgStatus枚举类型:0未发送,1已发送,2已送达,3已读取) |
| content   | string     | 消息内容，字符串类型 |

请参看 [protobuf/msg.proto](https://github.com/iotames/easyim/blob/master/protobuf/msg.proto)文件


## 在线文档

全局安装 `docsify` 文档生成工具: 

```
npm i docsify-cli -g
```

- 文档初始化: `docsify init ./docs`

- 本地预览: `docsify serve docs`