/* Go (cgo) interface to libgeoip */
package geoip

/*
#cgo CFLAGS: -I/opt/local/include -L/opt/local/lib
#cgo LDFLAGS: -lGeoIP -L/opt/local/lib
#include <stdio.h>
#include <errno.h>
#include <GeoIP.h>

//typedef GeoIP* GeoIP_pnt
*/
import "C"

import (
	"unsafe"
)

type GeoIP struct {
	gi *C.GeoIP
}

// Open opens a GeoIP database, all formats supported by libgeoip are supported though
// there are only functions to access the country information in this API.
// The database is opened in MEMORY_CACHE mode, if you need to optimize for memory
// instead of performance you should change this.
func GeoIP_Open(base string) *GeoIP {
	cbase := C.CString(base)
	gi := C.GeoIP_open(cbase, C.GEOIP_MEMORY_CACHE)
	C.GeoIP_set_charset(gi, C.GEOIP_CHARSET_UTF8)
	C.free(unsafe.Pointer(cbase))
	if gi == nil {
		return nil
	}
	return &GeoIP{gi}
}

// GetCountry takes an IPv4 address string and returns the country code for that IP.
func (gi *GeoIP) GetCountry(ip string) string {
	if gi == nil {
		return ""
	}
	cip := C.CString(ip)
	ccountry := C.GeoIP_country_code_by_addr(gi.gi, cip)
	C.free(unsafe.Pointer(cip))
	if ccountry != nil {
		rets := C.GoString(ccountry)
		return rets
	}
	return ""
}

// GetCountry_v6 works the same as GetCountry except for IPv6 addresses, be sure to
// load a database with IPv6 data to get any results.
func (gi *GeoIP) GetCountry_v6(ip string) string {
	if gi == nil {
		return ""
	}
	cip := C.CString(ip)
	ccountry := C.GeoIP_country_code_by_addr_v6(gi.gi, cip)
	C.free(unsafe.Pointer(cip))
	if ccountry != nil {
		rets := C.GoString(ccountry)
		return rets
	}
	return ""
}
