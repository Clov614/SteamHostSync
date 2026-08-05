package fileio

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAtomicWriteCreates(t *testing.T) {
	path := filepath.Join(t.TempDir(), "Hosts_github")
	if err := AtomicWrite(path, []byte("content\n")); err != nil {
		t.Fatalf("AtomicWrite() error = %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(data) != "content\n" {
		t.Errorf("content = %q, want %q", data, "content\n")
	}
}

func TestAtomicWriteOverwrites(t *testing.T) {
	path := filepath.Join(t.TempDir(), "Hosts")
	if err := AtomicWrite(path, []byte("old")); err != nil {
		t.Fatalf("first write: %v", err)
	}
	if err := AtomicWrite(path, []byte("new")); err != nil {
		t.Fatalf("second write: %v", err)
	}
	data, _ := os.ReadFile(path)
	if string(data) != "new" {
		t.Errorf("content = %q, want %q", data, "new")
	}
}

func TestAtomicWriteNoTempLeftover(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "Hosts")
	if err := AtomicWrite(path, []byte("x")); err != nil {
		t.Fatalf("AtomicWrite() error = %v", err)
	}
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".hosts-tmp-") {
			t.Errorf("leftover temp file: %s", e.Name())
		}
	}
}

func TestAtomicWriteMissingDir(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing", "Hosts")
	if err := AtomicWrite(path, []byte("x")); err == nil {
		t.Fatal("expected error when parent dir does not exist")
	}
}

// TestAtomicWriteRenameFails 覆盖目标为不可替换对象（非空目录）时
// Rename 失败返回错误、原目标保留、临时文件被清理的路径。
func TestAtomicWriteRenameFails(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "Hosts")
	if err := os.Mkdir(path, 0755); err != nil {
		t.Fatal(err)
	}
	// 放入文件使目录非空，令 Rename 文件覆盖目录失败。
	if err := os.WriteFile(filepath.Join(path, "keep.txt"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := AtomicWrite(path, []byte("x")); err == nil {
		t.Fatal("expected error when target is a non-empty directory")
	}
	// 临时文件应被清理
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".hosts-tmp-") {
			t.Errorf("leftover temp file: %s", e.Name())
		}
	}
}
