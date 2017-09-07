# git-gsub [![Gem Version](https://badge.fury.io/rb/git-gsub.svg)](http://badge.fury.io/rb/git-gsub)[![Build Status](https://travis-ci.org/fujimura/git-gsub.svg)](https://travis-ci.org/fujimura/git-gsub)

A Git subcommand to do gsub in a repository

## Usage

```
Usage: git gsub [options] FROM TO [PATHS]
    -v, --version                    Print version
        --snake                      Substitute snake-cased expressions
        --camel                      Substitute camel-cased expressions
        --kebab                      Substitute kebab-cased expressions
        --rename                     Rename files
        --dry-run                    Print commands to be run
```

## Example

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

To substitute `CommonLisp` with `VisualBasic` with case-awareness:

```
$ git gsub CommonLisp VisualBasic --kebab --snake --camel
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

Substitute and rename file:

```
$ git gsub ruby haskell --rename --camel
```

Result:

```diff
diff --git a/haskell.hs b/haskell.hs
new file mode 100644
index 0000000..ae86ce3
--- /dev/null
+++ b/haskell.hs
@@ -0,0 +1 @@
+putStrLn "Haskell"
diff --git a/haskell.rb b/haskell.rb
new file mode 100644
index 0000000..9f363d3
--- /dev/null
+++ b/haskell.rb
@@ -0,0 +1 @@
+puts "haskell"
diff --git a/ruby.hs b/ruby.hs
deleted file mode 100644
index 70db14d..0000000
--- a/ruby.hs
+++ /dev/null
@@ -1 +0,0 @@
-putStrLn "Ruby"
diff --git a/ruby.rb b/ruby.rb
deleted file mode 100644
index 966eb68..0000000
--- a/ruby.rb
+++ /dev/null
@@ -1 +0,0 @@
```

## Installation

```
$ gem install git-gsub
```

## Supported platforms

Maybe on many *nix like OS which has Perl and sed.

## Contributing

1. Fork it ( http://github.com/fujimura/git-gsub/fork )
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
