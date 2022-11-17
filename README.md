# `file-sync` is a file sync command

[![license](https://img.shields.io/github/license/yi-ge/file-sync.svg?style=flat-square)](https://github.com/yi-ge/file-sync/blob/master/LICENSE)
[![GitHub Actions](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Fyi-ge%2Ffile-sync%2Fbadge%3Fref%3Dmain&style=flat-square)](https://actions-badge.atrox.dev/yi-ge/file-sync/goto?ref=main)
[![Test Results](https://gist.github.com/yi-ge/00fdcacb47689d14b8e9fdf7fb0f7288/raw/badge.svg)](https://github.com/yi-ge/file-sync)
[![Coveralls github](https://img.shields.io/coveralls/github/yi-ge/file-sync?style=flat-square)](https://coveralls.io/github/yi-ge/file-sync?branch=main)
[![GitHub last commit](https://img.shields.io/github/last-commit/yi-ge/file-sync.svg?style=flat-square)](https://github.com/yi-ge/file-sync)

[简体中文](README_CHS.md)

⚠️⛔️ This program is not yet under development !!! ⛔️ ⚠️

Automatically sync single file. Sync `.env` file or `.config` file for single user.

The design principle of `file-sync`: the server is not trusted, the client (local) is trusted.

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

### List devices

```bash
file-sync list --device
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
