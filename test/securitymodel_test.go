package test

import (
	"libsm/objmodel"
	"libsm/yamlmodel"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSMWithNullYaml(t *testing.T) {
	var th TestHarness
	var sm objmodel.SecurityModel
	
	err := sm.Init(nil, th.Resolve)
	assert.NotNil(t, err)
}

func TestSMWithEmptyYamlStructure(t *testing.T) {
	var th TestHarness
	var ysm yamlmodel.SecurityModel
	var sm objmodel.SecurityModel
	
	err := sm.Init(&ysm, th.Resolve)
	assert.Nil(t, err)
	assert.Empty(t, sm.Title)
	assert.Empty(t, sm.DesignDocument)
	assert.Empty(t, sm.AddbPath)
	assert.Empty(t, sm.Externals)
	assert.Empty(t, sm.Entities)
	assert.Empty(t, sm.Flows)
}

func TestSMWithFullModel(t *testing.T) {
	var th TestHarness
	th.YamlStructures = make(map[string]interface{})
	th.YamlStructures["user"] = &yamlmodel.Entity {
		Id: "user",
		Type: yamlmodel.Human,
		Description: "User interacting with the system",
		Interface: "user-browser",
		ADM: []string{"user.adm"},
	}
	th.YamlStructures["user-browser"] = &yamlmodel.Entity {
		Id: "user-browser",
		Type: yamlmodel.Program,
		Description: "Interface used by 'user' to talk to system",
		ADM: []string{"user-browser.adm"},
	}
	th.YamlStructures["system"] = &yamlmodel.Entity {
		Id: "system",
		Type: yamlmodel.System,
		Description: "System under consideration",
		Base: []string{"server"},
		ADM: []string{"system.adm"},
	}
	th.YamlStructures["interaction"] = &yamlmodel.Flow {
		Id: "interaction",
		Sender: "user",
		Receiver: "system",
		Description: "User interacting with the system",
		Protocol: []string{"manual"},
		ADM: []string{"interaction.adm"},
	}
	var sm objmodel.SecurityModel
	ysm := yamlmodel.SecurityModel {
		Title: "Simple security model",
		Externals: []*yamlmodel.Entity {
			th.YamlStructures["user"].(*yamlmodel.Entity),
			th.YamlStructures["user-browser"].(*yamlmodel.Entity),
		},
		Entities: []*yamlmodel.Entity {th.YamlStructures["system"].(*yamlmodel.Entity)},
		Flows: []*yamlmodel.Flow {th.YamlStructures["interaction"].(*yamlmodel.Flow)},
	}
	
	errs := sm.Init(&ysm, th.Resolve)
	assert.GreaterOrEqual(t, len(errs), 1) // 'server' and 'manual' cannot be resolved.
}

func TestSMWithSocialEnggModel(t *testing.T) {
	var th TestHarness
	th.YamlStructures = make(map[string]interface{})
	th.YamlStructures["adam"] = &yamlmodel.Entity {
		Id: "adam",
		Type: yamlmodel.Human,
		Description: "Adam is the target of social-engineering attacks",
		ADM: []string{"adam.adm"},
	}
	th.YamlStructures["bob"] = &yamlmodel.Entity {
		Id: "bob",
		Type: yamlmodel.Human,
		Description: "Bob is the attacker",
		Base: []string{"crook"},
		ADM: []string{"bob.adm"},
	}
	th.YamlStructures["soceng-attack"] = &yamlmodel.Flow {
		Id: "soceng-attack",
		Sender: "Bob",
		Receiver: "Adam",
		Description: "Bob social-engineers Adam to reveal his bank credentials",
		Protocol: []string{"manual"},
		ADM: []string{"soceng-creds.adm"},
	}
	var sm objmodel.SecurityModel
	ysm := yamlmodel.SecurityModel {
		Title: "Simple security model",
		Externals: nil,
		Entities: []*yamlmodel.Entity {
			th.YamlStructures["adam"].(*yamlmodel.Entity),
			th.YamlStructures["bob"].(*yamlmodel.Entity),
		},
		Flows: []*yamlmodel.Flow {th.YamlStructures["soceng-attack"].(*yamlmodel.Flow)},
	}
	
	errs := sm.Init(&ysm, th.Resolve)
	assert.GreaterOrEqual(t, len(errs), 1) // 'crook' and 'manual' cannot be resolved.
}

func TestManualADMAddition(t *testing.T) {
	var th TestHarness
	var ysm yamlmodel.SecurityModel
	var sm objmodel.SecurityModel
	
	err := sm.Init(&ysm, th.Resolve)
	assert.Nil(t, err)
	sm.SetADM([]string{"waste.adm", "pointless.adm"})
	assert.Contains(t, sm.GetADM()["sm"], "waste.adm")
	assert.Contains(t, sm.GetADM()["sm"], "pointless.adm")
}