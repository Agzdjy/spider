package util

import (
	"net/http"
	"net/url"

	"time"

	"github.com/Agzdjy/proxy-pool"
)

//func init() {
//
//	proxypool.InitData("../config/config.json")
//}

func PorxyGet(requestUrl string) (resp *http.Response) {
	proxypool.InitData("../config/config.json")
	proxyIp := proxypool.Range("http")
	proxy := proxyIp.Protocol + "://" + proxyIp.Address + ":" + proxyIp.Port

	proxyUrl, _ := url.Parse(proxy)

	tr := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Millisecond * 1000,
	}

	req, _ := http.NewRequest("GET", requestUrl, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")

	resq, err := client.Do(req)
	if err != nil {
		return PorxyGet(requestUrl)
	}

	return resq
}
