package ucache

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	s := NewSetExpired(time.Second*5, time.Second*2)
	s.Add("3s")
	if s.Has("3s") {
		t.Log("包含3s")
	}
	time.Sleep(time.Second * 10)
	// t.Skip
	if s.Has("3s") {
		t.Error("怎么没有删除3s")
	} else {
		t.Log("已经删除3s")
	}
}
