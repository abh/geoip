/* 
 * Original author: Stiletto <blasux@blasux.ru>
 *
 * This program is free software. It comes without any warranty, to
 * the extent permitted by applicable law. You can redistribute it
 * and/or modify it under the terms of the Do What The Fuck You Want
 * To Public License, Version 2, as published by Sam Hocevar. See
 * http://sam.zoy.org/wtfpl/COPYING for more details. */

package main

import (
	"fmt"
	"geoip"
)

func main() {

	file6 := "./GeoIPv6-20120730.dat"
	file := "/opt/local/share/GeoIP/GeoIP.dat"

	gi6 := geoip.GeoIP_Open(file6)
	if gi6 == nil {
		fmt.Printf("Could not open GeoIPv6 database\n")
		return
	}
	gi := geoip.GeoIP_Open(file)
	if gi == nil {
		fmt.Printf("Could not open GeoIP database\n")
		return
	}

	country := gi.GetCountry("207.171.7.51")
	display(country)

	country = gi6.GetCountry_v6("2607:f238:2::5")
	display(country)

}

func display(country *string) {
	switch country {
	case nil:
		fmt.Printf("Could not get country\n")
	default:
		fmt.Printf("Country %s\n", *country)
	}

}
