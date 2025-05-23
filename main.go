package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/pechorka/gostdlib/pkg/clipboard"
	"github.com/pechorka/gostdlib/pkg/errs"
	"github.com/pechorka/gostdlib/pkg/uuid"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	cmds := []Command{
		newB64Command(),
		newUuidCommand(),
	}
	err := run(ctx, cmds)
	if err != nil {
		usage(err, cmds)
		os.Exit(2)
	}
}

func run(ctx context.Context, cmds []Command) error {
	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		return errs.Wrap(err, "failed to pass cmdArgs")
	}

	rest := flag.Args()
	if len(rest) == 0 {
		return errs.New("could't figure out command")
	}

	cmdName, cmdArgs := rest[0], rest[1:]

	for _, c := range cmds {
		if c.Name == cmdName || c.Short == cmdName {
			c.FS.Parse(cmdArgs)
			return c.Run(ctx)
		}
	}

	return errs.Errorf("%s is unknown command", cmdName)
}

func usage(err error, cmds []Command) {
	fmt.Fprintf(os.Stderr, "%s\nUsage:\n", err.Error())
	fmt.Fprintf(os.Stderr, "%s [global flags] <command> [flags]\n\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "Global flags:")
	flag.PrintDefaults()

	fmt.Fprintln(os.Stderr, "\nCommands:")
	for _, c := range cmds {
		fmt.Fprintf(os.Stderr, "  -%s (or -%s)\n", c.Name, c.Short)
	}
	fmt.Fprintf(os.Stderr, "\nRun '%s <command> -h' for details.", os.Args[0])
}

type Command struct {
	Name, Short string
	FS          *flag.FlagSet
	Run         func(ctx context.Context) error
}

func newB64Command() Command {
	const name = "base64"
	fs := newFlagSet(name)

	encode := boolAlias(fs, "e", "encode", true, "encode input")
	decode := boolAlias(fs, "d", "decode", false, "decode input")
	inputFile := stringAlias(fs, "in", "input", "", "input file")
	inputText := stringAlias(fs, "t", "text", "", "input text")
	outputPath := stringAlias(fs, "o", "output", "", "output file path. If empty, output will be printed to stdout")

	return Command{
		Name:  name,
		Short: "b64",
		FS:    fs,
		Run: func(ctx context.Context) error {
			input, err := readInput(*inputFile, *inputText)
			if err != nil {
				return err
			}

			var output []byte
			if *decode {
				output, err = base64.RawStdEncoding.AppendDecode(nil, input)
				if err != nil {
					return errs.Wrap(err, "failed to encode content")
				}
			} else if *encode {
				output = base64.RawStdEncoding.AppendEncode(nil, input)
			}

			err = writeOutput(*outputPath, output)
			if err != nil {
				return err
			}

			return nil
		},
	}
}

func newUuidCommand() Command {
	const name = "uuid"
	fs := newFlagSet(name)

	v4 := fs.Bool("v4", true, "generate v4 uuid")
	v7 := fs.Bool("v7", true, "generate v7 uuid")
	crypto := boolAlias(fs, "crypto", "c", false, "use cryptographic random generator. Slower and may fail")

	return Command{
		Name:  name,
		Short: "b64",
		FS:    fs,
		Run: func(ctx context.Context) error {
			var output string
			var err error
			switch {
			case *v4 && !*crypto:
				output = uuid.MustV4PseudoString()
			case *v4 && *crypto:
				output, err = uuid.NewV4CryptoString()
			case *v7 && !*crypto:
				output = uuid.MustV7PseudoString()
			case *v7 && *crypto:
				output, err = uuid.NewV4CryptoString()
			}
			if err != nil {
				return err
			}

			err = writeOutput("", []byte(output))
			if err != nil {
				return err
			}

			return nil
		},
	}
}

func readInput(filePath, text string) ([]byte, error) {
	if text != "" {
		// TODO:implement custom flag that will allow to provide byte input
		return []byte(text), nil
	}

	if filePath != "" {
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return nil, errs.Wrapf(err, "failed to read file %s", filePath)
		}

		return fileContent, nil
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		return nil, errs.Wrap(err, "failed to stat stdin")
	}
	if fi.Mode()&os.ModeCharDevice == 0 {
		stdinContent, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, errs.Wrap(err, "failed to read text from stdin")
		}
		return stdinContent, nil
	}

	clipboardContent, err := clipboard.Read()
	if err != nil {
		return nil, errs.Wrap(err, "failed to read clipboard content")
	}
	if len(clipboardContent) > 0 {
		return clipboardContent, nil
	}

	return nil, errs.New("no input provided")
}

func writeOutput(filePath string, data []byte) error {
	if filePath != "" {
		err := os.WriteFile(filePath, data, os.ModePerm)
		if err != nil {
			return errs.Wrapf(err, "failed to write data to file %s", filePath)
		}
	}

	_, err := io.Copy(os.Stdout, bytes.NewReader(data))
	if err != nil {
		return errs.Wrap(err, "failed to write output to stdout")
	}

	return nil
}

func newFlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "usage of %s: \n", name)
		fs.PrintDefaults()
	}

	return fs
}

func stringAlias(fs *flag.FlagSet, short, long string, value string, usage string) *string {
	var dst string
	fs.StringVar(&dst, short, value, usage)
	fs.StringVar(&dst, long, value, usage)

	return &dst
}

func boolAlias(fs *flag.FlagSet, short, long string, value bool, usage string) *bool {
	var dst bool
	fs.BoolVar(&dst, short, value, usage)
	fs.BoolVar(&dst, long, value, usage)

	return &dst
}
