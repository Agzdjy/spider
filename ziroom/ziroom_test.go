package ziroom

import (
	"testing"
)

func TestRun(t *testing.T) {
	ziroom := &Ziroom{}
	err := ziroom.Run()

	if err != nil {
		t.Error("ziroom error", err)
	}
}
