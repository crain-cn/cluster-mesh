export GOPROXY=https://goproxy.cn
go mod download
go mod vendor

GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o manager cmd/main.go
cp -R deploy/Dockerfile Dockerfile.amd64
docker build -t hub.xesv5.com/zhaoyu10/cluster-mesh:v1.0.2 -f Dockerfile.amd64 .
docker push hub.xesv5.com/zhaoyu10/cluster-mesh:v1.0.2