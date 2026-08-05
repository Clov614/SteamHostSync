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
// ctx 取消时尽快返回（不等待探测完成）。
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
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		mu.Lock()
		final := best
		mu.Unlock()
		return final
	case <-ctx.Done():
		// 取消时返回 DoH 首条 IP 并标记失败，调用方据此走降级；锁定避免与晚完成探测竞争。
		mu.Lock()
		final := best
		final.OK = false
		mu.Unlock()
		return final
	}
}
