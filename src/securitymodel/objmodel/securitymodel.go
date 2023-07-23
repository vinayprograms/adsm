package objmodel

import (
	"errors"
	"securitymodel/yamlmodel"
)

type SecurityModel struct {
	Title          string
	DesignDocument string
	AddbPath       string
	modelADM       []string
	Externals      map[string]ExternalSpec
	Entities       map[string]EntitySpec
	Flows          map[string]FlowSpec
}

// Collect all ADMs from program
func (t *SecurityModel) GetADM() (allADM map[string][]string) {
	allADM = make(map[string][]string)
	allADM["sm"] = t.modelADM

	if t.Entities != nil {
		for _, e := range t.Entities {
			if _, ok := e.(*Role); !ok {
				allADM = merge(allADM, "sm.entities", e.GetADM())
			}
		}
	}
	if t.Flows != nil {
		for _, f := range t.Flows {
			allADM = merge(allADM, "sm.flows", f.GetADM())
		}
	}

	return
}

// Collect all ADMs from program
func (t *SecurityModel) SetADM(adm []string) {
	t.modelADM = append(t.modelADM, adm...)
}

func (t *SecurityModel) Init(ysm *yamlmodel.SecurityModel, r Resolver) []error {

	var errs []error

	if ysm == nil {
		errs = append(errs, errors.New("cannot convert nil yaml to security-model"))
		return errs
	}

	t.Title = ysm.Title
	t.DesignDocument = ysm.DesignDocument
	t.AddbPath = ysm.AddbUri

	if ysm.AdmDir != "" {
		for _, adm := range ysm.ModelADM {
			adm = ysm.AdmDir + "/" + adm
			t.modelADM = append(t.modelADM, adm)
		}
	} else {
		t.modelADM = append(t.modelADM, ysm.ModelADM...)
	}

	// Build SM objects from indexed YAML content
	buildErrs := t.buildExternalEntities(ysm.Externals, r)
	if len(buildErrs) != 0 {
		errs = append(errs, buildErrs...)
	}
	buildErrs = t.buildEntities(ysm.Entities, r)
	if len(buildErrs) != 0 {
		errs = append(errs, buildErrs...)
	}
	buildErrs = t.buildFlows(ysm.Flows, r)
	if len(buildErrs) != 0 {
		errs = append(errs, buildErrs...)
	}

	return errs
}

func (t *SecurityModel) buildExternalEntities(externals []*yamlmodel.Entity, r Resolver) []error {
	var errs []error

	t.Externals = make(map[string]ExternalSpec)
	for _, entry := range externals {
		if entry == nil || entry.Id == "" {
			continue
		}
		if entry.Type == yamlmodel.Human {
			obj, extErrs := r(entry.Id)
			if len(extErrs) != 0 {
				errs = append(errs, extErrs...)
			}
			h, ok := obj.(*Human)
			if !ok { // control shouldn't reach this section. If it does, contact author!
				return []error{errors.New("error in creating human - " + entry.Id)}
			}
			t.Externals[h.GetID()] = h
		} else {
			obj, extErrs := r(entry.Id)
			if len(extErrs) != 0 {
				errs = append(errs, extErrs...)
			}
			p, ok := obj.(*Program)
			if !ok { // control shouldn't reach this section. If it does, contact author!
				return []error{errors.New("error in creating program - " + entry.Id)}
			}
			t.Externals[p.GetID()] = p
		}
	}
	return errs
}

func (t *SecurityModel) buildEntities(entities []*yamlmodel.Entity, r Resolver) []error {
	var errs []error

	t.Entities = make(map[string]EntitySpec)
	for _, entry := range entities {
		if entry == nil || entry.Id == "" {
			continue
		}
		switch entry.Type {
		case yamlmodel.Human:
			obj, entityErrs := r(entry.Id)
			if len(entityErrs) != 0 {
				errs = append(errs, entityErrs...)
			}
			h, ok := obj.(*Human)
			if !ok { // control shouldn't reach this section. If it does, contact author!
				return []error{errors.New("error in creating human - " + entry.Id)}
			}
			t.Entities[h.GetID()] = h
		case yamlmodel.Program, yamlmodel.System:
			obj, entityErrs := r(entry.Id)
			if len(entityErrs) != 0 {
				errs = append(errs, entityErrs...)
			}
			p, ok := obj.(*Program)
			if !ok {
				return []error{errors.New("error in creating program - " + entry.Id)}
			}

			t.Entities[p.GetID()] = p
		case yamlmodel.Role:
			obj, entityErrs := r(entry.Id)
			if len(entityErrs) != 0 {
				errs = append(errs, entityErrs...)
			}
			rol, ok := obj.(Role)
			if !ok {
				return []error{errors.New("error in creating program - " + entry.Id)}
			}

			t.Entities[rol.GetID()] = &rol
		}
	}
	return errs
}

func (t *SecurityModel) buildFlows(flows []*yamlmodel.Flow, r Resolver) []error {
	var errs []error

	t.Flows = make(map[string]FlowSpec)
	for _, entry := range flows {
		if entry == nil || entry.Id == "" {
			continue
		}
		obj, flowErrs := r(entry.Id)
		if len(flowErrs) != 0 {
			errs = append(errs, flowErrs...)
		}
		f, ok := obj.(*Flow)
		if !ok { // control shouldn't reach this section. If it does, contact author!
			return []error{errors.New("error in creating flow - " + entry.Id)}
		}
		t.Flows[f.GetID()] = f
	}
	return errs
}
