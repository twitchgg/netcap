package ngrep

import (
	"fmt"
	"os/exec"
)

// Application ngrep app应用程序
type Application struct {
	conf *AppConfig
	cmd  *exec.Cmd
}

// NewApplication 创建ngrep应用程序
func NewApplication(conf *AppConfig) (app *Application, err error) {
	if conf == nil {
		return nil, fmt.Errorf("应用配置未定义")
	}
	if err = conf.Check(); err != nil {
		return nil, fmt.Errorf("配置检查错误: %s", err.Error())
	}
	app = &Application{
		conf: conf,
	}
	return app, nil
}

// Start 后台启动应用程序
func (a *Application) Start() chan error {
	fmt.Println(a.conf.NgrepPath, a.conf.GenNgrepParams())
	a.cmd = exec.Command(a.conf.NgrepPath, a.conf.GenNgrepParams()...)
	errChan := make(chan error, 1)
	a.cmd.Start()
	if err := a.cmd.Wait(); err != nil {
		errChan <- fmt.Errorf("应用程序运行错误: %s", err.Error())
		return errChan
	}
	return errChan
}

// Stop 关闭应用程序
func (a *Application) Stop() error {
	if err := a.cmd.Process.Kill(); err != nil {
		return err
	}
	return nil
}
