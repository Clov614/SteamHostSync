# SteamHostSync

[![Stars](https://img.shields.io/github/stars/Clov614/SteamHostSync?style=flat-square)](https://github.com/Clov614/SteamHostSync/stargazers)
[![Go](https://img.shields.io/badge/go-1.22%2B-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![License](https://img.shields.io/github/license/Clov614/SteamHostSync?style=flat-square)](https://github.com/Clov614/SteamHostSync/blob/main/LICENSE)

SteamHostSync is a small Go utility and ready-to-use hosts source for syncing Steam, GitHub, Docker, GOG, Ubisoft, and related domains.

It is aimed at users who want to keep an organized, updatable hosts configuration with SwitchHosts or by hand.

> **Languages / 语言**: [English](README.en.md) · [简体中文](README.zh-CN.md)

## Table of Contents

- [Why This Repo Exists](#why-this-repo-exists)
- [What It Generates](#what-it-generates)
- [Quick Start](#quick-start)
- [Automatic updates on OpenWrt](#option-2-automatic-updates-on-openwrt)
- [Manual Setup](#manual-setup)
- [Refresh DNS](#refresh-dns)
- [Customize config.yaml](#customize-configyaml)
- [Build From Source](#build-from-source)
- [Repository Layout](#repository-layout)

## Why This Repo Exists

This project periodically resolves the IPs for common platform domains and publishes ready-to-use hosts files, so you can reach these services more reliably:

- Steam
- GitHub
- Docker Hub
- GOG Galaxy
- Ubisoft download services

## What It Generates

The repo publishes multiple host lists so you can subscribe to only what you need:

- `Hosts` for the full bundle
- `Hosts_steam` for Steam only
- `Hosts_steam_linux` for Linux Steam only (apt repo + client update domains)
- `Hosts_github` for GitHub only
- `Hosts_docker` for Docker-related domains
- `Hosts_gog` for GOG services
- `Hosts_ubisoft` for Ubisoft download services

## Quick Start

### Option 1: Use SwitchHosts (recommended)

[SwitchHosts](https://github.com/oldj/SwitchHosts) is the easiest way to keep the hosts file updated automatically.

Available subscription URLs:

**Primary (jsDelivr)**

1. ALL: `https://cdn.jsdelivr.net/gh/Clov614/SteamHostSync@main/Hosts`
2. Steam: `https://cdn.jsdelivr.net/gh/Clov614/SteamHostSync@main/Hosts_steam`
3. Linux Steam: `https://cdn.jsdelivr.net/gh/Clov614/SteamHostSync@main/Hosts_steam_linux`
4. GitHub: `https://cdn.jsdelivr.net/gh/Clov614/SteamHostSync@main/Hosts_github`

**Fallback (Statically)**

5. ALL: `https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts`
6. Steam: `https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts_steam`
7. Linux Steam: `https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts_steam_linux`
8. GitHub: `https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts_github`

If the primary CDN is unreachable, switch to the matching fallback URL.

### Option 2: Automatic updates on OpenWrt

OpenWrt users can install the architecture-neutral updater package and let dnsmasq load fresh hosts every 12 hours without modifying `/etc/hosts`:

- Mode A (default) subscribes to generated `Hosts*` files and uses no resident process.
- Mode B (optional) generates hosts on the router; releases include MIPS/MIPSLE softfloat archives containing the router binaries.

See the [OpenWrt guide](docs/openwrt.en.md) for installation, configuration, rollback, and security details. Routers using AdGuardHome can alternatively subscribe to an existing `Hosts` URL directly.

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

Notes:

- Each `platforms` entry produces a `Hosts_<name>` file; all platforms combined produce `Hosts`.
- Platform names are lowercased into safe filenames (illegal characters are replaced with `_`); lowercase names without spaces are recommended.
- Domains that fail to resolve are kept as `# domain` comment lines and do not affect the other entries.
- If `config.yaml` does not exist, the embedded default configuration is written automatically on the first run.

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
