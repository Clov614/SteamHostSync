// SteamHostSync 生成针对 Steam/GitHub/Docker/GOG/Ubisoft 的 hosts 文件。
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/SteamHostSync/internal/app"
)

func main() {
	var (
		cfgPath    = flag.String("config", app.DefaultConfigPath, "config.yaml 路径（不存在则写入默认配置）")
		outDir     = flag.String("out", app.DefaultOutputDir, "产物输出目录")
		readmeTmpl = flag.String("readme", "README_TEMP.md", "README 模板路径（空则跳过 README 生成）")
	)
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, *cfgPath, *outDir, *readmeTmpl); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
