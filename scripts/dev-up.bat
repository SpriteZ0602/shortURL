@echo off
docker start mysql redis >nul
go run cmd/shorturl/main.go