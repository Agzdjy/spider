package util

import (
	"fmt"
	"net/http"
	"net/url"

	"time"

	"github.com/Agzdjy/proxy-pool"
)

func init() {

	proxypool.InitData("./config/config.json")
}

func PorxyGet(requestUrl string) (resp *http.Response) {
	proxyIp := proxypool.Range("http")
	proxy := proxyIp.Protocol + "://" + proxyIp.Address + ":" + proxyIp.Port
	if !proxypool.Check(proxy) {
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

	resq, err := client.Do(req)
	if err != nil {
		fmt.Println("--->", err)
		return PorxyGet(requestUrl)
	}

	return resq
}
