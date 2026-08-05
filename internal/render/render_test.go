package render

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func fixedTime() time.Time {
	return time.Date(2026, 8, 5, 12, 0, 0, 0, time.FixedZone("CST", 8*3600))
}

func sampleResult() Result {
	return Result{
		Platform: "github",
		At:       fixedTime(),
		Entries: []Entry{
			{IP: "140.82.114.4", Domain: "github.com", OK: true},
			{IP: "140.82.114.5", Domain: "api.github.com", OK: true},
		},
	}
}

func TestRender(t *testing.T) {
	got := Render(sampleResult(), 1)
	want := `# SteamHostSync hosts v1
# Generated: 2026-08-05T12:00:00+08:00
# Project: https://github.com/Clov614/SteamHostSync

# github Start
140.82.114.4			github.com
140.82.114.5			api.github.com
# github End # Last Update Time : 2026-08-05T12:00:00+08:00
`
	if got != want {
		t.Errorf("Render() mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestRenderFailedEntry(t *testing.T) {
	r := Result{
		Platform: "steam",
		At:       fixedTime(),
		Entries: []Entry{
			{Domain: "store.steampowered.com", OK: false},
			{IP: "1.2.3.4", Domain: "steamcommunity.com", OK: true},
		},
	}
	got := Render(r, 1)
	if !strings.Contains(got, "# store.steampowered.com\n") {
		t.Errorf("failed entry should be a comment line:\n%s", got)
	}
	if !strings.Contains(got, "1.2.3.4\t\t\tsteamcommunity.com\n") {
		t.Errorf("ok entry should be an ip line:\n%s", got)
	}
}

func TestRenderLFOnly(t *testing.T) {
	got := Render(sampleResult(), 1)
	if strings.Contains(got, "\r\n") {
		t.Error("output must use LF line endings, found CRLF")
	}
	if !strings.HasSuffix(got, "\n") {
		t.Error("output should end with a newline")
	}
}

func TestRenderStartEndPair(t *testing.T) {
	got := Render(sampleResult(), 1)
	if !strings.Contains(got, "# github Start\n") {
		t.Error("missing Start marker")
	}
	if !strings.Contains(got, "# github End # Last Update Time :") {
		t.Error("missing End marker")
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := map[string]string{
		"gog galaxy":       "gog_galaxy",
		"Ubisoft_download": "ubisoft_download",
		"github":           "github",
		"a b c/../x":       "a_b_c_.._x",
	}
	for in, want := range tests {
		if got := sanitizeFilename(in); got != want {
			t.Errorf("sanitizeFilename(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestRenderReadme(t *testing.T) {
	tmpl := "header\nHOST_TARGET\nfooter\n"
	hosts := "# github Start\n..."
	got := RenderReadme(tmpl, hosts)
	if !strings.Contains(got, hosts) {
		t.Errorf("RenderReadme should embed hosts, got:\n%s", got)
	}
	if strings.Contains(got, "HOST_TARGET") {
		t.Error("placeholder should be replaced exactly once")
	}
}

func TestWriteAll(t *testing.T) {
	dir := t.TempDir()
	results := []Result{
		sampleResult(),
		{
			Platform: "steam",
			At:       fixedTime(),
			Entries:  []Entry{{IP: "1.2.3.4", Domain: "store.steampowered.com", OK: true}},
		},
	}
	tmpl := "README_TEMP\nHOST_TARGET\n"

	if err := WriteAll(results, dir, 1, tmpl); err != nil {
		t.Fatalf("WriteAll() error = %v", err)
	}

	// 平台文件
	for _, want := range []string{"Hosts_github", "Hosts_steam"} {
		data, err := os.ReadFile(filepath.Join(dir, want))
		if err != nil {
			t.Fatalf("read %s: %v", want, err)
		}
		if strings.Contains(string(data), "\r\n") {
			t.Errorf("%s must use LF endings", want)
		}
	}

	// 合并文件包含所有平台 + 项目尾部
	combined, err := os.ReadFile(filepath.Join(dir, "Hosts"))
	if err != nil {
		t.Fatalf("read Hosts: %v", err)
	}
	cs := string(combined)
	for _, sub := range []string{"# github Start", "# steam Start", "# Github: https://github.com/Clov614/SteamHostSync"} {
		if !strings.Contains(cs, sub) {
			t.Errorf("combined Hosts missing %q", sub)
		}
	}

	// README 渲染
	readme, err := os.ReadFile(filepath.Join(dir, "README.md"))
	if err != nil {
		t.Fatalf("read README.md: %v", err)
	}
	if strings.Contains(string(readme), "HOST_TARGET") {
		t.Error("README.md should not contain HOST_TARGET")
	}

	// 不应残留临时文件
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".hosts-tmp-") {
			t.Errorf("leftover temp file: %s", e.Name())
		}
	}
}

func TestWriteAllMissingDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "missing")
	if err := WriteAll([]Result{sampleResult()}, dir, 1, ""); err == nil {
		t.Fatal("expected error writing into missing dir")
	}
}

func TestWriteAllAllFailStillWrites(t *testing.T) {
	dir := t.TempDir()
	results := []Result{
		{
			Platform: "docker",
			At:       fixedTime(),
			Entries:  []Entry{{Domain: "docker.com", OK: false}},
		},
	}
	if err := WriteAll(results, dir, 1, ""); err != nil {
		t.Fatalf("WriteAll() error = %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "Hosts_docker"))
	if err != nil {
		t.Fatalf("read Hosts_docker: %v", err)
	}
	if !strings.Contains(string(data), "# docker.com\n") {
		t.Errorf("expected commented failed domain:\n%s", data)
	}
}

// TestWriteAllFilenameCollision 验证消毒后同名平台返回错误而非静默覆盖。
func TestWriteAllFilenameCollision(t *testing.T) {
	dir := t.TempDir()
	results := []Result{
		{Platform: "gog", At: fixedTime(), Entries: []Entry{{IP: "1.1.1.1", Domain: "gog.com", OK: true}}},
		{Platform: "GOG", At: fixedTime(), Entries: []Entry{{IP: "2.2.2.2", Domain: "gog.com", OK: true}}},
	}
	if err := WriteAll(results, dir, 1, ""); err == nil {
		t.Fatal("expected collision error for platforms 'gog' and 'GOG'")
	}
}
