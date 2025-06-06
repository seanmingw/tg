package middleware

import (
	"github.com/sirupsen/logrus"
	"trojan-panel-core/app/trojango"
	"trojan-panel-core/app/xray"
	"trojan-panel-core/app/hysteria"
	"trojan-panel-core/app/hysteria2"
)

// InitTrafficControl 初始化流量控制
func InitTrafficControl() {
	// 初始化流量控制器
	trafficController := GetTrafficController()
	
	// 启动清理任务
	trafficController.StartCleanupTask()
	
	// 初始化协议管理器
	protocolManager := GetProtocolManager()
	
	// 注册各协议处理器
	logrus.Info("正在注册协议处理器...")
	
	// 这里需要根据实际情况注册各协议处理器
	// 例如：
	// protocolManager.RegisterHandler("trojan-go", trojango.NewTrojanGoApi(8100))
	// protocolManager.RegisterHandler("xray-vmess", xray.NewXrayApi(8101))
	// 等等
	
	logrus.Info("流量控制初始化完成")
}