# Trojan Panel 安装指南

## 快速安装

使用以下命令安装Trojan Panel：

```bash
bash -c "$(curl -fsSL https://cdn.jsdelivr.net/gh/seanmingw/tg@main/install-script/install_script.sh)"
```

## 查看帮助

```bash
bash -c "$(curl -fsSL https://cdn.jsdelivr.net/gh/seanmingw/tg@main/install-script/install_script.sh)" -- -h
```

## 查看版本

```bash
bash -c "$(curl -fsSL https://cdn.jsdelivr.net/gh/seanmingw/tg@main/install-script/install_script.sh)" -- -v
```

## 注意事项

- 安装脚本使用jsDelivr CDN，解决了无法访问raw.githubusercontent.com的问题
- 静态资源文件也使用jsDelivr CDN加速

