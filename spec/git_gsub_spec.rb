require 'rspec'
require 'tmpdir'
require 'git/gsub'
require 'pry'

describe 'git-gsub' do
  def run_in_directory_with_a_file(filename, content)
    Dir.mktmpdir do |dir|
      Dir.chdir dir do
        File.open(filename, 'w') { |f| f << content }
        `git init`
        `git config --local user.email "you@example.com"`
        `git config --local user.name "Your Name"`
        `git add .`
        `git commit -m init`
        yield
      end
    end
  end

  it 'should substitute files' do
    run_in_directory_with_a_file 'README.md', 'Git Subversion Bzr' do
      Git::Gsub.run %w[Bzr Mercurial]
      expect(File.read('README.md')).to eq 'Git Subversion Mercurial'
    end
  end

  it 'should substitute files with case conversion' do
    run_in_directory_with_a_file 'README.md', 'GitGsub git_gsub git-gsub' do
      Git::Gsub.run %w[GitGsub SvnGsub --camel --kebab --snake]
      expect(File.read('README.md')).to eq 'SvnGsub svn_gsub svn-gsub'
    end
  end

  it 'should escape well' do
    run_in_directory_with_a_file 'README.md', %(<h1 class="foo">) do
      Git::Gsub.run [%(<h1 class="foo">), %(<h1 class="bar">)]
      expect(File.read('README.md')).to eq %(<h1 class="bar">)
    end

    run_in_directory_with_a_file 'README.md', %(git/gsub) do
      Git::Gsub.run [%(git/gsub), %(svn/sub)]
      expect(File.read('README.md')).to eq %(svn/sub)
    end

    run_in_directory_with_a_file 'README.md', %({git{svn}) do
      Git::Gsub.run [%({git{svn}), %({hg{svn})]
      expect(File.read('README.md')).to eq %({hg{svn})
    end
  end

  it 'should not create backup file' do
    run_in_directory_with_a_file 'README.md', 'Git Subversion Bzr' do
      Git::Gsub.run %w[Bzr Darcs]
      expect(`ls`).to eql "README.md\n"
    end
  end

  it 'should rename with --rename' do
    run_in_directory_with_a_file 'README-git_gsub.md', 'GitGsub git_gsub git-gsub' do
      Git::Gsub.run %w[GitGsub SvnGsub --snake --rename]
      expect(`ls`).to eql "README-svn_gsub.md\n"
      expect(File.read('README-svn_gsub.md')).to eq 'SvnGsub svn_gsub git-gsub'
    end
  end

  it 'should do nothing if no file found' do
    run_in_directory_with_a_file 'README-git_gsub.md', 'GitGsub git_gsub git-gsub' do
      expect {
        Git::Gsub.run %w[Atlanta Chicago --snake --rename]
      }.not_to raise_error
    end
  end
end
