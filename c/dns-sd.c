#include <ctype.h>
#include <stdio.h>          // For stdout, stderr
#include <stdlib.h>         // For exit()
#include <string.h>         // For strlen(), strcpy()
#include <stdarg.h>         // For va_start, va_arg, va_end, etc.
#include <errno.h>          // For errno, EINTR
#include <time.h>
#include <sys/types.h>      // For u_char
#include <unistd.h>
#include <sys/stat.h>



#ifndef __printflike
    #define __printflike(A, B)
#endif

#ifdef _WIN32
    #include <winsock2.h>
    #include <ws2tcpip.h>
    #include <Iphlpapi.h>
    #include <process.h>
    #include <stdint.h>
typedef int pid_t;
typedef int suseconds_t;
    #define getpid     _getpid
    #define strcasecmp _stricmp
    #define snprintf   _snprintf
static const char kFilePathSep = '\\';
    #ifndef HeapEnableTerminationOnCorruption
    #     define HeapEnableTerminationOnCorruption (HEAP_INFORMATION_CLASS)1
    #endif
    #if !defined(IFNAMSIZ)
     #define IFNAMSIZ 16
    #endif
    #define if_nametoindex if_nametoindex_win
    #define if_indextoname if_indextoname_win

typedef PCHAR (WINAPI * if_indextoname_funcptr_t)(ULONG index, PCHAR name);
typedef ULONG (WINAPI * if_nametoindex_funcptr_t)(PCSTR name);

unsigned if_nametoindex_win(const char *ifname)
{
    HMODULE library;
    unsigned index = 0;

    // Try and load the IP helper library dll
    if ((library = LoadLibrary(TEXT("Iphlpapi")) ) != NULL )
    {
        if_nametoindex_funcptr_t if_nametoindex_funcptr;

        // On Vista and above there is a Posix like implementation of if_nametoindex
        if ((if_nametoindex_funcptr = (if_nametoindex_funcptr_t) GetProcAddress(library, "if_nametoindex")) != NULL )
        {
            index = if_nametoindex_funcptr(ifname);
        }

        FreeLibrary(library);
    }

    return index;
}

char * if_indextoname_win( unsigned ifindex, char *ifname)
{
    HMODULE library;
    char * name = NULL;

    // Try and load the IP helper library dll
    if ((library = LoadLibrary(TEXT("Iphlpapi")) ) != NULL )
    {
        if_indextoname_funcptr_t if_indextoname_funcptr;

        // On Vista and above there is a Posix like implementation of if_indextoname
        if ((if_indextoname_funcptr = (if_indextoname_funcptr_t) GetProcAddress(library, "if_indextoname")) != NULL )
        {
            name = if_indextoname_funcptr(ifindex, ifname);
        }

        FreeLibrary(library);
    }

    return name;
}

static size_t _sa_len(const struct sockaddr *addr)
{
    if (addr->sa_family == AF_INET) return (sizeof(struct sockaddr_in));
    else if (addr->sa_family == AF_INET6) return (sizeof(struct sockaddr_in6));
    else return (sizeof(struct sockaddr));
}

#   define SA_LEN(addr) (_sa_len(addr))

typedef void (WINAPI* SystemTimeFunc)(LPFILETIME);

static const uint64_t epoch_diff = (UINT64)11644473600000000ULL;
static SystemTimeFunc fpTimeFunc;

int gettimeofday(struct timeval* tp, struct timezone* tzp)
{
    FILETIME ft;
    UINT64 us;

    if (!fpTimeFunc)
    {
        /* available on Windows 7 */
        fpTimeFunc = GetSystemTimeAsFileTime;

        HMODULE hKernel32 = LoadLibraryW(L"kernel32.dll");
        if (hKernel32)
        {
            FARPROC fp;

            /* available on Windows 8+ */
            fp = GetProcAddress(hKernel32, "GetSystemTimePreciseAsFileTime");
            if (fp)
            {
                fpTimeFunc = (SystemTimeFunc)fp;
            }
        }
    }

    fpTimeFunc(&ft);

    us = (((uint64_t)ft.dwHighDateTime << 32) | (uint64_t)ft.dwLowDateTime) / 10;
    us -= epoch_diff;

    tp->tv_sec = (long)(us / 1000000);
    tp->tv_usec = (long)(us % 1000000);

    return 0;
}

#else
    #include <unistd.h>         // For getopt() and optind
    #include <netdb.h>          // For getaddrinfo()
    #include <sys/time.h>       // For struct timeval
    #include <sys/socket.h>     // For AF_INET
    #include <netinet/in.h>     // For struct sockaddr_in()
    #include <arpa/inet.h>      // For inet_addr()
    #include <net/if.h>         // For if_nametoindex()
static const char kFilePathSep = '/';
// #ifndef NOT_HAVE_SA_LEN
//  #define SA_LEN(addr) ((addr)->sa_len)
// #else
    #define SA_LEN(addr) (((addr)->sa_family == AF_INET6) ? sizeof(struct sockaddr_in6) : sizeof(struct sockaddr_in))
// #endif
#endif

#if (TEST_NEW_CLIENTSTUB && !defined(__APPLE_API_PRIVATE))
#define __APPLE_API_PRIVATE 1
#endif

// DNSServiceSetDispatchQueue is not supported on 10.6 & prior
#if !TEST_NEW_CLIENTSTUB && defined(__ENVIRONMENT_MAC_OS_X_VERSION_MIN_REQUIRED__) && (__ENVIRONMENT_MAC_OS_X_VERSION_MIN_REQUIRED__ - (__ENVIRONMENT_MAC_OS_X_VERSION_MIN_REQUIRED__ % 10) <= 1060)
#undef _DNS_SD_LIBDISPATCH
#endif
#include "dns_sd.h"
#include "ClientCommon.h"


#if TEST_NEW_CLIENTSTUB
#include "../mDNSShared/dnssd_ipc.c"
#include "../mDNSShared/dnssd_clientlib.c"
#include "../mDNSShared/dnssd_clientstub.c"
#endif

/** 
 * Global
*/

#if _DNS_SD_LIBDISPATCH
static dispatch_queue_t main_queue;
static dispatch_source_t timer_source;
#endif

#if _DNS_SD_LIBDISPATCH
#define EXIT_IF_LIBDISPATCH_FATAL_ERROR(E) \
    if (main_queue && (E) == kDNSServiceErr_ServiceNotRunning) { fprintf(stderr, "Error code %d\n", (E)); exit(0); }
#else
#define EXIT_IF_LIBDISPATCH_FATAL_ERROR(E)
#endif

static int exitWhenNoMoreComing;

#define printtimestamp() printtimestamp_F(stdout)

static void printtimestamp_F(FILE *outstream)
{
    struct tm tm;
    int ms;
    static char date[16];
    static char new_date[16];
#ifdef _WIN32
    SYSTEMTIME sysTime;
    time_t uct = time(NULL);
    tm = *localtime(&uct);
    GetLocalTime(&sysTime);
    ms = sysTime.wMilliseconds;
#else
    struct timeval tv;
    gettimeofday(&tv, NULL);
    localtime_r((time_t*)&tv.tv_sec, &tm);
    ms = tv.tv_usec/1000;
#endif
    strftime(new_date, sizeof(new_date), "%a %d %b %Y", &tm);
    if (strncmp(date, new_date, sizeof(new_date)))
    {
        fprintf(outstream, "DATE: ---%s---\n", new_date); //display date only if it has changed
        strncpy(date, new_date, sizeof(date));
    }
    fprintf(outstream, "%2d:%02d:%02d.%03d  ", tm.tm_hour, tm.tm_min, tm.tm_sec, ms);
}

static void DNSSD_API MyRegisterRecordCallback(DNSServiceRef service, DNSRecordRef rec, const DNSServiceFlags flags,
                                               DNSServiceErrorType errorCode, void *context)
{
    char *name = (char *)context;

    (void)service;  // Unused
    (void)rec;      // Unused
    (void)flags;    // Unused
    EXIT_IF_LIBDISPATCH_FATAL_ERROR(errorCode);

    printtimestamp();
    printf("Got a reply for record %s: ", name);

    switch (errorCode)
    {
    case kDNSServiceErr_NoError:      printf("Name now registered and active\n"); break;
    case kDNSServiceErr_NameConflict: printf("Name in use, please choose another\n"); exit(-1);
    default:                          printf("Error %d\n", errorCode); break;
    }
    if (!(flags & kDNSServiceFlagsMoreComing))
    {
        fflush(stdout);
        if (exitWhenNoMoreComing) exit(0);
    }
}


static void getip(const char *const name, struct sockaddr_storage *result)
{
    struct addrinfo *addrs = NULL;
    int err = getaddrinfo(name, NULL, NULL, &addrs);
    if (err) fprintf(stderr, "getaddrinfo error %d for %s", err, name);
    else memcpy(result, addrs->ai_addr, SA_LEN(addrs->ai_addr));
    if (addrs) freeaddrinfo(addrs);
}

static DNSServiceErrorType RegisterProxyAddressRecord(DNSServiceRef sdref, const char *host, const char *ip, DNSServiceFlags flags)
{
    // Call getip() after the call DNSServiceCreateConnection().
    // On the Win32 platform, WinSock must be initialized for getip() to succeed.
    // Any DNSService* call will initialize WinSock for us, so we make sure
    // DNSServiceCreateConnection() is called before getip() is.
    struct sockaddr_storage hostaddr;
    static DNSRecordRef record = NULL;
    memset(&hostaddr, 0, sizeof(hostaddr));
    getip(ip, &hostaddr);
    if (!(flags & kDNSServiceFlagsShared))
    {
        flags |= kDNSServiceFlagsUnique;
    }
    if (hostaddr.ss_family == AF_INET)
        return(DNSServiceRegisterRecord(sdref, &record, flags, kDNSServiceInterfaceIndexLocalOnly, host,
                                        kDNSServiceType_A,    kDNSServiceClass_IN,  4, &((struct sockaddr_in *)&hostaddr)->sin_addr,  240, MyRegisterRecordCallback, (void*)host));
    else if (hostaddr.ss_family == AF_INET6)
        return(DNSServiceRegisterRecord(sdref, &record, flags, kDNSServiceInterfaceIndexLocalOnly, host,
                                        kDNSServiceType_AAAA, kDNSServiceClass_IN, 16, &((struct sockaddr_in6*)&hostaddr)->sin6_addr, 240, MyRegisterRecordCallback, (void*)host));
    else return(kDNSServiceErr_BadParam);
}

static DNSServiceRef client_pa = NULL;
static int exitTimeout;


static void HandleEvents(void)
#if _DNS_SD_LIBDISPATCH
{
    main_queue = dispatch_get_main_queue();
    if (client_pa) DNSServiceSetDispatchQueue(client_pa, main_queue);
    dispatch_main();
}
#else
{
    int dns_sd_fd  = client    ? DNSServiceRefSockFD(client   ) : -1;
    int dns_sd_fd2 = client_pa ? DNSServiceRefSockFD(client_pa) : -1;
    int nfds = dns_sd_fd + 1;
    fd_set readfds;
    struct timeval tv;
    int result;
    uint64_t timeout_when, now;
    int expectingMyTimer;

    if (dns_sd_fd2 > dns_sd_fd) nfds = dns_sd_fd2 + 1;

    if (exitTimeout != 0) {
        gettimeofday(&tv, NULL);
        timeout_when = tv.tv_sec * 1000ULL * 1000ULL + tv.tv_usec + exitTimeout * 1000ULL * 1000ULL;
    }

    while (!stopNow)
    {
        // 1. Set up the fd_set as usual here.
        // This example client has no file descriptors of its own,
        // but a real application would call FD_SET to add them to the set here
        FD_ZERO(&readfds);

        // 2. Add the fd for our client(s) to the fd_set
        if (client   ) FD_SET(dns_sd_fd, &readfds);
        if (client_pa) FD_SET(dns_sd_fd2, &readfds);

        // 3. Set up the timeout.
        expectingMyTimer = 1;
        if (exitTimeout > 0) {
            gettimeofday(&tv, NULL);
            now = tv.tv_sec * 1000ULL * 1000ULL + tv.tv_usec;
            if (timeout_when <= now) {
                exit(0);
            }
            if (timeout_when - now < timeOut * 1000ULL * 1000ULL) {
                tv.tv_sec = (time_t)(timeout_when - now) / 1000 / 1000;
                tv.tv_usec = (suseconds_t)(timeout_when % (1000 * 1000));
                expectingMyTimer = 0;
            }
        }
        if (expectingMyTimer) {
            tv.tv_sec  = timeOut;
            tv.tv_usec = 0;
        }
        result = select(nfds, &readfds, (fd_set*)NULL, (fd_set*)NULL, &tv);
        if (result > 0)
        {
            DNSServiceErrorType err = kDNSServiceErr_NoError;
            if      (client    && FD_ISSET(dns_sd_fd, &readfds)) err = DNSServiceProcessResult(client   );
            else if (client_pa && FD_ISSET(dns_sd_fd2, &readfds)) err = DNSServiceProcessResult(client_pa);
            if (err) { printtimestamp_F(stderr); fprintf(stderr, "DNSServiceProcessResult returned %d\n", err); stopNow = 1; }
        }
        else if (result == 0)
        {
            if (expectingMyTimer)
            {
                myTimerCallBack();
            }
            else
            {
                // exitTimeout has elapsed.
                exit(0);
            }
        }
        else
        {
            printf("select() returned %d errno %d %s\n", result, errno, strerror(errno));
            if (errno != EINTR) stopNow = 1;
        }
    }
}
#endif

static bool isEndedWithDotLocal(const char *const hostname) {
    const char *const dotLocal = ".local";
    const size_t hostnameLen = strlen(hostname);
    const size_t dotLocalLen = strlen(dotLocal);
    if (hostnameLen < dotLocalLen) return false;
    return (strcasecmp(hostname + hostnameLen - dotLocalLen, dotLocal) == 0);
}

int main(int argc, char *argv[]) {
    if (argc < 2) {
        printf("Usage: %s <list of hostnames>\n", argv[0]);
        return 1;
    }
    DNSServiceErrorType err;
    DNSServiceFlags flags = 0;
    printtimestamp();
    printf("...STARTING...\n");
    err = DNSServiceCreateConnection(&client_pa);
    if (err) { fprintf(stderr, "DNSServiceCreateConnection returned %d\n", err); return(err); }
    for (int i = 1; i < argc; i++) {
        char* host = argv[i];
        if (!isEndedWithDotLocal(host)) {
            printf("Adding .local to %s\n", host);
            host = malloc(strlen(argv[i]) + 7);
            strcpy(host, argv[i]);
            strcat(host, ".local");
        }
        printf("Registering %s\n", host);
        err = RegisterProxyAddressRecord(client_pa, host, "127.0.0.1", flags);
        if (err) { fprintf(stderr, "DNSServiceRegisterRecord returned %d\n", err); return(err); }
    }
    HandleEvents();

    if (client_pa) DNSServiceRefDeallocate(client_pa);
    return 0;
}
