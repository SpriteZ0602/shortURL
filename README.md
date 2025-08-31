# ShortURL - 云原生高性能短链服务

把长网址压缩成 11 位短码，支持 10 万 QPS 重定向 + JWT 用户体系 + 实时 PV/UV 统计。

### 前置

- Go 1.22+
- Docker Desktop
- MySQL & Redis 会自动拉镜像

### 起依赖
```bash
docker run -d --name mysql -e MYSQL_ROOT_PASSWORD=123456 -p 3306:3306 mysql:8
docker run -d --name redis -p 6379:6379 redis:7-alpine
docker run -d --name etcd -p 2379:2379 quay.io/coreos/etcd:v3.5.15 etcd --advertise-client-urls http://0.0.0.0:2379 --listen-client-urls http://0.0.0.0:2379
```

### 跑代码

```bash
go run cmd/shorturl/main.go 
```

### 骨架

```bash
shorturl/
├── cmd/shorturl/            # main.go
├── internal/
│   ├── handler/             # 路由（auth / shorten / redirect）
│   ├── middleware/          # 中间件（jwt）
│   ├── model/               # GORM 模型
│   ├── repo/                # 数据层
│   └── service/             # 业务逻辑
├── pkg/
│   ├── cache/               # Redis 客户端
│   ├── snowflake/           # 雪花 ID 生成器
│   └── jwt/                 # JWT 工具
├── scripts/
│   └── dev-up.bat           # Windows 一键脚本
├── k8s/                     # 暂未使用
├── go.mod & go.sum
└── README.md
```

### 一键脚本 

`scripts/dev-up.bat`



### 实现的服务

| 方法     | 路径               | 请求示例                                                     | 响应示例                   | 说明               |
| -------- | ------------------ | ------------------------------------------------------------ | -------------------------- | ------------------ |
| **POST** | `/api/v1/register` | `{"email":"a@b.com","password":"123456"}`                    | `{"message":"registered"}` | 用户注册           |
| **POST** | `/api/v1/login`    | `{"email":"a@b.com","password":"123456"}`                    | `{"token":"jwt..."}`       | 登录拿 JWT         |
| **POST** | `/api/v1/shorten`  | 头：`Authorization: Bearer <jwt>`<br>体：`{"url":"https://example.com"}` | `{"short_code":"5y4queDsELQ"}`  | 创建短链（需登录） |
| **GET**  | `/{short_code}`    | 浏览器访问 `http://localhost:8080/5y4queDsELQ`                    | 302 → 原长网址             | 公开跳转           |

### 分布式短码生成
分布式 Snowflake + Etcd 
用 Etcd 自动分配机器 ID，支持多实例无冲突发号。 
单表可支持 100 亿条记录，QPS 可达 10K+。

### 可观测链路（OpenTelemetry + Jaeger）

| Span 名称           | 说明                     | 耗时    |
|---------------------|--------------------------|-------|
| handler.Shorten     | 接收并解析请求           | 1 ms  |
| service.Shorten     | 业务/雪花 ID 生成        | 90 ms |
| repo.Save           | MySQL 插入               | 2 ms  |
| handler.Redirect    | 302 跳转                 | 19 ms |
| cache.Redis         | Redis 读写               | <1 ms |

访问 **http://localhost:16686** 即可在 Jaeger UI 查看完整链路图。