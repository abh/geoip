/* Go (cgo) interface to libgeoip */
package geoip

/*
#cgo CFLAGS: -I/opt/local/include -I/usr/local/include -I/usr/include
#cgo LDFLAGS: -lGeoIP -L/opt/local/lib -L/usr/local/lib -L/usr/lib
#include <stdio.h>
#include <errno.h>
#include <GeoIP.h>

//typedef GeoIP* GeoIP_pnt
*/
import "C"

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"unsafe"
)

type GeoIP struct {
	db *C.GeoIP
	mu sync.Mutex
}

func (gi *GeoIP) Free() {
	if gi == nil {
		return
	}
	if gi.db == nil {
		gi = nil
		return
	}
	C.GeoIP_delete(gi.db)
	gi = nil
	return
}

// Open opens a GeoIP database, all formats supported by libgeoip are supported though
// there are only functions to access the country information in this API.
// The database is opened in MEMORY_CACHE mode, if you need to optimize for memory
// instead of performance you should change this.
// If you don't pass a filename, it will try opening the database from
// a list of common paths.
func Open(files ...string) (*GeoIP, error) {
	if len(files) == 0 {
		files = []string{
			"/usr/share/GeoIP/GeoIP.dat",       // Linux default
			"/usr/share/local/GeoIP/GeoIP.dat", // source install?
			"/usr/local/share/GeoIP/GeoIP.dat", // FreeBSD
			"/opt/local/share/GeoIP/GeoIP.dat", // MacPorts
			"/usr/share/GeoIP/GeoIP.dat",       // ArchLinux
		}
	}

	g := &GeoIP{}
	runtime.SetFinalizer(g, (*GeoIP).Free)

	var err error

	for _, file := range files {

		// libgeoip prints errors if it can't open the file, so check first
		if _, err := os.Stat(file); err != nil {
			if os.IsExist(err) {
				log.Println(err)
			}
			continue
		}

		cbase := C.CString(file)
		defer C.free(unsafe.Pointer(cbase))

		g.db, err = C.GeoIP_open(cbase, C.GEOIP_MEMORY_CACHE)
		if g.db != nil && err != nil {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("Error opening GeoIP database (%s): %s", files, err)
	}

	if g.db == nil {
		return nil, fmt.Errorf("Didn't open GeoIP database (%s)", files)
	}

	C.GeoIP_set_charset(g.db, C.GEOIP_CHARSET_UTF8)
	return g, nil
}

// GetOrg takes an IPv4 address string and returns the org name for that IP.
// Requires the geoip org database
func (gi *GeoIP) GetOrg(ip string) string {
	if gi.db == nil {
		return ""
	}
	cip := C.CString(ip)
	defer C.free(unsafe.Pointer(cip))
	cname := C.GeoIP_name_by_addr(gi.db, cip)
	if cname != nil {
		rets := C.GoString(cname)
		return rets
	}
	return ""
}

// GetCountry takes an IPv4 address string and returns the country code for that IP
// and the netmask for that IP range.
func (gi *GeoIP) GetCountry(ip string) (cc string, netmask int) {
	if gi.db == nil {
		return
	}

	gi.mu.Lock() // Lock to make sure we get the right result from GeoIP_last_netmask
	defer gi.mu.Unlock()

	cip := C.CString(ip)
	defer C.free(unsafe.Pointer(cip))
	ccountry := C.GeoIP_country_code_by_addr(gi.db, cip)

	if ccountry != nil {
		cc = C.GoString(ccountry)
		netmask = int(C.GeoIP_last_netmask(gi.db))
		return
	}
	return
}

// GetCountry_v6 works the same as GetCountry except for IPv6 addresses, be sure to
// load a database with IPv6 data to get any results. Getting the netmask for IPv6
// addresses is not supported by the geoip library.
func (gi *GeoIP) GetCountry_v6(ip string) (cc string) {
	if gi.db == nil {
		return
	}

	cip := C.CString(ip)
	defer C.free(unsafe.Pointer(cip))
	ccountry := C.GeoIP_country_code_by_addr_v6(gi.db, cip)
	if ccountry != nil {
		cc = C.GoString(ccountry)
		return
	}
	return
}
