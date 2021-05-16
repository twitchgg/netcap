package ngrep

import (
	"fmt"
	"os/exec"
	"strconv"
)

// AppConfig ngrep app 配置
type AppConfig struct {
	NgrepPath string
	Dev       string
	Keyword   string
	Ports     []int
	DumpPath  string
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

// GenNgrepParams 生成ngrep参数
func (c *AppConfig) GenNgrepParams() []string {
	params := make([]string, 0)
	params = append(params, "-W")
	params = append(params, "byline")
	if c.Dev != "" {
		params = append(params, "-d")
		params = append(params, c.Dev)
	}
	params = append(params, "-i")
	params = append(params, c.Keyword)
	if c.Ports != nil && len(c.Ports) > 0 {
		for _, p := range c.Ports {
			params = append(params, "port")
			params = append(params, strconv.Itoa(p))
			params = append(params, "and")
		}
		params = params[:len(params)-1]
	}
	if c.HostIPS != nil && len(c.HostIPS) > 0 {
		for _, ip := range c.HostIPS {
			params = append(params, "and")
			params = append(params, "dst")
			params = append(params, "host")
			params = append(params, ip)
		}
	}
	if c.HostNets != nil && len(c.HostNets) > 0 {
		for _, n := range c.HostNets {
			params = append(params, "and")
			params = append(params, "dst")
			params = append(params, "net")
			params = append(params, n)
		}
	}
	params = append(params, "-O")
	params = append(params, c.DumpPath)
	return params
}
