# git-gsub

A Git subcommand to do gsub in a repository

## Usage

```

Options:
  -camel
        Substitute camel-cased expressions
  -kebab
        Substitute kebab-cased expressions
  -rename
        Rename files with expression
  -snake
        Substitute snake-cased expressions
  -version
        Show version
```

(Note that double hyphen before option is also allowed)

## Example

### Normal substutution

To substitute `Git` with `Subversion`:

```
$ git gsub Git Subversion
```

Result:

```diff
diff --git a/README.md b/README.md
index 2185dbf..393dbc6 100644
--- a/README.md
+++ b/README.md
@@ -1,4 +1,4 @@
-# Git::Gsub
+# Subversion::Gsub
 
 TODO: Write a gem description
 
diff --git a/bin/git-gsub b/bin/git-gsub
index c30f093..03b7c4c 100755
--- a/bin/git-gsub
+++ b/bin/git-gsub
@@ -1,4 +1,4 @@
 #! /usr/bin/env ruby
 
 require 'git/gsub'
-Git::Gsub.run
+Subversion::Gsub.run
```

### Case-converting substitution

To substitute `CommonLisp` with `VisualBasic` with case-awareness:

```
$ git gsub --kebab --snake --camel CommonLisp VisualBasic
```

Result:

```diff
diff --git a/README.md b/README.md
index 2185dbf..393dbc6 100644
--- a/README.md
+++ b/README.md
@@ -1,4 +1,4 @@
-# CommonLisp common_lisp common-lisp
+# VisualBasic visual_basic visual-basic
```

### Substitute and rename file

Substitute and rename file:

```
$ git gsub --rename --camel git svn
```

Result:

```diff
diff --git a/lib/git.rb b/lib/git.rb
deleted file mode 100644
index f0b8633..0000000
--- a/lib/git.rb
+++ /dev/null
@@ -1 +0,0 @@
-puts "GitGsub"
diff --git a/lib/svn.rb b/lib/svn.rb
new file mode 100644
index 0000000..312f137
--- /dev/null
+++ b/lib/svn.rb
@@ -0,0 +1 @@
+puts "SvnGsub"
```

## Installation

```
$ go get github.com/fujimura/git-gsub
```

## Contributing

1. Fork it ( http://github.com/fujimura/git-gsub/fork )
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
