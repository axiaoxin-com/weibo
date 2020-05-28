package weibo

import (
	"testing"
)

func TestSearch(t *testing.T) {
	results, err := Search("阿小信大人")
	if err != nil {
		t.Error(err)
	}
	if len(results) == 0 {
		t.Error("results len = 0")
	}
}
