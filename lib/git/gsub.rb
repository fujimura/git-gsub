require 'git/gsub/version'
require 'active_support/inflector'
require 'optparse'
require 'shellwords'
require 'English'

module Git
  module Gsub
    def self.run(argv)
      options = {}
      OptionParser.new([]) do |opts|
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

        opts.on('--rename') do |_v|
          options[:rename] = true
        end

        opts.on('--dry-run') do |_v|
          options[:dry] = true
        end
      end.parse!(argv)

      if options[:version]
        version
      else
        from, to, *paths = argv
        Commands::Gsub.new(from, to, paths, options).run
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

        def system_support_gsed?
          `which gsed`
          $CHILD_STATUS.success?
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
          commands = args.map { |from, to| build_commands(from, to, paths) }
          run_commands commands
        end

        def build_commands(from, to, paths = [], _options = {})
          abort 'No argument to gsub was given' if to.nil?

          from, to, *paths = [from, to, *paths].map { |s| Shellwords.escape s }

          target_files = `git grep -l #{from} #{paths.join ' '}`.each_line.map(&:chomp).join ' '

          if system_support_gsed?
            %(gsed -i s/#{from}/#{to}/g #{target_files})
          else
            %(sed -i "" -e s/#{from}/#{to}/g #{target_files})
          end
        end
      end
    end
  end
end
