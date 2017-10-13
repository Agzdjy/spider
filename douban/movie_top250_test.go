package douban

import (
	"testing"
)

func TestRun(t *testing.T) {
	err := Run()

	if err != nil {
		t.Error("douban movie top 250 error", err)
	}
}
