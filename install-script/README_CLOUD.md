# Trojan Panel 云端部署指南

## 简介

本指南提供了如何在云服务器上部署 Trojan Panel 的详细步骤。Trojan Panel 是一个支持 Xray/Trojan-Go/Hysteria/NaiveProxy 的多用户 Web 管理面板。

## 环境要求

- 一台运行 Linux 的云服务器（推荐 Ubuntu 20.04 LTS 或更高版本）
- 已安装 Docker 和 Docker Compose
- 一个已解析到服务器 IP 的域名（用于 HTTPS 证书）

## 快速部署

### 1. 克隆仓库

```bash
git clone https://github.com/yourusername/trojan-panel-cloud.git
cd trojan-panel-cloud
```

### 2. 配置环境变量

编辑 `.env` 文件，根据您的需求修改配置：

```bash
# 基本配置 - 必须修改
TP_DATA=/tpdata  # 数据存储路径，可以根据需要修改
DOMAIN=your-domain.com  # 替换为您的域名

# 安全配置 - 强烈建议修改
MARIADB_PASSWORD=your-secure-password  # 数据库密码
REDIS_PASSWORD=your-secure-password  # Redis密码
```

### 3. 创建必要的目录

```bash
mkdir -p ${TP_DATA}/{caddy,cert,web,mariadb,redis,trojan-panel,trojan-panel-ui,trojan-panel-core}
mkdir -p ${TP_DATA}/trojan-panel/{logs,config}
mkdir -p ${TP_DATA}/trojan-panel-ui/nginx
mkdir -p ${TP_DATA}/trojan-panel-core/{logs,config}
mkdir -p ${TP_DATA}/trojan-panel-core/bin/{xray,trojango,hysteria,naiveproxy,hysteria2}/config
```

### 4. 启动服务

```bash
docker-compose up -d
```

### 5. 访问面板

部署完成后，您可以通过以下地址访问 Trojan Panel：

- 管理面板：`https://your-domain.com:8888`
- 默认用户名：`sysadmin`
- 默认密码：`123456`

## 配置说明

### 端口配置

默认情况下，以下端口将被使用：

- 80: HTTP（Caddy）
- 443: HTTPS（Caddy）
- 8081: Trojan Panel 后端 API
- 8082: Trojan Panel Core API
- 8100: gRPC 端口
- 8888: Trojan Panel UI
- 9507: MariaDB
- 6378: Redis

如需修改端口，请编辑 `.env` 文件中的相应配置。

### 证书配置

默认情况下，系统会使用 Caddy 自动申请 Let's Encrypt 证书。如果您需要使用自己的证书，请将证书文件放置在 `${TP_DATA}/cert/` 目录下，并命名为 `${DOMAIN}.crt` 和 `${DOMAIN}.key`。

## 数据备份

所有数据都存储在 `${TP_DATA}` 目录中。定期备份此目录可以保护您的数据安全：

```bash
tar -czvf trojan-panel-backup-$(date +%Y%m%d).tar.gz ${TP_DATA}
```

## 故障排除

### 查看日志

```bash
# 查看所有容器状态
docker-compose ps

# 查看特定服务的日志
docker-compose logs trojan-panel
docker-compose logs trojan-panel-core
```

### 常见问题

1. **无法访问面板**：检查防火墙设置，确保相关端口已开放。
2. **证书问题**：确保域名正确解析到服务器 IP，且 80/443 端口可以被外部访问。
3. **数据库连接失败**：检查 MariaDB 容器是否正常运行，以及密码配置是否正确。

## 更新

```bash
# 拉取最新代码
git pull

# 重新构建并启动容器
docker-compose down
docker-compose up -d
```

## 安全建议

1. 修改默认的管理员密码
2. 使用强密码保护数据库和 Redis
3. 配置防火墙，只开放必要的端口
4. 定期更新系统和 Docker 镜像

## 贡献

欢迎提交 Issue 和 Pull Request 来帮助改进此项目。

## 许可证

本项目遵循 MIT 许可证。