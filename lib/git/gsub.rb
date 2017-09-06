require 'git/gsub/version'
require 'active_support/inflector'
require 'optparse'
require 'shellwords'
require 'English'

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

      opts.on('--snake', 'Substitute snake-cased expressions') do |_v|
        options[:snake] = true
      end

      opts.on('--camel', 'Substitute camel-cased expressions') do |_v|
        options[:camel] = true
      end

      opts.on('--kebab', 'Substitute kebab-cased expressions') do |_v|
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
            commands.each { |c| puts c }
          else
            commands.each { |c| system c }
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
          commands = args.map { |from, to| build_command(from, to, paths) }
          run_commands commands
        end

        def build_command(from, to, paths = [], _options = {})
          from, to, *paths = [from, to, *paths].map { |s| Shellwords.escape s }

          target_files = `git grep -l #{from} #{paths.join ' '}`.each_line.map(&:chomp).join ' '

          %(perl -pi -e 's/#{from}/#{to}/g' #{target_files})
        end
      end

      class Rename < Command
        def run
          commands = args.map { |from, to| build_command(from, to, paths) }
          run_commands commands
        end

        def build_command(from, to, paths = [], _options = {})
          from, to, *paths = [from, to, *paths].map { |s| Shellwords.escape s }

          %(rename -p 's/#{from}/#{to}/g' `git ls-files #{paths.join ' '}`)
        end
      end
    end
  end
end
