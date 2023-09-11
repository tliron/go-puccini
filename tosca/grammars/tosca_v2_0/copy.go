package tosca_v2_0

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/parsing"
)

// Implements the "copy" feature for node templates and relationship templates
func CopyTemplate(context *parsing.Context) {
	if map_, ok := context.Data.(ard.Map); ok {
		if copyFromName_, ok := map_["copy"]; ok {
			if copyFromName, ok := copyFromName_.(string); ok {
				if context.Parent != nil {
					templates_ := context.Parent.Data
					if templates, ok := templates_.(ard.Map); ok {
						if copyFromTemplate, ok := templates[copyFromName]; ok {
							if copied, ok := CopyAndMerge(copyFromTemplate, context.Data, nil, templates); ok {
								context.Data = copied
							} else {
								context.FieldChild("copy", copyFromName).ReportCopyLoop(copyFromName)
							}
						}
					}
				}
			}
		}
	}
}

func CopyAndMerge(target ard.Value, source ard.Value, copiedNames []string, targets ard.Map) (ard.Value, bool) {
	target = ard.Copy(target)

	if targetMap, ok := target.(ard.Map); ok {
		// Recurse target?
		if targetCopyFromName_, ok := targetMap["copy"]; ok {
			if targetCopyFromName, ok := targetCopyFromName_.(string); ok {
				for _, copiedName := range copiedNames {
					if targetCopyFromName == copiedName {
						return nil, false
					}
				}

				if targetCopyFromData, ok := targets[targetCopyFromName]; ok {
					copiedNames = append(copiedNames, targetCopyFromName)
					if target, ok = CopyAndMerge(targetCopyFromData, target, copiedNames, targets); ok {
						if targetMap_, ok := target.(ard.Map); ok {
							targetMap = targetMap_
						}
					} else {
						return nil, false
					}
				}
			}
		}

		if sourceMap, ok := source.(ard.Map); ok {
			ard.Merge(targetMap, sourceMap, false)
		}
	}

	return target, true
}
