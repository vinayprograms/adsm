package loaders

import (
	"addb"
	"errors"
	"securitymodel/objmodel"
	"securitymodel/yamlmodel"
	"strings"
)

func (t *Builder) ResolveYaml(id string, basePath string, addb *addb.ADDB) (interface{}, []error) {
	t.init()                                   // Initialize if not done already
	if _, exists := t.yamlIndex[id]; !exists { // yaml hasn't been indexed yet
		if strings.HasPrefix(id, "addb:") { // if entity is from ADDB
			yamlObj, errs := t.readandIndexYamlFromADDB(id, addb)
			if yamlObj == nil {
				errs = append(errs, errors.New("Entity '"+id+"' not found in ADDB"))
				return nil, errs
			}
			return yamlObj, nil
		} else {
			return nil, []error{errors.New("Entity '" + id + "' not found in model or ADDB")}
		}
	}

	return t.yamlIndex[id], nil
}

// If object is already indexed, return it. If not, check if
// its yaml is indexed. If so, build object and return it. If
// neither exist, return nil.
func (t *Builder) Resolve(id string) (interface{}, []error) {
	t.init() // Initialize if not done already

	if obj, exists := t.objectIndex[id]; exists { // object alredy indexed
		return obj, nil
	} else if yamlObj, exists := t.yamlIndex[id]; exists { // object has to be built from indexed yaml data
		// Build object, index it and return the built object.
		obj, errs := t.buildObjectFromYaml(yamlObj)
		if obj == nil {
			errs = append(errs, errors.New("unknown error when building object from yaml for '"+id+"'"))
		}
		t.objectIndex[id] = obj
		return obj, errs
	} else { // the YAML you are looking for is not indexed
		return nil, []error{errors.New("no YAML found for '" + id + "' in model or ADDB")}
	}
}

func (t *Builder) buildObjectFromYaml(yamlObj interface{}) (interface{}, []error) {
	if ent, ok := yamlObj.(*yamlmodel.Entity); ok {
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
	} else if flow, ok := yamlObj.(*yamlmodel.Flow); ok {
		var f objmodel.Flow
		errs := f.Init(flow, t.Resolve)
		return &f, errs
	}

	return nil, nil
}
