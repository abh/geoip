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
	}

	c.Check(gi, NotNil)

	country, netmask := gi.GetCountry("207.171.7.51")
	c.Check(country, Equals, "US")
	c.Check(netmask, Equals, 15)

	country, netmask = gi.GetCountry("64.235.248.1")
	c.Check(country, Equals, "US")
	c.Check(netmask, Equals, 20)
}
