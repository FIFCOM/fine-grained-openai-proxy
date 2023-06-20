<div align="center">

<h1 align="center">Fine-grained OpenAI API Proxy</h1>

English | [简体中文](README_ZH.md)

OpenAI API proxy with fine-grained access control, expiration control and remaining request limit.

⚠️ Note: Current project is still in the early development stage, and it is not recommended to use it in a production environment.

</div>

## Features

- 🔑 Create OpenAI API sub-keys. Sub-keys can protect your main key from abuse.

- ⏰ Set expiration time for sub-keys. After expiration, the sub-key will be unusable and needs to be recreated or have its validity extended.

- 🔢 Set remaining request count for sub-keys. Each correct request will consume one remaining request count.

- 🔒 Set permissions for sub-keys, such as only allowing the use of `gpt-3.5-turbo` and `gpt-4` models, or disallowing the use of `whisper-1` model.

- 🌐 OpenAI API proxy. The proxy will decide whether to allow requests based on the permissions and remaining request counts of the sub-key, and modify requests slightly before forwarding them to OpenAI API.

## Build

- **Clone repository**

```sh
git clone https://github.com/FIFCOM/fine-grained-openai-api-proxy.git
```

- Edit the default configuration in the **configuration file**
```sh
vim conf/config.go
```

```go
// AdminToken Admin Token, used to access the /admin interface.
var AdminToken = "123456789"

// Port App listening port.
var Port = ":8080"

// SqlitePath Sqlite database file path. DO NOT modify unless necessary.
var SqlitePath = "./data/oapi-proxy.db"

// OpenAIAPIBaseUrl OpenAI API URL. It can be modified to other mirror addresses.
var OpenAIAPIBaseUrl = "https://api.openai.com"
```

- **Compile** (requires installation of [Golang](https://go.dev/dl/), recommended to use version 1.18 or above).
```sh
cd fine-grained-openai-api-proxy
go build .
```

## Usage

- Create an OpenAI API Key. You can create one here: [https://platform.openai.com/account/api-keys](https://platform.openai.com/account/api-keys)

- Run the project

```sh
./fine-grained-openai-api-proxy -admin=true -port=8080
```

- As the current project does not have a frontend web page, so you need to use [Postman](https://www.postman.com/downloads/) to interact with the backend. Please import `./.postman/Fine-grained_OpenAI_API_Proxy.postman_collection.json` into Postman.

- For first-time use, please initialize the database.

    1. Select the `apikey/insert` interface and fill in your OpenAI API Key in the Body section. Click Send to send the request. This will save your OpenAI API Key to the database.
    2. Select the `model/init` interface and fill in `1` in the Body section, then send the request. The backend will retrieve all available models for your selected API Key and save them to the database.
    3. Select `fgkey/insert` interface and send a request. The backend will create a sub-key and return it to you (displayed only once). You can use this sub-key to access OpenAI API.
    4. Select `chat/completion` interface, fill in Authorization with Value as your previously obtained sub-key in Header section, then fill in your request parameters in Body section before clicking Send to send out a request. The backend will forward this request to OpenAI API and return its response back to you. If there is no error message returned, it means that this project is ready for normal use now.

- Now, you can configure the OpenAI API proxy address and your fine-grained sub-key in other projects to use the OpenAI API.

## TODO

- [ ] Front-end web page
- [ ] Docker support
- [ ] Project documentation
- [ ] ...

## LICENSE

This project is open source under the GPLv3 license. For more information, please refer to the [LICENSE](LICENSE) file.