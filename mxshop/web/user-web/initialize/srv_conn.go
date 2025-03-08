package initialize

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"

	"user-web/global"
	"user-web/proto"
)

func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo

	// 设置超时时间，避免连接阻塞
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 正确的 DialContext() 使用
	userConn, err := grpc.DialContext(
		ctx,
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()), // 替代 grpc.WithInsecure()
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】:", err)
	}

	// 创建 gRPC 客户端
	global.UserSrvClient = proto.NewUserClient(userConn)
}

func InitSrvConn2() {
	// 从注册中心获取到用户服务的信息
	cfg := api.DefaultConfig()
	consulInfo := global.ServerConfig.ConsulInfo
	cfg.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().Fatal("[InitSrvConn2] 创建 Consul 客户端失败", "error", err)
		return
	}

	// 使用过滤器获取用户服务信息
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.ServerConfig.UserSrvInfo.Name))
	if err != nil {
		zap.S().Fatal("[InitSrvConn2] 获取服务信息失败", "error", err)
		return
	}

	if len(data) == 0 {
		zap.S().Fatal("[InitSrvConn2] 未找到用户服务")
		return
	}

	// 假设第一个服务就是我们要找的用户服务，注意 data 是 map 类型
	var userSrvHost string
	var userSrvPort int
	for _, service := range data {
		userSrvHost = service.Address
		userSrvPort = service.Port
		break // 取第一个服务的信息
	}

	// 使用 context 设置超时，避免连接长时间阻塞
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 更新拨号连接，修复 URL 格式
	userConn, err := grpc.DialContext(
		ctx,
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()), // 替代 WithInsecure()
	)
	if err != nil {
		zap.S().Errorw("[InitSrvConn2] 连接用户服务失败",
			"host", userSrvHost,
			"port", userSrvPort,
			"error", err.Error(),
		)
		return
	}

	// 创建 gRPC 客户端
	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient

	zap.S().Info("[InitSrvConn2] 用户服务连接成功", "host", userSrvHost, "port", userSrvPort)
}
