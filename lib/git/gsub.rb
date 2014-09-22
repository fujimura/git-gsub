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

      if to.nil?
        abort "No argument to gsub was given"
      end

      target_files = (`git grep -l #{from} #{path}`).each_line.map(&:chomp).join ' '

      if system_support_gsed?
        system %|gsed -i s/#{from}/#{to}/g #{target_files}|
      else
        system %|sed -i "" -e s/#{from}/#{to}/g #{target_files}|
      end
    end

    private

    def self.system_support_gsed?
      `which gsed`
      $?.success?
    end
  end
end
