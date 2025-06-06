package vo

type AccountHysteriaVo struct {
	Id                 uint   `json:"id"`
	Username           string `json:"username"`
	IpLimit            int    `json:"ip_limit"`
	DownloadSpeedLimit int    `json:"download_speed_limit"`
	UploadSpeedLimit   int    `json:"upload_speed_limit"`
}

type AccountVo struct {
	Id       uint     `json:"id"`
	Username string   `json:"username"`
	RoleId   uint     `json:"roleId"`
	Deleted  uint     `json:"deleted"`
	Roles    []string `json:"roles"`
}
