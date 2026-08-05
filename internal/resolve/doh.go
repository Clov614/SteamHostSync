package resolve

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// dohAnswer 是 DNS-over-HTTPS JSON 应答中单条记录的结构。
type dohAnswer struct {
	Type int    `json:"type"`
	Data string `json:"data"`
}

// dohResponse 映射 DoH JSON API 的响应体。
type dohResponse struct {
	Status int         `json:"Status"` // 0=RCODE 0 (NOERROR), 3=NXDOMAIN, 2=SERVFAIL, 5=REFUSED
	Answer []dohAnswer `json:"Answer"`
}

// maxDoHBody 限制 DoH 响应体大小，防止异常上游回传超大 JSON。
const maxDoHBody = 1 << 20 // 1 MiB

// dohQuery 针对单个 DoH 上游发起一次 A 记录查询，返回过滤后的 A 记录 IP 列表。
// 过滤规则：仅保留 type==1 的 A 记录，其余（CNAME/AAAA 等）丢弃。
func dohQuery(ctx context.Context, client *http.Client, server, domain string) ([]string, error) {
	endpoint := fmt.Sprintf("%s?type=A&name=%s", server, url.QueryEscape(domain))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("build doh request for %s: %w", domain, err)
	}
	req.Header.Set("Accept", "application/dns-json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("doh query %s for %s: %w", server, domain, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("doh %s for %s returned status %s", server, domain, resp.Status)
	}

	var body dohResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxDoHBody)).Decode(&body); err != nil {
		return nil, fmt.Errorf("decode doh response from %s: %w", server, err)
	}

	// RCODE 语义：0=NOERROR, 3=NXDOMAIN（域名真实不存在，合法空结果），
	// 其余非零值（SERVFAIL=2/REFUSED=5 等）视为上游解析失败，
	// 交由 multi.go 尝试其他上游并记录日志。
	if body.Status != 0 && body.Status != 3 {
		return nil, fmt.Errorf("doh %s for %s returned RCODE %d", server, domain, body.Status)
	}

	var out []string
	for _, a := range body.Answer {
		if a.Type == 1 && a.Data != "" {
			out = append(out, a.Data)
		}
	}
	return out, nil
}
