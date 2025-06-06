#!/usr/bin/env bash

# TG Panel 云端部署脚本
# 用于在云服务器上快速部署 TG Panel

set -e

# 加载环境变量
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
else
  echo "错误：找不到 .env 文件，请先创建并配置环境变量"
  exit 1
fi

# 设置默认值
TP_DATA=${TP_DATA:-/tpdata}
DOMAIN=${DOMAIN:-example.com}
MARIADB_PASSWORD=${MARIADB_PASSWORD:-changeme}
REDIS_PASSWORD=${REDIS_PASSWORD:-changeme}

echo "==========================================================="
echo "              TG Panel 云端部署脚本                   "
echo "==========================================================="
echo ""
echo "数据目录: $TP_DATA"
echo "域名: $DOMAIN"
echo ""

# 检查 Docker 和 Docker Compose 是否已安装
if ! command -v docker &> /dev/null; then
  echo "错误：Docker 未安装，请先安装 Docker"
  exit 1
fi

if ! command -v docker-compose &> /dev/null; then
  echo "错误：Docker Compose 未安装，请先安装 Docker Compose"
  exit 1
fi

# 创建必要的目录结构
echo "正在创建目录结构..."
mkdir -p ${TP_DATA}/{caddy,cert,web,mariadb,redis,trojan-panel,trojan-panel-ui,trojan-panel-core}
mkdir -p ${TP_DATA}/caddy/logs
mkdir -p ${TP_DATA}/trojan-panel/{logs,config}
mkdir -p ${TP_DATA}/trojan-panel-ui/nginx
mkdir -p ${TP_DATA}/trojan-panel-core/{logs,config}
mkdir -p ${TP_DATA}/trojan-panel-core/bin/{xray,trojango,hysteria,naiveproxy,hysteria2}/config

# 创建 Caddy 配置文件
echo "正在创建 Caddy 配置文件..."
cat > ${TP_DATA}/caddy/config.json << EOF
{
  "admin": {
    "disabled": true
  },
  "logging": {
    "logs": {
      "default": {
        "writer": {
          "output": "file",
          "filename": "/tpdata/caddy/logs/caddy.log"
        },
        "level": "ERROR"
      }
    }
  },
  "apps": {
    "http": {
      "servers": {
        "srv0": {
          "listen": [":80"],
          "routes": [
            {
              "handle": [
                {
                  "handler": "static_response",
                  "headers": {
                    "Location": ["https://${DOMAIN}"]
                  },
                  "status_code": 301
                }
              ]
            }
          ]
        },
        "srv1": {
          "listen": [":443"],
          "routes": [
            {
              "handle": [
                {
                  "handler": "file_server",
                  "root": "/tpdata/web"
                }
              ]
            }
          ],
          "tls_connection_policies": [
            {
              "match": {
                "sni": ["${DOMAIN}"]
              }
            }
          ]
        }
      }
    },
    "tls": {
      "automation": {
        "policies": [
          {
            "subjects": ["${DOMAIN}"],
            "issuer": {
              "module": "acme"
            }
          }
        ]
      }
    }
  }
}
EOF

# 创建 Trojan Panel UI Nginx 配置文件
echo "正在创建 Trojan Panel UI Nginx 配置文件..."
cat > ${TP_DATA}/trojan-panel-ui/nginx/default.conf << EOF
server {
    listen 8888;
    
    root /usr/share/nginx/html;
    index index.html index.htm;
    
    location / {
        try_files \$uri \$uri/ /index.html;
    }
    
    location /api {
        proxy_pass http://trojan-panel:8081;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
    }
}
EOF

# 启动服务
echo "正在启动服务..."
docker-compose up -d

echo ""
echo "==========================================================="
echo "              TG Panel 部署完成                       "
echo "==========================================================="
echo ""
echo "管理面板: https://${DOMAIN}:8888"
echo "默认用户名: sysadmin"
echo "默认密码: 123456"
echo ""
echo "请务必修改默认密码以确保安全！"
echo ""