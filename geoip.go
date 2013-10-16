/* Go (cgo) interface to libgeoip */
package geoip

/*
#cgo CFLAGS: -I/opt/local/include -I/usr/local/include -I/usr/include
#cgo LDFLAGS: -lGeoIP -L/opt/local/lib -L/usr/local/lib -L/usr/lib
#include <stdio.h>
#include <errno.h>
#include <GeoIP.h>
#include <GeoIPCity.h>

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

func (gi *GeoIP) free() {
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

// Opens a GeoIP database by filename, all formats supported by libgeoip are
// supported though there are only functions to access some of the databases in this API.
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
	runtime.SetFinalizer(g, (*GeoIP).free)

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

// SetCustomDirectory sets the default location for the GeoIP .dat files used when
// calling OpenType()
func SetCustomDirectory(dir string) {
	cdir := C.CString(dir)
	// GeoIP doesn't copy the string, so don't free it when we're done here.
	// defer C.free(unsafe.Pointer(cdir))
	C.GeoIP_setup_custom_directory(cdir)
}

// OpenType opens a specified GeoIP database type in the default location. Constants
// are defined for each database type (for example GEOIP_COUNTRY_EDITION).
func OpenType(dbType int) (*GeoIP, error) {
	g := &GeoIP{}
	runtime.SetFinalizer(g, (*GeoIP).free)

	var err error

	g.db, err = C.GeoIP_open_type(C.int(dbType), C.GEOIP_MEMORY_CACHE)
	if err != nil {
		return nil, fmt.Errorf("Error opening GeoIP database (%d): %s", dbType, err)
	}

	if g.db == nil {
		return nil, fmt.Errorf("Didn't open GeoIP database (%d)", dbType)
	}

	C.GeoIP_set_charset(g.db, C.GEOIP_CHARSET_UTF8)

	return g, nil
}

// Takes an IPv4 address string and returns the organization name for that IP.
// Requires the GeoIP organization database.
func (gi *GeoIP) GetOrg(ip string) string {
	name, _ := gi.GetName(ip)
	return name
}

// Works on the ASN, Netspeed, Organization and probably other
// databases, takes and IP string and returns a "name" and the
// netmask.
func (gi *GeoIP) GetName(ip string) (name string, netmask int) {
	if gi.db == nil {
		return
	}

	gi.mu.Lock()
	defer gi.mu.Unlock()

	cip := C.CString(ip)
	defer C.free(unsafe.Pointer(cip))
	cname := C.GeoIP_name_by_addr(gi.db, cip)

	if cname != nil {
		name = C.GoString(cname)
		defer C.free(unsafe.Pointer(cname))
		netmask = int(C.GeoIP_last_netmask(gi.db))
		return
	}
	return
}

type GeoIPRecord struct {
	CountryCode  string
	CountryCode3 string
	CountryName  string
	Region       string
	City         string
	PostalCode   string
	Latitude     float32
	Longitude    float32
	// DMACode       int
	AreaCode      int
	CharSet       int
	ContinentCode string
}

// Returns the "City Record" for an IP address. Requires the GeoCity(Lite)
// database - http://www.maxmind.com/en/city
func (gi *GeoIP) GetRecord(ip string) *GeoIPRecord {
	if gi.db == nil {
		return nil
	}

	cip := C.CString(ip)
	defer C.free(unsafe.Pointer(cip))
	record := C.GeoIP_record_by_addr(gi.db, cip)
	if record == nil {
		return nil
	}
	// defer C.free(unsafe.Pointer(record))
	defer C.GeoIPRecord_delete(record)
	rec := new(GeoIPRecord)
	rec.CountryCode = C.GoString(record.country_code)
	rec.CountryCode3 = C.GoString(record.country_code3)
	rec.CountryName = C.GoString(record.country_name)
	rec.Region = C.GoString(record.region)
	rec.City = C.GoString(record.city)
	rec.PostalCode = C.GoString(record.postal_code)
	rec.Latitude = float32(record.latitude)
	rec.Longitude = float32(record.longitude)
	rec.AreaCode = int(record.area_code)
	rec.CharSet = int(record.charset)
	rec.ContinentCode = C.GoString(record.continent_code)

	return rec
}

// Returns the country code and region code for an IP address. Requires
// the GeoIP Region database.
func (gi *GeoIP) GetRegion(ip string) (string, string) {
	if gi.db == nil {
		return "", ""
	}

	cip := C.CString(ip)
	defer C.free(unsafe.Pointer(cip))
	region := C.GeoIP_region_by_addr(gi.db, cip)
	if region == nil {
		return "", ""
	}

	countryCode := C.GoString(&region.country_code[0])
	regionCode := C.GoString(&region.region[0])
	defer C.free(unsafe.Pointer(region))

	return countryCode, regionCode
}

// Returns the region name given a country code and region code
func GetRegionName(countryCode, regionCode string) string {

	cc := C.CString(countryCode)
	defer C.free(unsafe.Pointer(cc))

	rc := C.CString(regionCode)
	defer C.free(unsafe.Pointer(rc))

	region := C.GeoIP_region_name_by_code(cc, rc)
	if region == nil {
		return ""
	}

	// it's a static string constant, don't free this
	regionName := C.GoString(region)

	return regionName
}

// Same as GetName() but for IPv6 addresses.
func (gi *GeoIP) GetNameV6(ip string) (name string, netmask int) {
	if gi.db == nil {
		return
	}

	gi.mu.Lock()
	defer gi.mu.Unlock()

	cip := C.CString(ip)
	defer C.free(unsafe.Pointer(cip))
	cname := C.GeoIP_name_by_addr_v6(gi.db, cip)

	if cname != nil {
		name = C.GoString(cname)
		defer C.free(unsafe.Pointer(cname))
		netmask = int(C.GeoIP_last_netmask(gi.db))
		return
	}
	return
}

// Takes an IPv4 address string and returns the country code for that IP
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
// load a database with IPv6 data to get any results.
func (gi *GeoIP) GetCountry_v6(ip string) (cc string, netmask int) {
	if gi.db == nil {
		return
	}

	gi.mu.Lock()
	defer gi.mu.Unlock()

	cip := C.CString(ip)
	defer C.free(unsafe.Pointer(cip))
	ccountry := C.GeoIP_country_code_by_addr_v6(gi.db, cip)
	if ccountry != nil {
		cc = C.GoString(ccountry)
		netmask = int(C.GeoIP_last_netmask(gi.db))
		return
	}
	return
}

func (gi *GeoIP) GetContinent(ip string) (ccontinent, ccountry string, netmask int) {
	if gi.db == nil {
		return
	}

	gi.mu.Lock()
	defer gi.mu.Unlock()

	cip := C.CString(ip)
	defer C.free(unsafe.Pointer(cip))

	// First determine country ID and country string
	ccountryID, err := C.GeoIP_id_by_addr(gi.db, cip)
	if err != nil {
		return
	}
	ccountryName, err := C.GeoIP_country_name_by_id(gi.db, ccountryID)
	if err != nil {
		return
	}
	ccountry = C.GoString(ccountryName)

	// This is a simple static lookup of the continent from the ID
	ccontinentName := C.GeoIP_continent_by_id(ccountryID)
	ccontinent = C.GoString(ccontinentName)

	// Also return the netmask for convenience
	netmask = int(C.GeoIP_last_netmask(gi.db))
	return
}
