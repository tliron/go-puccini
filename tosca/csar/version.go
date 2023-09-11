package csar

import (
	"fmt"
	"strconv"
	"strings"
)

//
// Version
//

type Version struct {
	Major uint8 `json:"major" yaml:"major"`
	Minor uint8 `json:"minor" yaml:"minor"`
}

func ParseVersion(value string) (*Version, error) {
	split := strings.Split(value, ".")
	if len(split) != 2 {
		return nil, fmt.Errorf("malformed version in TOSCA.meta: %s", value)
	}

	major, err := parseDigit(value, split[0])
	if err != nil {
		return nil, err
	}

	minor, err := parseDigit(value, split[1])
	if err != nil {
		return nil, err
	}

	return &Version{major, minor}, nil
}

// ([fmt.Stringer] interface)
func (self *Version) String() string {
	return fmt.Sprintf("%d.%d", self.Major, self.Minor)
}

func parseDigit(value string, digit string) (uint8, error) {
	d, err := strconv.ParseUint(digit, 10, 8)

	// Seriously, the spec says this should be a single digit (?!)
	if (err != nil) || (d > 9) {
		return 0, fmt.Errorf("malformed version in TOSCA.meta: %s", value)
	}

	return uint8(d), nil
}
