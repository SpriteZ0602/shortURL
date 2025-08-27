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

shorturl
├── cmd/shorturl	   # 服务入口
├── internal
│   ├── model		# GORM 模型
│   ├── repo       	   # 数据层
│   └── service    	  # 业务逻辑
├── api              	    # Gin 路由
├── k8s                 	# 容器编排 yaml
├── scripts            	# 一键脚本
└── README.md