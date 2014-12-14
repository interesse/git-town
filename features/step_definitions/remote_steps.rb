Given(/^my repo has an upstream repo$/) do
  clone_repository remote_repository_path, upstream_remote_repository_path, bare: true
  clone_repository upstream_remote_repository_path, upstream_local_repository_path

  at_path upstream_local_repository_path do
    run 'git checkout main'
  end

  run "git remote add upstream #{upstream_remote_repository_path}"
end


Given(/^my remote origin is on (GitHub|Bitbucket) through (HTTPS|SSH)$/) do |domain, protocol|
  run "git remote set-url origin #{git_url domain, protocol}"
end
