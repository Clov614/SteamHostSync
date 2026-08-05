package resolve

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

// Resolver 抽象"给定域名得到解析候选 IP 列表"的解析器。
type Resolver interface {
	LookupA(ctx context.Context, domain string) ([]string, error)
}

// MultiResolver 并发查询多个 DoH 上游，合并去重各上游返回的 A 记录。
// 部分上游失败不影响整体结果；全部失败才返回 error。
type MultiResolver struct {
	servers []string
	client  *http.Client
	timeout time.Duration
}

// NewMultiResolver 构造一个多上游解析器。client 可为 nil（自动使用默认客户端）。
func NewMultiResolver(servers []string, client *http.Client, timeout time.Duration) *MultiResolver {
	if client == nil {
		client = &http.Client{Timeout: timeout}
	}
	return &MultiResolver{servers: servers, client: client, timeout: timeout}
}

// LookupA 并发查询所有上游并合并去重 A 记录（已剔除内网/保留 IP）。
func (m *MultiResolver) LookupA(ctx context.Context, domain string) ([]string, error) {
	if len(m.servers) == 0 {
		return nil, fmt.Errorf("resolve %s: no dns servers configured", domain)
	}

	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	seen := make(map[string]struct{})
	var mu sync.Mutex
	var anySuccess bool
	var errs []string

	var wg sync.WaitGroup
	for _, server := range m.servers {
		wg.Add(1)
		go func(server string) {
			defer wg.Done()
			ips, err := dohQuery(ctx, m.client, server, domain)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errs = append(errs, fmt.Sprintf("%s: %v", server, err))
				return
			}
			anySuccess = true
			for _, ip := range filterUnusable(ips) {
				seen[ip] = struct{}{}
			}
		}(server)
	}
	wg.Wait()

	if !anySuccess {
		return nil, fmt.Errorf("resolve %s: all dns servers failed (%s)", domain, strings.Join(errs, "; "))
	}
	// 部分上游失败时记录错误，便于排查劣质上游。
	if len(errs) > 0 {
		log.Printf("resolve %s: %d/%d servers failed (%s)",
			domain, len(errs), len(m.servers), strings.Join(errs, "; "))
	}

	out := make([]string, 0, len(seen))
	for ip := range seen {
		out = append(out, ip)
	}
	// 排序保证输出确定、可复现（fallback IP 与产物不随 map 迭代漂移）。
	sort.Strings(out)
	return out, nil
}
