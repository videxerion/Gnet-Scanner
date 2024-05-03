package main

import "sync"

func incCommonVar(target *uint64, varMu *sync.Mutex) {
	varMu.Lock()
	defer varMu.Unlock()
	*target += 1
}

func subCommonVar(target *uint64, varMu *sync.Mutex) {
	varMu.Lock()
	defer varMu.Unlock()
	*target -= 1
}

func setBoolCommonVar(target *bool, varMu *sync.Mutex, value bool) {
	varMu.Lock()
	defer varMu.Unlock()
	*target = value
}
