package hysteria2

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"
	"trojan-panel-core/dao"
	"trojan-panel-core/middleware"
	"trojan-panel-core/model/bo"
	"trojan-panel-core/model/constant"
	"trojan-panel-core/model/dto"
	"trojan-panel-core/util"
)

type hysteria2Api struct {
	apiPort         uint
	trafficController *middleware.TrafficController
	dao             *dao.AccountDao
}

func NewHysteria2Api(apiPort uint) *hysteria2Api {
	return &hysteria2Api{
		apiPort:         apiPort,
		trafficController: middleware.GetTrafficController(),
		dao:             dao.NewAccountDao(),
	}
}

// HandleConnection 处理用户连接（添加流量控制）
func (n *hysteria2Api) HandleConnection(hash, clientIP string, conn net.Conn) error {
	// 生成连接ID
	connectionId := util.GenerateUUID()
	
	// 添加连接到流量控制器
	err := n.trafficController.AddConnection(hash, connectionId, clientIP, "hysteria2", conn)
	if err != nil {
		logrus.Errorf("Hysteria2添加连接失败: %v", err)
		return err
	}
	
	// 连接处理逻辑
	go func() {
		defer func() {
			// 连接关闭时统计流量并移除连接
			download, upload := n.getConnectionTraffic(conn)
			n.trafficController.RemoveConnection(hash, connectionId, download, upload)
		}()
		
		// 这里可以添加具体的连接处理逻辑
		// 例如数据转发、协议处理等
	}()
	
	return nil
}

// getConnectionTraffic 获取连接的流量统计
func (n *hysteria2Api) getConnectionTraffic(conn net.Conn) (download, upload uint64) {
	// 如果是TrafficConn，直接获取流量统计
	if trafficConn, ok := conn.(*middleware.TrafficConn); ok {
		return trafficConn.GetDownloadBytes(), trafficConn.GetUploadBytes()
	}
	// 否则返回0（这种情况下会在RemoveConnection中重新获取）
	return 0, 0
}

func apiClient() *http.Client {
	return &http.Client{
		Timeout: 3 * time.Second,
	}
}

func (n *hysteria2Api) ListUsers(clear bool) (map[string]bo.Hysteria2UserTraffic, error) {
	client := apiClient()
	url := fmt.Sprintf("http://127.0.0.1:%d/traffic", n.apiPort)
	if clear {
		url = fmt.Sprintf("%s?clear=1", url)
	}
	resp, err := client.Get(url)
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	if err != nil || resp.StatusCode != http.StatusOK {
		logrus.Errorf("Hysteria2 ListUsers err: %v", err)
		return nil, errors.New(constant.HttpError)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Hysteria2 io read err: %v", err)
		return nil, errors.New(constant.HttpError)
	}
	var users map[string]bo.Hysteria2UserTraffic
	if err = json.Unmarshal(body, &users); err != nil {
		logrus.Errorf("Hysteria2 ListUsers Unmarshal err: %v", err)
		return nil, errors.New(constant.SysError)
	}
	return users, nil
}

func (n *hysteria2Api) GetUser(pass string, clear bool) (*bo.Hysteria2User, error) {
	users, err := n.ListUsers(clear)
	if err != nil {
		return nil, err
	}
	user := users[pass]
	return &bo.Hysteria2User{
		Pass: pass,
		Tx:   user.Tx,
		Rx:   user.Rx,
	}, nil
}
