require 'rspec'
require 'tmpdir'
require 'git/gsub'
require 'pry'

describe 'git-gsub' do
  def run_in_directory_with_a_file content
    Dir.mktmpdir do |dir|
      filename = "FOO"
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
    run_in_directory_with_a_file "Git Subversion Bzr" do |filename|
      Git::Gsub.gsub "Bzr", "Mercurial"
      file = File.read(filename).chomp
      expect(file).to eq "Git Subversion Mercurial"
    end
  end

  it 'should escape well' do
    run_in_directory_with_a_file %|<h1 class="foo">| do |filename|
      Git::Gsub.gsub %|<h1 class="foo">|, %|<h1 class="bar">|
      file = File.read(filename).chomp
      expect(file).to eq %|<h1 class="bar">|
    end
  end
end
