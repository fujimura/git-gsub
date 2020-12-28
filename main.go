package main

import (
	"fmt"
	"io/ioutil"
	"log"
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

type Substitution struct {
	re *regexp.Regexp
	to string
}

type Substitutions = map[string]Substitution

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

func runSubstitionsAndRenames(substitutions Substitutions, rename bool, path string) {
	if path == "" {
		return
	}

	info, err := os.Stat(path)

	if err != nil {
		log.Fatal(err)
	}

	if info.IsDir() {
		return
	}

	content, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
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
}

func getMaxProcs() int {
	var maxProcs int

	mp := os.Getenv("GIT_GSUB_MAX_PROCS")
	if mp == "" {
		maxProcs = 100
	} else {
		i, err := strconv.Atoi(mp)
		if err != nil {
			log.Fatal(err)
		}
		maxProcs = i
	}

	return maxProcs
}

func run(_args []string) (int, error) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	var snake = flag.Bool("snake", false, "Substitute snake-cased expressions")
	var kebab = flag.Bool("kebab", false, "Substitute kebab-cased expressions")
	var camel = flag.Bool("camel", false, "Substitute camel-cased expressions")
	var rename = flag.Bool("rename", false, "Rename files with expression")
	var fgrep = flag.BoolP("fgrep", "F", false, "Interpret given pattern as a fixed string")
	var version = flag.Bool("version", false, "Show version")

	flag.CommandLine.Parse(_args)
	args := flag.Args()

	if *version {
		fmt.Println(Version)
		return 0, nil
	}

	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage git gsub [options] FROM TO [PATHS]\n")
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
		return 1, nil
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

	substitutions[rawFrom] = Substitution{regexp.MustCompile(rawFrom), to}

	if *snake {
		snakeFrom := strcase.ToSnake(rawFrom)
		substitutions[snakeFrom] = Substitution{regexp.MustCompile(snakeFrom), strcase.ToSnake(to)}
	}
	if *kebab {
		kebabFrom := strcase.ToKebab(rawFrom)
		substitutions[kebabFrom] = Substitution{regexp.MustCompile(kebabFrom), strcase.ToKebab(to)}
	}
	if *camel {
		camelFrom := strcase.ToCamel(rawFrom)
		substitutions[camelFrom] = Substitution{regexp.MustCompile(camelFrom), strcase.ToCamel(to)}
	}

	files, err := getAllFiles(targetPaths)
	if err != nil {
		return 1, err

	}

	c := make(chan bool, getMaxProcs())
	var wg sync.WaitGroup

	for _, path := range files {
		wg.Add(1)
		go func(path_ string) {
			c <- true
			runSubstitionsAndRenames(substitutions, *rename, path_)
			<-c
			wg.Done()
		}(path)
	}
	wg.Wait()

	return 0, nil
}

func main() {
	exitcode, err := run(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(exitcode)
}
