require 'git/gsub/version'
require 'active_support/inflector'
require 'optparse'
require 'shellwords'
require 'English'
require "open3"

module Git
  module Gsub
    def self.run(argv)
      opts = OptionParser.new

      self.class.module_eval do
        define_method(:usage) do |msg = nil|
          puts opts.to_s
          puts "error: #{msg}" if msg
          exit 1
        end
      end

      options = {}
      opts.on('-v', '--version', 'Print version') do |_v|
        options[:version] = true
      end

      opts.on('--snake', 'Experimental: substitute snake-cased expressions') do |_v|
        options[:snake] = true
      end

      opts.on('--camel', 'Experimental: substitute camel-cased expressions') do |_v|
        options[:camel] = true
      end

      opts.on('--kebab', 'Experimental: substitute kebab-cased expressions') do |_v|
        options[:kebab] = true
      end

      opts.on('--rename', 'Rename files along with substitution') do |_v|
        options[:rename] = true
      end

      opts.on('--dry-run', 'Print commands to be run') do |_v|
        options[:dry] = true
      end

      opts.banner = 'Usage: git gsub [options] FROM TO [PATHS]'

      begin
        args = opts.parse(argv)
      rescue OptionParser::InvalidOption => e
        usage e.message
      end

      if options[:version]
        version
      else
        from, to, *paths = args
        usage 'No FROM given' unless from
        usage 'No TO given' unless to
        Commands::Gsub.new(from, to, paths, options).run
        Commands::Rename.new(from, to, paths, options).run if options[:rename]
      end
    end

    def self.version
      puts Git::Gsub::VERSION
    end

    module Commands
      class Command
        attr_accessor :from, :to, :paths, :options

        def initialize(from, to, paths = [], options = {})
          @from = from
          @to = to
          @paths = paths
          @options = options
        end

        def run_commands(commands)
          if options[:dry]
            commands.each { |args| puts args.join(' ') }
          else
            commands.each { |args| Open3.capture3(*args) }
          end
        end

        def args
          args = []
          args << [from, to]
          args << [from.camelcase, to.camelcase] if options[:camel]
          args << [from.underscore, to.underscore] if options[:snake]
          args << [from.underscore.dasherize, to.underscore.dasherize] if options[:kebab]

          args.compact
        end
      end

      class Gsub < Command
        def run
          commands = args.map { |from, to| build_command(from, to, paths) }.compact
          run_commands commands
        end

        def build_command(from, to, paths = [], _options = {})
          from, to, *paths = [from, to, *paths].map { |s| Shellwords.escape s }
          [from, to].each { |s| s.gsub! '@', '\@' }

          target_files = `git grep -l #{from} #{paths.join ' '}`.each_line.map(&:chomp).join ' '
          return if target_files.empty?


          ['perl', '-pi', '-e', "s{#{from}}{#{to}}g", *target_files]
        end
      end

      class Rename < Command
        def run
          commands = args.map { |from, to| build_command(from, to, paths) }
          run_commands commands
        end

        def build_command(from, to, paths = [], _options = {})
          from, to, *paths = [from, to, *paths].map { |s| Shellwords.escape s }

          %[for f in $(git ls-files #{paths.join ' '} | grep #{from}); do mkdir -p $(dirname $(echo $f | sed 's/#{from}/#{to}/g')); git mv $f $(echo $f | sed 's/#{from}/#{to}/g') 2>/dev/null; done]
        end
      end
    end
  end
end
