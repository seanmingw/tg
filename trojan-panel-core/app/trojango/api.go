package trojango

import (
	"context"
	"errors"
	"fmt"
	"github.com/p4gefau1t/trojan-go/api/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"net"
	"time"
	"trojan-panel-core/middleware"
	"trojan-panel-core/model/constant"
	"trojan-panel-core/model/dto"
	"trojan-panel-core/util"
)

type trojanGoApi struct {
	apiPort         uint
	trafficController *middleware.TrafficController
}

func NewTrojanGoApi(apiPort uint) *trojanGoApi {
	return &trojanGoApi{
		apiPort:         apiPort,
		trafficController: middleware.GetTrafficController(),
	}
}

func apiClient(apiPort uint) (clent service.TrojanServerServiceClient, ctx context.Context, clo func(), err error) {
	conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1:%d", apiPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	clent = service.NewTrojanServerServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	clo = func() {
		cancel()
		if conn != nil {
			conn.Close()
		}
	}
	if err != nil {
		logrus.Errorf("trojango apiClient init err: %v", err)
		err = errors.New(constant.GrpcError)
	}
	return
}

// ListUsers query all users on a node
func (t *trojanGoApi) ListUsers() ([]*service.UserStatus, error) {
	client, ctx, clo, err := apiClient(t.apiPort)
	if err != nil {
		return nil, err
	}
	stream, err := client.ListUsers(ctx, &service.ListUsersRequest{})
	if err != nil {
		logrus.Errorf("trojango ListUsers err: %v", err)
		return nil, errors.New(constant.GrpcError)
	}
	defer func() {
		if stream != nil {
			stream.CloseSend()
		}
		clo()
	}()
	var userStatus []*service.UserStatus
	for {
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			logrus.Errorf("trojango ListUsers recv err: %v", err)
			return nil, errors.New(constant.GrpcError)
		}
		if resp != nil {
			userStatus = append(userStatus, resp.Status)
		}
	}
	return userStatus, nil
}

// GetUser query users on a node
func (t *trojanGoApi) GetUser(hash string) (*service.UserStatus, error) {
	client, ctx, clo, err := apiClient(t.apiPort)
	if err != nil {
		return nil, err
	}
	stream, err := client.GetUsers(ctx)
	if err != nil {
		logrus.Errorf("trojango GetUser err: %v", err)
		return nil, errors.New(constant.GrpcError)
	}
	defer func() {
		if stream != nil {
			stream.CloseSend()
		}
		clo()
	}()
	if err = stream.Send(&service.GetUsersRequest{
		User: &service.User{
			Hash: hash,
		},
	}); err != nil {
		logrus.Errorf("trojango GetUser stream send err: %v", err)
		return nil, errors.New(constant.GrpcError)
	}
	resp, err := stream.Recv()
	if resp == nil || err != nil {
		logrus.Errorf("trojango GetUser stream recv err: %v", err)
		return nil, errors.New(constant.GrpcError)
	}
	return resp.Status, nil
}

// set user on node
func (t *trojanGoApi) setUser(setUsersRequest *service.SetUsersRequest) error {
	client, ctx, clo, err := apiClient(t.apiPort)
	if err != nil {
		return err
	}
	stream, err := client.SetUsers(ctx)
	if err != nil {
		logrus.Errorf("trojango setUser err: %v", err)
		return errors.New(constant.GrpcError)
	}
	defer func() {
		if stream == nil {
			stream.CloseSend()
		}
		clo()
	}()
	err = stream.Send(setUsersRequest)
	if err != nil {
		logrus.Errorf("trojango setUser send err: %v", err)
		return errors.New(constant.GrpcError)
	}
	resp, err := stream.Recv()
	if err != nil {
		logrus.Errorf("trojango setUser recv err: %v", err)
		return errors.New(constant.GrpcError)
	}
	if resp != nil && !resp.Success {
		logrus.Errorf("trojango setUser err resp info: %v", resp.Info)
		return errors.New(constant.GrpcError)
	}
	return nil
}

// ReSetUserTrafficByHash reset user traffic
func (t *trojanGoApi) ReSetUserTrafficByHash(hash string) error {
	req := &service.SetUsersRequest{
		Status: &service.UserStatus{
			User: &service.User{
				Hash: hash,
			},
			TrafficTotal: &service.Traffic{
				DownloadTraffic: 0,
				UploadTraffic:   0,
			},
		},
		Operation: service.SetUsersRequest_Modify,
	}
	return t.setUser(req)
}

// SetUserIpLimit set the number of user devices on the node
func (t *trojanGoApi) SetUserIpLimit(hash string, ipLimit uint) error {
	req := &service.SetUsersRequest{
		Status: &service.UserStatus{
			User: &service.User{
				Hash: hash,
			},
			IpLimit: int32(ipLimit),
		},
		Operation: service.SetUsersRequest_Modify,
	}
	return t.setUser(req)
}

// SetUserSpeedLimit set user speed limit on the node
func (t *trojanGoApi) SetUserSpeedLimit(hash string, uploadSpeedLimit int, downloadSpeedLimit int) error {
	req := &service.SetUsersRequest{
		Status: &service.UserStatus{
			User: &service.User{
				Hash: hash,
			},
			SpeedLimit: &service.Speed{
				UploadSpeed:   uint64(uploadSpeedLimit),
				DownloadSpeed: uint64(downloadSpeedLimit),
			},
		},
		Operation: service.SetUsersRequest_Modify,
	}
	return t.setUser(req)
}

// HandleConnection 处理用户连接（添加流量控制）
func (t *trojanGoApi) HandleConnection(hash, clientIP string, conn net.Conn) error {
	// 生成连接ID
	connectionId := util.GenerateUUID()
	
	// 添加连接到流量控制器
	err := t.trafficController.AddConnection(hash, connectionId, clientIP, "trojan-go", conn)
	if err != nil {
		logrus.Errorf("Trojan-Go添加连接失败: %v", err)
		return err
	}
	
	// 连接处理逻辑
	go func() {
		defer func() {
			// 连接关闭时统计流量并移除连接
			download, upload := t.getConnectionTraffic(conn)
			t.trafficController.RemoveConnection(hash, connectionId, download, upload)
		}()
		
		// 这里可以添加具体的连接处理逻辑
		// 例如数据转发、协议处理等
	}()
	
	return nil
}

// getConnectionTraffic 获取连接的流量统计
func (t *trojanGoApi) getConnectionTraffic(conn net.Conn) (download, upload uint64) {
	// 如果是TrafficConn，直接获取流量统计
	if trafficConn, ok := conn.(*middleware.TrafficConn); ok {
		return trafficConn.GetDownloadBytes(), trafficConn.GetUploadBytes()
	}
	// 否则返回0（这种情况下会在RemoveConnection中重新获取）
	return 0, 0
}

// DeleteUser delete user on node
func (t *trojanGoApi) DeleteUser(hash string) error {
	userStatus, err := t.GetUser(hash)
	if err != nil {
		return err
	}
	if userStatus == nil {
		return nil
	}
	req := &service.SetUsersRequest{
		Status: &service.UserStatus{
			User: &service.User{
				Hash: hash,
			},
		},
		Operation: service.SetUsersRequest_Delete,
	}
	return t.setUser(req)
}

// AddUser add user on node
func (t *trojanGoApi) AddUser(dto dto.TrojanGoAddUserDto) error {
	userStatus, err := t.GetUser(dto.Hash)
	if err != nil {
		return err
	}
	if userStatus != nil {
		return nil
	}
	req := &service.SetUsersRequest{
		Status: &service.UserStatus{
			User: &service.User{
				Hash: dto.Hash,
			},
			TrafficTotal: &service.Traffic{
				UploadTraffic:   uint64(dto.UploadTraffic),
				DownloadTraffic: uint64(dto.DownloadTraffic),
			},
			IpLimit: int32(dto.IpLimit),
			SpeedLimit: &service.Speed{
				UploadSpeed:   uint64(dto.UploadSpeedLimit),
				DownloadSpeed: uint64(dto.DownloadSpeedLimit),
			},
		},
		Operation: service.SetUsersRequest_Add,
	}
	return t.setUser(req)
}
