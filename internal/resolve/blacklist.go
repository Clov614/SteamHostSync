package resolve

import (
	"net"
	"strings"
)

// isUnusableIP 判断 IP 是否为不应写入 hosts 的内网/保留地址。
// 覆盖：0.0.0.0、回环 127/8、RFC1918 私网（10/8, 172.16/12, 192.168/16）、
// 链路本地 169.254/16 等。此类地址通常意味着 DNS 被劫持或污染。
func isUnusableIP(ip net.IP) bool {
	if ip == nil {
		return true
	}
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsUnspecified()
}

// filterUnusable 剔除内网/保留 IP，返回其余 IP。输入切片不被修改。
func filterUnusable(ips []string) []string {
	out := make([]string, 0, len(ips))
	for _, s := range ips {
		ip := net.ParseIP(strings.TrimSpace(s))
		if ip == nil || isUnusableIP(ip) {
			continue
		}
		out = append(out, s)
	}
	return out
}
