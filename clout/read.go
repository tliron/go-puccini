package clout

import (
	"fmt"
	"io"

	"github.com/tliron/go-ard"
)

func Read(reader io.Reader, format string) (*Clout, error) {
	if data, _, err := ard.Read(reader, format, false); err == nil {
		if map_, ok := data.(ard.Map); ok {
			if clout, err := Unpack(map_); err == nil {
				return clout, nil
			} else {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("not a map: %T", data)
		}
	} else {
		return nil, err
	}
}
