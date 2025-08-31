# ShortURL - 云原生高性能短链服务

## 项目介绍

ShortURL 是一个基于 Go 语言开发的高性能、云原生短链接服务，支持将长网址压缩成 11 位短码，并提供完整的用户管理、访问统计和监控链路。

## 核心特性

- **高性能重定向**：支持 10 万 QPS 的短链接重定向能力
- **完整用户体系**：基于 JWT 的用户认证和授权系统
- **分布式短码生成**：使用 Snowflake + Etcd 实现分布式无冲突发号
- **实时统计分析**：支持短链接的 PV/UV 实时统计
- **风险控制**：内置 URL 黑名单检查机制
- **使用配额限制**：支持用户每日创建短链接数量限制
- **可观测性**：集成 OpenTelemetry + Jaeger 链路追踪
- **容器化部署**：完整支持 Docker 容器化部署

## 技术栈

- **后端框架**：Gin Web Framework
- **数据库**：MySQL + Redis
- **认证**：JWT (JSON Web Token)
- **服务发现**：Etcd
- **链路追踪**：OpenTelemetry + Jaeger
- **ID 生成**：分布式 Snowflake 算法
- **构建工具**：Go Modules

## 快速开始

### 前置条件

- Go 1.22+ 开发环境
- Docker Desktop
- Git

### 环境准备

#### 1. 克隆项目代码

```bash
git clone https://your-repository-url/shortURL.git
cd shortURL
```

#### 2. 启动依赖服务

```bash
docker run -d --name mysql -e MYSQL_ROOT_PASSWORD=123456 -p 3306:3306 mysql:8
docker run -d --name redis -p 6379:6379 redis:7-alpine
docker run -d --name etcd -p 2379:2379 quay.io/coreos/etcd:v3.5.15 etcd --advertise-client-urls http://0.0.0.0:2379 --listen-client-urls http://0.0.0.0:2379
docker run -d --name jaeger -p 16686:16686 -p 14268:14268 jaegertracing/all-in-one:latest
```

#### 3. Windows 用户一键启动

```bash
scripts/dev-up.bat
```

### 运行服务

```bash
go run cmd/shortURL/main.go
```

服务启动后，可访问以下地址：
- API 服务：http://localhost:8080
- Jaeger 链路追踪 UI：http://localhost:16686

## 项目结构

```
shortURL/
├── cmd/shortURL/            # 程序入口
├── internal/                # 内部包（不对外暴露）
│   ├── handler/             # 控制器层，处理HTTP请求
│   │   ├── auth.go          # 认证相关接口
│   │   ├── short_handler.go # 短链接生成接口
│   │   └── redirect_handler.go # 短链接重定向接口
│   ├── middleware/          # 中间件
│   │   ├── auth.go          # JWT认证中间件
│   │   ├── quota.go         # 配额限制中间件
│   │   ├── risk.go          # 风险控制中间件
│   │   └── trace.go         # 链路追踪中间件
│   ├── model/               # 数据模型
│   │   ├── short_url.go     # 短链接模型
│   │   └── user.go          # 用户模型
│   ├── repo/                # 数据访问层
│   │   ├── short_url_repo.go # 短链接数据操作
│   │   └── user_repo.go     # 用户数据操作
│   └── service/             # 业务逻辑层
│       └── short_service.go # 短链接核心业务逻辑
├── pkg/                     # 公共包（可被其他项目引用）
│   ├── cache/               # Redis缓存客户端
│   ├── jwt/                 # JWT工具包
│   ├── snowflake/           # 雪花ID生成器
│   └── trace/               # 链路追踪初始化
├── scripts/                 # 脚本工具
│   └── dev-up.bat           # Windows开发环境一键启动脚本
├── k8s/                     # Kubernetes部署配置（预留）
├── go.mod & go.sum          # Go模块依赖
└── README.md                # 项目说明文档
```

## API 接口文档

### 1. 用户注册

**请求**:
- 方法: POST
- 路径: `/api/v1/register`
- 请求体: `{"email":"a@b.com","password":"123456"}`

**响应**:
- 成功: `{"message":"registered"}`
- 失败: 相应错误信息

### 2. 用户登录

**请求**:
- 方法: POST
- 路径: `/api/v1/login`
- 请求体: `{"email":"a@b.com","password":"123456"}`

**响应**:
- 成功: `{"token":"jwt..."}` （返回JWT令牌）
- 失败: 相应错误信息

### 3. 创建短链接

**请求**:
- 方法: POST
- 路径: `/api/v1/shorten`
- 请求头: `Authorization: Bearer <jwt_token>`
- 请求体: `{"url":"https://example.com"}`

**响应**:
- 成功: `{"short_code":"5y4queDsELQ"}`
- 失败: 相应错误信息

### 4. 短链接重定向

**请求**:
- 方法: GET
- 路径: `/{short_code}`
- 示例: `http://localhost:8080/5y4queDsELQ`

**响应**:
- 成功: 302 重定向到原始长链接
- 失败: 404 未找到

## 配置说明

### Redis 配置

系统默认连接到本地 Redis 服务（localhost:6379），使用 0 号数据库。如需修改，可在 `pkg/cache/redis.go` 文件中调整配置。

### JWT 配置

JWT 令牌默认使用硬编码密钥 `change_me`，并设置 24 小时有效期。生产环境中请修改为安全的密钥。

### 配额限制

系统默认限制每个用户每日创建 100 个短链接，可在主程序中通过 `middleware.NewQuota()` 函数参数调整。

## 风险控制

系统内置 URL 黑名单检查机制，会在创建短链接时验证 URL 是否在黑名单中。黑名单存储在 Redis 的 `blacklist` 集合中。

## 性能指标

### 核心业务耗时

| 操作               | 平均耗时    | 说明                      |
|------------------|---------|-------------------------|
| handler.Shorten  | 1 ms    | 接收并解析HTTP请求             |
| service.Shorten  | 90 ms   | 业务逻辑处理与雪花ID生成           |
| repo.Save        | 2 ms    | MySQL数据插入              |
| handler.Redirect | 19 ms   | 短链接重定向操作                |
| cache.Redis      | < 1 ms  | Redis缓存读写               |

### 容量规划
- 单表支持 100 亿条短链接记录
- 理论最大 QPS 可达 10K+
- 短码碰撞概率极低

## 可观测性

### 链路追踪

项目集成了 OpenTelemetry 和 Jaeger，可实时监控请求链路和性能瓶颈。

访问 **http://localhost:16686** 即可在 Jaeger UI 查看完整链路图。

### 常见问题排查

1. **Redis 缓存查看**：
   ```bash
docker exec -it redis redis-cli KEYS "*"
   ```
   可查看系统中的所有缓存键值，包括短链接映射、黑名单和用户配额信息。

2. **JWT 认证问题**：确保请求头中正确设置 `Authorization: Bearer <token>` 格式。

## 开发指南

### 代码规范
- 遵循 Go 标准代码风格
- 使用 Go Modules 管理依赖
- 新增功能需添加相应的单元测试

### 构建与部署

#### 本地构建

```bash
go build -o shortURL cmd/shortURL/main.go
```

#### 容器化部署（预留）

项目包含 k8s 目录，未来将支持 Kubernetes 集群部署。