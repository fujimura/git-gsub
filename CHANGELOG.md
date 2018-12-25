# Changelog
All notable changes to this project will be documented in this file.

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
