# SteamHostSync

[![Stars](https://img.shields.io/github/stars/Clov614/SteamHostSync?style=flat-square)](https://github.com/Clov614/SteamHostSync/stargazers)
[![Go](https://img.shields.io/badge/go-1.22%2B-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![License](https://img.shields.io/github/license/Clov614/SteamHostSync?style=flat-square)](https://github.com/Clov614/SteamHostSync/blob/main/LICENSE)

SteamHostSync 是一个小型 Go 工具兼现成的 hosts 源，用于同步 Steam、GitHub、Docker、GOG、Ubisoft 及相关域名的解析结果。

适合想通过 SwitchHosts 或手动维护 hosts 的用户，快速获取已整理好的可更新配置。

> **Languages / 语言**: [English](README.en.md) · [简体中文](README.zh-CN.md)

## 目录

- [为什么有这个项目](#为什么有这个项目)
- [生成的内容](#生成的内容)
- [快速开始](#快速开始)
- [手动配置](#手动配置)
- [刷新 DNS](#刷新-dns)
- [自定义 config.yaml](#自定义-configyaml)
- [从源码构建](#从源码构建)
- [仓库结构](#仓库结构)

## 为什么有这个项目

这个项目会定期整理常见平台域名对应的 IP，并输出可直接使用的 hosts 文件，帮助你更方便地访问：

- Steam
- GitHub
- Docker Hub
- GOG Galaxy
- Ubisoft download services

## 生成的内容

仓库会发布多个 hosts 列表，你可以只订阅所需的部分：

- `Hosts` —— 全部平台的合集
- `Hosts_steam` —— 仅 Steam
- `Hosts_github` —— 仅 GitHub
- `Hosts_docker` —— Docker 相关域名
- `Hosts_gog` —— GOG 服务
- `Hosts_ubisoft` —— Ubisoft 下载服务

## 快速开始

### 方式一：使用 SwitchHosts（推荐）

[SwitchHosts](https://github.com/oldj/SwitchHosts) 是最方便的自动更新 hosts 文件的方式。

可用的订阅地址：

**主源（jsDelivr）**

1. ALL：`https://cdn.jsdelivr.net/gh/Clov614/SteamHostSync@main/Hosts`
2. Steam：`https://cdn.jsdelivr.net/gh/Clov614/SteamHostSync@main/Hosts_steam`
3. GitHub：`https://cdn.jsdelivr.net/gh/Clov614/SteamHostSync@main/Hosts_github`

**备用源（Statically）**

4. ALL：`https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts`
5. Steam：`https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts_steam`
6. GitHub：`https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts_github`

若主源不可访问，请切换到对应的备用源。

### 可选：让 SwitchHosts 在 Windows 开机自启

1. 按 `Win + R`
2. 输入 `shell:startup`
3. 将 SwitchHosts 的快捷方式复制进该文件夹

![](/img/1.png)
![](/img/2.png)

### 在 SwitchHosts 中配置自动更新

![](/img/3.png)

## 手动配置

### hosts 文件位置

- Windows：`C:\Windows\System32\drivers\etc\hosts`
- Linux：`/etc/hosts`
- macOS：`/etc/hosts`

### 手动应用

将生成的某个文件（如 `Hosts` 或 `Hosts_steam`）的内容追加到系统 hosts 文件的末尾。

示例：

```bash
# Linux / macOS
cat Hosts_steam | sudo tee -a /etc/hosts
```

或直接打开文件，把生成的条目粘贴到底部。

## 刷新 DNS

大部分系统会很快生效。若未生效，可以刷新 DNS：

1. Windows：`ipconfig /flushdns`
2. Linux：`sudo nscd restart`
3. macOS：`sudo killall -HUP mDNSResponder`

## 自定义 config.yaml

运行一次可执行文件会生成 `config.yaml`，它定义了需要解析的平台与域名、DNS-over-HTTPS 上游、并发数与超时。

简化示例：

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

## 从源码构建

```bash
git clone https://github.com/Clov614/SteamHostSync.git
cd SteamHostSync
go run .
```

这会更新仓库根目录下生成的 hosts 文件，并在存在 `README_TEMP.md` 时刷新 `README.md`。

## 仓库结构

```text
main.go                  Entry point (flags + orchestration)
internal/                Layered packages: config / resolve / probe / render / fileio / app
config.yaml              Runtime configuration (auto-generated on first run)
Hosts*                   Generated host output files
img/                     README screenshots
README_TEMP.md           Template used when regenerating README (host content placeholder)
```