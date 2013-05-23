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
