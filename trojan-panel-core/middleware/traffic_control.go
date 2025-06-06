package middleware

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
	"trojan-panel-core/dao"
	"trojan-panel-core/dao/redis"
	"trojan-panel-core/model/constant"
	"github.com/sirupsen/logrus"
)

// ConnectionManager 连接管理器
type ConnectionManager struct {
	mutex       sync.RWMutex
	connections map[string]map[string]*ConnectionInfo // userHash -> connectionId -> ConnectionInfo
}

// AddConnection 添加连接到管理器
func (cm *ConnectionManager) AddConnection(userHash, connectionId, clientIP, protocol string, conn net.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	if cm.connections[userHash] == nil {
		cm.connections[userHash] = make(map[string]*ConnectionInfo)
	}
	
	connInfo := &ConnectionInfo{
		UserHash:     userHash,
		ConnectionId: connectionId,
		ClientIP:     clientIP,
		Protocol:     protocol,
		StartTime:    time.Now(),
		Download:     0,
		Upload:       0,
		Conn:         conn,
	}
	
	cm.connections[userHash][connectionId] = connInfo
	
	logrus.Infof("添加连接: 用户=%s, 连接ID=%s, 协议=%s, IP=%s", userHash, connectionId, protocol, clientIP)
}

// ConnectionInfo 连接信息
type ConnectionInfo struct {
	UserHash     string
	ConnectionId string
	ClientIP     string
	Protocol     string
	StartTime    time.Time
	Download     uint64
	Upload       uint64
	Conn         net.Conn
}

// TrafficController 流量控制器
type TrafficController struct {
	connManager *ConnectionManager
	redisClient *redis.RedisClient
}

var (
	globalTrafficController *TrafficController
	once                    sync.Once
)

// GetTrafficController 获取全局流量控制器实例
func GetTrafficController() *TrafficController {
	once.Do(func() {
		globalTrafficController = &TrafficController{
			connManager: &ConnectionManager{
				connections: make(map[string]map[string]*ConnectionInfo),
			},
			redisClient: redis.NewRedisClient(),
		}
	})
	return globalTrafficController
}

// CheckUserQuota 检查用户流量配额
func (tc *TrafficController) CheckUserQuota(userHash string) error {
	account, err := dao.SelectAccountByHash(userHash)
	if err != nil {
		return err
	}
	
	if account == nil {
		return errors.New("用户不存在")
	}
	
	// 检查流量配额（quota < 0 表示无限制）
	if account.Quota != nil && *account.Quota >= 0 {
		totalUsed := 0
		if account.Download != nil {
			totalUsed += *account.Download
		}
		if account.Upload != nil {
			totalUsed += *account.Upload
		}
		
		if totalUsed >= *account.Quota {
			return errors.New("流量配额已用完")
		}
	}
	
	return nil
}

// CheckConnectionLimit 检查用户连接数限制
func (tc *TrafficController) CheckConnectionLimit(userHash string) error {
	account, err := dao.SelectAccountByHash(userHash)
	if err != nil {
		return err
	}
	
	if account == nil {
		return errors.New("用户不存在")
	}
	
	// 检查IP连接数限制
	if account.IpLimit != nil && *account.IpLimit > 0 {
		tc.connManager.mutex.RLock()
		userConnections := tc.connManager.connections[userHash]
		connCount := len(userConnections)
		tc.connManager.mutex.RUnlock()
		
		if connCount >= int(*account.IpLimit) {
			return fmt.Errorf("连接数超过限制（最大%d个连接）", *account.IpLimit)
		}
	}
	
	return nil
}

// AddConnection 添加连接
func (tc *TrafficController) AddConnection(userHash, connectionId, clientIP, protocol string, conn net.Conn) error {
	// 检查用户配额
	if err := tc.CheckUserQuota(userHash); err != nil {
		return err
	}
	
	// 检查连接数限制
	if err := tc.CheckConnectionLimit(userHash, clientIP); err != nil {
		return err
	}
	
	// 包装连接以进行流量统计
	trafficConn := NewTrafficConn(conn)
	
	// 添加连接到管理器
	tc.connManager.AddConnection(userHash, connectionId, clientIP, protocol, trafficConn)
	
	return nil
}

// RemoveConnection 移除连接并统计流量
func (tc *TrafficController) RemoveConnection(userHash, connectionId string, download, upload uint64) {
	tc.connManager.mutex.Lock()
	defer tc.connManager.mutex.Unlock()
	
	if userConnections, exists := tc.connManager.connections[userHash]; exists {
		if connInfo, exists := userConnections[connectionId]; exists {
			// 如果传入的流量为0，尝试从TrafficConn获取流量统计
			if download == 0 && upload == 0 {
				if trafficConn, ok := connInfo.Conn.(*TrafficConn); ok {
					download = trafficConn.GetDownloadBytes()
					upload = trafficConn.GetUploadBytes()
				}
			}
			
			// 更新流量统计
			connInfo.Download += download
			connInfo.Upload += upload
			
			// 异步更新数据库
			go func() {
				err := dao.NewAccountDao().UpdateAccountFlowByHash(userHash, download, upload)
				if err != nil {
					logrus.Errorf("更新用户流量失败: %v", err)
				}
			}()
			
			// 移除连接
			delete(userConnections, connectionId)
			
			// 如果用户没有其他连接，清理用户记录
			if len(userConnections) == 0 {
				delete(tc.connManager.connections, userHash)
			}
			
			logrus.Infof("移除连接: 用户=%s, 连接ID=%s, 下载=%d, 上传=%d", userHash, connectionId, download, upload)
		}
	}
}

// GetUserConnections 获取用户连接数
func (tc *TrafficController) GetUserConnections(userHash string) int {
	tc.connManager.mutex.RLock()
	defer tc.connManager.mutex.RUnlock()
	
	if userConnections, exists := tc.connManager.connections[userHash]; exists {
		return len(userConnections)
	}
	return 0
}

// GetAllConnections 获取所有连接信息
func (tc *TrafficController) GetAllConnections() map[string]map[string]*ConnectionInfo {
	tc.connManager.mutex.RLock()
	defer tc.connManager.mutex.RUnlock()
	
	// 深拷贝避免并发问题
	result := make(map[string]map[string]*ConnectionInfo)
	for userHash, userConns := range tc.connManager.connections {
		result[userHash] = make(map[string]*ConnectionInfo)
		for connId, connInfo := range userConns {
			// 创建副本
			result[userHash][connId] = &ConnectionInfo{
				UserHash:     connInfo.UserHash,
				ConnectionId: connInfo.ConnectionId,
				ClientIP:     connInfo.ClientIP,
				Protocol:     connInfo.Protocol,
				StartTime:    connInfo.StartTime,
				Download:     connInfo.Download,
				Upload:       connInfo.Upload,
				// 注意：不复制 Conn 字段，避免并发访问
			}
		}
	}
	return result
}

// CleanupExpiredConnections 清理过期连接
func (tc *TrafficController) CleanupExpiredConnections(timeout time.Duration) {
	tc.connManager.mutex.Lock()
	defer tc.connManager.mutex.Unlock()
	
	now := time.Now()
	for userHash, userConnections := range tc.connManager.connections {
		for connectionId, connInfo := range userConnections {
			if now.Sub(connInfo.StartTime) > timeout {
				logrus.Warnf("清理过期连接: 用户=%s, 连接ID=%s, 持续时间=%v", userHash, connectionId, now.Sub(connInfo.StartTime))
				
				// 关闭连接
				if connInfo.Conn != nil {
					connInfo.Conn.Close()
				}
				
				delete(userConnections, connectionId)
			}
		}
		
		// 清理空的用户记录
		if len(userConnections) == 0 {
			delete(tc.connManager.connections, userHash)
		}
	}
}

// StartCleanupTask 启动清理任务
func (tc *TrafficController) StartCleanupTask(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute) // 每5分钟清理一次
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			tc.CleanupExpiredConnections(30 * time.Minute) // 清理30分钟无活动的连接
		}
	}
}