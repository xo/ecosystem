package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/xo/ecosystem/proto"
	"github.com/xo/ecosystem/types"
	"github.com/xo/xo/cmd"
	"github.com/xo/xo/templates"
	xo "github.com/xo/xo/types"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/compiler/protogen"
)

var knownTargets = map[string]struct{}{
	"postgres":  {},
	"mysql":     {},
	"sqlite3":   {},
	"sqlserver": {},
}

func main() {
	plugin := &xoPlugin{
		protobufNames: make(map[string]typeEntry),
	}
	protogen.Options{}.Run(plugin.Run)
}

type xoPlugin struct {
	targets       []string
	templateDir   string
	templateName  string
	schema        types.Schema
	skipPrefixes  []string
	plugin        *protogen.Plugin
	queryFile     string
	protobufNames map[string]typeEntry

	cpuProfileOutput string
}

func (x *xoPlugin) Run(p *protogen.Plugin) error {
	x.plugin = p
	// Parse parameters for template directory, name and prefixes to skip.
	_, err := x.ParseParams(p.Request.GetParameter(), "")
	if err != nil {
		return err
	}

	if x.cpuProfileOutput != "" {
		f, err := os.Create(x.cpuProfileOutput)
		if err != nil {
			return fmt.Errorf("error creating CPU profile file: %w", err)
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			return fmt.Errorf("error starting CPU profile: %w", err)
		}
		defer pprof.StopCPUProfile()
	}

	files := p.Request.FileToGenerate
	converter := proto.Converter{
		SkipPrefixes: x.skipPrefixes,
		PackageNames: make(map[string]string, len(p.FilesByPath)),
	}
	for path, file := range p.FilesByPath {
		converter.PackageNames[path] = string(file.GoPackageName)
	}
	// Sort p.Files to make the generation deterministic.
	sort.Slice(p.Files, func(i int, j int) bool {
		return p.Files[i].Proto.GetName() < p.Files[j].Proto.GetName()
	})
	// process proto packages
	for _, f := range p.Files {
		if !slices.Contains(files, f.Proto.GetName()) {
			continue
		}
		for _, msg := range f.Messages {
			pkgName := string(f.GoPackageName)
			tables, err := converter.ConvertMessage(pkgName, msg)
			if err != nil {
				return err
			}
			if len(tables) == 0 {
				continue
			}
			x.schema.Tables = append(x.schema.Tables, tables...)
			// Add to protobuf name map.
			x.addMessageType(converter, msg, tables[0])
		}
		for _, enum := range f.Enums {
			goEnum, err := converter.ConvertEnum(enum)
			if err != nil {
				return err
			}
			x.schema.Enums = append(x.schema.Enums, goEnum)
			// Add to protobuf name map.
			x.addEnumType(enum, goEnum)
		}
	}
	pkg := pathPackage(files[0])
	gf := p.NewGeneratedFile(pkg+".yaml", "")
	if err := yaml.NewEncoder(gf).Encode(x.schema); err != nil {
		return err
	}

	// Reuse xo schema command for parsing but inject our own schema set.
	ctx := context.Background()

	// Generate for all targets.
	var errs multiError
	for _, target := range x.targets {
		tSet, err := cmd.NewTemplateSet(ctx, x.templateDir, x.templateName)
		if err != nil {
			return err
		}

		if _, ok := knownTargets[target]; !ok {
			return fmt.Errorf("unknown sql output target %q", target)
		}
		xoArgs := cmd.NewArgs(tSet.Target(), tSet.Targets()...)
		cmd, err := x.Generate(ctx, tSet, xoArgs, target)
		if err != nil {
			return err
		}

		// Reparse to properly replace {{ .DB }}.
		flags, err := x.ParseParams(p.Request.GetParameter(), target)
		if err != nil {
			return err
		}
		cmd.SetArgs(flags)
		if err := cmd.Execute(); err != nil {
			errs = append(errs, fmt.Errorf("error while executing xo(%q): %w", target, err))
		}
	}
	if errs != nil {
		return errs
	}
	return nil
}

func (x *xoPlugin) ParseParams(params string, target string) ([]string, error) {
	split := strings.Split(params, ",")
	var flags []string
	var emitOk, pbOk bool
	for _, param := range split {
		key, val, _ := strings.Cut(param, "=")
		value, err := dbTpl(val, target)
		if err != nil {
			return nil, fmt.Errorf("could not parse %s value as template: %w", param, err)
		}
		switch key {
		case "cpu_profile":
			x.cpuProfileOutput = value
		case "skip_prefix":
			x.skipPrefixes = strings.Split(value, " ")
		case "src":
			x.templateDir = value
		case "template":
			x.templateName = value
		case "emit":
			flags = append(flags, "--out", value)
			emitOk = true
		case "targets":
			x.targets = strings.Split(value, " ")
		case "pb-names":
			pbOk, err = strconv.ParseBool(value)
			if err != nil {
				return nil, fmt.Errorf("unknown value for pb-names")
			}
		case "query":
			x.queryFile = value
		case "schema":
			x.schema.Name = value
			fallthrough
		default:
			if value != "" {
				flags = append(flags, "--"+key+"="+value)
			} else {
				flags = append(flags, "--"+key)
			}
		}
	}

	var missing []string
	if !emitOk {
		missing = append(missing, "emit")
	}
	if x.templateName == "" && x.templateDir == "" {
		missing = append(missing, "(src or template)")
	}
	if x.targets == nil {
		missing = append(missing, "targets")
	}
	if missing != nil {
		return nil, fmt.Errorf("missing required parameters: %v", missing)
	}

	if pbOk {
		// Marshal map[string]typeEntry and pass as parameter.
		val, _ := json.Marshal(x.protobufNames)
		tplName := x.templateName
		if tplName == "" {
			tplName = filepath.Base(x.templateDir)
		}
		flags = append(flags, "--"+tplName+"-protobuf", string(val))
	}
	return flags, nil
}

func (x *xoPlugin) Generate(ctx context.Context, ts *templates.Set, xoArgs *cmd.Args, target string) (*cobra.Command, error) {
	schemaCmd, err := cmd.SchemaCommand(ctx, ts, xoArgs)
	if err != nil {
		return nil, err
	}
	schemaCmd.RunE = func(c *cobra.Command, args []string) error {
		// Create output directory.
		err := os.MkdirAll(xoArgs.OutParams.Out, 0o755)
		if err != nil {
			return fmt.Errorf(
				"error creating output directory %q: %w",
				xoArgs.OutParams.Out, err,
			)
		}
		// Generate context as if we read from a DB.
		ctx = cmd.BuildContext(ctx, xoArgs)
		ctx = context.WithValue(ctx, xo.DriverKey, target)
		ctx = context.WithValue(ctx, xo.SchemaKey, xoArgs.LoaderParams.Schema)
		// Create a new SQL schema here based on db type.
		dbSchema := types.ToSQL(x.schema, target)
		set := &xo.Set{
			Schemas: []xo.Schema{dbSchema},
		}
		err = cmd.Generate(ctx, "schema", ts, set, xoArgs)
		if err != nil {
			return multiError(ts.Errors())
		}
		return nil
	}
	schemaCmd.Args = cobra.ExactArgs(0)
	return schemaCmd, nil
}

func dbTpl(tpl string, db string) (string, error) {
	if !strings.Contains(tpl, "{{") {
		return tpl, nil
	}
	t, err := template.New("db").Parse(tpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, map[string]string{
		"DB": db,
	})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// expected file format: xx/xx/xx/{packagename}/{filename}.proto
func pathPackage(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-2]
}

type multiError []error

func (m multiError) Error() string {
	errorDescriptions := make([]string, 0, len(m))
	for _, v := range m {
		errorDescriptions = append(errorDescriptions, strings.Split(v.Error(), "\n")...)
	}
	return strings.Join(errorDescriptions, "\n\t")
}
