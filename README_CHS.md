# `file-sync` 是一个文件同步命令

[![license](https://img.shields.io/github/license/yi-ge/file-sync.svg?style=flat-square)](https://github.com/yi-ge/file-sync/blob/master/LICENSE)
[![GitHub Actions](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Fyi-ge%2Ffile-sync%2Fbadge%3Fref%3Dmain&style=flat-square)](https://actions-badge.atrox.dev/yi-ge/file-sync/goto?ref=main)
[![Test Results](https://gist.github.com/yi-ge/00fdcacb47689d14b8e9fdf7fb0f7288/raw/badge.svg)](https://github.com/yi-ge/file-sync)
[![Coveralls github](https://img.shields.io/coveralls/github/yi-ge/file-sync?style=flat-square)](https://coveralls.io/github/yi-ge/file-sync?branch=main)
[![GitHub last commit](https://img.shields.io/github/last-commit/yi-ge/file-sync.svg?style=flat-square)](https://github.com/yi-ge/file-sync)

[ENGLISH](README.md)

自动同步单个文件。为单个用户同步 `.env` 或 `.config` 文件。

`file-sync`设计原则：服务器是不可信的，客户端（本地）是可信的。

## 安装

## 使用

### 设置服务器URL

```bash
file-sync --config https://example.com
```

默认使用命令作者提供的免费公共服务器：<https://file-sync.openapi.site/>

### 登录/注册设备

使用您的邮箱登录，将自动注册当前设备。

```bash
file-sync --login email@example.com
```

### 列出设备

```bash
file-sync list --device
```

### 移除当前设备

```bash
file-sync --remove-device
```

移除该设备后，将会删除该设备中所有的配置信息（但不包括用户文件）。此设备如需再次同步文件，需要重新登录进行设备注册。

### 移除指定设备

```bash
file-sync --remove-device <device id>
```

移除该设备后，该设备下次同步时将会删除该设备中所有的配置信息（但不包括用户文件）。此设备如需再次同步文件，需要重新登录进行设备注册。

### 列出文件列表

```bash
file-sync list
```

### 添加文件同步项

```bash
file-sync add <file path>
```

### 添加已经存在的文件同步项

```bash
file-sync add <file id> <file path>
```

### 移除文件同步项

```bash
file-sync remove <file id>
```
