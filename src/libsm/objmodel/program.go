package objmodel

import (
	"errors"
	"libsm/yamlmodel"
	"net/url"
)

type Program struct {
	CoreObject
	base []ProgramEntitySpec
	codeRepository string
	roles map[string]EntitySpec
	languages map[string]ProgramEntitySpec
	dependencies map[string]ProgramEntitySpec
}

func (p *Program) Init(e *yamlmodel.Entity, r Resolver) []error {
	var errs []error
	
	if e == nil {
		return []error{errors.New("cannot convert nil yaml to program specification")}
	}

	if e.Type != yamlmodel.Role && e.Type != yamlmodel.Program && e.Type != yamlmodel.System {
		return []error{errors.New("cannot initialize program with '" + string(e.Type) + "'")}
	}

	err := p.SetID(e.Id)
	if err != nil { return []error{err} }
	err = p.SetName(e.Name)
	if err != nil { return []error{err} }
	err = p.SetDescription(e.Description)
	if err != nil { return []error{err} }

	if e.AdmDir != "" {
		for _, adm := range e.ADM {
			adm = e.AdmDir + "/" + adm
			p.AddADM(adm)
		}
	} else {
		p.SetADM(e.ADM)
	}

	p.SetMitigations(e.Mitigations)
	p.SetRecommendations(e.Recommendations)

	err = p.SetRepository(e.CodeRepository)
	if err != nil {
		errs = append(errs, err)
	}

	if e.Base != nil && len(e.Base) > 0 {
		for _, base := range e.Base {
			obj, baseErrs := r(base)
			if len(baseErrs) != 0 {
				errs = append(errs, baseErrs...)
			}
			if b, ok := obj.(ProgramEntitySpec); ok {
				p.AddBase(b)
			} else {
				errs = append(errs, errors.New("error in resolving base '" + base + "' for entity '" + p.id + "'"))
			}
		}
	}

	for _, role := range e.Roles {
		obj, roleErrs := r(role)
		if len(roleErrs) != 0 {
			errs = append(errs, roleErrs...)
		}
		if rol, ok := obj.(Role); ok {
			err := p.AddRole(role, &rol)
			if err != nil { errs = append(errs, err)}
		} else {
			errs = append(errs, errors.New("error in resolving role '" + role + "' for entity '" + p.id + "'"))
		}
	}

	for _, lang := range e.Languages {
		obj, langErrs := r(lang)
		if len(langErrs) != 0 {
			errs = append(errs, langErrs...)
		}
		if l, ok := obj.(*Program); ok {
			err := p.AddLanguage(lang, l)
			if err != nil { errs = append(errs, err)}
		} else {
			errs = append(errs, errors.New("error in resolving language '" + lang + "' for entity '" + p.id + "'"))
		}
	}

	for _, dep := range e.Dependencies {
		obj, depErrs := r(dep)
		if len(depErrs) != 0 {
			errs = append(errs, depErrs...)
		}
		if d, ok := obj.(*Program); ok {
			err := p.AddDependency(d.id, d)
			if err != nil { errs = append(errs, err) }
		} else {
			errs = append(errs, errors.New("error in resolving dependency '" + dep + "' for entity '" + p.id + "'"))
		}
	}

	return errs
}

// Collect all ADMs from program
func (p *Program) GetADM() (allADM map[string][]string) {
	allADM = make(map[string][]string)
	allADM[p.id] = p.adm
	if p.base != nil && len(p.base) > 0 {
		for _, base := range p.base {
			allADM = merge(allADM, p.id + ".base", base.GetADM())
		}
	}
	if p.roles != nil {
		for _, r := range p.roles {
			allADM = merge(allADM, p.id + ".roles", r.GetADM())
		}
	}
	if p.dependencies != nil {
		for _, d := range p.dependencies {
			allADM = merge(allADM, p.id + ".dependencies", d.GetADM())
		}
	}
	if p.languages != nil {
		for _, l := range p.languages {
			allADM = merge(allADM, p.id + ".languages", l.GetADM())
		}
	}

	return
}

func (p *Program) AddADM(adm string) error {
	p.adm = append(p.adm, adm)
	return nil
}

func (p *Program) GetMitigations() map[string][]string {
	allMitigations := make(map[string][]string)
	if len(p.mitigations) > 0 {
		allMitigations[p.GetName()] = append(allMitigations[p.GetName()], p.mitigations...)
	}
	
	if p.base != nil && len(p.base) > 0 {
		for _, base := range p.base {
			for key, mitigations := range base.GetMitigations() {
				if len(mitigations) > 0 {
					allMitigations[p.GetName() + " -> (Base)" + base.GetName() + ":" + key] = append(allMitigations[p.GetName() + " -> (Base)" + base.GetName() + ":" + key], mitigations...)
				}
			}
		}
	}
	if p.roles != nil {
		for _, role := range p.roles {
			for key, mitigations := range role.GetMitigations() {
				if len(mitigations) > 0 {
					allMitigations[p.GetName() + " -> Role:" + key] = append(allMitigations[p.GetName() + " -> Role:" + key], mitigations...)
				}
			}
		}
	}
	if p.languages != nil {
		for _, lang := range p.languages {
			for key, mitigations := range lang.GetMitigations() {
				if len(mitigations) > 0 {
					allMitigations[p.GetName() + " -> Language:" + key] = append(allMitigations[p.GetName() + " -> Language:" + key], mitigations...)
				}
			}
		}
	}
	if p.dependencies != nil {
		for _, dep := range p.dependencies {
			for key, mitigations := range dep.GetMitigations() {
				if len(mitigations) > 0 {
					allMitigations[p.GetName() + " -> Dependency:" + key] = append(allMitigations[p.GetName() + " -> Dependency:" + key], mitigations...)
				}
			}
		}
	}
	
	return allMitigations
}

func (p *Program) GetRecommendations() map[string][]string {
	allRecommendations := make(map[string][]string)
	if len(p.recommendations) > 0 {
		allRecommendations[p.GetName()] = append(allRecommendations[p.GetName()], p.recommendations...)
	}
	
	if p.base != nil && len(p.base) > 0 {
		for _, base := range p.base {
			for key, recos := range base.GetRecommendations() {
				if len(recos) > 0 {
					allRecommendations[p.GetName() + " -> (Base)" + base.GetName() + ":" + key] = append(allRecommendations[p.GetName() + " -> (Base)" + base.GetName() + ":" + key], recos...)
				}
			}
		}
	}
	if p.roles != nil {
		for _, role := range p.roles {
			for key, recos := range role.GetRecommendations() {
				if len(recos) > 0 {
					allRecommendations[p.GetName() + " -> Role:" + key] = append(allRecommendations[p.GetName() + " -> Role:" + key], recos...)
				}
			}
		}
	}
	if p.languages != nil {
		for _, lang := range p.languages {
			for key, recos := range lang.GetRecommendations() {
				if len(recos) > 0 {
					allRecommendations[p.GetName() + " -> Language:" + key] = append(allRecommendations[p.GetName() + " -> Language:" + key], recos...)
				}
			}
		}
	}
	if p.dependencies != nil {
		for _, dep := range p.dependencies {
			for key, recos := range dep.GetRecommendations() {
				if len(recos) > 0 {
					allRecommendations[p.GetName() + " -> Dependency:" + key] = append(allRecommendations[p.GetName() + " -> Dependency:" + key], recos...)
				}
			}
		}
	}
	
	return allRecommendations
}

func (p *Program) AddMitigation(mitigation string) error {
	p.mitigations = append(p.mitigations, mitigation)
	return nil
}

func (p *Program) AddRecommendation(reco string) error {
	p.recommendations = append(p.recommendations, reco)
	return nil
}

func (p *Program) GetRoles() map[string]EntitySpec {
	return p.roles
}

func (p *Program) AddRole(id string, role EntitySpec) error {
	if p.roles == nil { p.roles = make(map[string]EntitySpec) }
	if _, present := p.roles[id]; present {
		return errors.New("Role '" + id + "' is already part of roles list for '" + p.id + "'")
	} else {
		p.roles[id] = role
	}
	return nil
}

func (p *Program) GetCodeRepository() string {
	return p.codeRepository
}

func (p *Program) SetRepository(urlString string)	error {
	_, err := url.Parse(urlString)
	if err != nil {
		return err
	}
	p.codeRepository = urlString
	return nil
}

func (p *Program) GetBase() []ProgramEntitySpec {
	return p.base
}

func (p *Program) SetBase(spec []ProgramEntitySpec) error {
	p.base = spec
	return nil
}

func (p *Program) AddBase(spec ProgramEntitySpec) error {
	p.base = append(p.base, spec)
	return nil
}

func (p *Program) GetLanguages() map[string]ProgramEntitySpec {
	return p.languages
}

func (p* Program) AddLanguage(id string, language ProgramEntitySpec) error {
	if p.languages == nil { p.languages = make(map[string]ProgramEntitySpec) }
	if _, present := p.languages[id]; present {
		return errors.New("Language '" + id + "' is already listed for entity '" + p.id + "'")
	} else {
		p.languages[id] = language
	}

	return nil
}

func (p *Program) GetDependencies() map[string]ProgramEntitySpec {
	return p.dependencies
}

func (p* Program) AddDependency(id string, dependency ProgramEntitySpec) error {
	if p.dependencies == nil { p.dependencies = make(map[string]ProgramEntitySpec) }
	if _, present := p.dependencies[id]; present {
		return errors.New("Dependency '" + id + "' is already listed for entity '" + p.id + "'")
	} else {
		p.dependencies[id] = dependency
	}

	return nil
}