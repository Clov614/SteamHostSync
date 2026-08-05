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
HOST_TARGET
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

