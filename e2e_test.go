package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func RunInTmpDir(run func()) {
	dir, err := ioutil.TempDir("", "")
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

	err = ioutil.WriteFile(name, []byte(content), os.ModePerm)
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
	out, err := exec.Command(GitGsubPath(), args...).Output()
	if err != nil {
		fmt.Println(string(out))
		panic(err)
	}
	return out, err
}

func TestVersion(t *testing.T) {
	out, _ := RunGitGsub("--version")
	if string(out) != "v0.0.14\n" {
		t.Errorf("Failed: %s", string(out))
	}
}

func TestSimpleSubstitution(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README.md", "Git Subversion Bzr")
		RunGitGsub("Bzr", "Mercurial")

		dat, _ := ioutil.ReadFile("./README.md")
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
		RunGitGsub("Bzr", "Mercurial")

		dat, _ := ioutil.ReadFile("./README_1.md")
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

		RunGitGsub("Git", "Svn", "foo")

		dat, _ := ioutil.ReadFile("./README.md")
		if string(dat) != "Git Subversion Bzr" {
			t.Errorf("Failed: %s", string(dat))
		}

		dat, _ = ioutil.ReadFile("./foo/git")
		if string(dat) != "Svn Subversion Bzr" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestSubstitutionWithCaseConversion(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README.md", "GitGsub git_gsub git-gsub")
		RunGitGsub("--camel", "--kebab", "--snake", "git-gsub", "svn-gsub")
		dat, _ := ioutil.ReadFile("./README.md")
		if string(dat) != "SvnGsub svn_gsub svn-gsub" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestOptionsCanBePutAfterArguments(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README.md", "GitGsub git_gsub git-gsub")
		RunGitGsub("git-gsub", "svn-gsub", "--camel", "--kebab", "--snake")
		dat, _ := ioutil.ReadFile("./README.md")
		if string(dat) != "SvnGsub svn_gsub svn-gsub" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestSubstitutionWithFixedStringOption(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("hello.rb", "puts('hello')")
		RunGitGsub("--fgrep", "(", " ")
		RunGitGsub("-F", ")", "")
		dat, _ := ioutil.ReadFile("./hello.rb")
		if string(dat) != "puts 'hello'" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestEscape(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README.md", `<h1 class="foo">`)
		RunGitGsub(`<h1 class="foo">`, `<h1 class="bar">`)
		dat, _ := ioutil.ReadFile("./README.md")
		if string(dat) != `<h1 class="bar">` {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestAtMark(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README.md", "foo@example.com")
		RunGitGsub("foo@example.com", "bar@example.com")
		dat, _ := ioutil.ReadFile("./README.md")
		if string(dat) != "bar@example.com" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestConsequesingAtMark(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README.md", "foo@example.com")
		RunGitGsub("foo@example.com", "bar@@example.com")
		dat, _ := ioutil.ReadFile("./README.md")
		if string(dat) != "bar@@example.com" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestDoubleQuoteToSingleQuote(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README.md", `hello this is "git"`)
		RunGitGsub(`"git"`, `'svn'`)
		dat, _ := ioutil.ReadFile("./README.md")
		if string(dat) != "hello this is 'svn'" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestSingleQuoteToDoubleQuote(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README.md", `hello this is 'git'`)
		RunGitGsub(`'git'`, `"svn"`)
		dat, _ := ioutil.ReadFile("./README.md")
		if string(dat) != `hello this is "svn"` {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestBracket(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README.md", `{git{svn}}`)
		RunGitGsub("{git{svn}}", "{hg{svn}}")
		dat, _ := ioutil.ReadFile("./README.md")
		if string(dat) != "{hg{svn}}" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestSubmatch(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README.md", "git-foo-1 git-bar-22 git-baz-3")
		RunGitGsub(`git-([a-z]+)-([\d]{1,2})`, `$2-$1`)
		dat, _ := ioutil.ReadFile("./README.md")
		if string(dat) != "1-foo 22-bar 3-baz" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestSubstituteToEmptyString(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README.md", "Git Svn Hg")
		RunGitGsub("Svn ", "")
		dat, _ := ioutil.ReadFile("./README.md")
		if string(dat) != "Git Hg" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestUTF8Filename(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("よんでね.txt", "よんでね")
		RunGitGsub("でね", "だよ")
		dat, _ := ioutil.ReadFile("./よんでね.txt")
		if string(dat) != "よんだよ" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestSimpleRename(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("README-git_gsub.md", "GitGsub git_gsub git-gsub")
		RunGitGsub("--snake", "--rename", "GitGsub", "SvnGsub")
		dat, _ := ioutil.ReadFile("./README-svn_gsub.md")
		if string(dat) != "SvnGsub svn_gsub git-gsub" {
			t.Errorf("Failed: %s", string(dat))
		}
	})
}

func TestRenameWithPath(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("foo/git.rb", "puts 'Git'")
		CommitFile("bar/git.rb", "puts 'Git'")
		RunGitGsub("git", "svn", "bar", "--rename")

		dat, _ := ioutil.ReadFile("./foo/git.rb")
		if string(dat) != "puts 'Git'" {
			t.Errorf("Failed: %s", string(dat))
		}

		dat, _ = ioutil.ReadFile("./bar/svn.rb")
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
		RunGitGsub("--rename", "(git|svn|bzr)/lib", "lib/$1")

		for _, path := range []string{"git", "svn", "bzr"} {
			files, _ := ioutil.ReadDir(path)
			if len(files) != 0 {
				t.Errorf("Failed: %d", len(files))
			}
		}

		files, _ := ioutil.ReadDir("./lib")
		if len(files) != 3 {
			t.Errorf("Failed: %d", len(files))
		}

		for _, path := range []string{"lib/git.rb", "lib/svn.rb", "lib/bzr.rb"} {
			dat, _ := ioutil.ReadFile(path)
			if string(dat) != "puts 'Git'" {
				t.Errorf("Failed: %s", string(dat))
			}
		}
	})
}

func TestRenameWithSpaceInPath(t *testing.T) {
	RunInTmpRepo(func() {
		CommitFile("git/l b.rb", "puts 'Git'")
		RunGitGsub("--rename", "l b.rb", "lib.rb")
		_, err := ioutil.ReadFile("git/lib.rb")
		if err != nil {
			t.Errorf("Failed")
		}
	})
}
