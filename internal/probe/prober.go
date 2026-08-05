// Package probe 对候选 IP 做 TCP 连接测速，为 DNS 解析结果选优。
package probe

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strconv"
	"time"
)

// Prober 抽象"测量单个 IP 的连接延迟"。
type Prober interface {
	Latency(ctx context.Context, ip string, port int) (time.Duration, error)
}

// DialFunc 抽象 TCP 连接函数，便于测试注入。
type DialFunc func(ctx context.Context, network, addr string) (net.Conn, error)

// TCPProber 通过多次 TCP 握手测速并取中位数，抵抗偶发抖动。
type TCPProber struct {
	dial     DialFunc
	attempts int
	timeout  time.Duration
}

// NewTCPProber 构造 TCP 测速器。dial 为 nil 时使用标准 net.Dialer。
func NewTCPProber(dial DialFunc, attempts int, timeout time.Duration) *TCPProber {
	if dial == nil {
		dial = (&net.Dialer{}).DialContext
	}
	if attempts < 1 {
		attempts = 1
	}
	return &TCPProber{dial: dial, attempts: attempts, timeout: timeout}
}

// Latency 测量到 ip:port 的 TCP 连接延迟，返回 attempts 次握手中位数。
func (p *TCPProber) Latency(ctx context.Context, ip string, port int) (time.Duration, error) {
	addr := net.JoinHostPort(ip, strconv.Itoa(port))
	samples := make([]time.Duration, 0, p.attempts)
	for i := 0; i < p.attempts; i++ {
		dctx, cancel := context.WithTimeout(ctx, p.timeout)
		start := time.Now()
		conn, err := p.dial(dctx, "tcp", addr)
		cancel()
		if err != nil {
			return 0, fmt.Errorf("probe %s: %w", addr, err)
		}
		_ = conn.Close()
		samples = append(samples, time.Since(start))
	}
	sort.Slice(samples, func(i, j int) bool { return samples[i] < samples[j] })
	return samples[len(samples)/2], nil
}
