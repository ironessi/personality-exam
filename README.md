# Personality Exam API

一个基于 GoFrame 的人格/心理测评后端服务，提供用户认证、测评题目、测评结果、收藏、短信验证码、订单支付与退款等接口能力，可作为心理测评、小程序、H5 或管理后台的后端接口演示项目。

## 项目描述

Personality Exam API 是一个人格测评系统后端示例项目，内置多种测评量表与结果计算逻辑，支持用户登录、题目获取、答案提交、结果生成、收藏记录、支付下单、二维码支付、退款处理和每日对账等业务流程。项目使用 GoFrame v2 构建，集成 MySQL、Redis、JWT、阿里云短信、支付宝和微信支付相关模块，适合用于心理测评平台、在线问卷测评、付费测试产品的后端服务展示与二次开发参考。

## 功能特性

- 用户认证：支持短信验证码、登录认证和 JWT 中间件保护。
- 测评管理：支持测评分类、测评类型列表、测评题目查询。
- 答案提交：支持用户提交测评答案并生成对应结果。
- 测评历史：支持查看用户历史测评记录和结果详情。
- 收藏功能：支持收藏测评结果、查询收藏列表、统计收藏数据。
- 支付能力：支持订单创建、订单查询、二维码支付、支付状态查询、取消订单。
- 退款能力：支持退款申请、退款记录查询和待退款任务定时处理。
- 支付回调：预留支付宝、微信支付异步通知接口。
- 定时任务：包含每日对账任务和周期性退款处理任务。
- 配置隔离：真实配置文件、证书和日志默认不进入 Git 仓库，只提交示例配置。

## 技术栈

- Go 1.25+
- GoFrame v2
- MySQL
- Redis
- JWT
- 阿里云短信服务
- 支付宝支付 SDK
- 微信支付 PowerWeChat

## 目录结构

```text
.
├── api/                    # API 请求/响应结构定义
├── config/                 # 本地配置目录，仅提交 example 配置
├── hack/                   # GoFrame CLI 与构建辅助配置
├── internal/
│   ├── cmd/                # 服务启动和路由注册
│   ├── controller/         # HTTP 控制器
│   ├── dao/                # 数据访问对象
│   ├── logic/              # 核心业务逻辑
│   ├── model/              # 数据模型
│   └── service/            # 服务接口定义
├── manifest/
│   ├── config/             # 服务运行配置，仅提交 example 配置
│   └── docker/             # Docker 构建文件
├── main.go                 # 程序入口
├── public.sql              # 示例数据库结构/数据
└── README.md
```

## 配置说明

为了避免泄露数据库密码、Redis 密码、短信密钥、支付私钥、支付证书和日志文件，真实配置文件已经被 `.gitignore` 忽略。仓库中只保留示例配置文件：

```text
config/config.example.yaml
hack/config.example.yaml
manifest/config/config.example.yaml
manifest/config/payment.example.yaml
```

本地运行前，复制示例配置并替换其中的 `YOUR_...` 占位符：

```bash
cp config/config.example.yaml config/config.yaml
cp hack/config.example.yaml hack/config.yaml
cp manifest/config/config.example.yaml manifest/config/config.yaml
cp manifest/config/payment.example.yaml manifest/config/payment.yaml
```

需要重点配置的内容：

- MySQL 地址、端口、账号、密码、数据库名。
- Redis 地址、数据库编号和密码。
- 阿里云短信 `accessKeyId`、`accessKeySecret`、签名和模板编号。
- 支付宝 `app_id`、应用私钥、证书路径、回调地址、跳转地址。
- 微信支付商户号、API v3 密钥、证书路径、证书序列号和回调地址。

支付证书建议放在本地 `certs/` 目录下，该目录默认不会提交到 Git。

## 快速开始

### 1. 克隆项目

```bash
git clone git@github.com:ironessi/personality-exam.git
cd personality-exam
```

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 准备数据库

创建 MySQL 数据库后，可以根据项目中的 `public.sql` 初始化表结构和演示数据。

```bash
mysql -u YOUR_DATABASE_USER -p YOUR_DATABASE_NAME < public.sql
```

### 4. 准备配置文件

```bash
cp config/config.example.yaml config/config.yaml
cp hack/config.example.yaml hack/config.yaml
cp manifest/config/config.example.yaml manifest/config/config.yaml
cp manifest/config/payment.example.yaml manifest/config/payment.yaml
```

然后根据自己的本地环境修改真实配置文件。

### 5. 启动服务

```bash
go run main.go
```

默认服务地址由 `manifest/config/config.yaml` 中的 `server.address` 控制，示例配置为：

```text
0.0.0.0:8000
```

启动后可以访问健康检查接口：

```bash
curl http://127.0.0.1:8000/api/health
```

返回：

```text
OK
```

## 初始化测评数据

项目内置多种测评量表的初始化逻辑。配置好数据库后，可以通过启动参数导入数据：

```bash
go run main.go -init
```

如果需要强制重新初始化：

```bash
go run main.go -init -force
```

## 常用接口

以下是部分路由概览，具体请求和响应结构可以参考 `api/` 目录。

```text
GET    /api/health                              健康检查
POST   /api/v1/auth/register                    用户注册
POST   /api/v1/auth/login                       用户登录
POST   /api/v1/sms/send                         发送短信验证码
POST   /api/v1/sms/verify                       验证短信验证码
GET    /api/v1/sms/status/{mobile}              查询短信状态
GET    /api/categories/index                    获取测评分类和分类下测评
GET    /api/tests/types                         获取测评类型列表
GET    /api/tests/{id}                          获取测评题目
GET    /api/user/tests                          获取用户测评历史
GET    /api/user/tests/{record_id}              获取测评结果详情
POST   /api/user/tests/save                     保存测评进度
POST   /api/user/tests/submit                   提交测评结果
GET    /api/user/favorites                      获取收藏列表
POST   /api/user/favorites                      添加收藏
POST   /api/v1/payment/orders                   创建订单
GET    /api/v1/payment/orders/{orderNo}         查询订单
POST   /api/v1/payment/qrcode                   创建支付二维码
GET    /api/v1/payment/qrcode/status            查询支付状态
POST   /api/v1/payment/refunds                  申请退款
POST   /api/v1/payment/notify/alipay            支付宝支付回调
POST   /api/v1/payment/notify/wechat            微信支付回调
```

需要登录的接口会经过 JWT 中间件校验，请在请求头中携带登录后获得的认证信息。

## 支付说明

演示配置中支付渠道默认可以关闭或使用占位配置。若需要完整测试支付流程，请完成以下准备：

- 开通支付宝或微信支付商户能力。
- 配置公网可访问的支付异步回调地址。
- 准备支付平台要求的证书文件。
- 将证书放入本地 `certs/` 目录。
- 在 `manifest/config/payment.yaml` 中填写真实支付参数。

`certs/` 已被 Git 忽略，请不要提交真实证书或私钥。

## 安全注意事项

请不要提交以下内容：

- `config/config.yaml`
- `hack/config.yaml`
- `manifest/config/config.yaml`
- `manifest/config/payment.yaml`
- `certs/`
- `logs/`
- 数据库密码
- Redis 密码
- 阿里云 AccessKey
- 支付宝应用私钥
- 微信支付 API 密钥
- 任何生产环境证书或日志

如果配置文件曾经被提交到公开仓库，建议立即更换相关密钥和密码。

## 开发检查

运行全部 Go 包测试/编译检查：

```bash
go test ./...
```

当前仓库主要是业务接口和逻辑代码，测试文件较少，因此该命令更多用于确认项目可以正常编译。

## 适用场景

- 心理测评平台后端接口演示。
- 人格测试、小程序或 H5 项目的服务端参考。
- GoFrame 单体服务项目结构参考。
- 付费测评、订单支付、退款、对账流程的后端示例。

## License

本项目用于学习、演示和二次开发参考。实际商用前请自行检查依赖协议、支付平台要求、隐私合规和数据安全要求。
