package probe

import (
	"context"
	"sync"
	"time"
)

// Result 描述对一个 IP 的探测结果。OK=false 表示所有候选探测均失败。
type Result struct {
	IP  string
	Lat time.Duration
	OK  bool
}

// Best 并发探测所有候选 IP，返回延迟最低者。
// 全部失败时返回第一个 IP 且 OK=false，供上层降级（直接使用 DoH 首条记录）。
func Best(ctx context.Context, p Prober, ips []string, port int) Result {
	if len(ips) == 0 {
		return Result{}
	}
	best := Result{IP: ips[0], Lat: time.Duration(1<<63 - 1)}
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, ip := range ips {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			lat, err := p.Latency(ctx, ip, port)
			if err != nil {
				return
			}
			mu.Lock()
			if !best.OK || lat < best.Lat {
				best = Result{IP: ip, Lat: lat, OK: true}
			}
			mu.Unlock()
		}(ip)
	}
	wg.Wait()
	return best
}
