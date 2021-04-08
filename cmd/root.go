package cmd

import (
	"fmt"

	"github.com/10gen/realm-cli/internal/cli"
	"github.com/10gen/realm-cli/internal/commands"

	"github.com/spf13/cobra"
	"honnef.co/go/tools/version"

	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra/doc"
)

type Option struct {
	Name         string
	Shorthand    string
	DefaultValue string
	Usage        string
}
type Command struct {
	CommandRef       string
	CommandName      string
	Synopsis         string
	Description      string
	Usage            string
	Options          []Option
	InheritedOptions []Option
	SeeAlso          []string
}

func getFilenames() []string {
	var files []string

	root := "./yaml"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

// Get a command from yaml
func getCommand(command_name string) (*Command, error) {
	buf, err := ioutil.ReadFile(fmt.Sprintf("./yaml/%q.yaml", command_name))
	if err != nil {
		return nil, err
	}

	c := &Command{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q.yaml: %v", command_name, err)
	}

	return c, nil
}

// Does template stuff (automatically)
func AutoTemplateStuff() {
	tpl, parseErr := template.New("command-template.rst").Funcs(template.FuncMap{
		"toChars": func(s string) []rune {
			return []rune(s)
		},
	}).ParseFiles("command-template.rst")
	if parseErr != nil {
		log.Fatalln(parseErr)
	}

	var filenames = getFilenames()
	for _, filename := range filenames {
		var command_name string = strings.TrimSuffix(filename, filepath.Ext(filename))
		command, err := getCommand(command_name)
		if err != nil {
			log.Fatal(err)
		}

		outfile, outfileErr := os.Create(fmt.Sprintf("./generated-docs/%q.rst", command.CommandRef))
		if outfileErr != nil {
			log.Println("outfile err: ", err)
		}

		tplErr := tpl.Execute(outfile, command)
		if tplErr != nil {
			panic(tplErr)
		}
	}
}

// Does template stuff (manually)
func TemplateStuff() {
	tpl, parseErr := template.New("command-template.rst").Funcs(template.FuncMap{
		"toChars": func(s string) []rune {
			return []rune(s)
		},
	}).ParseFiles("command-template.rst")
	if parseErr != nil {
		log.Fatalln(parseErr)
	}
	command := Command{
		CommandRef:       "realm-cli_pull",
		CommandName:      "realm-cli pull",
		Synopsis:         "Pull the latest version of your Realm app into your local directory",
		Description:      "Pull the latest version of your Realm app into your local directory\n\nUpdates a remote Realm app with your local directory by pulling changes from the former into the latter. Input a Realm app that you would like to have changes pushed from. If applicable, hosting and/or dependencies associated with your Realm app will be exported as well.",
		Usage:            "realm-cli pull [flags]",
		Options:          make([]Option, 0),
		InheritedOptions: make([]Option, 0),
		SeeAlso:          []string{"realm-cli - CLI tool to manage your MongoDB Realm application"},
	}
	command.Options = append(command.Options, Option{
		Name:         "include-hosting",
		Shorthand:    "s",
		DefaultValue: "false",
		Usage:        "include to push Realm app hosting changes as well",
	})
	command.Options = append(command.Options, Option{
		Name:         "app-version",
		DefaultValue: "0",
		Usage:        "specify the app config version to pull changes down as",
	})
	command.Options = append(command.Options, Option{
		Name:         "dry-run",
		Shorthand:    "x",
		DefaultValue: "false",
		Usage:        "include to run without writing any changes to the file system",
	})
	command.Options = append(command.Options, Option{
		Name:         "help",
		Shorthand:    "h",
		DefaultValue: "false",
		Usage:        "help for pull",
	})
	command.Options = append(command.Options, Option{
		Name:         "include-dependencies",
		Shorthand:    "d",
		DefaultValue: "false",
		Usage:        "include to to push Realm app dependencies changes as well",
	})
	command.Options = append(command.Options, Option{
		Name:         "include-hosting",
		Shorthand:    "s",
		DefaultValue: "false",
		Usage:        "include to push Realm app hosting changes as well",
	})
	command.Options = append(command.Options, Option{
		Name:  "local",
		Usage: "specify the local path to export a Realm app to",
	})
	command.Options = append(command.Options, Option{
		Name:  "project",
		Usage: "the MongoDB cloud project id",
	})
	command.Options = append(command.Options, Option{
		Name:  "remote",
		Usage: "specify the remote app to pull changes down from",
	})

	command.InheritedOptions = append(command.InheritedOptions, Option{
		Name:  "atlas-url",
		Usage: "specify the base Atlas server URL",
	})
	command.InheritedOptions = append(command.InheritedOptions, Option{
		Name:         "disable-colors",
		DefaultValue: "false",
		Usage:        "disable output styling",
	})
	command.InheritedOptions = append(command.InheritedOptions, Option{
		Name:      "output-format",
		Shorthand: "f",
		Usage:     "set the output format, available InheritedOptions: [json]",
	})
	command.InheritedOptions = append(command.InheritedOptions, Option{
		Name:      "output-target",
		Shorthand: "o",
		Usage:     "write output to the specified filepath",
	})
	command.InheritedOptions = append(command.InheritedOptions, Option{
		Name:         "profile",
		Shorthand:    "i",
		DefaultValue: "default",
		Usage:        "this is the --profile, -p usage",
	})
	command.InheritedOptions = append(command.InheritedOptions, Option{
		Name:  "realm-url",
		Usage: "specify the base Realm server URL",
	})
	command.InheritedOptions = append(command.InheritedOptions, Option{
		Name:      "telemetry",
		Shorthand: "m",
		Usage:     "enable or disable telemetry (this setting is remembered), available InheritedOptions: [\"off\", \"on\"]",
	})
	command.InheritedOptions = append(command.InheritedOptions, Option{
		Name:         "yes",
		Shorthand:    "y",
		DefaultValue: "false",
		Usage:        "set to automatically proceed through command confirmations",
	})

	tplErr := tpl.Execute(os.Stdout, command)
	if tplErr != nil {
		panic(tplErr)
	}
}

// Generates reference documentation
func GenerateDocs(cmd *cobra.Command) {
	const fmTemplate = `
:date: %s
:title: "%s"
:slug: %s
:url: %s
`
	docsFilePrepender := func(filename string) string {
		now := time.Now().Format(time.RFC3339)
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		url := "/commands/" + strings.ToLower(base) + "/"
		return fmt.Sprintf(fmTemplate, now, strings.Replace(base, "_", " ", -1), base, url)
	}
	docsLinkHandler := func(name, ref string) string {
		return fmt.Sprintf(":ref:`%s <%s>`", name, ref)
	}
	docsErr := doc.GenReSTTreeCustom(cmd, "./docs", docsFilePrepender, docsLinkHandler)
	if docsErr != nil {
		log.Fatal(docsErr)
	}
}

// Run runs the CLI
func Run() {
	// print commands in help/usage text in the order they are declared
	cobra.EnableCommandSorting = false

	cmd := &cobra.Command{
		Version:       cli.Version,
		Use:           cli.Name,
		Short:         "CLI tool to manage your MongoDB Realm application",
		Long:          fmt.Sprintf("Use %s command help for information on a specific command", cli.Name),
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	factory := cli.NewCommandFactory()
	cobra.OnInitialize(factory.Setup)
	defer factory.Close()

	cmd.Flags().SortFlags = false // ensures CLI help text displays global flags unsorted
	factory.SetGlobalFlags(cmd.PersistentFlags())

	cmd.AddCommand(factory.Build(commands.Login))
	cmd.AddCommand(factory.Build(commands.Logout))
	cmd.AddCommand(factory.Build(commands.Push))
	cmd.AddCommand(factory.Build(commands.Pull))
	cmd.AddCommand(factory.Build(commands.App))
	cmd.AddCommand(factory.Build(commands.Secrets))
	cmd.AddCommand(factory.Build(commands.User))
	cmd.AddCommand(factory.Build(commands.Whoami))
	cmd.AddCommand(factory.Build(commands.Function))

	// GenerateDocs(cmd)
	doc.GenYamlTree(cmd, "./yaml")
	TemplateStuff()
	// AutoTemplateStuff()

	factory.Run(cmd)
}
