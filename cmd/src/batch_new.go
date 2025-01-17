package main

import (
	"context"
	"flag"
	"fmt"
	cliLog "log"

	"github.com/khulnasoft/src-cli/internal/api"
	"github.com/khulnasoft/src-cli/internal/batches/service"
	"github.com/khulnasoft/src-cli/internal/cmderrors"
)

func init() {
	usage := `
'src batch new' creates a new batch spec YAML, prefilled with all required
fields.

Usage:

    src batch new [-f FILE]

Examples:


    $ src batch new -f batch.spec.yaml

`

	flagSet := flag.NewFlagSet("new", flag.ExitOnError)
	apiFlags := api.NewFlags(flagSet)

	var (
		fileFlag   = flagSet.String("f", "batch.yaml", "The name of the batch spec file to create.")
		skipErrors bool
	)
	flagSet.BoolVar(
		&skipErrors, "skip-errors", false,
		"If true, errors encountered won't stop the program, but only log them.",
	)

	handler := func(args []string) error {
		ctx := context.Background()

		if err := flagSet.Parse(args); err != nil {
			return err
		}

		if len(flagSet.Args()) != 0 {
			return cmderrors.Usage("additional arguments not allowed")
		}

		svc := service.New(&service.Opts{
			Client: cfg.apiClient(apiFlags, flagSet.Output()),
		})

		_, ffs, err := svc.DetermineLicenseAndFeatureFlags(ctx, skipErrors)
		if err != nil {
			return err
		}

		if err := validateSourcegraphVersionConstraint(ffs); err != nil {
			if !skipErrors {
				return err
			} else {
				cliLog.Printf("WARNING: %s", err)
			}
		}

		if err := svc.GenerateExampleSpec(ctx, *fileFlag); err != nil {
			return err
		}

		fmt.Printf("%s created.\n", *fileFlag)
		return nil
	}

	batchCommands = append(batchCommands, &command{
		flagSet: flagSet,
		aliases: []string{},
		handler: handler,
		usageFunc: func() {
			fmt.Fprintf(flag.CommandLine.Output(), "Usage of 'src batch %s':\n", flagSet.Name())
			flagSet.PrintDefaults()
			fmt.Println(usage)
		},
	})
}
