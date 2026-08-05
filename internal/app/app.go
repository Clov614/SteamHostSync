// Package app 编排完整的"解析→探测→渲染→写盘"流水线，并实现失败降级策略。
package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/SteamHostSync/internal/config"
	"github.com/SteamHostSync/internal/probe"
	"github.com/SteamHostSync/internal/render"
	"github.com/SteamHostSync/internal/resolve"
)

// DefaultConfigPath 与 DefaultOutputDir 是 main 使用的默认路径。
const (
	DefaultConfigPath = "config.yaml"
	DefaultOutputDir  = "."
)

// Run 执行一次完整的 hosts 生成。
//
//	cfgPath    — config.yaml 路径（不存在时自动写入内嵌默认配置）
//	outDir     — 产物输出目录
//	readmeTmpl — README 模板路径；为空则跳过 README 生成
func Run(ctx context.Context, cfgPath, outDir, readmeTmpl string) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	resolver := resolve.NewMultiResolver(cfg.DNSServers, nil, cfg.Timeout.Resolve.Value())
	prober := probe.NewTCPProber(nil, cfg.Probe.Attempts, cfg.Timeout.Probe.Value())
	return run(ctx, cfg, resolver, prober, outDir, readmeTmpl)
}

// run 是流水线核心，接受注入的 resolver/prober 以便测试。
func run(ctx context.Context, cfg *config.Config, resolver resolve.Resolver, prober probe.Prober, outDir, readmeTmpl string) error {
	results, err := gather(ctx, cfg, resolver, prober)
	if err != nil {
		return err
	}

	var tmpl string
	if readmeTmpl != "" {
		data, rerr := os.ReadFile(readmeTmpl)
		if rerr != nil {
			return fmt.Errorf("read readme template %s: %w", readmeTmpl, rerr)
		}
		tmpl = string(data)
	}

	if err := render.WriteAll(results, outDir, cfg.Version, tmpl); err != nil {
		return fmt.Errorf("render hosts: %w", err)
	}
	return nil
}

// gather 并发解析所有平台的所有域名，组装 render.Result 列表。
// 降级策略：
//   - 单域名 DoH 全失败 → 生成注释行，不阻塞
//   - 单平台全部域名失败 → 记录 warn，跳过该平台文件
//   - 全部平台全失败 → 返回 error
func gather(ctx context.Context, cfg *config.Config, resolver resolve.Resolver, prober probe.Prober) ([]render.Result, error) {
	results := make([]render.Result, len(cfg.Platforms))
	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(cfg.Concurrency)

	var mu sync.Mutex
	platformOK := 0

	for pi := range cfg.Platforms {
		pi := pi
		g.Go(func() error {
			entries, ok := resolvePlatform(gctx, cfg.Platforms[pi], cfg.Probe.Port, resolver, prober)
			mu.Lock()
			defer mu.Unlock()
			if ok {
				results[pi] = render.Result{
					Platform: cfg.Platforms[pi].Name,
					Entries:  entries,
					At:       time.Now(),
				}
				platformOK++
			} else {
				log.Printf("warn: platform %q 全部域名失败，跳过该平台文件", cfg.Platforms[pi].Name)
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// 抽出有有效条目的平台，保持配置顺序。
	out := make([]render.Result, 0, len(results))
	for _, r := range results {
		if len(r.Entries) > 0 {
			out = append(out, r)
		}
	}
	if platformOK == 0 {
		return nil, fmt.Errorf("所有平台的全部域名解析失败")
	}
	return out, nil
}

// resolvePlatform 解析单个平台的所有域名。返回条目列表与"是否至少有一个有效条目"。
func resolvePlatform(ctx context.Context, p config.Platform, port int, resolver resolve.Resolver, prober probe.Prober) ([]render.Entry, bool) {
	entries := make([]render.Entry, 0, len(p.Domains))
	for _, raw := range p.Domains {
		dom := strings.TrimSpace(raw)
		if dom == "" {
			continue
		}
		ips, err := resolver.LookupA(ctx, dom)
		if err != nil {
			entries = append(entries, render.Entry{Domain: dom})
			continue
		}
		best := probe.Best(ctx, prober, ips, port)
		e := render.Entry{IP: best.IP, Domain: dom, OK: best.OK}
		if !best.OK && best.IP != "" {
			// 探测失败（如 CI 出站受限）：仍用 DoH 首条 IP，标注降级。
			e.OK = true
			e.Degraded = true
		}
		entries = append(entries, e)
	}

	ok := false
	for _, e := range entries {
		if e.OK {
			ok = true
			break
		}
	}
	return entries, ok
}
