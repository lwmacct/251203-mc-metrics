// Package config 提供应用配置管理
//
// 配置加载优先级 (从低到高)：
//  1. 默认值 - DefaultConfig() 函数中定义
//  2. 配置文件 - 通过 --config 指定，或按顺序搜索默认路径
//  3. 环境变量 - 以 <AppRawName> 为前缀，下划线分隔嵌套路径
//  4. CLI flags - 最高优先级
package config

import "time"

// Config 应用配置
type Config struct {
	Server ServerConfig `koanf:"server" comment:"服务器配置"`
	Auth   AuthConfig   `koanf:"auth" comment:"认证配置"`
	TLS    TLSConfig    `koanf:"tls" comment:"TLS 配置"`
	Output OutputConfig `koanf:"output" comment:"输出配置"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	URL        string        `koanf:"url" comment:"VictoriaMetrics 服务器地址"`
	PathPrefix string        `koanf:"path_prefix" comment:"API 路径前缀 (如 /victoria)"`
	Timeout    time.Duration `koanf:"timeout" comment:"请求超时时间"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	Type     string `koanf:"type" comment:"认证类型: basic, bearer"`
	User     string `koanf:"user" comment:"Basic 认证用户名"`
	Password string `koanf:"password" comment:"Basic 认证密码"`
	Token    string `koanf:"token" comment:"Bearer Token"`
}

// TLSConfig TLS 配置
type TLSConfig struct {
	CA         string `koanf:"ca" comment:"CA 证书路径"`
	Cert       string `koanf:"cert" comment:"客户端证书路径"`
	Key        string `koanf:"key" comment:"客户端密钥路径"`
	SkipVerify bool   `koanf:"skip_verify" comment:"跳过证书验证"`
}

// OutputConfig 输出配置
type OutputConfig struct {
	Format    string `koanf:"format" comment:"输出格式: table, json, csv, graph"`
	NoHeaders bool   `koanf:"no_headers" comment:"禁用表头输出"`
}

// DefaultConfig 返回默认配置
// 注意：这里的默认值应对齐 internal/command/*/command.go 中的默认值
func DefaultConfig() Config {
	return Config{
		Server: ServerConfig{
			URL:     "http://localhost:8428",
			Timeout: 30 * time.Second,
		},
		Auth: AuthConfig{
			Type: "",
		},
		TLS: TLSConfig{
			SkipVerify: false,
		},
		Output: OutputConfig{
			Format:    "table",
			NoHeaders: false,
		},
	}
}
