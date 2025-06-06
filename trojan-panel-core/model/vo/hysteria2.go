package vo

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// response object
type hysteria2Result struct {
	Ok                 bool   `json:"ok"`
	Id                 string `json:"id"`
	IpLimit            int    `json:"ip_limit"`
	DownloadSpeedLimit int    `json:"download_speed_limit"`
	UploadSpeedLimit   int    `json:"upload_speed_limit"`
}

func Hysteria2ApiSuccess(id string, ipLimit int, downloadSpeedLimit int, uploadSpeedLimit int, c *gin.Context) {
	c.JSON(http.StatusOK, hysteria2Result{
		Ok:                 true,
		Id:                 id,
		IpLimit:            ipLimit,
		DownloadSpeedLimit: downloadSpeedLimit,
		UploadSpeedLimit:   uploadSpeedLimit,
	})
}

func Hysteria2ApiFail(id string, c *gin.Context) {
	c.JSON(http.StatusOK, hysteria2Result{
		Ok: false,
		Id: id,
	})
}
