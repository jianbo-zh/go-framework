package mmux

import (
	"testing"
)

func TestLock(t *testing.T) {

	key := "hello"
	Lock(key)
	defer Unlock(key)

	key1 := "world"
	Lock(key1)
	defer Unlock(key1)
}
