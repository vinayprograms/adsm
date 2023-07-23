package objmodel

import (
	"errors"
	"securitymodel/yamlmodel"
)

type Human struct {
	CoreObject
	base          []HumanEntitySpec
	userInterface ProgramSpec
}

func (h *Human) Init(e *yamlmodel.Entity, r Resolver) []error {
	var errs []error

	if e == nil {
		return []error{errors.New("cannot convert nil yaml to human specification")}
	}

	if e.Type != yamlmodel.Human {
		return []error{errors.New("cannot initialize human with '" + string(e.Type) + "'")}
	}

	err := h.SetID(e.Id)
	if err != nil {
		return []error{err}
	}
	err = h.SetName(e.Name)
	if err != nil {
		return []error{err}
	}
	err = h.SetDescription(e.Description)
	if err != nil {
		return []error{err}
	}

	if e.AdmDir != "" {
		for _, adm := range e.ADM {
			adm = e.AdmDir + "/" + adm
			h.AddADM(adm)
		}
	} else {
		h.SetADM(e.ADM)
	}

	h.SetMitigations(e.Mitigations)
	h.SetRecommendations(e.Recommendations)

	if e.Base != nil && len(e.Base) > 0 {
		for _, base := range e.Base {
			obj, baseErrs := r(base)
			if len(baseErrs) != 0 {
				errs = append(errs, baseErrs...)
			}
			if b, ok := obj.(HumanEntitySpec); ok {
				h.AddBase(b)
			} else {
				errs = append(errs, errors.New("error in resolving base '"+base+"' for human '"+h.id+"'"))
			}
		}
	}

	if e.Interface != "" {
		obj, ifaceErrs := r(e.Interface)
		if len(ifaceErrs) != 0 {
			errs = append(errs, ifaceErrs...)
		}
		if iface, ok := obj.(ProgramEntitySpec); ok {
			h.SetUserInterface(iface)
		} else {
			errs = append(errs, errors.New("error in resolving interface '"+e.Interface+"' for human '"+h.id+"'"))
		}
	}

	return errs
}

func (h *Human) AddADM(adm string) error {
	h.adm = append(h.adm, adm)
	return nil
}

// Collect all ADMs from program
func (h *Human) GetADM() (allADM map[string][]string) {
	allADM = make(map[string][]string)
	allADM[h.id] = h.adm
	if h.base != nil && len(h.base) > 0 {
		for _, base := range h.base {
			allADM = merge(allADM, h.id+".base", base.GetADM())
		}
	}
	if h.userInterface != nil {
		if _, ok := h.userInterface.(EntitySpec); !ok {
			return
		}
		allADM = merge(allADM, h.id+".interface", h.userInterface.(EntitySpec).GetADM())
	}

	return
}

func (h *Human) GetMitigations() map[string][]string {
	allMitigations := make(map[string][]string)
	if len(h.mitigations) > 0 {
		allMitigations[h.GetName()] = append(allMitigations[h.GetName()], h.mitigations...)
	}
	if h.base != nil && len(h.base) > 0 {
		for _, base := range h.base {
			for key, mitigations := range base.GetMitigations() {
				if len(mitigations) > 0 {
					allMitigations[h.GetName()+" -> (Base)"+base.GetName()+":"+key] = append(allMitigations[h.GetName()+" -> (Base)"+base.GetName()+":"+key], mitigations...)
				}
			}
		}
	}
	return allMitigations
}

func (h *Human) GetRecommendations() map[string][]string {
	allRecommendations := make(map[string][]string)
	if len(h.recommendations) > 0 {
		allRecommendations[h.GetName()] = append(allRecommendations[h.GetName()], h.recommendations...)
	}
	if h.base != nil && len(h.base) > 0 {
		for _, base := range h.base {
			for key, recos := range base.GetRecommendations() {
				if len(recos) > 0 {
					allRecommendations[h.GetName()+" -> (Base)"+base.GetName()+":"+key] = append(allRecommendations[h.GetName()+" -> (Base)"+base.GetName()+":"+key], recos...)
				}
			}
		}
	}
	return allRecommendations
}

func (h *Human) AddMitigation(mitigation string) error {
	h.mitigations = append(h.mitigations, mitigation)
	return nil
}

func (h *Human) AddRecommendation(reco string) error {
	h.recommendations = append(h.recommendations, reco)
	return nil
}

func (h *Human) GetUserInterface() ProgramSpec {
	return h.userInterface
}

func (h *Human) SetUserInterface(p ProgramSpec) error {
	h.userInterface = p
	return nil
}

func (h *Human) GetBase() []HumanEntitySpec {
	return h.base
}

func (h *Human) SetBase(spec []HumanEntitySpec) error {
	h.base = spec
	return nil
}

func (h *Human) AddBase(spec HumanEntitySpec) error {
	h.base = append(h.base, spec)
	return nil
}
