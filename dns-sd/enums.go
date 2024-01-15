package dnssd

// #cgo CFLAGS: -g -Wall
// #include "dns-sd.h"
import "C"

const (
	kDNSServiceErr_NoError                   = C.kDNSServiceErr_NoError
	kDNSServiceErr_Unknown                   = C.kDNSServiceErr_Unknown
	kDNSServiceErr_NoSuchName                = C.kDNSServiceErr_NoSuchName
	kDNSServiceErr_NoMemory                  = C.kDNSServiceErr_NoMemory
	kDNSServiceErr_BadParam                  = C.kDNSServiceErr_BadParam
	kDNSServiceErr_BadReference              = C.kDNSServiceErr_BadReference
	kDNSServiceErr_BadState                  = C.kDNSServiceErr_BadState
	kDNSServiceErr_BadFlags                  = C.kDNSServiceErr_BadFlags
	kDNSServiceErr_Unsupported               = C.kDNSServiceErr_Unsupported
	kDNSServiceErr_NotInitialized            = C.kDNSServiceErr_NotInitialized
	kDNSServiceErr_AlreadyRegistered         = C.kDNSServiceErr_AlreadyRegistered
	kDNSServiceErr_NameConflict              = C.kDNSServiceErr_NameConflict
	kDNSServiceErr_Invalid                   = C.kDNSServiceErr_Invalid
	kDNSServiceErr_Firewall                  = C.kDNSServiceErr_Firewall
	kDNSServiceErr_Incompatible              = C.kDNSServiceErr_Incompatible
	kDNSServiceErr_BadInterfaceIndex         = C.kDNSServiceErr_BadInterfaceIndex
	kDNSServiceErr_Refused                   = C.kDNSServiceErr_Refused
	kDNSServiceErr_NoSuchRecord              = C.kDNSServiceErr_NoSuchRecord
	kDNSServiceErr_NoAuth                    = C.kDNSServiceErr_NoAuth
	kDNSServiceErr_NoSuchKey                 = C.kDNSServiceErr_NoSuchKey
	kDNSServiceErr_NATTraversal              = C.kDNSServiceErr_NATTraversal
	kDNSServiceErr_DoubleNAT                 = C.kDNSServiceErr_DoubleNAT
	kDNSServiceErr_BadTime                   = C.kDNSServiceErr_BadTime
	kDNSServiceErr_BadSig                    = C.kDNSServiceErr_BadSig
	kDNSServiceErr_BadKey                    = C.kDNSServiceErr_BadKey
	kDNSServiceErr_Transient                 = C.kDNSServiceErr_Transient
	kDNSServiceErr_ServiceNotRunning         = C.kDNSServiceErr_ServiceNotRunning
	kDNSServiceErr_NATPortMappingUnsupported = C.kDNSServiceErr_NATPortMappingUnsupported
	kDNSServiceErr_NATPortMappingDisabled    = C.kDNSServiceErr_NATPortMappingDisabled
	kDNSServiceErr_NoRouter                  = C.kDNSServiceErr_NoRouter
	kDNSServiceErr_PollingMode               = C.kDNSServiceErr_PollingMode
	kDNSServiceErr_Timeout                   = C.kDNSServiceErr_Timeout
	kDNSServiceErr_DefunctConnection         = C.kDNSServiceErr_DefunctConnection
	kDNSServiceErr_PolicyDenied              = C.kDNSServiceErr_PolicyDenied
	kDNSServiceErr_NotPermitted              = C.kDNSServiceErr_NotPermitted
)
