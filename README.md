# SteamHostSync

[![Stars](https://img.shields.io/github/stars/Clov614/SteamHostSync?style=flat-square)](https://github.com/Clov614/SteamHostSync/stargazers)
[![Go](https://img.shields.io/badge/go-1.23%2B-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![License](https://img.shields.io/github/license/Clov614/SteamHostSync?style=flat-square)](https://github.com/Clov614/SteamHostSync/blob/main/LICENSE)

SteamHostSync is a small Go utility and ready-to-use hosts source for syncing Steam, GitHub, Docker, GOG, Ubisoft, and related domains.

适合想通过 SwitchHosts 或手动维护 hosts 的用户，快速获取已整理好的可更新配置。

## Table of Contents

- [Why This Repo Exists](#why-this-repo-exists)
- [What It Generates](#what-it-generates)
- [Quick Start](#quick-start)
- [Manual Setup](#manual-setup)
- [Refresh DNS](#refresh-dns)
- [Customize Source.yaml](#customize-sourceyaml)
- [Build From Source](#build-from-source)
- [Repository Layout](#repository-layout)

## Why This Repo Exists

这个项目会定期整理常见平台域名对应的 IP，并输出可直接使用的 hosts 文件，帮助你更方便地访问：

- Steam
- GitHub
- Docker Hub
- GOG Galaxy
- Ubisoft download services

## What It Generates

The repo publishes multiple host lists so you can subscribe to only what you need:

- `Hosts` for the full bundle
- `Hosts_steam` for Steam only
- `Hosts_github` for GitHub only
- `Hosts_docker` for Docker-related domains
- `Hosts_gog galaxy` for GOG services
- `Hosts_Ubisoft_download` for Ubisoft download services

## Quick Start

### Option 1: Use SwitchHosts (recommended)

[SwitchHosts](https://github.com/oldj/SwitchHosts) is the easiest way to keep the hosts file updated automatically.

备用下载源: <https://nas.iaimi.info/s/nT5pb8jMQp32QwB>

Available subscription URLs:

**Primary (jsDelivr)**

1. ALL: `https://cdn.jsdelivr.net/gh/Clov614/SteamHostSync@main/Hosts`
2. Steam: `https://cdn.jsdelivr.net/gh/Clov614/SteamHostSync@main/Hosts_steam`
3. GitHub: `https://cdn.jsdelivr.net/gh/Clov614/SteamHostSync@main/Hosts_github`

**Fallback (Statically)**

4. ALL: `https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts`
5. Steam: `https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts_steam`
6. GitHub: `https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts_github`

If the primary CDN is unreachable, switch to the matching fallback URL.

### Optional: Auto-start SwitchHosts on Windows

1. Press `Win + R`
2. Run `shell:startup`
3. Copy your SwitchHosts shortcut into that folder

![](/img/1.png)
![](/img/2.png)

### Configure automatic updates in SwitchHosts

![](/img/3.png)

## Manual Setup

### Hosts file locations

- Windows: `C:\Windows\System32\drivers\etc\hosts`
- Linux: `/etc/hosts`
- macOS: `/etc/hosts`

### Apply manually

Append the contents of one of the generated files, such as `Hosts` or `Hosts_steam`, to the end of your system hosts file.

Example:

```bash
# Linux / macOS
cat Hosts_steam | sudo tee -a /etc/hosts
```

Or open the file directly and paste the generated entries at the bottom.

## Refresh DNS

Most systems pick up the change quickly. If not, refresh DNS:

1. Windows: `ipconfig /flushdns`
2. Linux: `sudo nscd restart`
3. macOS: `sudo killall -HUP mDNSResponder`

## Customize Source.yaml

Running the executable once will generate `Source.yaml`, which defines the platforms and domains to resolve.

A simplified example:

```yaml
ua: Mozilla/5.0 ...
platforms:
  -
    - github
    - alive.github.com
    - github.com
    - raw.githubusercontent.com
  -
    - steam
    - steamcommunity.com
    - store.steampowered.com
```

This is useful if you want to add or remove platforms for your own host bundle.

## Build From Source

```bash
git clone https://github.com/Clov614/SteamHostSync.git
cd SteamHostSync
go run .
```

This updates the generated host files in the repository root and refreshes `README.md` from `README_TEMP.md` when present.

## Repository Layout

```text
main.go          Entry point
fileIO/          Config, file read, and write helpers
Source.yaml      Domain source definition
Hosts*           Generated host output files
img/             README screenshots
README_TEMP.md   Template used when regenerating README
```
