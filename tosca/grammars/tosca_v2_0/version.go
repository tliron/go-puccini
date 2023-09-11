package tosca_v2_0

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/parsing"
)

var VersionRE = regexp.MustCompile(
	`^(?P<major>\d+)\.(?P<minor>\d+)(?:\.(?P<fix>\d+)` +
		`(?:(?:\.(?P<qualifier>\w+))(?:-(?P<build>\d+))?)?)?$`)

//
// Version
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.3.2
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.3.2
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.2.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.2.2
//

type Version struct {
	CanonicalString string `json:"$string" yaml:"$string"`
	OriginalString  string `json:"$originalString" yaml:"$originalString"`
	Comparer        string `json:"$comparer" yaml:"$comparer"`

	Major     uint32 `json:"major" yaml:"major"`
	Minor     uint32 `json:"minor" yaml:"minor"`
	Fix       uint32 `json:"fix" yaml:"fix"`
	Qualifier string `json:"qualifier" yaml:"qualifier"`
	Build     uint32 `json:"build" yaml:"build"`
}

// ([parsing.Reader] signature)
func ReadVersion(context *parsing.Context) parsing.EntityPtr {
	var self Version

	if context.Is(ard.TypeString) {
		self.OriginalString = *context.ReadString()
		self.CanonicalString = self.OriginalString
	} else if context.HasQuirk(parsing.QuirkDataTypesStringPermissive) {
		if context.Is(ard.TypeFloat) {
			value := *context.ReadFloat()
			self.OriginalString = strconv.FormatFloat(value, 'g', -1, 64)
			if strings.Index(self.CanonicalString, "e") == -1 {
				context.ReportValueMalformed("version", "too many floating point digits")
				return &self
			}
			if strings.Index(self.CanonicalString, ".") == -1 {
				// Assume minor version is 0
				self.CanonicalString += ".0"
			}
		} else if context.Is(ard.TypeInteger) {
			value := *context.ReadInteger()
			// Assume minor version is 0
			self.OriginalString = strconv.FormatInt(value, 10) + ".0"
			self.CanonicalString = self.OriginalString
		} else {
			context.ReportValueWrongType(ard.TypeString, ard.TypeFloat, ard.TypeInteger)
			return &self
		}
	} else {
		context.ReportValueWrongType(ard.TypeString)
		return &self
	}

	matches := VersionRE.FindStringSubmatch(self.CanonicalString)
	length := len(matches)
	if length == 0 {
		context.ReportValueMalformed("version", "")
		return &self
	}

	if length > 1 {
		self.Major = parseVersionUint(matches[1])
	}
	if length > 2 {
		self.Minor = parseVersionUint(matches[2])
	}
	if length > 3 {
		if matches[3] != "" {
			self.Fix = parseVersionUint(matches[3])
		}
	}
	if length > 4 {
		self.Qualifier = matches[4]
	}
	if length > 5 {
		if matches[5] != "" {
			self.Build = parseVersionUint(matches[5])
		}
	}

	return &self
}

// ([fmt.Stringer] interface)
func (self *Version) String() string {
	return self.CanonicalString
}

// HasComparer interface
func (self *Version) SetComparer(comparer string) {
	self.Comparer = comparer
}

func (self *Version) Compare(data any) (int, error) {
	if version, ok := data.(*Version); ok {
		d := CompareUint32(self.Major, version.Major)
		if d != 0 {
			return d, nil
		}
		d = CompareUint32(self.Minor, version.Minor)
		if d != 0 {
			return d, nil
		}
		d = CompareUint32(self.Fix, version.Fix)
		if d != 0 {
			return d, nil
		}
		// Note: the qualifier is compared alphabetically, *not* semantically
		d = strings.Compare(strings.ToLower(self.Qualifier), strings.ToLower(version.Qualifier))
		if d != 0 {
			return d, nil
		}
		return CompareUint32(self.Build, version.Build), nil
	}
	return 0, errors.New("incompatible comparison")
}

func parseVersionUint(value string) uint32 {
	if u, err := strconv.ParseUint(value, 10, 32); err == nil {
		return uint32(u)
	}
	panic("as long as the regexp does its job we should never get here")
}
