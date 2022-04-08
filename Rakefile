# frozen_string_literal: true

require 'open3'

task :command_exists, [:command] do |_, args|
  abort "#{args.command} doesn't exists" if `command -v #{args.command} > /dev/null 2>&1 && echo $?`.chomp.empty?
end

task :repo_clean do
  abort 'please commit your changes first!' unless `git status -s | wc -l`.strip.to_i.zero?
end

task :current_version do
  version_file = File.open('.bumpversion.cfg', 'r')
  data = version_file.read
  version_file.close
  match = /current_version = (\d+).(\d+).(\d+)/.match(data)
  "#{match[1]}.#{match[2]}.#{match[3]}"
end

task :has_bumpversion do
  Rake::Task['command_exists'].invoke('bumpversion')
end

task :get_package_name do
  origin_url = `git config --get remote.origin.url`.strip

  abort "Error! your origin url (#{origin_url}) is not valid!\nShould be something like:\n\n  git@github.com:USER-ORG/REPO.git\n\n" if origin_url.include?('http')
  origin_url.split('@').last.split('.git').first.sub(':','/')
end

AVAILABLE_REVISIONS = %w[major minor patch].freeze
task :bump, [:revision] => [:has_bumpversion] do |_, args|
  args.with_defaults(revision: 'patch')
  unless AVAILABLE_REVISIONS.include?(args.revision)
    abort "Please provide valid revision: #{AVAILABLE_REVISIONS.join(',')}"
  end

  system "bumpversion #{args.revision}"
end

desc 'show avaliable tasks (default task)'
task :default do
  system('rake -sT')
end

desc 'run tests'
task :test, [:verbose] do |_, args|
  args.with_defaults(verbose: '')
  system "go test #{args.verbose} ./..."
end

DEFAULT_DOC_SERVER_PORT = 6060
desc "run doc server at :port (default: #{DEFAULT_DOC_SERVER_PORT})"
task :serve_doc, [:port] do |_, args|
  args.with_defaults(port: DEFAULT_DOC_SERVER_PORT)
  system "godoc -http=:#{args.port}"
end

desc "Verify package by tag"
task :verify, [:tag] do |_, args|
  origin_url = `git config --get remote.origin.url`.strip
  abort "Error! your origin url (#{origin_url}) is not valid!\nShould be something like:\n\n  git@github.com:USER-ORG/REPO.git\n\n" if origin_url.include?('http')

  latest_git_tag = `git describe --tags --abbrev=0`.strip
  args.with_defaults(tag: latest_git_tag)
  
  go_package_name = origin_url.split('@').last.split('.git').first.gsub(':','/')
  stdout, stderr, status = Open3.capture3("go list -m #{go_package_name}@#{args.tag}")
  abort "Error, #{stderr}" unless stderr.empty?
  
  puts "Package verified: #{stdout}"
end

namespace :release do
  desc 'do release check'
  task check: [:repo_clean] do
    system 'go mod tidy'
    Rake::Task['test'].invoke('-v')
    
    puts "-> you are good to go..."
  end

  desc "Publish project with revision: #{AVAILABLE_REVISIONS.join(',')}, default: patch"
  task :publish, [:revision] => [:repo_clean] do |_, args|
    args.with_defaults(revision: 'patch')

    Rake::Task['bump'].invoke(args.revision)

    go_package_name = "#{Rake::Task['get_package_name'].invoke.first.call}"
    current_git_tag = "v#{Rake::Task['current_version'].invoke.first.call}"
    current_branch = `git rev-parse --abbrev-ref HEAD`.strip

    puts ""\
      "-> new version: #{current_git_tag}\n"\
      "-> pushing tag #{current_git_tag} to remote...\n"\
      "-> updating/pushing #{current_branch} branch\n"\

    system %(
      git push origin #{current_git_tag} &&
      git push origin #{current_branch} &&
      go list -m #{go_package_name}@#{current_git_tag}
    )
    puts '-> all complete!'
  end  

end

namespace :docker do
  desc "Build"
  task :build do
    abort "Please set GITHUB_TOKEN env-var" unless ENV['GITHUB_TOKEN']
    system %{
      docker build --build-arg="github_token=${GITHUB_TOKEN}" .
    }
  end
end