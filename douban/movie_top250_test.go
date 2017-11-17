package douban

import (
	"testing"
)

func TestRun(t *testing.T) {
	douban := &Douban{}
	err := douban.Run()

	if err != nil {
		t.Error("douban movie top 250 error", err)
	}
}
