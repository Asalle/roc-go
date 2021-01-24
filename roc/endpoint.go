package roc

/*
#include <roc/endpoint.h>
*/
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"
)

// Endpoint represents network endpoint endpoint.
//
// Consists of IP endpoint plus UDP or TCP port number.
// Similar to net.Addr in Go and struct sockaddr in C.
//
// Should not be used concurrently
type Endpoint struct {
	raw *C.roc_endpoint
	mem []byte
}

// NewEndpoint creates and initializes a new Endpoint.
//
// The IP endpoint is parsed from a string representation. If family is AfAuto, the
// endpoint family is auto-detected from string format. Otherwise, the string format
// should correspond to the family specified.
//
// The port number should be in range [0; 65536).
//
// When Endpoint is used to bind a sender or receiver port, the "0.0.0.0" (for IPv4)
// or "::" (for IPv6) may be used to bind the port to all network interfaces, and
// zero port number may be used to bind the port to a randomly chosen ephemeral port.
func NewEndpoint() (*Endpoint, error) {
	e := new(Endpoint)
	e.mem = make([]byte, C.sizeof_roc_endpoint)
	e.raw = (*C.roc_endpoint)(unsafe.Pointer(&a.mem[0]))

	errCode := C.roc_endpoint_allocate(&e.raw)
	if errCode == 0 {
		return e, nil
	}
	if errCode < 0 {
		return nil, ErrInvalidArgs
	}

	panic(fmt.Sprintf(
		"unexpected return code %d from roc_endpoint_init()", errCode))
}

// func NewEndpoint(family Family, ip string, port int) (*Endpoint, error) {
// a := new(Endpoint)
// a.mem = make([]byte, C.sizeof_roc_endpoint)
// a.raw = (*C.roc_endpoint)(unsafe.Pointer(&a.mem[0]))
//
// cfamily := (C.roc_family)(family)
// cip := toCStr(ip)
// cport := (C.int)(port)
// errCode := C.roc_endpoint_init(a.raw, cfamily, (*C.char)(unsafe.Pointer(&cip[0])), cport)
//
// if errCode == 0 {
// return a, nil
// }
// if errCode < 0 {
// return nil, ErrInvalidArgs
// }
//
// panic(fmt.Sprintf(
// "unexpected return code %d from roc_endpoint_init()", errCode))
// }

// Family returns endpoint family.
//
// If AfAuto was used to construct endpoint, the actually selected family, i.e.
// either AfIPv4 or AfIPv6, is reported.
func (a *Endpoint) Family() Family {
	f := C.roc_endpoint_family(a.raw)
	family := (Family)(f)
	if family == afInvalid {
		panic("unexpected failure in roc_endpoint_family()")
	}
	return family
}

// IP returns IP endpoint formatted to string.
func (a *Endpoint) IP() string {
	const buflen = 255
	sIP := make([]byte, buflen)
	res := C.roc_endpoint_ip(a.raw, (*C.char)(unsafe.Pointer(&sIP[0])), buflen)
	if res == nil {
		panic("unexpected failure in roc_endpoint_ip()")
	}
	ret := C.GoString(res)
	runtime.KeepAlive(sIP)
	return ret
}

// Port return UDP or TCP port number.
//
// If Endpoint was passed to sender or receiver bind and the initial port number was
// zero, which means "use random port", this function will return the actually
// selected port number.
func (a *Endpoint) Port() int {
	res := C.roc_endpoint_port(a.raw)
	if res < 0 {
		panic("unexpected failure in roc_endpoint_port()")
	}
	return (int)(res)
}
