package args

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"securitymodel/loaders"
)

func statsInvoker(x bool, e bool, r bool, f bool, path string) error {
	err := checkPath(path)
	if err != nil {
		return errors.New("error when verifying path - '" + path + "'")
	}

	// Special case: If all flags are 'false' (i.e., none were specified) turn all of them to 'true'
	if !(x || e || r || f) {
		x = true
		e = true
		r = true
		f = true
	}

	models, err := getContent(path)
	if err != nil {
		return err
	}
	for _, modelText := range models {
		var l loaders.Loader
		model, errs := l.LoadSecurityModel(modelText, filepath.Dir(path))
		PrintErrors(errs)                    // send errors to STDOUT
		fmt.Println("MODEL: " + model.Title) // Print the title once (not for each flag)
		for _, adm := range model.GetADM()["sm"] {
			printADMStatLine(adm)
		}
		if x {
			externalStatsCommand{model: *model}.execute()
		}
		if e {
			entityStatsCommand{model: *model}.execute()
		}
		if r {
			roleStatsCommand{model: *model}.execute()
		}
		if f {
			flowStatsCommand{model: *model}.execute()
		}
	}

	return nil
}

func diagInvoker(sm bool, adm bool, outPath string, path string) error {
	err := checkPath(path)
	if err != nil {
		return errors.New("error when verifying path - '" + path + "'")
	}
	models, err := getContent(path)
	if err != nil {
		return err
	}
	for _, modelText := range models {
		var l loaders.Loader
		model, errs := l.LoadSecurityModel(modelText, filepath.Dir(path))
		PrintErrors(errs) // send errors to STDOUT

		if sm {
			generateSmCommand{model: *model, outputpath: outPath}.execute()
		}
		if adm {
			generateAdmCommand{model: *model, outputpath: outPath}.execute()
		}
	}

	return nil
}

func reportInvoker(outPath string, path string) error {
	err := checkPath(path)
	if err != nil {
		return errors.New("error when verifying path - '" + path + "'")
	}
	models, err := getContent(path)
	if err != nil {
		return err
	}
	for _, modelText := range models {
		var l loaders.Loader
		model, errs := l.LoadSecurityModel(modelText, filepath.Dir(path))
		PrintErrors(errs) // send errors to STDOUT

		generateReportCommand{model: *model, outputpath: outPath}.execute()
	}

	return nil
}

/*
func exportInvoker(e bool, r bool, f bool, outPath string, path string) error {
	err := checkPath(path)
	if err != nil {
		return errors.New("error when verifying path - '" + path + "'")
	}

	return nil
}*/

////////////////////////////////////////
// Helper functions

func checkPath(path string) error {
	_, err := os.Stat(path)
	return err
}
