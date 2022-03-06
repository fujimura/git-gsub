package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/iancoleman/strcase"
	flag "github.com/spf13/pflag"
)

const Version string = "v0.0.15"

type CLI struct {
	outStream, errStream io.Writer
}

type Substitution struct {
	re *regexp.Regexp
	to string
}

func getAllFiles(paths []string) ([]string, error) {
	var args []string
	args = append(args, "ls-files")
	args = append(args, "-z")
	args = append(args, paths...)
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	exitCode := cmd.ProcessState.ExitCode()
	if exitCode != 1 && err != nil {
		return []string{}, err
	}
	lines := strings.Split(string(out), "\x00")
	return lines, nil
}

func runSubstitionsAndRenames(substitutions map[string]Substitution, rename bool, path string) error {
	if path == "" {
		return nil
	}

	info, err := os.Stat(path)

	if err != nil {
		return err
	}

	if info.IsDir() {
		return err
	}

	content, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	replaced := false

	for _, sub := range substitutions {
		if sub.re.Match(content) {
			replaced = true
			content = sub.re.ReplaceAll(content, []byte(sub.to))
		}
	}

	if replaced {
		ioutil.WriteFile(path, content, os.ModePerm)
	}

	if rename {
		for _, sub := range substitutions {
			// TODO PERF
			newpath := sub.re.ReplaceAllString(path, sub.to)
			if newpath != path {
				os.MkdirAll(filepath.Dir(newpath), os.ModePerm)
				os.Rename(path, newpath)
			}
		}
	}
	return nil
}

func getMaxProcs() (int, error) {
	var maxProcs int

	mp := os.Getenv("GIT_GSUB_MAX_PROCS")
	if mp == "" {
		maxProcs = 100
	} else {
		i, err := strconv.Atoi(mp)
		if err != nil {
			return 0, err
		}
		maxProcs = i
	}

	return maxProcs, nil
}

func ToRubyDirectory(str string) string {
	result := strcase.ToSnake(str)
	return strings.Replace(result, "::", "/", -1)
}

func ToRubyModule(str string) string {
	result := strcase.ToCamel(str)
	return strings.Replace(result, "/", "::", -1)
}

func addSub(substitutions *map[string]Substitution, from string, to string, conv func(string) string) {
	f := conv(from)
	t := conv(to)
	(*substitutions)[f] = Substitution{regexp.MustCompile(f), t}
}

func (cli *CLI) Run(_args []string) int {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	var snake = flag.Bool("snake", false, "Substitute snake-cased expressions")
	var kebab = flag.Bool("kebab", false, "Substitute kebab-cased expressions")
	var camel = flag.Bool("camel", false, "Substitute camel-cased expressions")
	var ruby = flag.Bool("ruby", false, "Substitute Ruby module and directory expressions")
	var all = flag.BoolP("all", "a", false, "Substitute snake, kebab, camel and Ruby expressions")
	var rename = flag.BoolP("rename", "r", false, "Rename files with expression")
	var fgrep = flag.BoolP("fgrep", "F", false, "Interpret given pattern as a fixed string")
	var version = flag.BoolP("version", "v", false, "Show version")

	flag.CommandLine.Parse(_args)
	flag.CommandLine.SetOutput(cli.errStream)
	args := flag.Args()

	if *version {
		fmt.Fprintf(cli.outStream, Version)
		return 0
	}

	if len(args) < 2 {
		fmt.Fprintf(cli.errStream, "Usage git gsub [options] FROM TO [PATHS]\n")
		fmt.Fprintf(cli.errStream, "\nOptions:\n")
		fmt.Fprintf(cli.errStream, flag.CommandLine.FlagUsages())
		return 1
	}

	rawFrom := args[0]
	if *fgrep {
		rawFrom = regexp.QuoteMeta(rawFrom)
	}
	to := args[1]

	var targetPaths []string
	if len(args) > 2 {
		targetPaths = args[2:]
	}

	substitutions := map[string]Substitution{}

	addSub(&substitutions, rawFrom, to, func(x string) string { return x })

	if *snake || *all {
		addSub(&substitutions, rawFrom, to, strcase.ToSnake)
	}
	if *kebab || *all {
		addSub(&substitutions, rawFrom, to, strcase.ToKebab)
	}
	if *camel || *all {
		addSub(&substitutions, rawFrom, to, strcase.ToCamel)
	}

	if *ruby || *all {
		addSub(&substitutions, rawFrom, to, ToRubyDirectory)
		addSub(&substitutions, rawFrom, to, ToRubyModule)
	}

	files, err := getAllFiles(targetPaths)
	if err != nil {
		fmt.Fprint(cli.errStream, err)
		return 1

	}

	maxProcs, err := getMaxProcs()

	if err != nil {
		fmt.Fprint(cli.errStream, err)
		return 1
	}

	cn := make(chan bool, maxProcs)
	errCn := make(chan error, maxProcs)
	var wg sync.WaitGroup

	for _, path := range files {
		wg.Add(1)
		go func(path_ string) {
			cn <- true
			err := runSubstitionsAndRenames(substitutions, *rename, path_)
			if err == nil {
				<-cn
			} else {
				errCn <- err
				close(cn)
			}
			wg.Done()
		}(path)
	}

	wg.Wait()
	close(errCn)

	failed := false

	for err := range errCn {
		fmt.Fprint(cli.errStream, err)
		failed = true
	}

	if failed {
		return 1
	} else {
		return 0
	}
}

func main() {
	cli := &CLI{outStream: os.Stdout, errStream: os.Stderr}
	exitcode := cli.Run(os.Args[1:])
	os.Exit(exitcode)
}
