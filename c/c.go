package c

// #cgo CFLAGS: -g -Wall -I./mDNSShared
// #include "dns-sd.h"
import "C"
import (
	"sync"
)

func StartDNSService(host string) {
	mu := sync.Mutex{}
	go func() {
		mu.Lock()
		C.startDNSService(C.CString(host))
	}()
}
