package loaders

import (
	"errors"

	"libaddb"
	"libsm/objmodel"
	"libsm/yamlmodel"

	"gopkg.in/yaml.v3"
)

type Loader struct {
	builder Builder
}

func (l *Loader) LoadSecurityModel(yamlText string, admDir string) (*objmodel.SecurityModel, []error) {
	var errs []error
	
	if yamlText == "" {
		return nil, []error{errors.New("cannot work with empty YAML content")}
	}

	var m yamlmodel.SecurityModel
	err := yaml.Unmarshal([]byte(yamlText), &m)
	if err != nil {
		return nil, []error{err}
	}

	var addb libaddb.ADDB
	err = addb.Init(m.AddbUri)
	if err != nil {
		errs = append(errs, err)
	}

	idxErrs := l.builder.Index("", &m, admDir, &addb)
	if len(idxErrs) != 0 {
		errs = append(errs, idxErrs...)
	}

	m.AdmDir = admDir
	
	var model objmodel.SecurityModel
	initErrs := model.Init(&m, l.builder.Resolve)
	if len(initErrs) != 0 {
		errs = append(errs, initErrs...)
	}

	return &model, errs
}

func (l *Loader) LoadHuman(yamlText string, admDir string, addbUrl string) (*objmodel.Human, []error) {
	var errs []error
	
	if yamlText == "" {
		return nil, []error{errors.New("cannot work with empty YAML content")}
	}

	var h yamlmodel.Entity
	err := yaml.Unmarshal([]byte(yamlText), &h)
	if err != nil {
		return nil, []error{err}
	}

	if h.Type != yamlmodel.Human {
		return nil, []error{errors.New("Expected human, got " + string(h.Type))}
	}

	var addb libaddb.ADDB
	err = addb.Init(addbUrl)
	if err != nil {
		errs = append(errs, err)
	}
	
	// Index entity and its parts
	idxErrs := l.builder.Index(h.Id, &h, admDir, &addb)
	if len(idxErrs) != 0 { 
		errs = append(errs, idxErrs...)
	}

	var human objmodel.Human
	initErrs := human.Init(&h, l.builder.Resolve)
	if len(initErrs) != 0 {
		errs = append(errs, initErrs...)
	}

	return &human, errs
}

func (l *Loader) LoadProgram(yamlText string, admDir string, addbUrl string) (*objmodel.Program, []error) {
	var errs []error

	if yamlText == "" {
		return nil, []error{errors.New("cannot work with empty YAML content")}
	}

	var p yamlmodel.Entity
	err := yaml.Unmarshal([]byte(yamlText), &p)
	if err != nil {
		return nil, []error{err}
	}

	if p.Type != yamlmodel.Program && p.Type != yamlmodel.System {
		return nil, []error{errors.New("Expected program, got " + string(p.Type))}
	}

	var addb libaddb.ADDB
	err = addb.Init(addbUrl)
	if err != nil {
		errs = append(errs, err)
	}

	// Index entity and its parts
	idxErrs := l.builder.Index(p.Id, &p,  admDir, &addb)
	if idxErrs != nil { 
		errs = append(errs, idxErrs...)
	}

	var program objmodel.Program
	initErrs := program.Init(&p, l.builder.Resolve)
	if len(initErrs) != 0 {
		errs = append(errs, initErrs...)
	}

	return &program, errs
}

func (l *Loader) LoadFlow(yamlText string, admDir string, addbUrl string) (*objmodel.Flow, []error) {
	var errs []error

	if yamlText == "" {
		return nil, []error{errors.New("cannot work with empty YAML content")}
	}

	var f yamlmodel.Flow
	err := yaml.Unmarshal([]byte(yamlText), &f)
	if err != nil {
		return nil, []error{err}
	}

	var addb libaddb.ADDB
	err = addb.Init(addbUrl)
	if err != nil {
		errs = append(errs, err)
	}

	// Index entity and its parts
	idxErrs := l.builder.Index(f.Id, &f,  admDir, &addb)
	if len(idxErrs) != 0 { 
		errs = append(errs, idxErrs...)
	}

	var flow objmodel.Flow
	initErrs := flow.Init(&f, l.builder.Resolve)
	if initErrs != nil {
		errs = append(errs, initErrs...)
	}

	return &flow, errs
}