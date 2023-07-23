package objmodel

import (
	"errors"
	"securitymodel/yamlmodel"
)

type Role struct {
	CoreObject
}

func (rol *Role) Init(e *yamlmodel.Entity, r Resolver) []error {
	var errs []error

	if e == nil {
		return []error{errors.New("cannot convert nil yaml to role specification")}
	}

	if e.Type != yamlmodel.Role {
		return []error{errors.New("cannot initialize role with '" + string(e.Type) + "'")}
	}

	err := rol.SetID(e.Id)
	if err != nil {
		return []error{err}
	}
	err = rol.SetName(e.Name)
	if err != nil {
		return []error{err}
	}
	err = rol.SetDescription(e.Description)
	if err != nil {
		return []error{err}
	}

	if e.AdmDir != "" {
		for _, adm := range e.ADM {
			adm = e.AdmDir + "/" + adm
			rol.AddADM(adm)
		}
	} else {
		rol.SetADM(e.ADM)
	}

	rol.SetMitigations(e.Mitigations)
	rol.SetRecommendations(e.Recommendations)

	return errs
}

// Collect all ADMs from program
func (rol *Role) GetADM() (allADM map[string][]string) {
	allADM = make(map[string][]string)
	allADM[rol.id] = rol.adm
	return
}

func (rol *Role) AddADM(adm string) error {
	rol.adm = append(rol.adm, adm)
	return nil
}

func (rol *Role) GetMitigations() map[string][]string {
	allMitigations := make(map[string][]string)
	if len(rol.mitigations) > 0 {
		allMitigations[rol.GetName()] = append(allMitigations[rol.GetName()], rol.mitigations...)
	}
	return allMitigations
}

func (rol *Role) GetRecommendations() map[string][]string {
	allRecommendations := make(map[string][]string)
	if len(rol.recommendations) > 0 {
		allRecommendations[rol.GetName()] = append(allRecommendations[rol.GetName()], rol.recommendations...)
	}
	return allRecommendations
}

func (rol *Role) AddMitigation(mitigation string) error {
	rol.mitigations = append(rol.mitigations, mitigation)
	return nil
}

func (rol *Role) AddRecommendation(reco string) error {
	rol.recommendations = append(rol.recommendations, reco)
	return nil
}
