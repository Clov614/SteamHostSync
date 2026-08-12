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
- [Fork 客制化（让你的 fork 自动更新）](#fork-客制化让你的-fork-自动更新)
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
- `Hosts_steam_linux` —— 仅 Linux Steam（安装/更新所需仓库与客户端域名）
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
3. Linux Steam：`https://cdn.jsdelivr.net/gh/Clov614/SteamHostSync@main/Hosts_steam_linux`
4. GitHub：`https://cdn.jsdelivr.net/gh/Clov614/SteamHostSync@main/Hosts_github`

**备用源（Statically）**

5. ALL：`https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts`
6. Steam：`https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts_steam`
7. Linux Steam：`https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts_steam_linux`
8. GitHub：`https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts_github`

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

## Fork 客制化（让你的 fork 自动更新）

> 以下内容面向「把本仓库 fork 一份当作自己的 hosts 源」的用户，不是开发者协作指南。

fork 后，仓库自带的 CI（`.github/workflows/Update.yml`）并不会自动运行，需要按下面几步配置一次，之后它就会每 12 小时自动更新一次你的 hosts 文件：

### 1. 启用 Actions

fork 的定时触发（`schedule`）**默认是关闭的**，需要手动开启：

1. 打开你的 fork 仓库，点击 **Actions** 页签。
2. 首次会提示启用 Actions，点击 **I understand my workflows, go ahead and enable them**。
3. 侧边栏里若 `Update` 工作流显示 disabled，点开它的「…」菜单选择 **Enable workflow**（这会立即跑一次并重置定时）。

> ⚠️ **60 天无活动会被再次禁用**：GitHub 会对公共仓库中连续 60 天无活动的定时工作流自动停用。只要你的 fork 有新的 push 或每次手动触发，就能持续保活；真的被禁用后，重新执行上面的步骤即可。

### 2. 授予 GITHUB_TOKEN 写权限（关键步骤）

`Update.yml` 会把生成的 hosts 文件直接 `git push` 回你的 `main` 分支，这需要仓库默认的 `GITHUB_TOKEN` 有写权限。而新仓库默认只读，否则 CI 会在 push 时报 `Permission denied`。

1. 打开 **Settings → Actions → General**。
2. 找到 **Workflow permissions**，勾选 **Read and write permissions**，保存。

### 3. 自定义要解析的平台与域名（可选）

修改仓库根目录的 `config.yaml` 即可（`platforms` / `dns_servers` / `concurrency` / `timeout` / `probe`），字段说明见上一节。CI 每次运行会把这个文件一起提交，所以你的自定义配置会保留在 git 历史中。

> 提醒：请勿手动提交 CI 生成的 `Hosts*` 与 `README.md`，它们由工作流自动维护；本地构建若要生效同样的配置，记得同步修改 `internal/config/default_config.yaml`（详见 `CLAUDE.md`）。

### 4. 更新订阅地址（如使用 SwitchHosts）

把 README 里所有 CDN 链接中的 `Clov614` 换成你自己的 GitHub 用户名：

```
https://cdn.jsdelivr.net/gh/<你的用户名>/SteamHostSync@main/Hosts
```

> 注意：jsDelivr / Statically 只缓存公共仓库的内容，所以请保持你的 fork 为公开仓库；否则订阅地址无法访问。

### 5. （可选）发布 Release

不打 `v*` 标签就不会触发 `release.yml`（goreleaser），无需任何配置。若你想发布自己的版本，直接打一个 `v*` 标签即可，工作流会自动用 `GITHUB_TOKEN` 发布，无需额外密钥。

### 触发方式小结

- 定时：每 12 小时一次（需先完成第 1 步）
- `push` 到 `main`
- 手动：Actions 页签 → `Update` → **Run workflow**

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