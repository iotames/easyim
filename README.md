## 简介

easyim 是一个简单易用，二开友好，方便部署的即时通讯服务器。如数据量大、对性能有要求，请自行扩展，并将`sqlite3` 替换为其他数据库。

在线文档: [https://imdocs.catmes.com](https://imdocs.catmes.com)

代码源于刘丹冰老师视频教程：[8小时转职Golang工程师](https://www.bilibili.com/video/BV1gf4y1r79E/) - 即时通讯系统


## 客户端

- Flutter [https://github.com/dou23/easy_im](https://github.com/dou23/easy_im)


## 开发环境

下载并安装Go: https://golang.google.cn/doc/install

```
# 开启 module 功能
go env -w GO111MODULE=on
# 设置GO国内代理. 若执行 go mod tidy 命令提示模块下载失败. 请更换模块代理再重试 go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOPROXY=https://goproxy.io,direct
```

Windows环境下，如发现 `cgo` 报错，可能为 `sqlite3` 组件编译错误。请安装C编译器 [TDM-GCC](https://jmeubank.github.io/tdm-gcc/download/) 或 [Mingw-w64](https://github.com/niXman/mingw-builds-binaries/releases/latest)


## 快速开始

```
# 加载依赖包
go mod tidy

# 首次运行，添加初始化参数--init，初始化数据库
go run . --init

# 正常开发时调试运行
go run .

# 编译为二进制文件. 然后直接运行./easyim(easyim.exe)
go build .
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
| msg_type   | MsgType     | 消息类型(MsgType枚举类型:0文本,1图片,2语音,3视频,4事件,5系统通知) |
| status   | MsgStatus     | 消息状态(MsgStatus枚举类型:0未发送,1已发送,2已送达,3已读取) |
| content   | string     | 消息内容，字符串类型 |
| access_token | string | 鉴权令牌，字符串类型。用户登录成功后获取|

请参看 [protobuf/msg.proto](https://github.com/iotames/easyim/blob/master/protobuf/msg.proto)文件

```
protoc --go_out=model protobuf/msg.proto
```

### 事件消息

当 `msg_type=4` 时，代表当前通讯数据为一条`事件消息`。 定义如下表所示:

| 事件ID | 释义  |
| ----- | ----- |
|  KEEP_ALIVE  | 心跳事件 |

`事件ID` 为字符串类型，代表事件类型(如心跳事件)

客户端上报事件时 `to_user_id` 填写 `事件ID`, `content` 填写 `事件值`（事件内容）

服务端下发事件时 `from_user_id` 填写 `事件ID`, `content` 填写 `事件值`（事件内容）

具体如下所示:

| 事件              | from_user_id  | to_user_id    | content |
| -----             | -----        | -----      | -----       |
| 客户端上报心跳事件 | 客户端user_id | KEEP_ALIVE | 客户端IP |
| 服务端响应心跳事件 | KEEP_ALIVE | 客户端user_id | SUCCESS |


### 心跳消息

`心跳消息` 是一种特殊的 `事件消息`， 用于客户端向服务端发送周期性消息。因为服务端长期未收到消息，会主动断开连接。

用户一直未主动发消息，又要保持与服务端的长连接不断开，才能持续接收在线消息。故客户端要周期性地发送心跳事件消息。

发送方式，请参看 `事件消息` 介绍。

服务端设有因长期未收到消息，而主动断开连接的 `等待时间` , 客户端`心跳的发送周期`，稍小于该时间即可。


## 在线调试

1. 在 `.env` 文件中，设置 `MSG_FORMAT = "json"`
2. 进入 websocket 在线调试网页: [https://websocketking.com/](https://websocketking.com/) 或 [http://www.websocket-test.com/](http://www.websocket-test.com/)
3. 输入连接地址. 如: `ws://127.0.0.1:8888`. 然后点击连接
4. 连接成功后，在对话框中，发送符合规范的 json 格式数据，查看服务器数据响应。
5. 如发送数据不符合规范，服务器会断开连接。

json格式数据示例:

```
{"id":"1630392697653039104","seq":0,"from_user_id":"1630381388895096832","to_user_id":"1630381388895666666","chat_type":0,"msg_type":0,"status":1,"content":"hello word.can you received ?","access_token":"aa.bbb.cc"}
```


## 在线文档

### docsify 文档工具

```
# 全局安装 docsify 文档生成工具
npm i docsify-cli -g
# 文档初始化
docsify init ./docs
# 本地预览. 默认地址: http://localhost:3000
docsify serve docs
```

### apidoc 文档工具

```
# 全局安装 apidoc 文档工具，可以根据代码注释，生成API文档
npm install -g apidoc
# 扫描 server/handler 目录中的代码注释，在 docs/apidoc 目录生成API文档
apidoc -i server/handler -o docs/apidoc
# 配置文件 apidoc.json
# 默认地址: http://localhost:3000/apidoc/index.html
```