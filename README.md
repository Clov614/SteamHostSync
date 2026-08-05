# SteamHostSync

[![Stars](https://img.shields.io/github/stars/Clov614/SteamHostSync?style=flat-square)](https://github.com/Clov614/SteamHostSync/stargazers)
[![Go](https://img.shields.io/badge/go-1.22%2B-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![License](https://img.shields.io/github/license/Clov614/SteamHostSync?style=flat-square)](https://github.com/Clov614/SteamHostSync/blob/main/LICENSE)

SteamHostSync is a small Go utility and ready-to-use hosts source for syncing Steam, GitHub, Docker, GOG, Ubisoft, and related domains.

适合想通过 SwitchHosts 或手动维护 hosts 的用户，快速获取已整理好的可更新配置。

## Table of Contents

- [Why This Repo Exists](#why-this-repo-exists)
- [What It Generates](#what-it-generates)
- [Quick Start](#quick-start)
- [Manual Setup](#manual-setup)
- [Refresh DNS](#refresh-dns)
- [Customize config.yaml](#customize-configyaml)
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
- `Hosts_gog` for GOG services
- `Hosts_ubisoft` for Ubisoft download services

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

## Customize config.yaml

Running the executable once will generate `config.yaml`, which defines the platforms and domains to resolve, the DNS-over-HTTPS upstreams, concurrency and timeouts.

A simplified example:

```yaml
version: 1
concurrency: 8
timeout:
  resolve: 5s
  probe: 2s
probe:
  port: 443
  attempts: 3
dns_servers:
  - https://dns.alidns.com/resolve
  - https://doh.pub/dns-query
  - https://dns.google/resolve
platforms:
  - name: github
    domains:
      - alive.github.com
      - github.com
      - raw.githubusercontent.com
  - name: steam
    domains:
      - steamcommunity.com
      - store.steampowered.com
```

说明：

- 每个 `platforms` 项生成一个 `Hosts_<name>` 文件，所有平台合并生成 `Hosts`。
- 平台名会转为小写安全文件名（非法字符替换为 `_`），推荐使用无空格的小写名称。
- 解析失败的域名会以 `# domain` 注释行保留，不影响其他域名。
- 若 config.yaml 不存在，首次运行会自动写入内嵌默认配置。

## Build From Source

```bash
git clone https://github.com/Clov614/SteamHostSync.git
cd SteamHostSync
go run .
```

This updates the generated host files in the repository root and refreshes `README.md` from `README_TEMP.md` when present.

## Repository Layout

```text
main.go                  Entry point (flags + orchestration)
internal/                Layered packages: config / resolve / probe / render / fileio / app
config.yaml              Runtime configuration (auto-generated on first run)
Hosts*                   Generated host output files
img/                     README screenshots
README_TEMP.md           Template used when regenerating README (host content placeholder)
```

## Current hosts

```
# SteamHostSync hosts v1
# Generated: 2026-08-05T16:33:57Z
# Project: https://github.com/Clov614/SteamHostSync

# github Start
140.82.113.25			alive.github.com
140.82.113.25			live.github.com
185.199.108.215			github.githubassets.com
140.82.112.21			central.github.com
185.199.111.133			desktop.githubusercontent.com
# assets-cdn.github.com
185.199.108.133			camo.githubusercontent.com
185.199.109.133			github.map.fastly.net
151.101.129.194			github.global.ssl.fastly.net
140.82.116.4			gist.github.com
185.199.111.153			github.io
140.82.116.4			github.com
192.0.66.2			github.blog
140.82.116.5			api.github.com
185.199.111.133			raw.githubusercontent.com
185.199.111.133			user-images.githubusercontent.com
185.199.109.133			favicons.githubusercontent.com
185.199.110.133			avatars5.githubusercontent.com
185.199.108.133			avatars4.githubusercontent.com
185.199.111.133			avatars3.githubusercontent.com
185.199.111.133			avatars2.githubusercontent.com
185.199.110.133			avatars1.githubusercontent.com
185.199.108.133			avatars0.githubusercontent.com
185.199.110.133			avatars.githubusercontent.com
140.82.116.9			codeload.github.com
52.217.236.121			github-cloud.s3.amazonaws.com
16.15.199.183			github-com.s3.amazonaws.com
16.15.246.228			github-production-release-asset-2e65be.s3.amazonaws.com
16.15.213.78			github-production-user-asset-6210df.s3.amazonaws.com
16.182.64.73			github-production-repository-file-5c1aeb.s3.amazonaws.com
185.199.111.153			githubstatus.com
140.82.114.18			github.community
20.99.227.183			github.dev
185.199.109.133			media.githubusercontent.com
# github End # Last Update Time : 2026-08-05T16:33:57Z
# SteamHostSync hosts v1
# Generated: 2026-08-05T16:33:57Z
# Project: https://github.com/Clov614/SteamHostSync

# steam Start
96.6.236.246			steamcommunity.com
23.213.145.49			www.steamcommunity.com
23.199.21.199			store.steampowered.com
96.6.236.246			api.steampowered.com
104.68.104.163			help.steampowered.com
23.62.46.217			store.akamai.steamstatic.com
23.67.33.209			steamcdn-a.akamaihd.net
23.46.216.87			cdn.akamai.steamstatic.com
23.213.145.49			steam-chat.com
23.46.216.70			community.akamai.steamstatic.com
# steam End # Last Update Time : 2026-08-05T16:33:57Z
# SteamHostSync hosts v1
# Generated: 2026-08-05T16:33:57Z
# Project: https://github.com/Clov614/SteamHostSync

# docker Start
23.185.0.4			docker.com
104.18.43.187			hub.docker.com
65.8.54.82			docs.docker.com
172.64.144.74			login.docker.com
23.21.206.0			registry.hub.docker.com
3.215.0.184			docker.io
23.21.206.0			registry-1.docker.io
35.172.43.215			index.docker.io
# docker End # Last Update Time : 2026-08-05T16:33:57Z
# SteamHostSync hosts v1
# Generated: 2026-08-05T16:33:57Z
# Project: https://github.com/Clov614/SteamHostSync

# gog Start
151.101.1.241			auth.gog.com
151.101.1.241			www.gogalaxy.com
151.101.1.241			remote-config.gog.com
151.101.65.241			insights-collector.gog.com
151.101.193.241			gameplay.gog.com
151.101.1.241			gamesdb.gog.com
151.101.129.241			external-accounts.gog.com
151.101.65.241			www.gog.com
# gog End # Last Update Time : 2026-08-05T16:33:57Z
# SteamHostSync hosts v1
# Generated: 2026-08-05T16:33:57Z
# Project: https://github.com/Clov614/SteamHostSync

# ubisoft Start
23.40.25.65			static3.cdn.Ubi.com
23.37.17.207			static2.cdn.Ubi.com
# static1.cdn.Ubi.com
# ubisoft End # Last Update Time : 2026-08-05T16:33:57Z
# Github: https://github.com/Clov614/SteamHostSync

```
