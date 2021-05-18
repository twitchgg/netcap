package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"anyun.bitbucket.com/netcap/internal/api"
)

func TestGenNGrepParams(t *testing.T) {
	params := api.GenNgrepParams("", "ali|baidu", []string{"10.200.200.1", "192.168.0.1"}, "./dump.pcap")
	fmt.Println(params)
}

func TestFileSize(t *testing.T) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://127.0.0.1:8081/cap/stat", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("post a data successful.")
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("response data:%v\n", string(respBody))
}
func TestHTTPClient(t *testing.T) {
	client := &http.Client{}
	form := url.Values{}
	form.Add("ip_begin", "172.28.4.237")
	form.Add("ip_end", "172.28.4.238")
	form.Add("stop_time", "2021-05-18 15:18:00")
	form.Add("dump_size", "500")

	req, err := http.NewRequest("POST", "http://127.0.0.1:8081/cap/ip", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("post a data successful.")
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("response data:%v\n", string(respBody))
}
func TestKeyword(t *testing.T) {
	client := &http.Client{}
	form := url.Values{}
	form.Add("width", "dynamic")
	form.Add("keyword", "baidu.com")
	form.Add("stop_time", "2021-05-18 02:58:00")
	form.Add("dump_size", "500")

	req, err := http.NewRequest("POST", "http://127.0.0.1:8081/cap/keyword", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("post a data successful.")
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("response data:%v\n", string(respBody))
}

func TestIPRange(t *testing.T) {
	api.GetIpsFromRange("172.28.4.235", "172.28.4.240")
}

func TestTimeFormat(t *testing.T) {
	format := "2006-01-02 15:04:05"
	l, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatal(err)
	}
	t1, err := time.ParseInLocation(format, "2021-05-18 01:35:00", l)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(t1)
}
