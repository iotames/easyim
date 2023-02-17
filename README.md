## 简介

easyim 是一个简单易用，对二次开发友好，方便部署的即时通讯服务器。

如数据量大、对性能有要求，请将`sqlite3` 替换为其他数据库。

代码源于刘丹冰老师视频：[8小时转职Golang工程师](https://www.bilibili.com/video/BV1gf4y1r79E/) - 即时通讯系统


## 开发环境

下载并安装Go: https://golang.google.cn/doc/install

设置GO国内代理:

```
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
```


## 快速开始

```
# 加载依赖包
go mod tidy

# 首次运行，添加初始化参数--init，初始化数据库
go run . --init
```


## 配置文件

复制 `env.default` 文件为 `.env`, 并更改新配置文件 `.env` 的配置项，以覆盖 `env.default` 配置文件的默认值