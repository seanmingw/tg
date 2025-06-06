package middleware

import (
	"net"
	"sync/atomic"
	"time"
)

// TrafficConn 流量统计连接包装器
type TrafficConn struct {
	net.Conn
	downloadBytes uint64
	uploadBytes   uint64
	startTime     time.Time
}

// NewTrafficConn 创建流量统计连接
func NewTrafficConn(conn net.Conn) *TrafficConn {
	return &TrafficConn{
		Conn:      conn,
		startTime: time.Now(),
	}
}

// Read 读取数据并统计下载流量
func (tc *TrafficConn) Read(b []byte) (n int, err error) {
	n, err = tc.Conn.Read(b)
	if n > 0 {
		atomic.AddUint64(&tc.downloadBytes, uint64(n))
	}
	return
}

// Write 写入数据并统计上传流量
func (tc *TrafficConn) Write(b []byte) (n int, err error) {
	n, err = tc.Conn.Write(b)
	if n > 0 {
		atomic.AddUint64(&tc.uploadBytes, uint64(n))
	}
	return
}

// GetDownloadBytes 获取下载字节数
func (tc *TrafficConn) GetDownloadBytes() uint64 {
	return atomic.LoadUint64(&tc.downloadBytes)
}

// GetUploadBytes 获取上传字节数
func (tc *TrafficConn) GetUploadBytes() uint64 {
	return atomic.LoadUint64(&tc.uploadBytes)
}

// GetTotalBytes 获取总流量
func (tc *TrafficConn) GetTotalBytes() uint64 {
	return tc.GetDownloadBytes() + tc.GetUploadBytes()
}

// GetDuration 获取连接持续时间
func (tc *TrafficConn) GetDuration() time.Duration {
	return time.Since(tc.startTime)
}

// GetTrafficStats 获取流量统计信息
func (tc *TrafficConn) GetTrafficStats() (download, upload uint64, duration time.Duration) {
	return tc.GetDownloadBytes(), tc.GetUploadBytes(), tc.GetDuration()
}