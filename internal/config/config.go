package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Version 是当前配置文件结构的版本号，用于未来迁移。
const Version = 1

// Duration 是可被 YAML 字符串（如 "5s"）解析的时长包装类型。
type Duration time.Duration

// UnmarshalYAML 支持将 YAML 中的字符串时长解析为 Duration。
func (d *Duration) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if err := node.Decode(&s); err != nil {
		return fmt.Errorf("invalid duration value: %w", err)
	}
	v, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", s, err)
	}
	*d = Duration(v)
	return nil
}

// Value 返回底层 time.Duration。
func (d Duration) Value() time.Duration { return time.Duration(d) }

// TimeoutConfig 配置 DNS 查询与 TCP 探测的超时时间。
type TimeoutConfig struct {
	Resolve Duration `yaml:"resolve"`
	Probe   Duration `yaml:"probe"`
}

// ProbeConfig 配置 TCP 测速的参数。
type ProbeConfig struct {
	Port     int `yaml:"port"`
	Attempts int `yaml:"attempts"`
}

// Platform 描述一个需要生成独立 hosts 文件的平台。
type Platform struct {
	Name    string   `yaml:"name"`
	Domains []string `yaml:"domains"`
}

// Config 是 SteamHostSync 的完整运行配置。
type Config struct {
	Version     int           `yaml:"version"`
	Concurrency int           `yaml:"concurrency"`
	Timeout     TimeoutConfig `yaml:"timeout"`
	Probe       ProbeConfig   `yaml:"probe"`
	DNSServers  []string      `yaml:"dns_servers"`
	Platforms   []Platform    `yaml:"platforms"`
}

// Load 从 path 读取配置。若文件不存在，则写入内嵌默认配置后再次加载。
// 解析或校验失败时返回 error。
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		if werr := os.WriteFile(path, Default(), 0644); werr != nil {
			return nil, fmt.Errorf("create default config %s: %w", path, werr)
		}
		data, err = os.ReadFile(path)
	}
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}

	cfg, err := Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}
	return cfg, nil
}

// Parse 将配置字节解析为 Config 并校验。返回的 Config 非 nil 时必已通过校验。
func Parse(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid yaml: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Validate 校验配置字段，确保配置可被安全使用。
func (c *Config) Validate() error {
	if c.Version != Version {
		return fmt.Errorf("unsupported config version %d (want %d)", c.Version, Version)
	}
	if c.Concurrency < 1 {
		return fmt.Errorf("concurrency must be >= 1, got %d", c.Concurrency)
	}
	if c.Timeout.Resolve <= 0 {
		return errors.New("timeout.resolve must be positive")
	}
	if c.Timeout.Probe <= 0 {
		return errors.New("timeout.probe must be positive")
	}
	if c.Probe.Port <= 0 || c.Probe.Port > 65535 {
		return fmt.Errorf("probe.port out of range: %d", c.Probe.Port)
	}
	if c.Probe.Attempts < 1 {
		return fmt.Errorf("probe.attempts must be >= 1, got %d", c.Probe.Attempts)
	}
	if len(c.DNSServers) == 0 {
		return errors.New("dns_servers must not be empty")
	}
	for _, s := range c.DNSServers {
		if strings.TrimSpace(s) == "" {
			return errors.New("dns_servers contains an empty URL")
		}
		u, err := url.Parse(s)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
			return fmt.Errorf("dns_servers contains invalid URL %q", s)
		}
	}
	if len(c.Platforms) == 0 {
		return errors.New("platforms must not be empty")
	}
	seen := make(map[string]struct{}, len(c.Platforms))
	for i, p := range c.Platforms {
		name := strings.TrimSpace(p.Name)
		if name == "" {
			return fmt.Errorf("platform #%d has an empty name", i+1)
		}
		if _, dup := seen[name]; dup {
			return fmt.Errorf("duplicate platform name %q", name)
		}
		seen[name] = struct{}{}
		if len(p.Domains) == 0 {
			return fmt.Errorf("platform %q has no domains", name)
		}
	}
	return nil
}
