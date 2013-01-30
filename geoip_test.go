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
	gi := Open()
	if gi == nil {
		fmt.Printf("Could not open GeoIP database\n")
	}

	c.Check(gi, NotNil)

	country := gi.GetCountry("207.171.7.51")
	c.Check(country, Equals, "US")
}
