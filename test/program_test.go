package test

import (
	"libsm/objmodel"
	"libsm/yamlmodel"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProgramWithNullYaml(t *testing.T) {
	var p objmodel.Program
	var th TestHarness
	errs := p.Init(nil, th.Resolve)
	assert.Equal(t, len(errs), 1)
	assert.Equal(t, errs[0].Error(),"cannot convert nil yaml to program specification")
}

func TestProgramWithEmptySpec(t *testing.T) {
	var p1, p2 objmodel.Program
	var th TestHarness
	prog1 := yamlmodel.Entity {
		Id: "",
		Type: yamlmodel.Program,
		Description: "",
	}
	prog2 := yamlmodel.Entity {
		Id: "some-program",
		Type: yamlmodel.Program,
		Description: "",
	}
	errs := p1.Init(&prog1, th.Resolve)
	assert.Equal(t, "empty IDs are not allowed", errs[0].Error())
	
	errs = p2.Init(&prog2, th.Resolve)
	assert.Equal(t, "empty names are not allowed", errs[0].Error())
}

func TestProgramWithEmptyDescription(t *testing.T) {
	var p objmodel.Program
	var th TestHarness
	program := yamlmodel.Entity {
		Id: "some-program",
		Type: yamlmodel.Program,
		Name: "Some Program",
		Description: "",
	}
	errs := p.Init(&program, th.Resolve)
	assert.Equal(t, len(errs), 1)
	assert.Equal(t, errs[0].Error(),"warning... empty descriptions are useless")
}

func TestProgramWithNonProgramYaml(t *testing.T) {
	var p objmodel.Program
	var th TestHarness
	prog := yamlmodel.Entity {
		Id: "iron-man",
		Type: yamlmodel.Human,
		Description: "Human masquerading as an Android",
		ADM: []string{"careless.adm"},
	}
	errs := p.Init(&prog, th.Resolve)
	assert.Equal(t, len(errs), 1)
	assert.Equal(t, errs[0].Error(),"cannot initialize program with 'human'")
}

func TestProgramInheritedFromAnotherProgram(t *testing.T) {
	var p objmodel.Program
	var th TestHarness
	th.YamlStructures = make(map[string]interface{})
	th.YamlStructures["web-service"] = &yamlmodel.Entity {
		Id: "web-service",
		Name: "Web Service",
		Type: yamlmodel.Program,
		Description: "A Web Service",
		Base: []string{"service"},
	}
	th.YamlStructures["service"] = &yamlmodel.Entity {
		Id: "service",
		Name: "Generic Service",
		Type: yamlmodel.Program,
		Description: "A generic service program",
	}

	errs := p.Init(th.YamlStructures["web-service"].(*yamlmodel.Entity), th.Resolve)
	assert.Empty(t, errs)
	assert.Equal(t, "service", p.GetBase()[0].GetID())
}

func TestProgramInheritedFromHuman(t *testing.T) {
	var h objmodel.Program
	var th TestHarness
	th.YamlStructures = make(map[string]interface{})
	th.YamlStructures["sonny"] = &yamlmodel.Entity {
		Id: "sonny",
		Name: "Sonny",
		Type: yamlmodel.Program,
		Description: "The robot 'Sonny' from iRobot",
		Base: []string{"father"},
	}
	th.YamlStructures["father"] = &yamlmodel.Entity {
		Id: "father",
		Name: "Dr. Alfred Lanning",
		Type: yamlmodel.Human, // wrong base type
		Description: "Dr. Alfred Lanning, the scientist who create Sonny",
	}

	errs := h.Init(th.YamlStructures["sonny"].(*yamlmodel.Entity), th.Resolve)
	for _, err := range errs { // Confirm that all errors are related to the 'father' base reference
		assert.Contains(t, err.Error(), "father")
	}
}

func TestProgramMissingBase(t *testing.T) {
	var p objmodel.Program
	var th TestHarness
	th.YamlStructures = make(map[string]interface{})
	th.YamlStructures["program"] = &yamlmodel.Entity {
		Id: "program",
		Name: "Generic Program",
		Type: yamlmodel.Program,
		Description: "A generic program",
		Base: []string{"code"},
	}

	errs := p.Init(th.YamlStructures["program"].(*yamlmodel.Entity), th.Resolve)
	for _, err := range errs { // Confirm that all errors are related to the 'code' base reference
		assert.Contains(t, err.Error(), "code")
	}
}

func TestProgramWithInvalidRepoURL(t *testing.T) {
	var p objmodel.Program
	var th TestHarness
	th.YamlStructures = make(map[string]interface{})
	th.YamlStructures["adsm"] = &yamlmodel.Entity {
		Id: "adsm",
		Name: "ADSM",
		Type: yamlmodel.Program,
		Description: "ADM based Security Modeling",
		CodeRepository: "https://github%.com",
	}

	errs := p.Init(th.YamlStructures["adsm"].(*yamlmodel.Entity), th.Resolve)
	for _, err := range errs { // Confirm that all errors are related to code repository
		assert.Contains(t, err.Error(), "https://github%.com")
	}
}

func TestProgramWithMissingLanguageReference(t *testing.T) {
	var p objmodel.Program
	var th TestHarness
	th.YamlStructures = make(map[string]interface{})
	th.YamlStructures["server"] = &yamlmodel.Entity {
		Id: "server",
		Name: "Generic Server",
		Type: yamlmodel.Program,
		Description: "A generic server",
		Languages: []string{"go"},
	}

	errs := p.Init(th.YamlStructures["server"].(*yamlmodel.Entity), th.Resolve)
	for _, err := range errs { // Confirm that all errors are related to the 'go' language reference
		assert.Contains(t, err.Error(), "go")
	}
}

func TestProgramWithMissingDependencies(t *testing.T) {
	var p objmodel.Program
	var th TestHarness
	th.YamlStructures = make(map[string]interface{})
	th.YamlStructures["adsm"] = &yamlmodel.Entity {
		Id: "adsm",
		Name: "ADSM",
		Type: yamlmodel.Program,
		Description: "ADM based Security Modeling",
		Dependencies: []string{"libadm"},
	}

	errs := p.Init(th.YamlStructures["adsm"].(*yamlmodel.Entity), th.Resolve)
	for _, err := range errs { // Confirm that all errors are related to the 'libadm' dependency reference
		assert.Contains(t, err.Error(), "libadm")
	}
}

func TestProgramWithMissingRoles(t *testing.T) {
	var p objmodel.Program
	var th TestHarness
	th.YamlStructures = make(map[string]interface{})
	th.YamlStructures["adsm"] = &yamlmodel.Entity {
		Id: "adsm",
		Name: "ADSM",
		Type: yamlmodel.Program,
		Description: "ADM based Security Modeling",
		Roles: []string{"tool"},
	}

	errs := p.Init(th.YamlStructures["adsm"].(*yamlmodel.Entity), th.Resolve)
	for _, err := range errs { // Confirm that all errors are related to the 'tool' role reference
		assert.Contains(t, err.Error(), "tool")
	}
}

func TestProgramWithMultipleRoles(t *testing.T) {
	var p objmodel.Program
	var th TestHarness
	th.YamlStructures = make(map[string]interface{})
	th.YamlStructures["security-frontend"] = &yamlmodel.Entity {
		Id: "security-frontend",
		Type: yamlmodel.Program,
		Name: "Security Frontend",
		Description: "Frontend of a physical security monitoring system. It is the sole interface. It plays different roles depending on who's connected.",
		Roles: []string{"admin", "operator", "responder"},
	}
	th.YamlStructures["admin"] = &yamlmodel.Entity {
		Id: "admin",
		Type: yamlmodel.Role,
		Name: "Administrator",
		Description: "Represents the administrator of the security system",
		ADM: []string{"admin.adm"},
	}
	th.YamlStructures["operator"] = &yamlmodel.Entity {
		Id: "operator",
		Type: yamlmodel.Role,
		Name: "Operator",
		Description: "Represents the person manning the security monitoring station",
		ADM: []string{"operator.adm"},
	}
	th.YamlStructures["responder"] = &yamlmodel.Entity {
		Id: "responder",
		Type: yamlmodel.Role,
		Name: "Emergency Responder",
		Description: "Represents the person responding to security incidents",
		ADM: []string{"responder.adm"},
	}

	errs := p.Init(th.YamlStructures["security-frontend"].(*yamlmodel.Entity), th.Resolve)
	assert.Empty(t, errs)
	for _, r := range p.GetRoles() {
		assert.Contains(t, []string{"admin", "operator", "responder"}, r.GetID())
	}
}

func TestProgramWithRepeatReferences(t *testing.T) {
	var p objmodel.Program
	var th TestHarness
	th.YamlStructures = make(map[string]interface{})
	th.YamlStructures["server"] = &yamlmodel.Entity {
		Id: "server",
		Type: yamlmodel.Program,
		Name: "Server",
		Description: "A generic server",
		Roles: []string{"dummy", "dummy"},
		Languages: []string{"dummy", "dummy"},
		Dependencies: []string{"dummy", "dummy"},
	}
	th.YamlStructures["dummy"] = &yamlmodel.Entity { // to avoid 'missing reference' error
		Id: "dummy",
		Name: "Dummy program",
		Type: yamlmodel.Program,
		Description: "A dummy program specification",
	}

	errs := p.Init(th.YamlStructures["server"].(*yamlmodel.Entity), th.Resolve)
	for _, err := range errs { // Confirm that all errors are related to the 'go' language reference
		assert.Contains(t, err.Error(), "dummy")
	}
}

func TestFullyDefinedProgram(t *testing.T) {
	var p objmodel.Program
	var th TestHarness
	th.YamlStructures = make(map[string]interface{})
	th.YamlStructures["cloud-service"] = &yamlmodel.Entity {
		Id: "cloud-service",
		Type: yamlmodel.Program,
		Name: "Web Service in the Cloud",
		Description: "A web-service hosted in the cloud",
		Base: []string{"web-service"},
		AdmDir: "~/cloud-service",
		CodeRepository: "https://github.com/dummy/cloud-service",
		Recommendations: []string{"Cloud service must always be fronted by a firewall & IDS"},
		Roles: []string{"frontend-role"}, // Roles are for downstream services (not modelled here)
		Languages: []string{"go", "json"},
		Dependencies: []string{"go.net", "go.encoding.json"},
		ADM: []string{"cloud-service.adm"},
	}
	th.YamlStructures["web-service"] = &yamlmodel.Entity {
		Id: "web-service",
		Type: yamlmodel.Program,
		Name: "Web Service",
		Description: "A generic web-service",
		Recommendations: []string{"Never directly use URL path to resolve local paths on the server"},
		ADM: []string{"web-service.adm"},
	}
	th.YamlStructures["frontend-role"] = &yamlmodel.Entity {
		Id: "frontend-role",
		Type: yamlmodel.Role,
		Name: "Frontend",
		Description: "A frontend's role",
		Recommendations: []string{"A frontend's role must always follow least-privilege since it is susceptible to attacks from internet."},
		ADM: []string{"go.adm"},
	}
	th.YamlStructures["go"] = &yamlmodel.Entity {
		Id: "go",
		Name: "Go",
		Type: yamlmodel.Program,
		Description: "Generic security model for a program written in Go",
		ADM: []string{"go.adm"},
	}
	th.YamlStructures["json"] = &yamlmodel.Entity {
		Id: "json",
		Name: "JSON",
		Type: yamlmodel.Program,
		Description: "Generic security model for a program that uses JSON for data representation",
		ADM: []string{"json.adm"},
	}
	th.YamlStructures["go.net"] = &yamlmodel.Entity {
		Id: "go.net",
		Name: "'net' package (Go)",
		Type: yamlmodel.Program,
		Description: "Security model for Go's 'net' package",
		ADM: []string{"go.net.adm"},
	}
	th.YamlStructures["go.encoding.json"] = &yamlmodel.Entity {
		Id: "go.encoding.json",
		Type: yamlmodel.Program,
		Name: "encoding/json package (Go)",
		Description: "Security model for Go's 'encoding/json' package",
		ADM: []string{"go.encoding.json.adm"},
	}
	
	errs := p.Init(th.YamlStructures["cloud-service"].(*yamlmodel.Entity), th.Resolve)

	// Add additional recommendation to cloud-service
	p.AddRecommendation("Capture logs from cloud-service and store them in SIEM")

	assert.Empty(t, errs)
	assert.Equal(t, "cloud-service", p.GetID())
	assert.Equal(t, "A web-service hosted in the cloud", p.GetDescription())
	assert.Equal(t, "https://github.com/dummy/cloud-service", p.GetCodeRepository())
	assert.Contains(t, p.GetADM()[p.GetID()], "~/cloud-service/cloud-service.adm")
	for _, recos := range p.GetRecommendations() {
		for _, reco := range recos {
			assert.Contains(t, []string{"Cloud service must always be fronted by a firewall & IDS",
			"Never directly use URL path to resolve local paths on the server",
			"A frontend's role must always follow least-privilege since it is susceptible to attacks from internet.",
			"Capture logs from cloud-service and store them in SIEM",
			}, reco)
		}
	}
	for role := range p.GetRoles() {
		assert.Contains(t, []string{"frontend-role"}, role)
	}
	for lang := range p.GetLanguages() {
		assert.Contains(t, []string{"go","json"}, lang)
	}
	for dep := range p.GetDependencies() {
		assert.Contains(t, []string{"go.net","go.encoding.json"}, dep)
	}
}

func TestProgramWithAdmInAnotherDirectory(t *testing.T) {
	var p objmodel.Program
	var th TestHarness
	prog := yamlmodel.Entity {
		Id: "some-prog",
		Name: "Some program",
		Type: yamlmodel.Program,
		Description: "Some Program",
		AdmDir: "/root",
		ADM: []string{"oh-no.adm"},
	}
	errs := p.Init(&prog, th.Resolve)
	assert.Empty(t, errs)
	assert.Equal(t, []string{"/root/oh-no.adm"}, p.GetADM()[p.GetID()])
}