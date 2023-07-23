package objmodel

import (
	"errors"
)

////////////////////////////////////////
// Interface for builders. Object Model entities use these two functions to
// resolve references to other entities.
type Resolver func(string) (interface{}, []error)
type Indexer func(string, interface{}) error


////////////////////////////////////////
// Base Interfaces. These are not directly instantiated in the object model.
type CoreSpec interface {
	GetID() string
	SetID(string) error
	GetName() string
	SetName(string) error
	GetDescription() string
	SetDescription(string) error
}

type EntitySpec interface {
	CoreSpec
	GetADM() map[string][]string
	SetADM([]string) error
	AddADM(string) error
	GetMitigations() map[string][]string
	SetMitigations([]string) error
	AddMitigation(reco string) error
	GetRecommendations() map[string][]string
	SetRecommendations([]string) error
	AddRecommendation(reco string) error
}

type ExternalSpec interface {
	CoreSpec
}

type HumanSpec interface {
	CoreSpec
	GetUserInterface() ProgramSpec
	SetUserInterface(ProgramSpec) error
}

type ProgramSpec interface {
	CoreSpec
	GetRoles() map[string]EntitySpec				// Entity spec. because roles must be part of 'entities' section of the Security model.
	AddRole(string, EntitySpec) error
}

////////////////////////////////////////
// Security Model Interfaces. These are used directly when building object model
// and for processing all security model entities.
type ExternalHumanSpec interface {
	ExternalSpec		// Basic identity information
	HumanSpec				// Basic information about humans external to the system.
}
type ExternalProgramSpec interface {
	ExternalSpec		// Basic identity information
	ProgramSpec
}

type HumanEntitySpec interface {
	EntitySpec
	HumanSpec
	GetBase() []HumanEntitySpec
	SetBase([]HumanEntitySpec) error
	AddBase(HumanEntitySpec) error
}
type ProgramEntitySpec interface {
	EntitySpec
	ProgramSpec
	GetCodeRepository() string
	SetRepository(string) error
	GetBase() []ProgramEntitySpec
	SetBase([]ProgramEntitySpec) error
	AddBase(ProgramEntitySpec) error
	GetLanguages() map[string]ProgramEntitySpec
	AddLanguage(id string, language ProgramEntitySpec) error
	GetDependencies() map[string]ProgramEntitySpec
	AddDependency(id string, dependency ProgramEntitySpec) error
}

type FlowSpec interface {
	EntitySpec
	GetProtocol() map[string]FlowSpec
	AddProtocol(string, FlowSpec) error		
	GetSender() CoreSpec
	SetSender(CoreSpec) error
	GetReceiver() CoreSpec
	SetReceiver(CoreSpec) error
}

////////////////////////////////////////
// Common data structures used across models

type CoreObject struct {
	id string
	name string
	description string
	adm []string
	mitigations []string
	recommendations []string
}

func (c *CoreObject) GetID() string {
	return c.id
}

func (c *CoreObject) SetID(id string) error {
	if id == "" {
		return errors.New("empty IDs are not allowed")
	}
	c.id = id
	return nil
}

func (c *CoreObject) GetName() string {
	return c.name
}

func (c *CoreObject) SetName(name string) error {
	if name == "" {
		return errors.New("empty names are not allowed")
	}
	c.name = name
	return nil
}

func (c *CoreObject) GetDescription() string {
	return c.description
}

func (c *CoreObject) SetDescription(desc string) error {
	if desc == "" {
		return errors.New("warning... empty descriptions are useless")
	}
	c.description = desc
	return nil
}

func (c *CoreObject) SetADM(admlist []string) error {
	c.adm = admlist
	return nil
}

// Implemented so that it can adhere to EntitySpec.
func (c *CoreObject) SetMitigations(mitigations []string) error {
	c.mitigations = mitigations
	return nil
}


// Implemented so that it can adhere to EntitySpec.
func (c *CoreObject) SetRecommendations(reco []string) error {
	c.recommendations = reco
	return nil
}

////////////////////////////////////////
// Helper functions

// Used by all GetADM() implementations in this package.
func merge(existing map[string][]string, baseNamespace string, new map[string][]string) (updated map[string][]string) {
	updated = existing
	for k, v := range new {
		key := baseNamespace + "." + k
		if _, present := updated[key]; present {
			updated[key] = append(updated[key], v...)
		} else {
			updated[key] = v
		}
	}

	return
}
