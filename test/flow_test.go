package test

import (
	"libsm/objmodel"
	"libsm/yamlmodel"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlowWithNullYaml(t *testing.T) {
	var f objmodel.Flow
	var th TestHarness
	errs := f.Init(nil, th.Resolve)
	assert.Equal(t, len(errs), 1)
	assert.Equal(t, errs[0].Error(),"cannot convert nil yaml to flow specification")
}

func TestFlowWithEmptySpec(t *testing.T) {
	var f1, f2 objmodel.Flow
	var th TestHarness
	flow1 := yamlmodel.Flow {
		Id: "",
		Description: "",
	}
	flow2 := yamlmodel.Flow {
		Id: "some-flow",
		Description: "",
	}
	errs := f1.Init(&flow1, th.Resolve)
	assert.Equal(t, "empty IDs are not allowed", errs[0].Error())
	
	errs = f2.Init(&flow2, th.Resolve)
	assert.Equal(t, "empty names are not allowed", errs[0].Error())
}

func TestFlowWithEmptyDescription(t *testing.T) {
	var f objmodel.Flow
	var th TestHarness
	flow := yamlmodel.Flow {
		Id: "some-flow",
		Name: "Some Flow",
		Description: "",
	}
	errs := f.Init(&flow, th.Resolve)
	assert.Equal(t, len(errs), 1)
	assert.Equal(t, errs[0].Error(),"warning... empty descriptions are useless")
}

func TestFlowWithMissingSenderReference(t *testing.T) {
	var f objmodel.Flow
	var th TestHarness
	th.YamlStructures = make(map[string]interface{})
	th.YamlStructures["https"] = &yamlmodel.Flow {
		Id: "https",
		Name: "HTTP(S)",
		Description: "HTTP on top of TLS",
		Sender: "browser",
		Receiver: "web-server",
	}
	th.YamlStructures["web-server"] = &yamlmodel.Entity {
		Id: "web-server",
		Name: "Generic Web Server",
		Type: yamlmodel.Program,
		Description: "A generic web server",
	}

	errs := f.Init(th.YamlStructures["https"].(*yamlmodel.Flow), th.Resolve)
	for _, err := range errs { // Confirm that all errors are related to the 'browser' sender reference
		assert.Contains(t, err.Error(), "browser")
	}
}

func TestFlowWithMissingReceiverReference(t *testing.T) {
	var f objmodel.Flow
	var th TestHarness
	th.YamlStructures = make(map[string]interface{})
	th.YamlStructures["https"] = &yamlmodel.Flow {
		Id: "https",
		Name: "HTTP(S)",
		Description: "HTTP on top of TLS",
		Sender: "browser",
		Receiver: "web-server",
	}
	th.YamlStructures["browser"] = &yamlmodel.Entity {
		Id: "browser",
		Name: "Generic Web Browser",
		Type: yamlmodel.Program,
		Description: "A generic web browser",
	}

	errs := f.Init(th.YamlStructures["https"].(*yamlmodel.Flow), th.Resolve)
	for _, err := range errs { // Confirm that all errors are related to the 'web-server' receiver reference
		assert.Contains(t, err.Error(), "web-server")
	}
}

func TestFlowWithMissingProtocolReference(t *testing.T) {
	var f objmodel.Flow
	var th TestHarness
	th.YamlStructures = make(map[string]interface{})
	th.YamlStructures["https"] = &yamlmodel.Flow {
		Id: "https",
		Name: "HTTPS(S)",
		Description: "HTTP on top of TLS",
		Sender: "browser",
		Receiver: "web-server",
		Protocol: []string{"tls"},
	}
	th.YamlStructures["browser"] = &yamlmodel.Entity {
		Id: "browser",
		Name: "Generic Web Browser",
		Type: yamlmodel.Program,
		Description: "A generic web browser",
	}
	th.YamlStructures["web-server"] = &yamlmodel.Entity {
		Id: "web-server",
		Name: "Generic Web Server",
		Type: yamlmodel.Program,
		Description: "A generic web server",
	}

	errs := f.Init(th.YamlStructures["https"].(*yamlmodel.Flow), th.Resolve)
	for _, err := range errs { // Confirm that all errors are related to the 'tls' protocol reference
		assert.Contains(t, err.Error(), "tls")
	}
}

func TestFlowWithValidProtocolReference(t *testing.T) {
	var f objmodel.Flow
	var th TestHarness
	th.YamlStructures = make(map[string]interface{})
	th.YamlStructures["https"] = &yamlmodel.Flow {
		Id: "https",
		Name: "HTTP(S)",
		Description: "HTTP on top of TLS",
		Sender: "browser",
		Receiver: "web-server",
		Protocol: []string{"tls@1.2"},
	}
	th.YamlStructures["browser"] = &yamlmodel.Entity {
		Id: "browser",
		Name: "Generic Web Browser",
		Type: yamlmodel.Program,
		Description: "A generic web browser",
	}
	th.YamlStructures["web-server"] = &yamlmodel.Entity {
		Id: "web-server",
		Name: "Generic Web Server",
		Type: yamlmodel.Program,
		Description: "A generic web server",
	}
	th.YamlStructures["tls@1.2"] = &yamlmodel.Flow {
		Id: "tls@1.2",
		Name: "TLS v1.2",
		Description: "TLS v1.2",
	}

	errs := f.Init(th.YamlStructures["https"].(*yamlmodel.Flow), th.Resolve)
	assert.Empty(t, errs)
}

func TestProgramWithRepeatProtocolReferences(t *testing.T) {
	var f objmodel.Flow
	var th TestHarness
	th.YamlStructures = make(map[string]interface{})
	th.YamlStructures["https"] = &yamlmodel.Flow {
		Id: "https",
		Name: "HTTP(S)",
		Description: "HTTP on top of TLS",
		Sender: "browser",
		Receiver: "web-server",
		Protocol: []string{"tls@1.2", "tls@1.2"},
	}
	th.YamlStructures["browser"] = &yamlmodel.Entity {
		Id: "browser",
		Name: "Generic Web Browser",
		Type: yamlmodel.Program,
		Description: "A generic web browser",
	}
	th.YamlStructures["web-server"] = &yamlmodel.Entity {
		Id: "web-server",
		Name: "Generic Web Server",
		Type: yamlmodel.Program,
		Description: "A generic web server",
	}
	th.YamlStructures["tls@1.2"] = &yamlmodel.Flow {
		Id: "tls@1.2",
		Name: "TLS v1.2",
		Description: "TLS v1.2",
	}

	errs := f.Init(th.YamlStructures["https"].(*yamlmodel.Flow), th.Resolve)
	for _, err := range errs { // Confirm that all errors are related to the 'tls' protocol reference
		assert.Contains(t, err.Error(), "tls@1.2")
	}
}

func TestFullyDefinedFlow(t *testing.T) {
	var f objmodel.Flow
	var th TestHarness
	th.YamlStructures = make(map[string]interface{})
	th.YamlStructures["https"] = &yamlmodel.Flow {
		Id: "https",
		Name: "HTTP(S)",
		Description: "HTTP on top of TLS",
		Sender: "browser",
		Receiver: "web-server",
		Protocol: []string{"http@2", "tls@1.2"},
		ADM: []string{"https.adm"},
	}
	th.YamlStructures["browser"] = &yamlmodel.Entity {
		Id: "browser",
		Name: "Generic Web Browser",
		Type: yamlmodel.Program,
		Description: "A generic web browser",
	}
	th.YamlStructures["web-server"] = &yamlmodel.Entity {
		Id: "web-server",
		Name: "Generic Web Server",
		Type: yamlmodel.Program,
		Description: "A generic web server",
		Recommendations: []string{"Always enforce HTTPS (i.e., HSTS)"},
	}
	th.YamlStructures["tls@1.2"] = &yamlmodel.Flow {
		Id: "tls@1.2",
		Name: "TLS v1.2",
		Description: "TLS v1.2",
		Recommendations: []string{"Don't support weak ciphers in TLS 1.2"},
		ADM: []string{"tls@1.2.adm"},
	}
	th.YamlStructures["http@2"] = &yamlmodel.Flow {
		Id: "http@2",
		Name: "HTTP/2",
		Description: "HTTP/2 protocol",
		ADM: []string{"http@2.adm"},
	}
	
	errs := f.Init(th.YamlStructures["https"].(*yamlmodel.Flow), th.Resolve)
	assert.Empty(t, errs)

	// Add additional recommendation to https
	f.AddRecommendation("Don't accept certificates that don't belong to your target domain")

	
	assert.Equal(t, "https", f.GetID())
	assert.Equal(t, "HTTP(S)", f.GetName())
	assert.Equal(t, "HTTP on top of TLS", f.GetDescription())
	assert.Equal(t, "browser", f.GetSender().GetID())
	assert.Equal(t, "Generic Web Browser", f.GetSender().GetName())
	assert.Equal(t, "web-server", f.GetReceiver().GetID())
	assert.Equal(t, "Generic Web Server", f.GetReceiver().GetName())
	protocols := f.GetProtocol()
	assert.Equal(t, "HTTP/2 protocol", protocols["http@2"].GetDescription())
	assert.Equal(t, "TLS v1.2", protocols["tls@1.2"].GetDescription())
	adm := f.GetADM()[f.GetID()]
	assert.Contains(t, adm, "https.adm")
	for _, recos := range f.GetRecommendations() {
		for _, reco := range recos {
			assert.Contains(t, []string{
				"Don't accept certificates that don't belong to your target domain",
				"Always enforce HTTPS (i.e., HSTS)",
				"Don't support weak ciphers in TLS 1.2",
			}, reco)
		}
	}
}

func TestFlowWithAdmInAnotherDirectory(t *testing.T) {
	var f objmodel.Flow
	var th TestHarness
	th.YamlStructures = make(map[string]interface{})
	th.YamlStructures["https"] = &yamlmodel.Flow {
		Id: "https",
		Name: "HTTP(S)",
		Description: "HTTP on top of TLS",
		Sender: "browser",
		Receiver: "web-server",
		Protocol: []string{"tls@1.2"},
		AdmDir: "/root",
		ADM: []string{"https_oh_no.adm"},
	}
	th.YamlStructures["browser"] = &yamlmodel.Entity {
		Id: "browser",
		Name: "Generic Web Browser",
		Type: yamlmodel.Program,
		Description: "A generic web browser",
	}
	th.YamlStructures["web-server"] = &yamlmodel.Entity {
		Id: "web-server",
		Name: "Generic Web Server",
		Type: yamlmodel.Program,
		Description: "A generic web server",
	}
	th.YamlStructures["tls@1.2"] = &yamlmodel.Flow {
		Id: "tls@1.2",
		Name: "TLS v1.2",
		Description: "TLS v1.2",
	}

	errs := f.Init(th.YamlStructures["https"].(*yamlmodel.Flow), th.Resolve)
	assert.Empty(t, errs)
	assert.Equal(t, []string{"/root/https_oh_no.adm"}, f.GetADM()[f.GetID()])
}

func TestFlowInterface(t *testing.T) {
	var f objmodel.Flow
	var th TestHarness
	th.YamlStructures = make(map[string]interface{})
	th.YamlStructures["https"] = &yamlmodel.Flow {
		Id: "https",
		Name: "HTTP(S)",
		Description: "HTTP on top of TLS",
		Sender: "browser",
		Receiver: "web-server",
		Protocol: []string{"tls@1.2"},
		AdmDir: "/root",
		ADM: []string{"https_oh_no.adm"},
	}
	th.YamlStructures["browser"] = &yamlmodel.Entity {
		Id: "browser",
		Name: "Generic Web Browser",
		Type: yamlmodel.Program,
		Description: "A generic web browser",
	}
	th.YamlStructures["web-server"] = &yamlmodel.Entity {
		Id: "web-server",
		Name: "Generic Web Server",
		Type: yamlmodel.Program,
		Description: "A generic web server",
	}
	th.YamlStructures["tls@1.2"] = &yamlmodel.Flow {
		Id: "tls@1.2",
		Name: "TLS v1.2",
		Description: "TLS v1.2",
	}

	errs := f.Init(th.YamlStructures["https"].(*yamlmodel.Flow), th.Resolve)
	assert.Empty(t, errs)
	
}