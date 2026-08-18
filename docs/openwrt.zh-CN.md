# 在 OpenWrt 上自动更新 SteamHostSync

SteamHostSync 的 OpenWrt 集成通过 dnsmasq 的 `addnhosts` 加载独立 hosts 文件，不会追加、覆盖或接管 `/etc/hosts`。更新器以计划任务形式短暂运行，没有常驻 SteamHostSync 进程。

## 支持范围

当前集成面向使用 UCI 和 dnsmasq 的 OpenWrt。hosts 状态和临时输出写入 `/var/lib/steamhostsync`；更新锁位于 `/var/lock/steamhostsync.lock`。OpenWrt 的 `/var` 通常位于内存文件系统中，可以避免周期性写入闪存。

提供两种来源模式：

- `remote`（默认）：下载仓库已经生成的 `Hosts` 或某个 `Hosts_<平台>` 文件。资源占用最低，与 CPU 架构无关。
- `local`（可选）：在路由器本地运行 SteamHostSync，执行 DoH 解析和 TCP 探测，然后通过相同的校验、发布和回滚流程应用结果。

Steam++ 式“hosts 指向回环地址 + HTTPS 代理”不在本功能范围内。本集成不会安装证书、监听 80/443 或修改防火墙。

## 构建并安装 OpenWrt 包

将 `contrib/openwrt` 复制到 OpenWrt SDK 或源码树的 `package/steamhostsync-openwrt`，然后执行：

```sh
make package/steamhostsync-openwrt/compile V=s
```

把生成的 `steamhostsync-openwrt_*.ipk` 上传到路由器并安装：

```sh
opkg install steamhostsync-openwrt_*.ipk
steamhostsyncctl setup
steamhostsyncctl update
steamhostsyncctl status
```

`setup` 会进行以下幂等操作：

1. 创建 `/var/lib/steamhostsync/Hosts`。
2. 将该路径加入 dnsmasq `addnhosts` 列表，不删除已有条目。
3. 安装每天两次的 cron 任务（每 12 小时一次）。
4. 启用开机 one-shot 更新服务。

重复执行 `setup` 不会重复添加 UCI 或 cron 条目。

## A 模式：远程订阅

默认配置位于 `/etc/config/steamhostsync`：

```text
config steamhostsync 'main'
        option source 'remote'
        option url 'https://cdn.jsdelivr.net/gh/Clov614/SteamHostSync@main/Hosts'
        option fallback_url 'https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts'
```

例如只订阅 Steam：

```sh
uci set steamhostsync.main.url='https://cdn.jsdelivr.net/gh/Clov614/SteamHostSync@main/Hosts_steam'
uci set steamhostsync.main.fallback_url='https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts_steam'
uci commit steamhostsync
steamhostsyncctl update
```

下载内容必须：

- 使用 HTTPS 并通过证书验证。
- 不超过配置的 `max_bytes`。
- 包含 SteamHostSync 文件头和至少一条有效记录。
- 只包含有效的公网 IPv4 与 hosts 域名。
- 在严格模式下，域名必须存在于 `/etc/steamhostsync/allowed-domains`。

默认情况下，主源和备用源都必须下载成功且字节完全一致，才会发布；这样可以降低单一 CDN 缓存或链路被污染的风险。下载或验证失败时，当前有效文件保持不变，dnsmasq 不会 reload。仓库发布者和 GitHub Actions 仍是内容信任根。

如需使用自定义域名，应把域名加入 allowlist。也可以设置 `strict_domains='0'`，但这会允许远程订阅控制任意域名的局域网解析，不建议使用。

## B 模式：路由器本地生成

本地模式需要先安装与路由器 CPU 匹配的静态 SteamHostSync 二进制。Release 提供：

- `SteamHostSync_linux_amd64.tar.gz`
- `SteamHostSync_linux_386.tar.gz`
- `SteamHostSync_linux_armv7.tar.gz`
- `SteamHostSync_linux_arm64.tar.gz`
- `SteamHostSync_linux_mips_softfloat.tar.gz`
- `SteamHostSync_linux_mipsle_softfloat.tar.gz`

传统 OpenWrt MIPS 设备必须确认端序：`mips` 是大端，`mipsle` 是小端。若不能确认，请不要尝试执行；错误架构通常会报 `Exec format error` 或 `Illegal instruction`。

解压后安装：

```sh
install -m 0755 SteamHostSync /usr/bin/steamhostsync
steamhostsyncctl mode local
steamhostsyncctl update
```

本地模式使用 `/etc/steamhostsync/router-config.yaml`，默认限制为：

- 并发数 `1`
- 每个 IP 探测 `1` 次
- 只解析 GitHub、Steam 和 Linux Steam 的核心域名
- 不生成 README
- 临时输出位于 `/var/lib/steamhostsync`

如需增加域名，请同时更新路由器 YAML 和 allowlist。Go 进程只在更新时运行，完成后退出。

资源较少的设备建议继续使用 A 模式。当前 MIPS/MIPSLE 产物已通过交叉编译，但仍应在具体设备上验证内存和指令集兼容性。

## 更新、回滚和停用

```sh
# 立即更新
steamhostsyncctl update

# 查看来源模式、状态文件和 dnsmasq 配置状态
steamhostsyncctl status

# 恢复本次启动以来的上一有效版本
steamhostsyncctl rollback

# 删除本项目的 cron 和 dnsmasq 条目，停用更新和 init 服务
steamhostsyncctl disable
```

更新器使用原子锁阻止开机任务、cron 和手动更新并发执行。新内容会先下载或生成到临时路径，完整验证后才替换。dnsmasq reload 失败时会自动恢复上一版本。

`Hosts` 和 `Hosts.previous` 位于内存文件系统，路由器重启后会丢失。init 服务会在启动后重新获取；如果当时 WAN 不可用，普通 DNS 仍能工作，但在下一次成功更新前没有 SteamHostSync 加速记录。

## 常见问题

### 局域网设备没有使用这些记录

确认客户端的 DNS 指向路由器。浏览器或系统启用了独立 DoH/DoT 时，查询会绕过 dnsmasq，本功能无法生效。

### 路由器使用 AdGuardHome

本包的自动接入只支持 dnsmasq，不会改写 AdGuardHome。可以在 AdGuardHome 中直接把项目的 `Hosts` URL 添加为 hosts 格式的自定义过滤源。

### TLS 下载失败

确认路由器时间正确，并安装了 `ca-bundle`。更新器不会通过关闭证书校验来绕过 TLS 错误。

### hosts 仍不能完全解决 Steam 访问问题

hosts 只影响域名解析。地区限制、TLS、HTTP 代理、客户端网络或上游服务问题不在它的处理范围内。
