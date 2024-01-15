package dnssd

// #cgo CFLAGS: -g -Wall
// #include "dns-sd.h"
import "C"

type DNSRecord interface {
	ref() C.DNSRecordRef

	implementsOpaque()
}

type dnsRecord struct {
	_ref C.DNSRecordRef
}

func (r *dnsRecord) ref() C.DNSRecordRef {
	return r._ref
}

func (r *dnsRecord) implementsOpaque() {}
