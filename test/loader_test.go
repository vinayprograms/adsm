package test

import (
	"errors"
	"fmt"
	admgraph "libadm/graph"
	"libadm/graphviz"
	admloaders "libadm/loaders"
	admmodel "libadm/model"
	"os"
	smloaders "securitymodel/loaders"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithFile(t *testing.T) {
	var l smloaders.Loader
	yaml, err := GetYaml("./examples/simple_addb.smspec")
	assert.Nil(t, err)
	m, err := l.LoadSecurityModel(yaml, "./examples/adm/")
	assert.Zero(t, len(err))
	assert.Equal(t, "Sample security model", m.Title)
	assert.Equal(t, "AnInvalidPath.md", m.DesignDocument)
}

func TestValidSecurityModel(t *testing.T) {
	var l smloaders.Loader
	yaml := `
design-document: "AnInvalidPath.md"
title: Sample security model
addb: ~/addb

externals:
  - id: user
    type: human
    description: User of the system
    interface: user-browser
    roles: [add-record, modify-record, read-record]
  - id: admin
    type: human
    description: Administrator of the system. Controls system's configuration and operations.
    interface: admin-browser
    roles: [add-record, modify-record, read-record, delete-record]
  - id: user-browser
    type: program
    description: User's web-browser
    base: [web-browser]
  - id: admin-browser
    type: program
    description: Administrator's web-browser
    base: [web-browser]

entities:
  - id: frontend
    type: program
    description: System UI consumed by the user.
    repo: nil
    languages: [go]
    adm: null
  - id: backend
    type: program
    description: Server that processes all interactions
    repo: nil
    languages: [go]
    adm: null
  - id: db
    type: program
    description: Database used by backend to persist important data
    repo: nil
    languages: [sql]
    adm: null
flows:
  - id: user-interaction
    description: User's interaction with the frontend
    sender: user-browser
    receiver: frontend
    protocol:   # We are specifying a stack here.
        - http
        - tls
        - tcp
        - ip
    adm: []
  - id: process-requests
    description: System processes user request
    sender: frontend
    receiver: backend
    protocol: [https] # Instead of stack, we refer to the https flow spec, which internally specifies the rest of its stack.
    adm: []
  - id: upsert-user-data
    description: Update or Insert user data
    sender: backend
    receiver: db
    protocol: [sql-flow] # Predefined SQL flow specification
    adm: []
`
	m, err := l.LoadSecurityModel(yaml, "./")
	assert.NotNil(t, err) // Some of the entity dependencies do not resolve (to ADDB or other in-model entities)
	assert.Equal(t, "AnInvalidPath.md", m.DesignDocument)
}

func TestLoadHuman(t *testing.T) {
	var l smloaders.Loader
	yaml1 := `id: user-browser
type: program
name: Generic Web Browser
description: User's web-browser
base: [web-browser]
`
	yaml2 := `id: user
type: human
name: Generic User
description: User of the system
interface: user-browser
`
	_, err := l.LoadProgram(yaml1, "", "")
	assert.GreaterOrEqual(t, len(err), 1) // addb path and 'web-browser' cannot be resolved
	m, err := l.LoadHuman(yaml2, "", "")
	assert.GreaterOrEqual(t, len(err), 1) // addb path and 'web-browser' under interface cannot be resolved
	assert.Equal(t, "user-browser", m.GetUserInterface().GetID())
	assert.Equal(t, "Generic Web Browser", m.GetUserInterface().GetName())
	assert.Equal(t, "User of the system", m.GetDescription())
}

func TestLoadProgram(t *testing.T) {
	var l smloaders.Loader
	yaml := `id: db
type: program
name: Database
description: Database used by backend to persist important data
repo: null
languages: [sql]
adm: null
`
	m, err := l.LoadProgram(yaml, "", "")
	assert.GreaterOrEqual(t, len(err), 1) // 'sql' cannot be resolved
	assert.Equal(t, "db", m.GetID())
	assert.Equal(t, "Database", m.GetName())
	assert.Equal(t, "Database used by backend to persist important data", m.GetDescription())
	assert.Equal(t, "", m.GetCodeRepository())
	assert.ElementsMatch(t, []string{}, m.GetLanguages())
}

func TestLoadFlow(t *testing.T) {
	var l smloaders.Loader
	yaml1 := `id: ui
type: program
name: User Interface
description: System UI consumed by the user.
repo: nil
languages: [go]
adm: null
`
	yaml2 := `id: server
type: program
name: Server
description: Server that processes all interactions
repo: nil
languages: [go]
adm: null
`
	yaml3 := `id: login
name: Login
description: User tries to log into the system
sender: ui
receiver: server
protocol: [https] # Instead of stack, we refer to the https flow spec, which internally specifies the rest of its stack.
adm: []
`
	_, err := l.LoadProgram(yaml1, "", "~/addb")
	assert.GreaterOrEqual(t, len(err), 1) // 'go' cannot be resolved
	_, err = l.LoadProgram(yaml2, "", "~/addb")
	assert.GreaterOrEqual(t, len(err), 1) // 'go' cannot be resolved
	f, err := l.LoadFlow(yaml3, "", "~/addb")
	assert.GreaterOrEqual(t, len(err), 1) // 'https' & 'go' (from sender & receiver) cannot be resolved
	assert.Equal(t, "login", f.GetID())
	assert.Equal(t, "Login", f.GetName())
	assert.Equal(t, "User tries to log into the system", f.GetDescription())
	assert.Equal(t, "ui", f.GetSender().GetID())
	assert.Equal(t, "User Interface", f.GetSender().GetName())
	assert.Equal(t, "server", f.GetReceiver().GetID())
	assert.Equal(t, "Server", f.GetReceiver().GetName())
}

func TestHumanWithADDBReference(t *testing.T) {
	yaml := `id: user
type: human
name: User
description: User of the system
interface: addb:generic.browser`
	var l smloaders.Loader
	_, err := l.LoadHuman(yaml, "", "~/addb")
	assert.Nil(t, err)
}

func TestSecurityModelWithADDBReferences(t *testing.T) {
	// Step-1: Load security model
	var l smloaders.Loader
	yaml, err := GetYaml("./examples/simple_addb.smspec")
	assert.Nil(t, err)
	sm, err := l.LoadSecurityModel(yaml, "./examples")
	assert.Nil(t, err)
	assert.NotNil(t, sm)
	allADM := sm.GetADM()
	assert.Greater(t, len(allADM), 1)

	// Step-2: Load all ADM graphs from the security model (including ADDB)
	var graph admgraph.Graph
	graph.Init()
	for _, adm := range allADM {
		for _, filepath := range adm {
			file, err := os.ReadFile(filepath)
			if err == nil {
				gherkinModel, err1 := admloaders.LoadGherkinContent(string(file))
				if err1 != nil {
					fmt.Println(err1)
				}
				var m admmodel.Model
				err = m.Init(gherkinModel.Feature)
				if err == nil {
					graph.AddModel(&m)
				}
			}
		}
	}

	// Step-3: Generate graph
	config := graphviz.GraphvizConfig{
		Assumption: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "white", FillColor: "dimgray", BorderColor: "dimgray"},
			Font:  graphviz.TextProperties{FontName: "Times", FontSize: "18"},
		},
		Policy: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "black", FillColor: "darkolivegreen3", BorderColor: "darkolivegreen3"},
			Font:  graphviz.TextProperties{FontName: "Times", FontSize: "18"},
		},
		PreConditions: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "black", FillColor: "lightgray", BorderColor: "gray"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "14"},
		},

		// Defense config
		PreEmptiveDefense: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "white", FillColor: "purple", BorderColor: "purple"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "14"},
		},
		IncidentResponse: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "white", FillColor: "blue", BorderColor: "blue"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "14"},
		},
		EmptyDefense: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "black", FillColor: "transparent", BorderColor: "blue"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "14"},
		},

		// Attack config
		Attack: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "white", FillColor: "red", BorderColor: "red"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "14"},
		},
		EmptyAttack: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "black", FillColor: "transparent", BorderColor: "red"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "14"},
		},

		// Start and end node config
		Reality: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "white", FillColor: "black", BorderColor: "black"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "20"},
		},
		AttackerWins: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "red", FillColor: "yellow", BorderColor: "yellow"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "20"},
		},

		Subgraph: graphviz.TextProperties{FontName: "Arial", FontSize: "24"},
	}
	lines, graphErr := graphviz.GenerateGraphvizCode(&graph, config)
	assert.Nil(t, graphErr)
	output := fmt.Sprintf("%v", lines)
	writeErr := os.WriteFile("/tmp/result.dot", []byte(output[1:len(output)-1]), 0777)
	assert.Nil(t, writeErr)

	// Step-4: Generate security model report
	// TODO
}

////////////////////////////////////////
// Helper functions

func GetYaml(path string) (string, []error) {
	yamlData, err := getFileContents(path)
	if err != nil {
		return "", []error{err}
	}

	return string(yamlData), nil
}

// Load content from a file path
func getFileContents(filepath string) ([]byte, error) {
	if filepath == "" {
		return nil, errors.New("no path provided or invalid path")
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("no content found in file")
	}

	return data, nil
}
