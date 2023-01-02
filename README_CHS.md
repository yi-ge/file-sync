# `file-sync` 是一个文件同步命令

[![license](https://img.shields.io/github/license/yi-ge/file-sync.svg?style=flat-square)](https://github.com/yi-ge/file-sync/blob/master/LICENSE)
[![GitHub Actions](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Fyi-ge%2Ffile-sync%2Fbadge%3Fref%3Dmain&style=flat-square)](https://actions-badge.atrox.dev/yi-ge/file-sync/goto?ref=main)
[![npm version](https://img.shields.io/npm/v/file-sync-cli/latest?style=flat-square)](https://www.npmjs.com/package/file-sync-cli)
[![GitHub last commit](https://img.shields.io/github/last-commit/yi-ge/file-sync.svg?style=flat-square)](https://github.com/yi-ge/file-sync)
<!-- [![Test Results](https://gist.github.com/yi-ge/00fdcacb47689d14b8e9fdf7fb0f7288/raw/badge.svg)](https://github.com/yi-ge/file-sync)
[![Coveralls github](https://img.shields.io/coveralls/github/yi-ge/file-sync?style=flat-square)](https://coveralls.io/github/yi-ge/file-sync?branch=main) -->

[ENGLISH](README.md)

⚠️⛔️ 此程序尚在开发！！！ ⛔️ ⚠️

自动同步单个文件。为单个用户同步 `.env` 或 `config` 文件。

`file-sync`是跨平台的，非常适合在你的多个设备中自动同步`.env`、`.bashrc`、`.profile`、`.zshrc`、`.ssh/config`之类的文件。

`file-sync`设计原则：服务器是不可信的，客户端（本地）是可信的。

## 特性

- 安全 - 同步数据时，你的数据经过加密后进行传输
- 灵活 - 支持跨平台工作。在不同系统\不同路径\不同文件名\不同文件权限的情况下同步同一个文件的内容
- 自动化 - 文件被自动同步，无需手工干预，而且你不会丢失任何数据

## 安装

```bash
go install https://github.com/yi-ge/file-sync
```

## 使用

### 设置服务器URL

```bash
file-sync --config https://example.com
```

默认使用命令作者提供的免费公共服务器：<https://api.yizcore.xyz>

### 显示系统信息

```bash
file-sync --info
```

### 登录/注册设备

使用您的邮箱登录，将自动注册当前设备。

```bash
file-sync --login email@example.com
```

### 列出设备

```bash
file-sync --list-device
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

### 设置为开机自启服务

```bash
file-sync service --enable
```

### 禁用开机自启服务

```bash
file-sync service --disable
```

### 启动同步服务

```bash
file-sync service --start
```

### 关闭同步服务

```bash
file-sync service --stop
```

### 移除文件同步项

```bash
file-sync remove <file id>
```

注意，这将会移除该文件所有的同步项目。

### 移除单个设备的文件同步项

```bash
file-sync remove <file id> <number>
```

或

```bash
file-sync remove <file id> <device id>
```

提示：所有的`<device id>`均可用简写。

## 使用自托管服务器

可以选择借助Docker或者自行搭建PHP运行环境在自己的服务器中部署。

服务器端采用标准的HTTP API方式进行交互，没有用到罕见模块，极具兼容性，因此可以将程序部署在绝大部分虚拟主机中。

默认提供的PHP代码需要搭配MySQL 5.4+数据库使用。

### Docker

```bash
docker run xx:file-sync-server
```

### PHP

require PHP >= v5.4 （64位）

#### Server Configuration

<details><summary>CLICK ME</summary>
<p>

##### Apache

You may need to add the following snippet in your Apache HTTP server virtual host configuration or **.htaccess** file.

```apacheconf
RewriteEngine on
RewriteCond %{REQUEST_FILENAME} !-f
RewriteCond %{REQUEST_FILENAME} !-d
RewriteCond $1 !^(index\.php)
RewriteRule ^(.*)$ /index.php/$1 [L]
```

Alternatively, if you’re lucky enough to be using a version of Apache greater than 2.2.15, then you can instead just use this one, single line:

```apacheconf
FallbackResource /index.php
```

##### IIS

For IIS you will need to install URL Rewrite for IIS and then add the following rule to your `web.config`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<configuration>
    <system.webServer>
        <rewrite>
          <rule name="Toro" stopProcessing="true">
            <match url="^(.*)$" ignoreCase="false" />
              <conditions logicalGrouping="MatchAll">
                <add input="{REQUEST_FILENAME}" matchType="IsFile" ignoreCase="false" negate="true" />
                <add input="{REQUEST_FILENAME}" matchType="IsDirectory" ignoreCase="false" negate="true" />
                <add input="{R:1}" pattern="^(index\.php)" ignoreCase="false" negate="true" />
              </conditions>
            <action type="Rewrite" url="/index.php/{R:1}" />
          </rule>
        </rewrite>
    </system.webServer>
</configuration>
```

##### Nginx

Under the `server` block of your virtual host configuration, you only need to add three lines.

```conf
location / {
  try_files $uri $uri/ /index.php?$args;
}
```

</p>
</details>

## 服务器端API

`file-sync`程序目前使用`HTTP API`完成同步交互。目前~~已完成~~PHP版本的服务器端API。

<details><summary>CLICK ME</summary>
<p>

### 基础数据结构

```text
user: email, emailSha1, verify, publicKey, privateKey, createdAt
device: email, machineId, machineName, machineKey, createdAt
config: email, machineId, fileId, path, attribute, createdAt
file: email, fileId, fileName, content, sha256, fromMachineId, updateAt
log: email, machineId, action, content, createdAt
```

### 登录/注册用户并注册设备

POST /device/add

```json
{
  "email": "sha1(email)",
  "machineId": "sha1(sha256(machineId))",
  "machineName": "verify加密的machineName",
  "verify": "sha1(密码的sha256中的前16个字符取sha1)取后32个字符",
  "publicKey": "新生成的publicKey",
  "privateKey": "密码的sha256中的第二段16个字符进行加密的私钥（私钥密码是第三段16个字符）"
}
```

Return：

```json
{
  "status": 1, // 新用户及新设备注册成功， 2： 老用户新设备注册成功
  // -2：设备已注册， -3：此用户已经存在但verify值验证失败
  "result": {
    "publicKey": "verify作为密码 时间戳@公钥 加密", // 如果解密后的内容中的publicKey和传输的publicKey相同，则说明该用户是新用户。
    "privateKey": "verify作为密码 时间戳@私钥 加密 - 密码的sha256中的第二段16个字符进行加密的私钥（私钥密码是第三段16个字符）", // 如果不相同，说明该用户是老用户，则需要以返回回来的publicKey和privateKey为准。
    "machineKey": "由服务器端生成的machineKey"
  }
}
```

### 获取设备列表

POST /device/list

```json
{
  "timestamp": "时间戳",
  "machineId": "sha1(sha256(machineId))",
  "token": "签名[所有字段按json的key的ASCII字符顺序进行升序排列]",
  "email": "sha1(email)",
}
```

Return：

```json
{
  "status": 1,
  "result": [
    {
      "machineId": "sha1(sha256(machineId))",
      "machineName": "verify加密的machineName"
    }
  ]
}
```

### 移除设备

POST /device/remove

```json
{
  "timestamp": "时间戳",
  "machineId": "sha1(sha256(machineId))",
  "token": "签名[所有字段按json的key的ASCII字符顺序进行升序排列]",
  "email": "sha1(email)",
  "machineKey": "sha1(由服务器端生成的machineKey)",
  "removeMachineId": "sha1(sha256(machineId))",
}
```

Return：

```json
{
  "status": 1
}
```

### 获取文件同步配置信息

POST /file/configs

```json
{
  "timestamp": "时间戳",
  "machineId": "sha1(sha256(machineId))",
  "token": "签名[所有字段按json的key的ASCII字符顺序进行升序排列]",
  "email": "sha1(email)"
}
```

Return：

```json
{
  "status": 1,
  "result": [
    {
      "fileId": "",
      "fileName": "私钥加密后的fileName",
      "config": [
        {
          "machineId":, "machineId",
          "path": "私钥加密后的path"
        }
      ]
    }
  ]
}
```

### 设置文件配置

POST /file/config

```json
{
  "timestamp": "时间戳",
  "machineId": "sha1(sha256(machineId))",
  "token": "签名[所有字段按json的key的ASCII字符顺序进行升序排列]",
  "email": "sha1(email)",
  "config": {
    "fileId": "fileId",
    "fileName": "私钥加密后的fileName，如果是添加新同步则为空",
    "action": "add/remove",
    "machineId": "sha1(sha256(machineId))",
    "path": "私钥加密后的path"
  }
}
```

Return：

```json
{
  "status": 1
}
```

### 检查单个文件是否存在更新

POST /file/check

```json
{
  "timestamp": "时间戳",
  "machineId": "sha1(sha256(machineId))",
  "token": "签名[所有字段按json的key的ASCII字符顺序进行升序排列]",
  "email": "sha1(email)",
  "fileId": "fileId",
  "sha256": "文件sha256",
  "updateAt": "文件最后一次编辑时间"
}
```

Return：

```json
{
  "status": 0 // 0：文件不存在， 1：存在更新， 2：无更新
}
```

### 下载/上载文件数据

POST /file/sync

不带"content"和"updateAt"为下载，否则是上载。

```json
{
  "timestamp": "时间戳",
  "machineId": "sha1(sha256(machineId))",
  "token": "签名[所有字段按json的key的ASCII字符顺序进行升序排列]",
  "email": "sha1(email)",
  "fileId": "fileId，由首次添加的文件sha256取sha1生成",
  "sha256": "文件sha256",
  "updateAt": "文件最后一次编辑时间",
  "content": "私钥加密的文件内容"
}
```

Return：

```json
{
  "status": 1
}
```

</p>
</details>

## 开发

<details><summary>CLICK ME</summary>
<p>

### 进入开发调试环境

需要手工设置环境变量`GO_ENV`为`development`。

例如`Windows`平台中`PowerShell`环境下：

安装`xampp`，vscode配置：`"php.debug.executablePath": "C:\\xampp\\php\\php.exe"`。

```bash
$Env:GO_ENV = 'development'
```

例如`*unix`：

```bash
export GO_ENV="development"
```

</p>
</details>

## 关于安全性

按照`服务器是不可信的`的设计原则，服务器中存储的数据均为加密数据。由于`file-sync`使用了非对称加密且在外网使用的时候配合SSL，在传输过程中也是安全的。

服务器中以加密的形式保存所有版本的文件。

如果在编辑文件的时候处于离线状态，`file-sync`会自动记录最后完成编辑的时间，待接入网络后与服务器中的版本进行同步。在此期间，如果在另外一台设备中编辑了同一个文件且编辑内容不一致，将会不可避免的产生冲突，`file-sync`不会自动合并冲突，但是会自动应用最后一个修改的版本，产生冲突的版本将会以`[文件名].[日期].backup`的方式存储到所有同步设备同一目录下。

## TODO

- 支持生成guest账号，在可信设备和不可信设备之间单向/双向同步文件。
