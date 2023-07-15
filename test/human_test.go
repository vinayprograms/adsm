package test

import (
	"libsm/objmodel"
	"libsm/yamlmodel"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHumanWithNullYaml(t *testing.T) {
	var h objmodel.Human
	var th TestHarness
	errs := h.Init(nil, th.Resolve)
	assert.Equal(t, len(errs), 1)
	assert.Equal(t, errs[0].Error(),"cannot convert nil yaml to human specification")
}

func TestHumanWithEmptySpec(t *testing.T) {
	var h1, h2 objmodel.Human
	var th TestHarness
	person1 := yamlmodel.Entity {
		Id: "",
		Type: yamlmodel.Human,
		Description: "",
	}
	person2 := yamlmodel.Entity {
		Id: "someone",
		Type: yamlmodel.Human,
		Description: "",
	}
	errs := h1.Init(&person1, th.Resolve)
	assert.Equal(t, "empty IDs are not allowed", errs[0].Error())
	
	errs = h2.Init(&person2, th.Resolve)
	assert.Equal(t, "empty names are not allowed", errs[0].Error())
}


func TestHumanWithEmptyDescription(t *testing.T) {
	var h objmodel.Human
	var th TestHarness
	human := yamlmodel.Entity {
		Id: "someone",
		Type: yamlmodel.Human,
		Name: "Someone",
		Description: "",
	}
	errs := h.Init(&human, th.Resolve)
	assert.Equal(t, len(errs), 1)
	assert.Equal(t, errs[0].Error(),"warning... empty descriptions are useless")
}

func TestHumanWithNonHumanYaml(t *testing.T) {
	var h objmodel.Human
	var th TestHarness
	prog := yamlmodel.Entity {
		Id: "android",
		Name: "Android",
		Type: yamlmodel.System,
		Description: "AI system masquerading as a human",
		ADM: []string{"skynet.adm", "matrix.adm", "ultron.adm", "ARIIA.adm"},
	}
	errs := h.Init(&prog, th.Resolve)
	assert.Equal(t, len(errs), 1)
	assert.Equal(t, errs[0].Error(),"cannot initialize human with 'system'")
}

func TestHumanInheritedFromAnotherHuman(t *testing.T) {
	var h objmodel.Human
	var th TestHarness
	th.YamlStructures = make(map[string]interface{})
	th.YamlStructures["teen"] = &yamlmodel.Entity {
		Id: "teen",
		Name: "Teenager",
		Type: yamlmodel.Human,
		Description: "A teenager",
		Base: []string{"father"},
	}
	th.YamlStructures["father"] = &yamlmodel.Entity {
		Id: "father",
		Name: "Teen's Father",
		Type: yamlmodel.Human,
		Description: "Teen's father",
	}

	errs := h.Init(th.YamlStructures["teen"].(*yamlmodel.Entity), th.Resolve)
	assert.Empty(t, errs)
	assert.Equal(t, "father", h.GetBase()[0].GetID())
	assert.Equal(t, "Teen's Father", h.GetBase()[0].GetName())
}

func TestHumanInheritanceProblem(t *testing.T) {
	var h objmodel.Human
	var th TestHarness
	th.YamlStructures = make(map[string]interface{})
	th.YamlStructures["adult"] = &yamlmodel.Entity {
		Id: "adult",
		Name: "An adult",
		Type: yamlmodel.Human,
		Description: "Someone who just turned 18.",
		Base: []string{"father"},
	}

	errs := h.Init(th.YamlStructures["adult"].(*yamlmodel.Entity), th.Resolve)
	for _, err := range errs { // Confirm that all errors are related to the 'father' base reference
		assert.Contains(t, err.Error(), "father")
	}
}

func TestHumanInheritedFromProgram(t *testing.T) {
	var h objmodel.Human
	var th TestHarness
	th.YamlStructures = make(map[string]interface{})
	th.YamlStructures["teen"] = &yamlmodel.Entity {
		Id: "teen",
		Name: "Teenager",
		Type: yamlmodel.Human,
		Description: "A teenager",
		Base: []string{"father"},
	}
	th.YamlStructures["father"] = &yamlmodel.Entity {
		Id: "father",
		Name: "Teen's father",
		Type: yamlmodel.Program, // wrong base type
		Description: "Teen's father who is apparently a robot!",
	}

	errs := h.Init(th.YamlStructures["teen"].(*yamlmodel.Entity), th.Resolve)
	for _, err := range errs { // Confirm that all errors are related to the 'father' base reference
		assert.Contains(t, err.Error(), "father")
	}
}

func TestHumanWithInterface(t *testing.T) {
	var h objmodel.Human
	var th TestHarness
	th.YamlStructures = make(map[string]interface{})
	th.YamlStructures["agent"] = &yamlmodel.Entity {
		Id: "agent",
		Name: "Customer Care Agent",
		Type: yamlmodel.Human,
		Description: "A customer care agent",
		Interface: "agent-browser",
	}
	th.YamlStructures["agent-browser"] = &yamlmodel.Entity {
		Id: "agent-browser",
		Name: "Web Browser used by Customer Care Agent",
		Type: yamlmodel.Program,
		Description: "Browser used by customer-care agent",
		Roles: []string{"agent-role"},
	}

	th.YamlStructures["agent-role"] = &yamlmodel.Entity {
		Id: "agent-role",
		Name: "Agen't role in the system",
		Type: yamlmodel.Role,
		Description: "Agent's role when browsing the system",
	}

	errs := h.Init(th.YamlStructures["agent"].(*yamlmodel.Entity), th.Resolve)
	assert.Nil(t, errs)
	assert.Equal(t, "agent-browser", h.GetUserInterface().GetID())
}

func TestHumanWithMissingInterface(t *testing.T) {
	var h objmodel.Human
	var th TestHarness
	th.YamlStructures = make(map[string]interface{})
	th.YamlStructures["agent"] = &yamlmodel.Entity {
		Id: "agent",
		Name: "Customer Care Agent",
		Type: yamlmodel.Human,
		Description: "A customer care agent",
		Interface: "agent-browser",
	}

	errs := h.Init(th.YamlStructures["agent"].(*yamlmodel.Entity), th.Resolve)
	for _, err := range errs { // Confirm that all errors are related to the 'agent-browser' base reference
		assert.Contains(t, err.Error(), "agent-browser")
	}
}

func TestComplexHuman(t *testing.T) {
	var h objmodel.Human
	var th TestHarness
	th.YamlStructures = make(map[string]interface{})
	th.YamlStructures["agent"] = &yamlmodel.Entity {
		Id: "agent",
		Type: yamlmodel.Human,
		Name: "Customer Care Agent",
		Description: "A customer care agent",
		Base: []string{"employee"},
		AdmDir: "/home/agent",
		Recommendations: []string{"Always start conversation with a warm greeting."},
		Interface: "agent-browser",
		ADM: []string{"agent.adm"},
	}
	th.YamlStructures["employee"] = &yamlmodel.Entity {
		Id: "employee",
		Type: yamlmodel.Human,
		Name: "Company Employee",
		Description: "Employee of a company",
		AdmDir: "/home/agent/company",
		Recommendations: []string{"In god we trust. Rest have to authenticate."},
		ADM: []string{"employee.adm"},
	}
	th.YamlStructures["agent-browser"] = &yamlmodel.Entity {
		Id: "agent-browser",
		Type: yamlmodel.Program,
		Name: "Generic Web Browser",
		Description: "Agent interacts with the system via web-browser.",
		AdmDir: "/home/agent",
		Recommendations: []string{"Always run the latest version."},
		ADM: []string{"agent-browser.adm"},
	}
	errs := h.Init(th.YamlStructures["agent"].(*yamlmodel.Entity), th.Resolve)

	// Add browser recommendation to agent. Interface is a separate entity and 
	// its recommendations are not inherited by agent
	h.AddRecommendation("Always run the latest version of your web-browser.")

	assert.Empty(t, errs)
	assert.Equal(t, "agent", h.GetID())
	assert.Equal(t, "Customer Care Agent", h.GetName())
	assert.Equal(t, "A customer care agent", h.GetDescription())
	assert.Contains(t, h.GetADM()[h.GetID()], "/home/agent/agent.adm")
	assert.Contains(t, h.GetADM()[h.GetID()+".base."+h.GetBase()[0].GetID()], "/home/agent/company/employee.adm")
	for _, recos := range h.GetRecommendations() {
		for _, reco := range recos {
			assert.Contains(t, []string{"Always start conversation with a warm greeting.",
			"In god we trust. Rest have to authenticate.",
			"Always run the latest version of your web-browser.",
			}, reco)
		}
	}
}

func TestHumanWithAdmInAnotherDirectory(t *testing.T) {
	var h objmodel.Human
	var th TestHarness
	prog := yamlmodel.Entity {
		Id: "someone",
		Name: "Anonymous",
		Type: yamlmodel.Human,
		Description: "Someone",
		AdmDir: "/root",
		ADM: []string{"oh-no.adm"},
	}
	errs := h.Init(&prog, th.Resolve)
	assert.Empty(t, errs)
	assert.Equal(t, []string{"/root/oh-no.adm"}, h.GetADM()[h.GetID()])
}