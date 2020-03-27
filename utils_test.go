package weibo

import "testing"

func _TestTerminalOpen(t *testing.T) {
	if err := TerminalOpen("example/pic.jpg"); err != nil {
		t.Error(err)
	}
}
