# Changelog
All notable changes to this project will be documented in this file.

## [0.0.14]

### Added

- Add -F option for "fixed string" (a.k.a. fgrep) search ([#13](https://github.com/fujimura/git-gsub/pull/13)) Thanks @amatsuda!

### Fixed

- Handle multiplyte filenames properly ([#16](https://github.com/fujimura/git-gsub/pull/16)) Thanks @amatsuda!

## [0.0.13]

- Internal changes only

## [0.0.12]
### Added
- Substitution and renaming can be done with submatch
### Removed
- Deprecate `--dry-run` since git-gsub is no longer an one-liner generator

### Internal
- Rewrite in Go

## [0.0.11]
### Changed
- Support newer and safer activesupport ([#13](https://github.com/fujimura/git-gsub/pull/13))
- Add deprecation message on installation ([#13](https://github.com/fujimura/git-gsub/pull/13))

## [0.0.10]
### Changed
- Support activesupport >= 4 & < 6 ([#12](https://github.com/fujimura/git-gsub/pull/12)) Thanks @hanachin!

## [0.0.9]
### Fixed
- Make `--dry-run` works with `--rename`([#10](https://github.com/fujimura/git-gsub/pull/10))

## [0.0.8]
### Fixed
- Fix bug on passing file names to `Open3.capture3`[c05a05c](https://github.com/fujimura/git-gsub/commit/c05a05cd413d5a389c781b6649b42a46a825c4db)

## [0.0.7]
### Fixed
- Fix escaping @([#8](https://github.com/fujimura/git-gsub/pull/8)) Thanks @amatsuda!

## [0.0.6]
### Added
- Add option to substitute in camel/kebab/snake case
- Add option to rename files
### Changed
- Do substitution with Perl, not sed to reduce platform issues
- No escape required to command line arguments

## [0.0.5] - 2017-05-23
### Fixed
- Fix bug on assigning target directories

## (yanked)[0.0.4] - 2017-05-23
### Fixed
- Target directories can be specified correctly.
