package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/xo/xo/cmd"
)

func main() {
	out := flag.String("out", "", "output directory")
	tplDir := flag.String("src", "", "directory containing template")
	tplName := flag.String("template", "", "template name")
	flag.Parse()
	remainingArgs := flag.Args()
	args := Args{
		TemplateDir:  *tplDir,
		TemplateName: *tplName,
		Out:          *out,
		Args:         remainingArgs,
	}
	if err := run(context.Background(), args); err != nil {
		fmt.Fprintln(os.Stderr, "Error running command:", err)
		os.Exit(1)
	}
}

type Args struct {
	// TemplateDir is the directory containing the template to generate from.
	TemplateDir string
	// TemplateName is the name of the template to generate from.
	TemplateName string
	// Out is the output directory.
	Out string
	// Args is the remaining arguments of the directory.
	Args []string
}

func run(ctx context.Context, args Args) error {
	if args.Out == "" {
		return fmt.Errorf("output directory must be specified")
	}
	if len(args.Args) != 1 {
		return fmt.Errorf(
			"expected exactly 1 arguments for input directory, got %d instead",
			len(args.Args),
		)
	}
	tSet, err := cmd.NewTemplateSet(ctx, args.TemplateDir, args.TemplateName)
	if err != nil {
		return fmt.Errorf("error reading template: %w", err)
	}
	files, err := os.ReadDir(args.Args[0])
	if err != nil {
		return fmt.Errorf("cannot fetch input directory: %w", err)
	}
	for _, file := range files {
		fileName := file.Name()
		if file.IsDir() {
			fmt.Println("Skipping directory", fileName)
			continue
		}
		f, err := os.Open(file.Name())
		if err != nil {
			return fmt.Errorf("error opening file %q: %w", fileName, err)
		}
		defer f.Close()
		decoder := yaml.NewDecoder(f)
		decoded := File{}
		if err := decoder.Decode(&decoded); err != nil {
			return fmt.Errorf("error reading file %q: %w", fileName, err)
		}
		if err := generate(ctx, tSet, args.Out, decoded); err != nil {
			return fmt.Errorf("error generating for file %q: %w", fileName, err)
		}
	}
	return nil
}

type multiError []error

func (m multiError) Error() string {
	errorDescriptions := make([]string, 0, len(m))
	for _, v := range m {
		errorDescriptions = append(errorDescriptions, strings.Split(v.Error(), "\n")...)
	}
	return strings.Join(errorDescriptions, "\n\t")
}
