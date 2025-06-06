package hysteria

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"
	"trojan-panel-core/dao"
	"trojan-panel-core/middleware"
	"trojan-panel-core/model/constant"
	"trojan-panel-core/model/dto"
	"trojan-panel-core/util"
)

type hysteriaApi struct {
	apiPort         uint
	trafficController *middleware.TrafficController
	dao             *dao.AccountDao
}

func NewHysteriaApi(apiPort uint) *hysteriaApi {
	return &hysteriaApi{
		apiPort:         apiPort,
		trafficController: middleware.GetTrafficController(),
		dao:             dao.NewAccountDao(),
	}
}

// HandleConnection 处理用户连接（添加流量控制）
func (h *hysteriaApi) HandleConnection(hash, clientIP string, conn net.Conn) error {
	// 生成连接ID
	connectionId := util.GenerateUUID()
	
	// 添加连接到流量控制器
	err := h.trafficController.AddConnection(hash, connectionId, clientIP, "hysteria", conn)
	if err != nil {
		logrus.Errorf("Hysteria添加连接失败: %v", err)
		return err
	}
	
	// 连接处理逻辑
	go func() {
		defer func() {
			// 连接关闭时统计流量并移除连接
			download, upload := h.getConnectionTraffic(conn)
			h.trafficController.RemoveConnection(hash, connectionId, download, upload)
		}()
		
		// 这里可以添加具体的连接处理逻辑
		// 例如数据转发、协议处理等
	}()
	
	return nil
}

// getConnectionTraffic 获取连接的流量统计
func (h *hysteriaApi) getConnectionTraffic(conn net.Conn) (download, upload uint64) {
	// 如果是TrafficConn，直接获取流量统计
	if trafficConn, ok := conn.(*middleware.TrafficConn); ok {
		return trafficConn.GetDownloadBytes(), trafficConn.GetUploadBytes()
	}
	// 否则返回0（这种情况下会在RemoveConnection中重新获取）
	return 0, 0
}

// 新增流量统计方法：定期上报用户流量到数据库
func (h *hysteriaApi) ReportUserTraffic(hash string, download, upload uint64) error {
    // 调用DAO层更新account表的download和upload字段
    return h.dao.UpdateAccountTraffic(hash, download, upload)
}

// 在连接断开或定时任务中调用ReportUserTraffic
func (h *hysteriaApi) handleConnection(conn net.Conn, userHash string) {
    // ... existing 连接处理逻辑 ...
    
    // 连接关闭时统计流量
    defer func() {
        totalDownload := getDownloadBytes(conn)
        totalUpload := getUploadBytes(conn)
        _ = h.ReportUserTraffic(userHash, totalDownload, totalUpload)
    }()
}