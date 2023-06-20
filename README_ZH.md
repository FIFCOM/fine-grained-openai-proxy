<div align="center">

<h1 align="center">Fine-grained OpenAI API Proxy</h1>

[English](README.md) | 简体中文

具有 细粒度权限控制、有效期控制 和 剩余请求数限制 的OpenAI API代理

⚠️ 注意：当前项目仍处于早期开发阶段，不建议在生产环境中使用

</div>

## 主要功能

- 🔑 创建 OpenAI API 子 Key。子 Key 可以保护您的主 Key，防止被滥用。
- ⏰ 为子 Key 设置有效期。过期后子 Key 将无法使用，需要重新创建或者延长有效期。
- 🔢 为子 Key 设置剩余请求数。每次正确的请求会消耗一个剩余请求数。
- 🔒 为子 Key 设置权限，例如只允许使用 `gpt-3.5-turbo` 和 `gpt-4` 模型，或者不允许使用 `whisper-1` 模型。
- 🌐 OpenAI API 代理。代理会根据子 Key 的权限和剩余请求数来决定是否允许请求，并将请求稍作修改后转发到 OpenAI API。

## 构建

- **克隆本仓库**
```sh
git clone https://github.com/FIFCOM/fine-grained-openai-api-proxy.git
```

- **编辑配置文件**中的默认配置项
```sh
vim conf/config.go
```

```go
// AdminToken 管理员Token，用于访问/admin接口
var AdminToken = "123456789"

// Port 应用监听端口
var Port = ":8080"

// SqlitePath Sqlite数据库文件路径。如非必要，请勿修改
var SqlitePath = "./data/oapi-proxy.db"

// OpenAIAPIBaseUrl OpenAI API 的地址。如有需要，可以修改为其他镜像地址
var OpenAIAPIBaseUrl = "https://api.openai.com"
```

- **编译项目** (需要安装[Golang](https://go.dev/dl/), 推荐使用1.18以上版本)
```sh
cd fine-grained-openai-api-proxy
# 若网络环境不佳，可以使用Go Proxy
# go env -w GOPROXY=https://goproxy.cn,direct
go build .
```

## 使用

-  创建一个OpenAI API Key。可以在此处创建：[https://platform.openai.com/account/api-keys](https://platform.openai.com/account/api-keys)

-  运行项目
```sh
./fine-grained-openai-api-proxy -admin=true -port=8080
```

- 当前项目未实现前端Web节目，因此需使用[Postman](https://www.postman.com/downloads/)与后端进行交互。请导入`./.postman/Fine-grained_OpenAI_API_Proxy.postman_collection.json`到Postman中。

- 初次使用，请初始化数据库。
    1. 选择`apikey/insert`接口，在Body中填入您的OpenAI API Key，点击Send发送请求。这将会保存您的OpenAI API Key到数据库中。
    2. 选择`model/init`接口，在Body中填入 `1` ，并发送请求。后端将获取您所选API Key的所有可用模型列表，并保存到数据库中。
    3. 选择`fgkey/insert`接口，发送请求。后端将会创建一个子 Key，并返回给您(仅展示一次)。您可以使用此子 Key 来访问OpenAI API。
    4. 选择`chat/completion`接口，在Header中填入Authorization，Value为您上一步获取的子 Key。在Body中填入您的请求参数，点击Send发送请求。后端将会将请求转发到OpenAI API，并将OpenAI API的响应返回给您。如果返回结果正常，说明本项目已经可以正常使用了。

- 接下来，您可以在其他项目中配置本项目的OpenAI API代理地址，以及您的子 Key，来使用OpenAI API。

## TODO

- [ ] 前端Web界面
- [ ] Docker镜像支持
- [ ] 项目文档
- [ ] ...

## LICENSE

本项目基于GPLv3协议开源, 详细信息请参考[LICENSE](./LICENSE)文件