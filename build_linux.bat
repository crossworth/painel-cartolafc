@echo off
REM cd web/frontend && npm run build & cd ../../

go get -u github.com/gobuffalo/packr/v2/packr2
cd web && packr2 && cd ..

set GOOS=linux
set GOARCH=amd64
go build -o cartola_web_admin cmd/server/main.go
