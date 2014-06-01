require "git/gsub/version"
require 'shellwords'

module Git
  module Gsub
    def self.run
      case ARGV.first
      when '-v', '--version'
        version
      else
        gsub *ARGV
      end
    end

    def self.version
      puts Git::Gsub::VERSION
    end

    def self.gsub *args
      from, to, path, = args

      target_files = (`git grep -l #{from} #{path}`).each_line.map(&:chomp).join ' '

      system %|gsed -i s/#{Shellwords.escape from}/#{Shellwords.escape to}/g #{target_files}|
    end
  end
end
