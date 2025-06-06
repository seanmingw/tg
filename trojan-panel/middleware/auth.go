// ... existing code ...

// 新增设备数检查逻辑
func checkDeviceLimit(c *gin.Context, username string) error {
    // 查询当前用户已连接的设备数
    currentDevices, err := service.GetUserConnectedDevices(username)
    if err != nil {
        return err
    }
    
    // 查询用户的ip_limit（设备数限制）
    user, err := service.GetAccountByUsername(username)
    if err != nil {
        return err
    }
    
    if currentDevices >= int(user.IpLimit) {
        return fmt.Errorf("设备数超过限制（最大%d台）", user.IpLimit)
    }
    return nil
}

// 在认证中间件中添加设备限制检查
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // ... existing 认证逻辑 ...
        
        // 认证通过后检查设备限制
        if err := checkDeviceLimit(c, username); err != nil {
            vo.Fail(constant.DeviceLimitExceeded, c)
            c.Abort()
            return
        }
        
        c.Next()
    }
}

// ... existing code ...