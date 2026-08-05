package resolve

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestDohQueryFiltersNonARecords 验证仅保留 type==1 的 A 记录。
func TestDohQueryFiltersNonARecords(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/dns-json")
		_, _ = w.Write([]byte(`{
			"Answer": [
				{"name": "example.com", "type": 5, "TTL": 60, "data": "example.com.cdn."},
				{"name": "example.com", "type": 1, "TTL": 60, "data": "1.2.3.4"},
				{"name": "example.com", "type": 1, "TTL": 60, "data": "5.6.7.8"},
				{"name": "example.com", "type": 28, "TTL": 60, "data": "2001:db8::1"}
			]
		}`))
	}))
	defer srv.Close()

	ips, err := dohQuery(context.Background(), srv.Client(), srv.URL, "example.com")
	if err != nil {
		t.Fatalf("dohQuery() error = %v", err)
	}
	want := 2
	if len(ips) != want {
		t.Fatalf("dohQuery() returned %d IPs, want %d: %v", len(ips), want, ips)
	}
	// 应包含两个 A 记录
	if ips[0] != "1.2.3.4" || ips[1] != "5.6.7.8" {
		t.Errorf("unexpected IPs: %v", ips)
	}
}

func TestDohQueryNon200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer srv.Close()

	if _, err := dohQuery(context.Background(), srv.Client(), srv.URL, "example.com"); err == nil {
		t.Fatal("expected error for non-200 response")
	}
}

func TestDohQueryMalformedJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/dns-json")
		_, _ = w.Write([]byte(`{not json`))
	}))
	defer srv.Close()

	if _, err := dohQuery(context.Background(), srv.Client(), srv.URL, "example.com"); err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}

// TestDohQueryNXDOMAIN 验证 NXDOMAIN(RCODE=3) 是合法空结果而非错误。
func TestDohQueryNXDOMAIN(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/dns-json")
		_, _ = w.Write([]byte(`{"Status":3,"Answer":[]}`))
	}))
	defer srv.Close()

	ips, err := dohQuery(context.Background(), srv.Client(), srv.URL, "nonexistent.example")
	if err != nil {
		t.Fatalf("NXDOMAIN should be a valid empty result, got error: %v", err)
	}
	if len(ips) != 0 {
		t.Errorf("ips = %v, want empty", ips)
	}
}

// TestDohQueryServfail 验证 SERVFAIL(RCODE=2) 视为上游失败返回错误。
func TestDohQueryServfail(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/dns-json")
		_, _ = w.Write([]byte(`{"Status":2,"Answer":[]}`))
	}))
	defer srv.Close()

	if _, err := dohQuery(context.Background(), srv.Client(), srv.URL, "example.com"); err == nil {
		t.Fatal("expected error for RCODE SERVFAIL")
	}
}

func TestFilterUnusable(t *testing.T) {
	in := []string{
		"1.2.3.4",
		"0.0.0.0",     // unspecified
		"127.0.0.1",   // loopback
		"10.1.2.3",    // private
		"172.16.0.1",  // private
		"192.168.1.1", // private
		"169.254.0.1", // link-local
		"not-an-ip",   // invalid
		"8.8.8.8",
	}
	got := filterUnusable(in)
	want := []string{"1.2.3.4", "8.8.8.8"}
	if len(got) != len(want) {
		t.Fatalf("filterUnusable() = %v, want %v", got, want)
	}
	for i := range want {
		if !hasStr(got, want[i]) {
			t.Errorf("filterUnusable() missing %q, got %v", want[i], got)
		}
	}
}

// TestMultiResolverPartialFailure 验证部分上游失败仍返回成功上游的 IP。
func TestMultiResolverPartialFailure(t *testing.T) {
	good := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/dns-json")
		_, _ = w.Write([]byte(`{"Answer":[{"type":1,"data":"1.2.3.4"}]}`))
	}))
	defer good.Close()
	bad := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "down", http.StatusBadGateway)
	}))
	defer bad.Close()

	r := NewMultiResolver([]string{bad.URL, good.URL}, nil, 3*time.Second)
	ips, err := r.LookupA(context.Background(), "example.com")
	if err != nil {
		t.Fatalf("LookupA() error = %v", err)
	}
	if !hasStr(ips, "1.2.3.4") {
		t.Errorf("LookupA() = %v, want to include 1.2.3.4", ips)
	}
}

// TestMultiResolverAllFailure 验证全部上游失败返回 error。
func TestMultiResolverAllFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "down", http.StatusBadGateway)
	}))
	defer srv.Close()

	r := NewMultiResolver([]string{srv.URL, srv.URL}, nil, 3*time.Second)
	if _, err := r.LookupA(context.Background(), "example.com"); err == nil {
		t.Fatal("expected error when all upstreams fail")
	}
}

func TestMultiResolverEmptyServers(t *testing.T) {
	r := NewMultiResolver(nil, nil, time.Second)
	if _, err := r.LookupA(context.Background(), "example.com"); err == nil {
		t.Fatal("expected error when no servers configured")
	}
}

// TestMultiResolverTimeout 验证超时会取消并发查询并返回 error。
func TestMultiResolverTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			return
		case <-time.After(5 * time.Second):
			w.Header().Set("Content-Type", "application/dns-json")
			_, _ = w.Write([]byte(`{"Answer":[{"type":1,"data":"1.2.3.4"}]}`))
		}
	}))
	defer srv.Close()

	r := NewMultiResolver([]string{srv.URL}, nil, 100*time.Millisecond)
	_, err := r.LookupA(context.Background(), "example.com")
	if err == nil {
		t.Fatal("expected timeout error")
	}
}

// TestMultiResolverDedupAndSorted 验证多上游返回相同 IP 时去重，且输出按字典序稳定。
func TestMultiResolverDedupAndSorted(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/dns-json")
		_, _ = w.Write([]byte(`{"Answer":[{"type":1,"data":"5.6.7.8"},{"type":1,"data":"1.2.3.4"}]}`))
	}))
	defer srv.Close()

	r := NewMultiResolver([]string{srv.URL, srv.URL}, nil, 3*time.Second)
	ips, err := r.LookupA(context.Background(), "example.com")
	if err != nil {
		t.Fatalf("LookupA() error = %v", err)
	}
	// 两个上游返回相同集合 → 去重后仅 2 个，且排序稳定
	if len(ips) != 2 || ips[0] != "1.2.3.4" || ips[1] != "5.6.7.8" {
		t.Errorf("LookupA() = %v, want sorted deduped [1.2.3.4 5.6.7.8]", ips)
	}
}

// TestResolverInterfaceSatisfied 编译期断言 MultiResolver 实现 Resolver。
func TestResolverInterfaceSatisfied(t *testing.T) {
	var _ Resolver = (*MultiResolver)(nil)
}

func hasStr(s []string, target string) bool {
	for _, v := range s {
		if v == target {
			return true
		}
	}
	return false
}
