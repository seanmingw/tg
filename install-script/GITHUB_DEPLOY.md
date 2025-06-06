# GitHub 部署指南

本文档提供了如何将 Trojan Panel 云端部署版本推送到 GitHub 并使用 GitHub Actions 进行自动构建的步骤。

## 准备工作

1. 在 GitHub 上创建一个新的仓库
2. 确保你已经在本地安装了 Git

## 推送到 GitHub

### 1. 配置 Git 用户信息

如果你还没有配置 Git 用户信息，请运行以下命令：

```bash
git config --global user.name "你的名字"
git config --global user.email "你的邮箱"
```

### 2. 添加远程仓库

```bash
git remote add origin https://github.com/你的用户名/你的仓库名.git
```

### 3. 推送代码到 GitHub

```bash
git push -u origin main
```

## 配置 GitHub Actions

代码推送到 GitHub 后，GitHub Actions 将自动运行。但是，要使 Docker 镜像构建和推送功能正常工作，你需要配置以下 Secrets：

1. 在 GitHub 仓库页面，点击 "Settings" 标签
2. 在左侧菜单中，点击 "Secrets and variables" -> "Actions"
3. 点击 "New repository secret" 按钮
4. 添加以下 Secrets：
   - 名称：`DOCKERHUB_USERNAME`，值：你的 Docker Hub 用户名
   - 名称：`DOCKERHUB_TOKEN`，值：你的 Docker Hub 访问令牌（可以在 Docker Hub 账户设置中生成）

## 使用 Docker 镜像

一旦 GitHub Actions 工作流程成功运行，你可以使用以下命令拉取和运行 Docker 镜像：

```bash
# 拉取镜像
docker pull 你的用户名/trojan-panel-cloud:latest

# 运行容器
docker run -it --rm \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /你想存储数据的路径:/tpdata \
  -p 80:80 -p 443:443 -p 8888:8888 \
  你的用户名/trojan-panel-cloud:latest
```

进入容器后，你可以运行 `./cloud_deploy.sh` 来部署 Trojan Panel。

## 注意事项

1. 确保你的服务器已经安装了 Docker 和 Docker Compose
2. 在运行容器时，需要将宿主机的 Docker 套接字挂载到容器内，以便容器内的 Docker 命令可以控制宿主机的 Docker 守护进程
3. 根据需要修改端口映射，确保所需的端口都已正确映射
4. 在生产环境中，建议使用 HTTPS 并配置适当的安全措施

## 故障排除

如果遇到问题，请检查：

1. GitHub Actions 日志中是否有错误信息
2. Docker Hub 中是否成功创建了镜像
3. 服务器上的 Docker 和 Docker Compose 是否正确安装
4. 防火墙设置是否允许所需的端口通信