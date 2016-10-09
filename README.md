# GeoIP API for Go

This package wraps the [libgeoip C library](http://www.maxmind.com/app/c) for
access from Go (golang). [![Build Status](https://travis-ci.org/abh/geoip.png?branch=master)](https://travis-ci.org/abh/geoip)

Install with `go get github.com/abh/geoip` and use [godoc
geoip](http://godoc.org/github.com/abh/geoip) to read the documentation.

There's a small example in the `ex/` subdirectory.

You can download the free [GeoLite
Country](http://www.maxmind.com/app/geoip_country) database or you can
[subscribe to updates](http://www.maxmind.com/app/country).

## Examples

	file := "/usr/share/GeoIP/GeoIP.dat"

	gi, err := geoip.Open(file)
	if err != nil {
		fmt.Printf("Could not open GeoIP database\n")
	}

	if gi != nil {
		country, netmask := gi.GetCountry("207.171.7.51")
	}

	// Setup gi6 by opening the optional IPv6 database and then...
	country := gi6.GetCountry_v6("2607:f238:2::5")
	fmt.Println(country)

## build

###windows use MSYS/MSYS2

1. build libGeoIP (geoip-api-c)
	
		see https://github.com/maxmind/geoip-api-c/releases page
		download latest version
		tar xvf ...

		./configure --disable-shared --disable-dependency-tracking --disable-data-files CFLAGS=-D_WIN32_WINNT=0x0501

	`CFLAGS=-D_WIN32_WINNT=0x0501` enable link winsock2 `ws32`
	without this will show build error :
	
		Undefined reference to getaddrinfo
		...
	 
		
2. build(install) geoip static

		CGO_ENABLED=1 GOOS=windows CGO_CFLAGS=-I/usr/local/include CGO_LDFLAGS='-L/usr/local/lib -lGeoIP -lwsock32 -lws2_32' go install -x --ldflags '-extldflags "-static"' github.com/chennqqi/geoip
	
	
	github.com/chennqqi/geoip or source repo github.com/abh/geoip

	my MSYS don't contain `pkg-config`, so i comment in geoip.go


## Contact

Copyright 2012-2013 Ask Bjørn Hansen <ask@develooper.com>. The package
is MIT licensed, see the LICENSE file. Originally based on example code
from blasux@blasux.ru.
