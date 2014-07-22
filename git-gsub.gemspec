# coding: utf-8
lib = File.expand_path('../lib', __FILE__)
$LOAD_PATH.unshift(lib) unless $LOAD_PATH.include?(lib)
require 'git/gsub/version'

Gem::Specification.new do |spec|
  spec.name          = "git-gsub"
  spec.version       = Git::Gsub::VERSION
  spec.authors       = ["Fujimura Daisuke"]
  spec.email         = ["me@fujimuradaisuke.com"]
  spec.summary       = %q{A Git subcommand to do gsub in a repository}
  spec.description   = %q{A Git subcommand to do gsub in a repository}
  spec.homepage      = "https://github.com/fujimura/git-gsub"
  spec.license       = "MIT"

  spec.files         = `git ls-files -z`.split("\x0")
  spec.executables   = spec.files.grep(%r{^bin/}) { |f| File.basename(f) }
  spec.test_files    = spec.files.grep(%r{^(test|spec|features)/})
  spec.require_paths = ["lib"]

  spec.add_development_dependency "bundler", "~> 1.5"
  spec.add_development_dependency "rake"
  spec.add_development_dependency "rspec"
  spec.add_development_dependency "pry"
end
