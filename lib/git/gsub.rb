require 'git/gsub/version'
require 'active_support/inflector'
require 'optparse'
require 'shellwords'
require 'English'

module Git
  module Gsub
    def self.run
      options = {}
      OptionParser.new do |opts|
        # TODO
        opts.banner = 'Usage: example.rb [options]'

        opts.on('-v', '--version') do |_v|
          options[:version] = true
        end

        opts.on('--snake') do |_v|
          options[:snake] = true
        end

        opts.on('--camel') do |_v|
          options[:camel] = true
        end

        opts.on('--kebab') do |_v|
          options[:kebab] = true
        end

        opts.on('--dry-run') do |_v|
          options[:dry] = true
        end
      end.parse!

      if options[:version]
        version
      else
        from, to, *paths = ARGV
        gsub(from, to, paths, options)
      end
    end

    def self.version
      puts Git::Gsub::VERSION
    end

    def self.gsub(from, to, paths=[], options={})
      commands = []
      commands << build_commands(from, to, paths)

      commands << build_commands(from.camelcase, to.camelcase, paths) if options[:camel]
      commands << build_commands(from.underscore, to.underscore, paths) if options[:snake]
      commands << build_commands(from.underscore.dasherize, to.underscore.dasherize, paths) if options[:kebab]

      if options[:dry]
        commands.each { |c| puts c }
      else
        commands.each { |c| system c }
      end
    end

    def self.build_commands(from, to, paths)
      abort 'No argument to gsub was given' if to.nil?

      from, to, *paths = [from, to, *paths].map {|s| Shellwords.escape s }

      target_files = `git grep -l #{from} #{paths.join ' '}`.each_line.map(&:chomp).join ' '

      if system_support_gsed?
        %(gsed -i s/#{from}/#{to}/g #{target_files})
      else
        %(sed -i "" -e s/#{from}/#{to}/g #{target_files})
      end
    end

    private

    def self.system_support_gsed?
      `which gsed`
      $CHILD_STATUS.success?
    end
  end
end
