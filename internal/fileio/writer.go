// Package fileio 提供 hosts 产物文件的安全写盘能力。
package fileio

import (
	"fmt"
	"os"
	"path/filepath"
)

// AtomicWrite 将 content 原子写入 path：先写同目录临时文件再重命名，
// 避免进程中断留下半截文件。
func AtomicWrite(path string, content []byte) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".hosts-tmp-*")
	if err != nil {
		return fmt.Errorf("create temp file in %s: %w", dir, err)
	}
	tmpName := tmp.Name()
	defer func() {
		if tmpName != "" {
			_ = os.Remove(tmpName)
		}
	}()

	if _, err := tmp.Write(content); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}

	if err := os.Rename(tmpName, path); err != nil {
		// 失败时保留原文件不动（临时文件由 defer 清理），避免数据丢失。
		return fmt.Errorf("rename %s to %s: %w", tmpName, path, err)
	}
	tmpName = ""
	return nil
}
