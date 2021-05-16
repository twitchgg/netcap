package test

import (
	"fmt"
	"strings"
	"testing"

	"anyun.bitbucket.com/netcap/pkg/ngrep"
)

func TestGenNGrepParams(t *testing.T) {
	conf := &ngrep.AppConfig{
		Dev:     "eth0",
		Keyword: "ali|baidu",
		Ports:   []int{52, 53, 54},
	}
	fmt.Println(strings.Join(conf.GenNgrepParams(), " "))
}
