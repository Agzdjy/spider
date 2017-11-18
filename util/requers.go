package util

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"runtime"

	"path"

	"github.com/Agzdjy/proxy-pool"
)

func init() {
	proxypool.InitData(path.Join(getCurrentDir(), "../config/config.json"))
}

func getCurrentDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Join(path.Dir(filename))
}

func PorxyGet(requestUrl string) (res *http.Response) {
	proxyIp := proxypool.Range("http")
	proxy := proxyIp.Protocol + "://" + proxyIp.Address + ":" + proxyIp.Port
	if !proxypool.Check(proxy) {
		//fmt.Println("check error")
		return PorxyGet(requestUrl)
	}

	proxyUrl, _ := url.Parse(proxy)
	tr := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 5,
	}

	req, _ := http.NewRequest("GET", requestUrl, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36")

	res, err := client.Do(req)
	if err != nil || res == nil {
		return PorxyGet(requestUrl)
	}
	fmt.Println("success--->", requestUrl)

	return res
}
