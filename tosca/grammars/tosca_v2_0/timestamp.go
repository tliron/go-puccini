package tosca_v2_0

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/parsing"
)

// See the Timestamp Language-Independent Type for YAML Version 1.1 (Working Draft 2005-01-18)
// http://yaml.org/type/timestamp.html

var TimestampShortRE = regexp.MustCompile(`^(?P<year>[0-9][0-9][0-9][0-9])-(?P<month>[0-9][0-9])-(?P<day>[0-9][0-9])$`)

var TimestampLongRE = regexp.MustCompile(
	`^(?P<year>[0-9][0-9][0-9][0-9])-(?P<month>[0-9][0-9]?)-(?P<day>[0-9][0-9]?)` +
		`(?:[Tt]|[ \t]+)` +
		`(?P<hour>[0-9][0-9]?):(?P<minute>[0-9][0-9]):(?P<second>[0-9][0-9])(?:(?P<fraction>\.[0-9]*))?` +
		`(?:(?:[ \t]*)(?:Z|(?P<tzhour>[-+][0-9][0-9]?)(?::(?P<tzminute>[0-9][0-9]))?))?$`)

const TimestampFormat = "%04d-%02d-%02dT%02d:%02d:%02d%s%s"

const TimestampTimezoneFormat = "%s%02d:%02d"

//
// Timestamp
//

type Timestamp struct {
	CanonicalNumber int64  `json:"$number" yaml:"$number"`
	CanonicalString string `json:"$string" yaml:"$string"`
	OriginalString  string `json:"$originalString" yaml:"$originalString"`

	Year     uint32  `json:"year" yaml:"year"`
	Month    uint32  `json:"month" yaml:"month"`
	Day      uint32  `json:"day" yaml:"day"`
	Hour     uint32  `json:"hour" yaml:"hour"`
	Minute   uint32  `json:"minute" yaml:"minute"`
	Second   uint32  `json:"second" yaml:"second"`
	Fraction float64 `json:"fraction" yaml:"fraction"`
	TZSign   string  `json:"tzSign" yaml:"tzSign"`
	TZHour   uint32  `json:"tzHour" yaml:"tzHour"`
	TZMinute uint32  `json:"tzMinute" yaml:"tzMinute"`
}

// ([parsing.Reader] signature)
func ReadTimestamp(context *parsing.Context) parsing.EntityPtr {
	var self Timestamp

	if context.Is(ard.TypeString) {
		self.OriginalString = *context.ReadString()
		matches := TimestampShortRE.FindStringSubmatch(self.OriginalString)
		length := len(matches)
		if length == 0 {
			matches = TimestampLongRE.FindStringSubmatch(self.OriginalString)
			length = len(matches)
		}
		if length == 0 {
			context.ReportValueMalformed("timestamp", "")
			return &self
		}

		valid := true
		var err error

		if length > 1 {
			if self.Year, err = parseTimestampUint(matches[1]); err != nil {
				context.ReportValueMalformed("timestamp", "cannot parse year")
				valid = false
			}
		}

		if length > 2 {
			if self.Month, err = parseTimestampUint(matches[2]); err != nil {
				context.ReportValueMalformed("timestamp", "cannot parse month")
				valid = false
			} else if (self.Month == 0) || (self.Month > 12) {
				context.ReportValueMalformed("timestamp", "invalid month")
				valid = false
			}
		}

		if length > 3 {
			if self.Day, err = parseTimestampUint(matches[3]); err != nil {
				context.ReportValueMalformed("timestamp", "cannot parse day")
				valid = false
			} else if self.Day > 31 {
				context.ReportValueMalformed("timestamp", "invalid day")
				valid = false
			}
		}

		if length > 4 {
			if self.Hour, err = parseTimestampUint(matches[4]); err != nil {
				context.ReportValueMalformed("timestamp", "cannot parse hour")
				valid = false
			} else if self.Hour > 23 {
				context.ReportValueMalformed("timestamp", "invalid hour")
				valid = false
			}
		}

		if length > 5 {
			if self.Minute, err = parseTimestampUint(matches[5]); err != nil {
				context.ReportValueMalformed("timestamp", "cannot parse minute")
				valid = false
			} else if self.Minute > 59 {
				context.ReportValueMalformed("timestamp", "invalid minute")
				valid = false
			}
		}

		if length > 6 {
			if self.Second, err = parseTimestampUint(matches[6]); err != nil {
				context.ReportValueMalformed("timestamp", "cannot parse second")
				valid = false
			} else if self.Second > 59 {
				context.ReportValueMalformed("timestamp", "invalid second")
				valid = false
			}
		}

		if (length > 7) && (matches[7] != "") {
			if self.Fraction, err = parseTimestampFloat(matches[7]); err != nil {
				context.ReportValueMalformed("timestamp", "cannot parse fraction")
				valid = false
			} else if (self.Fraction < 0.0) || (self.Fraction >= 1.0) {
				context.ReportValueMalformed("timestamp", "invalid fraction")
				valid = false
			}
		}

		if (length > 8) && (matches[8] != "") {
			v := matches[8]
			self.TZSign = v[:1]
			if self.TZHour, err = parseTimestampUint(v[1:]); err != nil {
				context.ReportValueMalformed("timestamp", "cannot parse timezone year")
				valid = false
			}

			// Real-world timezones go from -12 to +14:
			// https://en.wikipedia.org/wiki/List_of_UTC_time_offsets
			// But there's probably no benefit in validating that here, because the math of timezone
			// translation would "just work" even if the specified timezone doesn't exist in the real
			// world
		}

		if (length > 9) && (matches[9] != "") {
			if self.TZMinute, err = parseTimestampUint(matches[9]); err != nil {
				context.ReportValueMalformed("timestamp", "cannot parse timezone minute")
				valid = false
			} else if self.TZMinute > 59 {
				context.ReportValueMalformed("timestamp", "invalid timezone minute")
				valid = false
			}
		}

		if !valid {
			return &self
		}
	} else if context.HasQuirk(parsing.QuirkDataTypesTimestampPermissive) && context.Is(ard.TypeTimestamp) {
		// Note: OriginalString will be empty because it is not preserved by our YAML parser
		time := context.Data.(time.Time)
		_, tzSeconds := time.Zone()

		self.Year = uint32(time.Year())
		self.Month = uint32(time.Month())
		self.Day = uint32(time.Day())
		self.Hour = uint32(time.Hour())
		self.Minute = uint32(time.Minute())
		self.Second = uint32(time.Second())
		self.Fraction = float64(time.Nanosecond()) / 1_000_000_000.0
		if tzSeconds >= 0 {
			self.TZSign = "+"
		} else {
			self.TZSign = "-"
			tzSeconds = -tzSeconds
		}
		self.TZHour = uint32(tzSeconds / 3600)
		self.TZMinute = uint32((tzSeconds % 3600) / 60)
	} else {
		if context.HasQuirk(parsing.QuirkDataTypesTimestampPermissive) {
			context.ReportValueWrongType(ard.TypeString, ard.TypeTimestamp)
		} else {
			context.ReportValueWrongType(ard.TypeString)
		}
		return &self
	}

	// Canonical string
	var tz string
	if (self.TZHour == 0) && (self.TZMinute == 0) {
		tz = "Z"
	} else {
		tz = fmt.Sprintf(TimestampTimezoneFormat, self.TZSign, self.TZHour, self.TZMinute)
	}
	fraction := fmt.Sprintf("%g", self.Fraction)[1:]
	self.CanonicalString = fmt.Sprintf(TimestampFormat, self.Year, self.Month, self.Day, self.Hour, self.Minute, self.Second, fraction, tz)

	// Canonical number is nanoseconds since (or before if negative) Jan 1 1970 UTC
	self.CanonicalNumber = int64(self.Time().UnixNano())

	return &self
}

// ([fmt.Stringer] interface)
func (self *Timestamp) String() string {
	return self.CanonicalString
}

func (self *Timestamp) Compare(data any) (int, error) {
	if timestamp, ok := data.(*Timestamp); ok {
		return CompareInt64(self.CanonicalNumber, timestamp.CanonicalNumber), nil
	}
	return 0, errors.New("incompatible comparison")
}

// Convert timezone to Go time.Location
func (self *Timestamp) Location() *time.Location {
	seconds := int(self.TZHour*3600 + self.TZMinute*60)
	if self.TZSign == "-" {
		seconds = -seconds
	}
	return time.FixedZone("", seconds)
}

// Convert to Go time.Time
func (self *Timestamp) Time() time.Time {
	return time.Date(
		int(self.Year),
		time.Month(self.Month),
		int(self.Day),
		int(self.Hour),
		int(self.Minute),
		int(self.Second),
		int(self.Fraction*1_000_000_000.0),
		self.Location(),
	)
}

// Utils

func parseTimestampUint(value string) (uint32, error) {
	if u, err := strconv.ParseUint(value, 10, 32); err == nil {
		return uint32(u), nil
	} else {
		return 0, err
	}
}

func parseTimestampFloat(value string) (float64, error) {
	if u, err := strconv.ParseFloat(value, 64); err == nil {
		return u, nil
	} else {
		return 0.0, err
	}
}
