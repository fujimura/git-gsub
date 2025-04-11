package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func RunInTmpDir(run func()) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		log.Fatal(err)
	}

	current, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	err = os.Chdir(dir)
	if err != nil {
		panic(err)
	}

	run()

	defer func() {
		os.RemoveAll(dir) // clean up
		os.Chdir(current)
	}()

}

func RunInTmpRepo(run func()) {
	commands := []string{
		"git init",
		"git config --local user.email \"you@example.com\"",
		"git config --local user.name \"Your Name\"",
	}
	RunInTmpDir(func() {
		for _, command := range commands {
			err := exec.Command("sh", "-c", command).Run()
			if err != nil {
				panic(err)
			}
		}

		run()
	})
}

func CommitFile(name string, content string) {
	err := os.MkdirAll(filepath.Dir(name), os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(name, []byte(content), os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = exec.Command("git", "add", ".").Run()
	if err != nil {
		panic(err)
	}

	err = exec.Command("git", "commit", "-m", fmt.Sprintf("\"Add\" %s", name)).Run()
	if err != nil {
		panic(err)
	}
}

func GitGsubPath() string {
	_, filename, _, _ := runtime.Caller(0)

	return filepath.Clean(fmt.Sprintf("%s/../bin/git-gsub", filename))
}

func RunGitGsub(args ...string) ([]byte, error) {
	var out []byte
	var err error

	_, e2e := os.LookupEnv("E2E")

	if e2e {
		out, err = exec.Command(GitGsubPath(), args...).Output()
	} else {
		outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)

		cli := &CLI{outStream: outStream, errStream: errStream}
		exitcode := cli.Run(args)
		if exitcode != 0 {
			err = errors.New(errStream.String())
		}
		out = outStream.Bytes()
	}
	return out, err
}

func TestVersion(t *testing.T) {
	out, err := RunGitGsub("--version")

	if err != nil {
		t.Errorf("Command failed: %s", err)
	}

	if string(out) != "v0.1.2\n" {
		t.Errorf("Failed: %s", string(out))
	}
}

func TestSimpleSubstitution(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README.md", "Git Subversion Bzr")
		_, err := RunGitGsub("Bzr", "Mercurial")

		if err != nil {
			t.Errorf("Command failed: %s", err)
		}

		dat, _ := os.ReadFile("./README.md")
		if string(dat) != "Git Subversion Mercurial" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestSimpleSubstitutionManyFiles(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README_1.md", "Git Subversion Bzr")
		CommitFile("README_2.md", "Git Subversion Bzr")
		CommitFile("README_3.md", "Git Subversion Bzr")
		CommitFile("README_4.md", "Git Subversion Bzr")
		CommitFile("README_5.md", "Git Subversion Bzr")
		CommitFile("README_6.md", "Git Subversion Bzr")
		_, err := RunGitGsub("Bzr", "Mercurial")

		if err != nil {
			t.Errorf("Command failed: %s", err)
		}

		dat, _ := os.ReadFile("./README_1.md")
		if string(dat) != "Git Subversion Mercurial" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestSubstitutionWithPath(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README.md", "Git Subversion Bzr")
		CommitFile("foo/git", "Git Subversion Bzr")
		CommitFile("bar/git", "Git Subversion Bzr")

		_, err := RunGitGsub("Git", "Svn", "foo")

		if err != nil {
			t.Errorf("Command failed: %s", err)
		}

		dat, _ := os.ReadFile("./README.md")
		if string(dat) != "Git Subversion Bzr" {
			t.Errorf("Failed: %s", string(dat))
		}

		dat, _ = os.ReadFile("./foo/git")
		if string(dat) != "Svn Subversion Bzr" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestSubstitutionWithCaseConversion(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README.md", "GitGsub git_gsub git-gsub GIT_GSUB")
		_, err := RunGitGsub("--camel", "--kebab", "--snake", "--screaming-snake", "git-gsub", "svn-gsub")

		if err != nil {
			t.Errorf("Command failed: %s", err)
		}

		dat, _ := os.ReadFile("./README.md")
		if string(dat) != "SvnGsub svn_gsub svn-gsub SVN_GSUB" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestSubstitutionWithUpperAndLowerCamelCaseConversion(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README.md", "GitGsub gitGsub")
		_, err := RunGitGsub("--upper-camel", "--lower-camel", "git-gsub", "svn-gsub")

		if err != nil {
			t.Errorf("Command failed: %s", err)
		}

		dat, _ := os.ReadFile("./README.md")
		if string(dat) != "SvnGsub svnGsub" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestSubstitutionOfAllUnderscoredPhraseWithCaseConversion(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README.md", "activerecord")
		_, err := RunGitGsub("activerecord", "inactiverecord", "--kebab", "--snake", "--camel")

		if err != nil {
			t.Errorf("Command failed: %s", err)
		}

		dat, _ := os.ReadFile("./README.md")
		if string(dat) != "inactiverecord" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestOptionsCanBePutAfterArguments(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README.md", "GitGsub git_gsub git-gsub")
		_, err := RunGitGsub("git-gsub", "svn-gsub", "--camel", "--kebab", "--snake")

		if err != nil {
			t.Errorf("Command failed: %s", err)
		}

		dat, _ := os.ReadFile("./README.md")
		if string(dat) != "SvnGsub svn_gsub svn-gsub" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestSubstitutionWithFixedStringOption(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("hello.rb", "puts('hello')")

		_, err := RunGitGsub("--fgrep", "(", " ")
		if err != nil {
			t.Errorf("Command failed: %s", err)
		}

		_, err = RunGitGsub("-F", ")", "")
		if err != nil {
			t.Errorf("Command failed: %s", err)
		}

		dat, _ := os.ReadFile("./hello.rb")
		if string(dat) != "puts 'hello'" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestEscape(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README.md", `<h1 class="foo">`)
		_, err := RunGitGsub(`<h1 class="foo">`, `<h1 class="bar">`)

		if err != nil {
			t.Errorf("Command failed: %s", err)
		}

		dat, _ := os.ReadFile("./README.md")
		if string(dat) != `<h1 class="bar">` {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestAtMark(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README.md", "foo@example.com")
		_, err := RunGitGsub("foo@example.com", "bar@example.com")
		if err != nil {
			t.Errorf("Command failed: %s", err)
		}

		dat, _ := os.ReadFile("./README.md")
		if string(dat) != "bar@example.com" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestConsequesingAtMark(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README.md", "foo@example.com")
		_, err := RunGitGsub("foo@example.com", "bar@@example.com")
		if err != nil {
			t.Errorf("Command failed: %s", err)
		}

		dat, _ := os.ReadFile("./README.md")
		if string(dat) != "bar@@example.com" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestDoubleQuoteToSingleQuote(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README.md", `hello this is "git"`)
		_, err := RunGitGsub(`"git"`, `'svn'`)
		if err != nil {
			t.Errorf("Command failed: %s", err)
		}

		dat, _ := os.ReadFile("./README.md")
		if string(dat) != "hello this is 'svn'" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestSingleQuoteToDoubleQuote(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README.md", `hello this is 'git'`)
		_, err := RunGitGsub(`'git'`, `"svn"`)
		if err != nil {
			t.Errorf("Command failed: %s", err)
		}

		dat, _ := os.ReadFile("./README.md")
		if string(dat) != `hello this is "svn"` {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestBracket(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README.md", `{git{svn}}`)
		_, err := RunGitGsub("{git{svn}}", "{hg{svn}}")
		if err != nil {
			t.Errorf("Command failed: %s", err)
		}

		dat, _ := os.ReadFile("./README.md")
		if string(dat) != "{hg{svn}}" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestSubmatch(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README.md", "git-foo-1 git-bar-22 git-baz-3")
		_, err := RunGitGsub(`git-([a-z]+)-([\d]{1,2})`, `$2-$1`)
		if err != nil {
			t.Errorf("Command failed: %s", err)
		}

		dat, _ := os.ReadFile("./README.md")
		if string(dat) != "1-foo 22-bar 3-baz" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestSubstituteToEmptyString(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README.md", "Git Svn Hg")
		_, err := RunGitGsub("Svn ", "")
		if err != nil {
			t.Errorf("Command failed: %s", err)
		}

		dat, _ := os.ReadFile("./README.md")
		if string(dat) != "Git Hg" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestUTF8Filename(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("よんでね.txt", "よんでね")
		_, err := RunGitGsub("でね", "だよ")
		if err != nil {
			t.Errorf("Command failed: %s", err)
		}

		dat, _ := os.ReadFile("./よんでね.txt")
		if string(dat) != "よんだよ" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestSimpleRename(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README-git_gsub.md", "GitGsub git_gsub git-gsub")
		_, err := RunGitGsub("--snake", "--rename", "GitGsub", "SvnGsub")
		if err != nil {
			t.Errorf("Command failed: %s", err)
		}

		dat, _ := os.ReadFile("./README-svn_gsub.md")
		if string(dat) != "SvnGsub svn_gsub git-gsub" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestRuby(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("./foo_bar/baz.rb", "module FooBar::Baz; foo_bar baz # foo_bar/baz; end")
		_, err := RunGitGsub("--ruby", "--rename", "FooBar::Baz", "QuxQuux::Quuz")
		if err != nil {
			t.Errorf("Command failed: %s", err)
		}

		dat, _ := os.ReadFile("./qux_quux/quuz.rb")
		if string(dat) != "module QuxQuux::Quuz; foo_bar baz # qux_quux/quuz; end" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}
func TestAll(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("./foo_bar.rb", "module FooBar; foo_bar foo-bar; end")
		_, err := RunGitGsub("--all", "--rename", "FooBar", "BazQux")
		if err != nil {
			t.Errorf("Command failed: %s", err)
		}

		dat, _ := os.ReadFile("./baz_qux.rb")
		if string(dat) != "module BazQux; baz_qux baz-qux; end" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestAllDoesntImplyRuby(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("./foo_bar/baz.rb", "module FooBar::Baz; foo_bar baz # foo_bar/baz; end")
		_, err := RunGitGsub("--all", "--rename", "FooBar::Baz", "QuxQuux::Quuz")
		if err != nil {
			t.Errorf("Command failed: %s", err)
		}

		dat, _ := os.ReadFile("./foo_bar/baz.rb")
		if string(dat) != "module QuxQuux::Quuz; foo_bar baz # foo_bar/baz; end" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestAllPlusRuby(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("./foo_bar/baz.rb", "FOO_BAR_BAZ=1; module FooBar::Baz; foo_bar baz # foo_bar/baz; end")
		_, err := RunGitGsub("--all", "--ruby", "--rename", "FooBar::Baz", "QuxQuux::Quuz")
		if err != nil {
			t.Errorf("Command failed: %s", err)
		}

		dat, _ := os.ReadFile("./qux_quux/quuz.rb")
		if string(dat) != "FOO_BAR_BAZ=1; module QuxQuux::Quuz; foo_bar baz # qux_quux/quuz; end" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestRenameWithPath(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("foo/git.rb", "puts 'Git'")
		CommitFile("bar/git.rb", "puts 'Git'")
		_, err := RunGitGsub("git", "svn", "bar", "--rename")
		if err != nil {
			t.Errorf("Command failed: %s", err)
		}

		dat, _ := os.ReadFile("./foo/git.rb")
		if string(dat) != "puts 'Git'" {
			t.Errorf("Failed: %s", string(dat))
		}

		dat, _ = os.ReadFile("./bar/svn.rb")
		if string(dat) != "puts 'Git'" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestRenameWithSubmatch(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("git/lib.rb", "puts 'Git'")
		CommitFile("svn/lib.rb", "puts 'Git'")
		CommitFile("bzr/lib.rb", "puts 'Git'")
		_, err := RunGitGsub("--rename", "(git|svn|bzr)/lib", "lib/$1")
		if err != nil {
			t.Errorf("Command failed: %s", err)
		}

		for _, path := range []string{"git", "svn", "bzr"} {
			files, _ := os.ReadDir(path)
			if len(files) != 0 {
				t.Errorf("Failed: %d", len(files))
			}
		}

		files, _ := os.ReadDir("./lib")
		if len(files) != 3 {
			t.Errorf("Failed: %d", len(files))
		}

		for _, path := range []string{"lib/git.rb", "lib/svn.rb", "lib/bzr.rb"} {
			dat, _ := os.ReadFile(path)
			if string(dat) != "puts 'Git'" {
				t.Errorf("Failed: %s", string(dat))
			}
		}
	})
}

func TestRenameWithSpaceInPath(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("git/l b.rb", "puts 'Git'")
		_, err := RunGitGsub("--rename", "l b.rb", "lib.rb")
		if err != nil {
			t.Errorf("Command failed: %s", err)
		}

		_, err = os.ReadFile("git/lib.rb")
		if err != nil {
			t.Errorf("Failed")
		}
	})
}
