package main

import (
	"sync"
	"sync/atomic"
)

func incCommonVar(target *uint64) {
	atomic.AddUint64(target, 1)
}

func subCommonVar(target *uint64) {
	atomic.AddUint64(target, ^uint64(0))
}

func setBoolCommonVar(target *bool, varMu *sync.Mutex, value bool) {
	varMu.Lock()
	defer varMu.Unlock()
	*target = value
}
