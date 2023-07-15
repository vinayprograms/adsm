package test

import (
	"args"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShowHelp(t *testing.T) {

	parser := args.Args{}
	parser.InitArgs()

	harness := output_interceptor{}
	harness.Hook()

	parser.PrintHelpToStdout()

	out, _ := harness.ReadAndRelease()

	expected :=
		"\nUsage: adsm [OPTIONS] [PATH]\n" +
			"\n[PATH]: The path to a directory or a single smspec file\n" +
			"\nstat: List model components\n" +
			"  -e\tList in-scope entities only.\n" +
			"  -f\tList flows only.\n" +
			"  -r\tList roles only.\n" +
			"  -x\tList external entities only.\n" + 
			"\n" +
			"diag: Generate security model and ADM diagrams.\n" +
			"  -adm\n" +
			"    \tGenerate ADM decision graph only.\n" +
			"  -d string\n" +
			"    \tOutput directory for diagrams. (default \"./\")\n" +
			"  -sm\n" + 
			"    \tGenerate security model diagram only.\n" +
			"\n" +
			"report: Generate security report as markdown file.\n" +
			"  -d string\n" +
			"    	Output directory for generated report. (default \"./\")\n"/* +
			"\n" +
			"export: Export security model and report to other formats.\n" +
			"  -d string\n" +
			"    	Output directory for exported files. (default \"./\")\n" +
			"  -f string\n" +
			"    	Output format. Supported values - json,xml. (default \"json\")\n"*/

	assert.Equal(t, out, expected)
}

func TestUninitializedArgsStruct_ShowHelp(t *testing.T) {
	parser := args.Args{}
	// initArgs not called.

	harness := output_interceptor{}
	harness.Hook()
	parser.ParseArgs([]string{})
	out, _ := harness.ReadAndRelease()
	expected :=
		"\nUsage: adsm [OPTIONS] [PATH]\n" +
			"\n[PATH]: The path to a directory or a single smspec file\n" +
			"\nstat: List model components\n" +
			"  -e\tList in-scope entities only.\n" +
			"  -f\tList flows only.\n" +
			"  -r\tList roles only.\n" +
			"  -x\tList external entities only.\n" + 
			"\n" +
			"diag: Generate security model and ADM diagrams.\n" +
			"  -adm\n" +
			"    \tGenerate ADM decision graph only.\n" +
			"  -d string\n" +
			"    \tOutput directory for diagrams. (default \"./\")\n" +
			"  -sm\n" + 
			"    \tGenerate security model diagram only.\n" +
			"\n" +
			"report: Generate security report as markdown file.\n" +
			"  -d string\n" +
			"    	Output directory for generated report. (default \"./\")\n"/* +
			"\n" +
			"export: Export security model and report to other formats.\n" +
			"  -d string\n" +
			"    	Output directory for exported files. (default \"./\")\n" +
			"  -f string\n" +
			"    	Output format. Supported values - json,xml. (default \"json\")\n"*/

	assert.Equal(t, out, expected)
}

func TestParseArgsWithoutSubCommand(t *testing.T) {
	args := []string{"./some_dir_that_doesn't_exist"}
	err := sendToParseArgs(args)
	assert.Equal(t, "require atleast two parameters - 'sub-command' and 'path'", err.Error())
}

func TestParseArgs_MissingParams(t *testing.T) {
	testVectors := map[string][]string{ // The last item in args list is the expected value
		"StatsInvoke":		{"stat", "require atleast two parameters - 'sub-command' and 'path'"},
		"DiagInvoke":			{"graph", "require atleast two parameters - 'sub-command' and 'path'"},
		"ReportInvoke":		{"report", "require atleast two parameters - 'sub-command' and 'path'"},
		//"ExportInvoke":		{"export", "require atleast two parameters - 'sub-command' and 'path'"},
	}

	for name, args := range testVectors {
		t.Run(name, func(t *testing.T) {

			expected := args[len(args)-1]
			params := args[:len(args)-1]

			harness := output_interceptor{}
			harness.Hook()

			err := sendToParseArgs(params)
			out, _ := harness.ReadAndRelease()

			// only the last line contains the required output
			out = strings.Split(out, "\n")[1]
			fmt.Println(out)

			assert.Equal(t, expected, err.Error())
		})
	}
}

func TestParseArgsWithoutPathAndOnlyFlags(t *testing.T) {
	testVectors := map[string][]string{ // The last item in args list is the expected value
		"StatsInvoke":		{"stat", "-x", "-e", "error when verifying path - '-e'"},
		"DiagInvoke":			{"diag", "-sm", "-adm", "error when verifying path - '-adm'"},
		//"ReportInvoke":		{"report", "-d", "./", "error when verifying path - './'"}, // Cannot test this. It has only one flag.
		//"ExportInvoke":		{"export", "-f", "-d", "error when verifying path - '-d'"},
	}

	for name, args := range testVectors {
		t.Run(name, func(t *testing.T) {
			expected := args[len(args)-1]
			params := args[:len(args)-1]
			err := sendToParseArgs(params)
			assert.Equal(t, expected, err.Error())
		})
	}
}

func TestParseArgsWithPathAndWrongSubcommand(t *testing.T) {
	args := []string{"-z", "./examples/basic"}
	err := sendToParseArgs(args)
	assert.Equal(t, "INVALID ARGUMENT - \"-z\"", err.Error())
}

func TestParseArgsWithUnsupportedPath(t *testing.T) {
	testVectors := map[string][]string{ // The last item in args list is the expected value
		"Stats":	{"stat", "dummy://dummy.dummy/test.adm", "error when verifying path - 'dummy://dummy.dummy/test.adm'"},
		"Diag":		{"diag", "dummy://dummy.dummy/test.adm", "error when verifying path - 'dummy://dummy.dummy/test.adm'"},
		"Report":	{"report", "dummy://dummy.dummy/test.adm", "error when verifying path - 'dummy://dummy.dummy/test.adm'"},
		//"Export":	{"export", "dummy://dummy.dummy/test.adm", "error when verifying path - 'dummy://dummy.dummy/test.adm'"},
	}

	for name, args := range testVectors {
		t.Run(name, func(t *testing.T) {
			expected := args[len(args)-1]
			params := args[:len(args)-1]
			err := sendToParseArgs(params)
			assert.Equal(t, expected, err.Error())
		})
	}
}

func TestParseArgsWithValidPath(t *testing.T) {
	testVectors := map[string][]string{ // The last item in args list is the expected value
		"StatsExternal":	{"stat", "-x", "examples/simple.smspec"},
		"StatsEntity":		{"stat", "-e", "examples/simple.smspec"},
		"StatsRole":			{"stat", "-r", "examples/simple.smspec"},
		"StatsFlow":			{"stat", "-f", "examples/simple.smspec"},
		"StatsAllNonFlows":{"stat", "-x", "-e", "-r", "-f", "examples/simple.smspec"},
		"DiagSM":					{"diag", "-sm", "examples/simple.smspec"},
		"DiagSMWithPath":	{"diag", "-d", "./examples/sm", "-sm", "examples/simple.smspec"},
		"DiagADMWithPath":{"diag", "-d", "./examples/adm", "-adm", "examples/simple.smspec"},
		"Report":					{"report", "examples/simple.smspec"},
		"ReportWithPath":	{"report", "-d", "examples", "examples/simple_addb.smspec"},
	}

	for name, args := range testVectors {
		t.Run(name, func(t *testing.T) {
			err := sendToParseArgs(args)
			assert.Nil(t, err)
		})
	}
}