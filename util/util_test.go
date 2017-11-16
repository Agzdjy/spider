package util

import (
	"testing"
)

func TestProxyGet(t *testing.T) {
	resp := PorxyGet("http://www.baidu.com/")

	if resp.StatusCode != 200 {
		t.Error("proxy get resp error", resp.StatusCode)
	}
}
