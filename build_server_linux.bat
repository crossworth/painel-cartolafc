@echo off
REM cd web/frontend && npm run build & cd ../../
REM go get -u github.com/gobuffalo/packr/v2/packr2

set GOOS=linux
set GOARCH=amd64
cd web && packr2 && cd ..
go build -o cartola_web_admin cmd/server/main.go
