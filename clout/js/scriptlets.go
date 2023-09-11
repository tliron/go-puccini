package js

import (
	"fmt"
	"strings"

	"github.com/tliron/go-ard"
	cloutpkg "github.com/tliron/puccini/clout"
)

func CleanupScriptlet(scriptlet string) string {
	scriptlet = strings.TrimSpace(scriptlet)
	scriptlet = strings.Replace(scriptlet, "\t", "  ", -1)
	return scriptlet
}

func GetScriptlet(name string, clout *cloutpkg.Clout) (string, error) {
	section, err := GetScriptletsMetadataSection(name, clout)
	if err != nil {
		return "", err
	}

	scriptlet, ok := section.(string)
	if !ok {
		return "", NewScriptletNotFoundError("scriptlet found in metadata but not a string: %s", name)
	}

	return scriptlet, nil
}

func SetScriptlet(name string, scriptlet string, clout *cloutpkg.Clout) error {
	metadata, err := GetScriptletsMetadata(clout)
	if err != nil {
		return err
	}

	if ard.With(metadata).ForceGetPath(name, ".").Set(scriptlet) {
		return nil
	} else {
		return fmt.Errorf("could not set scriptlet at path %q", name)
	}
}

func GetScriptletNamesInSection(baseName string, clout *cloutpkg.Clout) ([]string, error) {
	section, err := GetScriptletsMetadataSection(baseName, clout)
	if err != nil {
		return nil, err
	}

	scriptlets, ok := section.(ard.StringMap)
	if !ok {
		return nil, NewScriptletNotFoundError("scriptlet metadata found but not a map: %s", baseName)
	}

	list := make([]string, 0, len(scriptlets))
	for name := range scriptlets {
		list = append(list, fmt.Sprintf("%s.%s", baseName, name))
	}

	return list, nil
}

func GetScriptletsMetadataSection(name string, clout *cloutpkg.Clout) (ard.Value, error) {
	segments, final, err := parseScriptletName(name)
	if err != nil {
		return nil, err
	}

	metadata, err := GetScriptletsMetadata(clout)
	if err != nil {
		return nil, err
	}

	m := metadata
	for _, s := range segments {
		o := m[s]
		var ok bool
		if m, ok = o.(ard.StringMap); !ok {
			return nil, NewScriptletNotFoundError("scriptlet metadata not found: %s", name)
		}
	}

	section, ok := m[final]
	if !ok {
		return nil, NewScriptletNotFoundError("scriptlet metadata not found: %s", name)
	}

	return section, nil
}

func GetScriptletsMetadata(clout *cloutpkg.Clout) (ard.StringMap, error) {
	// TODO: check that version=1.0
	if scriptlets, ok := ard.With(clout.Metadata).Get("puccini", "scriptlets").StringMap(); ok {
		return scriptlets, nil
	} else {
		return nil, NewScriptletNotFoundError("%s", "no \"puccini.scriptlets\" metadata in Clout")
	}
}

//
// ScriptletNotFoundError
//

type ScriptletNotFoundError struct {
	string
}

func NewScriptletNotFoundError(format string, args ...any) *ScriptletNotFoundError {
	return &ScriptletNotFoundError{fmt.Sprintf(format, args...)}
}

func IsScriptletNotFoundError(err error) bool {
	_, ok := err.(*ScriptletNotFoundError)
	return ok
}

// (error interface)
func (self *ScriptletNotFoundError) Error() string {
	return self.string
}

// Utils

func parseScriptletName(name string) ([]string, string, error) {
	segments := strings.Split(name, ".")
	last := len(segments) - 1
	if last == -1 {
		return nil, "", fmt.Errorf("malformed scriptlet name: %s", name)
	}

	return segments[:last], segments[last], nil
}
