echo "开始构建"
export GOPATH=$WORKSPACE/..
export PATH=$PATH:$GOROOT/bin

# Print Go version
go version

export GO111MODULE=on
export GOPROXY=https://goproxy.io
export ENV=local

echo "current: ${USER}"
#拷贝配置文件到target下
mkdir -vp goods/target/goods
cp goods/config-pro.yaml goods/target/goods/config-pro.yaml

cd $WORKSPACE
go build -o goods/target/main goods/main.go
echo "构建结束2"