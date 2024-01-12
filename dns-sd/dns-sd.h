#ifndef _DNSSD_H
#define _DNSSD_H

#include "dns_sd.h"

DNSServiceErrorType RegisterProxyAddressRecord(DNSServiceRef sdref, DNSRecordRef *RecordRef, const char *host, const char *ip, DNSServiceFlags flags);

#endif
