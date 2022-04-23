# SteamHostSync
第一次用go写的项目，写的比较烂，欢迎大佬指出错误。

## 1. 实现
对Hosts进行一个新的更  
解决Steam、github访问问题

## 2. 使用方法
## 自动方法(使用工具)
推荐使用Hosts管理工具[SwitchHosts](https://github.com/oldj/SwitchHosts)
### 开机自启动SwitchHosts
win + R 后执行 `shell:startup`    
![](/img/1.png)  
将快捷方式复制进去即可  
![](/img/2.png)

## 手动方式
#### 1. hosts 文件在每个系统的位置不一，详情如下:
- Windows 系统：`C:\Windows\System32\drivers\etc\hosts`
- Linux 系统：`/etc/hosts`
- Mac（苹果电脑）系统：`/etc/hosts`

#### 2. 修改方法
复制下面的内容至hosts尾部(追加在文本末尾)

```
#github Start
140.82.114.26			alive.github.com
140.82.113.25			live.github.com
185.199.108.154			github.githubassets.com
140.82.113.22			central.github.com
185.199.108.133			desktop.githubusercontent.com
185.199.108.153			assets-cdn.github.com
185.199.108.133			camo.githubusercontent.com
185.199.108.133			github.map.fastly.net
199.232.69.194			github.global.ssl.fastly.net
140.82.112.3			gist.github.com
185.199.108.153			github.io
140.82.112.4			github.com
192.0.66.2			github.blog
140.82.113.5			api.github.com
185.199.108.133			raw.githubusercontent.com
185.199.108.133			user-images.githubusercontent.com
185.199.108.133			favicons.githubusercontent.com
185.199.108.133			avatars5.githubusercontent.com
185.199.108.133			avatars4.githubusercontent.com
185.199.108.133			avatars3.githubusercontent.com
185.199.108.133			avatars2.githubusercontent.com
185.199.108.133			avatars1.githubusercontent.com
185.199.108.133			avatars0.githubusercontent.com
185.199.108.133			avatars.githubusercontent.com
140.82.114.9			codeload.github.com
52.217.135.81			github-cloud.s3.amazonaws.com
52.216.186.75			github-com.s3.amazonaws.com
52.217.201.201			github-production-release-asset-2e65be.s3.amazonaws.com
54.231.130.161			github-production-user-asset-6210df.s3.amazonaws.com
52.217.92.116			github-production-repository-file-5c1aeb.s3.amazonaws.com
185.199.108.153			githubstatus.com
64.71.144.211			github.community
23.100.27.125			github.dev
185.199.108.133			media.githubusercontent.com
#github End
# Last Update Time : 2022-04-23 14:37:59 

#steam Start
23.204.9.127			steamcommunity.com
23.204.9.127			www.steamcommunity.com
23.35.76.32			store.steampowered.com
23.45.0.161			api.steampowered.com
104.98.68.175			help.steampowered.com
23.33.29.72			store.akamai.steamstatic.com
23.220.246.175			steamcdn-a.akamaihd.net
23.33.29.68			cdn.akamai.steamstatic.com
23.3.117.102			steam-chat.com
23.222.236.17			community.akamai.steamstatic.com
#steam End
# Last Update Time : 2022-04-23 14:38:10 

Github: https://github.com/Clov614/SteamHostSync

```

#### 3. 激活生效
大部分情况下是直接生效，如未生效可尝试下面的办法，刷新 DNS：
1. Windows：在 CMD 窗口输入：`ipconfig /flushdns`
2. Linux 命令：`sudo nscd restart`
3. Mac 命令：`sudo killall -HUP mDNSResponder`  
