package diagram

import (
	"fmt"
	"libadm/graph"
	admloaders "libadm/loaders"
	"libadm/model"
	"libsm/objmodel"
	"os"
)

func GenerateSMDiagram(model objmodel.SecurityModel) ([]string, error) {
	// Generate graphviz code
	var lines []string
	lines = append(lines, generateHeader()...)
	lines = append(lines, generateBody(model)...)
	lines = append(lines, generateFooter()...)
	return lines, nil
}

func GenerateExternalEntityCode(id string, ext objmodel.ExternalSpec) string {
	extProperties := " style=\"rounded\" shape=\"box\" fontname=\"Arial\"];"
	return id + "[label=\"" + wrap(ext.GetName()) + "\" " + extProperties
}

func GenerateEntityCode(id string, entity objmodel.EntitySpec) string {
	riskyEntityProperties := " style=\"rounded\" shape=\"record\" fontname=\"Arial\" color=\"red\" penwidth=\"2\"];"
	safeEntityProperties := " style=\"rounded\" shape=\"record\" fontname=\"Arial\"];"
	
	extraA := 0
	extraD := 0
	extraM := 0
	extraR := 0
	if prog, ok := entity.(*objmodel.Program); ok {
		for _, role := range prog.GetRoles() { // Consolidate role stats into entity that uses it.
			a, d, m, r, _ := getStatsForEntity(role)
			extraA += a
			extraD += d
			extraM += m
			extraR += r
		}
	}
	a, d, m, r, hasRisks := getStatsForEntity(entity)
	label := "{" + wrap(entity.GetName()) + "} | {A: " + fmt.Sprint(a + extraA) + "| D: " + fmt.Sprint(d + extraD) + "| M: " + fmt.Sprint(m + extraM) + "| R: " + fmt.Sprint(r + extraR) + "}"
	if hasRisks {
		return GenerateID(id) + "[label=\"" + label + "}\" " + riskyEntityProperties
	} else {
		return GenerateID(id) + "[label=\"" + label + "}\" " + safeEntityProperties
	}
}

func GenerateFlowCode(flow objmodel.FlowSpec, externalIDs []string) string {
	externalFlowProperties := " fontname=\"Arial\" fontcolor=\"blue\" fontsize=\"10\" decorate=\"true\""
	riskyFlowProperties := " fontname=\"Arial\" fontcolor=\"red\" fontsize=\"10\" decorate=\"true\"  color=\"red\" penwidth=\"2\""
	safeFlowProperties := " fontname=\"Arial\" fontcolor=\"blue\" fontsize=\"10\" decorate=\"true\""

	a, d, m, r, hasRisks := getStatsForFlow(flow)
	label := "<<b>" + htmlwrap(flow.GetName()) + "</b><br/>A: " + fmt.Sprint(a) + " | D: " + fmt.Sprint(d) + " | M: " + fmt.Sprint(m) + " | R: " + fmt.Sprint(r) + ">"
	
	sender := flow.GetSender()
	if sender == nil { return "" }
	senderID := sender.GetID()
	if senderID == "" { return ""}
	senderID = GenerateID(senderID)
	
	receiver := flow.GetReceiver()
	if receiver == nil { return "" }
	receiverID := receiver.GetID()
	if receiverID == "" { return "" }
	receiverID = GenerateID(receiverID)
	
	// flow from external entity, into the system
	if contains(senderID, externalIDs) {
		return senderID +  " -> " +  receiverID + "[label=" + label + " " + externalFlowProperties + "]"
	} else if hasRisks {
		return senderID +  " -> " +  receiverID + "[label=" + label + " " + riskyFlowProperties + "]"
	} else {
		return senderID +  " -> " +  receiverID + "[label=" + label + " " + safeFlowProperties + "]"
	}
}

////////////////////////////////////////
// Internal functions that build parts of the diagram

// Define main di-graph and its properties
func generateHeader() (header []string) {
	header = appendLine(header, 0, "digraph {")
	header = appendLine(header, 0, "// Base Styling")
	header = appendLine(header, 0, "compound=true")
	header = appendLine(header, 0, "graph[style=\"filled, rounded\" rankdir=\"LR\" splines=\"true\" overlap=\"false\" nodesep=\"0.5\" ranksep=\"0.5\" fontname=\"Arial\"];")
	return
}

// Add all items from the Security model and their flows.
func generateBody(model objmodel.SecurityModel) (body []string) {
	
	// Add externals
	body = appendLine(body, 1, "//externals")
	var externIDs []string
	for id, ext := range model.Externals {
		if id == "" || ext == nil { continue }
		externIDs = append(externIDs, GenerateID(id))
		body = appendLine(body, 1, GenerateExternalEntityCode(GenerateID(id), ext))
	}
	body = appendLineSpacer(body)

	// Add entities
	body = appendLine(body, 1, "//entities")
	body = appendLine(body, 1, "subgraph cluster_" + GenerateID(model.Title) + "{")
	graphProperties := " style=\"filled, rounded, dashed\" rankdir=\"LR\" splines=\"true\" overlap=\"false\" nodesep=\"0.5\" ranksep=\"0.5\" fontname=\"Arial\" fontcolor=\"black\"  fillcolor=\"transparent\"  color=\"red\"];"
	body = appendLine(body, 2, "graph[label=<<b>" + htmlwrap(model.Title) + "</b>>" + graphProperties)
	for id, entity := range model.Entities {
		if id == "" || entity == nil { continue }
		if _, ok := entity.(*objmodel.Role); ok {
			continue // Role data will be consolidated into the entity that uses it.
		}
		line := GenerateEntityCode(id, entity)
		body = appendLine(body, 2, line)
	}
	body = appendLine(body, 1, "}")
	body = appendLineSpacer(body)

	// Add flows
	body = appendLine(body, 1, "//flows")
	for _, flow := range model.Flows {
		if flow == nil { continue }
		body = appendLine(body, 1, GenerateFlowCode(flow, externIDs))
	}

	return
}

// Generate finishing code and close the digraph.
func generateFooter() (footer []string) {
	footer = appendLine(footer, 0, "}")
	return
}

////////////////////////////////////////
// Helper functions that handle all the logic related to generating code.

// For a set of ADM files, generate 'attacks, defenses' statistic

func getStats(allADM map[string][]string) (attacks int, defenses int, hasUnmitigatedAttacks bool) {
	var graph graph.Graph
	graph.Init()

	attacks = 0
	defenses = 0

	for _, admList := range allADM {
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
			
			var m model.Model
			err = m.Init(gherkinModel.Feature)
			if err != nil { 
				fmt.Println(err) 
				continue
			}
			attacks += len(m.Attacks)
			defenses += len(m.Defenses)

			err = graph.AddModel(&m)
			if err != nil { 
				fmt.Println(err) 
				continue
			}
		}
	}

	hasUnmitigatedAttacks = len(graph.AttackerWinsPredecessors) > 0

	return
}

// Generate statistics for an entity
func getStatsForEntity(m objmodel.EntitySpec) (attacks int, defenses int, mitigationsCount int, recommendationsCount int, hasOpenRisks bool) {
	attacks, defenses, hasOpenRisks = getStats(m.GetADM())
	mitigationsCount = 0
	recommendationsCount = 0
	for _, mitigations := range m.GetMitigations() {
		mitigationsCount += len(mitigations)
	}
	for _, recos := range m.GetRecommendations() {
		recommendationsCount += len(recos)
	}
	return 
}

// Generate statistics for a flow
func getStatsForFlow(f objmodel.FlowSpec) (attacks int, defenses int,  mitigationsCount int, recommendationsCount int,  hasOpenRisks bool) {
	attacks, defenses, hasOpenRisks = getStats(f.GetADM())
	mitigationsCount = 0
	recommendationsCount = 0
	for _, mitigations := range f.GetMitigations() {
		mitigationsCount += len(mitigations)
	}
	for _, recos := range f.GetRecommendations() {
		recommendationsCount += len(recos)
	}
	return 
}