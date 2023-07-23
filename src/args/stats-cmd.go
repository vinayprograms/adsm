package args

import (
	"fmt"
	"os"
	"securitymodel/objmodel"

	admloaders "libadm/loaders"
	admmodel "libadm/model"
)

type externalStatsCommand struct {
	model objmodel.SecurityModel
}

type entityStatsCommand struct {
	model objmodel.SecurityModel
}

type roleStatsCommand struct {
	model objmodel.SecurityModel
}

type flowStatsCommand struct {
	model objmodel.SecurityModel
}

////////////////////////////////////////
// 'execute()' implementation for each command

func (e externalStatsCommand) execute() error {
	for _, ext := range e.model.Externals {
		fmt.Println("\tExternal Entity: " + ext.GetName())
		fmt.Println("\t                 " + ext.GetDescription())
	}
	return nil
}

func (e entityStatsCommand) execute() error {
	for _, ent := range e.model.Entities {
		_, ok1 := ent.(*objmodel.Human)
		_, ok2 := ent.(*objmodel.Program)
		if !(ok1 || ok2) {
			continue
		}
		fmt.Println("\tEntity: " + ent.GetName())
		fmt.Println("\t        " + ent.GetDescription())
		for _, adm := range ent.GetADM() {
			for _, admFilePath := range adm {
				line := printADMStatLine(admFilePath)
				if line != "" {
					fmt.Println("\t        " + line)
				}
			}
		}
	}
	return nil
}

func (r roleStatsCommand) execute() error {
	for _, rol := range r.model.Entities {
		if _, ok := rol.(*objmodel.Role); ok {
			fmt.Println("\tRole: " + rol.GetName())
			fmt.Println("\t      " + rol.GetDescription())

			for _, adm := range rol.GetADM() {
				for _, admFilePath := range adm {
					line := printADMStatLine(admFilePath)
					if line != "" {
						fmt.Println("\t      " + line)
					}
				}
			}
		}
	}
	return nil
}

func (f flowStatsCommand) execute() error {
	for _, flo := range f.model.Flows {
		fmt.Println("\tFlow: " + flo.GetName())
		fmt.Println("\t      " + flo.GetDescription())
		for _, adm := range flo.GetADM() {
			for _, admFilePath := range adm {
				line := printADMStatLine(admFilePath)
				if line != "" {
					fmt.Println("\t      " + line)
				}
			}
		}
	}
	return nil
}

func printADMStatLine(file string) (line string) {
	content, err := os.ReadFile(file)
	if err == nil {
		gherkinModel, err1 := admloaders.LoadGherkinContent(string(content))
		if err1 != nil {
			fmt.Println(err1)
		}
		var m admmodel.Model
		err2 := m.Init(gherkinModel.Feature)
		if err2 == nil {
			line += "ADM: " + file
			if len(m.Assumptions) > 0 {
				line += ", ASSUMPTIONS:" + fmt.Sprint(len(m.Assumptions))
			}
			if len(m.Attacks) > 0 {
				line += ", ATTACKS:" + fmt.Sprint(len(m.Attacks))
			}
			if len(m.Defenses) > 0 {
				line += ", DEFENSES:" + fmt.Sprint(len(m.Defenses))
			}
			if len(m.Policies) > 0 {
				line += ", POLICIES:" + fmt.Sprint(len(m.Policies))
			}
		}
	}
	return
}
