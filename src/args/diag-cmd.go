package args

import (
	"fmt"
	"libadm/graph"
	"libadm/graphviz"
	admloaders "libadm/loaders"
	"libadm/model"
	"os"
	"securitymodel/diagram"
	"securitymodel/objmodel"
	"strings"
)

type generateAdmCommand struct {
	model      objmodel.SecurityModel
	outputpath string
}

type generateSmCommand struct {
	model      objmodel.SecurityModel
	outputpath string
}

////////////////////////////////////////
// 'execute()' implementation for each command

// Generate consolidated ADM graph from all ADM files listed in a security model
func (g generateAdmCommand) execute() error {
	var graph graph.Graph
	graph.Init()

	var allADM []string

	for _, entity := range g.model.Entities {
		for _, all := range entity.GetADM() {
			allADM = append(allADM, all...)
		}
	}
	for _, flow := range g.model.Flows {
		for _, all := range flow.GetADM() {
			allADM = append(allADM, all...)
		}
	}

	for _, model := range getADMModels(allADM) {
		err := graph.AddModel(model)
		if err != nil {
			fmt.Println(err)
		}
	}

	code, err := graphviz.GenerateGraphvizCode(&graph, getConfig())
	if err != nil {
		return err
	}

	output := strings.Join(code, "\n")
	if g.outputpath[len(g.outputpath)-1] != '/' { // append a '/' if path doesn't have it
		g.outputpath += "/"
	}
	checkAndCreateDirectory(g.outputpath)
	err = os.WriteFile(g.outputpath+diagram.GenerateID(g.model.Title)+".adm.dot", []byte(output), 0777)
	if err != nil {
		return err
	}

	return nil
}

// Generate security model diagram along with mitigated and unmitigated attacks count for each entity and flow.
func (g generateSmCommand) execute() error {
	lines, err := diagram.GenerateSMDiagram(g.model)
	if err != nil {
		return err
	}

	output := strings.Join(lines, "\n")
	if g.outputpath[len(g.outputpath)-1] != '/' { // append a '/' if path doesn't have it
		g.outputpath += "/"
	}
	checkAndCreateDirectory(g.outputpath)
	err = os.WriteFile(g.outputpath+diagram.GenerateID(g.model.Title)+".sm.dot", []byte(output), 0777)
	if err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////
// Helper functions

// Get a list of model objects from a list of adm file paths
func getADMModels(allADM []string) (models []*model.Model) {
	for _, admFile := range allADM {
		model := getADM(admFile)
		if model == nil {
			continue
		}
		models = append(models, model)
	}
	return
}

// Load a single ADM file into a model object
func getADM(file string) *model.Model {
	contents, err := getFileContent(file)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	if len(contents) == 0 { //no contents
		fmt.Println("No ADM content found in " + file)
		return nil
	}

	gherkinModel, err := admloaders.LoadGherkinContent(contents)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var m model.Model
	err = m.Init(gherkinModel.Feature)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &m
}

// Configuration data for use in graph diagram generation.
func getConfig() graphviz.GraphvizConfig {
	return graphviz.GraphvizConfig{
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
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "16"},
		},

		// Defense config
		PreEmptiveDefense: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "white", FillColor: "purple", BorderColor: "blue"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "16"},
		},
		IncidentResponse: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "white", FillColor: "blue", BorderColor: "blue"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "16"},
		},
		EmptyDefense: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "black", FillColor: "transparent", BorderColor: "blue"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "16"},
		},

		// Attack config
		Attack: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "white", FillColor: "red", BorderColor: "red"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "16"},
		},
		EmptyAttack: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "black", FillColor: "transparent", BorderColor: "red"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "16"},
		},

		// Start and end node config
		Reality: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "white", FillColor: "black", BorderColor: "black"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "20"},
		},
		AttackerWins: graphviz.NodeProperties{
			Color: graphviz.ColorSet{FontColor: "red", FillColor: "yellow", BorderColor: "red"},
			Font:  graphviz.TextProperties{FontName: "Arial", FontSize: "20"},
		},

		Subgraph: graphviz.TextProperties{FontName: "Arial", FontSize: "24"},
	}
}
