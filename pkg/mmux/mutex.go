package mmux

import (
	"sync"
	"time"
)

var mMap sync.Map

func Lock(key string) {
	ts := time.Now().UnixNano()
	for {
		_, loaded := mMap.LoadOrStore(key, ts)
		if !loaded {
			break
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func Unlock(key string) {
	mMap.Delete(key)
}
