package middleware

import (
	"errors"
	"net"
)

// ProtocolHandler 协议处理器接口
type ProtocolHandler interface {
	// HandleConnection 处理用户连接
	HandleConnection(hash, clientIP string, conn net.Conn) error
	// GetProtocolName 获取协议名称
	GetProtocolName() string
}

// ProtocolManager 协议管理器
type ProtocolManager struct {
	handlers map[string]ProtocolHandler
}

// NewProtocolManager 创建协议管理器
func NewProtocolManager() *ProtocolManager {
	return &ProtocolManager{
		handlers: make(map[string]ProtocolHandler),
	}
}

// RegisterHandler 注册协议处理器
func (pm *ProtocolManager) RegisterHandler(protocolName string, handler ProtocolHandler) {
	pm.handlers[protocolName] = handler
}

// GetHandler 获取协议处理器
func (pm *ProtocolManager) GetHandler(protocolName string) (ProtocolHandler, bool) {
	handler, exists := pm.handlers[protocolName]
	return handler, exists
}

// HandleConnection 统一处理连接
func (pm *ProtocolManager) HandleConnection(protocolName, hash, clientIP string, conn net.Conn) error {
	handler, exists := pm.GetHandler(protocolName)
	if !exists {
		return errors.New("unsupported protocol: " + protocolName)
	}
	return handler.HandleConnection(hash, clientIP, conn)
}

// 全局协议管理器实例
var globalProtocolManager *ProtocolManager

// GetProtocolManager 获取全局协议管理器
func GetProtocolManager() *ProtocolManager {
	if globalProtocolManager == nil {
		globalProtocolManager = NewProtocolManager()
	}
	return globalProtocolManager
}