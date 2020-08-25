@echo off
set GOOS=linux
set GOARCH=amd64
go build -o cartola_worker github.com/crossworth/cartola-web-admin/cmd/topicworker
