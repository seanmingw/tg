package api

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"trojan-panel-core/middleware"
	"trojan-panel-core/model/constant"
	"trojan-panel-core/model/dto"
	"trojan-panel-core/model/vo"
	"trojan-panel-core/service"
)

func HysteriaApi(c *gin.Context) {
	var hysteriaAuthDto dto.HysteriaAuthDto
	_ = c.ShouldBindJSON(&hysteriaAuthDto)
	if err := validate.Struct(&hysteriaAuthDto); err != nil {
		vo.HysteriaApiFail(constant.ValidateFailed, c)
		return
	}
	base64DecodeStr, err := base64.StdEncoding.DecodeString(*hysteriaAuthDto.Payload)
	if err != nil {
		vo.HysteriaApiFail(constant.ValidateFailed, c)
		return
	}
	pass := string(base64DecodeStr)
	accountHysteriaVo, err := service.SelectAccountByPass(pass)
	if err != nil || accountHysteriaVo == nil {
		vo.HysteriaApiFail(constant.UsernameOrPassError, c)
		return
	}
	
	// 获取客户端IP
	clientIP := c.ClientIP()
	
	// 检查连接数限制
	if accountHysteriaVo.IpLimit > 0 {
		// 获取当前用户连接数
		trafficController := middleware.GetTrafficController()
		userHash := pass // 在Hysteria中，我们使用pass作为hash
		currentConnections := trafficController.GetUserConnections(userHash)
		
		// 如果当前连接数已达到限制，拒绝连接
		if currentConnections >= accountHysteriaVo.IpLimit {
			vo.HysteriaApiFail("连接数超过限制", c)
			return
		}
	}
	
	vo.HysteriaApiSuccess("success", accountHysteriaVo.IpLimit, accountHysteriaVo.DownloadSpeedLimit, accountHysteriaVo.UploadSpeedLimit, c)
}
