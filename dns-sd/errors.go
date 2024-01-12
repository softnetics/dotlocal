package dnssd

// #cgo CFLAGS: -g -Wall
// #include "dns-sd.h"
import "C"
import "fmt"

type DNSServiceError struct {
	Code C.DNSServiceErrorType
}

func NewDNSServiceError(code C.DNSServiceErrorType) *DNSServiceError {
	return &DNSServiceError{code}
}

func (m *DNSServiceError) Error() string {
	return fmt.Sprintf("DNSServiceError %d", m.Code)
}
