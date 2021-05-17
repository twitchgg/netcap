package ngrep

import (
	"fmt"
	"os/exec"
)

// AppConfig ngrep app 配置
type AppConfig struct {
	NgrepPath string
	Dev       string
	Keyword   string
	Ports     []int
	DumpPath  string
	DumpSize  int
	HostIPS   []string
	HostNets  []string
}

// Check ngrep app 配置检查
func (c *AppConfig) Check() error {
	np, err := exec.LookPath("ngrep")
	if err != nil {
		return fmt.Errorf("ngrep执行程序路径检查错误: %s", err.Error())
	}
	c.NgrepPath = np
	return nil
}
