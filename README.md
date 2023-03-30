# `file-sync` is a file sync command

[![license](https://img.shields.io/github/license/yi-ge/file-sync.svg?style=flat-square)](https://github.com/yi-ge/file-sync/blob/master/LICENSE)
[![GitHub Actions](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Fyi-ge%2Ffile-sync%2Fbadge%3Fref%3Dmain&style=flat-square)](https://actions-badge.atrox.dev/yi-ge/file-sync/goto?ref=main)
[![Test Results](https://gist.github.com/yi-ge/535b9e71df85ad17d175c184f258b40a/raw/badge.svg)](https://github.com/yi-ge/file-sync)
[![GitHub last commit](https://img.shields.io/github/last-commit/yi-ge/file-sync.svg?style=flat-square)](https://github.com/yi-ge/file-sync)
<!-- [![Coveralls github](https://img.shields.io/coveralls/github/yi-ge/file-sync?style=flat-square)](https://coveralls.io/github/yi-ge/file-sync?branch=main) -->

[简体中文](README_CHS.md)

⚠️⛔️ This program is not yet under development !!! ⛔️ ⚠️

Automatically sync single file. Securely synchronize `.env` file or `.config` file for single user.

The design principle of `file-sync`: The server is untrustworthy, the client (local) is trusted.

## Feature

- Secure - when synchronizing data, your data is encrypted for transmission
- Flexible - supports cross-platform work. Synchronize the contents of the same file on different systems, different paths, different file names and different file permissions
- Automated - files are synchronized automatically without manual intervention, and you don't lose any data

## Install

### *unix system (Linux, macOS, BSD, etc.)

```bash
curl -sSL https://file-sync.yizcore.xyz/setup.sh | bash
```

### Windows

```powershell
iex ((New-Object System.Net.WebClient).DownloadString('https://file-sync.yizcore.xyz/setup.ps1'))
```

## Usage

### Set the server URL

```bash
file-sync --config https://example.com
```

The free public server provided by the command author is used by default: <https://file-sync.yizcore.xyz>

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
file-sync --remove-device <machine id>
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

- `--name`: Specify the name of the file to be recognized
- `--machineId`: Adding configurations for other devices

example:

```bash
file-sync add --name profile --machineId 123 /root/.profile
```

Configure the device with `machineId` as "123", add `/root/.profile` to the file sync item, and display the file configuration name as `profile`.

If the `~/.file-sync` folder is removed directly without removing the device, the device cannot be registered again to ensure data security unless it is removed from other devices in the same account.

### Add a file sync item that already exists

```bash
file-sync add <file id> <file path>
```

### Set as boot-up service

```bash
file-sync service enable
```

Note: The install script will automatically set the boot-up service.

### Disable boot-up service

```bash
file-sync service disable
```

### Turn on synchronization service

```bash
file-sync service start
```

Note: The install script will automatically start the synchronization service.

### Turn off synchronization services

```bash
file-sync service stop
```

### Remove file sync item

```bash
file-sync remove <file id>
```

Note that this will remove all sync items from the file.

### Remove file sync item from a single device

```bash
file-sync remove <file id / number>
```

or

```bash
file-sync remove --machineId <machine id> <file id>
```

Hint: All `<machine id>`, `<file id>` can be abbreviated.

## FAQ

**Q: How do I synchronize multiple files?**

A: This project is designed to synchronize a single configuration file, and you can synchronize multiple files by synchronizing the configuration entries of other file synchronization tools via `file-sync`.

**Q: Why don't work with NFS, SMB, FUSE, /proc, or /sys?**

A: `file-sync` require `fsnotify`, `fsnotify` requires support from underlying OS to work. The current NFS and SMB protocols does not provide network level support for file notifications, and neither do the /proc and /sys virtual filesystems.

## Use self-hosted server deploy (optional)

You can choose to deploy it in your own server with the Docker or build your own PHP runtime.

The server side uses standard HTTP API, no rare module is used, which is very compatible, so you can deploy the program in most of the Virtual Hosting.

The PHP code provided by default needs to be used with `MySQL 5.4+` database.

### Docker

```bash
docker run xx:file-sync-server
```

### PHP

require PHP >= v5.4 (64bit), It is recommended to turn on `shmop` and `mbstring` expansion for a better experience.

Upload the files in the `server/php` directory to the php root directory (excluding the `test` folder).

Note: Never upload the `.env` file from the PHP code directory to the `Server`/`Virtual Host` to avoid leaking the database configuration information.

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

The `file-sync` program currently uses the `HTTP API` to complete synchronization interactions, Server-to-client push is implemented using Server-Sent Events (SSE). Currently ~~ has completed ~~ the PHP version of the server-side API.

<details><summary>CLICK ME</summary>
<p>

Due to frequent changes, currently listed in the Chinese README： [简体中文](README_CHS.md)

</p>
</details>

## Development

<details><summary>CLICK ME</summary>
<p>

### Start the development and debugging environment

In the root file directory has `.env` environment variable configuration file, `GO_ENV` development environment value is `development` and production environment value is `production`.

#### Windows

Install `xampp` and configure `Zend Debugger`, change `DocumentRoot` and `Directory` in `httpd.conf` file to the absolute path where the `server/php` folder is located.

Start Apache, MySQL, and go to `http://localhost/phpmyadmin` to create a database named `file_sync`.

Modify the `.env.example` file in the root directory, and the environment variables in the `server/php/.htaccess.example` file.

**Note:** In `Windows` platform, `PHP_CLI_SERVER_WORKERS` environment variable is not supported, so please use the recommended latest version of `xampp` or `LAMP`, `LNMP` configuration for development and debugging in `Windows` platform. VSCode launch configuration is not applicable to `Windows` platform, do not use F5 to start the debugging environment.

### *unix

Install PHP 5.4+ and MySQL 5.4+, enable `shmop` extension, configure `Zend Debugger`, and create a database named `file_sync`.

Refer to the `.env.example` file in the root directory and the `server/php/.env.example` file for detailed environment variable configuration. Configure the `.htaccess` file according to the `Use self-hosted server` above.

Please set the `PHP_CLI_SERVER_WORKERS` environment variable to a value greater than `1` in order to test the working state of PHP in a Multi-threaded environment (relying on PHP CLI version >= 7.4.0, if you are developing with a lower version of PHP, please configure the `LNMP` or `LAMP` environment).

</p>
</details>

## Testing

<details><summary>CLICK ME</summary>
<p>

### Unit tests

```bash
go test . /...
```

### Integration tests

In the ``test`` folder, run the following command:

```bash
go run .
```

</p>
</details>

## About Safety

According to the design principle of `server is untrustworthy`, all the data stored in the server is encrypted. Since `file-sync` uses asymmetric encryption and is used with SSL on the extranet, it is also secure during transmission.

All versions of encrypted files are stored in the server.

If you are offline while editing a file, `file-sync` will automatically record the last edit time and synchronize it with the version on the server when it is connected to the network. In the meantime, if the same file is edited in another device and the edited content is not the same, a conflict will inevitably arise. The conflicting version will be stored in the same directory of all synced devices as `[filename].[date].backup`.

## TODO

- Support generating guest accounts to synchronize files between trusted and untrusted devices in one way/both ways.
