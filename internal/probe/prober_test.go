package probe

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"
)

// seqDial 依次返回给定延迟的拨号函数（非并发安全，仅用于单 goroutine 测试）。
func seqDial(delays []time.Duration) DialFunc {
	var i int
	return func(ctx context.Context, _ string, _ string) (net.Conn, error) {
		d := delays[i%len(delays)]
		i++
		select {
		case <-time.After(d):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
		c1, _ := net.Pipe()
		return c1, nil
	}
}

// addrDial 按 addr 返回固定延迟的拨号函数（并发安全，仅读 map）。
func addrDial(delays map[string]time.Duration) DialFunc {
	return func(ctx context.Context, _ string, addr string) (net.Conn, error) {
		d, ok := delays[addr]
		if !ok {
			return nil, errors.New("unexpected addr " + addr)
		}
		select {
		case <-time.After(d):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
		c1, _ := net.Pipe()
		return c1, nil
	}
}

func TestTCPProberMedian(t *testing.T) {
	p := NewTCPProber(
		seqDial([]time.Duration{100 * time.Millisecond, 50 * time.Millisecond, 200 * time.Millisecond}),
		3, time.Second)
	lat, err := p.Latency(context.Background(), "1.2.3.4", 443)
	if err != nil {
		t.Fatalf("Latency() error = %v", err)
	}
	// [50,100,200] 的中位数应落在 100ms 附近
	if lat < 90*time.Millisecond || lat > 110*time.Millisecond {
		t.Errorf("Latency() = %v, want ~100ms", lat)
	}
}

func TestTCPProberTimeout(t *testing.T) {
	dial := func(ctx context.Context, _ string, _ string) (net.Conn, error) {
		select {
		case <-time.After(5 * time.Second):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
		c1, _ := net.Pipe()
		return c1, nil
	}
	p := NewTCPProber(dial, 1, 100*time.Millisecond)
	if _, err := p.Latency(context.Background(), "1.2.3.4", 443); err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestBestPicksLowestLatency(t *testing.T) {
	d := addrDial(map[string]time.Duration{
		"1.2.3.4:443": 150 * time.Millisecond,
		"5.6.7.8:443": 50 * time.Millisecond,
	})
	p := NewTCPProber(d, 1, time.Second)
	res := Best(context.Background(), p, []string{"1.2.3.4", "5.6.7.8"}, 443)
	if !res.OK {
		t.Fatal("expected OK result")
	}
	if res.IP != "5.6.7.8" {
		t.Errorf("Best() IP = %s, want 5.6.7.8", res.IP)
	}
}

func TestBestAllFail(t *testing.T) {
	d := func(_ context.Context, _ string, _ string) (net.Conn, error) {
		return nil, errors.New("connection refused")
	}
	p := NewTCPProber(d, 1, 100*time.Millisecond)
	res := Best(context.Background(), p, []string{"1.2.3.4", "5.6.7.8"}, 443)
	if res.OK {
		t.Error("expected OK=false when all probes fail")
	}
	if res.IP != "1.2.3.4" {
		t.Errorf("Best() IP = %s, want first candidate 1.2.3.4", res.IP)
	}
}

func TestBestEmpty(t *testing.T) {
	p := NewTCPProber(nil, 1, time.Second)
	res := Best(context.Background(), p, nil, 443)
	if res.OK || res.IP != "" {
		t.Errorf("Best() empty = %+v, want zero Result", res)
	}
}

// TestBestContextCancel 验证即使探测不尊重 ctx，Best 也能在取消时尽快返回。
func TestBestContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	// 忽略 ctx 的拨号：永不返回，用于验证 Best 的 ctx.Done 分支。
	d := func(context.Context, string, string) (net.Conn, error) {
		select {}
	}
	p := NewTCPProber(d, 1, time.Second)

	done := make(chan struct{})
	go func() {
		Best(ctx, p, []string{"1.2.3.4", "5.6.7.8"}, 443)
		close(done)
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Best() did not return after ctx cancel")
	}
}

// TestProberInterfaceSatisfied 编译期断言 TCPProber 实现 Prober。
func TestProberInterfaceSatisfied(t *testing.T) {
	var _ Prober = (*TCPProber)(nil)
}
