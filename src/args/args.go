package args

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
)

type Args struct {
	statCmd			*flag.FlagSet
	diagCmd   	*flag.FlagSet
	reportCmd  	*flag.FlagSet
	//exportCmd  	*flag.FlagSet
	path      	string
}

func (a *Args) isInitialized() bool {
	return (a.statCmd != nil)
}

func (a *Args) InitArgs() {
	a.statCmd = flag.NewFlagSet("stat", flag.ExitOnError)
	a.statCmd.Bool("x", false, "List external entities only.")
	a.statCmd.Bool("e", false, "List in-scope entities only.")
	a.statCmd.Bool("r", false, "List roles only.")
	a.statCmd.Bool("f", false, "List flows only.")

	a.diagCmd = flag.NewFlagSet("diag", flag.ExitOnError)
	a.diagCmd.Bool("sm", false, "Generate security model diagram only.")
	a.diagCmd.Bool("adm", false, "Generate ADM decision graph only.")
	a.diagCmd.String("d", "./", "Output directory for diagrams.")

	a.reportCmd = flag.NewFlagSet("report", flag.ExitOnError)
	a.reportCmd.String("d", "./", "Output directory for generated report.")
	/*
	a.exportCmd = flag.NewFlagSet("export", flag.ExitOnError)
	a.exportCmd.String("f", "json", "Output format. Supported values - json,xml.")
	a.exportCmd.String("d", "./", "Output directory for exported files.")*/
}

func (a *Args) PrintHelpToStdout() {

	if !a.isInitialized() {
		a.InitArgs()
	}

	fmt.Println("\nUsage: adsm [OPTIONS] [PATH]")
	fmt.Println("\n[PATH]: The path to a directory or a single smspec file")
	
	fmt.Println("\nstat: List model components")
	a.statCmd.PrintDefaults()
	
	fmt.Println("\ndiag: Generate security model and ADM diagrams.")
	a.diagCmd.PrintDefaults()

	fmt.Println("\nreport: Generate security report as markdown file.")
	a.reportCmd.PrintDefaults()

	/*fmt.Println("\nexport: Export security model and report to other formats.")
	a.exportCmd.PrintDefaults()*/
}

func (a Args) ParseArgs(args []string) error {
	if !a.isInitialized() {
		a.InitArgs()
	}

	switch len(args) {
	case 0:
		a.PrintHelpToStdout()
		return nil
	case 1:
		return errors.New("require atleast two parameters - 'sub-command' and 'path'")
	}

	a.path = args[len(args)-1]	// Path must always be the last parameter

	// Process the remaining arguments.
	switch args[0] {
	case "stat":
		err := a.statCmd.Parse(args[1:len(args)-1])
		if err != nil {
			// Control should not reach here. Parse typically does a 'os.Exit()' if something goes wrong.
			// If you do reach, contact author.
			return err
		}
		xFlag, _ := strconv.ParseBool(a.statCmd.Lookup("x").Value.String())
		eFlag, _ := strconv.ParseBool(a.statCmd.Lookup("e").Value.String())
		rFlag, _ := strconv.ParseBool(a.statCmd.Lookup("r").Value.String())
		fFlag, _ := strconv.ParseBool(a.statCmd.Lookup("f").Value.String())
		
		return statsInvoker(xFlag, eFlag, rFlag, fFlag, a.path)

	case "diag":
		err := a.diagCmd.Parse(args[1:len(args)-1])
		if err != nil {
			// Control should not reach here. Parse typically does a 'os.Exit()' if something goes wrong.
			// If you do reach, contact author.
			return err
		}
		smFlag, _ := strconv.ParseBool(a.diagCmd.Lookup("sm").Value.String())
		admFlag, _ := strconv.ParseBool(a.diagCmd.Lookup("adm").Value.String())
		dFlag := a.diagCmd.Lookup("d").Value.String()

		return diagInvoker(smFlag, admFlag, dFlag, a.path)

	case "report":
		err := a.reportCmd.Parse(args[1:len(args)-1])
		if err != nil {
			// Control should not reach here. Parse typically does a 'os.Exit()' if something goes wrong.
			// If you do reach, contact author.
			return err
		}

		return reportInvoker(a.reportCmd.Lookup("d").Value.String(), a.path)
/*
	case "export":
		err := a.exportCmd.Parse(args[1:len(args)-1])
		if err != nil {
			// Control should not reach here. Parse typically does a 'os.Exit()' if something goes wrong.
			// If you do reach, contact author.
			return err
		}
		fFlag := a.exportCmd.Lookup("f").Value.String()
		dFlag := a.exportCmd.Lookup("d").Value.String()

		return exportInvoker(fFlag, dFlag, a.path)
*/
	default:
		return errors.New("INVALID ARGUMENT - \"" + args[0] + "\"")
	}
}	