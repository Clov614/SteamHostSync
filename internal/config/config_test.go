package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestParseValid(t *testing.T) {
	data := []byte(`version: 1
concurrency: 8
timeout:
  resolve: 5s
  probe: 2s
probe:
  port: 443
  attempts: 3
dns_servers:
  - https://dns.alidns.com/resolve
platforms:
  - name: github
    domains: [github.com, api.github.com]
`)
	cfg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if cfg.Concurrency != 8 {
		t.Errorf("Concurrency = %d, want 8", cfg.Concurrency)
	}
	if cfg.Timeout.Resolve.Value() != 5*time.Second {
		t.Errorf("Resolve timeout = %v, want 5s", cfg.Timeout.Resolve.Value())
	}
	if cfg.Probe.Attempts != 3 {
		t.Errorf("Probe.Attempts = %d, want 3", cfg.Probe.Attempts)
	}
	if len(cfg.Platforms) != 1 || cfg.Platforms[0].Name != "github" {
		t.Errorf("Platforms = %+v, want one github platform", cfg.Platforms)
	}
}

func TestLoadCreatesDefaultWhenMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, DefaultName)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("test precondition: %s should not exist", path)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(cfg.Platforms) < 2 {
		t.Errorf("default config should contain multiple platforms, got %d", len(cfg.Platforms))
	}
	// 文件应已被写入
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected default config file created: %v", err)
	}
}

func TestParseErrors(t *testing.T) {
	// validBase 是一份完整合法配置，各用例在其基础上改动单个字段。
	const validBase = `version: 1
concurrency: 8
timeout:
  resolve: 5s
  probe: 2s
probe:
  port: 443
  attempts: 3
dns_servers:
  - https://dns.alidns.com/resolve
platforms:
  - name: github
    domains: [github.com]
`

	tests := []struct {
		name string
		data string
		want string
	}{
		{
			name: "unsupported version",
			data: "version: 99\nconcurrency: 8\ndns_servers: [x]\nplatforms: [{name: a, domains: [x]}]\n",
			want: "unsupported config version",
		},
		{
			name: "missing version field",
			data: "concurrency: 8\ndns_servers: [x]\nplatforms: [{name: a, domains: [x]}]\n",
			want: "unsupported config version",
		},
		{
			name: "invalid yaml",
			data: "version: 1\n  bad-indent: [\n",
			want: "invalid yaml",
		},
		{
			name: "bad duration",
			data: "version: 1\nconcurrency: 8\ntimeout: {resolve: 5, probe: 2s}\ndns_servers: [x]\nplatforms: [{name: a, domains: [x]}]\n",
			want: "invalid duration",
		},
		{
			name: "concurrency zero",
			data: strings.Replace(validBase, "concurrency: 8", "concurrency: 0", 1),
			want: "concurrency must be",
		},
		{
			name: "empty dns_servers",
			data: `version: 1
concurrency: 8
timeout:
  resolve: 5s
  probe: 2s
probe:
  port: 443
  attempts: 3
platforms:
  - name: github
    domains: [github.com]
`,
			want: "dns_servers must not be empty",
		},
		{
			name: "empty platforms",
			data: strings.Replace(validBase, "platforms:\n  - name: github\n    domains: [github.com]\n", "", 1),
			want: "platforms must not be empty",
		},
		{
			name: "duplicate platform name",
			data: strings.Replace(validBase,
				"platforms:\n  - name: github\n    domains: [github.com]",
				"platforms:\n  - name: github\n    domains: [github.com]\n  - name: github\n    domains: [api.github.com]", 1),
			want: "duplicate platform name",
		},
		{
			name: "empty platform name",
			data: strings.Replace(validBase, "name: github", "name: \"  \"", 1),
			want: "empty name",
		},
		{
			name: "platform without domains",
			data: strings.Replace(validBase, "    domains: [github.com]", "", 1),
			want: "no domains",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse([]byte(tt.data))
			if err == nil {
				t.Fatal("Parse() error = nil, want error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error = %q, want it to contain %q", err, tt.want)
			}
		})
	}
}
