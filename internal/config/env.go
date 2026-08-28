package config

import (
	"fmt"
	"net"
	"path/filepath"
	"strings"
)

type Environment struct {
	Config  Config
	Host    string
	Port    string
	DataDir string
}

func ResolveEnvironment(cfg Config) (Environment, error) {
	if err := cfg.Validate(); err != nil {
		return Environment{}, err
	}
	host, port, err := net.SplitHostPort(cfg.Address)
	if err != nil {
		if strings.HasPrefix(cfg.Address, ":") {
			host = ""
			port = strings.TrimPrefix(cfg.Address, ":")
		} else {
			return Environment{}, fmt.Errorf("invalid listen address: %w", err)
		}
	}
	if port == "" {
		return Environment{}, fmt.Errorf("listen port is required")
	}
	dataDir := filepath.Dir(cfg.Database)
	if dataDir == "." {
		dataDir = "./"
	}
	return Environment{Config: cfg, Host: host, Port: port, DataDir: dataDir}, nil
}

func (e Environment) AddressLabel() string {
	host := e.Host
	if host == "" {
		host = "localhost"
	}
	return net.JoinHostPort(host, e.Port)
}

func (e Environment) IsLocal() bool {
	return e.Host == "" || e.Host == "127.0.0.1" || e.Host == "localhost"
}

func (e Environment) DataPath(name string) string {
	return filepath.Join(e.DataDir, strings.TrimSpace(name))
}

func (e Environment) Labels() []string {
	mode := "远程"
	if e.IsLocal() {
		mode = "本机"
	}
	return []string{e.AddressLabel(), mode, e.DataDir}
}

func SupportedAddresses() []string {
	return []string{":8080", ":9090", "127.0.0.1:8080"}
}

func (e Environment) DatabaseLabel() string {
	if strings.HasPrefix(e.Config.Database, ":") {
		return "内存数据库"
	}
	return filepath.Base(e.Config.Database)
}

func (e Environment) Describe() string {
	return fmt.Sprintf("服务 %s，数据 %s，%s", e.AddressLabel(), e.DatabaseLabel(), e.Config.Summary())
}

func (e Environment) HasCustomDatabase() bool {
	return e.Config.Database != "qin-culture.db" && e.Config.Database != ":memory:"
}

func (e Environment) IsDevelopment() bool {
	return e.Config.DevMode
}

func (e Environment) MaxBodyLabel() string {
	return fmt.Sprintf("请求上限 %d 字节", e.Config.MaxBody)
}

func (e Environment) Ready() bool {
	return e.Port != "" && e.DataDir != "" && e.Config.MaxBody > 0
}

func (e Environment) ModeLabel() string {
	if e.Config.DevMode {
		return "开发模式"
	}
	return "服务模式"
}

func (e Environment) SummaryLines() []string {
	return []string{e.Describe(), e.ModeLabel(), e.MaxBodyLabel(), "数据库目录：" + e.DataDir}
}

func (e Environment) UsesSQLiteFile() bool {
	return !strings.HasPrefix(e.Config.Database, ":")
}

func (e Environment) StorageKind() string {
	if e.UsesSQLiteFile() {
		return "SQLite 文件"
	}
	return "SQLite 内存"
}

func (e Environment) Endpoint(path string) string {
	return "http://" + e.AddressLabel() + "/" + strings.TrimPrefix(path, "/")
}

func (e Environment) HealthEndpoint() string {
	return e.Endpoint("healthz")
}
