package args

import (
	"fmt"

	"libadm/graph"
	admloaders "libadm/loaders"
	admmodel "libadm/model"
	"libsm/diagram"
	"libsm/objmodel"
	"os"
	"strings"

	"github.com/goccy/go-graphviz"
)

type generateReportCommand struct {
	model objmodel.SecurityModel
	outputpath string
}

////////////////////////////////////////
// 'execute()' implementation

func (g generateReportCommand) execute() error {
	// Generate report
	markdownReport := strings.Join(generateReport(g.model), "\n")
	outpath := checkAndCreateDirectory(g.outputpath)
	outpath = checkAndCreateDirectory(outpath + "report")
	err := os.WriteFile(outpath + diagram.GenerateID(g.model.Title) + ".sm.md", []byte(markdownReport), 0777)
	if err != nil {
		return err
	}
	
	// export security model diagram (used in report)
	generateSmCommand { model: g.model, outputpath: outpath + "resources"}.execute()
	exportPNG(outpath + "resources/" + diagram.GenerateID(g.model.Title))

	// export ADM (linked to in report)
	generateAdmCommand { model: g.model, outputpath: outpath + "resources"}.execute()
	
	return nil
}

////////////////////////////////////////
// Functions to generate report content

func exportPNG(dotFilepath string) error {
	// 1. Read 'dot' file
	bytes, err := os.ReadFile(dotFilepath + ".sm.dot")
	if err != nil { return err }
	graph, err := graphviz.ParseBytes(bytes)
	if err != nil { return err }
	
	// 2. Write png to target
	g := graphviz.New()
	if err := g.RenderFilename(graph, graphviz.PNG, dotFilepath + ".sm.png"); err != nil {
		return err
	}

	// 3. Remove 'dot' file
	err = os.Remove(dotFilepath + ".sm.dot")
	if err != nil { return err }

	return nil
}

func generateReport(model objmodel.SecurityModel) (markdownLines []string) {
	markdownLines = append(markdownLines, "# Security Report: " + model.Title)
	markdownLines = appendLineSpacer(markdownLines)
	markdownLines = append(markdownLines, "This report contains")
	markdownLines = appendLineSpacer(markdownLines)
	markdownLines = append(markdownLines, "* Existing mitigations implemented in specific entities/flows in this security model.")
	markdownLines = append(markdownLines, "* Security recommendations for specific entities/flows in this security model.")
	markdownLines = append(markdownLines, "* A list of un-mitigated risks for specific entities/flows.")
	markdownLines = appendLineSpacer(markdownLines)
	
	// Security model
	markdownLines = append(markdownLines, "## Security Model")
	markdownLines = appendLineSpacer(markdownLines)
	markdownLines = append(markdownLines, "![security-model](resources/" + diagram.GenerateID(model.Title) + ".sm.png)")
	markdownLines = appendLineSpacer(markdownLines)
	markdownLines = append(markdownLines, 
		"The Attack-Defense Graph for this model is available as a [graphviz file](" + 
		"resources/" + diagram.GenerateID(model.Title) + ".adm.dot). " + 
		"Please use a graphviz viewer or use [graphviz CLI tool](https://graphviz.org/download/) " +
		"to export it to an image format of your choice. " +
		"In case of CLI tool use `dot -Tpng resources/" + diagram.GenerateID(model.Title) + ".adm.dot` " + 
		"to generate a PNG image of the graph. " +
		"Detailed user documentation for CLI tool is available [here](https://graphviz.org/doc/info/command.html).")
		markdownLines = appendLineSpacer(markdownLines)

	// List risks
	risks := generateRisksSection(model)
	if len(risks) > 0 {
		markdownLines = append(markdownLines, "## Risks")
		markdownLines = appendLineSpacer(markdownLines)
		markdownLines = append(markdownLines, "This section lists all ADM attacks that have not been mitigated.")
		markdownLines = appendLineSpacer(markdownLines)	
		markdownLines = append(markdownLines, risks...)
		markdownLines = appendLineSpacer(markdownLines)	
	}

	// List mitigations
	mitigations := generateMitigationsSection(model)
	if len(mitigations) > 0 {
		markdownLines = append(markdownLines, "## Mitigations")
		markdownLines = appendLineSpacer(markdownLines)
		markdownLines = append(markdownLines, "This section lists mitigations present in entities/flows in this security model.")
		markdownLines = appendLineSpacer(markdownLines)	
		markdownLines = append(markdownLines, mitigations...)
		markdownLines = appendLineSpacer(markdownLines)	
	}
	
	// List recommendations
	recos := generateRecommendationsSection(model)
	if len(recos) > 0 {
		markdownLines = append(markdownLines, "## Recommendations")
		markdownLines = appendLineSpacer(markdownLines)
		markdownLines = append(markdownLines, "This section lists general security recommendations for specific entities/flows in this security model.")
		markdownLines = appendLineSpacer(markdownLines)	
		markdownLines = append(markdownLines, recos...)
		markdownLines = appendLineSpacer(markdownLines)	
	}
	
	return
}

func generateMitigationsSection(model objmodel.SecurityModel) (markdownLines []string) {
	// entities
	var entitiesSectionContent []string
	for _, entity := range model.Entities {
		if _, ok := entity.(*objmodel.Role); ok { // Roles will be processed as part of entity that uses it.
			continue
		}
		mitigations := entity.GetMitigations()
		if len(mitigations) > 0 {	
			entitiesSectionContent = appendLineSpacer(entitiesSectionContent)
			entitiesSectionContent = append(entitiesSectionContent, "#### " + entity.GetName())
			entitiesSectionContent = appendLineSpacer(entitiesSectionContent)
			for mitiSource, mitis := range mitigations {
				for _, miti := range mitis {
					if mitiSource == entity.GetName() { // Skip showing the source if it is the root entity.
						entitiesSectionContent = append(entitiesSectionContent, "* " + miti)
					} else {
						entitiesSectionContent = append(entitiesSectionContent, "* (`" + mitiSource + "`) " + miti)
					}
				}
			}
		}
	}
	if len(entitiesSectionContent) > 0 {
		markdownLines = append(markdownLines, "### Entities")
		markdownLines = append(markdownLines, entitiesSectionContent...)
	}

	// flows
	var flowsSectionContent []string
	for _, flow := range model.Flows {
		mitigations := flow.GetMitigations()
		if len(mitigations) > 0 {
			flowsSectionContent = appendLineSpacer(flowsSectionContent)	
			flowsSectionContent = append(flowsSectionContent, "#### " + flow.GetName())
			flowsSectionContent = appendLineSpacer(flowsSectionContent)
			for mitiSource, mitis := range mitigations {
				for _, miti := range mitis {
					if mitiSource == flow.GetName() { // Skip showing the source if it is the root entity.
						flowsSectionContent = append(flowsSectionContent, "* " + miti)
					} else {
						flowsSectionContent = append(flowsSectionContent, "* (`" + mitiSource + "`) " + miti)
					}
				}
			}
		}
	}
	if len(flowsSectionContent) > 0 {
		markdownLines = append(markdownLines, "### Flows")
		markdownLines = append(markdownLines, flowsSectionContent...)
	}
	return
}

func generateRecommendationsSection(model objmodel.SecurityModel) (markdownLines []string) {
	// entities
	var entitiesSectionContent []string
	for _, entity := range model.Entities {
		if _, ok := entity.(*objmodel.Role); ok { // Roles will be processed as part of entity that uses it.
			continue
		}
		recommendations := entity.GetRecommendations()
		if len(recommendations) > 0 {	
			entitiesSectionContent = appendLineSpacer(entitiesSectionContent)
			entitiesSectionContent = append(entitiesSectionContent, "#### " + entity.GetName())
			entitiesSectionContent = appendLineSpacer(entitiesSectionContent)
			for recoSource, recos := range recommendations {
				for _, reco := range recos {
					if recoSource == entity.GetName() { // Skip showing the source if it is the root entity.
						entitiesSectionContent = append(entitiesSectionContent, "* " + reco)
					} else {
						entitiesSectionContent = append(entitiesSectionContent, "* (`" + recoSource + "`) " + reco)
					}
				}
			}
		}
	}
	if len(entitiesSectionContent) > 0 {
		markdownLines = append(markdownLines, "### Entities")
		markdownLines = append(markdownLines, entitiesSectionContent...)
	}

	// flows
	var flowsSectionContent []string
	for _, flow := range model.Flows {
		recommendations := flow.GetRecommendations()
		if len(recommendations) > 0 {
			flowsSectionContent = appendLineSpacer(flowsSectionContent)	
			flowsSectionContent = append(flowsSectionContent, "#### " + flow.GetName())
			flowsSectionContent = appendLineSpacer(flowsSectionContent)
			for recoSource, recos := range recommendations {
				for _, reco := range recos {
					if recoSource == flow.GetName() { // Skip showing the source if it is the root entity.
						flowsSectionContent = append(flowsSectionContent, "* " + reco)
					} else {
						flowsSectionContent = append(flowsSectionContent, "* (`" + recoSource + "`) " + reco)
					}
				}
			}
		}
	}
	if len(flowsSectionContent) > 0 {
		markdownLines = append(markdownLines, "### Flows")
		markdownLines = append(markdownLines, flowsSectionContent...)
	}
	return
}

func generateRisksSection(model objmodel.SecurityModel) (markdownLines []string) {
	var graph graph.Graph
	graph.Init()

	attackMap := make(map[string][]string) // maps attack titles to the qualified-name of security-model item
	for qualifiedName, admList := range model.GetADM() {
		for _, admFile := range admList {
			contents, err := os.ReadFile(admFile)
			if err != nil { 
				fmt.Println(err.Error()) 
				continue
			}
			if len(contents) == 0 { //no contents
				fmt.Println("No ADM content found in " + admFile)
				continue
			}
			
			gherkinModel, err := admloaders.LoadGherkinContent(string(contents))
			if err != nil { 
				fmt.Println(err) 
				continue
			}
			
			var m admmodel.Model
			err = m.Init(gherkinModel.Feature)
			if err != nil { 
				fmt.Println(err) 
				continue
			}
			for attackTitle := range m.Attacks {
				attackMap[attackTitle] = append(attackMap[attackTitle], qualifiedName)
			}

			err = graph.AddModel(&m)
			if err != nil { 
				fmt.Println(err) 
				continue
			}
		}
	}
	for risk := range graph.AttackerWinsPredecessors {
		if attackMap[risk] == nil {
			// CAUTION: This line should never be reached. If it does, contact author.
			fmt.Println("ERROR: Cannot find attack - '" + risk + "' among all attacks listed for this security model.")
		}
		for _, qualifiedName := range attackMap[risk] {
			qualifiedName := strings.ReplaceAll(qualifiedName, "sm.", "")
			location := strings.ReplaceAll(qualifiedName, ".", " → ")
			markdownLines = append(markdownLines, "* " + risk + " (under `" + location + "`)")
		}
	}
	return
}

////////////////////////////////////////
// Helper Functions

func appendLineSpacer(document []string) []string {
	return append(document, "")
}