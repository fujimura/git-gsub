# coding: utf-8
lib = File.expand_path('../lib', __FILE__)
$LOAD_PATH.unshift(lib) unless $LOAD_PATH.include?(lib)
require 'git/gsub/version'

GIT_GSUB_RUBY_IS_DEPRECATED = %q{git-gsub: Newer version is written in Go and Ruby implementation is deprecated. See https://github.com/fujimura/git-gsub}

Gem::Specification.new do |spec|
  spec.name          = "git-gsub"
  spec.version       = Git::Gsub::VERSION
  spec.authors       = ["Fujimura Daisuke"]
  spec.email         = ["me@fujimuradaisuke.com"]
  spec.summary       = GIT_GSUB_RUBY_IS_DEPRECATED
  spec.description   = GIT_GSUB_RUBY_IS_DEPRECATED
  spec.homepage      = "https://github.com/fujimura/git-gsub"
  spec.license       = "MIT"
  spec.post_install_message = GIT_GSUB_RUBY_IS_DEPRECATED

  spec.files         = `git ls-files -z`.split("\x0")
  spec.executables   = spec.files.grep(%r{^bin/}) { |f| File.basename(f) }
  spec.test_files    = spec.files.grep(%r{^(test|spec|features)/})
  spec.add_dependency "activesupport", ">= 6.0.3.1"
  spec.require_paths = ["lib"]

  spec.add_development_dependency "bundler"
  spec.add_development_dependency "rake"
  spec.add_development_dependency "rspec"
  spec.add_development_dependency "pry"
end
