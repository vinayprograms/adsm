package test

import (
	"libsm/objmodel"
	"libsm/yamlmodel"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoleWithNullYaml(t *testing.T) {
	var r objmodel.Role
	var th TestHarness
	errs := r.Init(nil, th.Resolve)
	assert.Equal(t, len(errs), 1)
	assert.Equal(t, errs[0].Error(),"cannot convert nil yaml to role specification")
}

func TestRoleWithEmptySpec(t *testing.T) {
	var r1, r2 objmodel.Role
	var th TestHarness
	role1 := yamlmodel.Entity {
		Id: "",
		Type: yamlmodel.Role,
		Description: "",
	}
	role2 := yamlmodel.Entity {
		Id: "some-program",
		Type: yamlmodel.Role,
		Description: "",
	}
	errs := r1.Init(&role1, th.Resolve)
	assert.Equal(t, "empty IDs are not allowed", errs[0].Error())
	
	errs = r2.Init(&role2, th.Resolve)
	assert.Equal(t, "empty names are not allowed", errs[0].Error())
}

func TestRoleWithNonRoleYaml(t *testing.T) {
	var r objmodel.Role
	var th TestHarness
	role := yamlmodel.Entity {
		Id: "iron-man",
		Type: yamlmodel.Human,
		Description: "Human masquerading as an Android",
		ADM: []string{"careless.adm"},
	}
	errs := r.Init(&role, th.Resolve)
	assert.Equal(t, len(errs), 1)
	assert.Equal(t, errs[0].Error(),"cannot initialize role with 'human'")
}

func TestRoleWithEmptyDescription(t *testing.T) {
	var r objmodel.Role
	var th TestHarness
	role := yamlmodel.Entity {
		Id: "admin",
		Name: "Administrator",
		Type: yamlmodel.Role,
		Description: "",
		ADM: []string{"admin.adm"},
	}
	errs := r.Init(&role, th.Resolve)
	assert.Equal(t, len(errs), 1)
	assert.Equal(t, errs[0].Error(),"warning... empty descriptions are useless")
}

func TestFullyDefinedRole(t *testing.T) {
	var r objmodel.Role
	var th TestHarness
	role := yamlmodel.Entity {
		Id: "administrator",
		Type: yamlmodel.Role,
		Name: "Administrator",
		Description: "Administrator role. Has the full set of privileges to the system.",
		Recommendations: []string {"Do not give admin rights to everyone."},
		ADM: []string{"admin.adm"},
	}
	errs := r.Init(&role, th.Resolve)
	assert.Nil(t, errs)
	assert.Equal(t, r.GetID(),"administrator")
	assert.Equal(t, r.GetName(),"Administrator")
	assert.Equal(t, r.GetDescription(),"Administrator role. Has the full set of privileges to the system.")
	assert.Equal(t, 1, len(r.GetRecommendations()))
	assert.Contains(t, r.GetRecommendations()["Administrator"],"Do not give admin rights to everyone.")
	assert.Equal(t, 1, len(r.GetADM()))
	assert.Contains(t, r.GetADM()[r.GetID()],"admin.adm")
}

func TestRoleWithAdmInAnotherDirectory(t *testing.T) {
	var rol objmodel.Role
	var th TestHarness
	role := yamlmodel.Entity {
		Id: "administrator",
		Type: yamlmodel.Role,
		Name: "Administrator",
		Description: "Administrator role. Has the full set of privileges to the system.",
		Recommendations: []string {"Do not give admin rights to everyone."},
		AdmDir: "/root",
		ADM: []string{"admin.adm"},
	}
	errs := rol.Init(&role, th.Resolve)
	assert.Empty(t, errs)
	assert.Equal(t, []string{"/root/admin.adm"}, rol.GetADM()[rol.GetID()])
}

func TestInterfaceFunctions(t *testing.T) {
	var rol objmodel.Role
	var th TestHarness
	role := yamlmodel.Entity {
		Id: "administrator",
		Type: yamlmodel.Role,
		Name: "Administrator",
		Description: "Administrator role. Has the full set of privileges to the system.",
		Recommendations: []string {"Do not give admin rights to everyone."},
		ADM: []string{"admin.adm"},
	}
	errs := rol.Init(&role, th.Resolve)
	assert.Empty(t, errs)
	rol.AddADM("privileges.adm")
	assert.Equal(t, 2, len(rol.GetADM()[rol.GetID()]))
	assert.Contains(t, rol.GetADM()[rol.GetID()], "privileges.adm")
	rol.AddRecommendation("Administrator operations must be logged at higher level of detail than regular users.")
	assert.Equal(t, 2, len(rol.GetRecommendations()["Administrator"]))
	assert.Contains(t, rol.GetRecommendations()["Administrator"], "Administrator operations must be logged at higher level of detail than regular users.")
}