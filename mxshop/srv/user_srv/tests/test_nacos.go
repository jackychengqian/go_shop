package main

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"user_srv/global"
)

// InitConfig 从 Nacos 加载配置
func main() {
	// 创建 Nacos 配置
	serverConfigs := []constant.ServerConfig{
		*constant.NewServerConfig(global.NacosConfig.Host, global.NacosConfig.Port),
	}
	clientConfig := *constant.NewClientConfig(
		constant.WithNamespaceId(global.NacosConfig.Namespace),
		constant.WithUsername(global.NacosConfig.User),
		constant.WithPassword(global.NacosConfig.Password),
	)

	// 创建 Nacos 客户端
	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(fmt.Sprintf("创建 Nacos 客户端失败: %v", err))
	}

	// 获取配置
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})
	if err != nil {
		panic(fmt.Sprintf("获取 Nacos 配置失败: %v", err))
	}

	// **🔍 添加日志，检查返回内容**
	fmt.Println("📌 Nacos 配置内容:", content)

	// **如果配置为空，说明 Nacos 里没有正确配置**
	if content == "" {
		panic("❌ Nacos 配置为空！请检查 user-srv.json")
	}

	// 解析 JSON
	err = json.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		panic(fmt.Sprintf("❌ 解析 Nacos 配置失败: %v", err))
	}

	fmt.Println("✅ Nacos 配置加载成功！")
}
