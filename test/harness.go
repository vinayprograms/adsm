package test

import (
	"errors"
	"securitymodel/objmodel"
	"securitymodel/yamlmodel"
)

type TestHarness struct {
	YamlStructures map[string]interface{} // YAML structures store for use with 'Resolve' function
}

// Resolve function that will be used for testing 'objmodel' entities.
func (t *TestHarness) Resolve(id string) (interface{}, []error) {
	if value, exists := t.YamlStructures[id]; exists {
		if ent, ok := value.(*yamlmodel.Entity); ok {
			switch ent.Type {
			case yamlmodel.Human:
				var h objmodel.Human
				errs := h.Init(ent, t.Resolve)
				return &h, errs
			case yamlmodel.Role:
				var rol objmodel.Role
				errs := rol.Init(ent, t.Resolve)
				return rol, errs
			case yamlmodel.Program, yamlmodel.System:
				var p objmodel.Program
				errs := p.Init(ent, t.Resolve)
				return &p, errs
			}
		} else if flow, ok := value.(*yamlmodel.Flow); ok {
			var f objmodel.Flow
			errs := f.Init(flow, t.Resolve)
			return &f, errs
		} else {
			return nil, []error{errors.New("test harness found unsupported entity type for '" + id + "'")}
		}
	} else {
		return nil, []error{errors.New("test harness doesn't contain the YAML structure for - '" + id + "'")}
	}

	return nil, []error{errors.New("unknown error in test harness when resolving '" + id + "'")}
}
