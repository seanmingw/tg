package api

import (
	"github.com/gin-gonic/gin"
	"trojan-panel-core/middleware"
	"trojan-panel-core/model/constant"
	"trojan-panel-core/model/dto"
	"trojan-panel-core/model/vo"
	"trojan-panel-core/service"
)

func Hysteria2Api(c *gin.Context) {
	var hysteria2AuthDto dto.Hysteria2AuthDto
	_ = c.ShouldBindJSON(&hysteria2AuthDto)
	if err := validate.Struct(&hysteria2AuthDto); err != nil {
		vo.Hysteria2ApiFail(constant.ValidateFailed, c)
		return
	}
	//base64DecodeStr, err := base64.StdEncoding.DecodeString()
	//if err != nil {
	//	vo.Hysteria2ApiFail(constant.ValidateFailed, c)
	//	return
	//}
	//pass := string(base64DecodeStr)
	accountHysteria2Vo, err := service.SelectAccountByPass(*hysteria2AuthDto.Auth)
	if err != nil || accountHysteria2Vo == nil {
		vo.Hysteria2ApiFail("", c)
		return
	}
	
	// 获取客户端IP
	clientIP := c.ClientIP()
	
	// 检查连接数限制
	if accountHysteria2Vo.IpLimit > 0 {
		// 获取当前用户连接数
		trafficController := middleware.GetTrafficController()
		userHash := *hysteria2AuthDto.Auth // 在Hysteria2中，我们使用pass作为hash
		currentConnections := trafficController.GetUserConnections(userHash)
		
		// 如果当前连接数已达到限制，拒绝连接
		if currentConnections >= accountHysteria2Vo.IpLimit {
			vo.Hysteria2ApiFail("连接数超过限制", c)
			return
		}
	}
	
	vo.Hysteria2ApiSuccess(*hysteria2AuthDto.Auth, accountHysteria2Vo.IpLimit, accountHysteria2Vo.DownloadSpeedLimit, accountHysteria2Vo.UploadSpeedLimit, c)
}
