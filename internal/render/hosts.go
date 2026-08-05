package render

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Clov614/SteamHostSync/internal/fileio"
)

// HostTargetMarker 是 README 模板中的 hosts 内容占位符。
const HostTargetMarker = "HOST_TARGET"

var unsafeFilename = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

// sanitizeFilename 将平台名转为安全的文件名字段（小写、非法字符替换为下划线）。
func sanitizeFilename(name string) string {
	return strings.ToLower(unsafeFilename.ReplaceAllString(name, "_"))
}

// combined 将多个平台结果拼接为合并的 Hosts 内容。
func combined(results []Result, version int) string {
	var b strings.Builder
	for _, r := range results {
		b.WriteString(Render(r, version))
	}
	fmt.Fprintf(&b, "# Github: %s\n", ProjectURL)
	return b.String()
}

// RenderReadme 将生成的 hosts 内容填入 README 模板。
func RenderReadme(tmpl, hosts string) string {
	return strings.Replace(tmpl, HostTargetMarker, hosts, 1)
}

// WriteAll 渲染并写出每个平台的 Hosts_<name> 文件，以及合并的 Hosts 文件。
// readmeTmpl 非空时同时渲染出 README.md。所有写出均为原子写。
// 两个平台消毒后映射到同一文件名时返回错误，避免静默覆盖。
func WriteAll(results []Result, dir string, version int, readmeTmpl string) error {
	hosts := combined(results, version)

	seen := make(map[string]string, len(results)) // 输出文件名 -> 平台名
	for _, r := range results {
		fname := "Hosts_" + sanitizeFilename(r.Platform)
		if prev, dup := seen[fname]; dup {
			return fmt.Errorf("platform %q and %q collide on output file %s", prev, r.Platform, fname)
		}
		seen[fname] = r.Platform
	}

	for _, r := range results {
		content := Render(r, version)
		fname := "Hosts_" + sanitizeFilename(r.Platform)
		if err := fileio.AtomicWrite(filepath.Join(dir, fname), []byte(content)); err != nil {
			return fmt.Errorf("write %s: %w", fname, err)
		}
	}

	if err := fileio.AtomicWrite(filepath.Join(dir, "Hosts"), []byte(hosts)); err != nil {
		return fmt.Errorf("write Hosts: %w", err)
	}

	if readmeTmpl != "" {
		readme := RenderReadme(readmeTmpl, hosts)
		if err := fileio.AtomicWrite(filepath.Join(dir, "README.md"), []byte(readme)); err != nil {
			return fmt.Errorf("write README.md: %w", err)
		}
	}
	return nil
}
