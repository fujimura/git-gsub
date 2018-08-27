require 'rspec'
require 'tmpdir'
require 'git/gsub'
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

  around do |example|
    run_in_tmp_repo do
      example.run
    end
  end

  it 'should substitute files' do
    commit_file 'README.md', 'Git Subversion Bzr'
    Git::Gsub.run %w[Bzr Mercurial]

    expect(File.read('README.md')).to eq 'Git Subversion Mercurial'
  end

  it 'should substitute files with case conversion' do
    commit_file 'README.md', 'GitGsub git_gsub git-gsub'
    Git::Gsub.run %w[GitGsub SvnGsub --camel --kebab --snake]

    expect(File.read('README.md')).to eq 'SvnGsub svn_gsub svn-gsub'
  end

  it 'should escape well' do
    commit_file 'README.md', %(<h1 class="foo">)
    Git::Gsub.run [%(<h1 class="foo">), %(<h1 class="bar">)]

    expect(File.read('README.md')).to eq %(<h1 class="bar">)
  end

  it 'should substutute @' do
    commit_file 'README.md', %(foo@example.com)
    Git::Gsub.run [%(@example), %(bar@example)]

    expect(File.read('README.md')).to eq %(foobar@example.com)
  end

  it 'should substitute consequenting @' do
    commit_file 'README.md', %(Hello this is @git)
    Git::Gsub.run [%(@git), %(@@svn)]

    expect(File.read('README.md')).to eq %(Hello this is @@svn)
  end

  it %(should substitute " to ') do
    commit_file 'README.md', %(Hello this is "git")
    Git::Gsub.run [%("git"), %('svn')]

    expect(File.read('README.md')).to eq %(Hello this is 'svn')
  end

  it %(should substitute ' to ") do
    commit_file 'README.md', %(Hello this is 'git')
    Git::Gsub.run [%('git'), %("svn")]

    expect(File.read('README.md')).to eq %(Hello this is "svn")
  end

  it 'should substitute text including { and }'do
    commit_file 'README.md', %({git{svn})
    Git::Gsub.run [%({git{svn}), %({hg{svn})]

    expect(File.read('README.md')).to eq %({hg{svn})
  end

  it 'should substitute text to empty'do
    commit_file 'README.md', "Git Svn Hg"
    Git::Gsub.run [%(Svn ), %()]

    expect(File.read('README.md')).to eq %(Git Hg)
  end

  it 'should not create backup file' do
    commit_file 'README.md', 'Git Subversion Bzr'
    Git::Gsub.run %w[Bzr Darcs]

    expect(`ls`).to eql "README.md\n"
  end

  it 'should rename with --rename' do
    commit_file 'README-git_gsub.md', 'GitGsub git_gsub git-gsub'
    Git::Gsub.run %w[GitGsub SvnGsub --snake --rename]

    expect(`ls`).to eql "README-svn_gsub.md\n"
    expect(File.read('README-svn_gsub.md')).to eq 'SvnGsub svn_gsub git-gsub'
  end

  it 'should rename with --rename' do
    commit_file 'lib/git.rb', 'puts "Git"'
    Git::Gsub.run %w[git svn --camel --rename]

    expect(`ls lib`).to eql "svn.rb\n"
    expect(File.read('lib/svn.rb')).to eq 'puts "Svn"'
  end

  it 'should do nothing if no file found' do
    commit_file 'README-git_gsub.md', 'GitGsub git_gsub git-gsub'

    expect {
      Git::Gsub.run %w[Atlanta Chicago --snake --rename]
    }.not_to raise_error
  end

  it 'should output command with dry-run' do
    commit_file 'README-git_gsub.md', 'GitGsub git_gsub git-gsub'

    expect {
      Git::Gsub.run %w[GitGsub SvnGsub --snake --rename --dry-run]
    }.to output(/Svn/).to_stdout
  end
end
