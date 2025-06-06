package vo

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// response object
type hysteriaResult struct {
	Ok                 bool   `json:"ok"`
	Msg                string `json:"msg"`
	IpLimit            int    `json:"ip_limit"`
	DownloadSpeedLimit int    `json:"download_speed_limit"`
	UploadSpeedLimit   int    `json:"upload_speed_limit"`
}

func HysteriaApiSuccess(msg string, ipLimit int, downloadSpeedLimit int, uploadSpeedLimit int, c *gin.Context) {
	c.JSON(http.StatusOK, hysteriaResult{
		Ok:                 true,
		Msg:                msg,
		IpLimit:            ipLimit,
		DownloadSpeedLimit: downloadSpeedLimit,
		UploadSpeedLimit:   uploadSpeedLimit,
	})
}

func HysteriaApiFail(msg string, c *gin.Context) {
	c.JSON(http.StatusOK, hysteriaResult{
		Ok:  false,
		Msg: msg,
	})
}
