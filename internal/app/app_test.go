package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Clov614/SteamHostSync/internal/config"
)

// ---- fakes ----

// fakeResolver 用静态映射提供解析候选 IP。
type fakeResolver map[string][]string

func (f fakeResolver) LookupA(_ context.Context, domain string) ([]string, error) {
	ips, ok := f[domain]
	if !ok {
		return nil, fmt.Errorf("no such domain %s", domain)
	}
	return ips, nil
}

// okProber 恒定成功探测。
type okProber struct{}

func (okProber) Latency(_ context.Context, _ string, _ int) (time.Duration, error) {
	return 5 * time.Millisecond, nil
}

// failProber 恒定失败探测，用于触发降级。
type failProber struct{}

func (failProber) Latency(_ context.Context, _ string, _ int) (time.Duration, error) {
	return 0, errors.New("connection refused")
}

// blockingResolver 阻塞直到 ctx 取消，用于验证取消传播。
type blockingResolver struct{}

func (blockingResolver) LookupA(ctx context.Context, _ string) ([]string, error) {
	<-ctx.Done()
	return nil, ctx.Err()
}

// newTestConfig 构造一个含给定平台的合法配置。
func newTestConfig() *config.Config {
	return &config.Config{
		Version:     1,
		Concurrency: 4,
		Timeout:     config.TimeoutConfig{Resolve: config.Duration(time.Second), Probe: config.Duration(time.Second)},
		Probe:       config.ProbeConfig{Port: 443, Attempts: 1},
		DNSServers:  []string{"https://dns.google/resolve"},
		Platforms: []config.Platform{
			{Name: "test", Domains: []string{"alpha.example", "beta.example"}},
			{Name: "second", Domains: []string{"gamma.example"}},
		},
	}
}

func TestRunEndToEnd(t *testing.T) {
	dir := t.TempDir()
	tmpl := filepath.Join(dir, "README_TEMP.md")
	if err := os.WriteFile(tmpl, []byte("TMPL\nHOST_TARGET\n"), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := newTestConfig()
	resolver := fakeResolver{
		"alpha.example": {"1.1.1.1", "2.2.2.2"},
		"beta.example":  {"3.3.3.3"},
		"gamma.example": {"4.4.4.4"},
	}

	if err := run(context.Background(), cfg, resolver, okProber{}, dir, tmpl); err != nil {
		t.Fatalf("run() error = %v", err)
	}

	for _, name := range []string{"Hosts_test", "Hosts_second", "Hosts", "README.md"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Errorf("missing output %s: %v", name, err)
		}
	}

	// README 占位符应已被替换
	readme, _ := os.ReadFile(filepath.Join(dir, "README.md"))
	if strings.Contains(string(readme), "HOST_TARGET") {
		t.Error("README.md should not contain HOST_TARGET")
	}
}

func TestRunAllFailReturnsError(t *testing.T) {
	dir := t.TempDir()
	cfg := newTestConfig()
	resolver := fakeResolver{} // 所有域名都无法解析

	if err := run(context.Background(), cfg, resolver, okProber{}, dir, ""); err == nil {
		t.Fatal("expected error when all platforms fail to resolve")
	}
}

func TestRunSingleDomainFailureDoesNotBlock(t *testing.T) {
	dir := t.TempDir()
	cfg := newTestConfig()
	// alpha 失败，beta/gamma 成功
	resolver := fakeResolver{
		"beta.example":  {"3.3.3.3"},
		"gamma.example": {"4.4.4.4"},
	}

	if err := run(context.Background(), cfg, resolver, okProber{}, dir, ""); err != nil {
		t.Fatalf("run() error = %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "Hosts_test"))
	if err != nil {
		t.Fatalf("read Hosts_test: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "# alpha.example\n") {
		t.Error("failed domain should render as comment:\n" + content)
	}
	if !strings.Contains(content, "3.3.3.3\t\t\tbeta.example\n") {
		t.Error("ok domain should render as IP line:\n" + content)
	}
}

func TestRunProbeDegraded(t *testing.T) {
	dir := t.TempDir()
	cfg := newTestConfig()
	resolver := fakeResolver{
		"alpha.example": {"1.1.1.1", "2.2.2.2"},
		"beta.example":  {"3.3.3.3"},
		"gamma.example": {"4.4.4.4"},
	}
	// prober 全失败 → 仍用 DoH 首条 IP 并标注 # probe-failed
	if err := run(context.Background(), cfg, resolver, failProber{}, dir, ""); err != nil {
		t.Fatalf("run() error = %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "Hosts_test"))
	if err != nil {
		t.Fatalf("read Hosts_test: %v", err)
	}
	if !strings.Contains(string(data), "1.1.1.1\t\t\talpha.example\n# probe-failed\n") {
		t.Errorf("expected degraded entry with probe-failed note:\n%s", data)
	}
}

// TestRunContextCancelled 验证运行中取消：run 返回 context.Canceled 且不写出任何产物。
func TestRunContextCancelled(t *testing.T) {
	dir := t.TempDir()
	cfg := newTestConfig()
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() {
		done <- run(ctx, cfg, blockingResolver{}, okProber{}, dir, "")
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err == nil || !errors.Is(err, context.Canceled) {
			t.Fatalf("run() error = %v, want context.Canceled", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("run() did not return after cancellation")
	}

	// 不应写出任何产物
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		t.Errorf("unexpected output after cancellation: %s", e.Name())
	}
}

func TestRunWithHTTPServer(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")

	// 本地 DoH 上游：任何域名返回 A 记录 192.0.2.1（TEST-NET，探测必失败 → 降级路径）。
	doh := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/dns-json")
		_, _ = w.Write([]byte(`{"Answer":[{"type":1,"data":"192.0.2.1"}]}`))
	}))
	defer doh.Close()

	cfgYAML := fmt.Sprintf(`version: 1
concurrency: 2
timeout: {resolve: 5s, probe: 200ms}
probe: {port: 443, attempts: 1}
dns_servers: [%s]
platforms:
  - name: smoke
    domains: [alpha.example, beta.example]
`, doh.URL)
	if err := os.WriteFile(cfgPath, []byte(cfgYAML), 0644); err != nil {
		t.Fatal(err)
	}

	if err := Run(context.Background(), cfgPath, dir, ""); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "Hosts_smoke"))
	if err != nil {
		t.Fatalf("read Hosts_smoke: %v", err)
	}
	if !strings.Contains(string(data), "192.0.2.1\t\t\talpha.example\n# probe-failed\n") {
		t.Errorf("expected degraded entry from httptest pipeline:\n%s", data)
	}
}
