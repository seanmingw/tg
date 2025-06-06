FROM ubuntu:22.04

# 避免交互式提示
ARG DEBIAN_FRONTEND=noninteractive

# 设置时区
RUN apt-get update && apt-get install -y tzdata && \
    ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    dpkg-reconfigure -f noninteractive tzdata

# 安装基本工具
RUN apt-get update && apt-get install -y \
    curl \
    wget \
    git \
    lsof \
    tar \
    unzip \
    vim \
    net-tools \
    ca-certificates \
    gnupg \
    apt-transport-https \
    software-properties-common \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# 安装 Docker
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg && \
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null && \
    apt-get update && apt-get install -y docker-ce docker-ce-cli containerd.io && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# 安装 Docker Compose
RUN curl -L "https://github.com/docker/compose/releases/download/v2.15.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose && \
    chmod +x /usr/local/bin/docker-compose

# 复制项目文件
WORKDIR /app
COPY . /app/

# 设置执行权限
RUN chmod +x /app/cloud_deploy.sh

# 暴露端口
EXPOSE 80 443 8081 8082 8100 8888 9507 6378

# 设置入口点
ENTRYPOINT ["/bin/bash"]
CMD ["-c", "echo 'Trojan Panel Cloud 部署工具已准备就绪。运行 ./cloud_deploy.sh 开始部署。'; exec bash"]