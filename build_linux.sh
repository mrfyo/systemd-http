export GOARCH=amd64
export GOOS=linux
APP=build/systemd-http
go build -o $APP *.go
chmod +x $APP