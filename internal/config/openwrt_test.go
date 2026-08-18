package config

import (
	"os"
	"testing"
	"time"
)

func TestOpenWrtRouterConfigParsesAndIsBounded(t *testing.T) {
	data, err := os.ReadFile("../../contrib/openwrt/files/etc/steamhostsync/router-config.yaml")
	if err != nil {
		t.Fatalf("read OpenWrt router config: %v", err)
	}

	cfg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() OpenWrt router config error = %v", err)
	}
	if cfg.Concurrency != 1 {
		t.Errorf("Concurrency = %d, want 1", cfg.Concurrency)
	}
	if cfg.Timeout.Resolve.Value() != 3*time.Second {
		t.Errorf("Timeout.Resolve = %s, want 3s", cfg.Timeout.Resolve.Value())
	}
	if cfg.Timeout.Probe.Value() != time.Second {
		t.Errorf("Timeout.Probe = %s, want 1s", cfg.Timeout.Probe.Value())
	}
	if cfg.Probe.Port != 443 || cfg.Probe.Attempts != 1 {
		t.Errorf("Probe = %+v, want port 443 and one attempt", cfg.Probe)
	}
	if len(cfg.DNSServers) != 2 {
		t.Errorf("len(DNSServers) = %d, want 2", len(cfg.DNSServers))
	}
	wantPlatforms := []string{"github", "steam", "steam_linux"}
	if len(cfg.Platforms) != len(wantPlatforms) {
		t.Fatalf("len(Platforms) = %d, want %d", len(cfg.Platforms), len(wantPlatforms))
	}
	for i, want := range wantPlatforms {
		if cfg.Platforms[i].Name != want {
			t.Errorf("Platforms[%d].Name = %q, want %q", i, cfg.Platforms[i].Name, want)
		}
		if len(cfg.Platforms[i].Domains) == 0 {
			t.Errorf("Platforms[%d] %q has no domains", i, want)
		}
	}
}
