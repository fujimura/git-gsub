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
      from, to, path, = args.map do |arg|
        Shellwords.escape arg if arg
      end

      target_files = (`git grep -l #{from} #{path}`).each_line.map(&:chomp).join ' '

      system %|gsed -i s/#{from}/#{to}/g #{target_files}|
    end
  end
end
