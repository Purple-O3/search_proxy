package log

import (
	"sync"
	"testing"
)

func TestAll(t *testing.T) {
	t.Log("start")
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go writelog(i, &wg, t)
	}
	wg.Wait()
	CloseLogger()
	t.Log("stop")
}

func writelog(gonub int, wg *sync.WaitGroup, t *testing.T) {
	type s struct {
		nub int
	}
	ns := new(s)
	for i := 0; i < 1000; i++ {
		Debugf("gonub:%d, t:%v", gonub, ns)
		Infof("gonub:%d, t:%v", gonub, ns)
		Warnf("gonub:%d, t:%v", gonub, ns)
		Errorf("gonub:%d, t:%v", gonub, ns)
	}
	t.Log("gonub:", gonub)
	wg.Done()
}
