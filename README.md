# ShortURL - 云原生高性能短链服务

把长网址压缩成 6 位短码，支持 10 万 QPS 重定向 + 实时 PV/UV 统计。

### 前置
- Go 1.22+
- Docker Desktop
- MySQL & Redis 会自动拉镜像

### 起依赖
```bash
docker run -d --name mysql -e MYSQL_ROOT_PASSWORD=123456 -p 3306:3306 mysql:8
docker run -d --name redis -p 6379:6379 redis:7-alpine
```

### 跑代码

```bash
go run cmd/shorturl/main.go 
```

### 骨架

```bash
shorturl/
├── cmd/shorturl/            # 服务入口（main.go）
├── internal/
│   ├── handler/             # Gin 路由 / 控制器
│   │   ├── redirect_handler.go
│   │   └── shorten_handler.go
│   ├── model/               # GORM 模型
│   ├── repo/                # 数据层
│   ├── service/             # 业务逻辑
│   └── util/shortid/        # 短码生成器
├── pkg/cache/               # Redis 客户端
│   └── redis.go
├── scripts/                 # 一键脚本
│   └── dev-up.bat           # Windows 一键起服务
├── k8s/                     # 容器编排 yaml
├── go.mod & go.sum
└── README.md
```

### 一键脚本 

`scripts/dev-up.bat`

### 压测结果

wrk 压测 **9.5 万 QPS / 12.9 ms P99**