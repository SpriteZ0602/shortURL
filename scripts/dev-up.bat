@echo off
docker start mysql redis etcd jaeger >nul
go run cmd/shorturl/main.go