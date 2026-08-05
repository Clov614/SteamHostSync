# SteamHostSync
第一次用go写的项目，写的比较烂，欢迎大佬指出错误。

## 1. 实现
对Hosts进行一个新的更  
解决Steam、github访问问题

## 2. 使用方法
## 自动方法(使用工具)
推荐使用Hosts管理工具[SwitchHosts](https://github.com/oldj/SwitchHosts) 
[SwitchHosts备用下载源](https://nas.iaimi.info/s/nT5pb8jMQp32QwB)
### 开机自启动SwitchHosts
win + R 后执行 `shell:startup`    
![](/img/1.png)  
将快捷方式复制进去即可  
![](/img/2.png)  
### 配置SwitchHosts实现自动更新  
可选的URL有:
主源（jsDelivr）:
1. ALL: `https://cdn.jsdelivr.net/gh/Clov614/SteamHostSync@main/Hosts`
2. Steam: `https://cdn.jsdelivr.net/gh/Clov614/SteamHostSync@main/Hosts_steam`
3. github: `https://cdn.jsdelivr.net/gh/Clov614/SteamHostSync@main/Hosts_github`
备用源（Statically）:
4. ALL: `https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts`
5. Steam: `https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts_steam`
6. github: `https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts_github`
说明：若主源访问失败，请切换到对应的备用源链接。

![](/img/3.png)

## 手动方式
#### 1. hosts 文件在每个系统的位置不一，详情如下:
- Windows 系统：`C:\Windows\System32\drivers\etc\hosts`
- Linux 系统：`/etc/hosts`
- Mac（苹果电脑）系统：`/etc/hosts`

#### 2. 修改方法
复制下面的内容至hosts尾部(追加在文本末尾)

```
# SteamHostSync hosts v1
# Generated: 2026-08-05T15:39:12Z
# Project: https://github.com/Clov614/SteamHostSync

# github Start
140.82.113.25			alive.github.com
140.82.113.26			live.github.com
185.199.108.215			github.githubassets.com
140.82.114.22			central.github.com
185.199.110.133			desktop.githubusercontent.com
# assets-cdn.github.com
185.199.110.133			camo.githubusercontent.com
185.199.108.133			github.map.fastly.net
151.101.193.194			github.global.ssl.fastly.net
140.82.116.3			gist.github.com
185.199.110.153			github.io
140.82.116.4			github.com
192.0.66.2			github.blog
140.82.116.5			api.github.com
185.199.109.133			raw.githubusercontent.com
185.199.110.133			user-images.githubusercontent.com
185.199.111.133			favicons.githubusercontent.com
185.199.110.133			avatars5.githubusercontent.com
185.199.110.133			avatars4.githubusercontent.com
185.199.108.133			avatars3.githubusercontent.com
185.199.109.133			avatars2.githubusercontent.com
185.199.109.133			avatars1.githubusercontent.com
185.199.109.133			avatars0.githubusercontent.com
185.199.109.133			avatars.githubusercontent.com
140.82.116.9			codeload.github.com
16.15.223.210			github-cloud.s3.amazonaws.com
16.15.212.255			github-com.s3.amazonaws.com
16.15.223.197			github-production-release-asset-2e65be.s3.amazonaws.com
52.217.227.65			github-production-user-asset-6210df.s3.amazonaws.com
16.15.236.209			github-production-repository-file-5c1aeb.s3.amazonaws.com
185.199.108.153			githubstatus.com
140.82.114.17			github.community
20.99.227.183			github.dev
185.199.110.133			media.githubusercontent.com
# github End # Last Update Time : 2026-08-05T15:39:12Z
# SteamHostSync hosts v1
# Generated: 2026-08-05T15:39:12Z
# Project: https://github.com/Clov614/SteamHostSync

# steam Start
96.6.236.246			steamcommunity.com
104.68.112.116			www.steamcommunity.com
23.199.21.199			store.steampowered.com
96.6.236.246			api.steampowered.com
96.6.236.246			help.steampowered.com
23.192.228.19			store.akamai.steamstatic.com
23.62.46.211			steamcdn-a.akamaihd.net
23.192.228.24			cdn.akamai.steamstatic.com
23.213.145.49			steam-chat.com
23.192.228.27			community.akamai.steamstatic.com
# steam End # Last Update Time : 2026-08-05T15:39:12Z
# SteamHostSync hosts v1
# Generated: 2026-08-05T15:39:12Z
# Project: https://github.com/Clov614/SteamHostSync

# docker Start
23.185.0.4			docker.com
172.64.144.69			hub.docker.com
65.8.54.107			docs.docker.com
104.18.43.182			login.docker.com
35.168.108.241			registry.hub.docker.com
52.87.135.177			docker.io
44.194.160.0			registry-1.docker.io
104.244.43.229			index.docker.io
# docker End # Last Update Time : 2026-08-05T15:39:12Z
# SteamHostSync hosts v1
# Generated: 2026-08-05T15:39:12Z
# Project: https://github.com/Clov614/SteamHostSync

# gog Start
151.101.193.241			auth.gog.com
151.101.129.241			www.gogalaxy.com
151.101.129.241			remote-config.gog.com
151.101.193.241			insights-collector.gog.com
151.101.129.241			gameplay.gog.com
151.101.129.241			gamesdb.gog.com
151.101.193.241			external-accounts.gog.com
151.101.129.241			www.gog.com
# gog End # Last Update Time : 2026-08-05T15:39:12Z
# SteamHostSync hosts v1
# Generated: 2026-08-05T15:39:12Z
# Project: https://github.com/Clov614/SteamHostSync

# ubisoft Start
23.199.21.72			static3.cdn.Ubi.com
23.36.21.209			static2.cdn.Ubi.com
# static1.cdn.Ubi.com
# ubisoft End # Last Update Time : 2026-08-05T15:39:12Z
# Github: https://github.com/Clov614/SteamHostSync

```

## 激活生效
大部分情况下是直接生效，如未生效可尝试下面的办法，刷新 DNS：
1. Windows：在 CMD 窗口输入：`ipconfig /flushdns`
2. Linux 命令：`sudo nscd restart`
3. Mac 命令：`sudo killall -HUP mDNSResponder`  

## 手动配置 config.yaml 文件添加新 hosts  
手动下载可执行文件后，第一次执行会在当前目录生成 `config.yaml`，可手动配置。  

```
version: 1
concurrency: 8                 # 并发解析/探测上限
timeout:
  resolve: 5s                  # 单个 DoH 查询超时
  probe: 2s                    # 单次 TCP 握手超时
probe:
  port: 443
  attempts: 3                  # 测速次数，取中位数
dns_servers:                   # 至少 1 个，多上游并发合并去重
  - https://dns.alidns.com/resolve
  - https://doh.pub/dns-query
  - https://dns.google/resolve
platforms:
  - name: github               # 平台名，生成 Hosts_github 文件
    domains:
      - alive.github.com
      - github.com
      - api.github.com
  - name: steam                # 生成 Hosts_steam 文件
    domains:
      - store.steampowered.com
      - steamcommunity.com
```

说明：
- 每个 `platforms` 项生成一个 `Hosts_<name>` 文件，所有平台合并生成 `Hosts`。
- 平台名会转为小写安全文件名（非法字符替换为 `_`），因此配置 `gog galaxy` 会生成 `Hosts_gog_galaxy`。推荐直接使用无空格的小写名称。
- 解析失败的域名会以 `# domain` 注释行保留，不影响其他域名。

