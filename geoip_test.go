package geoip

import (
	"fmt"
	. "launchpad.net/gocheck"
	"testing"
)

// Hook up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type GeoIPSuite struct {
}

var _ = Suite(&GeoIPSuite{})

func (s *GeoIPSuite) Testv4(c *C) {
	gi, err := Open()
	if gi == nil || err != nil {
		fmt.Printf("Could not open GeoIP database: %s\n", err)
		return
	}

	c.Check(gi, NotNil)

	country, netmask := gi.GetCountry("207.171.7.51")
	c.Check(country, Equals, "US")
	c.Check(netmask, Equals, 15)

	country, netmask = gi.GetCountry("149.20.64.42")
	c.Check(country, Equals, "US")
	c.Check(netmask, Equals, 13)
}

func (s *GeoIPSuite) TestOpenType(c *C) {

	// SetCustomDirectory("/Users/ask/go/src/geoip/db")

	// Open Country database
	gi, err := OpenType(GEOIP_COUNTRY_EDITION)
	c.Check(err, IsNil)
	c.Assert(gi, NotNil)
	country, _ := gi.GetCountry("207.171.7.51")
	c.Check(country, Equals, "US")
}

func (s *GeoIPSuite) Benchmark_GetCountry(c *C) {
	gi, err := Open()
	if gi == nil || err != nil {
		fmt.Printf("Could not open GeoIP database: %s\n", err)
		return
	}

	for i := 0; i < c.N; i++ {
		gi.GetCountry("207.171.7.51")
	}
}

func (s *GeoIPSuite) Testv4Record(c *C) {
	gi, err := Open("db/GeoLiteCity.dat")
	if gi == nil || err != nil {
		fmt.Printf("Could not open GeoIP database: %s\n", err)
		return
	}

	c.Check(gi, NotNil)

	record := gi.GetRecord("207.171.7.51")
	c.Assert(record, NotNil)
	c.Check(record.CountryCode, Equals, "US")
	fmt.Printf("Record: %#v\n", record)

}

func (s *GeoIPSuite) Benchmark_GetRecord(c *C) {

	gi, err := Open("db/GeoIPCity.dat")
	if gi == nil || err != nil {
		fmt.Printf("Could not open GeoIP database: %s\n", err)
		return
	}

	for i := 0; i < c.N; i++ {
		record := gi.GetRecord("207.171.7.51")
		if record == nil {
			panic("")
		}
	}
}

func (s *GeoIPSuite) Testv4Region(c *C) {
	gi, err := Open("db/GeoIPRegion.dat")
	if gi == nil || err != nil {
		fmt.Printf("Could not open GeoIP database: %s\n", err)
		return
	}

	country, region := gi.GetRegion("207.171.7.51")
	c.Check(country, Equals, "US")
	c.Check(region, Equals, "CA")
}

func (s *GeoIPSuite) TestRegionName(c *C) {
	regionName := GetRegionName("NL", "07")
	c.Check(regionName, Equals, "Noord-Holland")
	regionName = GetRegionName("CA", "ON")
	c.Check(regionName, Equals, "Ontario")
}

func (s *GeoIPSuite) TestGetContinent(c *C) {
	gi, err := Open("db/GeoIP.dat")
	if gi == nil || err != nil {
		fmt.Printf("Could not open GeoIP database: %s\n", err)
		return
	}

	continentName, countryName, netmask := gi.GetContinent("207.171.7.51")
	c.Check(continentName, Equals, "NA")
	c.Check(countryName, Equals, "United States")
	c.Check(netmask, Equals, 15)

	continentName, countryName, netmask = gi.GetContinent("62.217.45.197")
	c.Check(continentName, Equals, "EU")
	c.Check(countryName, Equals, "Germany")
	c.Check(netmask, Equals, 19)
}

func (s *GeoIPSuite) Benchmark_GetContinent(c *C) {
	gi, err := Open("db/GeoIP.dat")
	if gi == nil || err != nil {
		fmt.Printf("Could not open GeoIP database: %s\n", err)
		return
	}

	for i := 0; i < c.N; i++ {
		continentName, countryName, netmask := gi.GetContinent("207.171.7.51")
		if continentName == "" || countryName == "" || netmask == 0 {
			panic("")
		}
	}
}
