require 'rspec'
require 'tmpdir'
require 'pry'

RSpec.configure do |config|
  config.filter_run :focus => true
  config.run_all_when_everything_filtered = true
end

describe 'git-gsub' do
  def run_in_tmp_repo
    Dir.mktmpdir do |dir|
      Dir.chdir dir do
        `git init`
        `git config --local user.email "you@example.com"`
        `git config --local user.name "Your Name"`
        `git add .`
        `git commit -m init`
        yield
      end
    end
  end

  def commit_file(name, content)
    FileUtils.mkdir_p(File.dirname(name))
    File.write(name, content)
    `git add .`
    `git commit -m 'Add #{name}'`
  end

  let!(:git_gsub_path) {
    if ENV['USE_RUBY']
      File.expand_path("../../bin/git-gsub-ruby", __FILE__)
    else
      File.expand_path("../../bin/git-gsub", __FILE__)
    end
  }

  around do |example|
    run_in_tmp_repo do
      example.run
    end
  end

  describe 'Version' do
    it 'should show version with --version' do
      output = `#{git_gsub_path} --version`

      expect(output).to eq "v0.0.1\n"
    end
  end

  describe 'Substituting' do
    it 'should substitute files' do
      commit_file 'README.md', 'Git Subversion Bzr'
      `#{git_gsub_path} Bzr Mercurial`

      expect(File.read('README.md')).to eq 'Git Subversion Mercurial'
    end

    it 'should substitute files in specified path' do
      commit_file 'README.md', 'Git Subversion Bzr'
      commit_file 'foo/git.rb', 'puts "Git"'
      commit_file 'bar/git.rb', 'puts "Git"'
      commit_file 'baz/git.rb', 'puts "Git"'
      `#{git_gsub_path} Git Svn foo baz`

      expect(File.read('README.md')).to eq 'Git Subversion Bzr'
      expect(File.read('foo/git.rb')).to eq 'puts "Svn"'
      expect(File.read('bar/git.rb')).to eq 'puts "Git"'
      expect(File.read('baz/git.rb')).to eq 'puts "Svn"'
    end

    it 'should substitute files with case conversion' do
      commit_file 'README.md', 'GitGsub git_gsub git-gsub'
      `#{git_gsub_path} --camel --kebab --snake GitGsub SvnGsub`

      expect(File.read('README.md')).to eq 'SvnGsub svn_gsub svn-gsub'
    end

    it 'should substitute files with case conversion' do
      commit_file 'README.md', 'GitGsub git_gsub git-gsub'
      `#{git_gsub_path} --camel --kebab --snake git-gsub svn-gsub`

      expect(File.read('README.md')).to eq 'SvnGsub svn_gsub svn-gsub'
    end

    it 'options can be put after arguments' do
      commit_file 'README.md', 'GitGsub git_gsub git-gsub'
      `#{git_gsub_path} git-gsub svn-gsub --camel --kebab --snake`

      expect(File.read('README.md')).to eq 'SvnGsub svn_gsub svn-gsub'
    end

    it 'should escape well' do
      commit_file 'README.md', %(<h1 class="foo">)
      `#{git_gsub_path} '<h1 class="foo">' '<h1 class="bar">'`

      expect(File.read('README.md')).to eq %(<h1 class="bar">)
    end

    it 'should substutute @' do
      commit_file 'README.md', %(foo@example.com)
      `#{git_gsub_path} foo@example bar@example`

      expect(File.read('README.md')).to eq %(bar@example.com)
    end

    it 'should substitute consequenting @' do
      commit_file 'README.md', %(Hello this is @git)
      `#{git_gsub_path} @git @@svn`

      expect(File.read('README.md')).to eq %(Hello this is @@svn)
    end

    it %(should substitute " to ') do
      commit_file 'README.md', %(Hello this is "git")
      `#{git_gsub_path} '"git"' "'svn'"`

      expect(File.read('README.md')).to eq %(Hello this is 'svn')
    end

    it %(should substitute ' to ") do
      commit_file 'README.md', %(Hello this is 'git')
      `#{git_gsub_path} "'git'" '"svn"'`

      expect(File.read('README.md')).to eq %(Hello this is "svn")
    end

    it 'should substitute text including { and }'do
      commit_file 'README.md', %({git{svn})
      `#{git_gsub_path} {git{svn} {hg{svn}}`

      expect(File.read('README.md')).to eq %({hg{svn}})
    end

    it 'should not create backup file' do
      commit_file 'README.md', 'Git Subversion Bzr'
      `#{git_gsub_path} Bzr Darcs`

      expect(`ls`).to eql "README.md\n"
    end

    it 'should substitute files using submatch' do
      commit_file 'README.md', 'git-foo-1 git-bar-22 git-baz-3'
      `#{git_gsub_path} 'git-([a-z]+)-([\\d]{1,2})' '$2-$1'`

      expect(File.read('README.md')).to eq '1-foo 22-bar 3-baz'
    end

  end

  describe 'Renaming' do
    it 'should rename with --rename' do
      commit_file 'README-git_gsub.md', 'GitGsub git_gsub git-gsub'
      `#{git_gsub_path} --snake --rename GitGsub SvnGsub`

      expect(`ls`).to eql "README-svn_gsub.md\n"
      expect(File.read('README-svn_gsub.md')).to eq 'SvnGsub svn_gsub git-gsub'
    end

    it 'should rename with --rename' do
      commit_file 'git.rb', 'puts "Git"'
      commit_file 'foo/git.rb', 'puts "Git"'
      commit_file 'bar/git.rb', 'puts "Git"'
      commit_file 'baz/git.rb', 'puts "Git"'
      `#{git_gsub_path} --rename git svn foo baz`

      expect(`ls foo`).to eql "svn.rb\n"
      expect(`ls bar`).to eql "git.rb\n"
      expect(`ls baz`).to eql "svn.rb\n"
      expect(`ls .`).to eql "bar\nbaz\nfoo\ngit.rb\n"
    end

    it 'should do nothing if no file found' do
      commit_file 'README-git_gsub.md', 'GitGsub git_gsub git-gsub'

      expect {
        `#{git_gsub_path} --snake --rename Atlanta Chicago`
      }.not_to raise_error
    end

    it 'should rename with --rename using submatch' do
      commit_file 'git/lib.rb', 'puts "Git"'
      commit_file 'svn/lib.rb', 'puts "Git"'
      commit_file 'bzr/lib.rb', 'puts "Git"'
      `#{git_gsub_path} --rename '(git|svn|bzr)/lib' 'lib/$1'`

      expect(`ls git`).to eql ""
      expect(`ls svn`).to eql ""
      expect(`ls bzr`).to eql ""
      expect(`ls lib`).to eql "bzr.rb\ngit.rb\nsvn.rb\n"
    end

    it 'should rename a file which has space in filename' do
      commit_file 'git/l b.rb', 'puts "Git"'
      `#{git_gsub_path} --rename git svn`

      expect(`ls git`).to eql ""
    end
  end
end
