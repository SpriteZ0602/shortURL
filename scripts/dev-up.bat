@echo off
docker start mysql redis etcd >nul
go run cmd/shorturl/main.go