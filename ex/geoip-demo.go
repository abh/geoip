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
	"github.com/abh/geoip"
)

func main() {

	geoip.Quiet = true

	file6 := "./GeoIPv6-20120730.dat"

	gi6 := geoip.Open(file6)
	if gi6 == nil {
		fmt.Printf("Could not open GeoIPv6 database\n")
	}

	gi := geoip.Open()
	if gi == nil {
		fmt.Printf("Could not open GeoIP database\n")
	}

	if gi != nil {
		test4(*gi, "207.171.7.51")
		test4(*gi, "127.0.0.1")
	}
	if gi6 != nil {
		ip := "2607:f238:2::5"
		country := gi6.GetCountry_v6(ip)
		fmt.Println(ip, ": ")
		display(country)
	}

}

func test4(g geoip.GeoIP, ip string) {
	test(func(s string) string { return g.GetCountry(s) }, ip)
}

func test(f func(string) string, ip string) {
	country := f(ip)
	fmt.Print(ip, ": ")
	display(country)

}
func display(country string) {
	switch country {
	case "":
		fmt.Printf("Could not get country\n")
	default:
		fmt.Printf("%s\n", country)
	}

}
