# git-gsub

A Git subcommand to do gsub in a repository

## Usage

To substitute `Git` with `Subversion`, run

```
$ git gsub Git Subversion
```

Then you will get

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

## Installation

```
$ gem install git-gsub
```

## Contributing

1. Fork it ( http://github.com/<my-github-username>/git-gsub/fork )
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
