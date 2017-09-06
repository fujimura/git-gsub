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
        yield filename
      end
    end
  end

  it 'should substitute files' do
    run_in_directory_with_a_file 'README.md', 'Git Subversion Bzr' do |filename|
      Git::Gsub.run %w[Bzr Mercurial]
      file = File.read(filename).chomp
      expect(file).to eq 'Git Subversion Mercurial'
    end
  end

  it 'should substitute files with case conversion' do
    run_in_directory_with_a_file 'README.md', 'GitGsub git_gsub git-gsub' do |filename|
      Git::Gsub.run %w[GitGsub SvnGsub --camel --kebab --snake]
      file = File.read(filename).chomp
      expect(file).to eq 'SvnGsub svn_gsub svn-gsub'
    end
  end

  it 'should escape well' do
    run_in_directory_with_a_file 'README.md', %(<h1 class="foo">) do |filename|
      Git::Gsub.run [%(<h1 class="foo">), %(<h1 class="bar">)]
      file = File.read(filename).chomp
      expect(file).to eq %(<h1 class="bar">)
    end
  end

  it 'should not create backup file' do
    run_in_directory_with_a_file 'README.md', 'Git Subversion Bzr' do |_filename|
      Git::Gsub.run %w[Bzr Darcs]
      expect(`ls`).to eql "README.md\n"
    end
  end

  it 'should rename with --rename' do
    run_in_directory_with_a_file 'README-git_gsub.md', 'GitGsub git_gsub git-gsub' do |_filename|
      Git::Gsub.run %w[GitGsub SvnGsub --snake --rename]
      expect(`ls`).to eql "README-svn_gsub.md\n"
    end
  end

  it 'should do nothing if no file found' do
    run_in_directory_with_a_file 'README-git_gsub.md', 'GitGsub git_gsub git-gsub' do |_filename|
      expect {
        Git::Gsub.run %w[Atlanta Chicago --snake --rename]
      }.not_to raise_error
    end
  end
end
