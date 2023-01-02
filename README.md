# `file-sync` is a file sync command

[![license](https://img.shields.io/github/license/yi-ge/file-sync.svg?style=flat-square)](https://github.com/yi-ge/file-sync/blob/master/LICENSE)
[![GitHub Actions](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Fyi-ge%2Ffile-sync%2Fbadge%3Fref%3Dmain&style=flat-square)](https://actions-badge.atrox.dev/yi-ge/file-sync/goto?ref=main)
[![npm version](https://img.shields.io/npm/v/file-sync-cli/latest?style=flat-square)](https://www.npmjs.com/package/file-sync-cli)
[![GitHub last commit](https://img.shields.io/github/last-commit/yi-ge/file-sync.svg?style=flat-square)](https://github.com/yi-ge/file-sync)
<!-- [![Test Results](https://gist.github.com/yi-ge/00fdcacb47689d14b8e9fdf7fb0f7288/raw/badge.svg)](https://github.com/yi-ge/file-sync)
[![Coveralls github](https://img.shields.io/coveralls/github/yi-ge/file-sync?style=flat-square)](https://coveralls.io/github/yi-ge/file-sync?branch=main) -->

[简体中文](README_CHS.md)

⚠️⛔️ This program is not yet under development !!! ⛔️ ⚠️

Automatically sync single file. Sync `.env` file or `.config` file for single user.

The design principle of `file-sync`: The server is untrustworthy, the client (local) is trusted.

## Feature

- Secure - when synchronizing data, your data is encrypted for transmission
- Flexible - supports cross-platform work. Synchronize the contents of the same file on different systems, different paths, different file names and different file permissions
- Automated - files are synchronized automatically without manual intervention, and you don't lose any data

## Install

```bash
go install https://github.com/yi-ge/file-sync
```

## Usage

### Set the server URL

```bash
file-sync --config https://example.com
```

The free public server provided by the command author is used by default: <https://api.yizcore.xyz>

### Login / Register Device

Signing in with your email will automatically register your current device.

```bash
file-sync --login email@example.com
```

### Display system information

```bash
file-sync --info
```

### List devices

```bash
file-sync --list-device
```

### Remove the current device

```bash
file-sync --remove-device
```

Removing this device will delete all configuration information (but not user files) from this device. If this device needs to synchronize files again, you need to log in again for device registration.

### Remove the specified device

```bash
file-sync --remove-device <device id>
```

or

```bash
file-sync --remove-device <number>
```

Removing this device will delete all configuration information (but not user files) from this device the next time it is synced. If this device needs to synchronize files again, you need to log in again for device registration.

### List of files

```bash
file-sync list
```

### Add file sync item

```bash
file-sync add <file path>
```

### Add a file sync item that already exists

```bash
file-sync add <file id> <file path>
```

### Set as boot-up service

```bash
file-sync service --enable
```

### Disable boot-up service

```bash
file-sync service --disable
```

### Turn on synchronization service

```bash
file-sync service --start
```

### Turn off synchronization services

```bash
file-sync service --stop
```

### Remove file sync item

```bash
file-sync remove <file id>
```

Note that this will remove all sync items from the file.

### Remove file sync item from a single device

```bash
file-sync remove <file id> <number>
```

or

```bash
file-sync remove <file id> <device id>
```

Hint: All `<device id>`s can be abbreviated.

## FAQ

**Q: How do I synchronize multiple files?**

A: This project is designed to synchronize a single configuration file, and you can synchronize multiple files by synchronizing the configuration entries of other file synchronization tools via `file-sync`.

## Use of self-hosted servers

You can choose to deploy it in your own server with the Docker or build your own PHP runtime.

The server side uses standard HTTP API, no rare module is used, which is very compatible, so you can deploy the program in most of the Virtual Hosting.

The PHP code provided by default needs to be used with `MySQL 5.4+` database.

### Docker

```bash
docker run xx:file-sync-server
```

### PHP

require PHP >= v5.4 (64bit)

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

## Server-side API

The `file-sync` program currently uses the `HTTP API` to complete synchronization interactions. Currently ~~ has completed ~~ the PHP version of the server-side API.

<details><summary>CLICK ME</summary>
<p>

Due to frequent changes, currently listed in the Chinese README： [简体中文](README_CHS.md)

</p>
</details>

## Development

<details><summary>CLICK ME</summary>
<p>

### Start the development and debugging environment

You need to set the environment variable `GO_ENV` to `development` manually.

For example, in the `PowerShell` environment in the `Windows` platform:

```bash
$Env:GO_ENV = 'development'
```

For example, `*unix`:

```bash
export GO_ENV="development"
```

</p>
</details>

## About Safety

According to the design principle of `server is untrustworthy`, all the data stored in the server is encrypted. Since `file-sync` uses asymmetric encryption and is used with SSL on the extranet, it is also secure during transmission.

All versions of encrypted files are stored in the server.

If you are offline while editing a file, `file-sync` will automatically record the last edit time and synchronize it with the version on the server when it is connected to the network. In the meantime, if the same file is edited in another device and the edited content is not the same, a conflict will inevitably arise. The conflicting version will be stored in the same directory of all synced devices as `[filename].[date].backup`.

## TODO

- Support generating guest accounts to synchronize files between trusted and untrusted devices in one way/both ways.
