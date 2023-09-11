package tosca_v2_0

import (
	"math"

	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// Range
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.3.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.3.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.2.3
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.2.3
//

type Range struct {
	Lower uint64 `json:"lower" yaml:"lower"`
	Upper uint64 `json:"upper" yaml:"upper"`
}

// ([parsing.Reader] signature)
func ReadRange(context *parsing.Context) parsing.EntityPtr {
	var self Range

	if !context.ValidateType(ard.TypeList) {
		return &self
	}

	list := context.Data.(ard.List)
	if len(list) != 2 {
		context.ReportValueMalformed("range", "list length not 2")
		return &self
	}

	lowerContext := context.ListChild(0, list[0])
	lowerOk := false
	if lowerContext.ValidateType(ard.TypeInteger) {
		lowerInt := *lowerContext.ReadInteger()
		if lowerInt < 0 {
			context.ReportValueMalformed("range", "lower bound negative")
		} else {
			self.Lower = uint64(lowerInt)
			lowerOk = true
		}
	}

	upperContext := context.ListChild(1, list[1])
	upperOk := false
	if upperContext.ValidateType(ard.TypeInteger, ard.TypeString) {
		if upperContext.Is(ard.TypeInteger) {
			upperInt := *upperContext.ReadInteger()
			if upperInt < 0 {
				context.ReportValueMalformed("range", "upper bound negative")
			} else {
				self.Upper = uint64(upperInt)
				upperOk = true
			}
		} else if upperContext.Data.(string) == "UNBOUNDED" {
			self.Upper = math.MaxUint64
			upperOk = true
		} else {
			context.ReportValueMalformed("range", "upper bound string not UNBOUNDED")
		}
	}

	// The TOSCA spec in 1.0 to 1.2 says that the upper bound *must* be greater than the lower
	// bound but that makes no sense to us (it would mean than ranges must include at least two
	// numbers) so we are ignoring that and assuming that upper must be >= lower
	// (The spec was fixed in 1.3)
	if upperOk && lowerOk && (self.Upper < self.Lower) {
		context.ReportValueMalformed("range", "upper bound lower than lower bound")
	}

	return &self
}

func (self *Range) InRange(number uint64) bool {
	return (number >= self.Lower) && (number <= self.Upper)
}

//
// RangeEntity
//

type RangeEntity struct {
	*Entity `name:"range"`

	Range *Range `traverse:"ignore"`
}

func NewRangeEntity(context *parsing.Context) *RangeEntity {
	return &RangeEntity{Entity: NewEntity(context)}
}

// ([parsing.Reader] signature)
func ReadRangeEntity(context *parsing.Context) parsing.EntityPtr {
	self := NewRangeEntity(context)
	self.Range = ReadRange(context).(*Range)
	return self
}
