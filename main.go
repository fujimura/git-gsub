package main

import (
	"flag"
	"fmt"
	"github.com/iancoleman/strcase"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func SubstituteFile(path string, re *regexp.Regexp, to string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	result := re.ReplaceAllString(string(content), to)
	ioutil.WriteFile(string(path), []byte(result), os.ModePerm)
}

func FindTargetFiles(from string, paths []string, options ...string) []string {
	var args []string
	args = append(args, "grep", "-l")
	args = append(args, options...)
	args = append(args, from)
	args = append(args, paths...)
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	exitCode := cmd.ProcessState.ExitCode()
	if exitCode != 1 && err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(out), "\n")
	return lines
}

func GetAllFiles(paths []string) []string {
	var args []string
	args = append(args, "ls-files")
	args = append(args, paths...)
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	exitCode := cmd.ProcessState.ExitCode()
	if exitCode != 1 && err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(out), "\n")
	return lines
}

func RunSubstitutions(from string, to string, paths []string) {
	for _, path := range paths {
		if path == "" {
			continue
		}
		re, err := regexp.Compile(from)
		if err != nil {
			log.Fatal(err)
		}
		SubstituteFile(path, re, to)
	}
}

func RunRenames(from string, to string, path string) {
	re, err := regexp.Compile(from)
	if err != nil {
		log.Fatal(err)
	}

	newpath := re.ReplaceAllString(path, to)
	if newpath != path {
		os.MkdirAll(filepath.Dir(newpath), os.ModePerm)
		os.Rename(path, newpath)
	}
}

func main() {
	var snake = flag.Bool("snake", false, "Substitute snake-cased expressions")
	var kebab = flag.Bool("kebab", false, "Substitute kebab-cased expressions")
	var camel = flag.Bool("camel", false, "Substitute camel-cased expressions")
	var rename = flag.Bool("rename", false, "Rename files with expression")

	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage git gsub [options] FROM TO [PATHS]\n")
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	from := args[0]
	to := args[1]
	var targetPaths []string
	if len(args) > 2 {
		targetPaths = args[2:]
	}

	files := FindTargetFiles(from, targetPaths)
	RunSubstitutions(from, to, files)

	snakeFrom := strcase.ToSnake(from)
	snakeTo := strcase.ToSnake(to)
	kebabFrom := strcase.ToKebab(from)
	kebabTo := strcase.ToKebab(to)
	camelFrom := strcase.ToCamel(from)
	camelTo := strcase.ToCamel(to)

	if *snake {
		snakePaths := FindTargetFiles(snakeFrom, targetPaths)
		RunSubstitutions(snakeFrom, snakeTo, snakePaths)
	}
	if *kebab {
		kebabPaths := FindTargetFiles(kebabFrom, targetPaths)
		RunSubstitutions(kebabFrom, kebabTo, kebabPaths)
	}
	if *camel {
		camelPaths := FindTargetFiles(camelFrom, targetPaths)
		RunSubstitutions(camelFrom, camelTo, camelPaths)
	}
	if *rename {
		allFiles := GetAllFiles(targetPaths)

		for _, path := range allFiles {
			RunRenames(from, to, path)
			if *snake {
				RunRenames(snakeFrom, snakeTo, path)
			}
			if *kebab {
				RunRenames(kebabFrom, kebabTo, path)
			}
			if *camel {
				RunRenames(camelFrom, camelTo, path)
			}
		}
	}
}
