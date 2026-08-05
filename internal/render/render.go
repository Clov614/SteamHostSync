// Package render 将解析/探测结果渲染为 hosts 文件内容。
package render

import (
	"fmt"
	"strings"
	"time"
)

// Entry 是一条 hosts 记录。
type Entry struct {
	IP       string
	Domain   string
	OK       bool // false 表示解析失败，渲染为注释行
	Degraded bool // true 表示探测失败，仅使用 DoH 首条 IP（CI 网络受限时）
}

// Result 是一个平台的全部解析结果。
type Result struct {
	Platform string
	Entries  []Entry
	At       time.Time
}

// ProjectURL 是项目主页，写入每个产物的文件头与合并文件尾部。
const ProjectURL = "https://github.com/Clov614/SteamHostSync"

// header 渲染文件头（版本 + 生成时间 + 项目链接）。
func header(version int, at time.Time) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# SteamHostSync hosts v%d\n", version)
	fmt.Fprintf(&b, "# Generated: %s\n", at.Format(time.RFC3339))
	fmt.Fprintf(&b, "# Project: %s\n\n", ProjectURL)
	return b.String()
}

// Render 将单个平台结果渲染为 hosts 文件内容（纯函数，无副作用）。
func Render(r Result, version int) string {
	var b strings.Builder
	b.WriteString(header(version, r.At))
	fmt.Fprintf(&b, "# %s Start\n", r.Platform)
	for _, e := range r.Entries {
		if e.OK {
			fmt.Fprintf(&b, "%s\t\t\t%s\n", e.IP, e.Domain)
			if e.Degraded {
				b.WriteString("# probe-failed\n")
			}
		} else {
			fmt.Fprintf(&b, "# %s\n", e.Domain)
		}
	}
	fmt.Fprintf(&b, "# %s End # Last Update Time : %s\n", r.Platform, r.At.Format(time.RFC3339))
	return b.String()
}
