package dnssd

// #cgo CFLAGS: -g -Wall
// #include "dns-sd.h"
import "C"
import (
	"context"
	"fmt"
	"os"

	"github.com/samber/lo"
	"golang.org/x/sys/unix"
)

type DNSService interface {
	RegisterProxyAddressRecord(host string, ip string, flags C.DNSServiceFlags) (DNSRecord, error)
	RemoveRecord(record DNSRecord, flags C.DNSServiceFlags) error
	Process(ctx context.Context) error
	Deallocate()

	implementsOpaque()
}

type dnsService struct {
	ref C.DNSServiceRef
}

func NewConnection() (DNSService, error) {
	var service dnsService
	res := C.DNSServiceCreateConnection(&service.ref)
	if res != C.kDNSServiceErr_NoError {
		return nil, NewDNSServiceError(res)
	}
	return &service, nil
}

func (s *dnsService) RegisterProxyAddressRecord(host string, ip string, flags C.DNSServiceFlags) (DNSRecord, error) {
	var record dnsRecord
	res := C.RegisterProxyAddressRecord(s.ref, &record._ref, C.CString(host), C.CString(ip), C.uint32_t(flags))
	if res != C.kDNSServiceErr_NoError {
		return nil, NewDNSServiceError(res)
	}
	return &record, nil
}

func (s *dnsService) RemoveRecord(record DNSRecord, flags C.DNSServiceFlags) error {
	res := C.DNSServiceRemoveRecord(s.ref, record.ref(), flags)
	if res != C.kDNSServiceErr_NoError {
		return NewDNSServiceError(res)
	}
	return nil
}

func (s *dnsService) Process(ctx context.Context) error {
	socket := s.useNonblockingSocket()

	fd := os.NewFile(uintptr(socket), "dnssd")
	defer fd.Close()

	for {
		if isContextDone(ctx) {
			return nil
		}

		readSet := unix.FdSet{}
		readSet.Zero()
		readSet.Set(int(socket))

		result, err := unix.Select(int(fd.Fd())+1, &readSet, nil, nil, &unix.Timeval{Sec: 10})
		if err != nil {
			if isContextDone(ctx) {
				return nil
			}
			return err
		}
		if result > 0 {
			res := C.DNSServiceProcessResult(s.ref)
			if res != C.kDNSServiceErr_NoError {
				return NewDNSServiceError(res)
			}
		} else if result == 0 {
			continue
		} else {
			panic(fmt.Sprintf("select error: %d", result))
		}
	}
}

func (s *dnsService) Deallocate() {
	C.DNSServiceRefDeallocate(s.ref)
}

func (s *dnsService) socket() C.dnssd_sock_t {
	return C.DNSServiceRefSockFD(s.ref)
}

func (s *dnsService) useNonblockingSocket() C.dnssd_sock_t {
	socket := s.socket()
	flags := lo.Must1(unix.FcntlInt(uintptr(socket), unix.F_GETFL, 0))
	if flags == -1 {
		flags = 0
	}
	_ = lo.Must1(unix.FcntlInt(uintptr(socket), unix.F_SETFL, flags|unix.O_NONBLOCK))
	return socket
}

func (s *dnsService) implementsOpaque() {}

func isContextDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
