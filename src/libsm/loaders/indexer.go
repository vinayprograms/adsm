package loaders

import (
	"errors"
	"libaddb"
	"libsm/objmodel"
	"libsm/yamlmodel"
	"strings"
)

func (b *Builder) Index(id string, obj interface{}, admDir string, addb *libaddb.ADDB) (errs []error) {
	b.init() // Initialize if not done already

	switch x := obj.(type) {
	case *yamlmodel.SecurityModel:
		e := b.indexModelParts(x, admDir, addb)
		errs = append(errs, e...)
	case *yamlmodel.Entity:
		if _, ok := b.yamlIndex[id]; !ok {
			b.yamlIndex[id] = obj
		}
		e := b.indexEntityParts(id, admDir, addb, x)
		errs = append(errs, e...)
	case *yamlmodel.Flow:
		if _, ok := b.yamlIndex[id]; !ok {
			b.yamlIndex[id] = obj
		}
		e := b.indexFlowParts(id, admDir, addb, x)
		errs = append(errs, e...)

	case *objmodel.SecurityModel, *objmodel.Human, *objmodel.Program, *objmodel.Flow:
		if _, ok := b.objectIndex[id]; ok {
			if b.objectIndex[id] == nil {
				b.objectIndex[id] = obj
				return
			} else {
				return []error{errors.New("multiple objects with name '" + id + "' found")}
			}
		} else {
			b.objectIndex[id] = obj
			return
		}
	default:
		return
	}

	return
}

////////////////////////////////////////
// internal functions

func (b *Builder) indexModelParts(m *yamlmodel.SecurityModel, admDir string, addb *libaddb.ADDB) (errors []error) {
	// Index all YAML entities from model first. 
	for _, obj := range m.Externals {
		if obj == nil || obj.Id == "" { continue }
		b.yamlIndex[obj.Id] = obj
	}
	for _, obj := range m.Entities { 
		if obj == nil || obj.Id == "" { continue }
		b.yamlIndex[obj.Id] = obj
	}
	for _, obj := range m.Flows { 
		if obj == nil || obj.Id == "" { continue }
		b.yamlIndex[obj.Id] = obj
	}

	m.AdmDir = admDir

	// Index parts of each entity. This is required by resolver later
	// to find yaml objects that are referred in other places in the doc.
	for _, obj := range m.Externals {
		if obj == nil || obj.Id == "" { continue }
		errs := b.indexEntityParts(obj.Id, admDir, addb, obj)
		errors = append(errors, errs...)
	}
	for _, obj := range m.Entities { 
		if obj == nil || obj.Id == "" { continue }
		errs := b.indexEntityParts(obj.Id, admDir, addb, obj) 
		errors = append(errors, errs...)
	}
	for _, obj := range m.Flows {
		if obj == nil || obj.Id == "" { continue }
		errs := b.indexFlowParts(obj.Id, admDir, addb, obj)
		errors = append(errors, errs...)
	}

	return errors
}

func (b *Builder) indexEntityParts(id string, basePath string, addb *libaddb.ADDB, entity *yamlmodel.Entity) []error {
	var allErrors []error

	entity.AdmDir = basePath

	// Index members that refer to other model objects
	if entity.Base != nil && len(entity.Base) > 0 { 
		for _, base := range entity.Base {
			_, err := b.ResolveYaml(base, basePath, addb)
			if err != nil { allErrors = append(allErrors, err...) }
		}
	}
	if entity.Interface != "" {
		_, err := b.ResolveYaml(entity.Interface, basePath, addb)
		if err != nil { allErrors = append(allErrors, err...) }
	}
	if len(entity.Roles) > 0 {
		for _, role := range entity.Roles {
			if role != "" { 
				_, err := b.ResolveYaml(role, basePath, addb)
				if err != nil { allErrors = append(allErrors, err...) }
			}
		}
	}
	if len(entity.Languages) > 0 {
		for _, lang := range entity.Languages {
			if lang != "" { 
				_, err := b.ResolveYaml(lang, basePath, addb)
				if err != nil { allErrors = append(allErrors, err...) }
			}
		}
	}
	if len(entity.Dependencies) > 0 {
		for _, dep := range entity.Dependencies {
			if dep != "" { 
				_, err := b.ResolveYaml(dep, basePath, addb)
				if err != nil { allErrors = append(allErrors, err...) }
			}
		}
	}

	if len(allErrors) > 0 {
		return allErrors
	} else {
		return nil
	}
}

func (b *Builder) indexFlowParts(id string, basePath string, addb *libaddb.ADDB, flow *yamlmodel.Flow) []error {
	var allErrors []error

	flow.AdmDir = basePath

	// Index members that refer to other model objects
	if flow.Sender != "" { 
		_, err := b.ResolveYaml(flow.Sender, basePath, addb)
		if err != nil { allErrors = append(allErrors, err...) }
	}
	if flow.Receiver != "" {
		_, err := b.ResolveYaml(flow.Receiver, basePath, addb)
		if err != nil { allErrors = append(allErrors, err...) }
	}
	if len(flow.Protocol) > 0 {
		for _, proto := range flow.Protocol {
			if proto != "" { 
				_, err := b.ResolveYaml(proto, basePath, addb)
				if err != nil { allErrors = append(allErrors, err...) }
			}
		}
	}

	if len(allErrors) > 0 {
		return allErrors
	} else {
		return nil
	}
}

func (b *Builder) readandIndexYamlFromADDB(id string, addb *libaddb.ADDB) (interface{}, []error) {
	component, err := addb.GetComponent(id)
	if err != nil {
		return nil, nil
	}

	switch strings.ToLower(string(component.Type)) {
	case "human", "program", "system":
		entity, err := b.translateADDBComponentToEntity(component)
		if err != nil { 
			return nil, nil
		}
		errs := b.Index(id, entity, "", addb)
		return entity, errs
	case "flow":
		flow, err := b.translateADDBComponentToFlow(component)
		if err != nil { 
			return nil, nil 
		}
		flow.AdmDir = addb.Location
		errs := b.Index(id, flow, "", addb)
		return flow, errs
	}

	return nil, nil
}

func (t *Builder) translateADDBComponentToEntity(component *libaddb.ADDBComponent) (*yamlmodel.Entity, error) {
	var entity yamlmodel.Entity

	if component == nil {
		return nil, errors.New("cannot translate null ADDB component")
	}

	entity.Id = component.Id
	entity.Name = component.Name
	switch component.Type {
	case libaddb.Human:		entity.Type = yamlmodel.Human
	case libaddb.Program:	entity.Type = yamlmodel.Program
	case libaddb.System:	entity.Type = yamlmodel.System

	}
	entity.Description = component.Description
	entity.Base = component.Base
	entity.Recommendations = component.Recommendations
	entity.ADM = component.ADM
	
	// Only for humans
	entity.Roles = component.Roles
	entity.Interface = component.Interface

	// Only for programs
	entity.CodeRepository = component.CodeRepository
	entity.Languages = component.Languages
	entity.Dependencies = component.Dependencies

	return &entity, nil
}

func (db *Builder) translateADDBComponentToFlow(component *libaddb.ADDBComponent) (*yamlmodel.Flow, error) {
	var flow yamlmodel.Flow

	if component == nil {
		return nil, errors.New("cannot translate null ADDB component")
	}

	flow.Id = component.Id
	flow.Name = component.Name
	flow.Description = component.Description
	flow.Protocol = component.Protocol
	
	// NOTE: Sender & Receiver are available for an ADDB component

	flow.Recommendations = component.Recommendations
	flow.ADM = component.ADM

	return &flow, nil
}