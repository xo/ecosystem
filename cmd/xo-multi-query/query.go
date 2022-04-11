package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/xo/xo/cmd"
	"github.com/xo/xo/templates"

	// Drivers included in xo/xo/main.go.
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/sijms/go-ora/v2"
)

// File is the input file containing multiple queries to generate for a
// database.
type File struct {
	// DB is the URL to use to connect to the database for query introspection.
	DB string `yaml:"db"`
	// Queries is a list of query to generate for.
	Queries map[string]Query `yaml:"queries"`
}

// Query is the query to generate for.
type Query struct {
	// Query is the actual SQL query.
	Query string `yaml:"query"`
	// Args is a map of arguments to apply when generating for the query.
	Args map[string]string `yaml:"args"`
}

func generate(ctx context.Context, ts *templates.Set, out string, f File) error {
	for queryName, query := range f.Queries {
		xoArgs := cmd.NewArgs(ts.Target(), ts.Targets()...)
		queryCmd, err := cmd.QueryCommand(ctx, ts, xoArgs)
		if err != nil {
			return err
		}
		args := make([]string, 0, len(query.Args)+7)
		for k, v := range query.Args {
			args = append(args, "--"+k+"="+v)
		}
		args = append(args, "-o", out, "-Q", query.Query, "-F", queryName, f.DB)
		queryCmd.SetUsageFunc(func(*cobra.Command) error {
			// On invalid arguments, dump the arguments used instead of the
			// full help.
			for k, v := range args {
				args[k] = strconv.Quote(v)
			}
			fmt.Printf("Arguments used: %v\n", args)
			return nil
		})
		queryCmd.SetArgs(args)
		if err := queryCmd.Execute(); err != nil {
			return fmt.Errorf("error generating query %q: %w", queryName, err)
		}
	}
	return nil
}
