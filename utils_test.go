package weibo

import "testing"

func TestTerminalOpen(t *testing.T) {
	if err := terminalOpen("example/pic.jpg"); err != nil {
		t.Error(err)
	}
}
