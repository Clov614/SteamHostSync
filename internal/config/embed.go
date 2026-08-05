// Package config 负责加载与校验 SteamHostSync 的运行时配置。
package config

import (
	_ "embed"
)

//go:embed default_config.yaml
var defaultConfig []byte

// Default 返回内嵌的默认配置内容（作为新 config.yaml 的单一数据源）。
func Default() []byte {
	return defaultConfig
}

// DefaultName 是运行时配置文件的默认文件名。
const DefaultName = "config.yaml"
