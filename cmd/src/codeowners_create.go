package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/khulnasoft/khulnasoft/lib/errors"

	"github.com/khulnasoft/src-cli/internal/api"
	"github.com/khulnasoft/src-cli/internal/cmderrors"
)

func init() {
	usage := `
Examples:

  Create a codeowners file for the repository "github.com/khulnasoft/khulnasoft":

    	$ src codeowners create -repo='github.com/sourcegraph/sourcegraph' -f CODEOWNERS

  Create a codeowners file for the repository "github.com/khulnasoft/khulnasoft" from stdin:

    	$ src codeowners create -repo='github.com/sourcegraph/sourcegraph' -f -
`

	flagSet := flag.NewFlagSet("create", flag.ExitOnError)
	usageFunc := func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of 'src codeowners %s':\n", flagSet.Name())
		flagSet.PrintDefaults()
		fmt.Println(usage)
	}
	var (
		repoFlag      = flagSet.String("repo", "", "The repository to attach the data to")
		fileFlag      = flagSet.String("file", "", "File path to read ownership information from (- for stdin)")
		fileShortFlag = flagSet.String("f", "", "File path to read ownership information from (- for stdin). Alias for -file")

		apiFlags = api.NewFlags(flagSet)
	)

	handler := func(args []string) error {
		if err := flagSet.Parse(args); err != nil {
			return err
		}

		if *repoFlag == "" {
			return errors.New("provide a repo name using -repo")
		}

		if *fileFlag == "" && *fileShortFlag == "" {
			return errors.New("provide a file using -file")
		}
		if *fileFlag != "" && *fileShortFlag != "" {
			return errors.New("have to provide either -file or -f")
		}
		if *fileShortFlag != "" {
			*fileFlag = *fileShortFlag
		}

		file, err := readFile(*fileFlag)
		if err != nil {
			return err
		}

		content, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		client := cfg.apiClient(apiFlags, flagSet.Output())

		query := `mutation CreateCodeownersFile(
	$repoName: String!,
	$content: String!
) {
	addCodeownersFile(input: {
		repoName: $repoName,
		fileContents: $content,
	}
	) {
		...CodeownersFileFields
	}
}
` + codeownersFragment

		var result struct {
			AddCodeownersFile CodeownersIngestedFile
		}
		if ok, err := client.NewRequest(query, map[string]interface{}{
			"repoName": *repoFlag,
			"content":  string(content),
		}).Do(context.Background(), &result); err != nil || !ok {
			var gqlErr api.GraphQlErrors
			if errors.As(err, &gqlErr) {
				for _, e := range gqlErr {
					if strings.Contains(e.Error(), "repo not found:") {
						return cmderrors.ExitCode(2, errors.Newf("repository %q not found", *repoFlag))
					}
					if strings.Contains(e.Error(), "codeowners file has already been ingested for this repository") {
						return cmderrors.ExitCode(2, errors.New("codeowners file has already been ingested for this repository"))
					}
				}
			}
			return err
		}

		return nil
	}

	// Register the command.
	codeownersCommands = append(codeownersCommands, &command{
		flagSet:   flagSet,
		handler:   handler,
		usageFunc: usageFunc,
	})
}
